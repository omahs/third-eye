package core

import (
	"context"
	"github.com/Gearbox-protocol/gearscan/ethclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"math/big"
)

const MaxUint = ^int64(0)

type SyncAdapter struct {
	*Contract
	LastSync int64 `gorm:"column:last_sync"`
}

func (SyncAdapter) TableName() string {
	return "sync_adapters"
}

type SyncAdapterI interface {
	OnLog(txLog types.Log)
	GetLastSync() int64
	SetLastSync(int64)
	GetAdapterState() *SyncAdapter
	GetAddress() string
	FirstSync() bool
	GetName() string
}

func (s *SyncAdapter) SetLastSync(lastSync int64) {
	s.LastSync = lastSync
}

func (s *SyncAdapter) FirstSync() bool {
	return s.FirstLogAt == s.LastSync
}

func NewSyncAdapter(addr, name string, discoveredAt int64,  client *ethclient.Client) *SyncAdapter {
	obj := &SyncAdapter{
			Contract: NewContract(addr, name, discoveredAt, client),
		}
	obj.LastSync = obj.FirstLogAt
	return obj
}
func (s *SyncAdapter) GetAdapterState() *SyncAdapter {
	return s
}
func (s *SyncAdapter) LoadState() {

}


// func (mdl *SyncAdapter) OnLog(txLog types.Log) {
// 	log.Infof("%s\n", reflect.TypeOf(mdl))
// 	log.Infof("%+v\n", txLog)
// }

func (mdl *SyncAdapter) GetLastSync() int64 {
	return mdl.LastSync
}

func (mdl *SyncAdapter) Monitor(startBlock, endBlock int64) (chan types.Log, event.Subscription, error) {
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetInt64(startBlock),
		ToBlock:   new(big.Int).SetInt64(endBlock),
		Addresses: []common.Address{common.HexToAddress(mdl.Address)},
	}
	var logs = make(chan types.Log, 2)
	s, err := mdl.Client.SubscribeFilterLogs(context.TODO(), query, logs)
	if err != nil {
		return logs, s, err
	}
	return logs, s, nil
}
