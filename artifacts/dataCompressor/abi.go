// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dataCompressor

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// DataTypesContractAdapter is an auto generated low-level Go binding around an user-defined struct.
type DataTypesContractAdapter struct {
	AllowedContract common.Address
	Adapter         common.Address
}

// DataTypesCreditAccountData is an auto generated low-level Go binding around an user-defined struct.
type DataTypesCreditAccountData struct {
	Addr                       common.Address
	Borrower                   common.Address
	InUse                      bool
	CreditManager              common.Address
	UnderlyingToken            common.Address
	BorrowedAmountPlusInterest *big.Int
	TotalValue                 *big.Int
	HealthFactor               *big.Int
	BorrowRate                 *big.Int
	Balances                   []DataTypesTokenBalance
}

// DataTypesCreditAccountDataExtended is an auto generated low-level Go binding around an user-defined struct.
type DataTypesCreditAccountDataExtended struct {
	Addr                       common.Address
	Borrower                   common.Address
	InUse                      bool
	CreditManager              common.Address
	UnderlyingToken            common.Address
	BorrowedAmountPlusInterest *big.Int
	TotalValue                 *big.Int
	HealthFactor               *big.Int
	BorrowRate                 *big.Int
	Balances                   []DataTypesTokenBalance
	RepayAmount                *big.Int
	LiquidationAmount          *big.Int
	CanBeClosed                bool
	BorrowedAmount             *big.Int
	CumulativeIndexAtOpen      *big.Int
	Since                      *big.Int
}

// DataTypesCreditManagerData is an auto generated low-level Go binding around an user-defined struct.
type DataTypesCreditManagerData struct {
	Addr               common.Address
	HasAccount         bool
	UnderlyingToken    common.Address
	IsWETH             bool
	CanBorrow          bool
	BorrowRate         *big.Int
	MinAmount          *big.Int
	MaxAmount          *big.Int
	MaxLeverageFactor  *big.Int
	AvailableLiquidity *big.Int
	AllowedTokens      []common.Address
	Adapters           []DataTypesContractAdapter
}

// DataTypesPoolData is an auto generated low-level Go binding around an user-defined struct.
type DataTypesPoolData struct {
	Addr                   common.Address
	IsWETH                 bool
	UnderlyingToken        common.Address
	DieselToken            common.Address
	LinearCumulativeIndex  *big.Int
	AvailableLiquidity     *big.Int
	ExpectedLiquidity      *big.Int
	ExpectedLiquidityLimit *big.Int
	TotalBorrowed          *big.Int
	DepositAPYRAY          *big.Int
	BorrowAPYRAY           *big.Int
	DieselRateRAY          *big.Int
	WithdrawFee            *big.Int
	CumulativeIndexRAY     *big.Int
	TimestampLU            *big.Int
}

// DataTypesTokenBalance is an auto generated low-level Go binding around an user-defined struct.
type DataTypesTokenBalance struct {
	Token   common.Address
	Balance *big.Int
}

// DataTypesTokenInfo is an auto generated low-level Go binding around an user-defined struct.
type DataTypesTokenInfo struct {
	Addr     common.Address
	Symbol   string
	Decimals uint8
}

