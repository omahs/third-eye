package repository

import (
	"github.com/Gearbox-protocol/third-eye/core"
	"github.com/Gearbox-protocol/third-eye/log"
	"math/big"
)

func (repo *Repository) loadLastDebtSync() int64 {
	data := core.DebtSync{}
	query := "SELECT max(last_calculated_at) as last_calculated_at FROM debt_sync"
	err := repo.db.Raw(query).Find(&data).Error
	if err != nil {
		log.Fatal(err)
	}
	return data.LastCalculatedAt
}

func (repo *Repository) loadLastAdapterSync() int64 {
	data := core.DebtSync{}
	query := "SELECT max(last_sync) as last_calculated_at FROM sync_adapters"
	err := repo.db.Raw(query).Find(&data).Error
	if err != nil {
		log.Fatal(err)
	}
	return data.LastCalculatedAt
}

func (repo *Repository) AddDebt(debt *core.Debt, forceAdd bool) {
	if repo.config.ThrottleDebtCal {
		lastDebt := repo.lastDebts[debt.SessionId]
		// add debt if throttle is enabled and (last debt is missing or forced add is set)
		if lastDebt == nil || forceAdd {
			repo.addLastDebt(debt)
			repo.debts = append(repo.debts, debt)
		} else if (debt.BlockNumber-lastDebt.BlockNumber) >= core.NoOfBlocksPerHr ||
			core.DiffMoreThanFraction(lastDebt.TotalValueBI, debt.TotalValueBI, big.NewFloat(0.05)) ||
			core.DiffMoreThanFraction(lastDebt.BorrowedAmountPlusInterestBI, debt.CalBorrowedAmountPlusInterestBI, big.NewFloat(0.05)) ||
			// add debt when the health factor is on different side of 10000 from the lastdebt
			(debt.CalHealthFactor >= 10000) != (lastDebt.CalHealthFactor >= 10000) {
			repo.addLastDebt(debt)
			repo.debts = append(repo.debts, debt)
		}
	} else {
		repo.debts = append(repo.debts, debt)
	}
}

func (repo *Repository) loadLastDebts() {
	data := []*core.Debt{}
	query := `SELECT debts.* FROM 
			(SELECT max(block_num), session_id FROM debts GROUP BY session_id) debt_max_block
			JOIN debts ON debt_max_block.max = debts.block_num AND debt_max_block.session_id = debts.session_id`
	err := repo.db.Raw(query).Find(&data).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, debt := range data {
		repo.addLastDebt(debt)
	}
}

func (repo *Repository) addLastDebt(debt *core.Debt) {
	repo.lastDebts[debt.SessionId] = debt
}
