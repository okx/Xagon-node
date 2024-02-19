package sequencesender

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/stretchr/testify/assert"
)

func TestIsSynced(t *testing.T) {
	const (
		retries   = 3
		waitRetry = 1 * time.Second
	)

	type IsSyncedTestCase = struct {
		name                   string
		lastVirtualBatchNum    uint64
		lastTrustedBatchClosed uint64
		lastSCBatchNum         []uint64
		expectedResult         bool
		err                    error
	}

	mockError := errors.New("error")

	stateMock := new(StateMock)
	ethermanMock := new(EthermanMock)
	ethTxManagerMock := new(EthTxManagerMock)
	ssender, err := New(Config{}, stateMock, ethermanMock, ethTxManagerMock, nil)
	assert.NoError(t, err)

	testCases := []IsSyncedTestCase{
		{
			name:                   "is synced",
			lastVirtualBatchNum:    10,
			lastTrustedBatchClosed: 12,
			lastSCBatchNum:         []uint64{10},
			expectedResult:         true,
			err:                    nil,
		},
		{
			name:                   "not synced",
			lastVirtualBatchNum:    9,
			lastTrustedBatchClosed: 12,
			lastSCBatchNum:         []uint64{10},
			expectedResult:         false,
			err:                    nil,
		},
		{
			name:                   "error virtual > trusted",
			lastVirtualBatchNum:    10,
			lastTrustedBatchClosed: 9,
			lastSCBatchNum:         []uint64{10},
			expectedResult:         false,
			err:                    ErrSyncVirtualGreaterTrusted,
		},
		{
			name:                   "error virtual > sc sequenced",
			lastVirtualBatchNum:    11,
			lastTrustedBatchClosed: 12,
			lastSCBatchNum:         []uint64{10, 10, 10, 10},
			expectedResult:         false,
			err:                    ErrSyncVirtualGreaterSequenced,
		},
		{
			name:                   "is synced, sc sequenced retries",
			lastVirtualBatchNum:    11,
			lastTrustedBatchClosed: 12,
			lastSCBatchNum:         []uint64{10, 10, 11},
			expectedResult:         true,
			err:                    nil,
		},
		{
			name:                   "is synced, sc sequenced retries (last)",
			lastVirtualBatchNum:    11,
			lastTrustedBatchClosed: 12,
			lastSCBatchNum:         []uint64{10, 10, 10, 11},
			expectedResult:         true,
			err:                    nil,
		},
		{
			name:                   "error state.GetLastVirtualBatchNum",
			lastVirtualBatchNum:    0,
			lastTrustedBatchClosed: 12,
			lastSCBatchNum:         []uint64{0},
			expectedResult:         false,
			err:                    nil,
		},
		{
			name:                   "error state.GetLastClosedBatch",
			lastVirtualBatchNum:    11,
			lastTrustedBatchClosed: 0,
			lastSCBatchNum:         []uint64{0},
			expectedResult:         false,
			err:                    nil,
		},
		{
			name:                   "error etherman.GetLatestBatchNumber",
			lastVirtualBatchNum:    11,
			lastTrustedBatchClosed: 12,
			lastSCBatchNum:         []uint64{0},
			expectedResult:         false,
			err:                    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var returnError error
			returnError = nil

			if tc.lastVirtualBatchNum == 0 {
				returnError = mockError
			}
			stateMock.On("GetLastVirtualBatchNum", context.Background(), nil).Return(tc.lastVirtualBatchNum, returnError).Once()

			if returnError == nil { // if previous call to mock function returns error then this function will be not called inside isSynced
				if tc.lastTrustedBatchClosed == 0 {
					returnError = mockError
				}
				stateMock.On("GetLastClosedBatch", context.Background(), nil).Return(&state.Batch{BatchNumber: tc.lastTrustedBatchClosed}, returnError).Once()
			}

			if returnError == nil { // if previous call to mock function returns error then this function will be not called inside isSynced
				for _, num := range tc.lastSCBatchNum {
					if num == 0 { // 0 means the function returns error
						returnError = mockError
					}
					ethermanMock.On("GetLatestBatchNumber").Return(num, returnError).Once()
				}
			}

			synced, err := ssender.isSynced(context.Background(), retries, waitRetry)

			assert.EqualValues(t, tc.expectedResult, synced)
			assert.EqualValues(t, tc.err, err)

			ethermanMock.AssertExpectations(t)
			stateMock.AssertExpectations(t)
		})
	}
}
