package debts

import (
	"github.com/Gearbox-protocol/third-eye/core"
	"github.com/Gearbox-protocol/third-eye/log"
)

// token threshold
func (eng *DebtEngine) loadAllowedTokenThreshold(lastDebtSync int64) {
	data := []*core.AllowedToken{}
	query := `SELECT * FROM allowed_tokens 
	JOIN (SELECT max(block_num) as bn, token, credit_manager FROM allowed_tokens 
		WHERE block_num <= ? group by token,credit_manager) as atokens
	ON atokens.bn = allowed_tokens.block_num
	AND atokens.credit_manager = allowed_tokens.credit_manager
	AND atokens.token = allowed_tokens.token
	WHERE block_num <= ? ORDER BY block_num;`
	err := eng.db.Raw(query, lastDebtSync, lastDebtSync).Find(&data).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, atoken := range data {
		eng.AddAllowedTokenThreshold(atoken)
	}
}

func (eng *DebtEngine) AddAllowedTokenThreshold(atoken *core.AllowedToken) {
	if eng.allowedTokensThreshold[atoken.CreditManager] == nil {
		eng.allowedTokensThreshold[atoken.CreditManager] = make(map[string]*core.BigInt)
	}
	eng.allowedTokensThreshold[atoken.CreditManager][atoken.Token] = atoken.LiquidityThreshold
}

// token price from feeds
func (eng *DebtEngine) loadTokenLastPrice(lastDebtSync int64) {
	data := []*core.PriceFeed{}
	query := `SELECT price_feeds.* FROM price_feeds
	JOIN (SELECT max(block_num) AS bn, token FROM price_feeds WHERE block_num <= ? GROUP BY token) AS max_pf
	ON max_pf.bn = price_feeds.block_num AND max_pf.token = price_feeds.token WHERE block_num <= ?`
	err := eng.db.Raw(query, lastDebtSync, lastDebtSync).Find(&data).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, tokenPrice := range data {
		eng.AddTokenLastPrice(tokenPrice)
	}
}

func (eng *DebtEngine) AddTokenLastPrice(pf *core.PriceFeed) {
	eng.tokenLastPrice[pf.Token] = pf
}