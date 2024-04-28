package state

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/jackc/pgx/v4"
)

// DebugBlock re-executes all block tx to generate its trace
func (s *State) DebugBlock(ctx context.Context, blockNumber uint64, traceConfig TraceConfig, dbTx pgx.Tx) (map[common.Hash]*runtime.ExecutionResult, error) {
	var err error

	// gets the l2 l2Block
	l2Block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	// the old state root is the previous block state root
	var oldStateRoot common.Hash
	previousL2BlockNumber := uint64(0)
	if blockNumber > 0 {
		previousL2BlockNumber = blockNumber - 1
	}
	previousL2Block, err := s.GetL2BlockByNumber(ctx, previousL2BlockNumber, dbTx)
	if err != nil {
		return nil, err
	}
	oldStateRoot = previousL2Block.Root()

	var lastTxInBlock *types.Transaction
	var lastTxReceiptInBlock *types.Receipt
	if len(l2Block.Transactions()) > 0 {
		// gets the last transaction
		lastTxInBlock = l2Block.Transactions()[len(l2Block.Transactions())-1]
		if err != nil {
			return nil, err
		}

		// gets the last tx receipt
		lastTxReceiptInBlock, err = s.GetTransactionReceipt(ctx, lastTxInBlock.Hash(), dbTx)
		if err != nil {
			return nil, err
		}
	}

	// since the executor only stores the state roots by block, we need to
	// execute all the txs in the block
	var txsToEncode []types.Transaction
	var effectivePercentage []uint8
	for i := 0; i < len(l2Block.Transactions()); i++ {
		txsToEncode = append(txsToEncode, *l2Block.Transactions()[i])
		txGasPrice := lastTxInBlock.GasPrice()
		effectiveGasPrice := lastTxReceiptInBlock.EffectiveGasPrice
		egpPercentage, err := CalculateEffectiveGasPricePercentage(txGasPrice, effectiveGasPrice)
		if errors.Is(err, ErrEffectiveGasPriceEmpty) {
			egpPercentage = MaxEffectivePercentage
		} else if err != nil {
			return nil, err
		}
		effectivePercentage = append(effectivePercentage, egpPercentage)
		log.Debugf("trace will reprocess tx: %v", l2Block.Transactions()[i].Hash().String())
	}

	// gets batch that including the l2 block
	batch, err := s.GetBatchByL2BlockNumber(ctx, l2Block.NumberU64(), dbTx)
	if err != nil {
		return nil, err
	}

	forkId := s.GetForkIDByBatchNumber(batch.BatchNumber)

	var responses []*ProcessTransactionResponse
	var startTime, endTime time.Time
	if forkId < FORKID_ETROG {
		traceConfigRequest := &executor.TraceConfig{
			TxHashToGenerateFullTrace: lastTxInBlock.Hash().Bytes(),
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
		responses = convertedResponse.BlockResponses[0].TransactionResponses
	} else {
		traceConfigRequestV2 := &executor.TraceConfigV2{
			//TxHashToGenerateFullTrace: lastTxInBlock.Hash().Bytes(),
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
			// build the raw batch, so we can get the index l1 info tree for the l2 block
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
			s.eventLog.LogExecutorErrorV2(ctx, processBatchResponseV2.Error, processBatchRequestV2)
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
		responses = convertedResponse.BlockResponses[0].TransactionResponses
	}

	// Sanity check
	if len(responses) != len(l2Block.Transactions()) {
		return nil, fmt.Errorf("tx hash not found in executor response")
	}

	results := make(map[common.Hash]*runtime.ExecutionResult)
	for index, response := range responses {
		result := &runtime.ExecutionResult{
			CreateAddress: response.CreateAddress,
			GasLeft:       response.GasLeft,
			GasUsed:       response.GasUsed,
			ReturnValue:   response.ReturnValue,
			StateRoot:     response.StateRoot.Bytes(),
			FullTrace:     response.FullTrace,
			Err:           response.RomError,
		}

		senderAddress, err := GetSender(response.Tx)
		if err != nil {
			return nil, err
		}

		context := instrumentation.Context{
			From:         senderAddress.String(),
			Input:        response.Tx.Data(),
			Gas:          response.Tx.Gas(),
			Value:        response.Tx.Value(),
			Output:       result.ReturnValue,
			GasPrice:     response.Tx.GasPrice().String(),
			OldStateRoot: oldStateRoot,
			Time:         uint64(endTime.Sub(startTime)),
			GasUsed:      result.GasUsed,
		}

		// Fill trace context
		if response.Tx.To() == nil {
			context.Type = "CREATE"
			context.To = result.CreateAddress.Hex()
		} else {
			context.Type = "CALL"
			context.To = response.Tx.To().Hex()
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
			BlockHash:   l2Block.Hash(),
			BlockNumber: l2Block.Number(),
			TxIndex:     index,
			TxHash:      response.TxHash,
		}

		if traceConfig.IsDefaultTracer() {
			structLoggerCfg := structlogger.Config{
				EnableMemory:     traceConfig.EnableMemory,
				DisableStack:     traceConfig.DisableStack,
				DisableStorage:   traceConfig.DisableStorage,
				EnableReturnData: traceConfig.EnableReturnData,
			}
			tracer := structlogger.NewStructLogger(structLoggerCfg)
			receipt, err := s.GetTransactionReceipt(ctx, response.TxHash, dbTx)
			if err != nil {
				return nil, err
			}
			traceResult, err := tracer.ParseTrace(result, *receipt)
			if err != nil {
				return nil, err
			}
			result.TraceResult = traceResult
			results[response.TxHash] = result
			continue
		} else if traceConfig.Is4ByteTracer() {
			tracer, err = native.NewFourByteTracer(tracerContext, traceConfig.TracerConfig)
			if err != nil {
				log.Errorf("debug block: failed to create 4byteTracer, err: %v", err)
				return nil, fmt.Errorf("failed to create 4byteTracer, err: %v", err)
			}
		} else if traceConfig.IsCallTracer() {
			tracer, err = native.NewCallTracer(tracerContext, traceConfig.TracerConfig)
			if err != nil {
				log.Errorf("debug block: failed to create callTracer, err: %v", err)
				return nil, fmt.Errorf("failed to create callTracer, err: %v", err)
			}
		} else if traceConfig.IsFlatCallTracer() {
			// xlayer handle
			tracer, err = native.NewFlatCallTracer(tracerContext, traceConfig.TracerConfig)
			if err != nil {
				log.Errorf("debug block: failed to create flatCallTracer, err: %v", err)
				return nil, fmt.Errorf("failed to create flatCallTracer, err: %v", err)
			}
			tracer = native.SetFlatCallTracerLimit(tracer, traceConfig.Limit)
		} else if traceConfig.IsNoopTracer() {
			tracer, err = native.NewNoopTracer(tracerContext, traceConfig.TracerConfig)
			if err != nil {
				log.Errorf("debug block: failed to create noopTracer, err: %v", err)
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
				log.Errorf("debug block: failed to create jsTracer, err: %v", err)
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
		fmt.Printf("DebugBlock  txhash:%s,  fullTrace:%+v,  traceResult:%s\n", response.TxHash.String(), result.FullTrace, string(traceResult))
		result.TraceResult = traceResult
		results[response.TxHash] = result
	}
	return results, nil
}