// DataCompressorABI is the input ABI used to generate the binding from.
const DataCompressorABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"creditManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"}],\"name\":\"calcExpectedHf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_creditManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"}],\"name\":\"getCreditAccountData\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"inUse\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"creditManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"underlyingToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"borrowedAmountPlusInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"healthFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.TokenBalance[]\",\"name\":\"balances\",\"type\":\"tuple[]\"}],\"internalType\":\"structDataTypes.CreditAccountData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"creditManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"}],\"name\":\"getCreditAccountDataExtended\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"inUse\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"creditManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"underlyingToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"borrowedAmountPlusInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"healthFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.TokenBalance[]\",\"name\":\"balances\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"repayAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"canBeClosed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"borrowedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cumulativeIndexAtOpen\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"since\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.CreditAccountDataExtended\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"}],\"name\":\"getCreditAccountList\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"inUse\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"creditManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"underlyingToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"borrowedAmountPlusInterest\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalValue\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"healthFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.TokenBalance[]\",\"name\":\"balances\",\"type\":\"tuple[]\"}],\"internalType\":\"structDataTypes.CreditAccountData[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_creditManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"}],\"name\":\"getCreditManagerData\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"hasAccount\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"underlyingToken\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isWETH\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"canBorrow\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxLeverageFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"availableLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"allowedTokens\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"allowedContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"adapter\",\"type\":\"address\"}],\"internalType\":\"structDataTypes.ContractAdapter[]\",\"name\":\"adapters\",\"type\":\"tuple[]\"}],\"internalType\":\"structDataTypes.CreditManagerData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"}],\"name\":\"getCreditManagersList\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"hasAccount\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"underlyingToken\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isWETH\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"canBorrow\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"borrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxLeverageFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"availableLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"allowedTokens\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"allowedContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"adapter\",\"type\":\"address\"}],\"internalType\":\"structDataTypes.ContractAdapter[]\",\"name\":\"adapters\",\"type\":\"tuple[]\"}],\"internalType\":\"structDataTypes.CreditManagerData[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_pool\",\"type\":\"address\"}],\"name\":\"getPoolData\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isWETH\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"underlyingToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dieselToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"linearCumulativeIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"availableLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedLiquidityLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalBorrowed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAPY_RAY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowAPY_RAY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dieselRate_RAY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cumulativeIndex_RAY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestampLU\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.PoolData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPoolsList\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isWETH\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"underlyingToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dieselToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"linearCumulativeIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"availableLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expectedLiquidityLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalBorrowed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"depositAPY_RAY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowAPY_RAY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"dieselRate_RAY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"cumulativeIndex_RAY\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestampLU\",\"type\":\"uint256\"}],\"internalType\":\"structDataTypes.PoolData[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addr\",\"type\":\"address[]\"}],\"name\":\"getTokenData\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"internalType\":\"structDataTypes.TokenInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"creditManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"}],\"name\":\"hasOpenedCreditAccount\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// DataCompressor is an auto generated Go binding around an Ethereum contract.
type DataCompressor struct {
	DataCompressorCaller     // Read-only binding to the contract
	DataCompressorTransactor // Write-only binding to the contract
	DataCompressorFilterer   // Log filterer for contract events
}

// DataCompressorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DataCompressorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataCompressorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DataCompressorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataCompressorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DataCompressorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataCompressorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DataCompressorSession struct {
	Contract     *DataCompressor   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DataCompressorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DataCompressorCallerSession struct {
	Contract *DataCompressorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// DataCompressorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DataCompressorTransactorSession struct {
	Contract     *DataCompressorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// DataCompressorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DataCompressorRaw struct {
	Contract *DataCompressor // Generic contract binding to access the raw methods on
}

// DataCompressorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DataCompressorCallerRaw struct {
	Contract *DataCompressorCaller // Generic read-only contract binding to access the raw methods on
}

// DataCompressorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DataCompressorTransactorRaw struct {
	Contract *DataCompressorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDataCompressor creates a new instance of DataCompressor, bound to a specific deployed contract.
func NewDataCompressor(address common.Address, backend bind.ContractBackend) (*DataCompressor, error) {
	contract, err := bindDataCompressor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DataCompressor{DataCompressorCaller: DataCompressorCaller{contract: contract}, DataCompressorTransactor: DataCompressorTransactor{contract: contract}, DataCompressorFilterer: DataCompressorFilterer{contract: contract}}, nil
}

// NewDataCompressorCaller creates a new read-only instance of DataCompressor, bound to a specific deployed contract.
func NewDataCompressorCaller(address common.Address, caller bind.ContractCaller) (*DataCompressorCaller, error) {
	contract, err := bindDataCompressor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DataCompressorCaller{contract: contract}, nil
}

// NewDataCompressorTransactor creates a new write-only instance of DataCompressor, bound to a specific deployed contract.
func NewDataCompressorTransactor(address common.Address, transactor bind.ContractTransactor) (*DataCompressorTransactor, error) {
	contract, err := bindDataCompressor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DataCompressorTransactor{contract: contract}, nil
}

// NewDataCompressorFilterer creates a new log filterer instance of DataCompressor, bound to a specific deployed contract.
func NewDataCompressorFilterer(address common.Address, filterer bind.ContractFilterer) (*DataCompressorFilterer, error) {
	contract, err := bindDataCompressor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DataCompressorFilterer{contract: contract}, nil
}

