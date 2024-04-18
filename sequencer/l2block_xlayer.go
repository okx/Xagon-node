package sequencer

func (f *finalizer) setWIPL2BlockCloseReason(closeReason BlockClosingReason) {
	if f.wipL2Block != nil {
		f.wipL2Block.metrics.closeReason = string(closeReason)
	}
}
