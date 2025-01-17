package handlers

import (
	"fmt"
	"sync"

	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/Gearbox-protocol/third-eye/config"
	"github.com/Gearbox-protocol/third-eye/ds"
	"github.com/Gearbox-protocol/third-eye/models/account_factory"
	"github.com/Gearbox-protocol/third-eye/models/account_manager"
	"github.com/Gearbox-protocol/third-eye/models/acl"
	"github.com/Gearbox-protocol/third-eye/models/address_provider"
	"github.com/Gearbox-protocol/third-eye/models/aggregated_block_feed"
	"github.com/Gearbox-protocol/third-eye/models/chainlink_price_feed"
	"github.com/Gearbox-protocol/third-eye/models/composite_chainlink"
	"github.com/Gearbox-protocol/third-eye/models/contract_register"
	"github.com/Gearbox-protocol/third-eye/models/credit_filter"
	"github.com/Gearbox-protocol/third-eye/models/credit_manager"
	"github.com/Gearbox-protocol/third-eye/models/gear_token"
	"github.com/Gearbox-protocol/third-eye/models/pool"
	"github.com/Gearbox-protocol/third-eye/models/pool_lmrewards"
	"github.com/Gearbox-protocol/third-eye/models/price_oracle"
	"github.com/Gearbox-protocol/third-eye/models/rebase_token"
	"github.com/Gearbox-protocol/third-eye/models/treasury"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SyncAdaptersRepo struct {
	r               ds.RepositoryI
	client          core.ClientI
	extras          *ExtrasRepo
	rollbackAllowed bool
	mu              *sync.Mutex
	*AdapterKitHandler
}

func NewSyncAdaptersRepo(client core.ClientI, repo ds.RepositoryI, cfg *config.Config, extras *ExtrasRepo) *SyncAdaptersRepo {
	obj := &SyncAdaptersRepo{
		client:            client,
		r:                 repo,
		extras:            extras,
		rollbackAllowed:   cfg.Rollback,
		mu:                &sync.Mutex{},
		AdapterKitHandler: NewAdpterKitHandler(client, repo, cfg),
	}
	return obj
}

// load/save
func (repo *SyncAdaptersRepo) LoadSyncAdapters(db *gorm.DB) {
	defer utils.Elapsed("loadSyncAdapters")()
	//
	data := []*ds.SyncAdapter{}
	err := db.Find(&data, "disabled = ? OR disabled is NULL OR type = 'PriceOracle' ORDER BY type", false).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, adapter := range data {
		p := repo.PrepareSyncAdapter(adapter)
		repo.addSyncAdapter(p)
	}
}

func (repo *SyncAdaptersRepo) Save(tx *gorm.DB) {
	defer utils.Elapsed("sync adapters sql statements")()
	adapters := make([]*ds.SyncAdapter, 0, repo.kit.Len())
	for lvlIndex := 0; lvlIndex < repo.kit.Len(); lvlIndex++ {
		for repo.kit.Next(lvlIndex) {
			adapter := repo.kit.Get(lvlIndex)
			if ds.IsWrapperAdapter(adapter.GetName()) {
				continue
			}
			adapters = append(adapters, adapter.GetAdapterState())
			if adapter.HasUnderlyingStateToSave() {
				err := tx.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(adapter.GetUnderlyingState()).Error
				log.CheckFatal(err)
			}
		}
		repo.kit.Reset(lvlIndex)
	}
	// save wrapper underlying states
	for _, adapter := range repo.AdapterKitHandler.GetAdaptersFromWrapper() {
		if adapter.HasUnderlyingStateToSave() {
			err := tx.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(adapter.GetUnderlyingState()).Error
			log.CheckFatal(err)
		}
		adapters = append(adapters, adapter.GetAdapterState())
	}
	//
	err := tx.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(adapters, 50).Error
	log.CheckFatal(err)

}

// external funcs
// for testing and load from db
func (repo *SyncAdaptersRepo) PrepareSyncAdapter(adapter *ds.SyncAdapter) ds.SyncAdapterI {
	adapter.Client = repo.client
	adapter.Repo = repo.r
	switch adapter.ContractName {
	case ds.ACL:
		return acl.NewACLFromAdapter(adapter)
	case ds.AddressProvider:
		ap := address_provider.NewAddressProviderFromAdapter(adapter)
		if ap.Details["dc"] != nil {
			repo.extras.GetDCWrapper().LoadMultipleDC(ap.Details["dc"])
		}
		return ap
	case ds.PoolLMRewards:
		return pool_lmrewards.NewPoolLMRewardsFromAdapter(adapter)
	case ds.AccountFactory:
		return account_factory.NewAccountFactoryFromAdapter(adapter)
	case ds.Pool:
		return pool.NewPoolFromAdapter(adapter)
	case ds.CreditManager:
		return credit_manager.NewCreditManagerFromAdapter(adapter)
	case ds.PriceOracle:
		return price_oracle.NewPriceOracleFromAdapter(adapter)
	case ds.ChainlinkPriceFeed:
		return chainlink_price_feed.NewChainlinkPriceFeedFromAdapter(adapter, false)
	case ds.CompositeChainlinkPF:
		return composite_chainlink.NewCompositeChainlinkPFFromAdapter(adapter)
	case ds.QueryPriceFeed:
		return aggregated_block_feed.NewQueryPriceFeedFromAdapter(adapter)
	case ds.ContractRegister:
		return contract_register.NewContractRegisterFromAdapter(adapter)
	case ds.GearToken:
		return gear_token.NewGearTokenFromAdapter(adapter)
	case ds.Treasury:
		return treasury.NewTreasuryFromAdapter(adapter)
	case ds.RebaseToken:
		return rebase_token.NewRebaseTokenFromAdapter(adapter)
	case ds.AccountManager:
		return account_manager.NewAccountManagerFromAdapter(adapter)
	case ds.CreditConfigurator:
		return credit_filter.NewCreditFilterFromAdapter(adapter)
	case ds.CreditFilter:
		if adapter.Details["creditManager"] != nil {
			cmAddr := adapter.Details["creditManager"].(string)
			repo.extras.GetDCWrapper().AddCreditManagerToFilter(cmAddr, adapter.GetAddress())
		} else {
			log.Fatal("Credit filter doesn't have credit manager", adapter.GetAddress())
		}
		return credit_filter.NewCreditFilterFromAdapter(adapter)
	}
	return nil
}

func (repo *SyncAdaptersRepo) AddSyncAdapter(newAdapterI ds.SyncAdapterI) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.rollbackAllowed {
		return
	}
	if newAdapterI.GetName() == ds.PriceOracle {
		oldPriceOracleAddrs := repo.kit.GetAdapterAddressByName(ds.PriceOracle)
		for _, addr := range oldPriceOracleAddrs {
			oldPriceOracle := repo.GetAdapter(addr)
			if !oldPriceOracle.IsDisabled() {
				oldPriceOracle.SetBlockToDisableOn(newAdapterI.GetDiscoveredAt())
			}
		}
	}
	repo.addSyncAdapter(newAdapterI)
}

