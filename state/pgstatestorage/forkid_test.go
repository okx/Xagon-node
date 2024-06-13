package pgstatestorage

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortIndexForForkdIDSortedByBlockNumber(t *testing.T) {
	forkIDs := []state.ForkIDInterval{
		{BlockNumber: 10, ForkId: 1},
		{BlockNumber: 5, ForkId: 2},
		{BlockNumber: 15, ForkId: 3},
		{BlockNumber: 1, ForkId: 4},
	}

	expected := []int{3, 1, 0, 2}
	actual := sortIndexForForkdIDSortedByBlockNumber(forkIDs)

	assert.Equal(t, expected, actual)

	// Ensure that the original slice is not modified
	assert.Equal(t, []state.ForkIDInterval{
		{BlockNumber: 10, ForkId: 1},
		{BlockNumber: 5, ForkId: 2},
		{BlockNumber: 15, ForkId: 3},
		{BlockNumber: 1, ForkId: 4},
	}, forkIDs)

	// Ensure that the sorted slice is sorted correctly
	sortedForkIDs := make([]state.ForkIDInterval, len(forkIDs))
	for i, idx := range actual {
		sortedForkIDs[i] = forkIDs[idx]
	}
	previousBlock := sortedForkIDs[0].BlockNumber
	for _, forkID := range sortedForkIDs {
		require.GreaterOrEqual(t, forkID.BlockNumber, previousBlock)
		previousBlock = forkID.BlockNumber
	}
}

func TestGetForkIDByBlockNumber(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		blockNumber uint64
		expected    uint64
	}{
		{
			name:        "Block number is less than the first interval",
			blockNumber: 1,
			expected:    1,
		},
		{
			name:        "Block number is equal to the first interval",
			blockNumber: 10,
			expected:    1,
		},
		{
			name:        "Block number is between two intervals",
			blockNumber: 11,
			expected:    1,
		},
		{
			name:        "Block number is equal to an interval",
			blockNumber: 200,
			expected:    2,
		},
		{
			name:        "Block number is greater to an interval",
			blockNumber: 201,
			expected:    2,
		},
		{
			name:        "Block number is greater than the last interval",
			blockNumber: 600,
			expected:    4,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := state.Config{
				ForkIDIntervals: []state.ForkIDInterval{
					{BlockNumber: 10, ForkId: 1},
					{BlockNumber: 200, ForkId: 2},
					{BlockNumber: 400, ForkId: 3},
					{BlockNumber: 500, ForkId: 4},
				},
			}
			storage := NewPostgresStorage(cfg, nil)
			// Create a new State instance with test data
			state := state.NewState(cfg, storage, nil, nil, nil, nil, nil)

			// Call the function being tested
			actual := state.GetForkIDByBlockNumber(tc.blockNumber)

			// Check the result
			assert.Equal(t, tc.expected, actual)
		})
	}
}
