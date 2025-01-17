package repository

import (
	"github.com/Gearbox-protocol/sdk-go/core/schemas"
	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/utils"
)

func (repo *Repository) loadPool() {
	defer utils.Elapsed("loadPool")()
	data := []*schemas.PoolState{}
	err := repo.db.Find(&data).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, pool := range data {
		adapter := repo.GetAdapter(pool.Address)
		adapter.SetUnderlyingState(pool)
	}
}

func (repo *Repository) AddPoolLedger(pl *schemas.PoolLedger) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if "AddLiquidity" == pl.Event {
		repo.AddPoolUniqueUser(pl.Pool, pl.User)
	}
	repo.SetAndGetBlock(pl.BlockNumber).AddPoolLedger(pl)
}
