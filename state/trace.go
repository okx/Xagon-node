package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/js"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers/native"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers/structlogger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/google/uuid"
	"github.com/holiman/uint256"
	"github.com/jackc/pgx/v4"
)

// DebugTransaction re-executes a tx to generate its trace
func (s *State) DebugTransaction(ctx context.Context, transactionHash common.Hash, traceConfig TraceConfig, dbTx pgx.Tx) (*runtime.ExecutionResult, error) {
	var err error

	// gets the transaction
	tx, err := s.GetTransactionByHash(ctx, transactionHash, dbTx)
	if err != nil {
		return nil, err
	}

	// gets the tx receipt
	receipt, err := s.GetTransactionReceipt(ctx, transactionHash, dbTx)
	if err != nil {
		return nil, err
	}

	// gets the l2 l2Block including the transaction
	l2Block, err := s.GetL2BlockByNumber(ctx, receipt.BlockNumber.Uint64(), dbTx)
	if err != nil {
		return nil, err
	}

	// the old state root is the previous block state root
	var oldStateRoot common.Hash
	previousL2BlockNumber := uint64(0)
	if receipt.BlockNumber.Uint64() > 0 {
		previousL2BlockNumber = receipt.BlockNumber.Uint64() - 1
	}
	previousL2Block, err := s.GetL2BlockByNumber(ctx, previousL2BlockNumber, dbTx)
	if err != nil {
		return nil, err
	}
	oldStateRoot = previousL2Block.Root()

	count := 0
	for _, tx := range l2Block.Transactions() {
		checkReceipt, err := s.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
		if err != nil {
			return nil, err
		}
		if checkReceipt.TransactionIndex < receipt.TransactionIndex {
			count++
		}
	}

	// since the executor only stores the state roots by block, we need to
	// execute all the txs in the block until the tx we want to trace
	var txsToEncode []types.Transaction
	var effectivePercentage []uint8
	for i := 0; i <= count; i++ {
		txsToEncode = append(txsToEncode, *l2Block.Transactions()[i])
		effectivePercentage = append(effectivePercentage, MaxEffectivePercentage)
		log.Debugf("trace will reprocess tx: %v", l2Block.Transactions()[i].Hash().String())
	}

	// gets batch that including the l2 block
	batch, err := s.GetBatchByL2BlockNumber(ctx, l2Block.NumberU64(), dbTx)
	if err != nil {
		return nil, err
	}

	forkId := s.GetForkIDByBatchNumber(batch.BatchNumber)

	var response *ProcessTransactionResponse
	var startTime, endTime time.Time
	if forkId < FORKID_ETROG {
		traceConfigRequest := &executor.TraceConfig{
			TxHashToGenerateFullTrace: transactionHash.Bytes(),
			// set the defaults to the maximum information we can have.
			// this is needed to process custom tracers later
			DisableStorage:   cFalse,
			DisableStack:     cFalse,
			EnableMemory:     cTrue,
			EnableReturnData: cTrue,
		}

		// if the default tracer is used, then we review the information
		// we want to have in the trace related to the parameters we received.
		if traceConfig.IsDefaultTracer() {
			if traceConfig.DisableStorage {
				traceConfigRequest.DisableStorage = cTrue
			}
			if traceConfig.DisableStack {
				traceConfigRequest.DisableStack = cTrue
			}
			if !traceConfig.EnableMemory {
				traceConfigRequest.EnableMemory = cFalse
			}
			if !traceConfig.EnableReturnData {
				traceConfigRequest.EnableReturnData = cFalse
			}
		}
		// generate batch l2 data for the transaction
		batchL2Data, err := EncodeTransactions(txsToEncode, effectivePercentage, forkId)
		if err != nil {
			return nil, err
		}

		// prepare process batch request
		processBatchRequest := &executor.ProcessBatchRequest{
			OldBatchNum:     batch.BatchNumber - 1,
			OldStateRoot:    oldStateRoot.Bytes(),
			OldAccInputHash: batch.AccInputHash.Bytes(),

			BatchL2Data:      batchL2Data,
			Coinbase:         batch.Coinbase.String(),
			UpdateMerkleTree: cFalse,
			ChainId:          s.cfg.ChainID,
			ForkId:           forkId,
			TraceConfig:      traceConfigRequest,
			ContextId:        uuid.NewString(),

			GlobalExitRoot: batch.GlobalExitRoot.Bytes(),
			EthTimestamp:   uint64(batch.Timestamp.Unix()),
		}

		// Send Batch to the Executor
		startTime = time.Now()
		processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
		endTime = time.Now()
		if err != nil {
			return nil, err
		} else if processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
			err = executor.ExecutorErr(processBatchResponse.Error)
			s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
			return nil, err
		}

		// Transactions are decoded only for logging purposes
		// as they are not longer needed in the convertToProcessBatchResponse function
		txs, _, _, err := DecodeTxs(batchL2Data, forkId)
		if err != nil && !errors.Is(err, ErrInvalidData) {
			return nil, err
		}

		for _, tx := range txs {
			log.Debugf(tx.Hash().String())
		}

		convertedResponse, err := s.convertToProcessBatchResponse(processBatchResponse)
		if err != nil {
			return nil, err
		}
		response = convertedResponse.BlockResponses[0].TransactionResponses[0]
	} else {
		traceConfigRequestV2 := &executor.TraceConfigV2{
			TxHashToGenerateFullTrace: transactionHash.Bytes(),
			// set the defaults to the maximum information we can have.
			// this is needed to process custom tracers later
			DisableStorage:   cFalse,
			DisableStack:     cFalse,
			EnableMemory:     cTrue,
			EnableReturnData: cTrue,
		}

		// if the default tracer is used, then we review the information
		// we want to have in the trace related to the parameters we received.
		if traceConfig.IsDefaultTracer() {
			if traceConfig.DisableStorage {
				traceConfigRequestV2.DisableStorage = cTrue
			}
			if traceConfig.DisableStack {
				traceConfigRequestV2.DisableStack = cTrue
			}
			if !traceConfig.EnableMemory {
				traceConfigRequestV2.EnableMemory = cFalse
			}
			if !traceConfig.EnableReturnData {
				traceConfigRequestV2.EnableReturnData = cFalse
			}
		}

		// if the l2 block number is 1, it means this is a network that started
		// at least on Etrog fork, in this case the l2 block 1 will contain the
		// injected tx that needs to be processed in a different way
		isInjectedTx := l2Block.NumberU64() == 1

		var transactions, batchL2Data []byte
		if isInjectedTx {
			transactions = append([]byte{}, batch.BatchL2Data...)
		} else {
			// build the raw batch so we can get the index l1 info tree for the l2 block
			rawBatch, err := DecodeBatchV2(batch.BatchL2Data)
			if err != nil {
				log.Errorf("error decoding BatchL2Data for batch %d, error: %v", batch.BatchNumber, err)
				return nil, err
			}

			// identify the first l1 block number so we can identify the
			// current l2 block index in the block array
			firstBlockNumberForBatch, err := s.GetFirstL2BlockNumberForBatchNumber(ctx, batch.BatchNumber, dbTx)
			if err != nil {
				log.Errorf("failed to get first l2 block number for batch %v: %v ", batch.BatchNumber, err)
				return nil, err
			}

			// computes the l2 block index
			rawL2BlockIndex := l2Block.NumberU64() - firstBlockNumberForBatch
			if rawL2BlockIndex > uint64(len(rawBatch.Blocks)-1) {
				log.Errorf("computed rawL2BlockIndex is greater than the number of blocks we have in the batch %v: %v ", batch.BatchNumber, err)
				return nil, err
			}

			// builds the ChangeL2Block transaction with the correct timestamp and IndexL1InfoTree
			rawL2Block := rawBatch.Blocks[rawL2BlockIndex]
			deltaTimestamp := uint32(l2Block.Time() - previousL2Block.Time())
			transactions = s.BuildChangeL2Block(deltaTimestamp, rawL2Block.IndexL1InfoTree)

			batchL2Data, err = EncodeTransactions(txsToEncode, effectivePercentage, forkId)
			if err != nil {
				log.Errorf("error encoding transaction ", err)
				return nil, err
			}

			transactions = append(transactions, batchL2Data...)
		}
		// prepare process batch request
		processBatchRequestV2 := &executor.ProcessBatchRequestV2{
			OldBatchNum:     batch.BatchNumber - 1,
			OldStateRoot:    oldStateRoot.Bytes(),
			OldAccInputHash: batch.AccInputHash.Bytes(),

			BatchL2Data:      transactions,
			Coinbase:         l2Block.Coinbase().String(),
			UpdateMerkleTree: cFalse,
			ChainId:          s.cfg.ChainID,
			ForkId:           forkId,
			TraceConfig:      traceConfigRequestV2,
			ContextId:        uuid.NewString(),

			// v2 fields
			L1InfoRoot:             GetMockL1InfoRoot().Bytes(),
			TimestampLimit:         uint64(time.Now().Unix()),
			SkipFirstChangeL2Block: cFalse,
			SkipWriteBlockInfoRoot: cTrue,
		}

		if isInjectedTx {
			virtualBatch, err := s.GetVirtualBatch(ctx, batch.BatchNumber, dbTx)
			if err != nil {
				log.Errorf("failed to load virtual batch %v", batch.BatchNumber, err)
				return nil, err
			}
			l1Block, err := s.GetBlockByNumber(ctx, virtualBatch.BlockNumber, dbTx)
			if err != nil {
				log.Errorf("failed to load l1 block %v", virtualBatch.BlockNumber, err)
				return nil, err
			}

			processBatchRequestV2.ForcedBlockhashL1 = l1Block.BlockHash.Bytes()
			processBatchRequestV2.SkipVerifyL1InfoRoot = 1
		} else {
			// gets the L1InfoTreeData for the transactions
			l1InfoTreeData, _, _, err := s.GetL1InfoTreeDataFromBatchL2Data(ctx, transactions, dbTx)
			if err != nil {
				return nil, err
			}

			// In case we have any l1InfoTreeData, add them to the request
			if len(l1InfoTreeData) > 0 {
				processBatchRequestV2.L1InfoTreeData = map[uint32]*executor.L1DataV2{}
				processBatchRequestV2.SkipVerifyL1InfoRoot = cTrue
				for k, v := range l1InfoTreeData {
					processBatchRequestV2.L1InfoTreeData[k] = &executor.L1DataV2{
						GlobalExitRoot: v.GlobalExitRoot.Bytes(),
						BlockHashL1:    v.BlockHashL1.Bytes(),
						MinTimestamp:   v.MinTimestamp,
					}
				}
			}
		}

		// Send Batch to the Executor
		startTime = time.Now()
		processBatchResponseV2, err := s.executorClient.ProcessBatchV2(ctx, processBatchRequestV2)
		endTime = time.Now()
		if err != nil {
			return nil, err
		} else if processBatchResponseV2.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
			err = executor.ExecutorErr(processBatchResponseV2.Error)
			s.eventLog.LogExecutorError(ctx, processBatchResponseV2.Error, processBatchRequestV2)
			return nil, err
		}

		if !isInjectedTx {
			// Transactions are decoded only for logging purposes
			// as they are no longer needed in the convertToProcessBatchResponse function
			txs, _, _, err := DecodeTxs(batchL2Data, forkId)
			if err != nil && !errors.Is(err, ErrInvalidData) {
				return nil, err
			}
			for _, tx := range txs {
				log.Debugf(tx.Hash().String())
			}
		}

		convertedResponse, err := s.convertToProcessBatchResponseV2(processBatchResponseV2)
		if err != nil {
			return nil, err
		}
		response = convertedResponse.BlockResponses[0].TransactionResponses[len(convertedResponse.BlockResponses[0].TransactionResponses)-1]
	}

	// Sanity check
	log.Debugf(response.TxHash.String())
	if response.TxHash != transactionHash {
		return nil, fmt.Errorf("tx hash not found in executor response")
	}

	result := &runtime.ExecutionResult{
		CreateAddress: response.CreateAddress,
		GasLeft:       response.GasLeft,
		GasUsed:       response.GasUsed,
		ReturnValue:   response.ReturnValue,
		StateRoot:     response.StateRoot.Bytes(),
		FullTrace:     response.FullTrace,
		Err:           response.RomError,
	}

	senderAddress, err := GetSender(*tx)
	if err != nil {
		return nil, err
	}

	context := instrumentation.Context{
		From:         senderAddress.String(),
		Input:        tx.Data(),
		Gas:          tx.Gas(),
		Value:        tx.Value(),
		Output:       result.ReturnValue,
		GasPrice:     tx.GasPrice().String(),
		OldStateRoot: oldStateRoot,
		Time:         uint64(endTime.Sub(startTime)),
		GasUsed:      result.GasUsed,
	}

	// Fill trace context
	if tx.To() == nil {
		context.Type = "CREATE"
		context.To = result.CreateAddress.Hex()
	} else {
		context.Type = "CALL"
		context.To = tx.To().Hex()
	}

	result.FullTrace.Context = context

	gasPrice, ok := new(big.Int).SetString(context.GasPrice, encoding.Base10)
	if !ok {
		log.Errorf("debug transaction: failed to parse gasPrice")
		return nil, fmt.Errorf("failed to parse gasPrice")
	}

	// select and prepare tracer
	var tracer tracers.Tracer
	tracerContext := &tracers.Context{
		BlockHash:   receipt.BlockHash,
		BlockNumber: receipt.BlockNumber,
		TxIndex:     int(receipt.TransactionIndex),
		TxHash:      transactionHash,
	}

	if traceConfig.IsDefaultTracer() {
		structLoggerCfg := structlogger.Config{
			EnableMemory:     traceConfig.EnableMemory,
			DisableStack:     traceConfig.DisableStack,
			DisableStorage:   traceConfig.DisableStorage,
			EnableReturnData: traceConfig.EnableReturnData,
		}
		tracer := structlogger.NewStructLogger(structLoggerCfg)
		traceResult, err := tracer.ParseTrace(result, *receipt)
		if err != nil {
			return nil, err
		}
		result.TraceResult = traceResult
		return result, nil
	} else if traceConfig.Is4ByteTracer() {
		tracer, err = native.NewFourByteTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create 4byteTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create 4byteTracer, err: %v", err)
		}
	} else if traceConfig.IsCallTracer() {
		tracer, err = native.NewCallTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create callTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create callTracer, err: %v", err)
		}
	} else if traceConfig.IsNoopTracer() {
		tracer, err = native.NewNoopTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create noopTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create noopTracer, err: %v", err)
		}
	} else if traceConfig.IsPrestateTracer() {
		tracer, err = native.NewPrestateTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create prestateTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create prestateTracer, err: %v", err)
		}
	} else if traceConfig.IsJSCustomTracer() {
		tracer, err = js.NewJsTracer(*traceConfig.Tracer, tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create jsTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create jsTracer, err: %v", err)
		}
	} else {
		return nil, fmt.Errorf("invalid tracer: %v, err: %v", traceConfig.Tracer, err)
	}

	fakeDB := &FakeDB{State: s, stateRoot: batch.StateRoot.Bytes()}
	evm := fakevm.NewFakeEVM(fakevm.BlockContext{BlockNumber: big.NewInt(1)}, fakevm.TxContext{GasPrice: gasPrice}, fakeDB, params.TestChainConfig, fakevm.Config{Debug: true, Tracer: tracer})

	traceResult, err := s.buildTrace(evm, result, tracer)
	if err != nil {
		log.Errorf("debug transaction: failed parse the trace using the tracer: %v", err)
		return nil, fmt.Errorf("failed parse the trace using the tracer: %v", err)
	}

	result.TraceResult = traceResult

	return result, nil
}

