package repository

import (
	"github.com/Gearbox-protocol/third-eye/core"
	"github.com/Gearbox-protocol/third-eye/log"
	"math/big"
)

// for credit filter
func (repo *Repository) AddAllowedProtocol(p *core.Protocol) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.blocks[p.BlockNumber].AddAllowedProtocol(p)
}

func (repo *Repository) AddAllowedToken(atoken *core.AllowedToken) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.blocks[atoken.BlockNumber].AddAllowedToken(atoken)
}

func (repo *Repository) loadAllowedTokenThreshold() {
	data := []*core.AllowedToken{}
	query := `SELECT * FROM allowed_tokens 
	JOIN (SELECT max(block_num) as bn, token, credit_manager FROM allowed_tokens group by token,credit_manager) as atokens
	ON atokens.bn = allowed_tokens.block_num
	AND atokens.credit_manager = allowed_tokens.credit_manager
	AND atokens.token = allowed_tokens.token;`
	err := repo.db.Raw(query).Find(&data).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, atoken := range data {
		repo.AddAllowedTokenThreshold(atoken)
	}
}

func (repo *Repository) AddAllowedTokenThreshold(atoken *core.AllowedToken) {
	if repo.allowedTokensThreshold[atoken.Token] == nil {
		repo.allowedTokensThreshold[atoken.Token] = make(map[string]*core.BigInt)
	}
	value, err := big.NewInt(0).SetString(atoken.LiquidityThreshold, 10)
	if !err {
		log.Fatal("Parsing liquidity threshold failed")
	}
	repo.allowedTokensThreshold[atoken.Token][atoken.CreditManager] = (*core.BigInt)(value)
}
