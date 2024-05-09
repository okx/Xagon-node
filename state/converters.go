package state

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// MaxTxGasLimit is the gas limit allowed per tx in a batch
	MaxTxGasLimit = uint64(30000000)
)

// TestConvertToProcessBatchResponse for test purposes
func (s *State) TestConvertToProcessBatchResponse(batchResponse *executor.ProcessBatchResponse) (*ProcessBatchResponse, error) {
	return s.convertToProcessBatchResponse(batchResponse)
}

func (s *State) convertToProcessBatchResponse(batchResponse *executor.ProcessBatchResponse) (*ProcessBatchResponse, error) {
	blockResponses, err := s.convertToProcessBlockResponse(batchResponse.Responses)
	if err != nil {
		return nil, err
	}

	readWriteAddresses, err := convertToReadWriteAddresses(batchResponse.ReadWriteAddresses)
	if err != nil {
		return nil, err
	}

	isExecutorLevelError := (batchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR)
	isRomLevelError := false
	isRomOOCError := false

	if batchResponse.Responses != nil {
		for _, resp := range batchResponse.Responses {
			if resp.Error != executor.RomError_ROM_ERROR_NO_ERROR {
				isRomLevelError = true
				break
			}
		}

		if len(batchResponse.Responses) > 0 {
			// Check out of counters
			errorToCheck := batchResponse.Responses[len(batchResponse.Responses)-1].Error
			isRomOOCError = executor.IsROMOutOfCountersError(errorToCheck)
		}
	}

	return &ProcessBatchResponse{
		NewStateRoot:         common.BytesToHash(batchResponse.NewStateRoot),
		NewAccInputHash:      common.BytesToHash(batchResponse.NewAccInputHash),
		NewLocalExitRoot:     common.BytesToHash(batchResponse.NewLocalExitRoot),
		NewBatchNumber:       batchResponse.NewBatchNum,
		UsedZkCounters:       convertToCounters(batchResponse),
		BlockResponses:       blockResponses,
		ExecutorError:        executor.ExecutorErr(batchResponse.Error),
		ReadWriteAddresses:   readWriteAddresses,
		FlushID:              batchResponse.FlushId,
		StoredFlushID:        batchResponse.StoredFlushId,
		ProverID:             batchResponse.ProverId,
		IsExecutorLevelError: isExecutorLevelError,
		IsRomLevelError:      isRomLevelError,
		IsRomOOCError:        isRomOOCError,
		ForkID:               batchResponse.ForkId,
	}, nil
}

// IsStateRootChanged returns true if the transaction changes the state root
func IsStateRootChanged(err executor.RomError) bool {
	return !executor.IsIntrinsicError(err) && !executor.IsROMOutOfCountersError(err) && err != executor.RomError_ROM_ERROR_INVALID_RLP
}

