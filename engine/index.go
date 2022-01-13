package engine

import (
	"context"
	"github.com/Gearbox-protocol/third-eye/config"
	"github.com/Gearbox-protocol/third-eye/core"
	"github.com/Gearbox-protocol/third-eye/ethclient"
	"github.com/Gearbox-protocol/third-eye/log"
	"github.com/Gearbox-protocol/third-eye/models/address_provider"
	"github.com/Gearbox-protocol/third-eye/utils"
	"github.com/ethereum/go-ethereum/common"
	"sync"
	"time"
)

type Engine struct {
	*core.Node
	config  *config.Config
	repo    core.RepositoryI
	debtEng core.DebtEngineI
}

var syncBlockBatchSize = 1000 * core.NoOfBlocksPerMin

func NewEngine(config *config.Config,
	ec *ethclient.Client,
	debtEng core.DebtEngineI,
	repo core.RepositoryI) core.EngineI {
	chaindId, err := ec.ChainID(context.TODO())
	log.CheckFatal(err)
	eng := &Engine{
		debtEng: debtEng,
		config:  config,
		repo:    repo,
		Node: &core.Node{
			Client:  ec,
			ChainId: chaindId.Int64(),
		},
	}
	eng.init()
	return eng
}

func (e *Engine) init() {
	// debt engine initialisation
	e.debtEng.ProcessBackLogs()
}

func (e *Engine) getLastSyncedTill() int64 {
	kit := e.repo.GetKit()
	kit.Details()
	if kit.LenOfLevel(0) == 0 {
		addr := common.HexToAddress(e.config.AddressProviderAddress).Hex()
		obj := address_provider.NewAddressProvider(addr, e.Client, e.repo)
		e.repo.AddSyncAdapter(obj)
		return obj.GetLastSync()
	} else {
		// it will allow syncing from scratch of least synced adapter in batches
		// NOTE: while syncing from scratch for some adapter disable the debt engine
		// as it might happen that some of the components for calculating debt are missing
		return e.repo.LoadLastAdapterSync()
	}
}

func (e *Engine) SyncHandler() {
	latestBlockNum := e.GetLatestBlockNumber()
	lastSyncedTill := e.getLastSyncedTill()
	// only do batch sync if latestblock is far from currently synced block
	if lastSyncedTill+syncBlockBatchSize <= latestBlockNum {
		syncedTill := e.syncLoop(lastSyncedTill, latestBlockNum)
		log.Infof("Synced till %d sleeping for 5 mins", syncedTill)
	}
	for {
		latestBlockNum = e.GetLatestBlockNumber()
		e.sync(latestBlockNum)
		log.Infof("Synced till %d sleeping for 5 mins", latestBlockNum)
		time.Sleep(5 * time.Minute) // on kovan 5 blocks in 1 min , sleep for 5 mins
	}
}

func (e *Engine) syncLoop(syncedTill, latestBlockNum int64) int64 {
	loopStartBlock := syncedTill
	syncTill := syncedTill + syncBlockBatchSize
	loopStartTime := time.Now()
	for syncTill <= latestBlockNum {
		roundStartTime := time.Now()
		e.sync(syncTill)
		syncedTill = syncTill
		roundSyncDur := (time.Now().Sub(roundStartTime).Minutes())
		timePerBlock := time.Now().Sub(loopStartTime).Minutes() / float64(syncedTill-loopStartBlock)
		remainingTime := (timePerBlock * float64(latestBlockNum-syncedTill)) / (60)
		log.Infof("Synced till %d in %f mins. Remaining time %f hrs ", syncedTill, roundSyncDur, remainingTime)
		// new sync target
		syncTill += syncBlockBatchSize
	}
	return syncedTill
}

func (e *Engine) sync(syncTill int64) {
	kit := e.repo.GetKit()
	log.Info("Sync till", syncTill)
	for lvlIndex := 0; lvlIndex < kit.Len(); lvlIndex++ {
		wg := &sync.WaitGroup{}
		for kit.Next(lvlIndex) {
			adapter := kit.Get(lvlIndex)
			// if utils.Contains([]string{core.AccountFactory, core.YearnPriceFeed, core.ChainlinkPriceFeed}, adapter.GetName()) {
			// 	continue
			// }
			if !adapter.IsDisabled() {
				wg.Add(1)
				if adapter.OnlyQueryAllowed() {
					adapter.Query(syncTill, wg)
				} else {
					e.SyncModel(adapter, syncTill, wg)
				}
			}
		}
		kit.Reset(lvlIndex)
		wg.Wait()
	}
	e.FlushAndDebt(syncTill)
	e.repo.CalCurrentTreasuryValue(syncTill)
}

func (e *Engine) SyncModel(mdl core.SyncAdapterI, syncTill int64, wg *sync.WaitGroup) {
	defer wg.Done()
	syncFrom := mdl.GetLastSync() + 1
	if syncFrom > syncTill {
		return
	}

	log.Infof("Sync %s(%s) from %d to %d", mdl.GetName(), mdl.GetAddress(), syncFrom, syncTill)
	logs, err := e.GetLogs(syncFrom, syncTill, []common.Address{common.HexToAddress(mdl.GetAddress())}, [][]common.Hash{})
	if err != nil {
		log.Fatal(err)
	}
	for _, txLog := range logs {
		blockNum := int64(txLog.BlockNumber)
		if mdl.GetBlockToDisableOn() < blockNum {
			break
		}
		e.repo.SetBlock(blockNum)
		mdl.OnLog(txLog)
	}
	// after sync
	mdl.AfterSyncHook(utils.Min(mdl.GetBlockToDisableOn(), syncTill))
}

func (e *Engine) FlushAndDebt(to int64) {
	e.repo.Flush()
	e.debtEng.CalculateDebtAndClear(to)
}
