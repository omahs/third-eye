package dc_wrapper

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"

	dcv2 "github.com/Gearbox-protocol/sdk-go/artifacts/dataCompressor/dataCompressorv2"
	"github.com/Gearbox-protocol/sdk-go/artifacts/dataCompressor/mainnet"
	"github.com/Gearbox-protocol/sdk-go/artifacts/multicall"
	"github.com/Gearbox-protocol/sdk-go/core"
	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/test"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type DataCompressorWrapper struct {
	mu *sync.RWMutex
	// blockNumbers of dc in asc order
	DCBlockNum         []int64
	BlockNumToName     map[int64]string
	discoveredAtToAddr map[int64]common.Address
	//
	v1DC    *MainnetDC
	testing *DCTesting

	client core.ClientI
}

var DCV1 = "DCV1"
var DCV2 = "DCV2"
var TESTING = "TESTING"

func NewDataCompressorWrapper(client core.ClientI) *DataCompressorWrapper {
	return &DataCompressorWrapper{
		mu:                 &sync.RWMutex{},
		BlockNumToName:     make(map[int64]string),
		discoveredAtToAddr: make(map[int64]common.Address),
		client:             client,
		v1DC:               NewMainnetDC(client),
		testing: &DCTesting{
			calls:  map[int64]*test.DCCalls{},
			client: client,
		},
	}
}

// testing
func (dcw *DataCompressorWrapper) SetCalls(blockNum int64, calls *test.DCCalls) {
	dcw.testing.calls[blockNum] = calls
}

func (dcw *DataCompressorWrapper) addDataCompressor(blockNum int64, addr string) {
	if len(dcw.DCBlockNum) > 0 && dcw.DCBlockNum[len(dcw.DCBlockNum)-1] >= blockNum {
		log.Fatalf("Current dc added at :%v, new dc:%s added at %d  ", dcw.DCBlockNum, addr, blockNum)
	}
	chainId, err := dcw.client.ChainID(context.TODO())
	log.CheckFatal(err)
	var key string
	if chainId.Int64() == 1 || chainId.Int64() == 5 { // or for goerli
		switch len(dcw.DCBlockNum) {
		case 0:
			key = DCV1
		case 1, 2: // 2 is for goerli added datacompressor v2, as datacompressor added second time
			key = DCV2
		}
	} else if chainId.Int64() == 42 {
		// 	switch len(dcw.DCBlockNum) {
		// for old address provider 0xA526311C39523F60b184709227875b5f34793bD4
		// we had a datacompressor which was used while first gearbox test deployment for users, that happened in nov 2021
		// later around july 2022 in redeployed whole kovan setup where there was only 1 dc per gearbox 1
		// 	case 0:
		// 		key = OLDKOVAN
		// 	case 1:
		// 		key = DCV1
		// 	case 2:
		// 		key = DCV2
		// }
		switch length := len(dcw.DCBlockNum); {
		case length == 0:
			key = DCV1
		case length < 9:
			key = DCV2
		}
	} else {
		key = TESTING
	}
	dcw.BlockNumToName[blockNum] = key
	dcw.discoveredAtToAddr[blockNum] = common.HexToAddress(addr)
	dcw.DCBlockNum = append(dcw.DCBlockNum, blockNum)
}

// the data compressor are added in increasing order of blockNum
func (dcw *DataCompressorWrapper) AddDataCompressor(blockNum int64, addr string) {
	dcw.mu.Lock()
	defer dcw.mu.Unlock()
	dcw.addDataCompressor(blockNum, addr)
}

func (dcw *DataCompressorWrapper) LoadMultipleDC(multiDCs interface{}) {
	dcw.mu.Lock()
	defer dcw.mu.Unlock()
	dcMap, ok := (multiDCs).(map[string]interface{})
	if !ok {
		log.Fatalf("Converting address provider() details for dc to map failed %v", multiDCs)
	}
	var blockNums []int64
	for k := range dcMap {
		blockNum, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		blockNums = append(blockNums, blockNum)
	}
	sort.Slice(blockNums, func(i, j int) bool { return blockNums[i] < blockNums[j] })
	for _, blockNum := range blockNums {
		k := fmt.Sprintf("%d", blockNum)
		dcAddr := dcMap[k]
		dcw.addDataCompressor(blockNum, dcAddr.(string))
	}
}

