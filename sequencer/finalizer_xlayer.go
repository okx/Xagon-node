package sequencer

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

func (f *finalizer) tryToSleep() {
	fullBatchSleepDuration := getFullBatchSleepDuration(f.cfg.FullBatchSleepDuration.Duration)
	if fullBatchSleepDuration > 0 {
		log.Infof("Slow down sequencer: %v", fullBatchSleepDuration)
		time.Sleep(fullBatchSleepDuration)
	}
}
