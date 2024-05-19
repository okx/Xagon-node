package sequencer

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/state"
)

// BlockClosingReason is the reason why a block is closed.
type BlockClosingReason string

const (
	// BlockMaxDeltaTimestamp is the closing reason when the max delta timestamp is reached.
	BlockMaxDeltaTimestamp BlockClosingReason = "Max delta timestamp"
)

func getReasonFromBatch(batchCloseReason state.ClosingReason) BlockClosingReason {
	return BlockClosingReason(fmt.Sprintf("Batch closed, %v", batchCloseReason))
}

// Summary returns the metrics summary.
func (m *metrics) Summary(blockNum, batchNum, timestamp uint64) string {
	TotalSequencerTime := "SequencerTime<" + fmt.Sprintf("%v", m.sequencerTime().Milliseconds()) +
		"ms, newL2Block<" + fmt.Sprintf("%v", m.newL2BlockTimes.sequencer.Milliseconds()) +
		"ms>, txs<" + fmt.Sprintf("%v", m.transactionsTimes.sequencer.Milliseconds()) +
		"ms>, l2Block<" + fmt.Sprintf("%v", m.l2BlockTimes.sequencer.Milliseconds()) + "ms>>, "

	TotalExecutorTime := "ExecutorTime<" + fmt.Sprintf("%v", m.executorTime().Milliseconds()) +
		"ms, newL2Block<" + fmt.Sprintf("%v", m.newL2BlockTimes.executor.Milliseconds()) +
		"ms>, txs<" + fmt.Sprintf("%v", m.transactionsTimes.executor.Milliseconds()) +
		"ms>, l2Block<" + fmt.Sprintf("%v", m.l2BlockTimes.executor.Milliseconds()) + "ms>>, "

	result := "BlockNumber<" + fmt.Sprintf("%v", blockNum) + ">, " +
		"BatchNum<" + fmt.Sprintf("%v", batchNum) + ">, " +
		"TxCount<" + fmt.Sprintf("%v", m.l2BlockTxsCount) + ">, " +
		"Gas<" + fmt.Sprintf("%v", m.gas) + ">, " +
		"TotalTime<" + fmt.Sprintf("%v", m.totalTime().Milliseconds()) + "ms>, " +
		"IdleTime<" + fmt.Sprintf("%v", m.idleTime.Milliseconds()) + "ms>, " +
		TotalSequencerTime +
		TotalExecutorTime +
		"CloseReason<" + m.closeReason + ">, " +
		"Timestamp<" + fmt.Sprintf("%v", timestamp) + ">, "

	return result
}
