package sequencer

func (w *Worker) deleteReadyTxCounter(addr string) {
	if w == nil || w.readyTxCounter == nil {
		return
	}
	delete(w.readyTxCounter, addr)
}

func (w *Worker) setReadyTxCounter(addr string, count uint64) {
	if w == nil || w.readyTxCounter == nil {
		return
	}
	w.readyTxCounter[addr] = count
}

// CountReadyTx returns the number of ready transactions
func (w *Worker) CountReadyTx() uint64 {
	if w == nil {
		return 0
	}
	w.workerMutex.Lock()
	defer w.workerMutex.Unlock()

	var count uint64
	for _, c := range w.readyTxCounter {
		count += c
	}
	return count
}
