package l1_parallel_sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	syncCommon "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
)

const (
	minNumIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData = 5
	minAcceptableTimeWaitingForNewRollupInfoData                       = 1 * time.Second
)

var (
	errContextCanceled                      = errors.New("consumer:context canceled")
	errConsumerStopped                      = errors.New("consumer:stopped by request")
	errConsumerStoppedBecauseIsSynchronized = errors.New("consumer:stopped because is synchronized")
	errConsumerAndProducerDesynchronized    = errors.New("consumer: consumer and producer are desynchronized")
)

// ConfigConsumer configuration for L1 sync parallel consumer
type ConfigConsumer struct {
	ApplyAfterNumRollupReceived int
	AceptableInacctivityTime    time.Duration
}

// synchronizerProcessBlockRangeInterface is the interface with synchronizer
// to execute blocks. This interface is used to mock the synchronizer in the tests
type synchronizerProcessBlockRangeInterface interface {
	ProcessBlockRange(blocks []etherman.Block, order map[common.Hash][]etherman.Order) error
}

// l1RollupInfoConsumer is the object that process the rollup info data incomming from channel chIncommingRollupInfo
type l1RollupInfoConsumer struct {
	mutex                 sync.Mutex
	synchronizer          synchronizerProcessBlockRangeInterface
	chIncommingRollupInfo chan L1SyncMessage
	ctx                   context.Context
	statistics            l1RollupInfoConsumerStatistics
	lastEthBlockSynced    *state.Block // Have been written in DB
	lastEthBlockReceived  *state.Block // is a memory cache
	highestBlockProcessed uint64
}

// NewL1RollupInfoConsumer creates a new l1RollupInfoConsumer
func NewL1RollupInfoConsumer(cfg ConfigConsumer,
	synchronizer synchronizerProcessBlockRangeInterface, ch chan L1SyncMessage) *l1RollupInfoConsumer {
	if cfg.AceptableInacctivityTime < minAcceptableTimeWaitingForNewRollupInfoData {
		log.Warnf("consumer: the AceptableInacctivityTime is too low (%s) minimum recommended %s", cfg.AceptableInacctivityTime, minAcceptableTimeWaitingForNewRollupInfoData)
	}
	if cfg.ApplyAfterNumRollupReceived < minNumIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData {
		log.Warnf("consumer: the ApplyAfterNumRollupReceived is too low (%d) minimum recommended  %d", cfg.ApplyAfterNumRollupReceived, minNumIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData)
	}

	return &l1RollupInfoConsumer{
		synchronizer:          synchronizer,
		chIncommingRollupInfo: ch,
		statistics: l1RollupInfoConsumerStatistics{
			startTime: time.Now(),
			cfg:       cfg,
		},
		highestBlockProcessed: invalidBlockNumber,
	}
}

func (l *l1RollupInfoConsumer) Start(ctx context.Context, lastEthBlockSynced *state.Block) error {
	l.ctx = ctx
	l.lastEthBlockSynced = lastEthBlockSynced
	if l.highestBlockProcessed == invalidBlockNumber && lastEthBlockSynced != nil {
		log.Infof("consumer: Starting consumer. setting HighestBlockProcessed: %d (lastEthBlockSynced)", lastEthBlockSynced.BlockNumber)
		l.highestBlockProcessed = lastEthBlockSynced.BlockNumber
	}
	log.Infof("consumer: Starting consumer. HighestBlockProcessed: %d", l.highestBlockProcessed)
	l.statistics.onStart()
	err := l.step()
	for ; err == nil; err = l.step() {
	}
	if err != errConsumerStopped && err != errConsumerStoppedBecauseIsSynchronized {
		return err
	}
	// The errConsumerStopped||errConsumerStoppedBecauseIsSynchronized are not an error, so we return nil meaning that the process finished in a normal way
	return nil
}

func (l *l1RollupInfoConsumer) Reset(startingBlockNumber uint64) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.highestBlockProcessed = startingBlockNumber
	l.lastEthBlockSynced = nil
	l.statistics.onReset()
}

