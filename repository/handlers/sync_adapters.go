package handlers

import (
	"fmt"
	"math"
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
	"github.com/Gearbox-protocol/third-eye/models/contract_register"
	"github.com/Gearbox-protocol/third-eye/models/credit_filter"
	"github.com/Gearbox-protocol/third-eye/models/credit_manager"
	"github.com/Gearbox-protocol/third-eye/models/gear_token"
	"github.com/Gearbox-protocol/third-eye/models/pool"
	"github.com/Gearbox-protocol/third-eye/models/price_oracle"
	"github.com/Gearbox-protocol/third-eye/models/treasury"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SyncAdaptersRepo struct {
	kit            *ds.AdapterKit
	AggregatedFeed *aggregated_block_feed.AggregatedBlockFeed
	r              ds.RepositoryI
	client         core.ClientI
	extras         *ExtrasRepo
	rollback       string
	mu             *sync.Mutex
}

func NewSyncAdaptersRepo(client core.ClientI, repo ds.RepositoryI, cfg *config.Config, extras *ExtrasRepo) *SyncAdaptersRepo {
	obj := &SyncAdaptersRepo{
		kit:      ds.NewAdapterKit(),
		client:   client,
		r:        repo,
		extras:   extras,
		rollback: cfg.Rollback,
		mu:       &sync.Mutex{},
	}
	// aggregated block feed
	obj.AggregatedFeed = aggregated_block_feed.NewAggregatedBlockFeed(client, repo, cfg.Interval)
	obj.kit.Add(obj.AggregatedFeed)
	return obj
}

func (repo *SyncAdaptersRepo) addSyncAdapter(adapterI ds.SyncAdapterI) {
	// if ds.GearToken == adapterI.GetName() {
	// 	repo.GearTokenAddr = adapterI.GetAddress()
	// }
	if adapterI.GetName() == ds.QueryPriceFeed {
		repo.AggregatedFeed.AddYearnFeed(adapterI)
	} else {
		repo.kit.Add(adapterI)
	}
}

// load/save
func (repo *SyncAdaptersRepo) LoadSyncAdapters(db *gorm.DB) {
	defer utils.Elapsed("loadSyncAdapters")()
	//
	data := []*ds.SyncAdapter{}
	err := db.Find(&data, "disabled = ? OR type = 'PriceOracle'", false).Error
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
			if adapter.GetName() != ds.AggregatedBlockFeed {
				adapters = append(adapters, adapter.GetAdapterState())
			}
			if adapter.HasUnderlyingState() {
				err := tx.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(adapter.GetUnderlyingState()).Error
				log.CheckFatal(err)
			}
		}
		repo.kit.Reset(lvlIndex)
	}
	// save qyery feeds from AggregatedFeed
	for _, adapter := range repo.AggregatedFeed.GetQueryFeeds() {
		adapters = append(adapters, adapter.GetAdapterState())
	}
	err := tx.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(adapters, 50).Error
	log.CheckFatal(err)

	if uniPools := repo.AggregatedFeed.GetUniswapPools(); len(uniPools) > 0 {
		err := tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).CreateInBatches(uniPools, 50).Error
		log.CheckFatal(err)
	}
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
	case ds.QueryPriceFeed:
		return aggregated_block_feed.NewQueryPriceFeedFromAdapter(adapter)
	case ds.ContractRegister:
		return contract_register.NewContractRegisterFromAdapter(adapter)
	case ds.GearToken:
		return gear_token.NewGearTokenFromAdapter(adapter)
	case ds.Treasury:
		return treasury.NewTreasuryFromAdapter(adapter)
	case ds.AccountManager:
		return account_manager.NewAccountManagerFromAdapter(adapter)
	case ds.CreditConfigurator:
		return credit_filter.NewCreditFilterFromAdapter(adapter)
	case ds.CreditFilter:
		if adapter.Details["creditManager"] != nil {
			cmAddr := adapter.Details["creditManager"].(string)
			repo.extras.AddCreditManagerToFilter(cmAddr, adapter.GetAddress())
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
	if repo.rollback == "1" {
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

func (repo *SyncAdaptersRepo) GetAdapter(addr string) ds.SyncAdapterI {
	adapter := repo.GetKit().GetAdapter(addr)
	if adapter == nil {
		feeds := repo.AggregatedFeed.GetQueryFeeds()
		for _, feed := range feeds {
			if feed.GetAddress() == addr {
				return feed
			}
		}
	}
	return adapter
}

func (repo *SyncAdaptersRepo) GetYearnFeedAddrs() (addrs []string) {
	feeds := repo.AggregatedFeed.GetQueryFeeds()
	for _, adapter := range feeds {
		addrs = append(addrs, adapter.GetAddress())
	}
	return
}

////////////////////
// for price oracle
////////////////////

// return the active first oracle under blockNum
// if all disabled return the last one
func (repo *SyncAdaptersRepo) GetActivePriceOracleByBlockNum(blockNum int64) (string, error) {
	var disabledLastOracle, activeFirstOracle string
	var disabledOracleBlock, activeOracleBlock int64
	activeOracleBlock = math.MaxInt64
	oracles := repo.kit.GetAdapterAddressByName(ds.PriceOracle)
	for _, addr := range oracles {
		oracleAdapter := repo.GetAdapter(addr)
		if oracleAdapter.GetDiscoveredAt() <= blockNum {
			if oracleAdapter.IsDisabled() {
				if disabledOracleBlock < oracleAdapter.GetDiscoveredAt() {
					disabledOracleBlock = oracleAdapter.GetDiscoveredAt()
					disabledLastOracle = addr
				}
			} else {
				if activeOracleBlock > oracleAdapter.GetDiscoveredAt() {
					activeOracleBlock = oracleAdapter.GetDiscoveredAt()
					activeFirstOracle = addr
				}
			}
		}
	}
	if activeFirstOracle != "" {
		return activeFirstOracle, nil
	} else if disabledLastOracle != "" {
		return disabledLastOracle, nil
	} else {
		return "", fmt.Errorf("Not Found")
	}
}

func (repo *SyncAdaptersRepo) GetPriceOracleByVersion(version int16) (string, error) {
	addrProviderAddr := repo.kit.GetAdapterAddressByName(ds.AddressProvider)
	addrProvider := repo.GetAdapter(addrProviderAddr[0])
	details := addrProvider.GetDetails()
	if details != nil {
		priceOracles, ok := details["priceOracles"].(map[string]interface{})
		if ok {
			value, ok := priceOracles[fmt.Sprintf("%d", version)].(string)
			if ok {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("Not Found")
}