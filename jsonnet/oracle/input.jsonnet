local utils = import '../utils.libsonnet';
local borrowedAmount = 4000;
local extraBorrowedAmount = 1000;
// test for disabling previous oracle for token and create new oracle
// test for changing the previous chainlink
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
      block: 4,
      feed: '#ChainlinkPriceFeed_3',
    }, {
      oracle: '#Oracle_2',
      block: 4,
      feed: '#ChainlinkPriceFeed_4',
    }],
  },
  blocks: {
    '3': {
      // price on usdc and yfi
      events: [
        {
          // price chainlink on usdc
          address: '#ChainlinkPriceFeed_1',
          txHash: '!#Hash_1',
          topics: [
            'AnswerUpdated(int256,uint256,uint256)',
            // 0.0004
            utils.bigIntTopic(0.0004, 18),
            // roundid
            utils.bigIntTopic(1, 0),
          ],
          data: [],
        },
        {
          // price chainlink on yfi
          address: '#ChainlinkPriceFeed_2',
          txHash: '!#Hash_2',
          topics: [
            'AnswerUpdated(int256,uint256,uint256)',
            // 8
            utils.bigIntTopic(8, 18),
            // roundid
            utils.bigIntTopic(1, 0),
          ],
          data: [],
        },
      ],
    },
    '5': {
      // new oracle on usdc
      events: [
        {
          address: '#PriceOracle_1',
          topics: [
            'NewPriceFeed(address,address)',
            '#Token_1',
            '#Oracle_3',
          ],
          txHash: '!#Hash_3',
        },
        // yfi
        {
          // price chainlink on yfi
          address: '#ChainlinkPriceFeed_4',
          txHash: '!#Hash_4',
          topics: [
            'AnswerUpdated(int256,uint256,uint256)',
            // 8
            utils.bigIntTopic(10, 18),
            // roundid
            utils.bigIntTopic(1, 0),
          ],
          data: [],
        },
      ],
    },
  },
}
