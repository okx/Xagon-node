package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// ZKEVMEndpoints contains implementations for the "zkevm" RPC endpoints
type ZKEVMEndpoints struct {
	cfg      Config
	pool     types.PoolInterface
	state    types.StateInterface
	etherman types.EthermanInterface
}

// NewZKEVMEndpoints returns ZKEVMEndpoints
func NewZKEVMEndpoints(cfg Config, pool types.PoolInterface, state types.StateInterface, etherman types.EthermanInterface) *ZKEVMEndpoints {
	return &ZKEVMEndpoints{
		cfg:      cfg,
		pool:     pool,
		state:    state,
		etherman: etherman,
	}
}

// ConsolidatedBlockNumber returns last block number related to the last verified batch
func (z *ZKEVMEndpoints) ConsolidatedBlockNumber() (interface{}, types.Error) {
	ctx := context.Background()
	lastBlockNumber, err := z.state.GetLastConsolidatedL2BlockNumber(ctx, nil)
	if err != nil {
		const errorMessage = "failed to get last consolidated block number from state"
		log.Errorf("%v:%v", errorMessage, err)
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	return hex.EncodeUint64(lastBlockNumber), nil
}

// IsBlockConsolidated returns the consolidation status of a provided block number
func (z *ZKEVMEndpoints) IsBlockConsolidated(blockNumber types.ArgUint64) (interface{}, types.Error) {
	ctx := context.Background()
	IsL2BlockConsolidated, err := z.state.IsL2BlockConsolidated(ctx, uint64(blockNumber), nil)
	if err != nil {
		const errorMessage = "failed to check if the block is consolidated"
		log.Errorf("%v: %v", errorMessage, err)
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	return IsL2BlockConsolidated, nil
}

// IsBlockVirtualized returns the virtualization status of a provided block number
func (z *ZKEVMEndpoints) IsBlockVirtualized(blockNumber types.ArgUint64) (interface{}, types.Error) {
	ctx := context.Background()
	IsL2BlockVirtualized, err := z.state.IsL2BlockVirtualized(ctx, uint64(blockNumber), nil)
	if err != nil {
		const errorMessage = "failed to check if the block is virtualized"
		log.Errorf("%v: %v", errorMessage, err)
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	return IsL2BlockVirtualized, nil
}

// BatchNumberByBlockNumber returns the batch number from which the passed block number is created
func (z *ZKEVMEndpoints) BatchNumberByBlockNumber(blockNumber types.ArgUint64) (interface{}, types.Error) {
	ctx := context.Background()
	batchNum, err := z.state.BatchNumberByL2BlockNumber(ctx, uint64(blockNumber), nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		const errorMessage = "failed to get batch number from block number"
		log.Errorf("%v: %v", errorMessage, err.Error())
		return nil, types.NewRPCError(types.DefaultErrorCode, errorMessage)
	}

	return hex.EncodeUint64(batchNum), nil
}

// BatchNumber returns the latest trusted batch number
func (z *ZKEVMEndpoints) BatchNumber() (interface{}, types.Error) {
	ctx := context.Background()
	lastBatchNumber, err := z.state.GetLastBatchNumber(ctx, nil)
	if err != nil {
		return "0x0", types.NewRPCError(types.DefaultErrorCode, "failed to get the last batch number from state")
	}

	return hex.EncodeUint64(lastBatchNumber), nil
}

// VirtualBatchNumber returns the latest virtualized batch number
func (z *ZKEVMEndpoints) VirtualBatchNumber() (interface{}, types.Error) {
	ctx := context.Background()
	lastBatchNumber, err := z.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil {
		return "0x0", types.NewRPCError(types.DefaultErrorCode, "failed to get the last virtual batch number from state")
	}

	return hex.EncodeUint64(lastBatchNumber), nil
}

// VerifiedBatchNumber returns the latest verified batch number
func (z *ZKEVMEndpoints) VerifiedBatchNumber() (interface{}, types.Error) {
	ctx := context.Background()
	lastBatch, err := z.state.GetLastVerifiedBatch(ctx, nil)
	if err != nil {
		return "0x0", types.NewRPCError(types.DefaultErrorCode, "failed to get the last verified batch number from state")
	}
	return hex.EncodeUint64(lastBatch.BatchNumber), nil
}

// GetBatchByNumber returns information about a batch by batch number
func (z *ZKEVMEndpoints) GetBatchByNumber(batchNumber types.BatchNumber, fullTx bool) (interface{}, types.Error) {
	ctx := context.Background()
	var err error
	numericBatchNumber, rpcErr := batchNumber.GetNumericBatchNumber(ctx, z.state, z.etherman, nil)
	if rpcErr != nil {
		return nil, rpcErr
	}

	batch, err := z.state.GetBatchByNumber(ctx, numericBatchNumber, nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load batch from state by number %v", numericBatchNumber), err, true)
	}
	batchTimestamp, err := z.state.GetBatchTimestamp(ctx, numericBatchNumber, nil, nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load batch timestamp from state by number %v", numericBatchNumber), err, true)
	}

	if batchTimestamp == nil {
		batch.Timestamp = time.Time{}
	} else {
		batch.Timestamp = *batchTimestamp
	}

	txs, _, err := z.state.GetTransactionsByBatchNumber(ctx, numericBatchNumber, nil)
	if !errors.Is(err, state.ErrNotFound) && err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load batch txs from state by number %v", numericBatchNumber), err, true)
	}

	receipts := make([]ethTypes.Receipt, 0, len(txs))
	for _, tx := range txs {
		receipt, err := z.state.GetTransactionReceipt(ctx, tx.Hash(), nil)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load receipt for tx %v", tx.Hash().String()), err, true)
		}
		receipts = append(receipts, *receipt)
	}

	virtualBatch, err := z.state.GetVirtualBatch(ctx, numericBatchNumber, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load virtual batch from state by number %v", numericBatchNumber), err, true)
	}

	verifiedBatch, err := z.state.GetVerifiedBatch(ctx, numericBatchNumber, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load virtual batch from state by number %v", numericBatchNumber), err, true)
	}

	ger, err := z.state.GetExitRootByGlobalExitRoot(ctx, batch.GlobalExitRoot, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load full GER from state by number %v", numericBatchNumber), err, true)
	} else if errors.Is(err, state.ErrNotFound) {
		ger = &state.GlobalExitRoot{}
	}

	blocks, err := z.state.GetL2BlocksByBatchNumber(ctx, numericBatchNumber, nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load blocks associated to the batch %v", numericBatchNumber), err, true)
	}

	batch.Transactions = txs
	rpcBatch, err := types.NewBatch(ctx, z.state, batch, virtualBatch, verifiedBatch, blocks, receipts, fullTx, true, ger, nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't build the batch %v response", numericBatchNumber), err, true)
	}
	return rpcBatch, nil
}

// GetFullBlockByNumber returns information about a block by block number
func (z *ZKEVMEndpoints) GetFullBlockByNumber(number types.BlockNumber, fullTx bool) (interface{}, types.Error) {
	ctx := context.Background()
	if number == types.PendingBlockNumber {
		lastBlock, err := z.state.GetLastL2Block(ctx, nil)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "couldn't load last block from state to compute the pending block", err, true)
		}
		l2Header := state.NewL2Header(&ethTypes.Header{
			ParentHash: lastBlock.Hash(),
			Number:     big.NewInt(0).SetUint64(lastBlock.Number().Uint64() + 1),
			TxHash:     ethTypes.EmptyRootHash,
			UncleHash:  ethTypes.EmptyUncleHash,
		})
		l2Block := state.NewL2BlockWithHeader(l2Header)
		rpcBlock, err := types.NewBlock(ctx, z.state, nil, l2Block, nil, fullTx, false, state.Ptr(true), nil)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "couldn't build the pending block response", err, true)
		}

		// clean fields that are not available for pending block
		rpcBlock.Hash = nil
		rpcBlock.Miner = nil
		rpcBlock.Nonce = nil
		rpcBlock.TotalDifficulty = nil

		return rpcBlock, nil
	}
	var err error
	blockNumber, rpcErr := number.GetNumericBlockNumber(ctx, z.state, z.etherman, nil)
	if rpcErr != nil {
		return nil, rpcErr
	}

	l2Block, err := z.state.GetL2BlockByNumber(ctx, blockNumber, nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load block from state by number %v", blockNumber), err, true)
	}

	txs := l2Block.Transactions()
	receipts := make([]ethTypes.Receipt, 0, len(txs))
	for _, tx := range txs {
		receipt, err := z.state.GetTransactionReceipt(ctx, tx.Hash(), nil)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load receipt for tx %v", tx.Hash().String()), err, true)
		}
		receipts = append(receipts, *receipt)
	}

	rpcBlock, err := types.NewBlock(ctx, z.state, state.Ptr(l2Block.Hash()), l2Block, receipts, fullTx, true, state.Ptr(true), nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't build block response for block by number %v", blockNumber), err, true)
	}

	return rpcBlock, nil
}

