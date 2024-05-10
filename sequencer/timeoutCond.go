package sequencer

import (
	"sync"
	"time"
)

type timeoutCond struct {
	L  sync.Locker
	ch chan bool
}

func newTimeoutCond(l sync.Locker) *timeoutCond {
	return &timeoutCond{ch: make(chan bool), L: l}
}

func (t *timeoutCond) Wait() {
	t.L.Unlock()
	<-t.ch
	t.L.Lock()
}

func (t *timeoutCond) WaitOrTimeout(d time.Duration) bool {
	timeout := time.NewTimer(d)
	t.L.Unlock()
	var r bool
	select {
	case <-timeout.C:
		r = false
	case <-t.ch:
		r = true
	}
	if !timeout.Stop() {
		select {
		case <-timeout.C:
		default:
		}
	}
	t.L.Lock()
	return r
}

func (t *timeoutCond) Signal() {
	t.signal()
}

func (t *timeoutCond) Broadcast() {
	for {
		// Stop when we run out of waiters
		if !t.signal() {
			return
		}
	}
}

func (t *timeoutCond) signal() bool {
	select {
	case t.ch <- true:
		return true
	default:
		return false
	}
}
