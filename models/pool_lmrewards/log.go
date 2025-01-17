package pool_lmrewards

import (
	"math/big"

	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/Gearbox-protocol/sdk-go/core/schemas"
	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/pkg"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (mdl *PoolLMRewards) GetAddrsForLogs() (addrs []common.Address) {
	addrs = append(addrs, common.HexToAddress(mdl.Address))
	return addrs
}

//

//

func (mdl *PoolLMRewards) Topics() [][]common.Hash {
	return [][]common.Hash{{core.Topic("Transfer(address,address,uint256)")}}
}

// onlog for transfer token events only, which is provided by Topics() func
func (mdl *PoolLMRewards) OnLog(txLog types.Log) {
	// conditions to return
	amount, ok := new(big.Int).SetString(common.BytesToHash(txLog.Data).Hex()[2:], 16)
	if !ok {
		log.Fatal("Failed parsing value")
	}
	if amount.Cmp(new(big.Int)) == 0 {
		return
	}
	// calculate rewards
	currentBlock := int64(txLog.BlockNumber)
	if currentBlock != mdl.pendingCalcBlock {
		mdl.calculateRewards(mdl.pendingCalcBlock, currentBlock-1)
		mdl.pendingCalcBlock = currentBlock
	}
	//
	//
	token := mdl.Repo.GetToken(txLog.Address.Hex())
	tokenSym := token.Symbol
	from := common.HexToAddress(txLog.Topics[1].Hex()).Hex()
	to := common.HexToAddress(txLog.Topics[2].Hex()).Hex()
	//
	if to == core.NULL_ADDR.Hex() {
		mdl.totalSupplies[tokenSym] = new(big.Int).Sub(
			utils.NotNilBigInt(mdl.totalSupplies[tokenSym]),
			amount,
		)
	} else {
		mdl.addBalance(tokenSym, to, amount)
	}
	//
	mdl.addDieselTransfer(&schemas.DieselTransfer{
		From:        from,
		To:          to,
		LogId:       int64(txLog.Index),
		TokenSymbol: tokenSym,
		Amount:      utils.GetFloat64Decimal(amount, token.Decimals),
		BlockNum:    int64(txLog.BlockNumber),
	})
	if from == core.NULL_ADDR.Hex() {
		mdl.totalSupplies[tokenSym] = new(big.Int).Add(
			utils.NotNilBigInt(mdl.totalSupplies[tokenSym]),
			amount,
		)
	} else {
		mdl.addBalance(tokenSym, from, new(big.Int).Neg(amount))
	}
}

func (mdl PoolLMRewards) addDieselTransfer(dt *schemas.DieselTransfer) {
	mdl.Repo.AddDieselTransfer(dt)
}

func (mdl PoolLMRewards) addBalance(tokenSym, user string, amount *big.Int) {
	if mdl.dieselBalances[tokenSym] == nil {
		mdl.dieselBalances[tokenSym] = map[string]*big.Int{}
	}
	mdl.dieselBalances[tokenSym][user] = new(big.Int).Add(
		utils.NotNilBigInt(mdl.dieselBalances[tokenSym][user]),
		amount,
	)
}

// inclusive of from and to
func (mdl PoolLMRewards) calculateRewards(from, to int64) {
	snapshots := pkg.GetRewardPerToken(mdl.chainId, from, to)

	snapStart := from
	if len(snapshots) > 0 && from < snapshots[0].Block {
		snapStart = snapshots[0].Block
	}
	for snapInd, snapshot := range snapshots {
		snapEnd := to
		if snapInd != len(snapshots)-1 {
			snapEnd = snapshots[snapInd+1].Block - 1
		}
		for dieselSym, userAndbalance := range mdl.dieselBalances {
			pool := mdl.decimalsAndPool[dieselSym].pool
			rewardPerBlock := utils.NotNilBigInt(snapshot.RewardPerBlock[dieselSym])
			for user, balance := range userAndbalance {
				norm := new(big.Int).Mul(balance, rewardPerBlock)
				reward := new(big.Int)
				if mdl.totalSupplies[dieselSym] != nil {
					userRewardNorm := new(big.Int).Mul(norm, big.NewInt(snapEnd-snapStart+1))
					reward = new(big.Int).Quo(userRewardNorm, mdl.totalSupplies[dieselSym])
				}
				mdl.addUserReward(pool, user, reward)
				// log.Info(snapStart, snapEnd+1, balance, rewardPerBlock, mdl.totalSupplies[dieselSym], mdl.rewards[pool][user])
				// update start
			}
		}
		snapStart = snapEnd + 1
	}
}

func (mdl *PoolLMRewards) addUserReward(pool, user string, reward *big.Int) {
	if mdl.rewards[pool] == nil {
		mdl.rewards[pool] = map[string]*big.Int{}
	}
	mdl.rewards[pool][user] = new(big.Int).Add(
		utils.NotNilBigInt(mdl.rewards[pool][user]),
		reward,
	)
}

// PoolLMRewards has fake address so no need for adding .Address value to addrs
func (mdl *PoolLMRewards) GetAllAddrsForLogs() (addrs []common.Address) {
	for addr, poolAndUToken := range mdl.Repo.GetDieselTokens() {
		addrs = append(addrs, common.HexToAddress(addr))
		token := mdl.Repo.GetToken(addr)
		mdl.decimalsAndPool[token.Symbol] = &_PoolAndDecimals{
			decimals: token.Decimals,
			pool:     poolAndUToken.Pool,
		}
	}
	return addrs
}