// GetFullBlockByHash returns information about a block by hash
func (z *ZKEVMEndpoints) GetFullBlockByHash(hash types.ArgHash, fullTx bool) (interface{}, types.Error) {
	ctx := context.Background()
	l2Block, err := z.state.GetL2BlockByHash(ctx, hash.Hash(), nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get block by hash from state", err, true)
	}

	txs := l2Block.Transactions()
	receipts := make([]ethTypes.Receipt, 0, len(txs))
	for _, tx := range txs {
		receipt, err := z.state.GetTransactionReceipt(ctx, tx.Hash(), nil)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't load receipt for tx %v", tx.Hash().String()), err, true)
		}
		receipts = append(receipts, *receipt)
	}

	rpcBlock, err := types.NewBlock(ctx, z.state, state.Ptr(l2Block.Hash()), l2Block, receipts, fullTx, true, state.Ptr(true), nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, fmt.Sprintf("couldn't build block response for block by hash %v", hash.Hash()), err, true)
	}

	return rpcBlock, nil
}

// GetNativeBlockHashesInRange return the state root for the blocks in range
func (z *ZKEVMEndpoints) GetNativeBlockHashesInRange(filter NativeBlockHashBlockRangeFilter) (interface{}, types.Error) {
	ctx := context.Background()
	fromBlockNumber, toBlockNumber, rpcErr := filter.GetNumericBlockNumbers(ctx, z.cfg, z.state, z.etherman, nil)
	if rpcErr != nil {
		return nil, rpcErr
	}

	nativeBlockHashes, err := z.state.GetNativeBlockHashesInRange(ctx, fromBlockNumber, toBlockNumber, nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if errors.Is(err, state.ErrMaxNativeBlockHashBlockRangeLimitExceeded) {
		errMsg := fmt.Sprintf(state.ErrMaxNativeBlockHashBlockRangeLimitExceeded.Error(), z.cfg.MaxNativeBlockHashBlockRange)
		return RPCErrorResponse(types.InvalidParamsErrorCode, errMsg, nil, false)
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get block by hash from state", err, true)
	}

	return nativeBlockHashes, nil
}

// GetTransactionByL2Hash returns a transaction by his l2 hash
func (z *ZKEVMEndpoints) GetTransactionByL2Hash(hash types.ArgHash) (interface{}, types.Error) {
	ctx := context.Background()
	// try to get tx from state
	tx, err := z.state.GetTransactionByL2Hash(ctx, hash.Hash(), nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to load transaction by l2 hash from state", err, true)
	}
	if tx != nil {
		receipt, err := z.state.GetTransactionReceipt(ctx, hash.Hash(), nil)
		if errors.Is(err, state.ErrNotFound) {
			return RPCErrorResponse(types.DefaultErrorCode, "transaction receipt not found", err, false)
		} else if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to load transaction receipt from state", err, true)
		}

		l2Hash, err := z.state.GetL2TxHashByTxHash(ctx, tx.Hash(), nil)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to get l2 transaction hash", err, true)
		}

		res, err := types.NewTransaction(*tx, receipt, false, l2Hash)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to build transaction response", err, true)
		}

		return res, nil
	}

	// if the tx does not exist in the state, look for it in the pool
	if z.cfg.SequencerNodeURI != "" {
		return z.getTransactionByL2HashFromSequencerNode(hash.Hash())
	}
	poolTx, err := z.pool.GetTransactionByL2Hash(ctx, hash.Hash())
	if errors.Is(err, pool.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to load transaction by l2 hash from pool", err, true)
	}
	if poolTx.Status == pool.TxStatusPending {
		tx = &poolTx.Transaction
		res, err := types.NewTransaction(*tx, nil, false, nil)
		if err != nil {
			return RPCErrorResponse(types.DefaultErrorCode, "failed to build transaction response", err, true)
		}
		return res, nil
	}
	return nil, nil
}