func (dcw *DataCompressorWrapper) GetCreditAccountData(blockNum int64, creditManager common.Address, borrower common.Address) (
	call multicall.Multicall2Call,
	resultFn func([]byte) (dcv2.CreditAccountData, error),
	errReturn error) {
	//
	key, discoveredAt := dcw.getDataCompressorIndex(blockNum)
	switch key {
	case DCV2:
		data, err := core.GetAbi("DataCompressorV2").Pack("getCreditAccountData", creditManager, borrower)
		call, errReturn = multicall.Multicall2Call{
			Target:   dcw.getDCAddr(discoveredAt),
			CallData: data,
		}, err
	case DCV1:
		data, err := core.GetAbi("DataCompressorMainnet").Pack("getCreditAccountDataExtended", creditManager, borrower)
		call, errReturn = multicall.Multicall2Call{
			Target:   dcw.getDCAddr(discoveredAt),
			CallData: data,
		}, err
	case TESTING:
		data, err := core.GetAbi("DataCompressorMainnet").Pack("getCreditAccountDataExtended", creditManager, borrower)
		call, errReturn = multicall.Multicall2Call{
			Target:   common.HexToAddress("0x0000000000000000000000000000000000000001"),
			CallData: data,
		}, err
	default:
		panic(fmt.Sprintf("data compressor number %s not found for credit account data extended", key))
	}
	resultFn = func(bytes []byte) (dcv2.CreditAccountData, error) {
		switch key {
		case DCV2:
			out, err := core.GetAbi("DataCompressorV2").Unpack("getCreditAccountData", bytes)
			if err != nil {
				return dcv2.CreditAccountData{}, err
			}
			accountData := *abi.ConvertType(out[0], new(dcv2.CreditAccountData)).(*dcv2.CreditAccountData)
			return accountData, nil
		case DCV1:
			out, err := core.GetAbi("DataCompressorMainnet").Unpack("getCreditAccountDataExtended", bytes)
			if err != nil {
				return dcv2.CreditAccountData{}, err
			}
			accountData := *abi.ConvertType(out[0], new(mainnet.DataTypesCreditAccountDataExtended)).(*mainnet.DataTypesCreditAccountDataExtended)
			return dcw.v1DC.getCreditAccountData(blockNum, accountData)
		case TESTING:
			return dcw.testing.getAccountData(blockNum, fmt.Sprintf("%s_%s", creditManager, borrower))
		}
		panic(fmt.Sprintf("data compressor number %s not found for pool data", key))
	}
	return
}

func (dcw *DataCompressorWrapper) GetCreditManagerData(blockNum int64, _creditManager common.Address) (
	call multicall.Multicall2Call,
	resultFn func([]byte) (dcv2.CreditManagerData, error),
	errReturn error) {
	//
	key, discoveredAt := dcw.getDataCompressorIndex(blockNum)
	switch key {
	case DCV2:
		data, err := core.GetAbi("DataCompressorV2").Pack("getCreditManagerData", _creditManager)
		call, errReturn = multicall.Multicall2Call{
			Target:   dcw.getDCAddr(discoveredAt),
			CallData: data,
		}, err
	case DCV1:
		data, err := core.GetAbi("DataCompressorMainnet").Pack("getCreditManagerData", _creditManager, _creditManager)
		call, errReturn = multicall.Multicall2Call{
			Target:   dcw.getDCAddr(discoveredAt),
			CallData: data,
		}, err
	case TESTING:
		data, err := core.GetAbi("DataCompressorMainnet").Pack("getCreditManagerData", _creditManager, _creditManager)
		call, errReturn = multicall.Multicall2Call{
			Target:   common.HexToAddress("0x0000000000000000000000000000000000000001"),
			CallData: data,
		}, err
	}
	//
	resultFn = func(bytes []byte) (dcv2.CreditManagerData, error) {
		switch key {
		case DCV2:
			out, err := core.GetAbi("DataCompressorV2").Unpack("getCreditManagerData", bytes)
			if err != nil {
				return dcv2.CreditManagerData{}, err
			}
			cmData := *abi.ConvertType(out[0], new(dcv2.CreditManagerData)).(*dcv2.CreditManagerData)
			return cmData, nil
		case DCV1:
			out, err := core.GetAbi("DataCompressorMainnet").Unpack("getCreditManagerData", bytes)
			if err != nil {
				return dcv2.CreditManagerData{}, err
			}
			cmData := *abi.ConvertType(out[0], new(mainnet.DataTypesCreditManagerData)).(*mainnet.DataTypesCreditManagerData)
			return getCMDataV1(cmData), nil
		case TESTING:
			return dcw.testing.getCMData(blockNum, _creditManager.Hex())
		}
		panic(fmt.Sprintf("data compressor number %s not found for pool data", key))
	}
	return
}

