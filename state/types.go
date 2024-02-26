package state

import (
	"encoding/json"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ProcessRequest represents the request of a batch process.
type ProcessRequest struct {
	BatchNumber               uint64
	GlobalExitRoot_V1         common.Hash
	L1InfoRoot_V2             common.Hash
	L1InfoTreeData_V2         map[uint32]L1DataV2
	OldStateRoot              common.Hash
	OldAccInputHash           common.Hash
	Transactions              []byte
	Coinbase                  common.Address
	ForcedBlockHashL1         common.Hash
	Timestamp_V1              time.Time
	TimestampLimit_V2         uint64
	Caller                    metrics.CallerLabel
	SkipFirstChangeL2Block_V2 bool
	SkipWriteBlockInfoRoot_V2 bool
	SkipVerifyL1InfoRoot_V2   bool
	ForkID                    uint64
	ExecutionMode             uint64
}

// L1DataV2 represents the L1InfoTree data used in ProcessRequest.L1InfoTreeData_V2 parameter
type L1DataV2 struct {
	GlobalExitRoot common.Hash
	BlockHashL1    common.Hash
	MinTimestamp   uint64
	SmtProof       [][]byte
}

// ProcessBatchResponse represents the response of a batch process.
type ProcessBatchResponse struct {
	NewStateRoot     common.Hash
	NewAccInputHash  common.Hash
	NewLocalExitRoot common.Hash
	NewBatchNumber   uint64
	UsedZkCounters   ZKCounters
	// TransactionResponses_V1 []*ProcessTransactionResponse
	BlockResponses       []*ProcessBlockResponse
	ExecutorError        error
	ReadWriteAddresses   map[common.Address]*InfoReadWrite
	IsRomLevelError      bool
	IsExecutorLevelError bool
	IsRomOOCError        bool
	FlushID              uint64
	StoredFlushID        uint64
	ProverID             string
	GasUsed_V2           uint64
	SMTKeys_V2           []merkletree.Key
	ProgramKeys_V2       []merkletree.Key
	ForkID               uint64
	InvalidBatch_V2      bool
	RomError_V2          error
	CloseBatch_V2        bool
}

// ProcessBlockResponse represents the response of a block
type ProcessBlockResponse struct {
	ParentHash           common.Hash
	Coinbase             common.Address
	GasLimit             uint64
	BlockNumber          uint64
	Timestamp            uint64
	GlobalExitRoot       common.Hash
	BlockHashL1          common.Hash
	GasUsed              uint64
	BlockInfoRoot        common.Hash
	BlockHash            common.Hash
	TransactionResponses []*ProcessTransactionResponse
	Logs                 []*types.Log
	RomError_V2          error
}

// ProcessTransactionResponse represents the response of a tx process.
type ProcessTransactionResponse struct {
	// TxHash is the hash of the transaction
	TxHash common.Hash
	// TxHashL2_V2 is the hash of the transaction in the L2
	TxHashL2_V2 common.Hash
	// Type indicates legacy transaction
	// It will be always 0 (legacy) in the executor
	Type uint32
	// ReturnValue is the returned data from the runtime (function result or data supplied with revert opcode)
	ReturnValue []byte
	// GasLeft is the total gas left as result of execution
	GasLeft uint64
	// GasUsed is the total gas used as result of execution or gas estimation
	GasUsed uint64
	// GasRefunded is the total gas refunded as result of execution
	GasRefunded uint64
	// RomError represents any error encountered during the execution
	RomError error
	// CreateAddress is the new SC Address in case of SC creation
	CreateAddress common.Address
	// StateRoot is the State Root
	StateRoot common.Hash
	// Logs emitted by LOG opcode
	Logs []*types.Log
	// ChangesStateRoot indicates if this tx affects the state
	ChangesStateRoot bool
	// Tx is the whole transaction object
	Tx types.Transaction
	// FullTrace contains the call trace.
	FullTrace instrumentation.FullTrace
	// EffectiveGasPrice effective gas price used for the tx
	EffectiveGasPrice string
	//EffectivePercentage effective percentage used for the tx
	EffectivePercentage uint32
	//HasGaspriceOpcode flag to indicate if opcode 'GASPRICE' has been called
	HasGaspriceOpcode bool
	//HasBalanceOpcode flag to indicate if opcode 'BALANCE' has been called
	HasBalanceOpcode bool
}

// EffectiveGasPriceLog contains all the data needed to calculate the effective gas price for logging purposes
type EffectiveGasPriceLog struct {
	Enabled        bool
	ValueFinal     *big.Int
	ValueFirst     *big.Int
	ValueSecond    *big.Int
	FinalDeviation *big.Int
	MaxDeviation   *big.Int
	GasUsedFirst   uint64
	GasUsedSecond  uint64
	GasPrice       *big.Int
	Percentage     uint8
	Reprocess      bool
	GasPriceOC     bool
	BalanceOC      bool
	L1GasPrice     uint64
	L2GasPrice     uint64
	Error          string
}

// StoreTxEGPData contains the data related to the effective gas price that needs to be stored when storing a tx
type StoreTxEGPData struct {
	EGPLog              *EffectiveGasPriceLog
	EffectivePercentage uint8
}

// ZKCounters counters for the tx
type ZKCounters struct {
	GasUsed              uint64
	UsedKeccakHashes     uint32
	UsedPoseidonHashes   uint32
	UsedPoseidonPaddings uint32
	UsedMemAligns        uint32
	UsedArithmetics      uint32
	UsedBinaries         uint32
	UsedSteps            uint32
	UsedSha256Hashes_V2  uint32
}

// SumUp sum ups zk counters with passed tx zk counters
func (z *ZKCounters) SumUp(other ZKCounters) {
	z.GasUsed += other.GasUsed
	z.UsedKeccakHashes += other.UsedKeccakHashes
	z.UsedPoseidonHashes += other.UsedPoseidonHashes
	z.UsedPoseidonPaddings += other.UsedPoseidonPaddings
	z.UsedMemAligns += other.UsedMemAligns
	z.UsedArithmetics += other.UsedArithmetics
	z.UsedBinaries += other.UsedBinaries
	z.UsedSteps += other.UsedSteps
	z.UsedSha256Hashes_V2 += other.UsedSha256Hashes_V2
}

// Sub subtract zk counters with passed zk counters (not safe). if there is a counter underflow it returns true and the name of the counter that caused the overflow
func (z *ZKCounters) Sub(other ZKCounters) (bool, string) {
	// ZKCounters
	if other.GasUsed > z.GasUsed {
		return true, "CumulativeGas"
	}
	if other.UsedKeccakHashes > z.UsedKeccakHashes {
		return true, "KeccakHashes"
	}
	if other.UsedPoseidonHashes > z.UsedPoseidonHashes {
		return true, "PoseidonHashes"
	}
	if other.UsedPoseidonPaddings > z.UsedPoseidonPaddings {
		return true, "PoseidonPaddings"
	}
	if other.UsedMemAligns > z.UsedMemAligns {
		return true, "UsedMemAligns"
	}
	if other.UsedArithmetics > z.UsedArithmetics {
		return true, "UsedArithmetics"
	}
	if other.UsedBinaries > z.UsedBinaries {
		return true, "UsedBinaries"
	}
	if other.UsedSteps > z.UsedSteps {
		return true, "UsedSteps"
	}
	if other.UsedSha256Hashes_V2 > z.UsedSha256Hashes_V2 {
		return true, "UsedSha256Hashes_V2"
	}

	z.GasUsed -= other.GasUsed
	z.UsedKeccakHashes -= other.UsedKeccakHashes
	z.UsedPoseidonHashes -= other.UsedPoseidonHashes
	z.UsedPoseidonPaddings -= other.UsedPoseidonPaddings
	z.UsedMemAligns -= other.UsedMemAligns
	z.UsedArithmetics -= other.UsedArithmetics
	z.UsedBinaries -= other.UsedBinaries
	z.UsedSteps -= other.UsedSteps
	z.UsedSha256Hashes_V2 -= other.UsedSha256Hashes_V2

	return false, ""
}

// BatchResources is a struct that contains the ZKEVM resources used by a batch/tx
type BatchResources struct {
	ZKCounters ZKCounters
	Bytes      uint64
}

// Sub subtracts the batch resources from other. if there is a resource underflow it returns true and the name of the resource that caused the overflow
func (r *BatchResources) Sub(other BatchResources) (bool, string) {
	// Bytes
	if other.Bytes > r.Bytes {
		return true, "Bytes"
	}
	bytesBackup := r.Bytes
	r.Bytes -= other.Bytes
	exhausted, resourceName := r.ZKCounters.Sub(other.ZKCounters)
	if exhausted {
		r.Bytes = bytesBackup
		return exhausted, resourceName
	}

	return false, ""
}

// SumUp sum ups the batch resources from other
func (r *BatchResources) SumUp(other BatchResources) {
	r.Bytes += other.Bytes
	r.ZKCounters.SumUp(other.ZKCounters)
}

// InfoReadWrite has information about modified addresses during the execution
type InfoReadWrite struct {
	Address common.Address
	Nonce   *uint64
	Balance *big.Int
}

// TraceConfig sets the debug configuration for the executor
type TraceConfig struct {
	DisableStorage   bool
	DisableStack     bool
	EnableMemory     bool
	EnableReturnData bool
	Tracer           *string
	TracerConfig     json.RawMessage

	Limit int // xlayer config
}

// IsDefaultTracer returns true when no custom tracer is set
func (t *TraceConfig) IsDefaultTracer() bool {
	return t.Tracer == nil || *t.Tracer == ""
}

// Is4ByteTracer returns true when should use 4byteTracer
func (t *TraceConfig) Is4ByteTracer() bool {
	return t.Tracer != nil && *t.Tracer == "4byteTracer"
}

// IsCallTracer returns true when should use callTracer
func (t *TraceConfig) IsCallTracer() bool {
	return t.Tracer != nil && *t.Tracer == "callTracer"
}

// IsNoopTracer returns true when should use noopTracer
func (t *TraceConfig) IsNoopTracer() bool {
	return t.Tracer != nil && *t.Tracer == "noopTracer"
}

// IsPrestateTracer returns true when should use prestateTracer
func (t *TraceConfig) IsPrestateTracer() bool {
	return t.Tracer != nil && *t.Tracer == "prestateTracer"
}

// IsJSCustomTracer returns true when should use js custom tracer
func (t *TraceConfig) IsJSCustomTracer() bool {
	return t.Tracer != nil && strings.Contains(*t.Tracer, "result") && strings.Contains(*t.Tracer, "fault")
}

// TrustedReorg represents a trusted reorg
type TrustedReorg struct {
	BatchNumber uint64
	Reason      string
}

// HexToAddressPtr create an address from a hex and returns its pointer
func HexToAddressPtr(hex string) *common.Address {
	a := common.HexToAddress(hex)
	return &a
}

// HexToHashPtr create a hash from a hex and returns its pointer
func HexToHashPtr(hex string) *common.Hash {
	h := common.HexToHash(hex)
	return &h
}
