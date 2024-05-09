package l1event_orders

import (
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// GetSequenceFromL1EventOrder returns the sequence of batches of  given event
// There are event that are Batch based or not, if not it returns a nil
func GetSequenceFromL1EventOrder(event etherman.EventOrder, l1Block *etherman.Block, position int) *state.Sequence {
	switch event {
	case etherman.InitialSequenceBatchesOrder:
		return getSequence(l1Block.SequencedBatches[position],
			func(batch etherman.SequencedBatch) uint64 { return batch.BatchNumber })
	case etherman.SequenceBatchesOrder:
		return getSequence(l1Block.SequencedBatches[position],
			func(batch etherman.SequencedBatch) uint64 { return batch.BatchNumber })
	case etherman.ForcedBatchesOrder:
		bn := l1Block.ForcedBatches[position].ForcedBatchNumber
		return &state.Sequence{FromBatchNumber: bn, ToBatchNumber: bn}
	case etherman.UpdateEtrogSequenceOrder:
		bn := l1Block.UpdateEtrogSequence.BatchNumber
		return &state.Sequence{FromBatchNumber: bn, ToBatchNumber: bn}
	case etherman.SequenceForceBatchesOrder:
		getSequence(l1Block.SequencedForceBatches[position],
			func(batch etherman.SequencedForceBatch) uint64 { return batch.BatchNumber })
	case etherman.TrustedVerifyBatchOrder:
		bn := l1Block.VerifiedBatches[position].BatchNumber
		return &state.Sequence{FromBatchNumber: bn, ToBatchNumber: bn}
	}
	return nil
}

func getSequence[T any](batches []T, getBatchNumber func(T) uint64) *state.Sequence {
	if len(batches) == 0 {
		return nil
	}
	res := state.Sequence{FromBatchNumber: getBatchNumber(batches[0]),
		ToBatchNumber: getBatchNumber(batches[0])}
	for _, batch := range batches {
		if getBatchNumber(batch) < res.FromBatchNumber {
			res.FromBatchNumber = getBatchNumber(batch)
		}
		if getBatchNumber(batch) > res.ToBatchNumber {
			res.ToBatchNumber = getBatchNumber(batch)
		}
	}
	return &res
}