func (dcw *DataCompressorWrapper) GetPoolData(blockNum int64, _pool common.Address) (
	call multicall.Multicall2Call,
	resultFn func([]byte) (dcv2.PoolData, error),
	errReturn error) {
	//
	key, discoveredAt := dcw.getDataCompressorIndex(blockNum)
	switch key {
	case DCV2:
		data, err := core.GetAbi("DataCompressorV2").Pack("getPoolData", _pool)
		call, errReturn = multicall.Multicall2Call{
			Target:   dcw.getDCAddr(discoveredAt),
			CallData: data,
		}, err
	case DCV1:
		data, err := core.GetAbi("DataCompressorMainnet").Pack("getPoolData", _pool)
		call, errReturn = multicall.Multicall2Call{
			Target:   dcw.getDCAddr(discoveredAt),
			CallData: data,
		}, err
	case TESTING:
		data, err := core.GetAbi("DataCompressorMainnet").Pack("getPoolData", _pool)
		call, errReturn = multicall.Multicall2Call{
			Target:   common.HexToAddress("0x0000000000000000000000000000000000000001"),
			CallData: data,
		}, err
	default:
		panic(fmt.Sprintf("data compressor number %s not found for pool data", key))
	}
	//
	resultFn = func(bytes []byte) (dcv2.PoolData, error) {
		switch key {
		case DCV2:
			out, err := core.GetAbi("DataCompressorV2").Unpack("getPoolData", bytes)
			if err != nil {
				return dcv2.PoolData{}, err
			}
			poolData := *abi.ConvertType(out[0], new(dcv2.PoolData)).(*dcv2.PoolData)
			return poolData, nil
		case DCV1:
			out, err := core.GetAbi("DataCompressorMainnet").Unpack("getPoolData", bytes)
			if err != nil {
				return dcv2.PoolData{}, err
			}
			poolData := *abi.ConvertType(out[0], new(mainnet.DataTypesPoolData)).(*mainnet.DataTypesPoolData)
			return getPoolDataV1(poolData), nil
		case TESTING:
			return dcw.testing.getPoolData(blockNum, _pool.Hex())
		}
		panic(fmt.Sprintf("data compressor number %s not found for pool data", key))
	}
	return
}

// get the last datacompressor added before blockNum
// blockNum to name
func (dcw *DataCompressorWrapper) getDataCompressorIndex(blockNum int64) (name string, discoveredAt int64) {
	dcw.mu.Lock()
	defer dcw.mu.Unlock()
	for _, num := range dcw.DCBlockNum {
		// dc should be deployed before it is queried
		if num < blockNum {
			name = dcw.BlockNumToName[num]
			discoveredAt = num
		} else {
			break
		}
	}
	return
}

func (dcw *DataCompressorWrapper) AddCreditManagerToFilter(cmAddr, cfAddr string) {
	dcw.mu.Lock()
	defer dcw.mu.Unlock()
	dcw.v1DC.AddCreditManagerToFilter(cmAddr, cfAddr)
}

func (dcw *DataCompressorWrapper) getDCAddr(discoveredAt int64) common.Address {
	dcw.mu.RLock()
	defer dcw.mu.RUnlock()
	return dcw.discoveredAtToAddr[discoveredAt]
}

func (dcw *DataCompressorWrapper) ToJson() string {
	return utils.ToJson(dcw)
}

// func (dcw *DataCompressorWrapper) GetCreditAccountDataForHack(opts *bind.CallOpts, creditManager common.Address, borrower common.Address) (*dcv2.CreditAccountData, error) {
// 	if opts == nil || opts.BlockNumber == nil {
// 		panic("opts or blockNumber is nil")
// 	}
// 	key, discoveredAt := dcw.getDataCompressorIndex(opts.BlockNumber.Int64())
// 	switch key {
// 	case MAINNET:
// 		dcw.setMainnet(discoveredAt)
// 		data, err := dcw.dcMainnet.GetCreditAccountData(opts, creditManager, borrower)
// 		if err != nil {
// 			return nil, err
// 		}
// 		account, err := creditAccount.NewCreditAccount(data.Addr, dcw.client)
// 		if err != nil {
// 			return nil, err
// 		}
// 		cumIndex, err := account.CumulativeIndexAtOpen(opts)
// 		if err != nil {
// 			return nil, err
// 		}
// 		borrowedAmount, err := account.BorrowedAmount(opts)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return &dcv2.CreditAccountData{
// 			Addr:                       data.Addr,
// 			Borrower:                   data.Borrower,
// 			InUse:                      data.InUse,
// 			CreditManager:              data.CreditManager,
// 			Underlying:                 data.Underlying,
// 			BorrowedAmountPlusInterest: data.BorrowedAmountPlusInterest,
// 			TotalValue:                 data.TotalValue,
// 			HealthFactor:               data.HealthFactor,
// 			BorrowRate:                 data.BorrowRate,

// 			RepayAmount:           borrowedAmount,
// 			LiquidationAmount:     borrowedAmount,
// 			CanBeClosed:           false,
// 			BorrowedAmount:        borrowedAmount,
// 			CumulativeIndexAtOpen: cumIndex,
// 			Since:                 new(big.Int),
// 			Balances:              data.Balances,
// 		}, nil
// 	}
// 	panic(fmt.Sprintf("data compressor number %s not found for credit account data", key))
// }