// ParseTheTraceUsingTheTracer parses the given trace with the given tracer.
func (s *State) buildTrace(evm *fakevm.FakeEVM, result *runtime.ExecutionResult, tracer tracers.Tracer) (json.RawMessage, error) {
	trace := result.FullTrace
	tracer.CaptureTxStart(trace.Context.Gas)
	contextGas := trace.Context.Gas - trace.Context.GasUsed
	if len(trace.Steps) > 0 {
		contextGas = trace.Steps[0].Gas
	}
	tracer.CaptureStart(evm, common.HexToAddress(trace.Context.From), common.HexToAddress(trace.Context.To), trace.Context.Type == "CREATE", trace.Context.Input, contextGas, trace.Context.Value)
	evm.StateDB.SetStateRoot(trace.Context.OldStateRoot.Bytes())

	var previousStep instrumentation.Step
	reverted := false
	internalTxSteps := NewStack[instrumentation.InternalTxContext]()
	memory := fakevm.NewMemory()

	for i, step := range trace.Steps {
		if step.OpCode == "SSTORE" {
			time.Sleep(time.Millisecond)
		}

		if step.OpCode == "SLOAD" {
			time.Sleep(time.Millisecond)
		}

		if step.OpCode == "RETURN" {
			time.Sleep(time.Millisecond)
		}

		// set Stack
		stack := fakevm.NewStack()
		for _, stackItem := range step.Stack {
			value, _ := uint256.FromBig(stackItem)
			stack.Push(value)
		}

		// set Memory
		memory.Resize(uint64(step.MemorySize))
		if len(step.Memory) > 0 {
			memory.Set(uint64(step.MemoryOffset), uint64(len(step.Memory)), step.Memory)
		}

		// Populate the step memory for future steps
		step.Memory = memory.Data()

		// set Contract
		contract := fakevm.NewContract(
			fakevm.NewAccount(step.Contract.Caller),
			fakevm.NewAccount(step.Contract.Address),
			step.Contract.Value, step.Gas)
		aux := step.Contract.Address
		contract.CodeAddr = &aux

		// set Scope
		scope := &fakevm.ScopeContext{
			Contract: contract,
			Memory:   memory,
			Stack:    stack,
		}

		// if the revert happens on an internal tx, we exit
		if previousStep.OpCode == "REVERT" && previousStep.Depth > 1 {
			gasUsed, err := s.getGasUsed(internalTxSteps, previousStep, step)
			if err != nil {
				return nil, err
			}
			tracer.CaptureExit(step.ReturnData, gasUsed, fakevm.ErrExecutionReverted)
		}

		// if the revert happens on top level, we break
		if step.OpCode == "REVERT" && step.Depth == 1 {
			reverted = true
			break
		}

		hasNextStep := i < len(trace.Steps)-1
		if step.OpCode != "CALL" || (hasNextStep && trace.Steps[i+1].Pc == 0) {
			if step.Error != nil {
				tracer.CaptureFault(step.Pc, fakevm.OpCode(step.Op), step.Gas, step.GasCost, scope, step.Depth, step.Error)
			} else {
				tracer.CaptureState(step.Pc, fakevm.OpCode(step.Op), step.Gas, step.GasCost, scope, nil, step.Depth, nil)
			}
		}

		previousStepStartedInternalTransaction := previousStep.OpCode == "CREATE" ||
			previousStep.OpCode == "CREATE2" ||
			previousStep.OpCode == "DELEGATECALL" ||
			previousStep.OpCode == "CALL" ||
			previousStep.OpCode == "STATICCALL" ||
			// deprecated ones
			previousStep.OpCode == "CALLCODE"

		// when an internal transaction is detected, the next step contains the context values
		if previousStepStartedInternalTransaction && previousStep.Error == nil {
			// if the previous depth is the same as the current one, this means
			// the internal transaction did not executed any other step and the
			// context is back to the same level. This can happen with pre compiled executions.
			if previousStep.Depth == step.Depth {
				addr, value, input, gas, gasUsed, err := s.getValuesFromInternalTxMemory(previousStep, step)
				if err != nil {
					return nil, err
				}
				from := previousStep.Contract.Address
				if previousStep.OpCode == "CALL" || previousStep.OpCode == "CALLCODE" {
					from = previousStep.Contract.Caller
				}
				tracer.CaptureEnter(fakevm.OpCode(previousStep.Op), from, addr, input, gas, value)
				tracer.CaptureExit(step.ReturnData, gasUsed, previousStep.Error)
			} else {
				value := step.Contract.Value
				if previousStep.OpCode == "STATICCALL" {
					value = nil
				}
				internalTxSteps.Push(instrumentation.InternalTxContext{
					OpCode:       previousStep.OpCode,
					RemainingGas: step.Gas,
				})
				tracer.CaptureEnter(fakevm.OpCode(previousStep.Op), step.Contract.Caller, step.Contract.Address, step.Contract.Input, step.Gas, value)
			}
		}

		// returning from internal transaction
		if previousStep.Depth > step.Depth && previousStep.OpCode != "REVERT" {
			var gasUsed uint64
			var err error
			if errors.Is(previousStep.Error, runtime.ErrOutOfGas) {
				itCtx, err := internalTxSteps.Pop()
				if err != nil {
					return nil, err
				}
				gasUsed = itCtx.RemainingGas
			} else {
				gasUsed, err = s.getGasUsed(internalTxSteps, previousStep, step)
				if err != nil {
					return nil, err
				}
			}
			tracer.CaptureExit(step.ReturnData, gasUsed, previousStep.Error)
		}

		// set StateRoot
		evm.StateDB.SetStateRoot(step.StateRoot.Bytes())

		// set previous step
		previousStep = step
	}

	var err error
	if reverted {
		err = fakevm.ErrExecutionReverted
	} else if result.Err != nil {
		err = result.Err
	}
	tracer.CaptureEnd(trace.Context.Output, trace.Context.GasUsed, err)
	restGas := trace.Context.Gas - trace.Context.GasUsed
	tracer.CaptureTxEnd(restGas)

	return tracer.GetResult()
}

