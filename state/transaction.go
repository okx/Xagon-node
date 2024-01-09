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
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/js"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers/native"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation/tracers/sentio"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/google/uuid"
	"github.com/holiman/uint256"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	two uint = 2
)

// GetSender gets the sender from the transaction's signature
func GetSender(tx types.Transaction) (common.Address, error) {
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(&tx)
	if err != nil {
		return common.Address{}, err
	}
	return sender, nil
}

// RlpFieldsToLegacyTx parses the rlp fields slice into a type.LegacyTx
// in this specific order:
//
// required fields:
// [0] Nonce    uint64
// [1] GasPrice *big.Int
// [2] Gas      uint64
// [3] To       *common.Address
// [4] Value    *big.Int
// [5] Data     []byte
//
// optional fields:
// [6] V        *big.Int
// [7] R        *big.Int
// [8] S        *big.Int
func RlpFieldsToLegacyTx(fields [][]byte, v, r, s []byte) (tx *types.LegacyTx, err error) {
	const (
		fieldsSizeWithoutChainID = 6
		fieldsSizeWithChainID    = 7
	)

	if len(fields) < fieldsSizeWithoutChainID {
		return nil, types.ErrTxTypeNotSupported
	}

	nonce := big.NewInt(0).SetBytes(fields[0]).Uint64()
	gasPrice := big.NewInt(0).SetBytes(fields[1])
	gas := big.NewInt(0).SetBytes(fields[2]).Uint64()
	var to *common.Address

	if fields[3] != nil && len(fields[3]) != 0 {
		tmp := common.BytesToAddress(fields[3])
		to = &tmp
	}
	value := big.NewInt(0).SetBytes(fields[4])
	data := fields[5]

	txV := big.NewInt(0).SetBytes(v)
	if len(fields) >= fieldsSizeWithChainID {
		chainID := big.NewInt(0).SetBytes(fields[6])

		// a = chainId * 2
		// b = v - 27
		// c = a + 35
		// v = b + c
		//
		// same as:
		// v = v-27+chainId*2+35
		a := new(big.Int).Mul(chainID, big.NewInt(double))
		b := new(big.Int).Sub(new(big.Int).SetBytes(v), big.NewInt(ether155V))
		c := new(big.Int).Add(a, big.NewInt(etherPre155V))
		txV = new(big.Int).Add(b, c)
	}

	txR := big.NewInt(0).SetBytes(r)
	txS := big.NewInt(0).SetBytes(s)

	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       to,
		Value:    value,
		Data:     data,
		V:        txV,
		R:        txR,
		S:        txS,
	}, nil
}

// StoreTransactions is used by the sequencer to add processed transactions into
// an open batch. If the batch already has txs, the processedTxs must be a super
// set of the existing ones, preserving order.
func (s *State) StoreTransactions(ctx context.Context, batchNumber uint64, processedTxs []*ProcessTransactionResponse, txsEGPLog []*EffectiveGasPriceLog, dbTx pgx.Tx) error {
	if dbTx == nil {
		return ErrDBTxNil
	}

	// check existing txs vs parameter txs
	existingTxs, err := s.GetTxsHashesByBatchNumber(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}
	if err := CheckSupersetBatchTransactions(existingTxs, processedTxs); err != nil {
		return err
	}

	// Check if last batch is closed. Note that it's assumed that only the latest batch can be open
	isBatchClosed, err := s.PostgresStorage.IsBatchClosed(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}
	if isBatchClosed {
		return ErrBatchAlreadyClosed
	}

	processingContext, err := s.GetProcessingContext(ctx, batchNumber, dbTx)
	if err != nil {
		return err
	}

	firstTxToInsert := len(existingTxs)

	for i := firstTxToInsert; i < len(processedTxs); i++ {
		processedTx := processedTxs[i]
		// if the transaction has an intrinsic invalid tx error it means
		// the transaction has not changed the state, so we don't store it
		// and just move to the next
		if executor.IsIntrinsicError(executor.RomErrorCode(processedTx.RomError)) || errors.Is(processedTx.RomError, executor.RomErr(executor.RomError_ROM_ERROR_INVALID_RLP)) {
			continue
		}

		lastL2Block, err := s.GetLastL2Block(ctx, dbTx)
		if err != nil {
			return err
		}

		header := &types.Header{
			Number:     new(big.Int).SetUint64(lastL2Block.Number().Uint64() + 1),
			ParentHash: lastL2Block.Hash(),
			Coinbase:   processingContext.Coinbase,
			Root:       processedTx.StateRoot,
			GasUsed:    processedTx.GasUsed,
			GasLimit:   s.cfg.MaxCumulativeGasUsed,
			Time:       uint64(processingContext.Timestamp.Unix()),
		}
		transactions := []*types.Transaction{&processedTx.Tx}

		receipt := generateReceipt(header.Number, processedTx)
		if !CheckLogOrder(receipt.Logs) {
			return fmt.Errorf("error: logs received from executor are not in order")
		}
		receipts := []*types.Receipt{receipt}

		// Create block to be able to calculate its hash
		block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
		block.ReceivedAt = processingContext.Timestamp

		receipt.BlockHash = block.Hash()

		storeTxsEGPData := []StoreTxEGPData{{EGPLog: nil, EffectivePercentage: uint8(processedTx.EffectivePercentage)}}
		if txsEGPLog != nil {
			storeTxsEGPData[0].EGPLog = txsEGPLog[i]
		}

		// Store L2 block and its transaction
		if err := s.AddL2Block(ctx, batchNumber, block, receipts, storeTxsEGPData, dbTx); err != nil {
			return err
		}
	}
	return nil
}

