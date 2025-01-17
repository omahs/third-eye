package debts

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/Gearbox-protocol/sdk-go/artifacts/dataCompressor/dataCompressorv2"
	"github.com/Gearbox-protocol/sdk-go/artifacts/multicall"
	"github.com/Gearbox-protocol/sdk-go/artifacts/yearnPriceFeed"
	"github.com/Gearbox-protocol/sdk-go/calc"
	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/Gearbox-protocol/sdk-go/core/schemas"
	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/Gearbox-protocol/third-eye/ds"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (eng *DebtEngine) SaveProfile(profile string) {
	err := eng.db.Create(&schemas.ProfileTable{Profile: profile}).Error
	if err != nil {
		log.Fatal(err)
	}
}

func (eng *DebtEngine) updateLocalState(blockNum int64, block *schemas.Block) (poolsUpdated []string, tokensUpdated map[string]bool, sessionsUpdated map[string]bool) {
	//////////////////////////////////////
	// Load data
	//////////////////////////////////////
	// L1 params
	// L2 poolinterestdata
	// L3 credit session snapshots
	// L4 lt
	// L5 pricefeeds
	// L5 rebaseToken details for stETH
	//
	// L1:update params
	for _, params := range block.GetParams() {
		eng.addLastParameters(params)
	}
	// L6: rebaseToken
	for _, params := range block.RebaseDetailsForDB {
		eng.lastRebaseDetails = params
	}

	///////////////////////////////////
	// calc debt conditions
	///////////////////////////////////
	// C1 If pool cumIndex changes
	// C2 If new css is present, meaning balance or new token added
	// C3 If account related token(either have balance or underlying token) has some data change, lt or price

	// update pool borrow rate and cumulative index
	// C1: UPDATE BASED ON POOL CUM INDEX
	for _, ps := range block.GetPoolStats() {
		poolsUpdated = append(poolsUpdated, ps.Address)
		// L2
		eng.AddPoolLastInterestData(&schemas.PoolInterestData{
			Address:              ps.Address,
			BlockNum:             ps.BlockNum,
			CumulativeIndexRAY:   ps.CumulativeIndexRAY,
			AvailableLiquidityBI: ps.AvailableLiquidityBI,
			BorrowAPYBI:          ps.BorrowAPYBI,
			Timestamp:            block.Timestamp,
		})
	}

	// update balance of last credit session snapshot and create credit session snapshots
	// C2: UPDATE BASED ON CSS
	sessionsUpdated = make(map[string]bool)
	for _, css := range block.GetCSS() {
		// L3
		eng.AddLastCSS(css)
		sessionsUpdated[css.SessionId] = true
	}

	// C3: UPDATE BASED ON TOKEN
	// C3.a: updated threshold
	tokensUpdated = make(map[string]bool)
	for _, allowedToken := range block.GetAllowedTokens() {
		// L4
		eng.AddAllowedTokenThreshold(allowedToken)
		tokensUpdated[allowedToken.Token] = true
	}

	// C3.b: updated price
	for _, pf := range block.GetPriceFeeds() {
		// L5
		eng.AddTokenLastPrice(pf)
		tokensUpdated[pf.Token] = true
	}
	// updates complete
	return poolsUpdated, tokensUpdated, sessionsUpdated
}

