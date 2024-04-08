package state

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// MaxL2BlockGasLimit is the gas limit allowed per L2 block in a batch
	MaxL2BlockGasLimit = uint64(1125899906842624)
)

var (
	errL2BlockInvalid = errors.New("A L2 block fails, that invalidate totally the batch")
)

// TestConvertToProcessBatchResponseV2 for test purposes
func (s *State) TestConvertToProcessBatchResponseV2(batchResponse *executor.ProcessBatchResponseV2) (*ProcessBatchResponse, error) {
	return s.convertToProcessBatchResponseV2(batchResponse)
}

func (s *State) convertToProcessBatchResponseV2(batchResponse *executor.ProcessBatchResponseV2) (*ProcessBatchResponse, error) {
	blockResponses, isRomLevelError, isRomOOCError, err := s.convertToProcessBlockResponseV2(batchResponse.BlockResponses)
	if err != nil {
		return nil, err
	}
	isRomOOCError = isRomOOCError || executor.IsROMOutOfCountersError(batchResponse.ErrorRom)
	readWriteAddresses, err := convertToReadWriteAddressesV2(batchResponse.ReadWriteAddresses)
	if err != nil {
		return nil, err
	}

	return &ProcessBatchResponse{
		NewStateRoot:         common.BytesToHash(batchResponse.NewStateRoot),
		NewAccInputHash:      common.BytesToHash(batchResponse.NewAccInputHash),
		NewLocalExitRoot:     common.BytesToHash(batchResponse.NewLocalExitRoot),
		NewBatchNumber:       batchResponse.NewBatchNum,
		UsedZkCounters:       convertToUsedZKCountersV2(batchResponse),
		ReservedZkCounters:   convertToReservedZKCountersV2(batchResponse),
		BlockResponses:       blockResponses,
		ExecutorError:        executor.ExecutorErr(batchResponse.Error),
		ReadWriteAddresses:   readWriteAddresses,
		FlushID:              batchResponse.FlushId,
		StoredFlushID:        batchResponse.StoredFlushId,
		ProverID:             batchResponse.ProverId,
		IsExecutorLevelError: batchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR,
		IsRomLevelError:      isRomLevelError,
		IsRomOOCError:        isRomOOCError,
		GasUsed_V2:           batchResponse.GasUsed,
		SMTKeys_V2:           convertToKeys(batchResponse.SmtKeys),
		ProgramKeys_V2:       convertToKeys(batchResponse.ProgramKeys),
		ForkID:               batchResponse.ForkId,
		InvalidBatch_V2:      batchResponse.InvalidBatch != 0,
		RomError_V2:          executor.RomErr(batchResponse.ErrorRom),
	}, nil
}

func (s *State) convertToProcessBlockResponseV2(responses []*executor.ProcessBlockResponseV2) ([]*ProcessBlockResponse, bool, bool, error) {
	isRomLevelError := false
	isRomOOCError := false

	results := make([]*ProcessBlockResponse, 0, len(responses))
	for _, response := range responses {
		result := new(ProcessBlockResponse)
		transactionResponses, respisRomLevelError, respisRomOOCError, err := s.convertToProcessTransactionResponseV2(response.Responses)
		isRomLevelError = isRomLevelError || respisRomLevelError
		isRomOOCError = isRomOOCError || respisRomOOCError
		if err != nil {
			return nil, isRomLevelError, isRomOOCError, err
		}

		result.ParentHash = common.BytesToHash(response.ParentHash)
		result.Coinbase = common.HexToAddress(response.Coinbase)
		result.GasLimit = response.GasLimit
		result.BlockNumber = response.BlockNumber
		result.Timestamp = response.Timestamp
		result.GlobalExitRoot = common.Hash(response.Ger)
		result.BlockHashL1 = common.Hash(response.BlockHashL1)
		result.GasUsed = response.GasUsed
		result.BlockInfoRoot = common.Hash(response.BlockInfoRoot)
		result.BlockHash = common.Hash(response.BlockHash)
		result.TransactionResponses = transactionResponses
		result.Logs = convertToLogV2(response.Logs)
		result.RomError_V2 = executor.RomErr(response.Error)

		results = append(results, result)
	}

	return results, isRomLevelError, isRomOOCError, nil
}

