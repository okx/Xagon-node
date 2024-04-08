package runtime

import (
	"encoding/json"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
)

var (

	// ROM ERRORS
	// ==========

	// ErrOutOfGas indicates there is not enough balance to continue the execution
	ErrOutOfGas = errors.New("out of gas")
	// ErrStackOverflow indicates a stack overflow has happened
	ErrStackOverflow = errors.New("stack overflow")
	// ErrStackUnderflow indicates a stack overflow has happened
	ErrStackUnderflow = errors.New("stack underflow")
	// ErrMaxCodeSizeExceeded indicates the code size is beyond the maximum
	ErrMaxCodeSizeExceeded = errors.New("evm: max code size exceeded")
	// ErrContractAddressCollision there is a collision regarding contract addresses
	ErrContractAddressCollision = errors.New("contract address collision")
	// ErrExecutionReverted indicates the execution has been reverted
	ErrExecutionReverted = errors.New("execution reverted")
	// ErrOutOfCountersStep indicates there are not enough step counters to continue the execution
	ErrOutOfCountersStep = errors.New("not enough step counters to continue the execution")
	// ErrOutOfCountersKeccak indicates there are not enough keccak counters to continue the execution
	ErrOutOfCountersKeccak = errors.New("not enough keccak counters to continue the execution")
	// ErrOutOfCountersBinary indicates there are not enough binary counters to continue the execution
	ErrOutOfCountersBinary = errors.New("not enough binary counters to continue the execution")
	// ErrOutOfCountersMemory indicates there are not enough memory align counters to continue the execution
	ErrOutOfCountersMemory = errors.New("not enough memory align counters to continue the execution")
	// ErrOutOfCountersArith indicates there are not enough arith counters to continue the execution
	ErrOutOfCountersArith = errors.New("not enough arith counters to continue the execution")
	// ErrOutOfCountersPadding indicates there are not enough padding counters to continue the execution
	ErrOutOfCountersPadding = errors.New("not enough padding counters to continue the execution")
	// ErrOutOfCountersPoseidon indicates there are not enough poseidon counters to continue the execution
	ErrOutOfCountersPoseidon = errors.New("not enough poseidon counters to continue the execution")
	// ErrOutOfCountersSha indicates there are not enough sha256 counters to continue the execution
	ErrOutOfCountersSha = errors.New("not enough sha256 counters to continue the execution")
	// ErrIntrinsicInvalidSignature indicates the transaction is failing at the signature intrinsic check
	ErrIntrinsicInvalidSignature = errors.New("signature intrinsic error")
	// ErrIntrinsicInvalidChainID indicates the transaction is failing at the chain id intrinsic check
	ErrIntrinsicInvalidChainID = errors.New("chain id intrinsic error")
	// ErrIntrinsicInvalidNonce indicates the transaction is failing at the nonce intrinsic check
	ErrIntrinsicInvalidNonce = errors.New("nonce intrinsic error")
	// ErrIntrinsicInvalidGasLimit indicates the transaction is failing at the gas limit intrinsic check
	ErrIntrinsicInvalidGasLimit = errors.New("gas limit intrinsic error")
	// ErrIntrinsicInvalidBalance indicates the transaction is failing at balance intrinsic check
	ErrIntrinsicInvalidBalance = errors.New("balance intrinsic error")
	// ErrIntrinsicInvalidBatchGasLimit indicates the batch is exceeding the batch gas limit
	ErrIntrinsicInvalidBatchGasLimit = errors.New("batch gas limit intrinsic error")
	// ErrIntrinsicInvalidSenderCode indicates the sender code is invalid
	ErrIntrinsicInvalidSenderCode = errors.New("invalid sender code intrinsic error")
	// ErrBatchDataTooBig indicates the batch_l2_data is too big to be processed
	ErrBatchDataTooBig = errors.New("batch data too big")
	// ErrInvalidJump indicates there is an invalid jump opcode
	ErrInvalidJump = errors.New("invalid jump opcode")
	// ErrInvalidOpCode indicates there is an invalid opcode
	ErrInvalidOpCode = errors.New("invalid opcode")
	// ErrInvalidStatic indicates there is an invalid static call
	ErrInvalidStatic = errors.New("invalid static call")
	// ErrInvalidByteCodeStartsEF indicates there is a byte code starting with 0xEF
	ErrInvalidByteCodeStartsEF = errors.New("byte code starting with 0xEF")
	// ErrIntrinsicInvalidTxGasOverflow indicates the transaction gasLimit*gasPrice > MAX_UINT_256 - 1
	ErrIntrinsicInvalidTxGasOverflow = errors.New("gas overflow")
	// ErrUnsupportedForkId indicates that the fork id is not supported
	ErrUnsupportedForkId = errors.New("unsupported fork id")
	// ErrInvalidRLP indicates that there has been an error while parsing the RLP
	ErrInvalidRLP = errors.New("invalid RLP")

	// Start of V2 errors

	// ErrInvalidDecodeChangeL2Block indicates that there has been an error while decoding a change l2 block transaction
	ErrInvalidDecodeChangeL2Block = errors.New("error while decoding a change l2 block transaction")
	// ErrInvalidNotFirstTxChangeL2Block indicates that there has been an error while decoding a create l2 block transaction
	ErrInvalidNotFirstTxChangeL2Block = errors.New("the first transaction in a batch is not a change l2 block transaction")
	// ErrInvalidTxChangeL2BlockLimitTimestamp indicates that the change l2 block transaction has trigger an error while executing
	ErrInvalidTxChangeL2BlockLimitTimestamp = errors.New("the change l2 block transaction has trigger an error while executing (limit timestamp)")
	// ErrInvalidTxChangeL2BlockMinTimestamp indicates that the change l2 block transaction has trigger an error while executing
	ErrInvalidTxChangeL2BlockMinTimestamp = errors.New("indicates that the change l2 block transaction has trigger an error while executing (min timestamp)")

	// EXECUTOR ERRORS
	// ===============

	// ErrExecutorDBError indicates that there is an error connecting to the database
	ErrExecutorDBError = errors.New("database error")
	// ErrExecutorSMMainCountersOverflowSteps indicates that the main execution exceeded the maximum number of steps
	ErrExecutorSMMainCountersOverflowSteps = errors.New("main execution exceeded the maximum number of steps")
	// ErrExecutorSMMainCountersOverflowKeccak indicates that the keccak counter exceeded the maximum
	ErrExecutorSMMainCountersOverflowKeccak = errors.New("keccak counter exceeded the maximum")
	// ErrExecutorSMMainCountersOverflowBinary indicates that the binary counter exceeded the maximum
	ErrExecutorSMMainCountersOverflowBinary = errors.New("binary counter exceeded the maximum")
	// ErrExecutorSMMainCountersOverflowMem indicates that the memory align counter exceeded the maximum
	ErrExecutorSMMainCountersOverflowMem = errors.New("memory align counter exceeded the maximum")
	// ErrExecutorSMMainCountersOverflowArith indicates that the arith counter exceeded the maximum
	ErrExecutorSMMainCountersOverflowArith = errors.New("arith counter exceeded the maximum")
	// ErrExecutorSMMainCountersOverflowPadding indicates that the padding counter exceeded the maximum
	ErrExecutorSMMainCountersOverflowPadding = errors.New("padding counter exceeded the maximum")
	// ErrExecutorSMMainCountersOverflowPoseidon indicates that the poseidon counter exceeded the maximum
	ErrExecutorSMMainCountersOverflowPoseidon = errors.New("poseidon counter exceeded the maximum")
	// ErrExecutorUnsupportedForkId indicates that the fork id is not supported
	ErrExecutorUnsupportedForkId = errors.New("unsupported fork id")
	// ErrExecutorBalanceMismatch indicates that there is a balance mismatch error in the ROM
	ErrExecutorBalanceMismatch = errors.New("balance mismatch")
	// ErrExecutorFEA2Scalar indicates that there is a fea2scalar error in the execution
	ErrExecutorFEA2Scalar = errors.New("fea2scalar error")
	// ErrExecutorTOS32 indicates that there is a TOS32 error in the execution
	ErrExecutorTOS32 = errors.New("TOS32 error")
	// ErrExecutorSMMainInvalidUnsignedTx indicates that there is an unsigned TX in a non-process batch (i.e. in a prover request)
	ErrExecutorSMMainInvalidUnsignedTx = errors.New("unsigned TX in a non-process batch")
	// ErrExecutorSMMainInvalidNoCounters indicates that there is a no-counters request in a non-process batch (i.e. in a prover request)
	ErrExecutorSMMainInvalidNoCounters = errors.New("no-counters request in a non-process batch")
	// ErrExecutorSMMainArithECRecoverDivideByZero indicates that there is a divide-by-zero situation during an ECRecover
	ErrExecutorSMMainArithECRecoverDivideByZero = errors.New("divide-by-zero situation during an ECRecover")
	// ErrExecutorSMMainAddressOutOfRange indicates that an address is out of valid memory space range
	ErrExecutorSMMainAddressOutOfRange = errors.New("address is out of valid memory space range")
	// ErrExecutorSMMainAddressNegative indicates that an address is negative
	ErrExecutorSMMainAddressNegative = errors.New("address is negative")
	// ErrExecutorSMMainStorageInvalidKey indicates that a register value is out of range while building storage key
	ErrExecutorSMMainStorageInvalidKey = errors.New("register value is out of range while building storage key")
	// ErrExecutorSMMainHashK indicates that a register value is out of range while calculating a Keccak hash
	ErrExecutorSMMainHashK = errors.New("register value is out of range while calculating a Keccak hash")
	// ErrExecutorSMMainHashKSizeOutOfRange indicates that a size register value is out of range while calculating a Keccak hash
	ErrExecutorSMMainHashKSizeOutOfRange = errors.New("size register value is out of range while calculating a Keccak hash")
	// ErrExecutorSMMainHashKPositionNegative indicates that a position register value is negative while calculating a Keccak hash
	ErrExecutorSMMainHashKPositionNegative = errors.New("position register value is negative while calculating a Keccak hash")
	// ErrExecutorSMMainHashKPositionPlusSizeOutOfRange indicates that a position register value plus a size register value is out of range while calculating a Keccak hash
	ErrExecutorSMMainHashKPositionPlusSizeOutOfRange = errors.New("position register value plus a size register value is out of range while calculating a Keccak hash")
	// ErrExecutorSMMainHashKDigestAddressNotFound indicates that an address has not been found while calculating a Keccak hash digest
	ErrExecutorSMMainHashKDigestAddressNotFound = errors.New("address has not been found while calculating a Keccak hash digest")
	// ErrExecutorSMMainHashKDigestNotCompleted indicates that the hash has not been completed while calling a Keccak hash digest
	ErrExecutorSMMainHashKDigestNotCompleted = errors.New("hash has not been completed while calling a Keccak hash digest")
	// ErrExecutorSMMainHashP indicates that a register value is out of range while calculating a Poseidon hash
	ErrExecutorSMMainHashP = errors.New("register value is out of range while calculating a Poseidon hash")
	// ErrExecutorSMMainHashPSizeOutOfRange indicates that a size register value is out of range while calculating a Poseidon hash
	ErrExecutorSMMainHashPSizeOutOfRange = errors.New("size register value is out of range while calculating a Poseidon hash")
	// ErrExecutorSMMainHashPPositionNegative indicates that a position register value is negative while calculating a Poseidon hash
	ErrExecutorSMMainHashPPositionNegative = errors.New("position register value is negative while calculating a Poseidon hash")
	// ErrExecutorSMMainHashPPositionPlusSizeOutOfRange indicates that a position register value plus a size register value is out of range while calculating a Poseidon hash
	ErrExecutorSMMainHashPPositionPlusSizeOutOfRange = errors.New("position register value plus a size register value is out of range while calculating a Poseidon hash")
	// ErrExecutorSMMainHashPDigestAddressNotFound indicates that an address has not been found while calculating a Poseidon hash digest
	ErrExecutorSMMainHashPDigestAddressNotFound = errors.New("address has not been found while calculating a Poseidon hash digest")
	// ErrExecutorSMMainHashPDigestNotCompleted indicates that the hash has not been completed while calling a Poseidon hash digest
	ErrExecutorSMMainHashPDigestNotCompleted = errors.New("hash has not been completed while calling a Poseidon hash digest")
	// ErrExecutorSMMainMemAlignOffsetOutOfRange indicates that the an offset register value is out of range while doing a mem align operation
	ErrExecutorSMMainMemAlignOffsetOutOfRange = errors.New("offset register value is out of range while doing a mem align operation")
	// ErrExecutorSMMainMultipleFreeIn indicates that we got more than one free inputs in one ROM instruction
	ErrExecutorSMMainMultipleFreeIn = errors.New("more than one free inputs in one ROM instruction")
	// ErrExecutorSMMainAssert indicates that the ROM assert instruction failed
	ErrExecutorSMMainAssert = errors.New("ROM assert instruction failed")
	// ErrExecutorSMMainMemory indicates that the memory instruction check failed
	ErrExecutorSMMainMemory = errors.New("memory instruction check failed")
	// ErrExecutorSMMainStorageReadMismatch indicates that the storage read instruction check failed
	ErrExecutorSMMainStorageReadMismatch = errors.New("storage read instruction check failed")
	// ErrExecutorSMMainStorageWriteMismatch indicates that the storage read instruction check failed
	ErrExecutorSMMainStorageWriteMismatch = errors.New("storage write instruction check failed")
	// ErrExecutorSMMainHashKValueMismatch indicates that the Keccak hash instruction value check failed
	ErrExecutorSMMainHashKValueMismatch = errors.New("keccak hash instruction value check failed")
	// ErrExecutorSMMainHashKPaddingMismatch indicates that the Keccak hash instruction padding check failed
	ErrExecutorSMMainHashKPaddingMismatch = errors.New("keccak hash instruction padding check failed")
	// ErrExecutorSMMainHashKSizeMismatch indicates that the Keccak hash instruction size check failed
	ErrExecutorSMMainHashKSizeMismatch = errors.New("keccak hash instruction check size failed")
	// ErrExecutorSMMainHashKLenLengthMismatch indicates that the Keccak hash length instruction length check failed
	ErrExecutorSMMainHashKLenLengthMismatch = errors.New("keccak hash length instruction length check failed")
	// ErrExecutorSMMainHashKLenCalledTwice indicates that the Keccak hash length instruction called once check failed
	ErrExecutorSMMainHashKLenCalledTwice = errors.New("keccak hash length instruction called once check failed")
	// ErrExecutorSMMainHashKDigestNotFound indicates that the Keccak hash digest instruction slot not found
	ErrExecutorSMMainHashKDigestNotFound = errors.New("keccak hash digest instruction slot not found")
	// ErrExecutorSMMainHashKDigestDigestMismatch indicates that the Keccak hash digest instruction digest check failed
	ErrExecutorSMMainHashKDigestDigestMismatch = errors.New("keccak hash digest instruction digest check failed")
	// ErrExecutorSMMainHashKDigestCalledTwice indicates that the Keccak hash digest instruction called once check failed
	ErrExecutorSMMainHashKDigestCalledTwice = errors.New("keccak hash digest instruction called once check failed")
	// ErrExecutorSMMainHashPValueMismatch indicates that the Poseidon hash instruction value check failed
	ErrExecutorSMMainHashPValueMismatch = errors.New("poseidon hash instruction value check failed")
	// ErrExecutorSMMainHashPPaddingMismatch indicates that the Poseidon hash instruction padding check failed
	ErrExecutorSMMainHashPPaddingMismatch = errors.New("poseidon hash instruction padding check failed")
	// ErrExecutorSMMainHashPSizeMismatch indicates that the Poseidon hash instruction size check failed
	ErrExecutorSMMainHashPSizeMismatch = errors.New("poseidon hash instruction size check failed")
	// ErrExecutorSMMainHashPLenLengthMismatch indicates that the Poseidon hash length instruction length check failed
	ErrExecutorSMMainHashPLenLengthMismatch = errors.New("poseidon hash length instruction length check failed")
	// ErrExecutorSMMainHashPLenCalledTwice indicates that the Poseidon hash length instruction called once check failed
	ErrExecutorSMMainHashPLenCalledTwice = errors.New("poseidon hash length instruction called once check failed")
	// ErrExecutorSMMainHashPDigestDigestMismatch indicates that the Poseidon hash digest instruction digest check failed
	ErrExecutorSMMainHashPDigestDigestMismatch = errors.New("poseidon hash digest instruction digest check failed")
	// ErrExecutorSMMainHashPDigestCalledTwice indicates that the Poseidon hash digest instruction called once check failed
	ErrExecutorSMMainHashPDigestCalledTwice = errors.New("poseidon hash digest instruction called once check failed")
	// ErrExecutorSMMainArithMismatch indicates that the arith instruction check failed
	ErrExecutorSMMainArithMismatch = errors.New("arith instruction check failed")
	// ErrExecutorSMMainArithECRecoverMismatch indicates that the arith ECRecover instruction check failed
	ErrExecutorSMMainArithECRecoverMismatch = errors.New("arith ECRecover instruction check failed")
	// ErrExecutorSMMainBinaryAddMismatch indicates that the binary add instruction check failed
	ErrExecutorSMMainBinaryAddMismatch = errors.New("binary add instruction check failed")
	// ErrExecutorSMMainBinarySubMismatch indicates that the binary sub instruction check failed
	ErrExecutorSMMainBinarySubMismatch = errors.New("binary sub instruction check failed")
	// ErrExecutorSMMainBinaryLtMismatch indicates that the binary less than instruction check failed
	ErrExecutorSMMainBinaryLtMismatch = errors.New("binary less than instruction check failed")
	// ErrExecutorSMMainBinarySLtMismatch indicates that the binary signed less than instruction check failed
	ErrExecutorSMMainBinarySLtMismatch = errors.New("binary signed less than instruction check failed")
	// ErrExecutorSMMainBinaryEqMismatch indicates that the binary equal instruction check failed
	ErrExecutorSMMainBinaryEqMismatch = errors.New("binary equal instruction check failed")
	// ErrExecutorSMMainBinaryAndMismatch indicates that the binary and instruction check failed
	ErrExecutorSMMainBinaryAndMismatch = errors.New("binary and instruction check failed")
	// ErrExecutorSMMainBinaryOrMismatch indicates that the binary or instruction check failed
	ErrExecutorSMMainBinaryOrMismatch = errors.New("binary or instruction check failed")
	// ErrExecutorSMMainBinaryXorMismatch indicates that the binary xor instruction check failed
	ErrExecutorSMMainBinaryXorMismatch = errors.New("binary xor instruction check failed")
	// ErrExecutorSMMainMemAlignWriteMismatch indicates that the memory align write instruction check failed
	ErrExecutorSMMainMemAlignWriteMismatch = errors.New("memory align write instruction check failed")
	// ErrExecutorSMMainMemAlignWrite8Mismatch indicates that the memory align write 8 instruction check failed
	ErrExecutorSMMainMemAlignWrite8Mismatch = errors.New("memory align write 8 instruction check failed")
	// ErrExecutorSMMainMemAlignReadMismatch indicates that the memory align read instruction check failed
	ErrExecutorSMMainMemAlignReadMismatch = errors.New("memory align read instruction check failed")
	// ErrExecutorSMMainJmpnOutOfRange indicates that the JMPN instruction found a jump position out of range
	ErrExecutorSMMainJmpnOutOfRange = errors.New("JMPN instruction found a jump position out of range")
	// ErrExecutorSMMainHashKReadOutOfRange indicates that the main execution Keccak check found read out of range
	ErrExecutorSMMainHashKReadOutOfRange = errors.New("main execution Keccak check found read out of range")
	// ErrExecutorSMMainHashPReadOutOfRange indicates that the main execution Poseidon check found read out of range
	ErrExecutorSMMainHashPReadOutOfRange = errors.New("main execution Poseidon check found read out of range")
	// ErrExecutorErrorInvalidOldStateRoot indicates that the input parameter old_state_root is invalid
	ErrExecutorErrorInvalidOldStateRoot = errors.New("old_state_root is invalid")
	// ErrExecutorErrorInvalidOldAccInputHash indicates that the input parameter old_acc_input_hash is invalid
	ErrExecutorErrorInvalidOldAccInputHash = errors.New("old_acc_input_hash is invalid")
	// ErrExecutorErrorInvalidChainId indicates that the input parameter chain_id is invalid
	ErrExecutorErrorInvalidChainId = errors.New("chain_id is invalid")
	// ErrExecutorErrorInvalidBatchL2Data indicates that the input parameter batch_l2_data is invalid
	ErrExecutorErrorInvalidBatchL2Data = errors.New("batch_l2_data is invalid")
	// ErrExecutorErrorInvalidGlobalExitRoot indicates that the input parameter global_exit_root is invalid
	ErrExecutorErrorInvalidGlobalExitRoot = errors.New("global_exit_root is invalid")
	// ErrExecutorErrorInvalidCoinbase indicates that the input parameter coinbase (i.e. sequencer address) is invalid
	ErrExecutorErrorInvalidCoinbase = errors.New("coinbase (i.e. sequencer address) is invalid")
	// ErrExecutorErrorInvalidFrom indicates that the input parameter from is invalid
	ErrExecutorErrorInvalidFrom = errors.New("from is invalid")
	// ErrExecutorErrorInvalidDbKey indicates that the input parameter db key is invalid
	ErrExecutorErrorInvalidDbKey = errors.New("db key is invalid")
	// ErrExecutorErrorInvalidDbValue indicates that the input parameter db value is invalid
	ErrExecutorErrorInvalidDbValue = errors.New("db value is invalid")
	// ErrExecutorErrorInvalidContractsBytecodeKey indicates that the input parameter contracts_bytecode key is invalid
	ErrExecutorErrorInvalidContractsBytecodeKey = errors.New("contracts_bytecode key is invalid")
	// ErrExecutorErrorInvalidContractsBytecodeValue indicates that the input parameter contracts_bytecode value is invalid
	ErrExecutorErrorInvalidContractsBytecodeValue = errors.New("contracts_bytecode value is invalid")
	// ErrExecutorErrorInvalidGetKey indicates that the input parameter key value is invalid
	ErrExecutorErrorInvalidGetKey = errors.New("key is invalid")

	// Start of V2 errors

	// ErrExecutorSMMainCountersOverflowSha256 indicates that the sha256 counter exceeded the maximum
	ErrExecutorSMMainCountersOverflowSha256 = errors.New("sha256 counter exceeded the maximum")
	// ErrExecutorSMMainHashS indicates that a register value is out of range while calculating a Sha256 hash
	ErrExecutorSMMainHashS = errors.New("register value is out of range while calculating a Sha256 hash")
	// ErrExecutorSMMainHashSSizeOutOfRange indicates that a size register value is out of range while calculating a Sha256 hash
	ErrExecutorSMMainHashSSizeOutOfRange = errors.New("size register value is out of range while calculating a Sha256 hash")
	// ErrExecutorSMMainHashSPositionNegative indicates that a position register value is negative while calculating a Sha256 hash
	ErrExecutorSMMainHashSPositionNegative = errors.New("position register value is negative while calculating a Sha256 hash")
	// ErrExecutorSMMainHashSPositionPlusSizeOutOfRange indicates that a position register value plus a size register value is out of range while calculating a Sha256 hash
	ErrExecutorSMMainHashSPositionPlusSizeOutOfRange = errors.New("position register value plus a size register value is out of range while calculating a Sha256 hash")
	// ErrExecutorSMMainHashSDigestAddressNotFound indicates that an address has not been found while calculating a Sha256 hash digest
	ErrExecutorSMMainHashSDigestAddressNotFound = errors.New("address has not been found while calculating a Sha256 hash digest")
	// ErrExecutorSMMainHashSDigestNotCompleted indicates that the hash has not been completed while calling a Sha256 hash digest
	ErrExecutorSMMainHashSDigestNotCompleted = errors.New("hash has not been completed while calling a Sha256 hash digest")
	// ErrExecutorSMMainHashSValueMismatch indicates that the Sha256 hash instruction value check failed
	ErrExecutorSMMainHashSValueMismatch = errors.New("sha256 hash instruction value check failed")
	// ErrExecutorSMMainHashSPaddingMismatch indicates that the Sha256 hash instruction padding check failed
	ErrExecutorSMMainHashSPaddingMismatch = errors.New("sha256 hash instruction padding check failed")
	// ErrExecutorSMMainHashSSizeMismatch indicates that the Sha256 hash instruction size check failed
	ErrExecutorSMMainHashSSizeMismatch = errors.New("sha256 hash instruction size check failed")
	// ErrExecutorSMMainHashSLenLengthMismatch indicates that the Sha256 hash length instruction length check failed
	ErrExecutorSMMainHashSLenLengthMismatch = errors.New("sha256 hash length instruction length check failed")
	// ErrExecutorSMMainHashSLenCalledTwice indicates that the Sha256 hash length instruction called once check failed
	ErrExecutorSMMainHashSLenCalledTwice = errors.New("sha256 hash length instruction called once check failed")
	// ErrExecutorSMMainHashSDigestNotFound indicates that the Sha256 hash digest instruction slot not found
	ErrExecutorSMMainHashSDigestNotFound = errors.New("sha256 hash digest instruction slot not found")
	// ErrExecutorSMMainHashSDigestDigestMismatch indicates that the Sha256 hash digest instruction digest check failed
	ErrExecutorSMMainHashSDigestDigestMismatch = errors.New("sha256 hash digest instruction digest check failed")
	// ErrExecutorSMMainHashSDigestCalledTwice indicates that the Sha256 hash digest instruction called once check failed
	ErrExecutorSMMainHashSDigestCalledTwice = errors.New("sha256 hash digest instruction called once check failed")
	// ErrExecutorSMMainHashSReadOutOfRange indicates that the main execution Sha256 check found read out of range
	ErrExecutorSMMainHashSReadOutOfRange = errors.New("main execution Sha256 check found read out of range")
	// ErrExecutorErrorInvalidL1InfoRoot indicates that the input parameter l1_info_root is invalid
	ErrExecutorErrorInvalidL1InfoRoot = errors.New("l1_info_root is invalid")
	// ErrExecutorErrorInvalidForcedBlockhashL1 indicates that the input parameter forced_blockhash_l1 is invalid
	ErrExecutorErrorInvalidForcedBlockhashL1 = errors.New("forced_blockhash_l1 is invalid")
	// ErrExecutorErrorInvalidL1DataV2GlobalExitRoot indicates that the input parameter l1_data_v2.global_exit_root is invalid
	ErrExecutorErrorInvalidL1DataV2GlobalExitRoot = errors.New("l1_data_v2.global_exit_root is invalid")
	// ErrExecutorErrorInvalidL1DataV2BlockHashL1 indicates that the input parameter l1_data_v2.block_hash_l1 is invalid
	ErrExecutorErrorInvalidL1DataV2BlockHashL1 = errors.New("l1_data_v2.block_hash_l1 is invalid")
	// ErrExecutorErrorInvalidL1SmtProof indicates that the input parameter l1_smt_proof is invalid
	ErrExecutorErrorInvalidL1SmtProof = errors.New("l1_smt_proof is invalid")
	// ErrExecutorErrorInvalidBalance indicates that the input parameter balance is invalid
	ErrExecutorErrorInvalidBalance = errors.New("balance is invalid")
	// ErrExecutorErrorSMMainBinaryLt4Mismatch indicates that the binary instruction less than four opcode failed
	ErrExecutorErrorSMMainBinaryLt4Mismatch = errors.New("the binary instruction less than four opcode failed")
	// ErrExecutorErrorInvalidNewStateRoot indicates that the input parameter new_state_root is invalid
	ErrExecutorErrorInvalidNewStateRoot = errors.New("new_state_root is invalid")
	// ErrExecutorErrorInvalidNewAccInputHash indicates that the input parameter new_acc_input_hash is invalid
	ErrExecutorErrorInvalidNewAccInputHash = errors.New("new_acc_input_hash is invalid")
	// ErrExecutorErrorInvalidNewLocalExitRoot indicates that the input parameter new_local_exit_root is invalid
	ErrExecutorErrorInvalidNewLocalExitRoot = errors.New("new_local_exit_root is invalid")
	// ErrExecutorErrorDBKeyNotFound indicates that the requested key was not found in the database
	ErrExecutorErrorDBKeyNotFound = errors.New("key not found in the database")
	// ErrExecutorErrorSMTInvalidDataSize indicates that the SMT data returned from the database does not have a valid size
	ErrExecutorErrorSMTInvalidDataSize = errors.New("invalid SMT data size")
	// ErrExecutorErrorHashDBGRPCError indicates that the executor failed calling the HashDB service via GRPC, when configured
	ErrExecutorErrorHashDBGRPCError = errors.New("HashDB GRPC error")
	// ErrExecutorErrorStateManager indicates an error in the State Manager
	ErrExecutorErrorStateManager = errors.New("state Manager error")
	// ErrExecutorErrorInvalidL1InfoTreeIndex indicates that the ROM asked for an L1InfoTree index that was not present in the input
	ErrExecutorErrorInvalidL1InfoTreeIndex = errors.New("invalid l1_info_tree_index")
	// ErrExecutorErrorInvalidL1InfoTreeSmtProofValue indicates that the ROM asked for an L1InfoTree SMT proof that was not present in the input
	ErrExecutorErrorInvalidL1InfoTreeSmtProofValue = errors.New("invalid l1_info_tree_smt_proof_value")
	// ErrExecutorErrorInvalidWitness indicates that the input parameter witness is invalid
	ErrExecutorErrorInvalidWitness = errors.New("invalid witness")
	// ErrExecutorErrorInvalidCBOR indicates that the input parameter cbor is invalid
	ErrExecutorErrorInvalidCBOR = errors.New("invalid cbor")
	// ErrExecutorErrorInvalidDataStream indicates that the input parameter data stream is invalid
	ErrExecutorErrorInvalidDataStream = errors.New("invalid data stream")
	// ErrExecutorErrorInvalidUpdateMerkleTree indicates that the input parameter update merkle tree is invalid
	ErrExecutorErrorInvalidUpdateMerkleTree = errors.New("invalid update merkle tree")
	// ErrExecutorErrorSMMainInvalidTxStatusError indicates that the TX has an invalid status-error combination
	ErrExecutorErrorSMMainInvalidTxStatusError = errors.New("tx has an invalid status-error combination")

	// GRPC ERRORS
	// ===========

	// ErrGRPCResourceExhaustedAsTimeout indicates a GRPC resource exhausted error
	ErrGRPCResourceExhaustedAsTimeout = errors.New("request timed out")
)

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	ReturnValue   []byte // Returned data from the runtime (function result or data supplied with revert opcode)
	GasLeft       uint64 // Total gas left as result of execution
	GasUsed       uint64 // Total gas used as result of execution
	Err           error  // Any error encountered during the execution, listed below
	CreateAddress common.Address
	StateRoot     []byte
	FullTrace     instrumentation.FullTrace
	TraceResult   json.RawMessage
}

// Succeeded indicates the execution was successful
func (r *ExecutionResult) Succeeded() bool {
	return r.Err == nil
}

// Failed indicates the execution was unsuccessful
func (r *ExecutionResult) Failed() bool {
	return r.Err != nil
}

// Reverted indicates the execution was reverted
func (r *ExecutionResult) Reverted() bool {
	return errors.Is(r.Err, ErrExecutionReverted)
}