func (eng *DebtEngine) CalculateDebt() {
	blocks := eng.repo.GetBlocks()
	sessions := eng.repo.GetSessions()
	// sort block numbers
	noOfBlock := len(blocks)
	blockNums := make([]int64, 0, noOfBlock)
	for blockNum := range blocks {
		blockNums = append(blockNums, blockNum)
	}
	sort.Slice(blockNums, func(i, j int) bool { return blockNums[i] < blockNums[j] })
	//
	for _, blockNum := range blockNums {
		block := blocks[blockNum]
		poolsUpdated, tokensUpdated, sessionsUpdated := eng.updateLocalState(blockNum, block)
		// get pool cumulative interest rate
		cmToPoolDetails := eng.GetCumulativeIndexAndDecimalForCMs(blockNum, block.Timestamp)

		//
		var caTotalValueInUSD float64 = 0
		// check if session's debt needs to be recalculated
		for _, session := range sessions {
			if (session.ClosedAt != 0 && session.ClosedAt <= blockNum) || session.Since > blockNum {
				continue
			}
			// #C1
			if utils.Contains(poolsUpdated, cmToPoolDetails[session.CreditManager].PoolAddr) {
				sessionsUpdated[session.ID] = true
			}
			// #C3
			sessionSnapshot := eng.lastCSS[session.ID]
			caTotalValueInUSD += utils.GetFloat64Decimal(
				eng.GetAmountInUSD(
					cmToPoolDetails[session.CreditManager].Token,
					sessionSnapshot.TotalValueBI.Convert(), session.Version,
				), 8)

			for token, tokenBalance := range *sessionSnapshot.Balances {
				if tokenBalance.IsEnabled && tokenBalance.HasBalanceMoreThanOne() {
					if tokensUpdated[token] {
						sessionsUpdated[session.ID] = true
					}
				}
			}
			// #C3
			underlyingToken := cmToPoolDetails[session.CreditManager].Token
			if tokensUpdated[underlyingToken] {
				sessionsUpdated[session.ID] = true
			}
		}
		//
		// calculate each session debt
		for sessionId := range sessionsUpdated {
			session := sessions[sessionId]
			if (session.ClosedAt != 0 && session.ClosedAt <= blockNum) || session.Since > blockNum {
				continue
			}
			cmAddr := session.CreditManager
			// pool cum index is when the pool is not registered
			if cmToPoolDetails[cmAddr] != nil {
				eng.SessionDebtHandler(blockNum, session, cmToPoolDetails[cmAddr])
				// send notification when account is liquidated
			} else {
				log.Fatalf("CM(%s):pool is missing stats at %d, so cumulative index of pool is unknown", cmAddr, blockNum)
			}
		}
		//
		eng.createTvlSnapshots(blockNum, caTotalValueInUSD)
		if len(sessionsUpdated) > 0 {
			log.Debugf("Calculated %d debts for block %d", len(sessionsUpdated), blockNum)
		}
	}
}

func (eng *DebtEngine) createTvlSnapshots(blockNum int64, caTotalValueInUSD float64) {
	if eng.lastTvlSnapshot != nil && blockNum-eng.lastTvlSnapshot.BlockNum < core.NoOfBlocksPerHr { // tvl snapshot every hour
		return
	}
	var totalAvailableLiquidityInUSD float64 = 0
	for _, entry := range eng.poolLastInterestData {
		adapter := eng.repo.GetAdapter(entry.Address)
		state := adapter.GetUnderlyingState()
		if state == nil {
			log.Fatal("State for pool not found for address: ", entry.Address)
		}
		//
		underlyingToken := state.(*schemas.PoolState).UnderlyingToken
		//
		version := core.NewVersion(1)
		if eng.tokenLastPriceV2[underlyingToken] != nil {
			version = core.NewVersion(2)
		}
		//
		totalAvailableLiquidityInUSD += utils.GetFloat64Decimal(
			eng.GetAmountInUSD(
				underlyingToken,
				entry.AvailableLiquidityBI.Convert(), version), 8)
	}
	// save as last tvl snapshot and add to db
	tvls := &schemas.TvlSnapshots{
		BlockNum:           blockNum,
		AvailableLiquidity: totalAvailableLiquidityInUSD,
		CATotalValue:       caTotalValueInUSD,
	}
	eng.tvlSnapshots = append(eng.tvlSnapshots, tvls)
	eng.lastTvlSnapshot = tvls
}

