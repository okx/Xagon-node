package l1_check_block

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/jackc/pgx/v4"
)

// StateForL1BlockCheckerIntegration is an interface for the state
type StateForL1BlockCheckerIntegration interface {
	GetPreviousBlockToBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*state.Block, error)
}

// L1BlockCheckerIntegration is a struct that integrates the L1BlockChecker with the synchronizer
type L1BlockCheckerIntegration struct {
	forceCheckOnStart  bool
	checker            syncinterfaces.AsyncL1BlockChecker
	preChecker         syncinterfaces.AsyncL1BlockChecker
	state              StateForL1BlockCheckerIntegration
	sync               SyncCheckReorger
	timeBetweenRetries time.Duration
}

// SyncCheckReorger is an interface that defines the methods required from Synchronizer object
type SyncCheckReorger interface {
	ExecuteReorgFromMismatchBlock(blockNumber uint64, reason string) error
	OnDetectedMismatchL1BlockReorg()
}

// NewL1BlockCheckerIntegration creates a new L1BlockCheckerIntegration
func NewL1BlockCheckerIntegration(checker syncinterfaces.AsyncL1BlockChecker, preChecker syncinterfaces.AsyncL1BlockChecker, state StateForL1BlockCheckerIntegration, sync SyncCheckReorger, forceCheckOnStart bool, timeBetweenRetries time.Duration) *L1BlockCheckerIntegration {
	return &L1BlockCheckerIntegration{
		forceCheckOnStart:  forceCheckOnStart,
		checker:            checker,
		preChecker:         preChecker,
		state:              state,
		sync:               sync,
		timeBetweenRetries: timeBetweenRetries,
	}
}

// OnStart is a method that is called before starting the synchronizer
func (v *L1BlockCheckerIntegration) OnStart(ctx context.Context) error {
	if v.forceCheckOnStart {
		log.Infof("%s Forcing L1BlockChecker check before start", logPrefix)
		result := v.runCheckerSync(ctx, v.checker)
		if result.ReorgDetected {
			v.executeResult(ctx, result)
		} else {
			log.Infof("%s Forcing L1BlockChecker check:OK ", logPrefix)
			if v.preChecker != nil {
				log.Infof("%s Forcing L1BlockChecker preCheck before start", logPrefix)
				result = v.runCheckerSync(ctx, v.preChecker)
				if result.ReorgDetected {
					v.executeResult(ctx, result)
				} else {
					log.Infof("%s Forcing L1BlockChecker preCheck:OK", logPrefix)
				}
			}
		}
	}
	v.launch(ctx)
	return nil
}

func (v *L1BlockCheckerIntegration) runCheckerSync(ctx context.Context, checker syncinterfaces.AsyncL1BlockChecker) syncinterfaces.IterationResult {
	for {
		result := checker.RunSynchronous(ctx)
		if result.Err == nil {
			return result
		} else {
			time.Sleep(v.timeBetweenRetries)
		}
	}
}

// OnStartL1Sync is a method that is called before starting the L1 sync
func (v *L1BlockCheckerIntegration) OnStartL1Sync(ctx context.Context) bool {
	return v.checkBackgroundResult(ctx, "before start L1 sync")
}

// OnStartL2Sync is a method that is called before starting the L2 sync
func (v *L1BlockCheckerIntegration) OnStartL2Sync(ctx context.Context) bool {
	return v.checkBackgroundResult(ctx, "before start 2 sync")
}

// OnResetState is a method that is called after a resetState
func (v *L1BlockCheckerIntegration) OnResetState(ctx context.Context) {
	log.Infof("%s L1BlockChecker: after a resetState relaunch background process", logPrefix)
	v.launch(ctx)
}

