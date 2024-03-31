package dataavailability

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// DABackender is an interface for components that store and retrieve batch data
type DABackender interface {
	SequenceRetriever
	SequenceSender
	// Init initializes the DABackend
	Init() error
}

// SequenceSender is used to send provided sequence of batches
type SequenceSender interface {
	// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
	// as expected by the contract
	PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error)
}

// SequenceRetriever is used to retrieve batch data
type SequenceRetriever interface {
	// GetSequence retrieves the sequence data from the data availability backend
	GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error)
}

// === Internal interfaces ===

type stateInterface interface {
	GetBatchL2DataByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]byte, error)
	GetBatchL2DataByNumbers(ctx context.Context, batchNumbers []uint64, dbTx pgx.Tx) (map[uint64][]byte, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
}

// BatchDataProvider is used to retrieve batch data
type BatchDataProvider interface {
	// GetBatchL2Data retrieve the data of a batch from the DA backend. The returned data must be the pre-image of the hash
	GetBatchL2Data(batchNum []uint64, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error)
}

// DataManager is an interface for components that send and retrieve batch data
type DataManager interface {
	BatchDataProvider
	SequenceSender
}

// ZKEVMClientTrustedBatchesGetter contains the methods required to interact with zkEVM-RPC
type ZKEVMClientTrustedBatchesGetter interface {
	BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error)
	BatchesByNumbers(ctx context.Context, numbers []*big.Int) ([]*types.BatchData, error)
}
