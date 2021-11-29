package price_oracle

import (
	"github.com/Gearbox-protocol/gearscan/core"
	"github.com/Gearbox-protocol/gearscan/log"
	"github.com/Gearbox-protocol/gearscan/models/price_feed"
	"github.com/ethereum/go-ethereum/core/types"
)

func (mdl *PriceOracle) OnLog(txLog types.Log) {
	blockNum := int64(txLog.BlockNumber)
	switch txLog.Topics[0] {
	case core.Topic("NewPriceFeed(address,address)"):

		newPriceFeedEvent, err := mdl.contractETH.ParseNewPriceFeed(txLog)
		if err != nil {
			log.Fatal("[PriceOracle]: Cant unpack NewPriceFeed event", err)
		}

		token := newPriceFeedEvent.Token.Hex()
		oracle := newPriceFeedEvent.PriceFeed.Hex()
		mdl.Repo.AddTokenOracle(token, oracle, blockNum)
		log.Info(token, oracle)
		obj := price_feed.NewPriceFeed(oracle, token, blockNum, mdl.SyncAdapter.Client, mdl.State.Repo)
		mdl.State.Repo.AddSyncAdapter(obj)
	}
}