// CheckReorgWrapper is a wrapper over reorg function of synchronizer.
// it checks the result of the function and the result of background process and decides which return
func (v *L1BlockCheckerIntegration) CheckReorgWrapper(ctx context.Context, reorgFirstBlockOk *state.Block, errReportedByReorgFunc error) (*state.Block, error) {
	resultBackground := v.getMergedResults()
	if resultBackground != nil && resultBackground.ReorgDetected {
		// Background process detected a reorg, decide which return
		firstOkBlockBackgroundCheck, err := v.state.GetPreviousBlockToBlockNumber(ctx, resultBackground.BlockNumber, nil)
		if err != nil {
			log.Warnf("%s Error getting previous block to block number where a reorg have been detected %d: %s. So we reorgFunc values", logPrefix, resultBackground.BlockNumber, err)
			return reorgFirstBlockOk, errReportedByReorgFunc
		}
		if reorgFirstBlockOk == nil || errReportedByReorgFunc != nil {
			log.Infof("%s Background checker detects bad block at block %d (first block ok %d) and regular reorg function no. Returning it", logPrefix,
				resultBackground.BlockNumber, firstOkBlockBackgroundCheck.BlockNumber)
			return firstOkBlockBackgroundCheck, nil
		}
		if firstOkBlockBackgroundCheck.BlockNumber < reorgFirstBlockOk.BlockNumber {
			// Background process detected a reorg at oldest block
			log.Warnf("%s Background checker detects bad block  at block %d (first block ok %d) and regular reorg function first block ok: %d. Returning from %d",
				logPrefix, resultBackground.BlockNumber, firstOkBlockBackgroundCheck.BlockNumber, reorgFirstBlockOk.BlockNumber, firstOkBlockBackgroundCheck.BlockNumber)
			return firstOkBlockBackgroundCheck, nil
		} else {
			// Regular reorg function detected a reorg at oldest block
			log.Warnf("%s Background checker detects bad block  at block %d (first block ok %d) and regular reorg function first block ok: %d. Executing from %d",
				logPrefix, resultBackground.BlockNumber, firstOkBlockBackgroundCheck.BlockNumber, reorgFirstBlockOk.BlockNumber, reorgFirstBlockOk.BlockNumber)
			return reorgFirstBlockOk, errReportedByReorgFunc
		}
	}
	if resultBackground != nil && !resultBackground.ReorgDetected {
		// Relaunch checker, if there is a reorg, It is going to be relaunched after (OnResetState)
		v.launch(ctx)
	}
	// Background process doesnt have anything to we return the regular reorg function result
	return reorgFirstBlockOk, errReportedByReorgFunc
}

func (v *L1BlockCheckerIntegration) checkBackgroundResult(ctx context.Context, positionMessage string) bool {
	log.Debugf("%s Checking L1BlockChecker %s", logPrefix, positionMessage)
	result := v.getMergedResults()
	if result != nil {
		if result.ReorgDetected {
			log.Warnf("%s Checking L1BlockChecker %s: reorg detected %s", logPrefix, positionMessage, result.String())
			v.executeResult(ctx, *result)
		}
		v.launch(ctx)
		return result.ReorgDetected
	}
	return false
}

func (v *L1BlockCheckerIntegration) getMergedResults() *syncinterfaces.IterationResult {
	result := v.checker.GetResult()
	var preResult *syncinterfaces.IterationResult
	preResult = nil
	if v.preChecker != nil {
		preResult = v.preChecker.GetResult()
	}
	if preResult == nil {
		return result
	}
	if result == nil {
		return preResult
	}
	// result and preResult have values
	if result.ReorgDetected && preResult.ReorgDetected {
		// That is the common case, checker must detect oldest blocks than preChecker
		if result.BlockNumber < preResult.BlockNumber {
			return result
		}
		return preResult
	}
	if preResult.ReorgDetected {
		return preResult
	}
	return result
}

func (v *L1BlockCheckerIntegration) onFinishChecker() {
	log.Infof("%s L1BlockChecker: finished background process, calling to synchronizer", logPrefix)
	// Stop both processes
	v.checker.Stop()
	if v.preChecker != nil {
		v.preChecker.Stop()
	}
	v.sync.OnDetectedMismatchL1BlockReorg()
}

func (v *L1BlockCheckerIntegration) launch(ctx context.Context) {
	log.Infof("%s L1BlockChecker: starting background process...", logPrefix)
	v.checker.Run(ctx, v.onFinishChecker)
	if v.preChecker != nil {
		log.Infof("%s L1BlockChecker: starting background precheck process...", logPrefix)
		v.preChecker.Run(ctx, v.onFinishChecker)
	}
}

func (v *L1BlockCheckerIntegration) executeResult(ctx context.Context, result syncinterfaces.IterationResult) bool {
	if result.ReorgDetected {
		for {
			err := v.sync.ExecuteReorgFromMismatchBlock(result.BlockNumber, result.ReorgMessage)
			if err == nil {
				return true
			}
			log.Errorf("%s Error executing reorg: %s", logPrefix, err)
			time.Sleep(v.timeBetweenRetries)
		}
	}
	return false
}
