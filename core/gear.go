package core

import (
	"github.com/Gearbox-protocol/third-eye/artifacts/dataCompressor/mainnet"
	"math/big"
)

type EngineI interface {
	SyncHandler()
}

type Protocol struct {
	Id            string `gorm:"primaryKey;column:id;autoincrement:true"`
	Protocol      string `gorm:"column:protocol"`
	Adapter       string `gorm:"column:adapter"`
	BlockNumber   int64  `gorm:"column:block_num"`
	CreditManager string `gorm:"column:credit_manager"`
}

func (Protocol) TableName() string {
	return "allowed_protocols"
}

type RepositoryI interface {
	// sync adapters
	GetKit() *AdapterKit
	AddSyncAdapter(adapterI SyncAdapterI)
	// saving to the db
	Flush() error
	// adding block/timestamp
	SetBlock(blockNum int64)
	GetBlocks() map[int64]*Block
	LoadBlocks(from, to int64)
	// credit account operations
	AddAccountOperation(accountOperation *AccountOperation)
	// for getting executeparser
	GetExecuteParser() ExecuteParserI
	// price feed/oracle funcs
	AddTokenOracle(token, oracle, feed string, blockNum int64)
	AddPriceFeed(blockNum int64, pf *PriceFeed)
	// token funcs
	AddAllowedProtocol(p *Protocol)
	AddToken(token string) *Token
	AddAllowedToken(atoken *AllowedToken)
	DisableAllowedToken(creditManager, token string, disableBlockNum int64)
	AddTokenObj(token *Token)
	GetToken(addr string) *Token
	ConvertToBalanceWithMask(balances []mainnet.DataTypesTokenBalance, mask *big.Int) (*JsonBalance, error)
	// credit session funcs
	AddCreditSession(session *CreditSession, loadedFromDB bool)
	GetCreditSession(sessionId string) *CreditSession
	GetSessions() map[string]*CreditSession
	GetPriceInUSD(blockNum int64, token string, amount *big.Int) *big.Int
	// credit session snapshots funcs
	AddCreditSessionSnapshot(css *CreditSessionSnapshot)
	AddEventBalance(eb EventBalance)
	// dc
	GetDCWrapper() *DataCompressorWrapper
	AddDataCompressor(blockNum int64, addr string)
	// pools
	AddPoolStat(ps *PoolStat)
	AddPoolLedger(pl *PoolLedger)
	GetPoolUniqueUserLen(pool string) int
	// weth
	SetWETHAddr(address string)
	GetWETHAddr() string
	GetUSDCAddr() string
	// credit manager
	AddCreditManagerToFilter(cmAddr, cfAddr string)
	GetMask(blockNum int64, cmAddr, accountAddr string) *big.Int
	AddCreditManagerStats(cms *CreditManagerStat)
	GetCMState(cmAddr string) *CreditManagerState
	GetUnderlyingDecimal(cmAddr string) int8
	AddRepayOnCM(blockNum int64, cm string, pnl PnlOnRepay)
	GetRepayOnCM(blockNum int64, cm string) *PnlOnRepay
	AddParameters(*Parameters)
	//
	LoadLastDebtSync() int64
	LoadLastAdapterSync() int64
	Clear()
}
