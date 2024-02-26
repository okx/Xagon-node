package metrics

import (
	"strconv"
	"sync"
	"time"
)

var instance *logStatisticsInstance
var once sync.Once

// GetLogStatistics is get log instance for statistic
func GetLogStatistics() logStatistics {
	once.Do(func() {
		instance = &logStatisticsInstance{}
		instance.init()
	})
	return instance
}

type logStatisticsInstance struct {
	timestamp  map[logTag]time.Time
	statistics map[logTag]int64 // value maybe the counter or time.Duration(ms)
	tags       map[logTag]string
}

func (l *logStatisticsInstance) init() {
	l.timestamp = make(map[logTag]time.Time)
	l.statistics = make(map[logTag]int64)
	l.tags = make(map[logTag]string)
}

func (l *logStatisticsInstance) CumulativeCounting(tag logTag) {
	l.statistics[tag]++
}

func (l *logStatisticsInstance) CumulativeValue(tag logTag, value int64) {
	l.statistics[tag] += value
}

func (l *logStatisticsInstance) CumulativeTiming(tag logTag, duration time.Duration) {
	l.statistics[tag] += duration.Milliseconds()
}

func (l *logStatisticsInstance) SetTag(tag logTag, value string) {
	l.tags[tag] = value
}

func (l *logStatisticsInstance) UpdateTimestamp(tag logTag, tm time.Time) {
	l.timestamp[tag] = tm
}

func (l *logStatisticsInstance) ResetStatistics() {
	l.statistics = make(map[logTag]int64)
	l.tags = make(map[logTag]string)
}

func (l *logStatisticsInstance) Summary() string {
	batchTotalDuration := "-"
	if key, ok := l.timestamp[NewRound]; ok {
		batchTotalDuration = strconv.Itoa(int(time.Since(key).Milliseconds()))
	}
	processTxTiming := "ProcessTx<" + strconv.Itoa(int(l.statistics[ProcessingTxTiming])) + "ms, " +
		"Commit<" + strconv.Itoa(int(l.statistics[ProcessingTxCommit])) + "ms>, " +
		"ProcessResponse<" + strconv.Itoa(int(l.statistics[ProcessingTxResponse])) + "ms>>, "

	finalizeBatchTiming := "FinalizeBatch<" + strconv.Itoa(int(l.statistics[FinalizeBatchTiming])) + "ms, " +
		"ReprocessFullBatch<" + strconv.Itoa(int(l.statistics[FinalizeBatchReprocessFullBatch])) + "ms>, " +
		"CloseBatch<" + strconv.Itoa(int(l.statistics[FinalizeBatchCloseBatch])) + "ms>, " +
		"OpenBatch<" + strconv.Itoa(int(l.statistics[FinalizeBatchOpenBatch])) + "ms>>, "

	result := "Batch<" + l.tags[FinalizeBatchNumber] + ">, " +
		"TotalDuration<" + batchTotalDuration + "ms>, " +
		"GasUsed<" + strconv.Itoa(int(l.statistics[BatchGas])) + ">, " +
		"Tx<" + strconv.Itoa(int(l.statistics[TxCounter])) + ">, " +
		"GetTx<" + strconv.Itoa(int(l.statistics[GetTx])) + "ms>, " +
		"GetTxPause<" + strconv.Itoa(int(l.statistics[GetTxPauseCounter])) + ">, " +
		"ReprocessTx<" + strconv.Itoa(int(l.statistics[ReprocessingTxCounter])) + ">, " +
		"FailTx<" + strconv.Itoa(int(l.statistics[FailTxCounter])) + ">, " +
		"InvalidTx<" + strconv.Itoa(int(l.statistics[ProcessingInvalidTxCounter])) + ">, " +
		processTxTiming +
		finalizeBatchTiming +
		"BatchCloseReason<" + l.tags[BatchCloseReason] + ">"

	return result
}

func (l *logStatisticsInstance) GetTag(tag logTag) string {
	return l.tags[tag]
}

func (l *logStatisticsInstance) GetStatistics(tag logTag) int64 {
	return l.statistics[tag]
}
