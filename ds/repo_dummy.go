package ds

import (
	"math/big"

	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/Gearbox-protocol/sdk-go/core/schemas"
	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/third-eye/ds/dc_wrapper"
)

type DummyRepo struct {
}

func (DummyRepo) Init() {

}

// sync adapters
func (DummyRepo) GetAdapter(addr string) SyncAdapterI {
	return nil
}

func (DummyRepo) GetAdapterAddressByName(name string) []string {
	return nil
}

func (DummyRepo) AddSyncAdapter(adapterI SyncAdapterI) {
}
func (DummyRepo) GetChainId() uint {
	return 0
}

// saving to the db
func (DummyRepo) Flush(syncTill int64) error {
	return nil
}

// adding block/timestamp
func (DummyRepo) SetBlock(blockNum int64) {
}
func (DummyRepo) SetAndGetBlock(blockNum int64) *schemas.Block {
	return nil
}
func (DummyRepo) GetBlocks() map[int64]*schemas.Block {
	return nil
}
func (DummyRepo) GetTokenOracles() map[core.VersionType]map[string]*schemas.TokenOracle {
	return nil
}
func (DummyRepo) GetDisabledTokens() []*schemas.AllowedToken {
	return nil
}
func (DummyRepo) LoadBlocks(from, to int64) {
}

// credit account operations
func (DummyRepo) AddAccountOperation(accountOperation *schemas.AccountOperation) {
}

// for getting executeparser
func (DummyRepo) GetExecuteParser() ExecuteParserI {
	return nil
}

// price feed/oracle funcs
func (DummyRepo) DirectlyAddTokenOracle(tokenOracle *schemas.TokenOracle) {
}
func (DummyRepo) GetPrice(token string) *big.Int {
	return nil
}
func (DummyRepo) AddPriceFeed(pf *schemas.PriceFeed) {
}

// token funcs
func (DummyRepo) AddAllowedProtocol(logID uint, txHash, creditFilter string, p *schemas.Protocol) {
}
func (DummyRepo) DisableProtocol(blockNum int64, logID uint, txHash, cm, creditFilter, protocol string) {
}
func (DummyRepo) AddAllowedToken(logID uint, txHash, creditFilter string, atoken *schemas.AllowedToken) {
}
func (DummyRepo) DisableAllowedToken(blockNum int64, logID uint, txHash string, creditManager, creditFilter, token string) {
}

// v2
func (DummyRepo) AddAllowedTokenV2(logID uint, txHash, creditFilter string, atoken *schemas.AllowedToken) {
}
func (DummyRepo) UpdateLimits(logID uint, txHash, creditConfigurator string, params *schemas.Parameters) {
}
func (DummyRepo) UpdateFees(logID uint, txHash, creditConfigurator string, params *schemas.Parameters) {
}
func (DummyRepo) UpdateEmergencyLiqDiscount(logID uint, txHash, creditConfigurator string, params *schemas.Parameters) {
}
func (DummyRepo) TransferAccountAllowed(*schemas.TransferAccountAllowed) {
}
func (DummyRepo) GetPricesInUSD(blockNum int64, tokenAddrs []string) core.JsonFloatMap {
	return nil
}

func (DummyRepo) GetToken(addr string) *schemas.Token {
	return nil
}
func (DummyRepo) GetTokens() []string {
	return nil
}

// credit session funcs
func (DummyRepo) AddCreditSession(session *schemas.CreditSession, loadedFromDB bool, txHash string, logID uint) {
}
func (DummyRepo) GetCreditSession(sessionId string) *schemas.CreditSession {
	return nil
}
func (DummyRepo) UpdateCreditSession(sessionId string, values map[string]interface{}) *schemas.CreditSession {
	return nil
}
func (DummyRepo) GetSessions() map[string]*schemas.CreditSession {
	return nil
}
func (DummyRepo) GetValueInCurrency(blockNum int64, version core.VersionType, token, currency string, amount *big.Int) *big.Int {
	return nil
}
func (DummyRepo) AddDieselToken(dieselToken, underlyingToken, pool string) {
}
func (DummyRepo) GetDieselTokens() map[string]*schemas.UTokenAndPool {
	return nil
}

// credit session snapshots funcs
func (DummyRepo) AddCreditSessionSnapshot(css *schemas.CreditSessionSnapshot) {
}

// dc
func (DummyRepo) GetDCWrapper() *dc_wrapper.DataCompressorWrapper {
	return nil
}

// pools
func (DummyRepo) AddPoolStat(ps *schemas.PoolStat) {
}
func (DummyRepo) AddDieselTransfer(dt *schemas.DieselTransfer) {
}

var Count int64

func (DummyRepo) AddRebaseDetailsForDB(transfer *schemas.RebaseDetailsForDB) {
	Count += 1
}
func (DummyRepo) AddPoolLedger(pl *schemas.PoolLedger) {
}
func (DummyRepo) GetPoolUniqueUserLen(pool string) int {
	return 0
}
func (DummyRepo) IsDieselToken(token string) bool {
	return false
}
func (DummyRepo) GetWETHAddr() string {
	return ""
}
func (DummyRepo) GetUSDCAddr() string {
	return ""
}
func (DummyRepo) GetGearTokenAddr() string {
	return ""
}

// credit manager
func (DummyRepo) AddAccountTokenTransfer(tt *schemas.TokenTransfer) {
}
func (DummyRepo) AddCreditManagerStats(cms *schemas.CreditManagerStat) {
}
func (DummyRepo) GetCMState(cmAddr string) *schemas.CreditManagerState {
	return nil
}
func (DummyRepo) GetUnderlyingDecimal(cmAddr string) int8 {
	return 0
}
func (DummyRepo) AddRepayOnCM(cm string, pnl schemas.PnlOnRepay) {
}
func (DummyRepo) AddParameters(logID uint, txHash string, params *schemas.Parameters, token string) {
}
func (DummyRepo) AddFastCheckParams(logID uint, txHash, cm, creditFilter string, fcParams *schemas.FastCheckParams) {
}
func (DummyRepo) AfterSync(blockNum int64) {
}
func (DummyRepo) GetAccountManager() *DirectTransferManager {
	return nil
}
func (DummyRepo) AddAccountAddr(account string) {
}

// dao
func (DummyRepo) AddDAOOperation(operation *schemas.DAOOperation) {
}
func (DummyRepo) CalCurrentTreasuryValue(syncTill int64) {
}
func (DummyRepo) AddTreasuryTransfer(blockNum int64, logID uint, token string, amount *big.Int, operationTransfer bool) {
}
func (DummyRepo) RecentMsgf(headers log.RiskHeader, msg string, args ...interface{}) {
}

// oracle
func (DummyRepo) GetYearnFeedAddrs() []string {
	return nil
}

// has mutex lock
func (DummyRepo) AddNewPriceOracleEvent(tokenOracle *schemas.TokenOracle, bounded bool) {
}
func (DummyRepo) GetOracleForV2Token(token string) *schemas.TokenOracle {
	return nil
}

func (DummyRepo) LoadLastDebtSync() int64 {
	return 0
}
func (DummyRepo) LoadLastAdapterSync() int64 {
	return 0
}
func (DummyRepo) Clear() {
}

// multicall
func (DummyRepo) ChainlinkPriceUpdatedAt(token string, blockNums []int64) {
}

// for testing
func (DummyRepo) AddTokenObj(token *schemas.Token) {
}
func (DummyRepo) PrepareSyncAdapter(adapter *SyncAdapter) SyncAdapterI {
	return nil
}

func (DummyRepo) GetTokenFromSdk(string) string {
	return ""
}
