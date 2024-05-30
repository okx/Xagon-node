package sequencer

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	seqMetrics "github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
)

func (f *finalizer) tryToSleep() {
	fullBatchSleepDuration := getFullBatchSleepDuration(f.cfg.FullBatchSleepDuration.Duration)
	if fullBatchSleepDuration > 0 {
		log.Infof("Slow down sequencer: %v", fullBatchSleepDuration)
		time.Sleep(fullBatchSleepDuration)
		seqMetrics.GetLogStatistics().CumulativeCounting(seqMetrics.GetTxPauseCounter)
	}
}

// keepHalting keeps the finalizer halted
func (f *finalizer) keepHalting(ctx context.Context, err error) {
	f.haltFinalizer.Store(true)

	f.LogEvent(ctx, event.Level_Critical, event.EventID_FinalizerHalt, fmt.Sprintf("finalizer halted due to error: %s", err), nil)

	for {
		seqMetrics.HaltCount()
		log.Errorf("keep halting finalizer, error: %v", err)
		time.Sleep(5 * time.Second) //nolint:gomnd
	}
}