func (eng *DebtEngine) ifAccountLiquidated(sessionId, cmAddr string, closedAt int64, status int) {
	sessionSnapshot := eng.lastCSS[sessionId]
	if schemas.IsStatusLiquidated(status) {
		account := eng.liquidableBlockTracker[sessionId]
		var liquidableSinceBlockNum int64
		if account != nil {
			liquidableSinceBlockNum = account.BlockNum
		} else {
			log.Warnf("Session(%s) liquidated at block:%d, but liquidable since block not stored", sessionId, closedAt)
		}
		urls := core.NetworkUIUrl(core.GetChainId(eng.client))
		eng.repo.RecentMsgf(log.RiskHeader{
			BlockNumber: closedAt - 1,
			EventCode:   "AMQP",
		}, `Liquidation Alert:
		CreditManager: %s/address/%s
		Tx: %s/tx/%s
		Borrower: %s
		LiquidatedAt: %d 
		Liquidable since: %d
		web: %s/accounts/history/%s`,
			urls.ExplorerUrl, cmAddr,
			urls.ExplorerUrl, eng.GetLiquidationTx(sessionId),
			sessionSnapshot.Borrower,
			closedAt, liquidableSinceBlockNum,
			urls.ChartUrl, sessionSnapshot.SessionId)
	} else if status != schemas.Active {
		// comment deletion form map
		// this way if the account is liquidable and liquidated in the same debt calculation cycle
		// we will be able to store in db at which block it became liquidable
		// delete(eng.liquidableBlockTracker, sessionId)
	}
}

func (eng *DebtEngine) GetCumulativeIndexAndDecimalForCMs(blockNum int64, ts uint64) map[string]*ds.CumIndexAndUToken {
	// this is assuming that credit managers are not disabled
	cmAddrs := eng.repo.GetAdapterAddressByName(ds.CreditManager)
	cmToCumIndex := make(map[string]*ds.CumIndexAndUToken)
	for _, cmAddr := range cmAddrs {
		cmState := eng.repo.GetCMState(cmAddr)
		// if block_num is before cm exists cmState will be nil
		if cmState == nil {
			continue
		}
		poolAddr := cmState.PoolAddress
		poolInterestData := eng.poolLastInterestData[poolAddr]
		var cumIndexNormalized *big.Int
		if poolInterestData != nil {
			tsDiff := new(big.Int).SetInt64(int64(ts - poolInterestData.Timestamp))
			newInterest := new(big.Int).Mul(poolInterestData.BorrowAPYBI.Convert(), tsDiff)
			newInterestPerSec := new(big.Int).Quo(newInterest, big.NewInt(3600*365*24))
			predicate := new(big.Int).Add(newInterestPerSec, utils.GetExpInt(27))
			cumIndex := new(big.Int).Mul(poolInterestData.CumulativeIndexRAY.Convert(), predicate)
			cumIndexNormalized = utils.GetInt64(cumIndex, 27)
		}
		// set fields
		tokenAddr := cmState.UnderlyingToken
		token := eng.repo.GetToken(tokenAddr)
		cmToCumIndex[cmAddr] = &ds.CumIndexAndUToken{
			CumulativeIndex: cumIndexNormalized,
			Token:           tokenAddr,
			Symbol:          token.Symbol,
			Decimals:        token.Decimals,
			PoolAddr:        poolAddr,
		}
	}
	return cmToCumIndex
}

func (eng *DebtEngine) getTokenPriceFeed(token string, version core.VersionType) *schemas.PriceFeed {
	if version.IsGBv1() {
		return eng.tokenLastPrice[token]
	} else { // v2 and above
		return eng.tokenLastPriceV2[token]
	}
}

func (eng *DebtEngine) SessionDebtHandler(blockNum int64, session *schemas.CreditSession, cumIndexAndUToken *ds.CumIndexAndUToken) {
	sessionId := session.ID
	sessionSnapshot := eng.lastCSS[sessionId]
	cmAddr := session.CreditManager
	debt, profile := eng.CalculateSessionDebt(blockNum, session, cumIndexAndUToken)
	// if profile is not null
	// yearn price feed might be stale as a result difference btw dc and calculated values
	// solution: fetch price again for all stale yearn feeds
	if profile != nil {
		yearnFeeds := eng.repo.GetYearnFeedAddrs()
		for tokenAddr, details := range *sessionSnapshot.Balances {
			if details.IsEnabled && details.HasBalanceMoreThanOne() {
				lastPriceEvent := eng.getTokenPriceFeed(tokenAddr, session.Version)
				//
				if tokenAddr != eng.repo.GetWETHAddr() && lastPriceEvent.BlockNumber != blockNum {
					feed := lastPriceEvent.Feed
					if utils.Contains(yearnFeeds, feed) {
						eng.requestPriceFeed(blockNum, feed, tokenAddr, lastPriceEvent.IsPriceInUSD)
					}
				}
			}
		}
		debt, profile = eng.CalculateSessionDebt(blockNum, session, cumIndexAndUToken)
		if profile != nil {
			log.Fatalf("Debt fields different from data compressor fields: %s", profile)
		}
	}
	// check if data compressor and calculated values match
	eng.liquidationCheck(debt, cmAddr, sessionSnapshot.Borrower, cumIndexAndUToken)
	if session.ClosedAt == blockNum+1 {
		eng.ifAccountLiquidated(sessionId, cmAddr, session.ClosedAt, session.Status)
		eng.addCurrentDebt(debt, cumIndexAndUToken.Decimals)
	}
	eng.AddDebt(debt, sessionSnapshot.BlockNum == blockNum)
}

