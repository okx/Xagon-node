package sequencer

import (
	"sync"
)

// PoolReadyTxCounter is the pool tx counter
type PoolReadyTxCounter struct {
	// Ready is the number of ready transactions
	Address map[string]int64

	// RwMutex is the mutex for the pool tx counter
	RwMutex sync.RWMutex
}

var poolReadyTxCounterInst *PoolReadyTxCounter
var poolReadyTxCounterOnce sync.Once

func getPoolReadyTxCounter() *PoolReadyTxCounter {
	poolReadyTxCounterOnce.Do(func() {
		poolReadyTxCounterInst = &PoolReadyTxCounter{}
	})
	return poolReadyTxCounterInst
}

func (ptx *PoolReadyTxCounter) set(addr string, count int64) {
	ptx.RwMutex.Lock()
	defer ptx.RwMutex.Unlock()
	if ptx.Address == nil {
		ptx.Address = make(map[string]int64)
	}
	ptx.Address[addr] = count
}

func (ptx *PoolReadyTxCounter) delete(addr string) {
	ptx.RwMutex.Lock()
	defer ptx.RwMutex.Unlock()
	delete(ptx.Address, addr)
}

// Sum returns the sum of the ready tx counter
func (ptx *PoolReadyTxCounter) Sum() int64 {
	ptx.RwMutex.RLock()
	defer ptx.RwMutex.RUnlock()
	var sum int64
	for _, v := range ptx.Address {
		sum += v
	}
	return sum
}