// DebugTransaction re-executes a tx to generate its trace
func (s *State) DebugTransaction(ctx context.Context, transactionHash common.Hash, traceConfig TraceConfig, dbTx pgx.Tx) (*runtime.ExecutionResult, error) {
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

	// gets the l2 block including the transaction
	block, err := s.GetL2BlockByNumber(ctx, receipt.BlockNumber.Uint64(), dbTx)
	if err != nil {
		return nil, err
	}

	// get the previous L2 Block
	previousBlockNumber := uint64(0)
	if receipt.BlockNumber.Uint64() > 0 {
		previousBlockNumber = receipt.BlockNumber.Uint64() - 1
	}
	previousBlock, err := s.GetL2BlockByNumber(ctx, previousBlockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	// gets batch that including the l2 block
	batch, err := s.GetBatchByL2BlockNumber(ctx, block.NumberU64(), dbTx)
	if err != nil {
		return nil, err
	}

	forkId := s.GetForkIDByBatchNumber(batch.BatchNumber)

	// gets batch that including the previous l2 block
	previousBatch, err := s.GetBatchByL2BlockNumber(ctx, previousBlock.NumberU64(), dbTx)
	if err != nil {
		return nil, err
	}

	// generate batch l2 data for the transaction
	batchL2Data, err := EncodeTransactions([]types.Transaction{*tx}, []uint8{MaxEffectivePercentage}, forkId)
	if err != nil {
		return nil, err
	}

	var txHashToGenerateCallTrace []byte
	var txHashToGenerateExecuteTrace []byte

	if traceConfig.IsDefaultTracer() {
		txHashToGenerateExecuteTrace = transactionHash.Bytes()
	} else {
		txHashToGenerateCallTrace = transactionHash.Bytes()
	}

	// Create Batch
	traceConfigRequest := &executor.TraceConfig{
		TxHashToGenerateCallTrace:    txHashToGenerateCallTrace,
		TxHashToGenerateExecuteTrace: txHashToGenerateExecuteTrace,
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

	oldStateRoot := previousBlock.Root()
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:     batch.BatchNumber - 1,
		OldStateRoot:    oldStateRoot.Bytes(),
		OldAccInputHash: previousBatch.AccInputHash.Bytes(),

		BatchL2Data:      batchL2Data,
		GlobalExitRoot:   batch.GlobalExitRoot.Bytes(),
		EthTimestamp:     uint64(batch.Timestamp.Unix()),
		Coinbase:         batch.Coinbase.String(),
		UpdateMerkleTree: cFalse,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkId,
		TraceConfig:      traceConfigRequest,
		ContextId:        uuid.NewString(),
	}

	// Send Batch to the Executor
	startTime := time.Now()
	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	endTime := time.Now()
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

	// Sanity check
	response := convertedResponse.Responses[0]
	log.Debugf(response.TxHash.String())
	if response.TxHash != transactionHash {
		return nil, fmt.Errorf("tx hash not found in executor response")
	}

	// const path = "/Users/thiago/github.com/0xPolygonHermez/zkevm-node/dist/%v.json"
	// filePath := fmt.Sprintf(path, "EXECUTOR_processBatchResponse")
	// c, _ := json.MarshalIndent(processBatchResponse, "", "    ")
	// os.WriteFile(filePath, c, 0644)

	// filePath = fmt.Sprintf(path, "NODE_execution_trace")
	// c, _ = json.MarshalIndent(response.ExecutionTrace, "", "    ")
	// os.WriteFile(filePath, c, 0644)

	// filePath = fmt.Sprintf(path, "NODE_call_trace")
	// c, _ = json.MarshalIndent(response.CallTrace, "", "    ")
	// os.WriteFile(filePath, c, 0644)

	result := &runtime.ExecutionResult{
		CreateAddress: response.CreateAddress,
		GasLeft:       response.GasLeft,
		GasUsed:       response.GasUsed,
		ReturnValue:   response.ReturnValue,
		StateRoot:     response.StateRoot.Bytes(),
		StructLogs:    response.ExecutionTrace,
		ExecutorTrace: response.CallTrace,
		Err:           response.RomError,
	}

	// if is the default trace, return the result
	if traceConfig.IsDefaultTracer() {
		return result, nil
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

	result.ExecutorTrace.Context = context

	gasPrice, ok := new(big.Int).SetString(context.GasPrice, encoding.Base10)
	if !ok {
		log.Errorf("debug transaction: failed to parse gasPrice")
		return nil, fmt.Errorf("failed to parse gasPrice")
	}

	tracerContext := &tracers.Context{
		BlockHash:   receipt.BlockHash,
		BlockNumber: receipt.BlockNumber,
		TxIndex:     int(receipt.TransactionIndex),
		TxHash:      transactionHash,
	}

	var customTracer tracers.Tracer
	if traceConfig.Is4ByteTracer() {
		customTracer, err = native.NewFourByteTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create 4byteTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create 4byteTracer, err: %v", err)
		}
	} else if traceConfig.IsCallTracer() {
		customTracer, err = native.NewCallTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create callTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create callTracer, err: %v", err)
		}
	} else if traceConfig.IsFlatCallTracer() {
		customTracer, err = native.NewFlatCallTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create flatCallTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create flatCallTracer, err: %v", err)
		}
		customTracer = native.SetFlatCallTracerLimit(customTracer, traceConfig.Limit)
	} else if traceConfig.IsNoopTracer() {
		customTracer, err = native.NewNoopTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create noopTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create noopTracer, err: %v", err)
		}
	} else if traceConfig.IsPrestateTracer() {
		customTracer, err = native.NewPrestateTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create prestateTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create prestateTracer, err: %v", err)
		}
	} else if traceConfig.IsJSCustomTracer() {
		customTracer, err = js.NewJsTracer(*traceConfig.Tracer, tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create jsTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create jsTracer, err: %v", err)
		}
	} else if *traceConfig.Tracer == "sentioTracer" {
		customTracer, err = sentio.NewSentioTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create sentioTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create sentioTracer, err: %v", err)
		}
	} else if *traceConfig.Tracer == "sentioPrestateTracer" {
		customTracer, err = sentio.NewSentioPrestateTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create sentioPrestateTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create sentioPrestateTracer, err: %v", err)
		}
	} else {
		return nil, fmt.Errorf("invalid tracer: %v, err: %v", traceConfig.Tracer, err)
	}

	fakeDB := &FakeDB{State: s, stateRoot: batch.StateRoot.Bytes()}
	evm := fakevm.NewFakeEVM(fakevm.BlockContext{BlockNumber: big.NewInt(1)}, fakevm.TxContext{GasPrice: gasPrice}, fakeDB, params.TestChainConfig, fakevm.Config{Debug: true, Tracer: customTracer})

	traceResult, err := s.buildTrace(evm, result, customTracer)
	if err != nil {
		log.Errorf("debug transaction: failed parse the trace using the tracer: %v", err)
		return nil, fmt.Errorf("failed parse the trace using the tracer: %v", err)
	}

	result.ExecutorTraceResult = traceResult

	return result, nil
}