func (l *l1RollupInfoConsumer) step() error {
	l.statistics.onStartStep()
	var err error
	select {
	case <-l.ctx.Done():
		return errContextCanceled
	case rollupInfo := <-l.chIncommingRollupInfo:
		if rollupInfo.dataIsValid {
			err = l.processIncommingRollupInfoData(rollupInfo.data)
			if err != nil {
				log.Error("consumer: error processing package.RollupInfoData. Error: ", err)
			}
		}
		if rollupInfo.ctrlIsValid {
			err = l.processIncommingRollupControlData(rollupInfo.ctrl)
			if err != nil && !errors.Is(err, errConsumerStoppedBecauseIsSynchronized) && !errors.Is(err, errConsumerStopped) {
				log.Error("consumer: error processing package.ControlData. Error: ", err)
			}
			log.Infof("consumer: processed ControlData[%s]. Result: %s", rollupInfo.ctrl.String(), err)
		}
	}
	return err
}
func (l *l1RollupInfoConsumer) processIncommingRollupControlData(control l1ConsumerControl) error {
	log.Debugf("consumer: processing controlPackage: %s", control.String())
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if control.event == eventStop {
		log.Infof("consumer: received a stop, so it stops processing. ignoring rest of items on channel len=%d", len(l.chIncommingRollupInfo))
		return errConsumerStopped
	}
	if control.event == eventProducerIsFullySynced {
		itemsInChannel := len(l.chIncommingRollupInfo)
		if itemsInChannel == 0 {
			consumerHigherBlockReceived := control.parameter
			log.Infof("consumer: received a fullSync and nothing pending in channel to process, so stopping consumer. lastBlock: %d", consumerHigherBlockReceived)
			if (l.highestBlockProcessed != invalidBlockNumber) && (l.highestBlockProcessed != consumerHigherBlockReceived) {
				log.Warnf("consumer: received a fullSync but highestBlockProcessed (%d) is not the same as consumerHigherBlockRequested (%d)", l.highestBlockProcessed, consumerHigherBlockReceived)
				return errConsumerAndProducerDesynchronized
			}
			return errConsumerStoppedBecauseIsSynchronized
		} else {
			log.Infof("consumer: received a fullSync but still have %d items in channel to process, so not stopping consumer", itemsInChannel)
		}
	}
	return nil
}

func checkPreviousBlocks(rollupInfo rollupInfoByBlockRangeResult, cachedBlock *state.Block) error {
	if cachedBlock == nil {
		return nil
	}
	if rollupInfo.previousBlockOfRange == nil {
		return nil
	}
	if cachedBlock.BlockNumber == rollupInfo.previousBlockOfRange.NumberU64() {
		if cachedBlock.BlockHash != rollupInfo.previousBlockOfRange.Hash() {
			err := fmt.Errorf("consumer: Previous block %d hash is not the same. state.Hash:%s != l1.Hash:%s",
				cachedBlock.BlockNumber, cachedBlock.BlockHash, rollupInfo.previousBlockOfRange.Hash())
			log.Errorf(err.Error())
			return syncCommon.NewReorgError(cachedBlock.BlockNumber, err)
		}

		log.Infof("consumer: Verified previous block %d  not the same: OK", cachedBlock.BlockNumber)
	}
	return nil
}

func (l *l1RollupInfoConsumer) processIncommingRollupInfoData(rollupInfo rollupInfoByBlockRangeResult) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	var err error
	if (l.highestBlockProcessed != invalidBlockNumber) && (l.highestBlockProcessed+1 != rollupInfo.blockRange.fromBlock) {
		log.Warnf("consumer: received a rollupInfo with a wrong block range.  Ignoring it. Highest block synced: %d. RollupInfo block range: %s",
			l.highestBlockProcessed, rollupInfo.blockRange.String())
		return nil
	}
	l.highestBlockProcessed = rollupInfo.getHighestBlockNumberInResponse()
	// Uncommented that line to produce a infinite loop of errors, and resets! (just for develop)
	//return errors.New("forcing an continuous error!")
	statisticsMsg := l.statistics.onStartProcessIncommingRollupInfoData(rollupInfo)
	log.Infof("consumer: processing rollupInfo #%000d: range:%s num_blocks [%d] highest_block [%d] statistics:%s", l.statistics.numProcessedRollupInfo, rollupInfo.blockRange.String(), len(rollupInfo.blocks), l.highestBlockProcessed, statisticsMsg)
	timeProcessingStart := time.Now()

	if l.lastEthBlockReceived != nil {
		err = checkPreviousBlocks(rollupInfo, l.lastEthBlockReceived)
		if err != nil {
			log.Errorf("consumer: error checking previous blocks: %s", err.Error())
			return err
		}
	}
	l.lastEthBlockReceived = rollupInfo.getHighestBlockReceived()

	lastBlockProcessed, err := l.processUnsafe(rollupInfo)
	if err == nil && lastBlockProcessed != nil {
		l.lastEthBlockSynced = lastBlockProcessed
	}
	l.statistics.onFinishProcessIncommingRollupInfoData(rollupInfo, time.Since(timeProcessingStart), err)
	if err != nil {
		log.Infof("consumer: error processing rollupInfo %s. Error: %s", rollupInfo.blockRange.String(), err.Error())
		return err
	}
	l.statistics.numProcessedBlocks += uint64(len(rollupInfo.blocks))
	return nil
}

