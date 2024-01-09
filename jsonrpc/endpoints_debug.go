package jsonrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/instrumentation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/jackc/pgx/v4"
)

var defaultTraceConfig = &traceConfig{
	DisableStorage:   false,
	DisableStack:     false,
	EnableMemory:     false,
	EnableReturnData: false,
	Tracer:           nil,
	Limit:            0,
}

// DebugEndpoints is the debug jsonrpc endpoint
type DebugEndpoints struct {
	cfg      Config
	state    types.StateInterface
	etherman types.EthermanInterface
	txMan    DBTxManager
}

// NewDebugEndpoints returns DebugEndpoints
func NewDebugEndpoints(cfg Config, state types.StateInterface, etherman types.EthermanInterface) *DebugEndpoints {
	return &DebugEndpoints{
		cfg:      cfg,
		state:    state,
		etherman: etherman,
	}
}

type traceConfig struct {
	DisableStorage   bool            `json:"disableStorage"`
	DisableStack     bool            `json:"disableStack"`
	EnableMemory     bool            `json:"enableMemory"`
	EnableReturnData bool            `json:"enableReturnData"`
	Tracer           *string         `json:"tracer"`
	TracerConfig     json.RawMessage `json:"tracerConfig"`
	Limit            int             `json:"limit"`
}

// StructLogRes represents the debug trace information for each opcode
type StructLogRes struct {
	Pc            uint64             `json:"pc"`
	Op            string             `json:"op"`
	Gas           uint64             `json:"gas"`
	GasCost       uint64             `json:"gasCost"`
	Depth         int                `json:"depth"`
	Error         string             `json:"error,omitempty"`
	Stack         *[]types.ArgBig    `json:"stack,omitempty"`
	Memory        *[]string          `json:"memory,omitempty"`
	Storage       *map[string]string `json:"storage,omitempty"`
	RefundCounter uint64             `json:"refund,omitempty"`
}

