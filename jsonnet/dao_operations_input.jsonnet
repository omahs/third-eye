local bigInt(x, decimals=0) = std.toString(x * std.pow(10, decimals));
local bigIntTopic(x, decimals) = 'bigint:' + bigInt(x, decimals);
local initialAmount = 1000;
local borrowedAmount = 4000;
{
  mocks: {
    syncAdapters: 'mocks/syncAdapter1.json',
  },
  states: {
    oracles: [{
      oracle: '#Oracle_1',
      block: 1,
      feed: '#ChainlinkPriceFeed_1',
    }, {
      oracle: '#Oracle_2',
      block: 1,
      feed: '#ChainlinkPriceFeed_2',
    }, {
      oracle: '#Oracle_3',
      block: 1,
      feed: '#ChainlinkPriceFeed_3',
    }],
  },
  blocks: {
    '3': {
      events: [
        {
          // credit manager on usdc
          address: '#CreditFilter_1',
          topics: [
            'TokenAllowed(address,uint256)',
            '#Token_1',
          ],
          data: [
            'bigint:7500',
          ],
          txHash: '!#Hash_1',
        },
        {
          // credit manager on usdc
          address: '#CreditFilter_1',
          topics: [
            'TokenForbidden(address)',
            '#Token_1',
          ],
          txHash: '!#Hash_2',
        },
        {
          // credit manager on usdc
          address: '#CreditFilter_1',
          topics: [
            'ContractAllowed(address,address)',
            '#Protocol_1',
            '#Adapter_1',
          ],
          txHash: '!#Hash_3',
        },
        {
          // credit manager on usdc
          address: '#CreditFilter_1',
          topics: [
            'ContractForbidden(address)',
            '#Protocol_1',
          ],
          txHash: '!#Hash_4',
        },
        {
          // credit manager on usdc
          address: '#CreditFilter_1',
          topics: [
            'NewFastCheckParameters(uint256,uint256)',
          ],
          data: [
            'bigint:7500',
            'bigint:7500',
          ],
          txHash: '!#Hash_5',
        },
        {
          // credit manager on usdc
          address: '#CreditFilter_1',
          topics: [
            'TransferPluginAllowed(address,bool)',
            '#Plugin_1',
          ],
          data: [
            'bool:1',  // for true
          ],
          txHash: '!#Hash_6',
        },
        {
          // credit manager on usdc
          address: '#CreditFilter_1',
          topics: [
            'PriceOracleUpdated(address)',
            '#PriceOracle_2',
          ],
          txHash: '!#Hash_7',
        },
        {
          // credit manager on usdc
          address: '#Pool_1',
          topics: [
            'NewInterestRateModel(address)',
            '#IntereestRateModel_1',
          ],
          txHash: '!#Hash_8',
        },
        {
          // credit manager on usdc
          address: '#Pool_1',
          topics: [
            'NewCreditManagerConnected(address)',
            '#CreditManager_1',
          ],
          txHash: '!#Hash_9',
        },
        {
          // credit manager on usdc
          address: '#Pool_1',
          topics: [
            'NewExpectedLiquidityLimit(uint256)',
          ],
          data: [
            bigIntTopic(10000, 6),
          ],
          txHash: '!#Hash_10',
        },
        {
          // credit manager on usdc
          address: '#Pool_1',
          topics: [
            'BorrowForbidden(address)',
            '#CreditManager_2',
          ],
          txHash: '!#Hash_11',
        },
        {
          // credit manager on usdc
          address: '#Pool_1',
          topics: [
            'NewWithdrawFee(uint256)',
          ],
          data: [
            bigIntTopic(100, 0),  //1%
          ],
          txHash: '!#Hash_12',
        },
        {
          address: '#PriceOracle_1',
          topics: [
            'NewPriceFeed(address,address)',
            '#Token_4',
            '#Oracle_3',
          ],
          txHash: '!#Hash_13',
        },
        {
          address: '#AccountFactory_1',
          topics: [
            'TakeForever(address,address)',
            '#Account_10',
            '#To_1',
          ],
          txHash: '!#Hash_14',
        },
        {
          address: '#ACL_1',
          topics: [
            'PausableAdminAdded(address)',
            '#Admin_1',
          ],
          txHash: '!#Hash_15',
        },
        {
          address: '#ACL_1',
          topics: [
            'PausableAdminRemoved(address)',
            '#Admin_1',
          ],
          txHash: '!#Hash_16',
        },
        {
          address: '#ACL_1',
          topics: [
            'UnpausableAdminAdded(address)',
            '#Admin_1',
          ],
          txHash: '!#Hash_17',
        },
        {
          address: '#ACL_1',
          topics: [
            'UnpausableAdminRemoved(address)',
            '#Admin_1',
          ],
          txHash: '!#Hash_18',
        },
        {
          address: '#ACL_1',
          topics: [
            'OwnershipTransferred(address,address)',
            '#Owner_1',
            '#Admin_2',
          ],
          txHash: '!#Hash_19',
        },
        {
          address: '#ACL_1',
          topics: [
            'Paused()',
          ],
          txHash: '!#Hash_20',
        },
        {
          address: '#ACL_1',
          topics: [
            'UnPaused()',
          ],
          txHash: '!#Hash_21',
        },
        {
          // credit manager on usdc
          address: '#CreditManager_1',
          topics: [
            'NewParameters(uint256,uint256,uint256,uint256,uint256,uint256)',
          ],
          data: [
            // minAnount
            bigIntTopic(1000, 6),
            // maxAmount
            bigIntTopic(5000, 6),
            // maxLeverage
            bigIntTopic(400, 6),
            // feeInterest
            bigIntTopic(1000, 0),
            // feeLiquidation
            bigIntTopic(200, 0),
            // liquidationDiscount
            bigIntTopic(9500, 0),
          ],
          txHash: '!#Hash_22',
        },
      ],
      calls:
        {
          pools: [{
            address: '#Pool_1',
            totalBorrowed: bigInt(borrowedAmount, 6),
            expectedLiquidity: bigInt(borrowedAmount + 1000, 6),
            availableLiquidity: bigInt(1000, 6),
            depositAPY: bigInt(0),
            borrowAPY: bigInt(0),
            dieselRate: bigInt(0),
            withdrawFee: '0',
            linearCumulativeIndex: bigInt(1, 27),
          }],
          cms: [{
            address: '#CreditManager_1',
            isWETH: false,
            minAmount: bigInt(1000, 6),
            maxAmount: bigInt(5000, 6),
            availableLiquidity: bigInt(1000, 6),
            borrowRate: '0',
          }],
        },
    },
  },
}