// bindDataCompressor binds a generic wrapper to an already deployed contract.
func bindDataCompressor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DataCompressorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataCompressor *DataCompressorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataCompressor.Contract.DataCompressorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataCompressor *DataCompressorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataCompressor.Contract.DataCompressorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataCompressor *DataCompressorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataCompressor.Contract.DataCompressorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataCompressor *DataCompressorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataCompressor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataCompressor *DataCompressorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataCompressor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataCompressor *DataCompressorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataCompressor.Contract.contract.Transact(opts, method, params...)
}

// CalcExpectedHf is a free data retrieval call binding the contract method 0xba3b7345.
//
// Solidity: function calcExpectedHf(address creditManager, address borrower, uint256[] balances) view returns(uint256)
func (_DataCompressor *DataCompressorCaller) CalcExpectedHf(opts *bind.CallOpts, creditManager common.Address, borrower common.Address, balances []*big.Int) (*big.Int, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "calcExpectedHf", creditManager, borrower, balances)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalcExpectedHf is a free data retrieval call binding the contract method 0xba3b7345.
//
// Solidity: function calcExpectedHf(address creditManager, address borrower, uint256[] balances) view returns(uint256)
func (_DataCompressor *DataCompressorSession) CalcExpectedHf(creditManager common.Address, borrower common.Address, balances []*big.Int) (*big.Int, error) {
	return _DataCompressor.Contract.CalcExpectedHf(&_DataCompressor.CallOpts, creditManager, borrower, balances)
}

// CalcExpectedHf is a free data retrieval call binding the contract method 0xba3b7345.
//
// Solidity: function calcExpectedHf(address creditManager, address borrower, uint256[] balances) view returns(uint256)
func (_DataCompressor *DataCompressorCallerSession) CalcExpectedHf(creditManager common.Address, borrower common.Address, balances []*big.Int) (*big.Int, error) {
	return _DataCompressor.Contract.CalcExpectedHf(&_DataCompressor.CallOpts, creditManager, borrower, balances)
}

// GetCreditAccountData is a free data retrieval call binding the contract method 0x0dbd616d.
//
// Solidity: function getCreditAccountData(address _creditManager, address borrower) view returns((address,address,bool,address,address,uint256,uint256,uint256,uint256,(address,uint256)[]))
func (_DataCompressor *DataCompressorCaller) GetCreditAccountData(opts *bind.CallOpts, _creditManager common.Address, borrower common.Address) (DataTypesCreditAccountData, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "getCreditAccountData", _creditManager, borrower)

	if err != nil {
		return *new(DataTypesCreditAccountData), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesCreditAccountData)).(*DataTypesCreditAccountData)

	return out0, err

}

// GetCreditAccountData is a free data retrieval call binding the contract method 0x0dbd616d.
//
// Solidity: function getCreditAccountData(address _creditManager, address borrower) view returns((address,address,bool,address,address,uint256,uint256,uint256,uint256,(address,uint256)[]))
func (_DataCompressor *DataCompressorSession) GetCreditAccountData(_creditManager common.Address, borrower common.Address) (DataTypesCreditAccountData, error) {
	return _DataCompressor.Contract.GetCreditAccountData(&_DataCompressor.CallOpts, _creditManager, borrower)
}

// GetCreditAccountData is a free data retrieval call binding the contract method 0x0dbd616d.
//
// Solidity: function getCreditAccountData(address _creditManager, address borrower) view returns((address,address,bool,address,address,uint256,uint256,uint256,uint256,(address,uint256)[]))
func (_DataCompressor *DataCompressorCallerSession) GetCreditAccountData(_creditManager common.Address, borrower common.Address) (DataTypesCreditAccountData, error) {
	return _DataCompressor.Contract.GetCreditAccountData(&_DataCompressor.CallOpts, _creditManager, borrower)
}

