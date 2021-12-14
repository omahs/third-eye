package core

import (
	"fmt"
	"github.com/Gearbox-protocol/third-eye/artifacts/dataCompressor"
	"github.com/Gearbox-protocol/third-eye/artifacts/dataCompressor/mainnet"
	"github.com/Gearbox-protocol/third-eye/ethclient"
	"github.com/Gearbox-protocol/third-eye/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"sort"
	"context"
)

type DataCompressorWrapper struct {
	// blockNumbers of dc in asc order
	dcBlockNum     []int64
	blockNumToName map[int64]string
	dcOldKovan     *dataCompressor.DataCompressor
	dcMainnet      *mainnet.DataCompressor
	nameToAddr     map[string]string
	client         *ethclient.Client
}

var OLDKOVAN = "OLDKOVAN"
var MAINNET = "MAINNET"

func NewDataCompressorWrapper(client *ethclient.Client) *DataCompressorWrapper {
	return &DataCompressorWrapper{
		blockNumToName: make(map[int64]string),
		nameToAddr:     make(map[string]string),
		client:         client,
	}
}

func (dcw *DataCompressorWrapper) getDataCompressorIndex(blockNum int64) string {
	var name string
	for _, num := range dcw.dcBlockNum {
		// dc should be deployed before it is queried
		if num < blockNum {
			name = dcw.blockNumToName[num]
		} else {
			break
		}
	}
	return name
}

func (dcw *DataCompressorWrapper) AddDataCompressor(blockNum int64, addr string) {
	chainId, err := dcw.client.ChainID(context.TODO())
	log.CheckFatal(err)
	var key string
	if chainId.Int64() == 1 {
		key = MAINNET
	} else {
		switch len(dcw.dcBlockNum) {
		case 0:
			key = OLDKOVAN
		case 1:
			key = MAINNET
		}
	}
	dcw.blockNumToName[blockNum] = key
	dcw.nameToAddr[key] = addr
	dcw.dcBlockNum = append(dcw.dcBlockNum, blockNum)
	arr := dcw.dcBlockNum
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	dcw.dcBlockNum = arr
}

func (dcw *DataCompressorWrapper) GetCreditAccountDataExtended(opts *bind.CallOpts, creditManager common.Address, borrower common.Address) (mainnet.DataTypesCreditAccountDataExtended, error) {
	if opts == nil || opts.BlockNumber == nil {
		panic("opts or blockNumber is nil")
	}
	key := dcw.getDataCompressorIndex(opts.BlockNumber.Int64())
	switch key {
	case OLDKOVAN:
		dcw.setOldKovan()
		data, err := dcw.dcOldKovan.GetCreditAccountDataExtended(opts, creditManager, borrower)
		log.CheckFatal(err)
		latestFormat := mainnet.DataTypesCreditAccountDataExtended{
			Addr:                       data.Addr,
			Borrower:                   data.Borrower,
			InUse:                      data.InUse,
			CreditManager:              data.CreditManager,
			UnderlyingToken:            data.UnderlyingToken,
			BorrowedAmountPlusInterest: data.BorrowedAmountPlusInterest,
			TotalValue:                 data.TotalValue,
			HealthFactor:               data.HealthFactor,
			BorrowRate:                 data.BorrowRate,

			RepayAmount:           data.RepayAmount,
			LiquidationAmount:     data.LiquidationAmount,
			CanBeClosed:           data.CanBeClosed,
			BorrowedAmount:        data.BorrowedAmount,
			CumulativeIndexAtOpen: data.CumulativeIndexAtOpen,
			Since:                 data.Since,
		}
		for _, balance := range data.Balances {
			latestFormat.Balances = append(latestFormat.Balances, mainnet.DataTypesTokenBalance{
				Token:   balance.Token,
				Balance: balance.Balance,
			})
		}
		return latestFormat, err
	case MAINNET:
		dcw.setMainnet()
		return dcw.dcMainnet.GetCreditAccountDataExtended(opts, creditManager, borrower)
	}
	panic(fmt.Sprintf("data compressor number %s not found for credit account data", key))
}