// TraceCall execute a tx to generate its trace
func (s *State) TraceCall(ctx context.Context, blockNumber uint64, sender common.Address, tx *types.Transaction, traceConfig TraceConfig, dbTx pgx.Tx) (*runtime.ExecutionResult, error) {
	// gets the l2 block including the transaction
	block, err := s.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	// get the previous L2 Block
	previousBlockNumber := uint64(0)
	if blockNumber > 0 {
		previousBlockNumber = blockNumber - 1
	}
	previousBlock, err := s.GetL2BlockByNumber(ctx, previousBlockNumber, dbTx)
	if err != nil {
		return nil, err
	}

	// gets batch that including the l2 block
	batch, err := s.GetBatchByL2BlockNumber(ctx, block.NumberU64(), dbTx)
	if err != nil {
		return nil, err
	}

	forkId := s.GetForkIDByBatchNumber(batch.BatchNumber)
	// gets batch that including the previous l2 block
	previousBatch, err := s.GetBatchByL2BlockNumber(ctx, previousBlock.NumberU64(), dbTx)
	if err != nil {
		return nil, err
	}

	loadedNonce, err := s.tree.GetNonce(ctx, sender, previousBlock.Root().Bytes())
	if err != nil {
		return nil, err
	}
	nonce := loadedNonce.Uint64()

	v, _ := new(big.Int).SetString("0x1c", 0)
	r, _ := new(big.Int).SetString("0xa54492cfacf71aef702421b7fbc70636537a7b2fbe5718c5ed970a001bb7756b", 0)
	ss, _ := new(big.Int).SetString("0x2e9fb27acc75955b898f0b12ec52aa34bf08f01db654374484b80bf12f0d841e", 0)
	tx = types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       tx.To(),
		Value:    tx.Value(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Data:     tx.Data(),
		V:        v,
		R:        r,
		S:        ss,
	})

	transactionHash := tx.Hash()

	// generate batch l2 data for the transaction
	batchL2Data, err := EncodeTransactions([]types.Transaction{*tx}, []uint8{MaxEffectivePercentage}, forkId)
	if err != nil {
		return nil, err
	}

	var txHashToGenerateCallTrace []byte
	var txHashToGenerateExecuteTrace []byte

	if traceConfig.IsDefaultTracer() {
		txHashToGenerateExecuteTrace = transactionHash.Bytes()
	} else {
		txHashToGenerateCallTrace = transactionHash.Bytes()
	}

	// Create Batch
	traceConfigRequest := &executor.TraceConfig{
		TxHashToGenerateCallTrace:    txHashToGenerateCallTrace,
		TxHashToGenerateExecuteTrace: txHashToGenerateExecuteTrace,
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

	oldStateRoot := previousBlock.Root()
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      batch.BatchNumber - 1,
		OldStateRoot:     oldStateRoot.Bytes(),
		OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
		From:             sender.String(),
		BatchL2Data:      batchL2Data,
		GlobalExitRoot:   batch.GlobalExitRoot.Bytes(),
		EthTimestamp:     uint64(batch.Timestamp.Unix()),
		Coinbase:         batch.Coinbase.String(),
		UpdateMerkleTree: cFalse,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkId,
		TraceConfig:      traceConfigRequest,
	}

	// Send Batch to the Executor
	startTime := time.Now()
	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	endTime := time.Now()
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

	// Sanity check
	response := convertedResponse.Responses[0]
	log.Debugf(response.TxHash.String())
	if response.TxHash != transactionHash {
		return nil, fmt.Errorf("tx hash not found in executor response")
	}

	// const path = "/Users/thiago/github.com/0xPolygonHermez/zkevm-node/dist/%v.json"
	// filePath := fmt.Sprintf(path, "EXECUTOR_processBatchResponse")
	// c, _ := json.MarshalIndent(processBatchResponse, "", "    ")
	// os.WriteFile(filePath, c, 0644)

	// filePath = fmt.Sprintf(path, "NODE_execution_trace")
	// c, _ = json.MarshalIndent(response.ExecutionTrace, "", "    ")
	// os.WriteFile(filePath, c, 0644)

	// filePath = fmt.Sprintf(path, "NODE_call_trace")
	// c, _ = json.MarshalIndent(response.CallTrace, "", "    ")
	// os.WriteFile(filePath, c, 0644)
	result := &runtime.ExecutionResult{
		CreateAddress: response.CreateAddress,
		GasLeft:       response.GasLeft,
		GasUsed:       response.GasUsed,
		ReturnValue:   response.ReturnValue,
		StateRoot:     response.StateRoot.Bytes(),
		StructLogs:    response.ExecutionTrace,
		ExecutorTrace: response.CallTrace,
		Err:           response.RomError,
	}

	if result.Err != nil && errors.Is(result.Err, runtime.ErrOutOfCountersStep) {
		result.Err = fmt.Errorf("%v. arith: %d, memAlign: %d, binary: %d, keccakF: %d, poseidonG: %d, paddingPG %d",
			result.Err, processBatchResponse.CntArithmetics, processBatchResponse.CntMemAligns,
			processBatchResponse.CntBinaries, processBatchResponse.CntKeccakHashes, processBatchResponse.CntPoseidonHashes,
			processBatchResponse.CntPoseidonPaddings)
	}

	// if is the default trace, return the result
	if traceConfig.IsDefaultTracer() {
		return result, nil
	}

	//senderAddress, err := GetSender(*tx)
	//if err != nil {
	//	return nil, err
	//}

	context := instrumentation.Context{
		From:         sender.String(),
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

	result.ExecutorTrace.Context = context

	gasPrice, ok := new(big.Int).SetString(context.GasPrice, encoding.Base10)
	if !ok {
		log.Errorf("debug transaction: failed to parse gasPrice")
		return nil, fmt.Errorf("failed to parse gasPrice")
	}

	tracerContext := &tracers.Context{
		BlockHash:   block.Hash(),
		BlockNumber: big.NewInt(int64(blockNumber)),
		TxIndex:     0,
		TxHash:      transactionHash,
	}

	var customTracer tracers.Tracer
	if traceConfig.Is4ByteTracer() {
		customTracer, err = native.NewFourByteTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create 4byteTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create 4byteTracer, err: %v", err)
		}
	} else if traceConfig.IsCallTracer() {
		customTracer, err = native.NewCallTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create callTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create callTracer, err: %v", err)
		}
	} else if traceConfig.IsNoopTracer() {
		customTracer, err = native.NewNoopTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create noopTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create noopTracer, err: %v", err)
		}
	} else if traceConfig.IsPrestateTracer() {
		customTracer, err = native.NewPrestateTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create prestateTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create prestateTracer, err: %v", err)
		}
	} else if traceConfig.IsJSCustomTracer() {
		customTracer, err = js.NewJsTracer(*traceConfig.Tracer, tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create jsTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create jsTracer, err: %v", err)
		}
	} else if *traceConfig.Tracer == "sentioTracer" {
		customTracer, err = sentio.NewSentioTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create sentioTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create sentioTracer, err: %v", err)
		}
	} else if *traceConfig.Tracer == "sentioPrestateTracer" {
		customTracer, err = sentio.NewSentioPrestateTracer(tracerContext, traceConfig.TracerConfig)
		if err != nil {
			log.Errorf("debug transaction: failed to create sentioPrestateTracer, err: %v", err)
			return nil, fmt.Errorf("failed to create sentioPrestateTracer, err: %v", err)
		}
	} else {
		return nil, fmt.Errorf("invalid tracer: %v, err: %v", traceConfig.Tracer, err)
	}

	fakeDB := &FakeDB{State: s, stateRoot: batch.StateRoot.Bytes()}
	evm := fakevm.NewFakeEVM(fakevm.BlockContext{BlockNumber: big.NewInt(1)}, fakevm.TxContext{GasPrice: gasPrice}, fakeDB, params.TestChainConfig, fakevm.Config{Debug: true, Tracer: customTracer})

	traceResult, err := s.buildTrace(evm, result, customTracer)
	if err != nil {
		log.Errorf("debug transaction: failed parse the trace using the tracer: %v", err)
		return nil, fmt.Errorf("failed parse the trace using the tracer: %v", err)
	}

	result.ExecutorTraceResult = traceResult

	return result, nil
}