func (eng *DebtEngine) CalculateSessionDebt(blockNum int64, session *schemas.CreditSession, cumIndexAndUToken *ds.CumIndexAndUToken) (*schemas.Debt, *ds.DebtProfile) {
	sessionId := session.ID
	sessionSnapshot := eng.lastCSS[sessionId]
	cmAddr := session.CreditManager
	//
	// calculating account fields
	calculator := calc.Calculator{Store: storeForCalc{inner: eng}}
	calHF, calDebt, calTotalValue, calThresholdValue, _calBorowedWithInterst := calculator.CalcAccountFields(
		session.Version,
		blockNum,
		sessionDetailsForCalc{
			CreditSessionSnapshot: sessionSnapshot,
			CM:                    session.CreditManager,
			rebaseDetails:         eng.lastRebaseDetails,
			stETH:                 eng.repo.GetTokenFromSdk("stETH"),
		},
		cumIndexAndUToken.CumulativeIndex,
		cumIndexAndUToken.Token,
		eng.lastParameters[session.CreditManager].FeeInterest,
	)
	// if session.ID == "0x57ed1ED84461bb2079f8575d06A6feC07F0a13B1_16159748_285" {
	// 	log.Fatal(blockNum, calHF, calDebt, calTotalValue, calThresholdValue, _calBorowedWithInterst)
	// }
	//
	// the value of credit account is in terms of underlying asset
	// set debt fields
	debt := &schemas.Debt{
		CommonDebtFields: schemas.CommonDebtFields{
			BlockNumber:     blockNum,
			CalHealthFactor: (*core.BigInt)(calHF),
			CalTotalValueBI: (*core.BigInt)(calTotalValue),
			// it has fees too for v2
			CalDebtBI: (*core.BigInt)(calDebt),
			// used for calculating the close amount and for comparing with data compressor results as v2 doesn't have borrowedAmountWithInterest
			CalBorrowedWithInterestBI: (*core.BigInt)(_calBorowedWithInterst),
			CalThresholdValueBI:       (*core.BigInt)(calThresholdValue),
			CollateralInUnderlying:    sessionSnapshot.CollateralInUnderlying,
		},
		SessionId: sessionId,
	}
	var notMatched bool
	profile := ds.DebtProfile{CreditSessionSnapshot: sessionSnapshot, Tokens: map[string]ds.TokenDetails{}}
	for tokenAddr, details := range *sessionSnapshot.Balances {
		if details.IsAllowed {
			profile.Tokens[tokenAddr] = ds.TokenDetails{
				Price:             eng.GetTokenLastPrice(tokenAddr, session.Version),
				Decimals:          eng.repo.GetToken(tokenAddr).Decimals,
				TokenLiqThreshold: eng.allowedTokensThreshold[session.CreditManager][tokenAddr],
			}
		}
	}

	// use data compressor if debt check is enabled
	if eng.config.DebtDCMatching {
		data := eng.SessionDataFromDC(blockNum, cmAddr, sessionSnapshot.Borrower)
		utils.ToJson(data)
		// set debt data fetched from dc
		// if healthfactor on diff
		if !CompareBalance(debt.CalTotalValueBI, (*core.BigInt)(data.TotalValue), cumIndexAndUToken) ||
			!CompareBalance(debt.CalBorrowedWithInterestBI, (*core.BigInt)(data.BorrowedAmountPlusInterest), cumIndexAndUToken) ||
			core.ValueDifferSideOf10000(debt.CalHealthFactor, (*core.BigInt)(data.HealthFactor)) {
			profile.DCData = &data
			notMatched = true
		}
		// even if data compressor matching is disabled check the calc values  with session data at block where last credit snapshot was taken
	} else if sessionSnapshot.BlockNum == blockNum {
		if !CompareBalance(debt.CalTotalValueBI, sessionSnapshot.TotalValueBI, cumIndexAndUToken) ||
			// hf value calculated are on different side of 1
			core.ValueDifferSideOf10000(debt.CalHealthFactor, sessionSnapshot.HealthFactor) ||
			// if healhFactor diff by 4 %
			core.DiffMoreThanFraction(debt.CalHealthFactor, sessionSnapshot.HealthFactor, big.NewFloat(0.04)) {
			// log.Info(debt.CalHealthFactor, sessionSnapshot.HealthFactor, blockNum)
			// log.Info(debt.CalTotalValueBI, sessionSnapshot.TotalValueBI, blockNum)
			notMatched = true
		}
	}
	eng.calAmountToPoolAndProfit(debt, session, cumIndexAndUToken)
	eng.farmingCalc.addFarmingVal(debt, session, eng.lastCSS[session.ID], storeForCalc{inner: eng})
	if notMatched {
		profile.CumIndexAndUToken = cumIndexAndUToken
		profile.Debt = debt
		profile.CreditSessionSnapshot = sessionSnapshot
		profile.UnderlyingDecimals = cumIndexAndUToken.Decimals
		return debt, &profile
	}
	return debt, nil
}