func (repo *SyncAdaptersRepo) GetKit() *ds.AdapterKit {
	return repo.kit
}

////////////////////
// for price oracle
////////////////////

// return the active first oracle under blockNum
// if all disabled return the last one
func (repo *SyncAdaptersRepo) GetActivePriceOracleByBlockNum(blockNum int64) (latestOracle string, version core.VersionType, err error) {
	var latestBlock int64 = 0
	oracles := repo.kit.GetAdapterAddressByName(ds.PriceOracle)
	for _, addr := range oracles {
		oracleAdapter := repo.GetAdapter(addr)
		// only get the oracles that are present at the blocknumg
		if oracleAdapter.GetDiscoveredAt() <= blockNum {
			// if the oracle is discoverdat later
			if latestBlock < oracleAdapter.GetDiscoveredAt() ||
				// if the oracles are discoveredat at same time but the version of oracle is more
				(latestBlock == oracleAdapter.GetDiscoveredAt() && oracleAdapter.GetVersion().MoreThan(version)) {
				latestBlock = oracleAdapter.GetDiscoveredAt()
				latestOracle = addr
				version = oracleAdapter.GetVersion()
			}
		}
	}
	if latestOracle == "" {
		err = fmt.Errorf("not Found")
	}
	return
}

func (repo *SyncAdaptersRepo) GetPriceOracleByDiscoveredAt(blockNum int64) (string, error) {
	addrProviderAddr := repo.kit.GetAdapterAddressByName(ds.AddressProvider)
	addrProvider := repo.GetAdapter(addrProviderAddr[0])
	priceOracle := addrProvider.GetDetailsByKey(fmt.Sprintf("%d", blockNum))
	if priceOracle == "" {
		return "", fmt.Errorf("not Found")
	}
	return priceOracle, nil
}
