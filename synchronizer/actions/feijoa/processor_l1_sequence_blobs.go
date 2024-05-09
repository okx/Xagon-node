package feijoa

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	commonsync "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// stateProcessorSequenceBlobsInterface interface required from state
type stateProcessorSequenceBlobsInterface interface {
	AddBlobSequence(ctx context.Context, blobSequence *state.BlobSequence, dbTx pgx.Tx) error
	GetLastBlobSequence(ctx context.Context, dbTx pgx.Tx) (*state.BlobSequence, error)
	AddBlobInner(ctx context.Context, blobInner *state.BlobInner, dbTx pgx.Tx) error
	GetL1InfoRecursiveRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error)
}

type stateBlobInnerProcessor interface {
	ProcessBlobInner(ctx context.Context, request state.ProcessBlobInnerProcessRequest, data []byte) (*state.ProcessBlobInnerResponse, error)
}

// ProcessorSequenceBlobs processor for SequenceBlobs
type ProcessorSequenceBlobs struct {
	actions.ProcessorBase[ProcessorL1InfoTreeUpdate]
	state                   stateProcessorSequenceBlobsInterface
	stateBlobInnerProcessor stateBlobInnerProcessor
	timeProvider            commonsync.TimeProvider
}

// NewProcessorSequenceBlobs new processor for SequenceBlobs
func NewProcessorSequenceBlobs(state stateProcessorSequenceBlobsInterface, stateBlobInnerProcessor stateBlobInnerProcessor, timeProvider commonsync.TimeProvider) *ProcessorSequenceBlobs {
	if timeProvider == nil {
		timeProvider = &commonsync.DefaultTimeProvider{}
	}
	return &ProcessorSequenceBlobs{
		ProcessorBase: *actions.NewProcessorBase[ProcessorL1InfoTreeUpdate](
			[]etherman.EventOrder{etherman.SequenceBlobsOrder},
			actions.ForksIdOnlyFeijoa),
		state:                   state,
		stateBlobInnerProcessor: stateBlobInnerProcessor,
		timeProvider:            timeProvider,
	}
}

// Process process event
// - Store BlobSequence
// - Split BlobInner into Batches (executor)
// - Store BlobInner
func (p *ProcessorSequenceBlobs) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	seqBlobs := &l1Block.SequenceBlobs[order.Pos]
	previousBlobSequence, newBlobSequence, err := p.doBlobSequence(ctx, seqBlobs, l1Block, dbTx)
	if err != nil {
		return err
	}

	for idx := range seqBlobs.Blobs {
		blobNum := newBlobSequence.FirstBlobSequenced + uint64(idx)
		log.Infof("Blob %d: blobNum:%d", idx, blobNum)
		err := p.doBlobInner(ctx, blobNum, &seqBlobs.Blobs[idx], newBlobSequence, previousBlobSequence, dbTx)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *ProcessorSequenceBlobs) doBlobInner(ctx context.Context, blobNum uint64, blob *etherman.SequenceBlob, newBlobSequence, previousBlobSequence *state.BlobSequence, dbTx pgx.Tx) error {
	// TODO: We have to choose which tree depending on ForkID?
	leaf, err := p.state.GetL1InfoRecursiveRootLeafByIndex(ctx, blob.Params.L1InfoLeafIndex, dbTx)
	if err != nil {
		return err
	}

	stateBlob, err := p.convertToStateBlobInner(blob, blobNum, newBlobSequence.BlobSequenceIndex, leaf.L1InfoTreeRoot)
	if err != nil {
		log.Errorf("Error converting blob to state: %v", err)
		return err
	}

	processRequest, err := state.NewProcessBlobInnerProcessRequest(uint64(actions.ForkIDFeijoa), stateBlob, previousBlobSequence, *newBlobSequence)
	if err != nil {
		return err
	}
	log.Infof("storing Blob %d: BlobInner: %v", blobNum, stateBlob)
	err = p.state.AddBlobInner(ctx, stateBlob, dbTx)
	if err != nil {
		log.Errorf("Error storing blobInner to state: %v", err)
		return err
	}
	response, err := p.stateBlobInnerProcessor.ProcessBlobInner(ctx, *processRequest, blob.Data)
	if err != nil {
		return err
	}
	if response == nil {
		return errors.New("response is nil")
	}
	log.Infof("Blob %d: response: %v", blobNum, response)
	if response.IsSuccessfulExecution() {
		// We need to store the batches
		outcomeData := response.GetSuccesfulData()
		for idx := 0; idx < outcomeData.HowManyBatches(); idx++ {
			log.Infof("storing Blob %d: Batch %d: Hash:%s", blobNum, idx, outcomeData.GetBatchHash(idx).String())
			// TODO: Store batch
		}
	} else {
		err := response.GetUnifiedError()
		log.Errorf("Blob %d: response is not successful: Err: %s", blobNum, err.Error())
		return err
	}

	return nil
}