// these values are calculated for the borrower not the liquidation, so the calrepayamount takes isliquidated = false
// only for block that the account is liquidated we use isLiquidated set to true so that we can calculate the true amountToPool
//
// amountToPool
// for v1,v2 amountToPool is calculated by calc.CalCloseAmount
//
// remainingfunds => is used for profit calculation, assets that user gets back.
// - v1 close or liquidated -- remainingfunds is taken from event
// - v1 repay + open -- remainingfunds is manually calculated with help of calCloseAmount in sdk-go
// - v2 for closeCreditAccount, remainingFunds is calculated from the account transfers
// - v2 for liquidateCreditAccount, remainingFunds is taken from event
// - v2 for openedAccounts, totalValue - amountToPool
//
// repayAmount => transfer from owner to account needed to close the account
// v1 - repayAmount = amountToPool, except the blockNum at which account is liquidated
//   - for liquidated account repay amount is amountToPool+ calcRemainginFunds https://github.com/Gearbox-protocol/gearbox-contracts/blob/master/contracts/credit/CreditManager.sol#L999
//   - NIT, closeAmount doesn't need repayAmount as all assets are converted to underlying token
//     so repayAmount is zero, https://github.com/Gearbox-protocol/gearbox-contracts/blob/master/contracts/credit/CreditManager.sol#L448-L465
//
// v2 - close repayAmount is transferred from borrower to account as underlying token
// v2 - for liquidation, repayAmount is zero.
// v2 - opened accounts
//
//	the account might be having some underlying token balance so repayAMount = amountToPool - underlyingToken balance
func (eng *DebtEngine) calAmountToPoolAndProfit(debt *schemas.Debt, session *schemas.CreditSession, cumIndexAndUToken *ds.CumIndexAndUToken) {
	var amountToPool, calRemainingFunds *big.Int
	sessionSnapshot := eng.lastCSS[session.ID]
	//
	status := schemas.Active
	if schemas.IsStatusLiquidated(session.Status) && session.ClosedAt == debt.BlockNumber+1 { // calc based on status if the block has liquidation
		status = session.Status
	}
	// amount to pool
	amountToPool, calRemainingFunds, _, _ = calc.CalCloseAmount(eng.lastParameters[session.CreditManager],
		session.Version, debt.CalTotalValueBI.Convert(), status,
		debt.CalBorrowedWithInterestBI.Convert(),
		sessionSnapshot.BorrowedAmountBI.Convert())

	// calculate profit
	debt.AmountToPoolBI = (*core.BigInt)(amountToPool)
	var remainingFunds *big.Int
	if session.Version == 2 {
		repayAmount := new(big.Int)
		// while close account on v2 we calculate remainingFunds from all the token transfer from the user
		if session.Status == schemas.Closed && session.ClosedAt == debt.BlockNumber+1 {
			prices := core.JsonFloatMap{}
			for token, transferAmt := range *session.CloseTransfers {
				tokenPrice := eng.GetTokenLastPrice(token, session.Version)
				price := utils.GetFloat64Decimal(tokenPrice, 8)
				prices[token] = price
				if transferAmt < 0 {
					// assuming there is only one transfer from borrower to account
					// this transfer will be in underlyingtoken. execute_parser.go:246 and
					// https://github.com/Gearbox-protocol/core-v2/blob/main/contracts/credit/CreditManager.sol#L359-L363
					amt := new(big.Float).Mul(big.NewFloat(transferAmt*-1), utils.GetExpFloat(eng.repo.GetToken(token).Decimals))
					repayAmount, _ = amt.Int(nil)
				}
			}
			// remainingFunds calculation
			// set price for underlying token
			prices[cumIndexAndUToken.Token] = utils.GetFloat64Decimal(
				eng.GetTokenLastPrice(cumIndexAndUToken.Token, session.Version), 8)
			remainingFunds = session.CloseTransfers.ValueInUnderlying(cumIndexAndUToken.Token, cumIndexAndUToken.Decimals, prices)
		} else if session.ClosedAt == debt.BlockNumber+1 && schemas.IsStatusLiquidated(session.Status) {
			remainingFunds = session.RemainingFunds.Convert()
			repayAmount = new(big.Int)
		} else {
			// repayamount
			// for account not closed or liquidated yet
			// get underlying balance
			underlying := (*eng.lastCSS[session.ID].Balances)[cumIndexAndUToken.Token]
			underlyingBalance := new(big.Int)
			if underlying.BI != nil {
				underlyingBalance = underlying.BI.Convert()
			}
			if new(big.Int).Sub(underlyingBalance, new(big.Int).Add(amountToPool, calRemainingFunds)).Cmp(big.NewInt(1)) > 0 {
				repayAmount = new(big.Int)
			} else {
				repayAmount = new(big.Int).Sub(new(big.Int).Add(amountToPool, calRemainingFunds), underlyingBalance)
			}

			// remainingfunds
			remainingFunds = new(big.Int).Sub(debt.CalTotalValueBI.Convert(), amountToPool)
		}
		debt.RepayAmountBI = (*core.BigInt)(repayAmount)
	} else {
		if session.ClosedAt == debt.BlockNumber+1 && (session.Status == schemas.Closed || session.Status == schemas.Liquidated) {
			remainingFunds = (*big.Int)(session.RemainingFunds)
		} else {
			remainingFunds = calRemainingFunds
		}
		if session.Status == schemas.Liquidated && session.ClosedAt == debt.BlockNumber+1 {
			debt.RepayAmountBI = (*core.BigInt)(new(big.Int).Add(amountToPool, remainingFunds))
		} else if session.Status == schemas.Closed && session.ClosedAt == debt.BlockNumber+1 {
			debt.RepayAmountBI = (*core.BigInt)(new(big.Int))
		} else {
			// https://github.com/Gearbox-protocol/gearbox-contracts/blob/master/contracts/credit/CreditManager.sol#L487-L490
			debt.RepayAmountBI = (*core.BigInt)(amountToPool)
		}
	}

	remainingFundsInUSD := eng.GetAmountInUSD(cumIndexAndUToken.Token, remainingFunds, session.Version)
	debt.ProfitInUnderlying = utils.GetFloat64Decimal(remainingFunds, cumIndexAndUToken.Decimals) - debt.CollateralInUnderlying
	debt.CollateralInUnderlying = sessionSnapshot.CollateralInUnderlying
	// fields in USD
	debt.CollateralInUSD = sessionSnapshot.CollateralInUSD
	debt.ProfitInUSD = utils.GetFloat64Decimal(remainingFundsInUSD, 8) - sessionSnapshot.CollateralInUSD
	debt.TotalValueInUSD = utils.GetFloat64Decimal(
		eng.GetAmountInUSD(cumIndexAndUToken.Token, debt.CalTotalValueBI.Convert(), session.Version), 8)
}

