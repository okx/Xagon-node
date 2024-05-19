package sequencer

func (a *addrQueue) GetTxCount() uint64 {
	if a == nil {
		return 0
	}
	var readyTxCount uint64
	if a.readyTx != nil {
		readyTxCount = 1
	}
	notReadyTxCount := uint64(len(a.notReadyTxs))
	forcedTxCount := uint64(len(a.forcedTxs))
	pendingTxsToStoreCount := uint64(len(a.pendingTxsToStore))

	return readyTxCount + notReadyTxCount + forcedTxCount + pendingTxsToStoreCount
}
