package sequencer

import (
	"sync"
	"sync/atomic"
)

// WaitGroupCount implements a sync.WaitGroup that also has a field to get the WaitGroup counter
type WaitGroupCount struct {
	sync.WaitGroup
	count atomic.Int32
}

// Add adds delta to the WaitGroup and increase the counter
func (wg *WaitGroupCount) Add(delta int) {
	wg.count.Add(int32(delta))
	wg.WaitGroup.Add(delta)
}

// Done decrements the WaitGroup and counter by one
func (wg *WaitGroupCount) Done() {
	wg.count.Add(-1)
	wg.WaitGroup.Done()
}

// Count returns the counter of the WaitGroup
func (wg *WaitGroupCount) Count() int {
	return int(wg.count.Load())
}