// GetTransactionReceiptByL2Hash returns a transaction receipt by his hash
func (z *ZKEVMEndpoints) GetTransactionReceiptByL2Hash(hash types.ArgHash) (interface{}, types.Error) {
	ctx := context.Background()
	tx, err := z.state.GetTransactionByL2Hash(ctx, hash.Hash(), nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get tx from state", err, true)
	}

	r, err := z.state.GetTransactionReceipt(ctx, hash.Hash(), nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get tx receipt from state", err, true)
	}

	l2Hash, err := z.state.GetL2TxHashByTxHash(ctx, tx.Hash(), nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get l2 transaction hash", err, true)
	}

	receipt, err := types.NewReceipt(*tx, r, l2Hash)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to build the receipt response", err, true)
	}

	return receipt, nil
}

func (z *ZKEVMEndpoints) getTransactionByL2HashFromSequencerNode(hash common.Hash) (interface{}, types.Error) {
	res, err := client.JSONRPCCall(z.cfg.SequencerNodeURI, "zkevm_getTransactionByL2Hash", hash.String())
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get tx from sequencer node by l2 hash", err, true)
	}

	if res.Error != nil {
		return RPCErrorResponse(res.Error.Code, res.Error.Message, nil, false)
	}

	var tx *types.Transaction
	err = json.Unmarshal(res.Result, &tx)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to read tx loaded by l2 hash from sequencer node", err, true)
	}
	return tx, nil
}