// GetCreditAccountDataExtended is a free data retrieval call binding the contract method 0x191482d4.
//
// Solidity: function getCreditAccountDataExtended(address creditManager, address borrower) view returns((address,address,bool,address,address,uint256,uint256,uint256,uint256,(address,uint256)[],uint256,uint256,bool,uint256,uint256,uint256))
func (_DataCompressor *DataCompressorCaller) GetCreditAccountDataExtended(opts *bind.CallOpts, creditManager common.Address, borrower common.Address) (DataTypesCreditAccountDataExtended, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "getCreditAccountDataExtended", creditManager, borrower)

	if err != nil {
		return *new(DataTypesCreditAccountDataExtended), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesCreditAccountDataExtended)).(*DataTypesCreditAccountDataExtended)

	return out0, err

}

// GetCreditAccountDataExtended is a free data retrieval call binding the contract method 0x191482d4.
//
// Solidity: function getCreditAccountDataExtended(address creditManager, address borrower) view returns((address,address,bool,address,address,uint256,uint256,uint256,uint256,(address,uint256)[],uint256,uint256,bool,uint256,uint256,uint256))
func (_DataCompressor *DataCompressorSession) GetCreditAccountDataExtended(creditManager common.Address, borrower common.Address) (DataTypesCreditAccountDataExtended, error) {
	return _DataCompressor.Contract.GetCreditAccountDataExtended(&_DataCompressor.CallOpts, creditManager, borrower)
}

// GetCreditAccountDataExtended is a free data retrieval call binding the contract method 0x191482d4.
//
// Solidity: function getCreditAccountDataExtended(address creditManager, address borrower) view returns((address,address,bool,address,address,uint256,uint256,uint256,uint256,(address,uint256)[],uint256,uint256,bool,uint256,uint256,uint256))
func (_DataCompressor *DataCompressorCallerSession) GetCreditAccountDataExtended(creditManager common.Address, borrower common.Address) (DataTypesCreditAccountDataExtended, error) {
	return _DataCompressor.Contract.GetCreditAccountDataExtended(&_DataCompressor.CallOpts, creditManager, borrower)
}

// GetCreditAccountList is a free data retrieval call binding the contract method 0xa80deda3.
//
// Solidity: function getCreditAccountList(address borrower) view returns((address,address,bool,address,address,uint256,uint256,uint256,uint256,(address,uint256)[])[])
func (_DataCompressor *DataCompressorCaller) GetCreditAccountList(opts *bind.CallOpts, borrower common.Address) ([]DataTypesCreditAccountData, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "getCreditAccountList", borrower)

	if err != nil {
		return *new([]DataTypesCreditAccountData), err
	}

	out0 := *abi.ConvertType(out[0], new([]DataTypesCreditAccountData)).(*[]DataTypesCreditAccountData)

	return out0, err

}

// GetCreditAccountList is a free data retrieval call binding the contract method 0xa80deda3.
//
// Solidity: function getCreditAccountList(address borrower) view returns((address,address,bool,address,address,uint256,uint256,uint256,uint256,(address,uint256)[])[])
func (_DataCompressor *DataCompressorSession) GetCreditAccountList(borrower common.Address) ([]DataTypesCreditAccountData, error) {
	return _DataCompressor.Contract.GetCreditAccountList(&_DataCompressor.CallOpts, borrower)
}

// GetCreditAccountList is a free data retrieval call binding the contract method 0xa80deda3.
//
// Solidity: function getCreditAccountList(address borrower) view returns((address,address,bool,address,address,uint256,uint256,uint256,uint256,(address,uint256)[])[])
func (_DataCompressor *DataCompressorCallerSession) GetCreditAccountList(borrower common.Address) ([]DataTypesCreditAccountData, error) {
	return _DataCompressor.Contract.GetCreditAccountList(&_DataCompressor.CallOpts, borrower)
}

// GetCreditManagerData is a free data retrieval call binding the contract method 0xb10b074e.
//
// Solidity: function getCreditManagerData(address _creditManager, address borrower) view returns((address,bool,address,bool,bool,uint256,uint256,uint256,uint256,uint256,address[],(address,address)[]))
func (_DataCompressor *DataCompressorCaller) GetCreditManagerData(opts *bind.CallOpts, _creditManager common.Address, borrower common.Address) (DataTypesCreditManagerData, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "getCreditManagerData", _creditManager, borrower)

	if err != nil {
		return *new(DataTypesCreditManagerData), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesCreditManagerData)).(*DataTypesCreditManagerData)

	return out0, err

}

