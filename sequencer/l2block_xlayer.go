package sequencer

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

func (f *finalizer) setWIPL2BlockCloseReason(closeReason BlockClosingReason) {
	if f.wipL2Block != nil {
		f.wipL2Block.metrics.closeReason = string(closeReason)
	}
}

func (f *finalizer) waitL2BlocksProcessed() {
	if f.pendingL2BlocksToProcessWG.Count() > 0 {
		startWait := time.Now()
		f.pendingL2BlocksToProcessWG.Wait()
		waitTime := time.Since(startWait)
		log.Debugf("waiting for previous L2 block to be processed took: %v", waitTime)
		f.wipL2Block.metrics.waitl2BlockTime = waitTime
	}
}