// ParseTheTraceUsingTheTracer parses the given trace with the given tracer.
func (s *State) buildTrace(evm *fakevm.FakeEVM, result *runtime.ExecutionResult, tracer tracers.Tracer) (json.RawMessage, error) {
	trace := result.ExecutorTrace
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
				tracer.CaptureState(step.Pc, fakevm.OpCode(step.Op), step.Gas, step.GasCost, scope, step.ReturnData, step.Depth, nil)
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

// PreProcessTransaction processes the transaction in order to calculate its zkCounters before adding it to the pool
func (s *State) PreProcessTransaction(ctx context.Context, tx *types.Transaction, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	sender, err := GetSender(*tx)
	if err != nil {
		return nil, err
	}

	response, err := s.internalProcessUnsignedTransaction(ctx, tx, sender, nil, false, dbTx)
	if err != nil {
		return response, err
	}

	return response, nil
}

// ProcessUnsignedTransaction processes the given unsigned transaction.
func (s *State) ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, noZKEVMCounters bool, dbTx pgx.Tx) (*runtime.ExecutionResult, error) {
	result := new(runtime.ExecutionResult)
	response, err := s.internalProcessUnsignedTransaction(ctx, tx, senderAddress, l2BlockNumber, noZKEVMCounters, dbTx)
	if err != nil {
		return nil, err
	}

	r := response.Responses[0]
	result.ReturnValue = r.ReturnValue
	result.GasLeft = r.GasLeft
	result.GasUsed = r.GasUsed
	result.CreateAddress = r.CreateAddress
	result.StateRoot = r.StateRoot.Bytes()

	if errors.Is(r.RomError, runtime.ErrExecutionReverted) {
		result.Err = constructErrorFromRevert(r.RomError, r.ReturnValue)
	} else {
		result.Err = r.RomError
	}

	return result, nil
}