func (dcw *DataCompressorWrapper) GetCreditManagerData(opts *bind.CallOpts, _creditManager common.Address, borrower common.Address) (mainnet.DataTypesCreditManagerData, error) {
	if opts == nil || opts.BlockNumber == nil {
		panic("opts or blockNumber is nil")
	}
	key := dcw.getDataCompressorIndex(opts.BlockNumber.Int64())
	switch key {
	case OLDKOVAN:
		dcw.setOldKovan()
		data, err := dcw.dcOldKovan.GetCreditManagerData(opts, _creditManager, borrower)
		log.CheckFatal(err)
		latestFormat := mainnet.DataTypesCreditManagerData{
			Addr:               data.Addr,
			HasAccount:         data.HasAccount,
			UnderlyingToken:    data.UnderlyingToken,
			IsWETH:             data.IsWETH,
			CanBorrow:          data.CanBorrow,
			BorrowRate:         data.BorrowRate,
			MinAmount:          data.MinAmount,
			MaxAmount:          data.MaxAmount,
			MaxLeverageFactor:  data.MaxLeverageFactor,
			AvailableLiquidity: data.AvailableLiquidity,
			AllowedTokens:      data.AllowedTokens,
		}
		for _, adapter := range data.Adapters {
			latestFormat.Adapters = append(latestFormat.Adapters, mainnet.DataTypesContractAdapter{
				Adapter:         adapter.Adapter,
				AllowedContract: adapter.AllowedContract,
			})
		}
		return latestFormat, err
	case MAINNET:
		dcw.setMainnet()
		return dcw.dcMainnet.GetCreditManagerData(opts, _creditManager, borrower)
	}
	panic(fmt.Sprintf("data compressor number %s not found for credit manager data", key))
}

func (dcw *DataCompressorWrapper) GetPoolData(opts *bind.CallOpts, _pool common.Address) (mainnet.DataTypesPoolData, error) {
	if opts == nil || opts.BlockNumber == nil {
		panic("opts or blockNumber is nil")
	}
	key := dcw.getDataCompressorIndex(opts.BlockNumber.Int64())
	switch key {
	case OLDKOVAN:
		dcw.setOldKovan()
		data, err := dcw.dcOldKovan.GetPoolData(opts, _pool)
		log.CheckFatal(err)
		latestFormat := mainnet.DataTypesPoolData{
			Addr:                   data.Addr,
			IsWETH:                 data.IsWETH,
			UnderlyingToken:        data.UnderlyingToken,
			DieselToken:            data.DieselToken,
			LinearCumulativeIndex:  data.LinearCumulativeIndex,
			AvailableLiquidity:     data.AvailableLiquidity,
			ExpectedLiquidity:      data.ExpectedLiquidity,
			ExpectedLiquidityLimit: data.ExpectedLiquidityLimit,
			TotalBorrowed:          data.TotalBorrowed,
			DepositAPYRAY:          data.DepositAPYRAY,
			BorrowAPYRAY:           data.BorrowAPYRAY,
			DieselRateRAY:          data.DieselRateRAY,
			WithdrawFee:            data.WithdrawFee,
			CumulativeIndexRAY:     data.CumulativeIndexRAY,
			TimestampLU:            data.TimestampLU,
		}
		return latestFormat, err
	case MAINNET:
		dcw.setMainnet()
		return dcw.dcMainnet.GetPoolData(opts, _pool)
	}
	panic(fmt.Sprintf("data compressor number %s not found for pool data", key))
}

func (dcw *DataCompressorWrapper) setOldKovan() {
	if dcw.dcOldKovan == nil {
		addr := dcw.nameToAddr[OLDKOVAN]
		var err error
		dcw.dcOldKovan, err = dataCompressor.NewDataCompressor(common.HexToAddress(addr), dcw.client)
		log.CheckFatal(err)
	}
}

func (dcw *DataCompressorWrapper) setMainnet() {
	if dcw.dcMainnet == nil {
		addr := dcw.nameToAddr[MAINNET]
		var err error
		dcw.dcMainnet, err = mainnet.NewDataCompressor(common.HexToAddress(addr), dcw.client)
		log.CheckFatal(err)
	}
}
