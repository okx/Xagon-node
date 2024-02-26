package metrics

import (
	"testing"
	"time"
)

func Test_logStatisticsInstance_Summary(t *testing.T) {
	type fields struct {
		timestamp  map[logTag]time.Time
		statistics map[logTag]int64
		tags       map[logTag]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{"1", fields{
			timestamp: map[logTag]time.Time{NewRound: time.Now().Add(-time.Second)},
			statistics: map[logTag]int64{
				BatchGas:                        111111,
				TxCounter:                       10,
				GetTx:                           time.Second.Milliseconds(),
				GetTxPauseCounter:               2,
				ReprocessingTxCounter:           3,
				FailTxCounter:                   1,
				ProcessingInvalidTxCounter:      2,
				ProcessingTxTiming:              time.Second.Milliseconds() * 30,
				ProcessingTxCommit:              time.Second.Milliseconds() * 10,
				ProcessingTxResponse:            time.Second.Milliseconds() * 15,
				FinalizeBatchTiming:             time.Second.Milliseconds() * 50,
				FinalizeBatchReprocessFullBatch: time.Second.Milliseconds() * 20,
				FinalizeBatchCloseBatch:         time.Second.Milliseconds() * 10,
				FinalizeBatchOpenBatch:          time.Second.Milliseconds() * 10,
			},
			tags: map[logTag]string{BatchCloseReason: "deadline", FinalizeBatchNumber: "123"},
		}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logStatisticsInstance{
				timestamp:  tt.fields.timestamp,
				statistics: tt.fields.statistics,
				tags:       tt.fields.tags,
			}
			t.Log(l.Summary())
		})
	}
}
