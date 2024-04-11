package aggregator

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

type proverInterface interface {
	Name() string
	ID() string
	Addr() string
	IsIdle() (bool, error)
	BatchProof(input *prover.InputProver) (*string, error)
	AggregatedProof(inputProof1, inputProof2 string) (*string, error)
	FinalProof(inputProof string, aggregatorAddr string) (*string, error)
	WaitRecursiveProof(ctx context.Context, proofID string) (string, error)
	WaitFinalProof(ctx context.Context, proofID string) (*prover.FinalProof, error)
}

// ethTxManager contains the methods required to send txs to
// ethereum.
type ethTxManager interface {
	Add(ctx context.Context, owner, id string, from common.Address, to *common.Address, value *big.Int, data []byte, gasOffset uint64, dbTx pgx.Tx) error
	Result(ctx context.Context, owner, id string, dbTx pgx.Tx) (ethtxmanager.MonitoredTxResult, error)
	ResultsByStatus(ctx context.Context, owner string, statuses []ethtxmanager.MonitoredTxStatus, dbTx pgx.Tx) ([]ethtxmanager.MonitoredTxResult, error)
	ProcessPendingMonitoredTxs(ctx context.Context, owner string, failedResultHandler ethtxmanager.ResultHandler, dbTx pgx.Tx)
}

// etherman contains the methods required to interact with ethereum
type etherman interface {
	GetRollupId() uint32
	GetLatestVerifiedBatchNum() (uint64, error)
	BuildTrustedVerifyBatchesTxData(lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs, beneficiary common.Address) (to *common.Address, data []byte, err error)
	GetLatestBlockHeader(ctx context.Context) (*types.Header, error)
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(context.Context, *big.Int) (bool, error)
}

// stateInterface gathers the methods to interact with the state.
type stateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	CheckProofContainsCompleteSequences(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) (bool, error)
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetProofReadyToVerify(ctx context.Context, lastVerfiedBatchNumber uint64, dbTx pgx.Tx) (*state.Proof, error)
	GetVirtualBatchToProve(ctx context.Context, lastVerfiedBatchNumber uint64, maxL1Block uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetProofsToAggregate(ctx context.Context, dbTx pgx.Tx) (*state.Proof, *state.Proof, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	AddGeneratedProof(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) error
	UpdateGeneratedProof(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) error
	DeleteGeneratedProofs(ctx context.Context, batchNumber uint64, batchNumberFinal uint64, dbTx pgx.Tx) error
	DeleteUngeneratedProofs(ctx context.Context, dbTx pgx.Tx) error
	CleanupGeneratedProofs(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	CleanupLockedProofs(ctx context.Context, duration string, dbTx pgx.Tx) (int64, error)
	GetL1InfoRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error)
	GetLeavesByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]state.L1InfoTreeExitRootStorageEntry, error)
	GetVirtualBatchParentHash(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetForcedBatchParentHash(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetVirtualBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.VirtualBatch, error)
}