type traceTransactionResponse struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue interface{}    `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

type traceBlockTransactionResponse struct {
	Result interface{} `json:"result"`
}

type traceBatchTransactionResponse struct {
	TxHash common.Hash `json:"txHash"`
	Result interface{} `json:"result"`
}

// TraceTransaction creates a response for debug_traceTransaction request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtracetransaction
func (d *DebugEndpoints) TraceTransaction(hash types.ArgHash, cfg *traceConfig) (interface{}, types.Error) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		return d.buildTraceTransaction(ctx, hash.Hash(), cfg, dbTx)
	})
}

// TraceBlockByNumber creates a response for debug_traceBlockByNumber request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtraceblockbynumber
func (d *DebugEndpoints) TraceBlockByNumber(number types.BlockNumber, cfg *traceConfig) (interface{}, types.Error) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		blockNumber, rpcErr := number.GetNumericBlockNumber(ctx, d.state, d.etherman, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		block, err := d.state.GetL2BlockByNumber(ctx, blockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("block #%d not found", blockNumber))
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get block by number", err, true)
		}

		traces, rpcErr := d.buildTraceBlock(ctx, block.Transactions(), cfg, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		return traces, nil
	})
}

// TraceBlockByHash creates a response for debug_traceBlockByHash request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtraceblockbyhash
func (d *DebugEndpoints) TraceBlockByHash(hash types.ArgHash, cfg *traceConfig) (interface{}, types.Error) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		block, err := d.state.GetL2BlockByHash(ctx, hash.Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("block %s not found", hash.Hash().String()))
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get block by hash", err, true)
		}

		traces, rpcErr := d.buildTraceBlock(ctx, block.Transactions(), cfg, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		return traces, nil
	})
}

// TraceBatchByNumber creates a response for debug_traceBatchByNumber request.
// this endpoint tries to help clients to get traces at once for all the transactions
// attached to the same batch.
//
// IMPORTANT: in order to take advantage of the infrastructure automatically scaling,
// instead of parallelizing the trace transaction internally and pushing all the load
// to a single jRPC and Executor instance, the code will redirect the trace transaction
// requests to the same url, making them external calls, so we can process in parallel
// with multiple jRPC and Executor instances.
//
// the request flow will work as follows:
// -> user do a trace batch request
// -> jRPC balancer picks a jRPC server to handle the trace batch request
// -> picked jRPC sends parallel trace transaction requests for each transaction in the batch
// -> jRPC balancer sends each request to a different jRPC to handle the trace transaction requests
// -> picked jRPC server group trace transaction responses from other jRPC servers
// -> picked jRPC respond the initial request to the user with all the tx traces
func (d *DebugEndpoints) TraceBatchByNumber(httpRequest *http.Request, number types.BatchNumber, cfg *traceConfig) (interface{}, types.Error) {
	type traceResponse struct {
		blockNumber uint64
		txIndex     uint64
		txHash      common.Hash
		trace       interface{}
		err         error
	}

	// the size of the buffer defines
	// how many txs it will process in parallel.
	const bufferSize = 10

	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		batchNumber, rpcErr := number.GetNumericBatchNumber(ctx, d.state, d.etherman, dbTx)
		if rpcErr != nil {
			return nil, rpcErr
		}

		batch, err := d.state.GetBatchByNumber(ctx, batchNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("batch #%d not found", batchNumber))
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get batch by number", err, true)
		}

		txs, _, err := d.state.GetTransactionsByBatchNumber(ctx, batch.BatchNumber, dbTx)
		if !errors.Is(err, state.ErrNotFound) && err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load batch txs from state by number %v to create the traces", batchNumber), err, true)
		}

		receipts := make([]ethTypes.Receipt, 0, len(txs))
		for _, tx := range txs {
			receipt, err := d.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
			if err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load receipt for tx %v to get trace", tx.Hash().String()), err, true)
			}
			receipts = append(receipts, *receipt)
		}

		requests := make(chan (ethTypes.Receipt), bufferSize)

		mu := &sync.Mutex{}
		wg := sync.WaitGroup{}
		wg.Add(len(receipts))
		responses := make([]traceResponse, 0, len(receipts))

		// gets the trace from the jRPC and adds it to the responses
		loadTraceByTxHash := func(d *DebugEndpoints, receipt ethTypes.Receipt, cfg *traceConfig) {
			response := traceResponse{
				blockNumber: receipt.BlockNumber.Uint64(),
				txIndex:     uint64(receipt.TransactionIndex),
				txHash:      receipt.TxHash,
			}

			defer wg.Done()
			trace, err := d.TraceTransaction(types.ArgHash(receipt.TxHash), cfg)
			if err != nil {
				err := fmt.Errorf("failed to get tx trace for tx %v, err: %w", receipt.TxHash.String(), err)
				log.Errorf(err.Error())
				response.err = err
			} else {
				response.trace = trace
			}

			// add to the responses
			mu.Lock()
			defer mu.Unlock()
			responses = append(responses, response)
		}

		// goes through the buffer and loads the trace
		// by all the transactions added in the buffer
		// then add the results to the responses map
		go func() {
			index := uint(0)
			for req := range requests {
				go loadTraceByTxHash(d, req, cfg)
				index++
			}
		}()

		// add receipts to the buffer
		for _, receipt := range receipts {
			requests <- receipt
		}

		// wait the traces to be loaded
		if waitTimeout(&wg, d.cfg.ReadTimeout.Duration) {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("failed to get traces for batch %v: timeout reached", batchNumber), nil, true)
		}

		close(requests)

		// since the txs are attached to a L2 Block and the L2 Block is
		// the struct attached to the Batch, in order to always respond
		// the traces in the same order, we need to order the transactions
		// first by block number and then by tx index, so we can have something
		// close to the txs being sorted by a tx index related to the batch
		sort.Slice(responses, func(i, j int) bool {
			if responses[i].txIndex != responses[j].txIndex {
				return responses[i].txIndex < responses[j].txIndex
			}
			return responses[i].blockNumber < responses[j].blockNumber
		})

		// build the batch trace response array
		traces := make([]traceBatchTransactionResponse, 0, len(receipts))
		for _, response := range responses {
			if response.err != nil {
				return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("failed to get traces for batch %v: failed to get trace for tx: %v, err: %v", batchNumber, response.txHash.String(), response.err.Error()), nil, true)
			}

			traces = append(traces, traceBatchTransactionResponse{
				TxHash: response.txHash,
				Result: response.trace,
			})
		}
		return traces, nil
	})
}

func (d *DebugEndpoints) buildTraceBlock(ctx context.Context, txs []*ethTypes.Transaction, cfg *traceConfig, dbTx pgx.Tx) (interface{}, types.Error) {
	traces := []traceBlockTransactionResponse{}
	for _, tx := range txs {
		traceTransaction, err := d.buildTraceTransaction(ctx, tx.Hash(), cfg, dbTx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get trace for transaction %v: %v", tx.Hash().String(), err.Error())
			return RPCErrorResponse(types.DefaultErrorCode, errMsg, err, true)
		}
		traceBlockTransaction := traceBlockTransactionResponse{
			Result: traceTransaction,
		}
		traces = append(traces, traceBlockTransaction)
	}

	return traces, nil
}

func (d *DebugEndpoints) buildTraceTransaction(ctx context.Context, hash common.Hash, cfg *traceConfig, dbTx pgx.Tx) (interface{}, types.Error) {
	traceCfg := cfg
	if traceCfg == nil {
		traceCfg = defaultTraceConfig
	}

	// check tracer
	if traceCfg.Tracer != nil && *traceCfg.Tracer != "" && !isBuiltInTracer(*traceCfg.Tracer) && !isJSCustomTracer(*traceCfg.Tracer) && !isSentioTracer(*traceCfg.Tracer) {
		return RPCErrorResponse(types.DefaultErrorCode, "invalid tracer", nil, false)
	}

	stateTraceConfig := state.TraceConfig{
		DisableStack:     traceCfg.DisableStack,
		DisableStorage:   traceCfg.DisableStorage,
		EnableMemory:     traceCfg.EnableMemory,
		EnableReturnData: traceCfg.EnableReturnData,
		Tracer:           traceCfg.Tracer,
		TracerConfig:     traceCfg.TracerConfig,
		Limit:            traceCfg.Limit,
	}
	result, err := d.state.DebugTransaction(ctx, hash, stateTraceConfig, dbTx)
	if errors.Is(err, state.ErrNotFound) {
		return RPCErrorResponse(types.DefaultErrorCode, "transaction not found", nil, false)
	} else if err != nil {
		errorMessage := fmt.Sprintf("failed to get trace: %v", err.Error())
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	// if a tracer was specified, then return the trace result
	if stateTraceConfig.Tracer != nil && *stateTraceConfig.Tracer != "" && len(result.ExecutorTraceResult) > 0 {
		return result.ExecutorTraceResult, nil
	}

	receipt, err := d.state.GetTransactionReceipt(ctx, hash, dbTx)
	if err != nil {
		const errorMessage = "failed to tx receipt"
		log.Errorf("%v: %v", errorMessage, err)
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	failed := receipt.Status == ethTypes.ReceiptStatusFailed
	var returnValue interface{}
	if stateTraceConfig.EnableReturnData {
		returnValue = common.Bytes2Hex(result.ReturnValue)
	}

	structLogs := d.buildStructLogs(result.StructLogs, *traceCfg)

	resp := traceTransactionResponse{
		Gas:         result.GasUsed,
		Failed:      failed,
		ReturnValue: returnValue,
		StructLogs:  structLogs,
	}

	return resp, nil
}

func (d *DebugEndpoints) buildStructLogs(stateStructLogs []instrumentation.StructLog, cfg traceConfig) []StructLogRes {
	structLogs := make([]StructLogRes, 0, len(stateStructLogs))
	memory := fakevm.NewMemory()
	for _, structLog := range stateStructLogs {
		errRes := ""
		if structLog.Err != nil {
			errRes = structLog.Err.Error()
		}

		op := structLog.Op
		if op == "SHA3" {
			op = "KECCAK256"
		} else if op == "STOP" && structLog.Pc == 0 {
			// this stop is generated for calls with single
			// step(no depth increase) and must be ignored
			continue
		}

		structLogRes := StructLogRes{
			Pc:            structLog.Pc,
			Op:            op,
			Gas:           structLog.Gas,
			GasCost:       structLog.GasCost,
			Depth:         structLog.Depth,
			Error:         errRes,
			RefundCounter: structLog.RefundCounter,
		}

		if !cfg.DisableStack {
			stack := make([]types.ArgBig, 0, len(structLog.Stack))
			for _, stackItem := range structLog.Stack {
				if stackItem != nil {
					stack = append(stack, types.ArgBig(*stackItem))
				}
			}
			structLogRes.Stack = &stack
		}

		if cfg.EnableMemory {
			memory.Resize(uint64(structLog.MemorySize))
			if len(structLog.Memory) > 0 {
				memory.Set(uint64(structLog.MemoryOffset), uint64(len(structLog.Memory)), structLog.Memory)
			}

			if structLog.MemorySize > 0 {
				// Populate the structLog memory
				structLog.Memory = memory.Data()

				// Convert memory to string array
				const memoryChunkSize = 32
				memoryArray := make([]string, 0, len(structLog.Memory))

				for i := 0; i < len(structLog.Memory); i = i + memoryChunkSize {
					slice32Bytes := make([]byte, memoryChunkSize)
					copy(slice32Bytes, structLog.Memory[i:i+memoryChunkSize])
					memoryStringItem := hex.EncodeToString(slice32Bytes)
					memoryArray = append(memoryArray, memoryStringItem)
				}

				structLogRes.Memory = &memoryArray
			} else {
				memory = fakevm.NewMemory()
				structLogRes.Memory = &[]string{}
			}
		}

		if !cfg.DisableStorage && len(structLog.Storage) > 0 {
			storage := make(map[string]string, len(structLog.Storage))
			for storageKey, storageValue := range structLog.Storage {
				k := hex.EncodeToString(storageKey.Bytes())
				v := hex.EncodeToString(storageValue.Bytes())
				storage[k] = v
			}
			structLogRes.Storage = &storage
		}

		structLogs = append(structLogs, structLogRes)
	}

	if cfg.Limit > 0 && len(structLogs) > cfg.Limit {
		structLogs = structLogs[:cfg.Limit]
	}
	return structLogs
}

func (d *DebugEndpoints) buildTraceCall(ctx context.Context, blockNumber uint64, sender common.Address, tx *ethTypes.Transaction, cfg *traceConfig, dbTx pgx.Tx) (interface{}, types.Error) {
	traceCfg := cfg
	if traceCfg == nil {
		traceCfg = defaultTraceConfig
	}

	// check tracer
	if traceCfg.Tracer != nil && *traceCfg.Tracer != "" && !isBuiltInTracer(*traceCfg.Tracer) && !isJSCustomTracer(*traceCfg.Tracer) && !isSentioTracer(*traceCfg.Tracer) {
		return RPCErrorResponse(types.DefaultErrorCode, "invalid tracer", nil, false)
	}

	stateTraceConfig := state.TraceConfig{
		DisableStack:     traceCfg.DisableStack,
		DisableStorage:   traceCfg.DisableStorage,
		EnableMemory:     traceCfg.EnableMemory,
		EnableReturnData: traceCfg.EnableReturnData,
		Tracer:           traceCfg.Tracer,
		TracerConfig:     traceCfg.TracerConfig,
	}
	result, err := d.state.TraceCall(ctx, blockNumber, sender, tx, stateTraceConfig, dbTx)
	if errors.Is(err, state.ErrNotFound) {
		return RPCErrorResponse(types.DefaultErrorCode, "transaction not found", nil, false)
	} else if err != nil {
		errorMessage := fmt.Sprintf("failed to get trace: %v", err.Error())
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	// if a tracer was specified, then return the trace result
	if stateTraceConfig.Tracer != nil && *stateTraceConfig.Tracer != "" && len(result.ExecutorTraceResult) > 0 {
		return result.ExecutorTraceResult, nil
	}

	//receipt, err := d.state.GetTransactionReceipt(ctx, tx.Hash(), dbTx)
	//if err != nil {
	//	const errorMessage = "failed to tx receipt"
	//	log.Errorf("%v: %v", errorMessage, err)
	//	return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	//}

	//failed := receipt.Status == ethTypes.ReceiptStatusFailed
	// TODO check is this enough or not
	failed := result.Failed()
	var returnValue interface{}
	if stateTraceConfig.EnableReturnData {
		returnValue = common.Bytes2Hex(result.ReturnValue)
	}

	structLogs := d.buildStructLogs(result.StructLogs, *traceCfg)

	resp := traceTransactionResponse{
		Gas:         result.GasUsed,
		Failed:      failed,
		ReturnValue: returnValue,
		StructLogs:  structLogs,
	}

	return resp, nil
}

// isBuiltInTracer checks if the tracer is one of the
// built-in tracers
func isBuiltInTracer(tracer string) bool {
	// built-in tracers
	switch tracer {
	case "callTracer", "flatCallTracer", "4byteTracer", "prestateTracer", "noopTracer":
		return true
	default:
		return false
	}
}

// isJSCustomTracer checks if the tracer contains the
// functions result and fault which are required for a custom tracer
// https://geth.ethereum.org/docs/developers/evm-tracing/custom-tracer
func isJSCustomTracer(tracer string) bool {
	return strings.Contains(tracer, "result") && strings.Contains(tracer, "fault")
}

// // getHttpScheme tries to get the scheme from the http request in different ways
// func getHttpScheme(r *http.Request) (string, error) {
// 	// scheme headers
// 	headers := []string{"X-Forwarded-Proto", "X-Forwarded-Protocol", "X-Url-Scheme"}
// 	for _, header := range headers {
// 		value := r.Header.Get(header)
// 		if value == "http" || value == "https" {
// 			return value, nil
// 		} else if value != "" {
// 			return "", fmt.Errorf("header %v must be set to HTTP or HTTPS, value found: %s", header, value)
// 		}
// 	}

// 	// https on/off headers
// 	headers = []string{"X-Forwarded-Ssl", "Front-End-Https"}
// 	for _, header := range headers {
// 		value := r.Header.Get(header)
// 		if value == "on" {
// 			return "https", nil
// 		} else if value == "off" {
// 			return "http", nil
// 		} else if value != "" {
// 			return "", fmt.Errorf("header %v must be set to ON or OFF, value found: %s", header, value)
// 		}
// 	}

//		// httpRequest TLS check
//		scheme := "http"
//		if r.TLS != nil {
//			scheme = "https"
//		}
//		return scheme, nil
//	}
func isSentioTracer(tracer string) bool {
	switch tracer {
	case "sentioTracer", "sentioPrestateTracer":
		return true
	default:
		return false
	}
}

// waitTimeout waits for the waitGroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

type BlockOverrides struct {
	BlockNumber *hexutil.Uint64         `json:"blockNumber"`
	Coinbase    *common.Address         `json:"coinbase"`
	Timestamp   *hexutil.Uint64         `json:"timestamp"`
	GasLimit    *hexutil.Uint           `json:"gasLimit"`
	Difficulty  *hexutil.Uint           `json:"difficulty"`
	BaseFee     *uint256.Int            `json:"baseFee"`
	BlockHash   *map[uint64]common.Hash `json:"blockHash"`
}

type Bundle struct {
	Transactions  []*types.TxArgs `json:"transactions"`
	BlockOverride BlockOverrides  `json:"blockOverride"`
}

type StateContext struct {
	BlockNumber      *types.BlockNumberOrHash `json:"blockNumber"`
	TransactionIndex int                      `json:"transactionIndex"`
}

type Account struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]uint256.Int `json:"state"`
	StateDiff *map[common.Hash]uint256.Int `json:"stateDiff"`
}

type StateOverrides map[common.Address]Account

// Copoed from endpoints_eth, TODO reuse old code
func (e *DebugEndpoints) getBlockByArg(ctx context.Context, blockArg *types.BlockNumberOrHash, dbTx pgx.Tx) (*ethTypes.Block, types.Error) {
	// If no block argument is provided, return the latest block
	if blockArg == nil {
		block, err := e.state.GetLastL2Block(ctx, dbTx)
		if err != nil {
			return nil, types.NewRPCError(types.DefaultErrorCode, "failed to get the last block number from state")
		}
		return block, nil
	}

	// If we have a block hash, try to get the block by hash
	if blockArg.IsHash() {
		block, err := e.state.GetL2BlockByHash(ctx, blockArg.Hash().Hash(), dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, types.NewRPCError(types.DefaultErrorCode, "header for hash not found")
		} else if err != nil {
			return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("failed to get block by hash %v", blockArg.Hash().Hash()))
		}
		return block, nil
	}

	// Otherwise, try to get the block by number
	blockNum, rpcErr := blockArg.Number().GetNumericBlockNumber(ctx, e.state, e.etherman, dbTx)
	if rpcErr != nil {
		return nil, rpcErr
	}
	block, err := e.state.GetL2BlockByNumber(context.Background(), blockNum, dbTx)
	if errors.Is(err, state.ErrNotFound) || block == nil {
		return nil, types.NewRPCError(types.DefaultErrorCode, "header not found")
	} else if err != nil {
		return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("failed to get block by number %v", blockNum))
	}

	return block, nil
}

func (d *DebugEndpoints) TraceCallMany(bundles []*Bundle, simulateContext *StateContext, cfg *traceConfig) (interface{}, types.Error) {
	return d.txMan.NewDbTxScope(d.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		block, err := d.getBlockByArg(ctx, simulateContext.BlockNumber, dbTx)
		if errors.Is(err, state.ErrNotFound) {
			return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("block %d %s not found", simulateContext.BlockNumber.Number(), simulateContext.BlockNumber.Hash()))
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get block by hash", err, false)
		}
		var traceCallManyResponse []interface{}

		for _, bundle := range bundles {
			var txs []*ethTypes.Transaction

			var sender common.Address
			for _, txArg := range bundle.Transactions {
				defaultSenderAddress := common.HexToAddress(DefaultSenderAddress)
				senderAddress, tx, err := txArg.ToTransaction(ctx, d.state, d.cfg.MaxCumulativeGasUsed, block.Root(), defaultSenderAddress, dbTx)
				sender = senderAddress
				if err != nil {
					return RPCErrorResponse(types.DefaultErrorCode, "failed to convert arguments into an unsigned transaction", err, false)
				}
				txs = append(txs, tx)
			}

			traces, rpcErr := d.buildTraceCallMany(ctx, block.NumberU64(), sender, txs, cfg, dbTx)
			if rpcErr != nil {
				return nil, rpcErr
			}
			traceCallManyResponse = append(traceCallManyResponse, traces)
		}

		return traceCallManyResponse, nil
	})
}

// TODO be able to override `BlockOverrides`, state overrides, etc
func (d *DebugEndpoints) buildTraceCallMany(ctx context.Context, blockNumber uint64, sender common.Address, txs []*ethTypes.Transaction, cfg *traceConfig, dbTx pgx.Tx) (interface{}, types.Error) {
	var traces []interface{}
	for _, tx := range txs {
		traceTransaction, err := d.buildTraceCall(ctx, blockNumber, sender, tx, cfg, dbTx)
		if err != nil {
			errMsg := fmt.Sprintf("failed to get trace for trace call %v, %v", tx.Hash().String(), err.Error())
			return RPCErrorResponse(types.DefaultErrorCode, errMsg, err, false)
		}
		traces = append(traces, traceTransaction)
	}
	return traces, nil
}
