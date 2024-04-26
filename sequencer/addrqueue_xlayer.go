package sequencer

func (a *addrQueue) GetTxCount() int64 {
	if a == nil {
		return 0
	}
	var readyTxCount int64
	if a.readyTx != nil {
		readyTxCount = 1
	}
	notReadyTxCount := int64(len(a.notReadyTxs))
	forcedTxCount := int64(len(a.forcedTxs))
	pendingTxsToStoreCount := int64(len(a.pendingTxsToStore))

	return readyTxCount + notReadyTxCount + forcedTxCount + pendingTxsToStoreCount
}