// GetCreditManagerData is a free data retrieval call binding the contract method 0xb10b074e.
//
// Solidity: function getCreditManagerData(address _creditManager, address borrower) view returns((address,bool,address,bool,bool,uint256,uint256,uint256,uint256,uint256,address[],(address,address)[]))
func (_DataCompressor *DataCompressorSession) GetCreditManagerData(_creditManager common.Address, borrower common.Address) (DataTypesCreditManagerData, error) {
	return _DataCompressor.Contract.GetCreditManagerData(&_DataCompressor.CallOpts, _creditManager, borrower)
}

// GetCreditManagerData is a free data retrieval call binding the contract method 0xb10b074e.
//
// Solidity: function getCreditManagerData(address _creditManager, address borrower) view returns((address,bool,address,bool,bool,uint256,uint256,uint256,uint256,uint256,address[],(address,address)[]))
func (_DataCompressor *DataCompressorCallerSession) GetCreditManagerData(_creditManager common.Address, borrower common.Address) (DataTypesCreditManagerData, error) {
	return _DataCompressor.Contract.GetCreditManagerData(&_DataCompressor.CallOpts, _creditManager, borrower)
}

// GetCreditManagersList is a free data retrieval call binding the contract method 0xb8169039.
//
// Solidity: function getCreditManagersList(address borrower) view returns((address,bool,address,bool,bool,uint256,uint256,uint256,uint256,uint256,address[],(address,address)[])[])
func (_DataCompressor *DataCompressorCaller) GetCreditManagersList(opts *bind.CallOpts, borrower common.Address) ([]DataTypesCreditManagerData, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "getCreditManagersList", borrower)

	if err != nil {
		return *new([]DataTypesCreditManagerData), err
	}

	out0 := *abi.ConvertType(out[0], new([]DataTypesCreditManagerData)).(*[]DataTypesCreditManagerData)

	return out0, err

}

// GetCreditManagersList is a free data retrieval call binding the contract method 0xb8169039.
//
// Solidity: function getCreditManagersList(address borrower) view returns((address,bool,address,bool,bool,uint256,uint256,uint256,uint256,uint256,address[],(address,address)[])[])
func (_DataCompressor *DataCompressorSession) GetCreditManagersList(borrower common.Address) ([]DataTypesCreditManagerData, error) {
	return _DataCompressor.Contract.GetCreditManagersList(&_DataCompressor.CallOpts, borrower)
}

// GetCreditManagersList is a free data retrieval call binding the contract method 0xb8169039.
//
// Solidity: function getCreditManagersList(address borrower) view returns((address,bool,address,bool,bool,uint256,uint256,uint256,uint256,uint256,address[],(address,address)[])[])
func (_DataCompressor *DataCompressorCallerSession) GetCreditManagersList(borrower common.Address) ([]DataTypesCreditManagerData, error) {
	return _DataCompressor.Contract.GetCreditManagersList(&_DataCompressor.CallOpts, borrower)
}

// GetPoolData is a free data retrieval call binding the contract method 0x13d21cdf.
//
// Solidity: function getPoolData(address _pool) view returns((address,bool,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_DataCompressor *DataCompressorCaller) GetPoolData(opts *bind.CallOpts, _pool common.Address) (DataTypesPoolData, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "getPoolData", _pool)

	if err != nil {
		return *new(DataTypesPoolData), err
	}

	out0 := *abi.ConvertType(out[0], new(DataTypesPoolData)).(*DataTypesPoolData)

	return out0, err

}

// GetPoolData is a free data retrieval call binding the contract method 0x13d21cdf.
//
// Solidity: function getPoolData(address _pool) view returns((address,bool,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_DataCompressor *DataCompressorSession) GetPoolData(_pool common.Address) (DataTypesPoolData, error) {
	return _DataCompressor.Contract.GetPoolData(&_DataCompressor.CallOpts, _pool)
}

// GetPoolData is a free data retrieval call binding the contract method 0x13d21cdf.
//
// Solidity: function getPoolData(address _pool) view returns((address,bool,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_DataCompressor *DataCompressorCallerSession) GetPoolData(_pool common.Address) (DataTypesPoolData, error) {
	return _DataCompressor.Contract.GetPoolData(&_DataCompressor.CallOpts, _pool)
}