func convertToReadWriteAddresses(addresses map[string]*executor.InfoReadWrite) (map[common.Address]*InfoReadWrite, error) {
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

func (s *State) convertToProcessBlockResponse(responses []*executor.ProcessTransactionResponse) ([]*ProcessBlockResponse, error) {
	results := make([]*ProcessBlockResponse, 0, len(responses))
	for _, response := range responses {
		blockResponse := new(ProcessBlockResponse)
		blockResponse.TransactionResponses = make([]*ProcessTransactionResponse, 0, 1)
		txResponse := new(ProcessTransactionResponse)
		txResponse.TxHash = common.BytesToHash(response.TxHash)
		txResponse.Type = response.Type
		txResponse.ReturnValue = response.ReturnValue
		txResponse.GasLeft = response.GasLeft
		txResponse.GasUsed = response.GasUsed
		txResponse.GasRefunded = response.GasRefunded
		txResponse.RomError = executor.RomErr(response.Error)
		txResponse.CreateAddress = common.HexToAddress(response.CreateAddress)
		txResponse.StateRoot = common.BytesToHash(response.StateRoot)
		txResponse.Logs = convertToLog(response.Logs)
		txResponse.ChangesStateRoot = IsStateRootChanged(response.Error)
		fullTrace, err := convertToFullTrace(response.FullTrace)
		if err != nil {
			return nil, err
		}
		txResponse.FullTrace = *fullTrace
		txResponse.EffectiveGasPrice = response.EffectiveGasPrice
		txResponse.EffectivePercentage = response.EffectivePercentage
		txResponse.HasGaspriceOpcode = (response.HasGaspriceOpcode == 1)
		txResponse.HasBalanceOpcode = (response.HasBalanceOpcode == 1)

		tx := new(types.Transaction)

		if response.Error != executor.RomError_ROM_ERROR_INVALID_RLP && len(response.GetRlpTx()) > 0 {
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

				return nil, err
			}
		} else {
			log.Warnf("ROM_ERROR_INVALID_RLP returned by the executor")
		}

		if tx != nil {
			txResponse.Tx = *tx
			log.Debugf("ProcessTransactionResponse[TxHash]: %v", txResponse.TxHash)
			if response.Error == executor.RomError_ROM_ERROR_NO_ERROR {
				log.Debugf("ProcessTransactionResponse[Nonce]: %v", txResponse.Tx.Nonce())
			}
			log.Debugf("ProcessTransactionResponse[StateRoot]: %v", txResponse.StateRoot.String())
			log.Debugf("ProcessTransactionResponse[Error]: %v", txResponse.RomError)
			log.Debugf("ProcessTransactionResponse[GasUsed]: %v", txResponse.GasUsed)
			log.Debugf("ProcessTransactionResponse[GasLeft]: %v", txResponse.GasLeft)
			log.Debugf("ProcessTransactionResponse[GasRefunded]: %v", txResponse.GasRefunded)
			log.Debugf("ProcessTransactionResponse[ChangesStateRoot]: %v", txResponse.ChangesStateRoot)
			log.Debugf("ProcessTransactionResponse[EffectiveGasPrice]: %v", txResponse.EffectiveGasPrice)
			log.Debugf("ProcessTransactionResponse[EffectivePercentage]: %v", txResponse.EffectivePercentage)
		}

		blockResponse.TransactionResponses = append(blockResponse.TransactionResponses, txResponse)
		blockResponse.GasLimit = MaxTxGasLimit
		results = append(results, blockResponse)
	}

	return results, nil
}

func convertToLog(protoLogs []*executor.Log) []*types.Log {
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

func convertToTopics(responses [][]byte) []common.Hash {
	results := make([]common.Hash, 0, len(responses))

	for _, response := range responses {
		results = append(results, common.BytesToHash(response))
	}
	return results
}

func convertToFullTrace(fullTrace *executor.FullTrace) (*instrumentation.FullTrace, error) {
	trace := new(instrumentation.FullTrace)
	if fullTrace != nil {
		trace.Context = convertToContext(fullTrace.Context)
		steps, err := convertToInstrumentationSteps(fullTrace.Steps)
		if err != nil {
			return nil, err
		}
		trace.Steps = steps
	}

	return trace, nil
}

func convertToContext(context *executor.TransactionContext) instrumentation.Context {
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

func convertToInstrumentationSteps(responses []*executor.TransactionStep) ([]instrumentation.Step, error) {
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
		step.Contract = convertToInstrumentationContract(response.Contract)
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

func convertToInstrumentationContract(response *executor.Contract) instrumentation.Contract {
	return instrumentation.Contract{
		Address: common.HexToAddress(response.Address),
		Caller:  common.HexToAddress(response.Caller),
		Value:   hex.DecodeBig(response.Value),
		Input:   response.Data,
		Gas:     response.Gas,
	}
}

func convertToCounters(resp *executor.ProcessBatchResponse) ZKCounters {
	return ZKCounters{
		GasUsed:          resp.CumulativeGasUsed,
		KeccakHashes:     resp.CntKeccakHashes,
		PoseidonHashes:   resp.CntPoseidonHashes,
		PoseidonPaddings: resp.CntPoseidonPaddings,
		MemAligns:        resp.CntMemAligns,
		Arithmetics:      resp.CntArithmetics,
		Binaries:         resp.CntBinaries,
		Steps:            resp.CntSteps,
	}
}
