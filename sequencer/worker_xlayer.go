package sequencer

// GetAddressCount returns the number of transactions for a given address
func (w *Worker) GetAddressCount(addr string) uint64 {
	if w == nil || w.pool == nil || w.pool[addr] == nil {
		return 0
	}

	return 0
}