// returns previousBlobSequence and new one
func (p *ProcessorSequenceBlobs) doBlobSequence(ctx context.Context,
	incommingSequenceBlobs *etherman.SequenceBlobs, l1Block *etherman.Block, dbTx pgx.Tx) (*state.BlobSequence, *state.BlobSequence, error) {
	previousBlobSequence, err := p.state.GetLastBlobSequence(ctx, dbTx)
	if err != nil {
		return nil, nil, err
	}
	blobSequenceIndex := p.calculateBlobSequenceIndex(previousBlobSequence)
	newBlobSequence := p.convertToStateBlobSequence(incommingSequenceBlobs, blobSequenceIndex, l1Block.ReceivedAt, p.timeProvider.Now(), l1Block.BlockNumber)
	log.Infof("storing BlobSequence: %v", newBlobSequence)
	err = p.state.AddBlobSequence(ctx, newBlobSequence, dbTx)
	if err != nil {
		return nil, nil, err
	}
	return previousBlobSequence, newBlobSequence, nil
}

func (p *ProcessorSequenceBlobs) calculateBlobSequenceIndex(previousBlobSequence *state.BlobSequence) uint64 {
	nextIndex := uint64(1)
	if previousBlobSequence != nil {
		nextIndex = previousBlobSequence.BlobSequenceIndex + 1
	}
	return nextIndex
}

func (p *ProcessorSequenceBlobs) convertToStateBlobInner(blobInner *etherman.SequenceBlob, blobInnerNum uint64, blobSequenceIndex uint64, l1InfoTreeRoot common.Hash) (*state.BlobInner, error) {
	res := &state.BlobInner{
		BlobSequenceIndex:    blobSequenceIndex,
		BlobInnerNum:         blobInnerNum, // ho trect del previousBlobSequence
		Type:                 p.convertBlobType(blobInner.Type),
		MaxSequenceTimestamp: time.Unix(int64(blobInner.Params.MaxSequenceTimestamp), 0),
		ZkGasLimit:           blobInner.Params.ZkGasLimit,
		L1InfoLeafIndex:      blobInner.Params.L1InfoLeafIndex,
		L1InfoTreeRoot:       l1InfoTreeRoot,
	}
	if res.Type == state.TypeBlobTransaction {
		if blobInner.BlobBlobTypeParams == nil {
			return nil, errors.New("BlobBlobTypeParams from etherman is required for BlobTransaction")
		}
		res.BlobBlobTypeParams = &state.BlobBlobTypeParams{
			BlobIndex:  blobInner.BlobBlobTypeParams.BlobIndex.Uint64(),
			Z:          blobInner.BlobBlobTypeParams.Z,
			Y:          blobInner.BlobBlobTypeParams.Y,
			Commitment: blobInner.BlobBlobTypeParams.Commitment,
			Proof:      blobInner.BlobBlobTypeParams.Proof,
		}
	}
	return res, nil
}

func (p *ProcessorSequenceBlobs) convertBlobType(value etherman.BlobType) state.BlobType {
	return state.BlobType(value)
}

func (p *ProcessorSequenceBlobs) convertToStateBlobSequence(etherSeqBlobs *etherman.SequenceBlobs,
	nextIndex uint64,
	createAt time.Time,
	receviedAt time.Time,
	l1BlockNumber uint64) *state.BlobSequence {
	return &state.BlobSequence{
		BlobSequenceIndex:  nextIndex,
		L2Coinbase:         etherSeqBlobs.L2Coinbase,
		FirstBlobSequenced: etherSeqBlobs.EventData.LastBlobSequenced - uint64(len(etherSeqBlobs.Blobs)),
		LastBlobSequenced:  etherSeqBlobs.EventData.LastBlobSequenced,
		FinalAccInputHash:  etherSeqBlobs.FinalAccInputHash,
		CreateAt:           createAt,
		ReceivedAt:         receviedAt,
		BlockNumber:        l1BlockNumber,
	}
}
