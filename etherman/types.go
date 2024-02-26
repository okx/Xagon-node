package etherman

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/oldpolygonzkevm"
	polygonzkevm "github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonvalidium_xlayer"
	"github.com/ethereum/go-ethereum/common"
)

// Block struct
type Block struct {
	BlockNumber           uint64
	BlockHash             common.Hash
	ParentHash            common.Hash
	ForcedBatches         []ForcedBatch
	SequencedBatches      [][]SequencedBatch
	UpdateEtrogSequence   UpdateEtrogSequence
	VerifiedBatches       []VerifiedBatch
	SequencedForceBatches [][]SequencedForceBatch
	ForkIDs               []ForkID
	ReceivedAt            time.Time
	// GER data
	GlobalExitRoots, L1InfoTree []GlobalExitRoot
}

// GlobalExitRoot struct
type GlobalExitRoot struct {
	BlockNumber       uint64
	MainnetExitRoot   common.Hash
	RollupExitRoot    common.Hash
	GlobalExitRoot    common.Hash
	Timestamp         time.Time
	PreviousBlockHash common.Hash
}

// SequencedBatch represents virtual batch
type SequencedBatch struct {
	BatchNumber   uint64
	L1InfoRoot    *common.Hash
	SequencerAddr common.Address
	TxHash        common.Hash
	Nonce         uint64
	Coinbase      common.Address
	// Struct used in preEtrog forks
	*oldpolygonzkevm.PolygonZkEVMBatchData
	// Struct used in Etrog
	*polygonzkevm.PolygonRollupBaseEtrogBatchData
}

// UpdateEtrogSequence represents the first etrog sequence
type UpdateEtrogSequence struct {
	BatchNumber   uint64
	SequencerAddr common.Address
	TxHash        common.Hash
	Nonce         uint64
	// Struct used in Etrog
	*polygonzkevm.PolygonRollupBaseEtrogBatchData
}

// ForcedBatch represents a ForcedBatch
type ForcedBatch struct {
	BlockNumber       uint64
	ForcedBatchNumber uint64
	Sequencer         common.Address
	GlobalExitRoot    common.Hash
	RawTxsData        []byte
	ForcedAt          time.Time
}

// VerifiedBatch represents a VerifiedBatch
type VerifiedBatch struct {
	BlockNumber uint64
	BatchNumber uint64
	Aggregator  common.Address
	StateRoot   common.Hash
	TxHash      common.Hash
}

// SequencedForceBatch is a sturct to track the ForceSequencedBatches event.
type SequencedForceBatch struct {
	BatchNumber uint64
	Coinbase    common.Address
	TxHash      common.Hash
	Timestamp   time.Time
	Nonce       uint64
	polygonzkevm.PolygonRollupBaseEtrogBatchData
}

// ForkID is a sturct to track the ForkID event.
type ForkID struct {
	BatchNumber uint64
	ForkID      uint64
	Version     string
}