// GetLastEthBlockSynced returns the last block synced, if true is returned, otherwise it returns false
func (l *l1RollupInfoConsumer) GetLastEthBlockSynced() (state.Block, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.lastEthBlockSynced == nil {
		return state.Block{}, false
	}
	return *l.lastEthBlockSynced, true
}

func (l *l1RollupInfoConsumer) StopAfterProcessChannelQueue() {
	log.Infof("consumer: Sending stop package: it will stop consumer (current channel len=%d)", len(l.chIncommingRollupInfo))
	l.sendStopPackage()
}

func (l *l1RollupInfoConsumer) sendStopPackage() {
	// Send a stop to the channel to stop the consumer when reach this point
	l.chIncommingRollupInfo <- *newL1SyncMessageControl(eventStop)
}

func (l *l1RollupInfoConsumer) processUnsafe(rollupInfo rollupInfoByBlockRangeResult) (*state.Block, error) {
	blocks := rollupInfo.blocks
	order := rollupInfo.order
	var lastEthBlockSynced *state.Block

	if len(blocks) == 0 {
		lb := rollupInfo.lastBlockOfRange
		if lb == nil {
			log.Info("consumer: Empty block range: ", rollupInfo.blockRange.String())
			return nil, nil
		}
		b := convertL1BlockToEthBlock(lb)
		err := l.synchronizer.ProcessBlockRange([]etherman.Block{b}, order)
		if err != nil {
			log.Error("consumer: Error processing last block of range: ", rollupInfo.blockRange, " err:", err)
			return nil, err
		}
		block := convertL1BlockToStateBlock(lb)
		lastEthBlockSynced = &block
		log.Debug("consumer: Storing empty block. BlockNumber: ", b.BlockNumber, ". BlockHash: ", b.BlockHash)
	} else {
		tmpStateBlock := convertEthmanBlockToStateBlock(&blocks[len(blocks)-1])
		lastEthBlockSynced = &tmpStateBlock
		logBlocks(blocks)
		err := l.synchronizer.ProcessBlockRange(blocks, order)
		if err != nil {
			log.Info("consumer: Error processing block range: ", rollupInfo.blockRange, " err:", err)
			return nil, err
		}
	}
	return lastEthBlockSynced, nil
}

func logBlocks(blocks []etherman.Block) {
	for i := range blocks {
		log.Debug("consumer: Position: [", i, "/", len(blocks), "] . BlockNumber: ", blocks[i].BlockNumber, ". BlockHash: ", blocks[i].BlockHash)
	}
}

func convertL1BlockToEthBlock(fb *types.Block) etherman.Block {
	return etherman.Block{
		BlockNumber: fb.NumberU64(),
		BlockHash:   fb.Hash(),
		ParentHash:  fb.ParentHash(),
		ReceivedAt:  time.Unix(int64(fb.Time()), 0),
	}
}

func convertL1BlockToStateBlock(fb *types.Block) state.Block {
	return state.Block{
		BlockNumber: fb.NumberU64(),
		BlockHash:   fb.Hash(),
		ParentHash:  fb.ParentHash(),
		ReceivedAt:  time.Unix(int64(fb.Time()), 0),
	}
}

func convertEthmanBlockToStateBlock(fb *etherman.Block) state.Block {
	return state.Block{
		BlockNumber: fb.BlockNumber,
		BlockHash:   fb.BlockHash,
		ParentHash:  fb.ParentHash,
		ReceivedAt:  fb.ReceivedAt,
	}
}
