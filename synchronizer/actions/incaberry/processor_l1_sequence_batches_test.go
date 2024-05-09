package incaberry

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	mocks "github.com/0xPolygonHermez/zkevm-node/synchronizer/mocks"
	"github.com/stretchr/testify/require"
)

func TestProcessorL1SequenceBatches_Process(t *testing.T) {
	ctx := context.Background()
	sut := NewProcessorL1SequenceBatches(nil, nil, nil, nil, nil)

	l1Block := &etherman.Block{
		//SequencedBatches: []Batch{}, // Mock sequenced batches
		BlockNumber: 123, // Mock block number
	}

	dbTx := mocks.NewDbTxMock(t)

	// Create an instance of ProcessorL1SequenceBatches

	// Test invalid call, no sequenced batches
	sut.ProcessorBase.SupportedEvents()
	err := sut.Process(ctx, etherman.Order{Name: sut.SupportedEvents()[0], Pos: 0}, l1Block, dbTx)
	require.Error(t, err)
}
