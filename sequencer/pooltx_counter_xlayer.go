package sequencer

import (
	"sync"
	"sync/atomic"
)

// PoolReadyTxCounter is the struct that holds the ready tx counter
type PoolReadyTxCounter struct {
	// Count is the number of ready transactions
	Count uint64
}

var poolReadyTxCounterInst *PoolReadyTxCounter
var poolReadyTxCounterOnce sync.Once

func getPoolReadyTxCounter() *PoolReadyTxCounter {
	poolReadyTxCounterOnce.Do(func() {
		poolReadyTxCounterInst = &PoolReadyTxCounter{}
	})
	return poolReadyTxCounterInst
}

func (ptx *PoolReadyTxCounter) setReadyTxCount(count uint64) {
	atomic.StoreUint64(&ptx.Count, count)
}

// Sum returns the sum of the ready tx counter
func (ptx *PoolReadyTxCounter) getReadyTxCount() uint64 {
	return atomic.LoadUint64(&ptx.Count)
}
