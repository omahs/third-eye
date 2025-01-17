package aggregated_block_feed

import (
	"fmt"
	"math/big"

	"github.com/Gearbox-protocol/sdk-go/artifacts/priceFeed"
	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/Gearbox-protocol/sdk-go/core/schemas"
	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type yearnPFInternal struct {
	mainPFAddress        common.Address // for yearn manual price calculation
	yVaultAddr           common.Address //for yearn manual price calculation
	underlyingPFContract *priceFeed.PriceFeed
	decimalDivider       *big.Int
}

func (mdl *yearnPFInternal) calculatePrice(blockNum int64, client core.ClientI, version core.VersionType) (*schemas.PriceFeed, error) {
	if mdl.underlyingPFContract == nil {
		if err := mdl.setContracts(blockNum, client); err != nil {
			return nil, err
		}
	}
	opts := &bind.CallOpts{
		BlockNumber: big.NewInt(blockNum),
	}
	//
	pps, err := core.CallFuncWithExtraBytes(client, "99530b06", mdl.yVaultAddr, blockNum, nil) // pps
	if err != nil {
		return nil, err
	}
	pricePerShare := new(big.Int).SetBytes(pps)
	//
	roundData, err := mdl.underlyingPFContract.LatestRoundData(opts)
	if err != nil {
		return nil, err
	}

	// for yearn it is based on the vault. https://github.com/Gearbox-protocol/integrations-v2/blob/main/contracts/oracles/yearn/YearnPriceFeed.sol#L62
	newAnswer := new(big.Int).Quo(
		new(big.Int).Mul(pricePerShare, roundData.Answer),
		mdl.decimalDivider,
	)
	return &schemas.PriceFeed{
		RoundId:      roundData.RoundId.Int64(),
		PriceBI:      (*core.BigInt)(newAnswer),
		Price:        utils.GetFloat64Decimal(newAnswer, version.Decimals()),
		IsPriceInUSD: version.IsPriceInUSD(),
	}, nil
}

func (mdl *yearnPFInternal) setContracts(blockNum int64, client core.ClientI) error {
	// set the price feed contract
	underlyingPFAddrBytes, err := core.CallFuncWithExtraBytes(client, "741bef1a", mdl.mainPFAddress, blockNum, nil) // priceFeed
	if err != nil {
		return err
	}
	// underlying price feed not found
	if common.BytesToAddress(underlyingPFAddrBytes) == core.NULL_ADDR {
		return fmt.Errorf("address for underlying pf for yearn feed(%d) not found at %d",
			mdl.mainPFAddress, blockNum)
	}
	mdl.underlyingPFContract, err = priceFeed.NewPriceFeed(common.BytesToAddress(underlyingPFAddrBytes), client)
	log.CheckFatal(err)

	// set the yvault contract
	yVaultAddrBytes, err := core.CallFuncWithExtraBytes(client, "33303f8e", mdl.mainPFAddress, blockNum, nil) // yVault
	if err != nil {
		return err
	}
	mdl.yVaultAddr = common.BytesToAddress(yVaultAddrBytes)
	//

	// set the decimals
	decimalsBytes, err := core.CallFuncWithExtraBytes(client, "313ce567", mdl.yVaultAddr, blockNum, nil) // decimals
	if err != nil {
		return err
	}
	mdl.decimalDivider = utils.GetExpInt(int8(
		new(big.Int).SetBytes(decimalsBytes).Int64(),
	))
	//
	return nil
}
