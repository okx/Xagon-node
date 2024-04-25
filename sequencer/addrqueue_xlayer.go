package sequencer

func (a *addrQueue) GetTxCount() int64 {
	if a == nil {
		return 0
	}
	var readyTxCount, notReadyTxCount int64
	if a.readyTx != nil {
		readyTxCount = 1
	}
	notReadyTxCount = int64(len(a.notReadyTxs))

	return readyTxCount + notReadyTxCount
}