// helper methods
func (eng *DebtEngine) GetAmountInUSD(tokenAddr string, amount *big.Int, version core.VersionType) *big.Int {
	usdcAddr := eng.repo.GetUSDCAddr()
	tokenPrice := eng.GetTokenLastPrice(tokenAddr, version)
	tokenDecimals := eng.repo.GetToken(tokenAddr).Decimals
	if version == 2 {
		return utils.GetInt64(new(big.Int).Mul(tokenPrice, amount), tokenDecimals)
	}
	usdcPrice := eng.GetTokenLastPrice(usdcAddr, version)
	usdcDecimals := eng.repo.GetToken(usdcAddr).Decimals

	value := new(big.Int).Mul(amount, tokenPrice)
	value = utils.GetInt64(value, tokenDecimals-usdcDecimals)
	value = new(big.Int).Quo(value, usdcPrice)
	return new(big.Int).Mul(value, big.NewInt(100))
}

func (eng *DebtEngine) GetTokenLastPrice(addr string, version core.VersionType, dontFail ...bool) *big.Int {
	switch version {
	case 1:
		if eng.tokenLastPrice[addr] != nil {
			return eng.tokenLastPrice[addr].PriceBI.Convert()
		} else if eng.repo.GetWETHAddr() == addr {
			return core.WETHPrice
		}
	case 2:
		if eng.tokenLastPriceV2[addr] != nil {
			return eng.tokenLastPriceV2[addr].PriceBI.Convert()
		}
	}
	if len(dontFail) > 0 && dontFail[0] {
		return nil
	}
	log.Fatal(fmt.Sprintf("Price not found for %s version: %d", addr, version))
	return nil
}

