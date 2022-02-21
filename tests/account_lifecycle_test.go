package tests

import (
	"testing"

	"github.com/Gearbox-protocol/third-eye/config"
	"github.com/Gearbox-protocol/third-eye/debts"
	"github.com/Gearbox-protocol/third-eye/engine"
	"github.com/Gearbox-protocol/third-eye/log"
	"github.com/Gearbox-protocol/third-eye/tests/framework"
	"github.com/Gearbox-protocol/third-eye/repository"
	"github.com/Gearbox-protocol/third-eye/utils"
)

func TestLifecycleCreditAccount(t *testing.T) {
	log.SetTestLogging(t)
	client := framework.NewTestClient()
	cfg := &config.Config{}
	repo := repository.GetRepository(nil, client, cfg, nil)
	debtEng := debts.NewDebtEngine(nil, client, cfg, repo)
	eng := engine.NewEngine(cfg, client, debtEng, repo)
	r := framework.NewMockRepo(repo, client,  "account_lifecycle/input.json", t, eng)
	r.Init()
	log.Info(utils.ToJson(r.AddressMap))
	eng.Sync(10)

	outputBlocks := repo.GetBlocks()
	delete(outputBlocks, 2)
	r.Check(t, outputBlocks, "account_lifecycle/blocks.json")
	debtEng.CalculateDebt()
	r.Check(t, debtEng.GetDebts(), "account_lifecycle/debts.json")
}