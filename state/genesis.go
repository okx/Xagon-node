package state

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// Genesis contains the information to populate state on creation
type Genesis struct {
	// RollupBlockNumber is the block number where the polygonZKEVM smc was deployed on L1
	RollupBlockNumber uint64
	// RollupManagerBlockNumber is the block number where the RollupManager smc was deployed on L1
	RollupManagerBlockNumber uint64
	// Root hash of the genesis block
	Root common.Hash
	// Actions is the data to populate into the state trie
	Actions []*GenesisAction
}

// GenesisAction represents one of the values set on the SMT during genesis.
type GenesisAction struct {
	Address         string `json:"address"`
	Type            int    `json:"type"`
	StoragePosition string `json:"storagePosition"`
	Bytecode        string `json:"bytecode"`
	Key             string `json:"key"`
	Value           string `json:"value"`
	Root            string `json:"root"`
}

// SetGenesis populates state with genesis information
func (s *State) SetGenesis(ctx context.Context, block Block, genesis Genesis, m metrics.CallerLabel, dbTx pgx.Tx) (common.Hash, error) {
	var (
		root             common.Hash
		genesisStateRoot []byte
		err              error
	)
	if dbTx == nil {
		return common.Hash{}, ErrDBTxNil
	}
	if s.tree == nil {
		return common.Hash{}, ErrStateTreeNil
	}

	uuid := uuid.New().String()

	err = s.tree.StartBlock(ctx, root, uuid)
	if err != nil {
		log.Errorf("error starting block before genesis: %v", err)
		return common.Hash{}, err
	}

	for _, action := range genesis.Actions {
		address := common.HexToAddress(action.Address)
		switch action.Type {
		case int(merkletree.LeafTypeBalance):
			balance, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return common.Hash{}, err
			}
			genesisStateRoot, _, err = s.tree.SetBalance(ctx, address, balance, genesisStateRoot, uuid)
			if err != nil {
				return common.Hash{}, err
			}
		case int(merkletree.LeafTypeNonce):
			nonce, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return common.Hash{}, err
			}
			genesisStateRoot, _, err = s.tree.SetNonce(ctx, address, nonce, genesisStateRoot, uuid)
			if err != nil {
				return common.Hash{}, err
			}
		case int(merkletree.LeafTypeCode):
			code, err := hex.DecodeHex(action.Bytecode)
			if err != nil {
				return common.Hash{}, fmt.Errorf("could not decode SC bytecode for address %q: %v", address, err)
			}
			genesisStateRoot, _, err = s.tree.SetCode(ctx, address, code, genesisStateRoot, uuid)
			if err != nil {
				return common.Hash{}, err
			}
		case int(merkletree.LeafTypeStorage):
			// Parse position and value
			positionBI, err := encoding.DecodeBigIntHexOrDecimal(action.StoragePosition)
			if err != nil {
				return common.Hash{}, err
			}
			valueBI, err := encoding.DecodeBigIntHexOrDecimal(action.Value)
			if err != nil {
				return common.Hash{}, err
			}
			// Store
			genesisStateRoot, _, err = s.tree.SetStorageAt(ctx, address, positionBI, valueBI, genesisStateRoot, uuid)
			if err != nil {
				return common.Hash{}, err
			}
		case int(merkletree.LeafTypeSCLength):
			log.Debug("Skipped genesis action of type merkletree.LeafTypeSCLength, these actions will be handled as part of merkletree.LeafTypeCode actions")
		default:
			return common.Hash{}, fmt.Errorf("unknown genesis action type %q", action.Type)
		}
	}

	root.SetBytes(genesisStateRoot)

	err = s.tree.FinishBlock(ctx, root, uuid)
	if err != nil {
		log.Errorf("error finishing block after genesis: %v", err)
		return common.Hash{}, err
	}

	// flush state db
	err = s.tree.Flush(ctx, root, uuid)
	if err != nil {
		log.Errorf("error flushing state tree after genesis: %v", err)
		return common.Hash{}, err
	}

	// store L1 block related to genesis batch
	err = s.AddBlock(ctx, &block, dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// store genesis batch
	batch := Batch{
		BatchNumber:    0,
		Coinbase:       ZeroAddress,
		BatchL2Data:    nil,
		StateRoot:      root,
		LocalExitRoot:  ZeroHash,
		Timestamp:      block.ReceivedAt,
		Transactions:   []types.Transaction{},
		GlobalExitRoot: ZeroHash,
		ForcedBatchNum: nil,
	}

	err = s.StoreGenesisBatch(ctx, batch, string(SyncGenesisBatchClosingReason), dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// mark the genesis batch as virtualized
	virtualBatch := &VirtualBatch{
		BatchNumber: batch.BatchNumber,
		TxHash:      ZeroHash,
		Coinbase:    ZeroAddress,
		BlockNumber: block.BlockNumber,
	}
	forkID := s.GetForkIDByBatchNumber(0)
	if forkID >= FORKID_ETROG {
		virtualBatch.TimestampBatchEtrog = &block.ReceivedAt
	}

	err = s.AddVirtualBatch(ctx, virtualBatch, dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// mark the genesis batch as verified/consolidated
	verifiedBatch := &VerifiedBatch{
		BatchNumber: batch.BatchNumber,
		TxHash:      ZeroHash,
		Aggregator:  ZeroAddress,
		BlockNumber: block.BlockNumber,
	}
	err = s.AddVerifiedBatch(ctx, verifiedBatch, dbTx)
	if err != nil {
		return common.Hash{}, err
	}

	// store L2 genesis block
	header := NewL2Header(&types.Header{
		Number:     big.NewInt(0),
		ParentHash: ZeroHash,
		Coinbase:   ZeroAddress,
		Root:       root,
		Time:       uint64(block.ReceivedAt.Unix()),
	})
	rootHex := root.Hex()
	log.Info("Genesis root ", rootHex)

	receipts := []*types.Receipt{}
	st := trie.NewStackTrie(nil)
	l2Block := NewL2Block(header, []*types.Transaction{}, []*L2Header{}, receipts, st)
	l2Block.ReceivedAt = block.ReceivedAt

	// Sanity check
	if len(l2Block.Transactions()) > 0 {
		return common.Hash{}, fmt.Errorf("genesis L2Block contains %d transactions and should have 0", len(l2Block.Transactions()))
	}

	storeTxsEGPData := []StoreTxEGPData{}
	txsL2Hash := []common.Hash{}

	err = s.AddL2Block(ctx, batch.BatchNumber, l2Block, receipts, txsL2Hash, storeTxsEGPData, []common.Hash{}, dbTx)
	if err != nil {
		return common.Hash{}, err
	}
	return root, nil
}