// GetPoolsList is a free data retrieval call binding the contract method 0x1bcd8fc0.
//
// Solidity: function getPoolsList() view returns((address,bool,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_DataCompressor *DataCompressorCaller) GetPoolsList(opts *bind.CallOpts) ([]DataTypesPoolData, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "getPoolsList")

	if err != nil {
		return *new([]DataTypesPoolData), err
	}

	out0 := *abi.ConvertType(out[0], new([]DataTypesPoolData)).(*[]DataTypesPoolData)

	return out0, err

}

// GetPoolsList is a free data retrieval call binding the contract method 0x1bcd8fc0.
//
// Solidity: function getPoolsList() view returns((address,bool,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_DataCompressor *DataCompressorSession) GetPoolsList() ([]DataTypesPoolData, error) {
	return _DataCompressor.Contract.GetPoolsList(&_DataCompressor.CallOpts)
}

// GetPoolsList is a free data retrieval call binding the contract method 0x1bcd8fc0.
//
// Solidity: function getPoolsList() view returns((address,bool,address,address,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_DataCompressor *DataCompressorCallerSession) GetPoolsList() ([]DataTypesPoolData, error) {
	return _DataCompressor.Contract.GetPoolsList(&_DataCompressor.CallOpts)
}

// GetTokenData is a free data retrieval call binding the contract method 0xbf2eb19e.
//
// Solidity: function getTokenData(address[] addr) view returns((address,string,uint8)[])
func (_DataCompressor *DataCompressorCaller) GetTokenData(opts *bind.CallOpts, addr []common.Address) ([]DataTypesTokenInfo, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "getTokenData", addr)

	if err != nil {
		return *new([]DataTypesTokenInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]DataTypesTokenInfo)).(*[]DataTypesTokenInfo)

	return out0, err

}

// GetTokenData is a free data retrieval call binding the contract method 0xbf2eb19e.
//
// Solidity: function getTokenData(address[] addr) view returns((address,string,uint8)[])
func (_DataCompressor *DataCompressorSession) GetTokenData(addr []common.Address) ([]DataTypesTokenInfo, error) {
	return _DataCompressor.Contract.GetTokenData(&_DataCompressor.CallOpts, addr)
}

// GetTokenData is a free data retrieval call binding the contract method 0xbf2eb19e.
//
// Solidity: function getTokenData(address[] addr) view returns((address,string,uint8)[])
func (_DataCompressor *DataCompressorCallerSession) GetTokenData(addr []common.Address) ([]DataTypesTokenInfo, error) {
	return _DataCompressor.Contract.GetTokenData(&_DataCompressor.CallOpts, addr)
}

// HasOpenedCreditAccount is a free data retrieval call binding the contract method 0xfc9914cb.
//
// Solidity: function hasOpenedCreditAccount(address creditManager, address borrower) view returns(bool)
func (_DataCompressor *DataCompressorCaller) HasOpenedCreditAccount(opts *bind.CallOpts, creditManager common.Address, borrower common.Address) (bool, error) {
	var out []interface{}
	err := _DataCompressor.contract.Call(opts, &out, "hasOpenedCreditAccount", creditManager, borrower)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasOpenedCreditAccount is a free data retrieval call binding the contract method 0xfc9914cb.
//
// Solidity: function hasOpenedCreditAccount(address creditManager, address borrower) view returns(bool)
func (_DataCompressor *DataCompressorSession) HasOpenedCreditAccount(creditManager common.Address, borrower common.Address) (bool, error) {
	return _DataCompressor.Contract.HasOpenedCreditAccount(&_DataCompressor.CallOpts, creditManager, borrower)
}

// HasOpenedCreditAccount is a free data retrieval call binding the contract method 0xfc9914cb.
//
// Solidity: function hasOpenedCreditAccount(address creditManager, address borrower) view returns(bool)
func (_DataCompressor *DataCompressorCallerSession) HasOpenedCreditAccount(creditManager common.Address, borrower common.Address) (bool, error) {
	return _DataCompressor.Contract.HasOpenedCreditAccount(&_DataCompressor.CallOpts, creditManager, borrower)
}