func (s *State) getGasUsed(internalTxContextStack *Stack[instrumentation.InternalTxContext], previousStep, step instrumentation.Step) (uint64, error) {
	itCtx, err := internalTxContextStack.Pop()
	if err != nil {
		return 0, err
	}
	var gasUsed uint64
	if itCtx.OpCode == "CREATE" || itCtx.OpCode == "CREATE2" {
		// if the context was initialized by a CREATE, we should use the contract gas
		gasUsed = previousStep.Contract.Gas - step.Gas
	} else {
		// otherwise we use the step gas
		gasUsed = itCtx.RemainingGas - previousStep.Gas + previousStep.GasCost
	}
	return gasUsed, nil
}

func (s *State) getValuesFromInternalTxMemory(previousStep, step instrumentation.Step) (common.Address, *big.Int, []byte, uint64, uint64, error) {
	if previousStep.OpCode == "DELEGATECALL" || previousStep.OpCode == "CALL" || previousStep.OpCode == "STATICCALL" || previousStep.OpCode == "CALLCODE" {
		gasPos := len(previousStep.Stack) - 1
		addrPos := gasPos - 1

		argsOffsetPos := addrPos - 1
		argsSizePos := argsOffsetPos - 1

		// read tx value if it exists
		var value *big.Int
		stackHasValue := previousStep.OpCode == "CALL" || previousStep.OpCode == "CALLCODE"
		if stackHasValue {
			valuePos := addrPos - 1
			// valueEncoded := step.Stack[valuePos]
			// value = hex.DecodeBig(valueEncoded)
			value = previousStep.Contract.Value

			argsOffsetPos = valuePos - 1
			argsSizePos = argsOffsetPos - 1
		}

		retOffsetPos := argsSizePos - 1
		retSizePos := retOffsetPos - 1

		addr := common.BytesToAddress(previousStep.Stack[addrPos].Bytes())
		argsOffset := previousStep.Stack[argsOffsetPos].Uint64()
		argsSize := previousStep.Stack[argsSizePos].Uint64()
		retOffset := previousStep.Stack[retOffsetPos].Uint64()
		retSize := previousStep.Stack[retSizePos].Uint64()

		input := make([]byte, argsSize)

		if argsOffset > uint64(previousStep.MemorySize) {
			// when none of the bytes can be found in the memory
			// do nothing to keep input as zeroes
		} else if argsOffset+argsSize > uint64(previousStep.MemorySize) {
			// when partial bytes are found in the memory
			// copy just the bytes we have in memory and complement the rest with zeroes
			copy(input[0:argsSize], previousStep.Memory[argsOffset:uint64(previousStep.MemorySize)])
		} else {
			// when all the bytes are found in the memory
			// read the bytes from memory
			copy(input[0:argsSize], previousStep.Memory[argsOffset:argsOffset+argsSize])
		}

		// Compute call memory expansion cost
		memSize := previousStep.MemorySize
		lastMemSizeWord := math.Ceil((float64(memSize) + 31) / 32)                          //nolint:gomnd
		lastMemCost := math.Floor(math.Pow(lastMemSizeWord, 2)/512) + (3 * lastMemSizeWord) //nolint:gomnd

		memSizeWord := math.Ceil((float64(argsOffset+argsSize+31) / 32))                    //nolint:gomnd
		newMemCost := math.Floor(math.Pow(memSizeWord, float64(2))/512) + (3 * memSizeWord) //nolint:gomnd
		callMemCost := newMemCost - lastMemCost

		// Compute return memory expansion cost
		retMemSizeWord := math.Ceil((float64(retOffset) + float64(retSize) + 31) / 32)      //nolint:gomnd
		retNewMemCost := math.Floor(math.Pow(retMemSizeWord, 2)/512) + (3 * retMemSizeWord) //nolint:gomnd
		retMemCost := retNewMemCost - newMemCost
		if retMemCost < 0 {
			retMemCost = 0
		}

		callGasCost := retMemCost + callMemCost + 100 //nolint:gomnd
		gasUsed := float64(previousStep.GasCost) - callGasCost

		// Compute gas sent to call
		gas := float64(previousStep.Gas) - callGasCost
		gas -= math.Floor(gas / 64) //nolint:gomnd

		return addr, value, input, uint64(gas), uint64(gasUsed), nil
	} else {
		createdAddressPos := len(step.Stack) - 1
		addr := common.BytesToAddress(step.Stack[createdAddressPos].Bytes())

		valuePos := len(previousStep.Stack) - 1
		value := previousStep.Stack[valuePos]

		offsetPos := valuePos - 1
		offset := previousStep.Stack[offsetPos].Uint64()

		sizePos := offsetPos - 1
		size := previousStep.Stack[sizePos].Uint64()

		input := make([]byte, size)

		if offset > uint64(previousStep.MemorySize) {
			// when none of the bytes can be found in the memory
			// do nothing to keep input as zeroes
		} else if offset+size > uint64(previousStep.MemorySize) {
			// when partial bytes are found in the memory
			// copy just the bytes we have in memory and complement the rest with zeroes
			copy(input[0:size], previousStep.Memory[offset:uint64(previousStep.MemorySize)])
		} else {
			// when all the bytes are found in the memory
			// read the bytes from memory
			copy(input[0:size], previousStep.Memory[offset:offset+size])
		}

		// Compute gas sent to call
		gas := float64(previousStep.Gas - previousStep.GasCost) //nolint:gomnd
		gas -= math.Floor(gas / 64)                             //nolint:gomnd

		return addr, value, input, uint64(gas), 0, nil
	}
}
