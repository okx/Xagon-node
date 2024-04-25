package sequencer

import (
	"sync"
)

// PoolTxCounter is the pool tx counter
type PoolTxCounter struct {
	// Ready is the number of ready transactions
	Address map[string]int64

	// RwMutex is the mutex for the pool tx counter
	RwMutex sync.RWMutex
}

var poolTxCounterInst *PoolTxCounter
var poolTxCounterOnce sync.Once

func getPoolTxCounter() *PoolTxCounter {
	poolTxCounterOnce.Do(func() {
		poolTxCounterInst = &PoolTxCounter{}
	})
	return poolTxCounterInst
}

func (ptx *PoolTxCounter) set(addr string, count int64) {
	ptx.RwMutex.Lock()
	defer ptx.RwMutex.Unlock()
	if ptx.Address == nil {
		ptx.Address = make(map[string]int64)
	}
	ptx.Address[addr] = count
}

func (ptx *PoolTxCounter) delete(addr string) {
	ptx.RwMutex.Lock()
	defer ptx.RwMutex.Unlock()
	delete(ptx.Address, addr)
}

func (ptx *PoolTxCounter) sum() int64 {
	ptx.RwMutex.RLock()
	defer ptx.RwMutex.RUnlock()
	var sum int64
	for _, v := range ptx.Address {
		sum += v
	}
	return sum
}