// GetExitRootsByGER returns the exit roots accordingly to the provided Global Exit Root
func (z *ZKEVMEndpoints) GetExitRootsByGER(globalExitRoot common.Hash) (interface{}, types.Error) {
	ctx := context.Background()
	exitRoots, err := z.state.GetExitRootByGlobalExitRoot(ctx, globalExitRoot, nil)
	if errors.Is(err, state.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to get exit roots by global exit root from state", err, true)
	}

	return types.ExitRoots{
		BlockNumber:     types.ArgUint64(exitRoots.BlockNumber),
		Timestamp:       types.ArgUint64(exitRoots.Timestamp.Unix()),
		MainnetExitRoot: exitRoots.MainnetExitRoot,
		RollupExitRoot:  exitRoots.RollupExitRoot,
	}, nil
}

// EstimateGasPrice returns an estimate gas price for the transaction.
func (z *ZKEVMEndpoints) EstimateGasPrice(arg *types.TxArgs, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	ctx := context.Background()
	gasPrice, _, err := z.internalEstimateGasPriceAndFee(ctx, arg, blockArg, nil)
	if err != nil {
		return nil, err
	}
	return hex.EncodeBig(gasPrice), nil
}

// EstimateFee returns an estimate fee for the transaction.
func (z *ZKEVMEndpoints) EstimateFee(arg *types.TxArgs, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	ctx := context.Background()
	_, fee, err := z.internalEstimateGasPriceAndFee(ctx, arg, blockArg, nil)
	if err != nil {
		return nil, err
	}
	return hex.EncodeBig(fee), nil
}