// ProcessUnsignedTransaction processes the given unsigned transaction.
func (s *State) internalProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, noZKEVMCounters bool, dbTx pgx.Tx) (*ProcessBatchResponse, error) {
	var attempts = 1

	if s.executorClient == nil {
		return nil, ErrExecutorNil
	}
	if s.tree == nil {
		return nil, ErrStateTreeNil
	}
	lastBatches, l2BlockStateRoot, err := s.PostgresStorage.GetLastNBatchesByL2BlockNumber(ctx, l2BlockNumber, two, dbTx)
	if err != nil {
		return nil, err
	}

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	stateRoot := l2BlockStateRoot
	timestamp := uint64(lastBatch.Timestamp.Unix())
	if l2BlockNumber != nil {
		l2Block, err := s.GetL2BlockByNumber(ctx, *l2BlockNumber, dbTx)
		if err != nil {
			return nil, err
		}
		stateRoot = l2Block.Root()

		latestL2BlockNumber, err := s.PostgresStorage.GetLastL2BlockNumber(ctx, dbTx)
		if err != nil {
			return nil, err
		}

		if *l2BlockNumber == latestL2BlockNumber {
			timestamp = uint64(time.Now().Unix())
		}
	}

	forkID := s.GetForkIDByBatchNumber(lastBatch.BatchNumber)
	loadedNonce, err := s.tree.GetNonce(ctx, senderAddress, stateRoot.Bytes())
	if err != nil {
		return nil, err
	}
	nonce := loadedNonce.Uint64()

	batchL2Data, err := EncodeUnsignedTransaction(*tx, s.cfg.ChainID, &nonce, forkID)
	if err != nil {
		log.Errorf("error encoding unsigned transaction ", err)
		return nil, err
	}

	// Create Batch
	processBatchRequest := &executor.ProcessBatchRequest{
		OldBatchNum:      lastBatch.BatchNumber,
		BatchL2Data:      batchL2Data,
		From:             senderAddress.String(),
		OldStateRoot:     stateRoot.Bytes(),
		GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
		OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
		EthTimestamp:     timestamp,
		Coinbase:         lastBatch.Coinbase.String(),
		UpdateMerkleTree: cFalse,
		ChainId:          s.cfg.ChainID,
		ForkId:           forkID,
		ContextId:        uuid.NewString(),
	}

	if noZKEVMCounters {
		processBatchRequest.NoCounters = cTrue
	}

	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.From]: %v", processBatchRequest.From)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.globalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
	log.Debugf("internalProcessUnsignedTransaction[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)

	// Send Batch to the Executor
	processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
	if err != nil {
		if status.Code(err) == codes.ResourceExhausted || (processBatchResponse != nil && processBatchResponse.Error == executor.ExecutorError(executor.ExecutorError_EXECUTOR_ERROR_DB_ERROR)) {
			log.Errorf("error processing unsigned transaction ", err)
			for attempts < s.cfg.MaxResourceExhaustedAttempts {
				time.Sleep(s.cfg.WaitOnResourceExhaustion.Duration)
				log.Errorf("retrying to process unsigned transaction")
				processBatchResponse, err = s.executorClient.ProcessBatch(ctx, processBatchRequest)
				if status.Code(err) == codes.ResourceExhausted || (processBatchResponse != nil && processBatchResponse.Error == executor.ExecutorError(executor.ExecutorError_EXECUTOR_ERROR_DB_ERROR)) {
					log.Errorf("error processing unsigned transaction ", err)
					attempts++
					continue
				}
				break
			}
		}

		if err != nil {
			if status.Code(err) == codes.ResourceExhausted || (processBatchResponse != nil && processBatchResponse.Error == executor.ExecutorError(executor.ExecutorError_EXECUTOR_ERROR_DB_ERROR)) {
				log.Error("reporting error as time out")
				return nil, runtime.ErrGRPCResourceExhaustedAsTimeout
			}
			// Log the error
			event := &event.Event{
				ReceivedAt:  time.Now(),
				Source:      event.Source_Node,
				Level:       event.Level_Error,
				EventID:     event.EventID_ExecutorError,
				Description: fmt.Sprintf("error processing unsigned transaction %s: %v", tx.Hash(), err),
			}

			err2 := s.eventLog.LogEvent(context.Background(), event)
			if err2 != nil {
				log.Errorf("error logging event %v", err2)
			}
			log.Errorf("error processing unsigned transaction ", err)
			return nil, err
		}
	}

	if err == nil && processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
		err = executor.ExecutorErr(processBatchResponse.Error)
		s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
		return nil, err
	}

	response, err := s.convertToProcessBatchResponse(processBatchResponse)
	if err != nil {
		return nil, err
	}

	if processBatchResponse.Responses[0].Error != executor.RomError_ROM_ERROR_NO_ERROR {
		err := executor.RomErr(processBatchResponse.Responses[0].Error)
		if !isEVMRevertError(err) {
			return response, err
		}
	}

	return response, nil
}

