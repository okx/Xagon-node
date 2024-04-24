package l1_check_block_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block"
	mock_l1_check_block "github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block/mocks"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSafeBlockNumber(t *testing.T) {
	ctx := context.Background()
	mockRequester := mock_l1_check_block.NewL1Requester(t)
	//safeBlockPoint := big.NewInt(50)
	offset := 10
	safeL1Block := l1_check_block.NewSafeL1BlockNumberFetch(l1_check_block.StringToL1BlockPoint("safe"), offset)

	mockRequester.EXPECT().HeaderByNumber(ctx, mock.Anything).Return(&types.Header{
		Number: big.NewInt(100),
	}, nil)
	blockNumber, err := safeL1Block.GetSafeBlockNumber(ctx, mockRequester)
	assert.NoError(t, err)
	expectedBlockNumber := uint64(100 + offset)
	assert.Equal(t, expectedBlockNumber, blockNumber)
}

func TestGetSafeBlockNumberMutliplesCases(t *testing.T) {
	tests := []struct {
		name                string
		blockPoint          string
		offset              int
		l1ReturnBlockNumber uint64
		expectedCallToGeth  *big.Int
		expectedBlockNumber uint64
	}{
		{
			name:                "SafeBlockNumber+10",
			blockPoint:          "safe",
			offset:              10,
			l1ReturnBlockNumber: 100,
			expectedCallToGeth:  big.NewInt(int64(rpc.SafeBlockNumber)),
			expectedBlockNumber: 110,
		},
		{
			name:                "FinalizedBlockNumber+10",
			blockPoint:          "finalized",
			offset:              10,
			l1ReturnBlockNumber: 100,
			expectedCallToGeth:  big.NewInt(int64(rpc.FinalizedBlockNumber)),
			expectedBlockNumber: 110,
		},
		{
			name:                "PendingBlockNumber+10",
			blockPoint:          "pending",
			offset:              10,
			l1ReturnBlockNumber: 100,
			expectedCallToGeth:  big.NewInt(int64(rpc.PendingBlockNumber)),
			expectedBlockNumber: 110,
		},
		{
			name:                "LastBlockNumber+10, can't add 10 to latest block number. So must return latest block number and ignore positive offset",
			blockPoint:          "latest",
			offset:              10,
			l1ReturnBlockNumber: 100,
			expectedCallToGeth:  nil,
			expectedBlockNumber: 100,
		},
		{
			name:                "FinalizedBlockNumber-1000. negative blockNumbers are not welcome. So must return 0",
			blockPoint:          "finalized",
			offset:              -1000,
			l1ReturnBlockNumber: 100,
			expectedCallToGeth:  big.NewInt(int64(rpc.FinalizedBlockNumber)),
			expectedBlockNumber: 0,
		},
		{
			name:                "FinalizedBlockNumber(1000)-1000. is 0 ",
			blockPoint:          "finalized",
			offset:              -1000,
			l1ReturnBlockNumber: 1000,
			expectedCallToGeth:  big.NewInt(int64(rpc.FinalizedBlockNumber)),
			expectedBlockNumber: 0,
		},
		{
			name:                "FinalizedBlockNumber(1001)-1000. is 1 ",
			blockPoint:          "finalized",
			offset:              -1000,
			l1ReturnBlockNumber: 1001,
			expectedCallToGeth:  big.NewInt(int64(rpc.FinalizedBlockNumber)),
			expectedBlockNumber: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRequester := mock_l1_check_block.NewL1Requester(t)
			safeL1Block := l1_check_block.NewSafeL1BlockNumberFetch(l1_check_block.StringToL1BlockPoint(tt.blockPoint), tt.offset)

			mockRequester.EXPECT().HeaderByNumber(ctx, tt.expectedCallToGeth).Return(&types.Header{
				Number: big.NewInt(int64(tt.l1ReturnBlockNumber)),
			}, nil)
			blockNumber, err := safeL1Block.GetSafeBlockNumber(ctx, mockRequester)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBlockNumber, blockNumber)
		})
	}
}
