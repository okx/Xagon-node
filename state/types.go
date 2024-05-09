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
	BatchNumber                uint64
	GlobalExitRoot_V1          common.Hash
	L1InfoRoot_V2              common.Hash
	L1InfoTreeData_V2          map[uint32]L1DataV2
	L1InfoTreeData_V3          map[uint32]L1DataV3
	OldStateRoot               common.Hash
	OldAccInputHash            common.Hash
	Transactions               []byte
	Coinbase                   common.Address
	ForcedBlockHashL1          common.Hash
	Timestamp_V1               time.Time
	TimestampLimit_V2          uint64
	Caller                     metrics.CallerLabel
	SkipFirstChangeL2Block_V2  bool
	SkipWriteBlockInfoRoot_V2  bool
	SkipVerifyL1InfoRoot_V2    bool
	ForkID                     uint64
	PreviousL1InfoTreeRoot_V3  common.Hash
	PreviousL1InfoTreeIndex_V3 uint32
}

// L1DataV2 represents the L1InfoTree data used in ProcessRequest.L1InfoTreeData_V2 parameter
type L1DataV2 struct {
	GlobalExitRoot common.Hash
	BlockHashL1    common.Hash
	MinTimestamp   uint64
	SmtProof       [][]byte
}

// L1DataV3 represents the L1InfoTree data used in ProcessRequest.L1InfoTreeData_V3 parameter
type L1DataV3 struct {
	GlobalExitRoot        common.Hash
	BlockHashL1           common.Hash
	MinTimestamp          uint64
	SmtProofPreviousIndex [][]byte
	InitialHistoricRoot   common.Hash
}

// ProcessBatchResponse represents the response of a batch process.
type ProcessBatchResponse struct {
	NewStateRoot       common.Hash
	NewAccInputHash    common.Hash
	NewLocalExitRoot   common.Hash
	NewBatchNumber     uint64
	UsedZkCounters     ZKCounters
	ReservedZkCounters ZKCounters
	// TransactionResponses_V1 []*ProcessTransactionResponse
	BlockResponses            []*ProcessBlockResponse
	ExecutorError             error
	ReadWriteAddresses        map[common.Address]*InfoReadWrite
	IsRomLevelError           bool
	IsExecutorLevelError      bool
	IsRomOOCError             bool
	FlushID                   uint64
	StoredFlushID             uint64
	ProverID                  string
	GasUsed_V2                uint64
	SMTKeys_V2                []merkletree.Key
	ProgramKeys_V2            []merkletree.Key
	ForkID                    uint64
	InvalidBatch_V2           bool
	RomError_V2               error
	OldStateRoot_V2           common.Hash
	NewLastTimestamp_V3       uint64
	CurrentL1InfoTreeRoot_V3  common.Hash
	CurrentL1InfoTreeIndex_V3 uint32
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
	// CumulativeGasUsed is the accumulated gas used (sum of tx GasUsed and CumulativeGasUsed of the previous tx in the L2 block)
	CumulativeGasUsed uint64
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
	// EffectivePercentage effective percentage used for the tx
	EffectivePercentage uint32
	// HasGaspriceOpcode flag to indicate if opcode 'GASPRICE' has been called
	HasGaspriceOpcode bool
	// HasBalanceOpcode flag to indicate if opcode 'BALANCE' has been called
	HasBalanceOpcode bool
	// Status of the transaction, 1 = success, 0 = failure
	Status uint32
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
	GasUsed          uint64
	KeccakHashes     uint32
	PoseidonHashes   uint32
	PoseidonPaddings uint32
	MemAligns        uint32
	Arithmetics      uint32
	Binaries         uint32
	Steps            uint32
	Sha256Hashes_V2  uint32
}

// SumUp sum ups zk counters with passed tx zk counters
func (z *ZKCounters) SumUp(other ZKCounters) {
	z.GasUsed += other.GasUsed
	z.KeccakHashes += other.KeccakHashes
	z.PoseidonHashes += other.PoseidonHashes
	z.PoseidonPaddings += other.PoseidonPaddings
	z.MemAligns += other.MemAligns
	z.Arithmetics += other.Arithmetics
	z.Binaries += other.Binaries
	z.Steps += other.Steps
	z.Sha256Hashes_V2 += other.Sha256Hashes_V2
}

// Fits checks if other zk counters fits in the zk counters. if there is a counter underflow it returns false and the name of the counter that caused the underflow
func (z *ZKCounters) Fits(other ZKCounters) (bool, string) {
	if other.GasUsed > z.GasUsed {
		return false, "CumulativeGas"
	}
	if other.KeccakHashes > z.KeccakHashes {
		return false, "KeccakHashes"
	}
	if other.PoseidonHashes > z.PoseidonHashes {
		return false, "PoseidonHashes"
	}
	if other.PoseidonPaddings > z.PoseidonPaddings {
		return false, "PoseidonPaddings"
	}
	if other.MemAligns > z.MemAligns {
		return false, "UsedMemAligns"
	}
	if other.Arithmetics > z.Arithmetics {
		return false, "UsedArithmetics"
	}
	if other.Binaries > z.Binaries {
		return false, "UsedBinaries"
	}
	if other.Steps > z.Steps {
		return false, "UsedSteps"
	}
	if other.Sha256Hashes_V2 > z.Sha256Hashes_V2 {
		return false, "UsedSha256Hashes_V2"
	}

	return true, ""
}

// Sub subtract zk counters with passed zk counters (not safe). if there is a counter underflow it returns true and the name of the counter that caused the underflow
func (z *ZKCounters) Sub(other ZKCounters) (bool, string) {
	if other.GasUsed > z.GasUsed {
		return true, "CumulativeGas"
	}
	if other.KeccakHashes > z.KeccakHashes {
		return true, "KeccakHashes"
	}
	if other.PoseidonHashes > z.PoseidonHashes {
		return true, "PoseidonHashes"
	}
	if other.PoseidonPaddings > z.PoseidonPaddings {
		return true, "PoseidonPaddings"
	}
	if other.MemAligns > z.MemAligns {
		return true, "UsedMemAligns"
	}
	if other.Arithmetics > z.Arithmetics {
		return true, "UsedArithmetics"
	}
	if other.Binaries > z.Binaries {
		return true, "UsedBinaries"
	}
	if other.Steps > z.Steps {
		return true, "UsedSteps"
	}
	if other.Sha256Hashes_V2 > z.Sha256Hashes_V2 {
		return true, "UsedSha256Hashes_V2"
	}

	z.GasUsed -= other.GasUsed
	z.KeccakHashes -= other.KeccakHashes
	z.PoseidonHashes -= other.PoseidonHashes
	z.PoseidonPaddings -= other.PoseidonPaddings
	z.MemAligns -= other.MemAligns
	z.Arithmetics -= other.Arithmetics
	z.Binaries -= other.Binaries
	z.Steps -= other.Steps
	z.Sha256Hashes_V2 -= other.Sha256Hashes_V2

	return false, ""
}

// BatchResources is a struct that contains the limited resources of a batch
type BatchResources struct {
	ZKCounters ZKCounters
	Bytes      uint64
}

// Fits check if the other batch resources fit in the batch resources. If there is a resource underflow it returns false and the name of the resource that caused the overflow
func (r *BatchResources) Fits(other BatchResources) (bool, string) {
	if other.Bytes > r.Bytes {
		return false, "Bytes"
	}
	return r.ZKCounters.Fits(other.ZKCounters)
}

// Sub subtracts the batch resources from "other". If there is a resource overflow it returns true and the name of the resource that caused the overflow
func (r *BatchResources) Sub(other BatchResources) (bool, string) {
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