// isContractCreation checks if the tx is a contract creation
func (s *State) isContractCreation(tx *types.Transaction) bool {
	return tx.To() == nil && len(tx.Data()) > 0
}

// StoreTransaction is used by the sequencer and trusted state synchronizer to add process a transaction.
func (s *State) StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *ProcessTransactionResponse, coinbase common.Address, timestamp uint64, egpLog *EffectiveGasPriceLog, dbTx pgx.Tx) (*types.Header, error) {
	if dbTx == nil {
		return nil, ErrDBTxNil
	}

	// if the transaction has an intrinsic invalid tx error it means
	// the transaction has not changed the state, so we don't store it
	if executor.IsIntrinsicError(executor.RomErrorCode(processedTx.RomError)) {
		return nil, nil
	}

	lastL2Block, err := s.GetLastL2Block(ctx, dbTx)
	if err != nil {
		return nil, err
	}

	header := &types.Header{
		Number:     new(big.Int).SetUint64(lastL2Block.Number().Uint64() + 1),
		ParentHash: lastL2Block.Hash(),
		Coinbase:   coinbase,
		Root:       processedTx.StateRoot,
		GasUsed:    processedTx.GasUsed,
		GasLimit:   s.cfg.MaxCumulativeGasUsed,
		Time:       timestamp,
	}
	transactions := []*types.Transaction{&processedTx.Tx}

	receipt := generateReceipt(header.Number, processedTx)
	receipts := []*types.Receipt{receipt}

	// Create block to be able to calculate its hash
	block := types.NewBlock(header, transactions, []*types.Header{}, receipts, &trie.StackTrie{})
	block.ReceivedAt = time.Unix(int64(timestamp), 0)

	receipt.BlockHash = block.Hash()

	storeTxsEGPData := []StoreTxEGPData{{EGPLog: egpLog, EffectivePercentage: uint8(processedTx.EffectivePercentage)}}

	// Store L2 block and its transaction
	if err := s.AddL2Block(ctx, batchNumber, block, receipts, storeTxsEGPData, dbTx); err != nil {
		return nil, err
	}

	return block.Header(), nil
}

// CheckSupersetBatchTransactions verifies that processedTransactions is a
// superset of existingTxs and that the existing txs have the same order,
// returns a non-nil error if that is not the case.
func CheckSupersetBatchTransactions(existingTxHashes []common.Hash, processedTxs []*ProcessTransactionResponse) error {
	if len(existingTxHashes) > len(processedTxs) {
		return ErrExistingTxGreaterThanProcessedTx
	}
	for i, existingTxHash := range existingTxHashes {
		if existingTxHash != processedTxs[i].TxHash {
			return ErrOutOfOrderProcessedTx
		}
	}
	return nil
}

