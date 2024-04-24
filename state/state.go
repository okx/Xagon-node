package state

import (
	"context"
	"math/big"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

const newL2BlockEventBufferSize = 500

var (
	// DefaultSenderAddress is the address that jRPC will use
	// to communicate with the state for eth_EstimateGas and eth_Call when
	// the From field is not specified because it is optional
	DefaultSenderAddress = "0x1111111111111111111111111111111111111111"
	// ZeroHash is the hash 0x0000000000000000000000000000000000000000000000000000000000000000
	ZeroHash = common.Hash{}
	// ZeroAddress is the address 0x0000000000000000000000000000000000000000
	ZeroAddress = common.Address{}
)

// State is an implementation of the state
type State struct {
	cfg Config
	storage
	executorClient executor.ExecutorServiceClient
	tree           *merkletree.StateTree
	eventLog       *event.EventLog
	l1InfoTree     *l1infotree.L1InfoTree

	newL2BlockEvents        chan NewL2BlockEvent
	newL2BlockEventHandlers []NewL2BlockEventHandler
}

// NewState creates a new State
func NewState(cfg Config, storage storage, executorClient executor.ExecutorServiceClient, stateTree *merkletree.StateTree, eventLog *event.EventLog, mt *l1infotree.L1InfoTree) *State {
	var once sync.Once
	once.Do(func() {
		metrics.Register()
	})

	state := &State{
		cfg:                     cfg,
		storage:                 storage,
		executorClient:          executorClient,
		tree:                    stateTree,
		eventLog:                eventLog,
		newL2BlockEvents:        make(chan NewL2BlockEvent, newL2BlockEventBufferSize),
		newL2BlockEventHandlers: []NewL2BlockEventHandler{},
		l1InfoTree:              mt,
	}

	return state
}

// StateTx is the state transaction that extends the database tx
type StateTx struct {
	pgx.Tx
	stateInstance      *State
	L1InfoTreeModified bool
}

// BeginStateTransaction starts a state transaction
func (s *State) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.Begin(ctx)
	if err != nil {
		return nil, err
	}
	res := &StateTx{
		Tx:            tx,
		stateInstance: s,
	}
	return res, nil
}

// Rollback do the dbTx rollback + modifications in cache mechanism
func (tx *StateTx) Rollback(ctx context.Context) error {
	if tx.L1InfoTreeModified {
		tx.stateInstance.ResetL1InfoTree()
	}
	return tx.Tx.Rollback(ctx)
}

// SetL1InfoTreeModified sets the flag to true to save that the L1InfoTree has been modified
func (tx *StateTx) SetL1InfoTreeModified() {
	tx.L1InfoTreeModified = true
}

// GetBalance from a given address
func (s *State) GetBalance(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	return s.tree.GetBalance(ctx, address, root.Bytes())
}

// GetCode from a given address
func (s *State) GetCode(ctx context.Context, address common.Address, root common.Hash) ([]byte, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	return s.tree.GetCode(ctx, address, root.Bytes())
}

// GetNonce returns the nonce of the given account at the given block number
func (s *State) GetNonce(ctx context.Context, address common.Address, root common.Hash) (uint64, error) {
	if s.tree == nil {
		return 0, ErrStateTreeNil
	}
	nonce, err := s.tree.GetNonce(ctx, address, root.Bytes())
	if err != nil {
		return 0, err
	}
	return nonce.Uint64(), nil
}

// GetStorageAt from a given address
func (s *State) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root common.Hash) (*big.Int, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	return s.tree.GetStorageAt(ctx, address, position, root.Bytes())
}

// GetLastStateRoot returns the latest state root
func (s *State) GetLastStateRoot(ctx context.Context, dbTx pgx.Tx) (common.Hash, error) {
	lastBlockHeader, err := s.GetLastL2BlockHeader(ctx, dbTx)
	if err != nil {
		return common.Hash{}, err
	}
	return lastBlockHeader.Root, nil
}

// GetBalanceByStateRoot gets balance from the MT Service using the provided state root
func (s *State) GetBalanceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	balance, err := s.tree.GetBalance(ctx, address, root.Bytes())
	if err != nil && balance == nil {
		balance = big.NewInt(0)
	}
	return balance, err
}

// GetNonceByStateRoot gets nonce from the MT Service using the provided state root
func (s *State) GetNonceByStateRoot(ctx context.Context, address common.Address, root common.Hash) (*big.Int, error) {
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	return s.tree.GetNonce(ctx, address, root.Bytes())
}

// GetTree returns State inner tree
func (s *State) GetTree() *merkletree.StateTree {
	return s.tree
}

// FlushMerkleTree persists updates in the Merkle tree
func (s *State) FlushMerkleTree(ctx context.Context, newStateRoot common.Hash) error {
	if s.tree == nil {
		return ErrStateTreeNil
	}
	return s.tree.Flush(ctx, newStateRoot, "")
}

// GetStoredFlushID returns the stored flush ID and Prover ID
func (s *State) GetStoredFlushID(ctx context.Context) (uint64, string, error) {
	if s.executorClient == nil {
		return 0, "", ErrExecutorNil
	}
	res, err := s.executorClient.GetFlushStatus(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, "", err
	}

	return res.StoredFlushId, res.ProverId, nil
}
