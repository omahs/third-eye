package credit_manager

import (
	"math/big"

	"github.com/Gearbox-protocol/sdk-go/artifacts/multicall"
	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/Gearbox-protocol/sdk-go/core/schemas"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/Gearbox-protocol/third-eye/ds"
	"github.com/ethereum/go-ethereum/core/types"
)

// v1 method
func (mdl *CreditManager) processExecuteEvents() {
	if len(mdl.executeParams) > 0 {
		mdl.saveExecuteEvents(mdl.executeParams)
		mdl.executeParams = []ds.ExecuteParams{}
	}
}

// x,i event
// fetch events x,i+1 to x,i+1000
//
// works for newBlockNum > mdl.lastEventBlock
func (mdl *CreditManager) OnBlockChange(lastBlockNum int64) (calls []multicall.Multicall2Call, processFns []func(multicall.Multicall2Result)) {
	// datacompressor works for cm address only after the address is registered with contractregister
	// i.e. discoveredAt
	if mdl.lastEventBlock != 0 && mdl.lastEventBlock == lastBlockNum && lastBlockNum >= mdl.DiscoveredAt {
		//// ON NEW TXHASH
		mdl.onTxHash("")
		// ON NEW BLOCK
		data := mdl.Repo.GetAccountManager().CheckTokenTransfer(mdl.GetAddress(), lastBlockNum, lastBlockNum+1)
		mdl.processDirectTransfersOnBlock(lastBlockNum, data[lastBlockNum])
		calls, processFns = mdl.FetchFromDCForChangedSessions(lastBlockNum)
		call, processFn := mdl.getCMCallAndProcessFn(lastBlockNum)
		if processFn != nil {
			calls = append(calls, call)
			processFns = append(processFns, processFn)
		}
		mdl.lastEventBlock = 0
	}
	return
}

func bytesToUInt16(data []byte) uint16 {
	return uint16(new(big.Int).SetBytes(data).Int64())
}

func (mdl *CreditManager) OnLog(txLog types.Log) {
	// creditConfigurator events for test
	// CreditFacadeUpgraded is emitted when creditconfigurator is initialized, so we will receive it on init
	// although we have already set creditfacadeUpgra
	if mdl.GetDetailsByKey("configurator") == txLog.Address.Hex() {
		switch txLog.Topics[0] {
		case core.Topic("CreditFacadeUpgraded(address)"):
			facade := utils.ChecksumAddr(txLog.Topics[1].Hex())
			mdl.setCreditFacadeSyncer(facade)
		case core.Topic("FeesUpdated(uint16,uint16,uint16,uint16,uint16)"):
			mdl.setParams(&schemas.Parameters{
				BlockNum:                   int64(txLog.BlockNumber),
				CreditManager:              mdl.Address,
				FeeInterest:                bytesToUInt16(txLog.Data[:32]),
				FeeLiquidation:             bytesToUInt16(txLog.Data[32:64]),
				LiquidationDiscount:        10000 - bytesToUInt16(txLog.Data[64:96]), // 10000- liqPremium
				FeeLiquidationExpired:      bytesToUInt16(txLog.Data[96:128]),
				LiquidationDiscountExpired: 10000 - bytesToUInt16(txLog.Data[128:160]), // 10000- liqPremiumExpired
			})
		}
		return
	}

	// if facade or cm , not configurator
	mdl.onTxHash(txLog.TxHash.Hex())
	mdl.lastEventBlock = int64(txLog.BlockNumber)
	//
	mdl.Repo.GetAccountManager().DeleteTxHash(int64(txLog.BlockNumber), txLog.TxHash.Hex())
	if mdl.GetVersion().IsGBv1() {
		mdl.checkLogV1(txLog)
	} else {
		mdl.checkLogV2(txLog)
	}
}

// handles for v2(for multicalls) and v1 (for executeorder)
func (mdl *CreditManager) onTxHash(newTxHash string) {
	// on txHash
	if mdl.LastTxHash != "" && mdl.LastTxHash != newTxHash {
		if mdl.GetVersion().IsGBv1() {
			// storing execute order in a single tx and processing them in a single go on next tx
			// for credit session stats
			//
			// execute events are matched with tenderly response to get transfers for each events
			mdl.processExecuteEvents()
		} else {
			mdl.onNewTxHashV2()
		}
	}
	mdl.LastTxHash = newTxHash
}
