package syncinterfaces

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// StateInterface contains the methods required to interact with the state.
type StateBeginTransactionInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
}

type StateGetBatchByNumberInterface interface {
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
}

type StateLastBlockGetter interface {
	GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error)
}

type StateBlobSequencer interface {
}

// StateFullInterface gathers the methods required to interact with the state.
type StateFullInterface interface {
	GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error)
	AddGlobalExitRoot(ctx context.Context, exitRoot *state.GlobalExitRoot, dbTx pgx.Tx) error
	AddForcedBatch(ctx context.Context, forcedBatch *state.ForcedBatch, dbTx pgx.Tx) error
	AddBlock(ctx context.Context, block *state.Block, dbTx pgx.Tx) error
	Reset(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error
	GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*state.Block, error)
	GetPreviousBlock(ctx context.Context, offset uint64, dbTx pgx.Tx) (*state.Block, error)
	GetFirstUncheckedBlock(ctx context.Context, fromBlockNumber uint64, dbTx pgx.Tx) (*state.Block, error)
	UpdateCheckedBlockByNumber(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx pgx.Tx) error
	GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	AddVirtualBatch(ctx context.Context, virtualBatch *state.VirtualBatch, dbTx pgx.Tx) error
	GetNextForcedBatches(ctx context.Context, nextForcedBatches int, dbTx pgx.Tx) ([]state.ForcedBatch, error)
	AddVerifiedBatch(ctx context.Context, verifiedBatch *state.VerifiedBatch, dbTx pgx.Tx) error
	ProcessAndStoreClosedBatch(ctx context.Context, processingCtx state.ProcessingContext, encodedTxs []byte, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
	ProcessAndStoreClosedBatchV2(ctx context.Context, processingCtx state.ProcessingContextV2, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
	SetGenesis(ctx context.Context, block state.Block, genesis state.Genesis, m metrics.CallerLabel, dbTx pgx.Tx) (common.Hash, error)
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ProcessBatch(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	ProcessBatchV2(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, string, error)
	StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, egpLog *state.EffectiveGasPriceLog, globalExitRoot, blockInfoRoot common.Hash, dbTx pgx.Tx) (*state.L2Header, error)
	GetStateRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error)
	ExecuteBatch(ctx context.Context, batch state.Batch, updateMerkleTree bool, dbTx pgx.Tx) (*executor.ProcessBatchResponse, error)
	ExecuteBatchV2(ctx context.Context, batch state.Batch, L1InfoTreeRoot common.Hash, l1InfoTreeData map[uint32]state.L1DataV2, timestampLimit time.Time, updateMerkleTree bool, skipVerifyL1InfoRoot uint32, forcedBlockHashL1 *common.Hash, dbTx pgx.Tx) (*executor.ProcessBatchResponseV2, error)
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	AddSequence(ctx context.Context, sequence state.Sequence, dbTx pgx.Tx) error
	AddAccumulatedInputHash(ctx context.Context, batchNum uint64, accInputHash common.Hash, dbTx pgx.Tx) error
	AddTrustedReorg(ctx context.Context, trustedReorg *state.TrustedReorg, dbTx pgx.Tx) error
	GetReorgedTransactions(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*ethTypes.Transaction, error)
	ResetForkID(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	GetForkIDs(ctx context.Context, dbTx pgx.Tx) ([]state.ForkIDInterval, error)
	AddForkIDInterval(ctx context.Context, newForkID state.ForkIDInterval, dbTx pgx.Tx) error
	SetLastBatchInfoSeenOnEthereum(ctx context.Context, lastBatchNumberSeen, lastBatchNumberVerified uint64, dbTx pgx.Tx) error
	SetInitSyncBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	GetForkIDByBlockNumber(blockNumber uint64) uint64
	GetStoredFlushID(ctx context.Context) (uint64, string, error)
	AddL1InfoTreeLeaf(ctx context.Context, L1InfoTreeLeaf *state.L1InfoTreeLeaf, dbTx pgx.Tx) (*state.L1InfoTreeExitRootStorageEntry, error)
	AddL1InfoTreeRecursiveLeaf(ctx context.Context, L1InfoTreeLeaf *state.L1InfoTreeLeaf, dbTx pgx.Tx) (*state.L1InfoTreeExitRootStorageEntry, error)
	StoreL2Block(ctx context.Context, batchNumber uint64, l2Block *state.ProcessBlockResponse, txsEGPLog []*state.EffectiveGasPriceLog, dbTx pgx.Tx) (common.Hash, error)
	GetL1InfoRootLeafByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error)
	UpdateWIPBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	GetL1InfoTreeDataFromBatchL2Data(ctx context.Context, batchL2Data []byte, dbTx pgx.Tx) (map[uint32]state.L1DataV2, common.Hash, common.Hash, error)
	GetExitRootByGlobalExitRoot(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (*state.GlobalExitRoot, error)
	GetForkIDInMemory(forkId uint64) *state.ForkIDInterval
	GetLastL2BlockByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.L2Block, error)
	UpdateForkIDBlockNumber(ctx context.Context, forkdID uint64, newBlockNumber uint64, updateMemCache bool, dbTx pgx.Tx) error
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetL2BlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*state.L2Block, error)
	GetLastBlobSequence(ctx context.Context, dbTx pgx.Tx) (*state.BlobSequence, error)
	AddBlobSequence(ctx context.Context, blobSequence *state.BlobSequence, dbTx pgx.Tx) error
	GetL1InfoRecursiveRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error)
	ProcessBlobInner(ctx context.Context, request state.ProcessBlobInnerProcessRequest, data []byte) (*state.ProcessBlobInnerResponse, error)
	AddBlobInner(ctx context.Context, blobInner *state.BlobInner, dbTx pgx.Tx) error
	GetUncheckedBlocks(ctx context.Context, fromBlockNumber uint64, toBlockNumber uint64, dbTx pgx.Tx) ([]*state.Block, error)
	GetPreviousBlockToBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*state.Block, error)
}