func (s *State) convertToProcessTransactionResponseV2(responses []*executor.ProcessTransactionResponseV2) ([]*ProcessTransactionResponse, bool, bool, error) {
	isRomLevelError := false
	isRomOOCError := false

	results := make([]*ProcessTransactionResponse, 0, len(responses))

	for _, response := range responses {
		if response.Error != executor.RomError_ROM_ERROR_NO_ERROR {
			isRomLevelError = true
		}
		if executor.IsROMOutOfCountersError(response.Error) {
			isRomOOCError = true
		}
		if executor.IsInvalidL2Block(response.Error) {
			err := fmt.Errorf("fails L2 block: romError %v error:%w", response.Error, errL2BlockInvalid)
			return nil, isRomLevelError, isRomOOCError, err
		}
		result := new(ProcessTransactionResponse)
		result.TxHash = common.BytesToHash(response.TxHash)
		result.TxHashL2_V2 = common.BytesToHash(response.TxHashL2)
		result.Type = response.Type
		result.ReturnValue = response.ReturnValue
		result.GasLeft = response.GasLeft
		result.GasUsed = response.GasUsed
		result.CumulativeGasUsed = response.CumulativeGasUsed
		result.GasRefunded = response.GasRefunded
		result.RomError = executor.RomErr(response.Error)
		result.CreateAddress = common.HexToAddress(response.CreateAddress)
		result.StateRoot = common.BytesToHash(response.StateRoot)
		result.Logs = convertToLogV2(response.Logs)
		result.ChangesStateRoot = IsStateRootChanged(response.Error)
		fullTrace, err := convertToFullTraceV2(response.FullTrace)
		if err != nil {
			return nil, isRomLevelError, isRomOOCError, err
		}
		result.FullTrace = *fullTrace
		result.EffectiveGasPrice = response.EffectiveGasPrice
		result.EffectivePercentage = response.EffectivePercentage
		result.HasGaspriceOpcode = (response.HasGaspriceOpcode == 1)
		result.HasBalanceOpcode = (response.HasBalanceOpcode == 1)
		result.Status = response.Status

		var tx *types.Transaction
		if response.Error != executor.RomError_ROM_ERROR_INVALID_RLP {
			if len(response.GetRlpTx()) > 0 {
				tx, err = DecodeTx(common.Bytes2Hex(response.GetRlpTx()))
				if err != nil {
					timestamp := time.Now()
					log.Errorf("error decoding rlp returned by executor %v at %v", err, timestamp)

					event := &event.Event{
						ReceivedAt: timestamp,
						Source:     event.Source_Node,
						Level:      event.Level_Error,
						EventID:    event.EventID_ExecutorRLPError,
						Json:       string(response.GetRlpTx()),
					}

					eventErr := s.eventLog.LogEvent(context.Background(), event)
					if eventErr != nil {
						log.Errorf("error storing payload: %v", err)
					}

					return nil, isRomLevelError, isRomOOCError, err
				}
			} else {
				log.Infof("no txs returned by executor")
			}
		} else {
			log.Warnf("ROM_ERROR_INVALID_RLP returned by the executor")
		}

		if tx != nil {
			result.Tx = *tx
		}

		results = append(results, result)
	}

	return results, isRomLevelError, isRomOOCError, nil
}

func convertToLogV2(protoLogs []*executor.LogV2) []*types.Log {
	logs := make([]*types.Log, 0, len(protoLogs))

	for _, protoLog := range protoLogs {
		log := new(types.Log)
		log.Address = common.HexToAddress(protoLog.Address)
		log.Topics = convertToTopics(protoLog.Topics)
		log.Data = protoLog.Data
		log.TxHash = common.BytesToHash(protoLog.TxHash)
		log.TxIndex = uint(protoLog.TxIndex)
		log.Index = uint(protoLog.Index)
		logs = append(logs, log)
	}

	return logs
}

func convertToFullTraceV2(fullTrace *executor.FullTraceV2) (*instrumentation.FullTrace, error) {
	trace := new(instrumentation.FullTrace)
	if fullTrace != nil {
		trace.Context = convertToContextV2(fullTrace.Context)
		steps, err := convertToInstrumentationStepsV2(fullTrace.Steps)
		if err != nil {
			return nil, err
		}
		trace.Steps = steps
	}

	return trace, nil
}

func convertToContextV2(context *executor.TransactionContextV2) instrumentation.Context {
	return instrumentation.Context{
		Type:         context.Type,
		From:         context.From,
		To:           context.To,
		Input:        context.Data,
		Gas:          context.Gas,
		Value:        hex.DecodeBig(context.Value),
		Output:       context.Output,
		GasPrice:     context.GasPrice,
		OldStateRoot: common.BytesToHash(context.OldStateRoot),
		Time:         uint64(context.ExecutionTime),
		GasUsed:      context.GasUsed,
	}
}