func (eng *DebtEngine) SessionDataFromDC(blockNum int64, cmAddr, borrower string) dataCompressorv2.CreditAccountData {
	call, resultFn, err := eng.repo.GetDCWrapper().GetCreditAccountData(blockNum,
		common.HexToAddress(cmAddr),
		common.HexToAddress(borrower),
	)
	if err != nil {
		log.Fatalf("Prepaing failed. cm:%s borrower:%s blocknum:%d err:%s", cmAddr, borrower, blockNum, err)
	}
	results := core.MakeMultiCall(eng.client, blockNum, false, []multicall.Multicall2Call{call})
	if !results[0].Success {
		log.Fatalf("Failed multicall. cm:%s borrower:%s blocknum:%d ", cmAddr, borrower, blockNum)
	}
	data, err := resultFn(results[0].ReturnData)
	if err != nil {
		log.Fatalf("cm:%s borrower:%s blocknum:%d err:%s", cmAddr, borrower, blockNum, err)
	}
	return data
}

func (eng *DebtEngine) requestPriceFeed(blockNum int64, feed, token string, isPriceInUSD bool) {
	// PFFIX
	yearnPFContract, err := yearnPriceFeed.NewYearnPriceFeed(common.HexToAddress(feed), eng.client)
	log.CheckFatal(err)
	opts := &bind.CallOpts{
		BlockNumber: big.NewInt(blockNum),
	}
	roundData, err := yearnPFContract.LatestRoundData(opts)
	if err != nil {
		log.Fatal(err)
	}
	var decimals int8 = 18 // for eth
	if isPriceInUSD {
		decimals = 8 // for usd
	}
	eng.AddTokenLastPrice(&schemas.PriceFeed{
		BlockNumber:  blockNum,
		Token:        token,
		Feed:         feed,
		RoundId:      roundData.RoundId.Int64(),
		PriceBI:      (*core.BigInt)(roundData.Answer),
		Price:        utils.GetFloat64Decimal(roundData.Answer, decimals),
		IsPriceInUSD: isPriceInUSD,
	})
}