// internalEstimateGasPriceAndFee computes the estimated gas price and the estimated fee for the transaction
func (z *ZKEVMEndpoints) internalEstimateGasPriceAndFee(ctx context.Context, arg *types.TxArgs, blockArg *types.BlockNumberOrHash, dbTx pgx.Tx) (*big.Int, *big.Int, types.Error) {
	if arg == nil {
		return nil, nil, types.NewRPCError(types.InvalidParamsErrorCode, "missing value for required argument 0")
	}

	block, respErr := z.getBlockByArg(ctx, blockArg, nil)
	if respErr != nil {
		return nil, nil, respErr
	}

	var blockToProcess *uint64
	if blockArg != nil {
		blockNumArg := blockArg.Number()
		if blockNumArg != nil && (*blockArg.Number() == types.LatestBlockNumber || *blockArg.Number() == types.PendingBlockNumber) {
			blockToProcess = nil
		} else {
			n := block.NumberU64()
			blockToProcess = &n
		}
	}

	defaultSenderAddress := common.HexToAddress(state.DefaultSenderAddress)
	sender, tx, err := arg.ToTransaction(ctx, z.state, z.cfg.MaxCumulativeGasUsed, block.Root(), defaultSenderAddress, nil)
	if err != nil {
		return nil, nil, types.NewRPCError(types.DefaultErrorCode, "failed to convert arguments into an unsigned transaction")
	}

	gasEstimation, returnValue, err := z.state.EstimateGas(tx, sender, blockToProcess, nil)
	if errors.Is(err, runtime.ErrExecutionReverted) {
		data := make([]byte, len(returnValue))
		copy(data, returnValue)
		return nil, nil, types.NewRPCErrorWithData(types.RevertedErrorCode, err.Error(), data)
	} else if err != nil {
		errMsg := fmt.Sprintf("failed to estimate gas: %v", err.Error())
		return nil, nil, types.NewRPCError(types.DefaultErrorCode, errMsg)
	}

	gasPrices, err := z.pool.GetGasPrices(ctx)
	if err != nil {
		return nil, nil, types.NewRPCError(types.DefaultErrorCode, "failed to get L2 gas price", err, false)
	}

	txGasPrice := new(big.Int).SetUint64(gasPrices.L2GasPrice) // by default we assume the tx gas price is the current L2 gas price
	txEGPPct := state.MaxEffectivePercentage
	egpEnabled := z.pool.EffectiveGasPriceEnabled()

	if egpEnabled {
		rawTx, err := state.EncodeTransactionWithoutEffectivePercentage(*tx)
		if err != nil {
			return nil, nil, types.NewRPCError(types.DefaultErrorCode, "failed to encode tx", err, false)
		}

		txEGP, err := z.pool.CalculateEffectiveGasPrice(rawTx, txGasPrice, gasEstimation, gasPrices.L1GasPrice, gasPrices.L2GasPrice)
		if err != nil {
			return nil, nil, types.NewRPCError(types.DefaultErrorCode, "failed to calculate effective gas price", err, false)
		}

		if txEGP.Cmp(txGasPrice) == -1 { // txEGP < txGasPrice
			// We need to "round" the final effectiveGasPrice to a 256 fraction of the txGasPrice
			txEGPPct, err = z.pool.CalculateEffectiveGasPricePercentage(txGasPrice, txEGP)
			if err != nil {
				return nil, nil, types.NewRPCError(types.DefaultErrorCode, "failed to calculate effective gas price percentage", err, false)
			}
			// txGasPriceFraction = txGasPrice/256
			txGasPriceFraction := new(big.Int).Div(txGasPrice, new(big.Int).SetUint64(256)) //nolint:gomnd
			// txGasPrice = txGasPriceFraction*(txEGPPct+1)
			txGasPrice = new(big.Int).Mul(txGasPriceFraction, new(big.Int).SetUint64(uint64(txEGPPct+1)))
		}

		log.Infof("[internalEstimateGasPriceAndFee] finalGasPrice: %d, effectiveGasPrice: %d, egpPct: %d, l2GasPrice: %d, len: %d, gas: %d, l1GasPrice: %d",
			txGasPrice, txEGP, txEGPPct, gasPrices.L2GasPrice, len(rawTx), gasEstimation, gasPrices.L1GasPrice)
	}

	fee := new(big.Int).Mul(txGasPrice, new(big.Int).SetUint64(gasEstimation))

	log.Infof("[internalEstimateGasPriceAndFee] egpEnabled: %t, fee: %d, gasPrice: %d, gas: %d", egpEnabled, fee, txGasPrice, gasEstimation)

	return txGasPrice, fee, nil
}

