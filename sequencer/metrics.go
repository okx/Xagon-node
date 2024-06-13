package sequencer

import (
	"fmt"
	"math"
	"time"
)

// SEQUENTIAL L2 BLOCK PROCESSING
//           |-----------------------------------------------------------------------------| -> totalTime
//                        |------------|    |-------------------------|                      -> transactionsTime
//           |-newL2Block-|----tx 1----|    |---tx 2---|-----tx 3-----|  |-----l2Block-----|
// sequencer |sssss     ss|sss       ss|    |sss     ss|sss         ss|  |ssss           ss| -> sequencerTime
//  executor |     xxxxx  |   xxxxxxx  |    |   xxxxx  |   xxxxxxxxx  |  |    xxxxxxxxxxx  | -> executorTime
//      idle |                         |iiii|          |              |ii|                 | -> idleTime
//

// PARALLEL L2 BLOCK PROCESSING
//             |---------------------------------------------------------------------------------------------| -> totalTime
//             |-----------------------L2 block 1-----------------------|  |-----------L2 block 2------------|
//                          |------------|    |-------------------------|               |--------------------| -> transactionsTime
//             |-newL2Block-|----tx 1----|    |---tx 2---|-----tx 3-----|  |-newL2Block-|--tx 4---|---tx 5---|
//   sequencer |sssss     ss|sss       ss|    |sss     ss|sss         ss|  |sssss     ss|ss     ss|sss     ss| -> sequencerTime
//    executor |     xxxxx  |   xxxxxxx  |    |   xxxxx  |   xxxxxxxxx  |  |     xxxxx  |  xxxxxx |   xxxxx  | -> executorTime
//        idle |                         |iiii|          |              |ii|                                 | -> idleTime

//                                                                         | -> L2 block 1   |
// seq-l2block |                                                           |ssss           ss|
// exe-l2block |                                                           |    xxxxxxxxxxx  |
//

type processTimes struct {
	sequencer time.Duration
	executor  time.Duration
}

func (p *processTimes) total() time.Duration {
	return p.sequencer + p.executor
}

func (p *processTimes) sub(ptSub processTimes) {
	p.sequencer -= ptSub.sequencer
	p.executor -= ptSub.executor
}

func (p *processTimes) sumUp(ptSumUp processTimes) {
	p.sequencer += ptSumUp.sequencer
	p.executor += ptSumUp.executor
}

type metrics struct {
	closedAt           time.Time
	processedTxsCount  int64
	l2BlockTxsCount    int64
	idleTime           time.Duration
	newL2BlockTimes    processTimes
	transactionsTimes  processTimes
	l2BlockTimes       processTimes
	waitl2BlockTime    time.Duration
	gas                uint64
	estimatedTxsPerSec float64
	estimatedGasPerSec uint64
	sequential         bool
}

func (m *metrics) sub(mSub metrics) {
	m.processedTxsCount -= mSub.processedTxsCount
	m.l2BlockTxsCount -= mSub.l2BlockTxsCount
	m.idleTime -= mSub.idleTime
	m.newL2BlockTimes.sub(mSub.newL2BlockTimes)
	m.transactionsTimes.sub(mSub.transactionsTimes)
	m.l2BlockTimes.sub(mSub.l2BlockTimes)
	m.waitl2BlockTime -= mSub.waitl2BlockTime
	m.gas -= mSub.gas
}

func (m *metrics) sumUp(mSumUp metrics) {
	m.processedTxsCount += mSumUp.processedTxsCount
	m.l2BlockTxsCount += mSumUp.l2BlockTxsCount
	m.idleTime += mSumUp.idleTime
	m.newL2BlockTimes.sumUp(mSumUp.newL2BlockTimes)
	m.transactionsTimes.sumUp(mSumUp.transactionsTimes)
	m.l2BlockTimes.sumUp(mSumUp.l2BlockTimes)
	m.waitl2BlockTime += mSumUp.waitl2BlockTime
	m.gas += mSumUp.gas
}

func (m *metrics) executorTime() time.Duration {
	if m.sequential {
		return m.newL2BlockTimes.executor + m.transactionsTimes.executor + m.l2BlockTimes.executor
	} else {
		return m.newL2BlockTimes.executor + m.transactionsTimes.executor + m.waitl2BlockTime
	}
}

func (m *metrics) sequencerTime() time.Duration {
	if m.sequential {
		return m.newL2BlockTimes.sequencer + m.transactionsTimes.sequencer + m.l2BlockTimes.sequencer
	} else {
		return m.newL2BlockTimes.sequencer + m.transactionsTimes.sequencer
	}
}

func (m *metrics) totalTime() time.Duration {
	if m.sequential {
		return m.newL2BlockTimes.total() + m.transactionsTimes.total() + m.l2BlockTimes.total() + m.idleTime
	} else {
		return m.newL2BlockTimes.total() + m.transactionsTimes.total() + m.waitl2BlockTime + m.idleTime
	}
}

