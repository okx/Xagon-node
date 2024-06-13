package l1_check_block

import (
	"context"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

// L1BlockChecker is an interface that defines the method to check L1 blocks
type L1BlockChecker interface {
	Step(ctx context.Context) error
}

const (
	defaultPeriodTime = time.Second
)

// AsyncCheck is a wrapper for L1BlockChecker to become asynchronous
type AsyncCheck struct {
	checker      L1BlockChecker
	mutex        sync.Mutex
	lastResult   *syncinterfaces.IterationResult
	onFinishCall func()
	periodTime   time.Duration
	// Wg is a wait group to wait for the result
	Wg        sync.WaitGroup
	ctx       context.Context
	cancelCtx context.CancelFunc
	isRunning bool
}

// NewAsyncCheck creates a new AsyncCheck
func NewAsyncCheck(checker L1BlockChecker) *AsyncCheck {
	return &AsyncCheck{
		checker:    checker,
		periodTime: defaultPeriodTime,
	}
}

// SetPeriodTime sets the period time between relaunch checker.Step
func (a *AsyncCheck) SetPeriodTime(periodTime time.Duration) {
	a.periodTime = periodTime
}

// Run is a method that starts the async check
func (a *AsyncCheck) Run(ctx context.Context, onFinish func()) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.onFinishCall = onFinish
	if a.isRunning {
		log.Infof("%s L1BlockChecker: already running, changing onFinish call", logPrefix)
		return
	}
	a.lastResult = nil
	a.ctx, a.cancelCtx = context.WithCancel(ctx)
	a.launchChecker(a.ctx)
}

// Stop is a method that stops the async check
func (a *AsyncCheck) Stop() {
	a.cancelCtx()
	a.Wg.Wait()
}

// RunSynchronous is a method that forces the check to be synchronous before starting the async check
func (a *AsyncCheck) RunSynchronous(ctx context.Context) syncinterfaces.IterationResult {
	return a.executeIteration(ctx)
}

// GetResult returns the last result of the check:
// - Nil -> still running
// - Not nil -> finished, and this is the result. You must call again Run to start a new check
func (a *AsyncCheck) GetResult() *syncinterfaces.IterationResult {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.lastResult
}

// https://stackoverflow.com/questions/32840687/timeout-for-waitgroup-wait
// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// GetResultBlockingUntilAvailable wait the time specific in timeout, if reach timeout returns current
// result, if not, wait until the result is available.
// if timeout is 0, it waits indefinitely
func (a *AsyncCheck) GetResultBlockingUntilAvailable(timeout time.Duration) *syncinterfaces.IterationResult {
	if timeout == 0 {
		a.Wg.Wait()
	} else {
		waitTimeout(&a.Wg, timeout)
	}
	return a.GetResult()
}

func (a *AsyncCheck) setResult(result syncinterfaces.IterationResult) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.lastResult = &result
}

func (a *AsyncCheck) launchChecker(ctx context.Context) {
	// add waitGroup  to wait for a result
	a.Wg.Add(1)
	a.isRunning = true
	go func() {
		log.Infof("%s L1BlockChecker: starting background process", logPrefix)
		for {
			result := a.step(ctx)
			if result != nil {
				a.setResult(*result)
				// Result is set wg is done
				break
			}
		}
		log.Infof("%s L1BlockChecker: finished background process", logPrefix)
		a.Wg.Done()
		a.mutex.Lock()
		onFinishCall := a.onFinishCall
		a.isRunning = false
		a.mutex.Unlock()
		// call onFinish function with no mutex
		if onFinishCall != nil {
			onFinishCall()
		}
	}()
}

// step is a method that executes until executeItertion
// returns an error or a reorg
func (a *AsyncCheck) step(ctx context.Context) *syncinterfaces.IterationResult {
	select {
	case <-ctx.Done():
		log.Debugf("%s L1BlockChecker: context done", logPrefix)
		return &syncinterfaces.IterationResult{Err: ctx.Err()}
	default:
		result := a.executeIteration(ctx)
		if result.ReorgDetected {
			return &result
		}
		log.Debugf("%s L1BlockChecker:returned %s waiting %s to relaunch", logPrefix, result.String(), a.periodTime)
		time.Sleep(a.periodTime)
	}
	return nil
}

// executeIteration executes a single iteration of the checker
func (a *AsyncCheck) executeIteration(ctx context.Context) syncinterfaces.IterationResult {
	res := syncinterfaces.IterationResult{}
	log.Debugf("%s calling checker.Step(...)", logPrefix)
	res.Err = a.checker.Step(ctx)
	log.Debugf("%s returned checker.Step(...) %w", logPrefix, res.Err)
	if res.Err != nil {
		log.Errorf("%s Fail check L1 Blocks: %w", logPrefix, res.Err)
		if common.IsReorgError(res.Err) {
			// log error
			blockNumber := common.GetReorgErrorBlockNumber(res.Err)
			log.Infof("%s Reorg detected at block %d", logPrefix, blockNumber)
			// It keeps blocked until the channel is read
			res.BlockNumber = blockNumber
			res.ReorgDetected = true
			res.ReorgMessage = res.Err.Error()
			res.Err = nil
		}
	}
	return res
}
