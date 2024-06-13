package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

var defaultTraceConfig = &traceConfig{
	DisableStorage:   false,
	DisableStack:     false,
	EnableMemory:     false,
	EnableReturnData: false,
	Tracer:           nil,
}

// DebugEndpoints is the debug jsonrpc endpoint
type DebugEndpoints struct {
	cfg      Config
	state    types.StateInterface
	etherman types.EthermanInterface
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
	ctx := context.Background()
	return d.buildTraceTransaction(ctx, hash.Hash(), cfg, nil)
}

// TraceBlockByNumber creates a response for debug_traceBlockByNumber request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtraceblockbynumber
func (d *DebugEndpoints) TraceBlockByNumber(number types.BlockNumber, cfg *traceConfig) (interface{}, types.Error) {
	ctx := context.Background()
	blockNumber, rpcErr := number.GetNumericBlockNumber(ctx, d.state, d.etherman, nil)
	if rpcErr != nil {
		return nil, rpcErr
	}

	block, err := d.state.GetL2BlockByNumber(ctx, blockNumber, nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("block #%d not found", blockNumber))
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get block by number", err, true)
	}

	traces, rpcErr := d.buildTraceBlock(ctx, block.Transactions(), cfg, nil)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return traces, nil
}

// TraceBlockByHash creates a response for debug_traceBlockByHash request.
// See https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtraceblockbyhash
func (d *DebugEndpoints) TraceBlockByHash(hash types.ArgHash, cfg *traceConfig) (interface{}, types.Error) {
	ctx := context.Background()
	block, err := d.state.GetL2BlockByHash(ctx, hash.Hash(), nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("block %s not found", hash.Hash().String()))
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get block by hash", err, true)
	}

	traces, rpcErr := d.buildTraceBlock(ctx, block.Transactions(), cfg, nil)
	if rpcErr != nil {
		return nil, rpcErr
	}

	return traces, nil
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

	ctx := context.Background()
	batchNumber, rpcErr := number.GetNumericBatchNumber(ctx, d.state, d.etherman, nil)
	if rpcErr != nil {
		return nil, rpcErr
	}

	batch, err := d.state.GetBatchByNumber(ctx, batchNumber, nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("batch #%d not found", batchNumber))
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get batch by number", err, true)
	}

	txs, _, err := d.state.GetTransactionsByBatchNumber(ctx, batch.BatchNumber, nil)
	if !errors.Is(err, state.ErrNotFound) && err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load batch txs from state by number %v to create the traces", batchNumber), err, true)
	}

	receipts := make([]ethTypes.Receipt, 0, len(txs))
	for _, tx := range txs {
		receipt, err := d.state.GetTransactionReceipt(ctx, tx.Hash(), nil)
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

	stateTraceConfig := state.TraceConfig{
		DisableStack:     traceCfg.DisableStack,
		DisableStorage:   traceCfg.DisableStorage,
		EnableMemory:     traceCfg.EnableMemory,
		EnableReturnData: traceCfg.EnableReturnData,
		Tracer:           traceCfg.Tracer,
		TracerConfig:     traceCfg.TracerConfig,
	}
	result, err := d.state.DebugTransaction(ctx, hash, stateTraceConfig, dbTx)
	if errors.Is(err, state.ErrNotFound) {
		return RPCErrorResponse(types.DefaultErrorCode, "transaction not found", nil, false)
	} else if err != nil {
		errorMessage := fmt.Sprintf("failed to get trace: %v", err.Error())
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	return result.TraceResult, nil
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