func (m *metrics) close(createdAt time.Time, l2BlockTxsCount int64, sequential bool) {
	// Compute pending fields
	m.closedAt = time.Now()
	totalTime := time.Since(createdAt)
	m.sequential = sequential
	m.l2BlockTxsCount = l2BlockTxsCount

	if m.sequential {
		m.transactionsTimes.sequencer = totalTime - m.idleTime - m.newL2BlockTimes.total() - m.transactionsTimes.executor - m.l2BlockTimes.total()
	} else {
		m.transactionsTimes.sequencer = totalTime - m.idleTime - m.newL2BlockTimes.total() - m.transactionsTimes.executor - m.waitl2BlockTime
	}

	// Compute performance
	if m.processedTxsCount > 0 {
		var timePerTxuS int64
		if m.sequential {
			// timePerTxuS is the average time spent per tx. This includes the l2Block time since the processing time of this section is proportional to the number of txs
			timePerTxuS = (m.transactionsTimes.total() + m.l2BlockTimes.total()).Microseconds() / m.processedTxsCount
		} else {
			// timePerTxuS is the average time spent per tx. This includes the waitl2Block
			timePerTxuS = (m.transactionsTimes.total() + m.waitl2BlockTime).Microseconds() / m.processedTxsCount
		}
		// estimatedTxs is the number of transactions that we estimate could have been processed in the block
		estimatedTxs := float64(totalTime.Microseconds()-m.newL2BlockTimes.total().Microseconds()) / float64(timePerTxuS)
		// estimatedTxxPerSec is the estimated transactions per second (rounded to 2 decimal digits)
		m.estimatedTxsPerSec = math.Ceil(estimatedTxs/totalTime.Seconds()*100) / 100 //nolint:gomnd

		// gasPerTx is the average gas used per tx
		gasPerTx := m.gas / uint64(m.processedTxsCount)
		// estimatedGasPerSec is the estimated gas per second
		m.estimatedGasPerSec = uint64(m.estimatedTxsPerSec * float64(gasPerTx))
	}
}

func (m *metrics) log() string {
	return fmt.Sprintf("blockTxs: %d, txs: %d, gas: %d, txsSec: %.2f, gasSec: %d, time: {total: %d, idle: %d, waitL2Block: %d, sequencer: {total: %d, newL2Block: %d, txs: %d, l2Block: %d}, executor: {total: %d, newL2Block: %d, txs: %d, l2Block: %d}",
		m.l2BlockTxsCount, m.processedTxsCount, m.gas, m.estimatedTxsPerSec, m.estimatedGasPerSec, m.totalTime().Microseconds(), m.idleTime.Microseconds(), m.waitl2BlockTime.Microseconds(),
		m.sequencerTime().Microseconds(), m.newL2BlockTimes.sequencer.Microseconds(), m.transactionsTimes.sequencer.Microseconds(), m.l2BlockTimes.sequencer.Microseconds(),
		m.executorTime().Microseconds(), m.newL2BlockTimes.executor.Microseconds(), m.transactionsTimes.executor.Microseconds(), m.l2BlockTimes.executor.Microseconds())
}

type intervalMetrics struct {
	l2Blocks    []*metrics
	maxInterval time.Duration
	metrics
	estimatedTxsPerSecAcc float64
	estimatedGasPerSecAcc uint64
	l2BlockCountAcc       int64
}

func newIntervalMetrics(maxInterval time.Duration) *intervalMetrics {
	return &intervalMetrics{
		l2Blocks:    []*metrics{},
		maxInterval: maxInterval,
		metrics:     metrics{},
	}
}

func (i *intervalMetrics) cleanUp() {
	now := time.Now()
	ct := 0
	for {
		if len(i.l2Blocks) == 0 {
			return
		}
		l2Block := i.l2Blocks[0]
		if l2Block.closedAt.Add(i.maxInterval).Before(now) {
			// Subtract l2Block metrics from accumulated values
			i.sub(*l2Block)
			if l2Block.processedTxsCount > 0 {
				i.estimatedTxsPerSecAcc -= l2Block.estimatedTxsPerSec
				i.estimatedGasPerSecAcc -= l2Block.estimatedGasPerSec
				i.l2BlockCountAcc--
			}
			// Remove from l2Blocks
			i.l2Blocks = i.l2Blocks[1:]
			ct++
		} else {
			break
		}
	}

	if ct > 0 {
		// Compute performance
		i.computePerformance()
	}
}

func (i *intervalMetrics) addL2BlockMetrics(l2Block metrics) {
	i.cleanUp()

	i.sumUp(l2Block)
	if l2Block.processedTxsCount > 0 {
		i.estimatedTxsPerSecAcc += l2Block.estimatedTxsPerSec
		i.estimatedGasPerSecAcc += l2Block.estimatedGasPerSec
		i.l2BlockCountAcc++
		i.computePerformance()
	}

	i.l2Blocks = append(i.l2Blocks, &l2Block)
}

func (i *intervalMetrics) computePerformance() {
	if i.l2BlockCountAcc > 0 {
		i.estimatedTxsPerSec = math.Ceil(i.estimatedTxsPerSecAcc/float64(i.l2BlockCountAcc)*100) / 100 //nolint:gomnd
		i.estimatedGasPerSec = i.estimatedGasPerSecAcc / uint64(i.l2BlockCountAcc)
	} else {
		i.estimatedTxsPerSec = 0
		i.estimatedGasPerSec = 0
	}
}

func (i *intervalMetrics) startsAt() time.Time {
	return time.Now().Add(-i.maxInterval)
}