func convertToInstrumentationStepsV2(responses []*executor.TransactionStepV2) ([]instrumentation.Step, error) {
	results := make([]instrumentation.Step, 0, len(responses))
	for _, response := range responses {
		step := new(instrumentation.Step)
		step.StateRoot = common.BytesToHash(response.StateRoot)
		step.Depth = int(response.Depth)
		step.Pc = response.Pc
		step.Gas = response.Gas
		step.OpCode = fakevm.OpCode(response.Op).String()
		step.Refund = response.GasRefund
		step.Op = uint64(response.Op)
		err := executor.RomErr(response.Error)
		if err != nil {
			step.Error = err
		}
		step.Contract = convertToInstrumentationContractV2(response.Contract)
		step.GasCost = response.GasCost
		step.Stack = make([]*big.Int, 0, len(response.Stack))
		for _, s := range response.Stack {
			if len(s)%2 != 0 {
				s = "0" + s
			}
			bi, ok := new(big.Int).SetString(s, hex.Base)
			if !ok {
				log.Debugf("error while parsing stack valueBigInt")
				return nil, ErrParsingExecutorTrace
			}
			step.Stack = append(step.Stack, bi)
		}
		step.MemorySize = response.MemorySize
		step.MemoryOffset = response.MemoryOffset
		step.Memory = make([]byte, len(response.Memory))
		copy(step.Memory, response.Memory)
		step.ReturnData = make([]byte, len(response.ReturnData))
		copy(step.ReturnData, response.ReturnData)
		step.Storage = make(map[common.Hash]common.Hash, len(response.Storage))
		for k, v := range response.Storage {
			addr := common.BytesToHash(hex.DecodeBig(k).Bytes())
			value := common.BytesToHash(hex.DecodeBig(v).Bytes())
			step.Storage[addr] = value
		}
		results = append(results, *step)
	}
	return results, nil
}

func convertToInstrumentationContractV2(response *executor.ContractV2) instrumentation.Contract {
	return instrumentation.Contract{
		Address: common.HexToAddress(response.Address),
		Caller:  common.HexToAddress(response.Caller),
		Value:   hex.DecodeBig(response.Value),
		Input:   response.Data,
		Gas:     response.Gas,
	}
}

func convertToUsedZKCountersV2(resp *executor.ProcessBatchResponseV2) ZKCounters {
	return ZKCounters{
		GasUsed:          resp.GasUsed,
		KeccakHashes:     resp.CntKeccakHashes,
		PoseidonHashes:   resp.CntPoseidonHashes,
		PoseidonPaddings: resp.CntPoseidonPaddings,
		MemAligns:        resp.CntMemAligns,
		Arithmetics:      resp.CntArithmetics,
		Binaries:         resp.CntBinaries,
		Steps:            resp.CntSteps,
		Sha256Hashes_V2:  resp.CntSha256Hashes,
	}
}

func convertToReservedZKCountersV2(resp *executor.ProcessBatchResponseV2) ZKCounters {
	return ZKCounters{
		// There is no "ReserveGasUsed" in the response, so we use "GasUsed" as it will make calculations easier
		GasUsed:          resp.GasUsed,
		KeccakHashes:     resp.CntReserveKeccakHashes,
		PoseidonHashes:   resp.CntReservePoseidonHashes,
		PoseidonPaddings: resp.CntReservePoseidonPaddings,
		MemAligns:        resp.CntReserveMemAligns,
		Arithmetics:      resp.CntReserveArithmetics,
		Binaries:         resp.CntReserveBinaries,
		Steps:            resp.CntReserveSteps,
		Sha256Hashes_V2:  resp.CntReserveSha256Hashes,
	}
}

func convertToReadWriteAddressesV2(addresses map[string]*executor.InfoReadWriteV2) (map[common.Address]*InfoReadWrite, error) {
	results := make(map[common.Address]*InfoReadWrite, len(addresses))

	for addr, addrInfo := range addresses {
		var nonce *uint64 = nil
		var balance *big.Int = nil
		var ok bool

		address := common.HexToAddress(addr)

		if addrInfo.Nonce != "" {
			bigNonce, ok := new(big.Int).SetString(addrInfo.Nonce, encoding.Base10)
			if !ok {
				log.Debugf("received nonce as string: %v", addrInfo.Nonce)
				return nil, fmt.Errorf("error while parsing address nonce")
			}
			nonceNp := bigNonce.Uint64()
			nonce = &nonceNp
		}

		if addrInfo.Balance != "" {
			balance, ok = new(big.Int).SetString(addrInfo.Balance, encoding.Base10)
			if !ok {
				log.Debugf("received balance as string: %v", addrInfo.Balance)
				return nil, fmt.Errorf("error while parsing address balance")
			}
		}

		results[address] = &InfoReadWrite{Address: address, Nonce: nonce, Balance: balance}
	}

	return results, nil
}

func convertToKeys(keys [][]byte) []merkletree.Key {
	result := make([]merkletree.Key, 0, len(keys))
	for _, key := range keys {
		result = append(result, merkletree.Key(key))
	}
	return result
}

func convertProcessingContext(p *ProcessingContextV2) (*ProcessingContext, error) {
	tstamp := time.Time{}
	if p.Timestamp != nil {
		tstamp = *p.Timestamp
	}
	result := ProcessingContext{
		BatchNumber:    p.BatchNumber,
		Coinbase:       p.Coinbase,
		ForcedBatchNum: p.ForcedBatchNum,
		BatchL2Data:    p.BatchL2Data,
		Timestamp:      tstamp,
		GlobalExitRoot: p.GlobalExitRoot,
		ClosingReason:  p.ClosingReason,
	}
	return &result, nil
}