// EstimateCounters returns an estimation of the counters that are going to be used while executing
// this transaction.
func (z *ZKEVMEndpoints) EstimateCounters(arg *types.TxArgs, blockArg *types.BlockNumberOrHash) (interface{}, types.Error) {
	ctx := context.Background()
	if arg == nil {
		return RPCErrorResponse(types.InvalidParamsErrorCode, "missing value for required argument 0", nil, false)
	}

	block, respErr := z.getBlockByArg(ctx, blockArg, nil)
	if respErr != nil {
		return nil, respErr
	}

	var blockToProcess *uint64
	if blockArg != nil {
		blockNumArg := blockArg.Number()
		if blockNumArg != nil && (*blockArg.Number() == types.LatestBlockNumber || *blockArg.Number() == types.PendingBlockNumber) {
			blockToProcess = nil
		} else {
			n := block.NumberU64()
			blockToProcess = &n
		}
	}

	defaultSenderAddress := common.HexToAddress(state.DefaultSenderAddress)
	sender, tx, err := arg.ToTransaction(ctx, z.state, z.cfg.MaxCumulativeGasUsed, block.Root(), defaultSenderAddress, nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "failed to convert arguments into an unsigned transaction", err, false)
	}

	var oocErr error
	processBatchResponse, err := z.state.PreProcessUnsignedTransaction(ctx, tx, sender, blockToProcess, nil)
	if err != nil {
		if executor.IsROMOutOfCountersError(executor.RomErrorCode(err)) {
			oocErr = err
		} else {
			errMsg := fmt.Sprintf("failed to estimate counters: %v", err.Error())
			return nil, types.NewRPCError(types.DefaultErrorCode, errMsg)
		}
	}

	var revert *types.RevertInfo
	if len(processBatchResponse.BlockResponses) > 0 && len(processBatchResponse.BlockResponses[0].TransactionResponses) > 0 {
		txResponse := processBatchResponse.BlockResponses[0].TransactionResponses[0]
		err = txResponse.RomError
		if errors.Is(err, runtime.ErrExecutionReverted) {
			returnValue := make([]byte, len(txResponse.ReturnValue))
			copy(returnValue, txResponse.ReturnValue)
			err := state.ConstructErrorFromRevert(err, returnValue)
			revert = &types.RevertInfo{
				Message: err.Error(),
				Data:    state.Ptr(types.ArgBytes(returnValue)),
			}
		}
	}

	limits := types.ZKCountersLimits{
		MaxGasUsed:          types.ArgUint64(state.MaxTxGasLimit),
		MaxKeccakHashes:     types.ArgUint64(z.cfg.ZKCountersLimits.MaxKeccakHashes),
		MaxPoseidonHashes:   types.ArgUint64(z.cfg.ZKCountersLimits.MaxPoseidonHashes),
		MaxPoseidonPaddings: types.ArgUint64(z.cfg.ZKCountersLimits.MaxPoseidonPaddings),
		MaxMemAligns:        types.ArgUint64(z.cfg.ZKCountersLimits.MaxMemAligns),
		MaxArithmetics:      types.ArgUint64(z.cfg.ZKCountersLimits.MaxArithmetics),
		MaxBinaries:         types.ArgUint64(z.cfg.ZKCountersLimits.MaxBinaries),
		MaxSteps:            types.ArgUint64(z.cfg.ZKCountersLimits.MaxSteps),
		MaxSHA256Hashes:     types.ArgUint64(z.cfg.ZKCountersLimits.MaxSHA256Hashes),
	}
	return types.NewZKCountersResponse(processBatchResponse.UsedZkCounters, limits, revert, oocErr), nil
}

func (z *ZKEVMEndpoints) getBlockByArg(ctx context.Context, blockArg *types.BlockNumberOrHash, dbTx pgx.Tx) (*state.L2Block, types.Error) {
	// If no block argument is provided, return the latest block
	if blockArg == nil {
		block, err := z.state.GetLastL2Block(ctx, nil)
		if err != nil {
			return nil, types.NewRPCError(types.DefaultErrorCode, "failed to get the last block number from state")
		}
		return block, nil
	}

	// If we have a block hash, try to get the block by hash
	if blockArg.IsHash() {
		block, err := z.state.GetL2BlockByHash(ctx, blockArg.Hash().Hash(), nil)
		if errors.Is(err, state.ErrNotFound) {
			return nil, types.NewRPCError(types.DefaultErrorCode, "header for hash not found")
		} else if err != nil {
			return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("failed to get block by hash %v", blockArg.Hash().Hash()))
		}
		return block, nil
	}

	// Otherwise, try to get the block by number
	blockNum, rpcErr := blockArg.Number().GetNumericBlockNumber(ctx, z.state, z.etherman, nil)
	if rpcErr != nil {
		return nil, rpcErr
	}
	block, err := z.state.GetL2BlockByNumber(context.Background(), blockNum, nil)
	if errors.Is(err, state.ErrNotFound) || block == nil {
		return nil, types.NewRPCError(types.DefaultErrorCode, "header not found")
	} else if err != nil {
		return nil, types.NewRPCError(types.DefaultErrorCode, fmt.Sprintf("failed to get block by number %v", blockNum))
	}

	return block, nil
}

// GetLatestGlobalExitRoot returns the last global exit root used by l2
func (z *ZKEVMEndpoints) GetLatestGlobalExitRoot() (interface{}, types.Error) {
	ctx := context.Background()
	var err error

	ger, err := z.state.GetLatestBatchGlobalExitRoot(ctx, nil)
	if err != nil {
		return RPCErrorResponse(types.DefaultErrorCode, "couldn't load the last global exit root", err, true)
	}

	return ger.String(), nil
}
