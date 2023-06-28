package wrappers

import (
	"math"

	"github.com/Gearbox-protocol/sdk-go/artifacts/multicall"
	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/Gearbox-protocol/third-eye/ds"
	"github.com/Gearbox-protocol/third-eye/models/pool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type OrderedMap struct {
	m    map[string]ds.SyncAdapterI
	allM map[string]ds.SyncAdapterI
	a    []ds.SyncAdapterI
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		m:    make(map[string]ds.SyncAdapterI), // adapter by its actual addr, like creditmanager uses cf , cc but it can be fetched only with creditmanager addr from outside
		allM: make(map[string]ds.SyncAdapterI),
		a:    make([]ds.SyncAdapterI, 0),
	}
}

func (x OrderedMap) Get(addr string) ds.SyncAdapterI {
	return x.m[addr]
}
func (x OrderedMap) GetFromLogAddr(name string) ds.SyncAdapterI {
	return x.allM[name]
}
func (x *OrderedMap) Add(addr string, allAddrsForAdapter []common.Address, val ds.SyncAdapterI) {
	// for
	if x.m[addr] == nil {
		x.a = append(x.a, val)
	}
	for _, addr := range allAddrsForAdapter {
		x.allM[addr.String()] = val
	}
	x.m[addr] = val
}

func (x OrderedMap) GetAll() []ds.SyncAdapterI {
	return x.a
}

// we are creating sync wrappers to wrap , chainlink, creditfilter, credit manager and pools to reduce the number of rpc calls
// only for HasOnLog = true
type SyncWrapper struct {
	Adapters       OrderedMap
	ViaDataProcess int
	name           string
	lastSync       int64
	Client         core.ClientI
	WillSyncTill   int64
}

func NewSyncWrapper(name string, client core.ClientI) *SyncWrapper {
	return &SyncWrapper{
		Adapters:       NewOrderedMap(),
		ViaDataProcess: -1,
		name:           name,
		lastSync:       math.MaxInt64 - 10,
		Client:         client,
	}
}

// extra methods
func (w SyncWrapper) GetAdapter(addr string) ds.SyncAdapterI {
	return w.Adapters.Get(addr)
}

func (w *SyncWrapper) AddSyncAdapter(adapter ds.SyncAdapterI) {
	if w.ViaDataProcess == -1 {
		log.Fatal("SyncWrapper: ViaDataProcess not set")
	}
	w.Adapters.Add(adapter.GetAddress(), adapter.GetAllAddrsForLogs(), adapter)
	w.lastSync = utils.Min(adapter.GetLastSync(), w.lastSync)
}

func (w *SyncWrapper) GetUnderlyingAdapterAddrs() (addrs []string) {
	for _, adapter := range w.Adapters.GetAll() {
		if !adapter.IsDisabled() {
			addrs = append(addrs, adapter.GetAddress())
		}
	}
	return
}

// //////////
// //////////
func (s SyncWrapper) Topics() [][]common.Hash {
	adapters := s.Adapters.GetAll()
	if len(adapters) == 0 {
		return nil
	}
	return adapters[0].Topics()
}

func (w *SyncWrapper) GetDataProcessType() int {
	if w.ViaDataProcess == -1 {
		return ds.ViaLog
	}
	return w.ViaDataProcess
}

func (s SyncWrapper) GetName() string {
	return s.name
}
func (s SyncWrapper) GetAddress() string {
	return s.name
}

func (SyncWrapper) HasUnderlyingStateToSave() bool {
	return false
}

func (SyncWrapper) GetUnderlyingState() interface{} {
	return nil
}

func (SyncWrapper) Query(queryTill int64) {
}

func (SyncWrapper) GetDetails() core.Json {
	return nil
}

func (SyncWrapper) GetDetailsByKey(key string) string {
	return ""
}

func (SyncWrapper) GetDiscoveredAt() int64 {
	return 0
}
func (SyncWrapper) GetBlockToDisableOn() int64 {
	return math.MaxInt64
}
func (SyncWrapper) IsDisabled() bool {
	return false
}

func (SyncWrapper) SetBlockToDisableOn(int64) {
}

// /
func (SyncWrapper) GetVersion() core.VersionType {
	return 1
}
func (w SyncWrapper) GetLastSync() int64 {
	return w.lastSync
}

func (s SyncWrapper) OnLogs(txLog []types.Log) {
	var lastBlockNum int64 = 0
	for _, txLog := range txLog {
		//
		newBlockNum := int64(txLog.BlockNumber)
		if lastBlockNum == 0 {
			lastBlockNum = newBlockNum
		}
		if lastBlockNum != newBlockNum {
			s.onBlockChange(lastBlockNum)
			lastBlockNum = newBlockNum
		}
		//
		s.OnLog(txLog)
	}
	if lastBlockNum != 0 {
		s.onBlockChange(lastBlockNum)
	}
}

func (s SyncWrapper) onBlockChange(lastBlockNum int64) {
	adapters := s.Adapters.GetAll()
	//
	calls := make([]multicall.Multicall2Call, 0, len(adapters))
	processFns := make([]func(multicall.Multicall2Result), 0, len(adapters))
	//
	for _, adapter := range adapters {
		if adapter.GetLastSync() >= lastBlockNum {
			continue
		}
		switch v := adapter.(type) {
		case *pool.Pool:
			call, processFn := v.OnBlockChange(lastBlockNum)
			// if process fn is not null
			if processFn != nil {
				processFns = append(processFns, processFn)
				calls = append(calls, call)
			}
		}
	}
	results := core.MakeMultiCall(s.Client, lastBlockNum, false, calls)
	for ind, result := range results {
		processFns[ind](result)
	}
}

func (s SyncWrapper) OnLog(txLog types.Log) {
	adapter := s.Adapters.GetFromLogAddr(txLog.Address.Hex())
	if adapter.GetLastSync() < int64(txLog.BlockNumber) {
		adapter.OnLog(txLog)
	}
}

func (s SyncWrapper) GetAdapters() (states []ds.SyncAdapterI) {
	return s.Adapters.GetAll()
}

// ///////
// if not disabled, then do the operation on the underlying sync adapter
// ///////
func (w *SyncWrapper) GetAllAddrsForLogs() (addrs []common.Address) {
	adapters := w.Adapters.GetAll()
	addrs = make([]common.Address, 0, len(adapters))
	for _, cf := range adapters {
		if !cf.IsDisabled() {
			addrs = append(addrs, cf.GetAllAddrsForLogs()...)
		}
	}
	return
}

func (s SyncWrapper) AfterSyncHook(syncTill int64) {
	adapters := s.Adapters.GetAll()
	for _, cf := range adapters {
		if !cf.IsDisabled() {
			cf.AfterSyncHook(syncTill)
		}
	}
}

func (s *SyncWrapper) WillBeSyncedTo(blockNum int64) {
	s.WillSyncTill = blockNum
	adapters := s.Adapters.GetAll()
	for _, adapter := range adapters {
		// if last sync is smaller then new sync till
		if adapter.GetLastSync() < blockNum && !adapter.IsDisabled() {
			adapter.WillBeSyncedTo(blockNum)
		}
	}
}

func (SyncWrapper) GetAdapterState() *ds.SyncAdapter {
	return nil
}

func (SyncWrapper) SetUnderlyingState(interface{}) {

}
