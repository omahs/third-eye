package repository

import (
	"sync"

	"github.com/Gearbox-protocol/third-eye/artifacts/dataCompressor/mainnet"
	"github.com/Gearbox-protocol/third-eye/config"
	"github.com/Gearbox-protocol/third-eye/core"
	"github.com/Gearbox-protocol/third-eye/ethclient"
	"github.com/Gearbox-protocol/third-eye/log"
	"github.com/Gearbox-protocol/third-eye/utils"
	"gorm.io/gorm"

	"context"
	"math/big"
)

type Repository struct {
	// mutex
	mu *sync.Mutex
	// object fx objects
	WETHAddr      string
	db            *gorm.DB
	kit           *core.AdapterKit
	client        *ethclient.Client
	executeParser core.ExecuteParserI
	config        *config.Config
	dcWrapper     *core.DataCompressorWrapper
	// blocks/token
	blocks map[int64]*core.Block
	tokens map[string]*core.Token
	// changed during syncing
	sessions            map[string]*core.CreditSession
	poolUniqueUsers     map[string]map[string]bool
	tokensCurrentOracle map[string]*core.TokenOracle
	// modified after sync loop
	lastCSS        map[string]*core.CreditSessionSnapshot
	tokenLastPrice map[string]*core.PriceFeed
	//// credit_manager -> token -> liquidity threshold
	allowedTokensThreshold map[string]map[string]*core.BigInt
	poolLastInterestData   map[string]*core.PoolInterestData
	debts                  []*core.Debt
	lastDebts              map[string]*core.Debt
}

func NewRepository(db *gorm.DB, client *ethclient.Client, config *config.Config, ep core.ExecuteParserI) core.RepositoryI {
	r := &Repository{
		db:                     db,
		mu:                     &sync.Mutex{},
		client:                 client,
		config:                 config,
		blocks:                 make(map[int64]*core.Block),
		executeParser:          ep,
		kit:                    core.NewAdapterKit(),
		tokens:                 make(map[string]*core.Token),
		sessions:               make(map[string]*core.CreditSession),
		lastCSS:                make(map[string]*core.CreditSessionSnapshot),
		poolUniqueUsers:        make(map[string]map[string]bool),
		tokensCurrentOracle:    make(map[string]*core.TokenOracle),
		tokenLastPrice:         make(map[string]*core.PriceFeed),
		allowedTokensThreshold: make(map[string]map[string]*core.BigInt),
		poolLastInterestData:   make(map[string]*core.PoolInterestData),
		dcWrapper:              core.NewDataCompressorWrapper(client),
		lastDebts:              make(map[string]*core.Debt),
	}
	r.init()
	return r
}

func (repo *Repository) GetDCWrapper() *core.DataCompressorWrapper {
	return repo.dcWrapper
}

func (repo *Repository) GetExecuteParser() core.ExecuteParserI {
	return repo.executeParser
}

func (repo *Repository) init() {
	lastDebtSync := repo.loadLastDebtSync()
	// token should be loaded before syncAdapters as credit manager adapter uses underlying token details
	repo.loadToken()
	// syncadapter state for cm and pool is set after loading of pool/credit manager table data from db
	repo.loadSyncAdapters()
	repo.loadCurrentTokenOracle()
	repo.loadPool()
	repo.loadCreditManagers()
	repo.loadCreditSessions(lastDebtSync)
	repo.debtInit()
}

func (repo *Repository) debtInit() {
	lastDebtSync := repo.loadLastDebtSync()
	repo.loadLastCSS(lastDebtSync)
	repo.loadTokenLastPrice(lastDebtSync)
	repo.loadAllowedTokenThreshold(lastDebtSync)
	repo.loadPoolLastInterestData(lastDebtSync)
	repo.loadLastDebts()
	// process blocks for calculating debts
	adaptersSyncedTill := repo.loadLastAdapterSync()
	var batchSize int64 = 1000
	for ; lastDebtSync+batchSize < adaptersSyncedTill; lastDebtSync += batchSize {
		repo.processBlocksInBatch(lastDebtSync, lastDebtSync+batchSize)
	}
	repo.processBlocksInBatch(lastDebtSync, adaptersSyncedTill)
}

func (repo *Repository) processBlocksInBatch(from, to int64) {
	repo.loadBlocks(from, to)
	if len(repo.blocks) > 0 {
		repo.calculateDebtAndClear()
	}
}

func (repo *Repository) AddAccountOperation(accountOperation *core.AccountOperation) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.blocks[accountOperation.BlockNumber].AddAccountOperation(accountOperation)
}

func (repo *Repository) SetBlock(blockNum int64) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.blocks[blockNum] == nil {
		b, err := repo.client.BlockByNumber(context.Background(), big.NewInt(blockNum))
		if err != nil {
			log.Fatal(err)
		}
		repo.blocks[blockNum] = &core.Block{BlockNumber: blockNum, Timestamp: b.Time()}
	}
}

func (repo *Repository) AddCreditManagerStats(cms *core.CreditManagerStat) {
	repo.blocks[cms.BlockNum].AddCreditManagerStats(cms)
}

func (repo *Repository) GetKit() *core.AdapterKit {
	return repo.kit
}

func (repo *Repository) AddEventBalance(eb core.EventBalance) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.blocks[eb.BlockNumber].AddEventBalance(&eb)
}

func (repo *Repository) ConvertToBalance(balances []mainnet.DataTypesTokenBalance) *core.JsonBalance {
	jsonBalance := core.JsonBalance{}
	for _, token := range balances {
		tokenAddr := token.Token.Hex()
		if token.Balance.Sign() != 0 {
			jsonBalance[tokenAddr] = &core.BalanceType{
				BI: (*core.BigInt)(token.Balance),
				F:  utils.GetFloat64Decimal(token.Balance, repo.GetToken(tokenAddr).Decimals),
			}
		}
	}
	return &jsonBalance
}

func (repo *Repository) loadBlocks(from, to int64) {
	data := []*core.Block{}
	err := repo.db.Preload("CSS").Preload("PoolStats").
		Preload("AllowedTokens").Preload("PriceFeeds").
		Find(&data, "id > ? AND id <= ?", from, to).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, block := range data {
		repo.blocks[block.BlockNumber] = block
	}
}

func (repo *Repository) FlushAndDebt() {
	repo.Flush()
	repo.calculateDebtAndClear()
}

func (repo *Repository) calculateDebtAndClear() {
	if !repo.config.DisableDebtEngine {
		repo.calculateDebt()
	}
	repo.clear()
}

func (repo *Repository) SetWETHAddr(addr string) {
	repo.WETHAddr = addr
}
