package aggregated_block_feed

import (
	"github.com/Gearbox-protocol/third-eye/artifacts/priceFeed"
	"github.com/Gearbox-protocol/third-eye/artifacts/yVault"
	"github.com/Gearbox-protocol/third-eye/artifacts/yearnPriceFeed"
	"github.com/Gearbox-protocol/third-eye/core"
	"github.com/Gearbox-protocol/third-eye/ethclient"
	"github.com/Gearbox-protocol/third-eye/log"
	"github.com/Gearbox-protocol/third-eye/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"sync"
)

type YearnPriceFeed struct {
	*core.SyncAdapter
	contractETH       *yearnPriceFeed.YearnPriceFeed
	YVaultContract    *yVault.YVault
	PriceFeedContract *priceFeed.PriceFeed
	DecimalDivider    *big.Int
	mu                *sync.Mutex
}

func NewYearnPriceFeed(token, oracle string, discoveredAt int64, client ethclient.ClientI, repo core.RepositoryI, version int64) *YearnPriceFeed {
	syncAdapter := &core.SyncAdapter{
		Contract: &core.Contract{
			Address:      oracle,
			DiscoveredAt: discoveredAt,
			FirstLogAt:   discoveredAt,
			ContractName: core.YearnPriceFeed,
			Client:       client,
		},
		Details:  map[string]interface{}{"token": token},
		LastSync: discoveredAt - 1,
		Repo:     repo,
		V: version,
	}
	// add token oracle for db
	// feed is also oracle address for yearn address
	// we don't relie on underlying feed
	repo.AddTokenOracle(token, oracle, oracle, discoveredAt)
	return NewYearnPriceFeedFromAdapter(
		syncAdapter,
	)
}

func NewYearnPriceFeedFromAdapter(adapter *core.SyncAdapter) *YearnPriceFeed {
	yearnPFContract, err := yearnPriceFeed.NewYearnPriceFeed(common.HexToAddress(adapter.Address), adapter.Client)
	if err != nil {
		log.Fatal(err)
	}
	obj := &YearnPriceFeed{
		SyncAdapter: adapter,
		contractETH: yearnPFContract,
		mu:          &sync.Mutex{},
	}
	obj.OnlyQuery = true
	return obj
}

func (mdl *YearnPriceFeed) OnLog(txLog types.Log) {

}
func (mdl *YearnPriceFeed) isNotified() bool {
	mdl.mu.Lock()
	defer mdl.mu.Unlock()
	if mdl.Details == nil || mdl.Details["notified"] == nil {
		return false
	}
	value, ok := mdl.Details["notified"].(bool)
	if !ok {
		log.Fatal("Notified not parsed")
	}
	return value
}

func (mdl *YearnPriceFeed) setNotified(notified bool) {
	mdl.mu.Lock()
	defer mdl.mu.Unlock()
	mdl.Details["notified"] = notified
}

func (mdl *YearnPriceFeed) GetTokenAddr() string {
	mdl.mu.Lock()
	defer mdl.mu.Unlock()
	tokenAddr, ok := mdl.Details["token"].(string)
	if !ok {
		log.Fatal("Failing in asserting to string: %s", mdl.Details["token"])
	}
	return tokenAddr
}

func (mdl *YearnPriceFeed) calculatePriceFeedInternally(blockNum int64) *core.PriceFeed {
	if mdl.YVaultContract == nil || mdl.PriceFeedContract == nil || mdl.DecimalDivider == nil {
		mdl.setContracts(blockNum)
	}
	opts := &bind.CallOpts{
		BlockNumber: big.NewInt(blockNum),
	}

	roundData, err := mdl.PriceFeedContract.LatestRoundData(opts)
	log.CheckFatal(err)
	pricePerShare, err := mdl.YVaultContract.PricePerShare(opts)
	log.CheckFatal(err)

	lowerBound, err := mdl.contractETH.LowerBound(opts)
	log.CheckFatal(err)
	uppwerBound, err := mdl.contractETH.UpperBound(opts)
	log.CheckFatal(err)
	if !(pricePerShare.Cmp(lowerBound) >= 0 && pricePerShare.Cmp(uppwerBound) <= 0) {
		if !mdl.isNotified() {
			mdl.setNotified(true)
			log.Warnf("PricePerShare(%d) is not btw lower limit(%d) and upper limit(%d).", pricePerShare, lowerBound, uppwerBound)
		}
	} else {
		mdl.Details["notified"] = false
	}

	newAnswer := new(big.Int).Quo(
		new(big.Int).Mul(pricePerShare, roundData.Answer),
		mdl.DecimalDivider,
	)
	isPriceInUSD := mdl.GetVersion() > 1
	var decimals int8 = 18 // for eth
	if isPriceInUSD {
		decimals = 8 // for usd
	}
	return &core.PriceFeed{
		RoundId:    roundData.RoundId.Int64(),
		PriceBI: (*core.BigInt)(newAnswer),
		Price:   utils.GetFloat64Decimal(newAnswer, decimals),
		IsPriceInUSD: isPriceInUSD,
	}
}

func (mdl *YearnPriceFeed) setContracts(blockNum int64) {
	opts := &bind.CallOpts{
		BlockNumber: big.NewInt(blockNum),
	}
	// set the price feed contract
	priceFeedAddr, err := mdl.contractETH.PriceFeed(opts)
	log.CheckFatal(err)
	priceFeedContract, err := priceFeed.NewPriceFeed(priceFeedAddr, mdl.Client)
	log.CheckFatal(err)
	mdl.PriceFeedContract = priceFeedContract

	// set the yvault contract
	yVaultAddr, err := mdl.contractETH.YVault(opts)
	log.CheckFatal(err)
	yVaultContract, err := yVault.NewYVault(yVaultAddr, mdl.Client)
	log.CheckFatal(err)
	mdl.YVaultContract = yVaultContract

	// set the decimals
	decimals, err := yVaultContract.Decimals(opts)
	log.CheckFatal(err)
	mdl.DecimalDivider = utils.GetExpInt(int8(decimals))
}
