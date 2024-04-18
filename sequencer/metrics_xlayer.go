package sequencer

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/state"
)

// BlockClosingReason is the reason why a block is closed.
type BlockClosingReason string

const (
	// BlockMaxDeltaTimestamp is the closing reason when the max delta timestamp is reached.
	BlockMaxDeltaTimestamp BlockClosingReason = "Block closed, max delta timestamp"
)

func getReasonFromBatch(batchCloseReason state.ClosingReason) BlockClosingReason {
	return BlockClosingReason(fmt.Sprintf("Batch closed, %v", batchCloseReason))
}

// Summary returns the metrics summary.
func (m *metrics) Summary(blockNum, batchNum, timestamp uint64) string {
	TotalSequencerTime := "sequencer<" + fmt.Sprintf("%v", m.sequencerTime().Microseconds()) +
		"ms, newL2Block<" + fmt.Sprintf("%v", m.newL2BlockTimes.sequencer.Microseconds()) +
		"ms>, txs<" + fmt.Sprintf("%v", m.transactionsTimes.sequencer.Microseconds()) +
		"ms>, l2Block<" + fmt.Sprintf("%v", m.l2BlockTimes.sequencer.Microseconds()) + ">>, "

	TotalExecutorTime := "executor<" + fmt.Sprintf("%v", m.executorTime().Microseconds()) +
		"ms, newL2Block<" + fmt.Sprintf("%v", m.newL2BlockTimes.executor.Microseconds()) +
		"ms>, txs<" + fmt.Sprintf("%v", m.transactionsTimes.executor.Microseconds()) +
		"ms>, l2Block<" + fmt.Sprintf("%v", m.l2BlockTimes.executor.Microseconds()) + ">>, "

	result := "BlockNumber<" + fmt.Sprintf("%v", blockNum) + ">, " +
		"BatchNum<" + fmt.Sprintf("%v", batchNum) + ">, " +
		"Timestamp<" + fmt.Sprintf("%v", timestamp) + ">, " +
		"TxCount<" + fmt.Sprintf("%v", m.l2BlockTxsCount) + ">, " +
		"Gas<" + fmt.Sprintf("%v", m.gas) + ">, " +
		"time<" + fmt.Sprintf("%v", m.totalTime().Microseconds()) + "ms>, " +
		"idleTime<" + fmt.Sprintf("%v", m.idleTime.Microseconds()) + "ms>, " +
		TotalSequencerTime +
		TotalExecutorTime +
		"closeReason<" + m.closeReason + ">, "

	return result
}