// EstimateGas for a transaction
func (s *State) EstimateGas(transaction *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, dbTx pgx.Tx) (uint64, []byte, error) {
	const ethTransferGas = 21000

	var lowEnd uint64
	var highEnd uint64

	ctx := context.Background()

	lastBatches, l2BlockStateRoot, err := s.PostgresStorage.GetLastNBatchesByL2BlockNumber(ctx, l2BlockNumber, two, dbTx)
	if err != nil {
		return 0, nil, err
	}

	stateRoot := l2BlockStateRoot
	if l2BlockNumber != nil {
		l2Block, err := s.GetL2BlockByNumber(ctx, *l2BlockNumber, dbTx)
		if err != nil {
			return 0, nil, err
		}
		stateRoot = l2Block.Root()
	}

	loadedNonce, err := s.tree.GetNonce(ctx, senderAddress, stateRoot.Bytes())
	if err != nil {
		return 0, nil, err
	}
	nonce := loadedNonce.Uint64()

	// Get latest batch from the database to get globalExitRoot and Timestamp
	lastBatch := lastBatches[0]

	// Get batch before latest to get state root and local exit root
	previousBatch := lastBatches[0]
	if len(lastBatches) > 1 {
		previousBatch = lastBatches[1]
	}

	lowEnd, err = core.IntrinsicGas(transaction.Data(), transaction.AccessList(), s.isContractCreation(transaction), true, false, false)
	if err != nil {
		return 0, nil, err
	}

	if lowEnd == ethTransferGas && transaction.To() != nil {
		code, err := s.tree.GetCode(ctx, *transaction.To(), stateRoot.Bytes())
		if err != nil {
			log.Warnf("error while getting transaction.to() code %v", err)
		} else if len(code) == 0 {
			return lowEnd, nil, nil
		}
	}

	if transaction.Gas() != 0 && transaction.Gas() > lowEnd {
		highEnd = transaction.Gas()
	} else {
		highEnd = s.cfg.MaxCumulativeGasUsed
	}

	var availableBalance *big.Int

	if senderAddress != ZeroAddress {
		senderBalance, err := s.tree.GetBalance(ctx, senderAddress, stateRoot.Bytes())
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				senderBalance = big.NewInt(0)
			} else {
				return 0, nil, err
			}
		}

		availableBalance = new(big.Int).Set(senderBalance)

		if transaction.Value() != nil {
			if transaction.Value().Cmp(availableBalance) > 0 {
				return 0, nil, ErrInsufficientFunds
			}

			availableBalance.Sub(availableBalance, transaction.Value())
		}
	}

	if transaction.GasPrice().BitLen() != 0 && // Gas price has been set
		availableBalance != nil && // Available balance is found
		availableBalance.Cmp(big.NewInt(0)) > 0 { // Available balance > 0
		gasAllowance := new(big.Int).Div(availableBalance, transaction.GasPrice())

		// Check the gas allowance for this account, make sure high end is capped to it
		if gasAllowance.IsUint64() && highEnd > gasAllowance.Uint64() {
			log.Debugf("Gas estimation high-end capped by allowance [%d]", gasAllowance.Uint64())
			highEnd = gasAllowance.Uint64()
		}
	}

	// Run the transaction with the specified gas value.
	// Returns a status indicating if the transaction failed, if it was reverted and the accompanying error
	testTransaction := func(gas uint64, nonce uint64, shouldOmitErr bool) (failed, reverted bool, gasUsed uint64, returnValue []byte, err error) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       transaction.To(),
			Value:    transaction.Value(),
			Gas:      gas,
			GasPrice: transaction.GasPrice(),
			Data:     transaction.Data(),
		})

		forkID := s.GetForkIDByBatchNumber(lastBatch.BatchNumber)

		batchL2Data, err := EncodeUnsignedTransaction(*tx, s.cfg.ChainID, nil, forkID)
		if err != nil {
			log.Errorf("error encoding unsigned transaction ", err)
			return false, false, gasUsed, nil, err
		}

		// Create a batch to be sent to the executor
		processBatchRequest := &executor.ProcessBatchRequest{
			OldBatchNum:      lastBatch.BatchNumber,
			BatchL2Data:      batchL2Data,
			From:             senderAddress.String(),
			OldStateRoot:     stateRoot.Bytes(),
			GlobalExitRoot:   lastBatch.GlobalExitRoot.Bytes(),
			OldAccInputHash:  previousBatch.AccInputHash.Bytes(),
			EthTimestamp:     uint64(lastBatch.Timestamp.Unix()),
			Coinbase:         lastBatch.Coinbase.String(),
			UpdateMerkleTree: cFalse,
			ChainId:          s.cfg.ChainID,
			ForkId:           forkID,
			ContextId:        uuid.NewString(),
		}

		log.Debugf("EstimateGas[processBatchRequest.OldBatchNum]: %v", processBatchRequest.OldBatchNum)
		// log.Debugf("EstimateGas[processBatchRequest.BatchL2Data]: %v", hex.EncodeToHex(processBatchRequest.BatchL2Data))
		log.Debugf("EstimateGas[processBatchRequest.From]: %v", processBatchRequest.From)
		log.Debugf("EstimateGas[processBatchRequest.OldStateRoot]: %v", hex.EncodeToHex(processBatchRequest.OldStateRoot))
		log.Debugf("EstimateGas[processBatchRequest.globalExitRoot]: %v", hex.EncodeToHex(processBatchRequest.GlobalExitRoot))
		log.Debugf("EstimateGas[processBatchRequest.OldAccInputHash]: %v", hex.EncodeToHex(processBatchRequest.OldAccInputHash))
		log.Debugf("EstimateGas[processBatchRequest.EthTimestamp]: %v", processBatchRequest.EthTimestamp)
		log.Debugf("EstimateGas[processBatchRequest.Coinbase]: %v", processBatchRequest.Coinbase)
		log.Debugf("EstimateGas[processBatchRequest.UpdateMerkleTree]: %v", processBatchRequest.UpdateMerkleTree)
		log.Debugf("EstimateGas[processBatchRequest.ChainId]: %v", processBatchRequest.ChainId)
		log.Debugf("EstimateGas[processBatchRequest.ForkId]: %v", processBatchRequest.ForkId)
		log.Debugf("EstimateGas[processBatchRequest.ContextId]: %v", processBatchRequest.ContextId)

		txExecutionOnExecutorTime := time.Now()
		processBatchResponse, err := s.executorClient.ProcessBatch(ctx, processBatchRequest)
		log.Debugf("executor time: %vms", time.Since(txExecutionOnExecutorTime).Milliseconds())
		if err != nil {
			log.Errorf("error estimating gas: %v", err)
			return false, false, gasUsed, nil, err
		}
		if processBatchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR {
			err = executor.ExecutorErr(processBatchResponse.Error)
			s.eventLog.LogExecutorError(ctx, processBatchResponse.Error, processBatchRequest)
			return false, false, gasUsed, nil, err
		}
		gasUsed = processBatchResponse.Responses[0].GasUsed

		// Check if an out of gas error happened during EVM execution
		if processBatchResponse.Responses[0].Error != executor.RomError_ROM_ERROR_NO_ERROR {
			err := executor.RomErr(processBatchResponse.Responses[0].Error)

			if (isGasEVMError(err) || isGasApplyError(err)) && shouldOmitErr {
				// Specifying the transaction failed, but not providing an error
				// is an indication that a valid error occurred due to low gas,
				// which will increase the lower bound for the search
				return true, false, gasUsed, nil, nil
			}

			if isEVMRevertError(err) {
				// The EVM reverted during execution, attempt to extract the
				// error message and return it
				returnValue := processBatchResponse.Responses[0].ReturnValue
				return true, true, gasUsed, returnValue, constructErrorFromRevert(err, returnValue)
			}

			return true, false, gasUsed, nil, err
		}

		return false, false, gasUsed, nil, nil
	}

	txExecutions := []time.Duration{}
	var totalExecutionTime time.Duration

	// Check if the highEnd is a good value to make the transaction pass
	failed, reverted, gasUsed, returnValue, err := testTransaction(highEnd, nonce, false)
	log.Debugf("Estimate gas. Trying to execute TX with %v gas", highEnd)
	if failed {
		if reverted {
			return 0, returnValue, err
		}

		// The transaction shouldn't fail, for whatever reason, at highEnd
		return 0, nil, fmt.Errorf(
			"unable to apply transaction even for the highest gas limit %d: %w",
			highEnd,
			err,
		)
	}

	if lowEnd < gasUsed {
		lowEnd = gasUsed
	}

	// Start the binary search for the lowest possible gas price
	for (lowEnd < highEnd) && (highEnd-lowEnd) > 4096 {
		txExecutionStart := time.Now()
		mid := (lowEnd + highEnd) / uint64(two)

		log.Debugf("Estimate gas. Trying to execute TX with %v gas", mid)

		failed, reverted, _, _, testErr := testTransaction(mid, nonce, true)
		executionTime := time.Since(txExecutionStart)
		totalExecutionTime += executionTime
		txExecutions = append(txExecutions, executionTime)
		if testErr != nil && !reverted {
			// Reverts are ignored in the binary search, but are checked later on
			// during the execution for the optimal gas limit found
			return 0, nil, testErr
		}

		if failed {
			// If the transaction failed => increase the gas
			lowEnd = mid + 1
		} else {
			// If the transaction didn't fail => make this ok value the high end
			highEnd = mid
		}
	}

	executions := int64(len(txExecutions))
	if executions > 0 {
		log.Infof("EstimateGas executed TX %v %d times in %d milliseconds", transaction.Hash(), executions, totalExecutionTime.Milliseconds())
	} else {
		log.Error("Estimate gas. Tx not executed")
	}
	return highEnd, nil, nil
}

// Checks if executor level valid gas errors occurred
func isGasApplyError(err error) bool {
	return errors.Is(err, ErrNotEnoughIntrinsicGas)
}

// Checks if EVM level valid gas errors occurred
func isGasEVMError(err error) bool {
	return errors.Is(err, runtime.ErrOutOfGas)
}

// Checks if the EVM reverted during execution
func isEVMRevertError(err error) bool {
	return errors.Is(err, runtime.ErrExecutionReverted)
}
