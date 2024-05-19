package sequencer

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//TODO: Fix tests ETROG
/*
const (
	forkId5 uint64 = 5
)
*/

var (
	f            *finalizer
	ctx          context.Context
	err          error
	nilErr       error
	poolMock     = new(PoolMock)
	stateMock    = new(StateMock)
	ethermanMock = new(EthermanMock)
	workerMock   = new(WorkerMock)
	dbTxMock     = new(DbTxMock)
	bc           = state.BatchConstraintsCfg{
		MaxTxsPerBatch:       300,
		MaxBatchBytesSize:    120000,
		MaxCumulativeGasUsed: 30000000,
		MaxKeccakHashes:      2145,
		MaxPoseidonHashes:    252357,
		MaxPoseidonPaddings:  135191,
		MaxMemAligns:         236585,
		MaxArithmetics:       236585,
		MaxBinaries:          473170,
		MaxSteps:             7570538,
		MaxSHA256Hashes:      1596,
	}
	cfg = FinalizerCfg{
		ForcedBatchesTimeout: cfgTypes.Duration{
			Duration: 60,
		},
		NewTxsWaitInterval: cfgTypes.Duration{
			Duration: 60,
		},
		ForcedBatchesCheckInterval: cfgTypes.Duration{
			Duration: 10 * time.Second,
		},
		ResourceExhaustedMarginPct: 10,
		SequentialBatchSanityCheck: true,
	}
	poolCfg = pool.Config{
		EffectiveGasPrice: pool.EffectiveGasPriceCfg{
			Enabled:                     false,
			L1GasPriceFactor:            0.25,
			ByteGasCost:                 16,
			ZeroByteGasCost:             4,
			NetProfit:                   1.0,
			BreakEvenFactor:             1.1,
			FinalDeviationPct:           10,
			EthTransferGasPrice:         0,
			EthTransferL1GasPriceFactor: 0,
			L2GasPriceSuggesterFactor:   0.5,
		},
		DefaultMinGasPriceAllowed: 1000000000,
	}
	// chainID         = new(big.Int).SetInt64(400)
	// pvtKey          = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	nonce1  = uint64(1)
	nonce2  = uint64(2)
	seqAddr = common.Address{}
	oldHash = common.HexToHash("0x01")
	newHash = common.HexToHash("0x02")
	// newHash2 = common.HexToHash("0x03")
	// stateRootHashes = []common.Hash{oldHash, newHash, newHash2}
	// txHash       = common.HexToHash("0xf9e4fe4bd2256f782c66cffd76acdb455a76111842bb7e999af2f1b7f4d8d092")
	// txHash2      = common.HexToHash("0xb281831a3401a04f3afa4ec586ef874f58c61b093643d408ea6aa179903df1a4")
	senderAddr   = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	receiverAddr = common.HexToAddress("0x1555324")
	isSynced     = func(ctx context.Context) bool {
		return true
	}
	testErrStr = "some err"
	// testErr                 = fmt.Errorf(testErrStr)
	// openBatchError          = fmt.Errorf("failed to open new batch, err: %v", testErr)
	// cumulativeGasErr        = state.GetZKCounterError("CumulativeGasUsed")
	testBatchL2DataAsString = "0xee80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e980801186622d03b6b8da7cf111d1ccba5bb185c56deae6a322cebc6dda0556f3cb9700910c26408b64b51c5da36ba2f38ef55ba1cee719d5a6c012259687999074321bff"
	decodedBatchL2Data      []byte
	// done                    chan bool
	// gasPrice                = big.NewInt(1000000)
	// effectiveGasPrice       = big.NewInt(1000000)
	// l1GasPrice              = uint64(1000000)
)

func testNow() time.Time {
	return time.Unix(0, 0)
}

func TestNewFinalizer(t *testing.T) {
	eventStorage, err := nileventstorage.NewNilEventStorage()
	require.NoError(t, err)
	eventLog := event.NewEventLog(event.Config{}, eventStorage)

	poolMock.On("GetLastSentFlushID", context.Background()).Return(uint64(0), nil)

	// arrange and act
	f = newFinalizer(cfg, poolCfg, workerMock, poolMock, stateMock, ethermanMock, seqAddr, isSynced, bc, eventLog, nil, newTimeoutCond(&sync.Mutex{}), nil)

	// assert
	assert.NotNil(t, f)
	assert.Equal(t, f.cfg, cfg)
	assert.Equal(t, f.workerIntf, workerMock)
	assert.Equal(t, poolMock, poolMock)
	assert.Equal(t, f.stateIntf, stateMock)
	assert.Equal(t, f.sequencerAddress, seqAddr)
	assert.Equal(t, f.batchConstraints, bc)
}

/*func TestFinalizer_handleProcessTransactionResponse(t *testing.T) {
	   f = setupFinalizer(true)
	   ctx = context.Background()

	   	txTracker := &TxTracker{
	   		Hash:              txHash,
	   		From:              senderAddr,
	   		Nonce:             1,
	   		GasPrice:          gasPrice,
	   		EffectiveGasPrice: effectiveGasPrice,
	   		L1GasPrice:        l1GasPrice,
	   		EGPLog: state.EffectiveGasPriceLog{
	   			ValueFinal:     new(big.Int).SetUint64(0),
	   			ValueFirst:     new(big.Int).SetUint64(0),
	   			ValueSecond:    new(big.Int).SetUint64(0),
	   			FinalDeviation: new(big.Int).SetUint64(0),
	   			MaxDeviation:   new(big.Int).SetUint64(0),
	   			GasPrice:       new(big.Int).SetUint64(0),
	   		},
	   		BatchResources: state.BatchResources{
	   			Bytes: 1000,
	   			ZKCounters: state.ZKCounters{
	   				GasUsed: 500,
	   			},
	   		},
	   		RawTx: []byte{0, 0, 1, 2, 3, 4, 5},
	   	}

	   	txResponse := &state.ProcessTransactionResponse{
	   		TxHash:    txHash,
	   		StateRoot: newHash2,
	   		RomError:  nil,
	   		GasUsed:   100000,
	   	}

	   	blockResponse := &state.ProcessBlockResponse{
	   		TransactionResponses: []*state.ProcessTransactionResponse{
	   			txResponse,
	   		},
	   	}

	   	batchResponse := &state.ProcessBatchResponse{
	   		BlockResponses: []*state.ProcessBlockResponse{
	   			blockResponse,
	   		},
	   	}

	   	txResponseIntrinsicError := &state.ProcessTransactionResponse{
	   		TxHash:    txHash,
	   		StateRoot: newHash2,
	   		RomError:  runtime.ErrIntrinsicInvalidNonce,
	   	}

	   	blockResponseIntrinsicError := &state.ProcessBlockResponse{
	   		TransactionResponses: []*state.ProcessTransactionResponse{
	   			txResponseIntrinsicError,
	   		},
	   	}

	   	txResponseOOCError := &state.ProcessTransactionResponse{
	   		TxHash:    txHash,
	   		StateRoot: newHash2,
	   		RomError:  runtime.ErrOutOfCountersKeccak,
	   	}

	   	blockResponseOOCError := &state.ProcessBlockResponse{
	   		TransactionResponses: []*state.ProcessTransactionResponse{
	   			txResponseOOCError,
	   		},
	   	}

	   	testCases := []struct {
	   		name                       string
	   		executorResponse           *state.ProcessBatchResponse
	   		oldStateRoot               common.Hash
	   		expectedStoredTx           transactionToStore
	   		expectedMoveToNotReadyCall bool
	   		expectedDeleteTxCall       bool
	   		expectedUpdateTxCall       bool
	   		expectedError              error
	   		expectedUpdateTxStatus     pool.TxStatus
	   	}{

	   		{
	   			name: "Successful transaction",
	   			executorResponse: &state.ProcessBatchResponse{
	   				BlockResponses: []*state.ProcessBlockResponse{
	   					blockResponse,
	   				},
	   				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
	   					senderAddr: {
	   						Address: senderAddr,
	   						Nonce:   &nonce2,
	   						Balance: big.NewInt(100),
	   					},
	   					receiverAddr: {
	   						Address: receiverAddr,
	   						Nonce:   nil,
	   						Balance: big.NewInt(100),
	   					},
	   				},
	   			},
	   			oldStateRoot: oldHash,
	   			expectedStoredTx: transactionToStore{
	   				hash:          txHash,
	   				from:          senderAddr,
	   				batchNumber:   f.wipBatch.batchNumber,
	   				coinbase:      f.wipBatch.coinbase,
	   				timestamp:     f.wipBatch.timestamp,
	   				oldStateRoot:  oldHash,
	   				batchResponse: batchResponse,
	   				response:      txResponse,
	   				isForcedBatch: false,
	   			},
	   		},
	   		{
	   			name: "Batch resources underflow err",
	   			executorResponse: &state.ProcessBatchResponse{
	   				UsedZkCounters: state.ZKCounters{
	   					GasUsed: f.wipBatch.remainingResources.ZKCounters.GasUsed + 1,
	   				},
	   				BlockResponses: []*state.ProcessBlockResponse{
	   					blockResponse,
	   				},
	   				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
	   					senderAddr: {
	   						Address: senderAddr,
	   						Nonce:   &nonce1,
	   						Balance: big.NewInt(100),
	   					},
	   				},
	   			},
	   			oldStateRoot:         oldHash,
	   			expectedUpdateTxCall: true,
	   			expectedError:        state.NewBatchRemainingResourcesUnderflowError(cumulativeGasErr, cumulativeGasErr.Error()),
	   		},
	   		{
	   			name: "Intrinsic err",
	   			executorResponse: &state.ProcessBatchResponse{
	   				IsRomOOCError: false,
	   				UsedZkCounters: state.ZKCounters{
	   					GasUsed: 1,
	   				},
	   				BlockResponses: []*state.ProcessBlockResponse{
	   					blockResponseIntrinsicError,
	   				},
	   				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
	   					senderAddr: {
	   						Address: senderAddr,
	   						Nonce:   &nonce1,
	   						Balance: big.NewInt(100),
	   					},
	   				},
	   			},
	   			oldStateRoot:               oldHash,
	   			expectedMoveToNotReadyCall: true,
	   			expectedError:              txResponseIntrinsicError.RomError,
	   		},
	   		{
	   			name: "Out Of Counters err",
	   			executorResponse: &state.ProcessBatchResponse{
	   				IsRomOOCError: true,
	   				UsedZkCounters: state.ZKCounters{
	   					UsedKeccakHashes: bc.MaxKeccakHashes + 1,
	   				},
	   				BlockResponses: []*state.ProcessBlockResponse{
	   					blockResponseOOCError,
	   				},
	   			},
	   			oldStateRoot:           oldHash,
	   			expectedError:          txResponseOOCError.RomError,
	   			expectedDeleteTxCall:   true,
	   			expectedUpdateTxStatus: pool.TxStatusInvalid,
	   		},
	   	}

	   	for _, tc := range testCases {
	   		t.Run(tc.name, func(t *testing.T) {
	   			storedTxs := make([]transactionToStore, 0)
	   			f.pendingL2BlocksToStore = make(chan transactionToStore)

	   			if tc.expectedStoredTx.batchResponse != nil {
	   				done = make(chan bool) // init a new done channel
	   				go func() {
	   					for tx := range f.pendingL2BlocksToStore {
	   						storedTxs = append(storedTxs, tx)
	   						f.pendingL2BlocksToStoreWG.Done()
	   					}
	   					done <- true // signal that the goroutine is done
	   				}()
	   			}
	   			if tc.expectedDeleteTxCall {
	   				workerMock.On("DeleteTx", txTracker.Hash, txTracker.From).Return().Once()
	   			}
	   			if tc.expectedMoveToNotReadyCall {
	   				addressInfo := tc.executorResponse.ReadWriteAddresses[senderAddr]
	   				workerMock.On("MoveTxToNotReady", txHash, senderAddr, addressInfo.Nonce, addressInfo.Balance).Return([]*TxTracker{}).Once()
	   			}
	   			if tc.expectedUpdateTxCall {
	   				workerMock.On("UpdateTxZKCounters", txTracker.Hash, txTracker.From, tc.executorResponse.UsedZkCounters).Return().Once()
	   			}
	   			if tc.expectedError == nil {
	   				//stateMock.On("GetGasPrices", ctx).Return(pool.GasPrices{L1GasPrice: 0, L2GasPrice: 0}, nilErr).Once()
	   				workerMock.On("DeleteTx", txTracker.Hash, txTracker.From).Return().Once()
	   				workerMock.On("UpdateAfterSingleSuccessfulTxExecution", txTracker.From, tc.executorResponse.ReadWriteAddresses).Return([]*TxTracker{}).Once()
	   				workerMock.On("AddPendingTxToStore", txTracker.Hash, txTracker.From).Return().Once()
	   			}
	   			if tc.expectedUpdateTxStatus != "" {
	   				stateMock.On("UpdateTxStatus", ctx, txHash, tc.expectedUpdateTxStatus, false, mock.Anything).Return(nil).Once()
	   			}

	   			errWg, err := f.handleProcessTransactionResponse(ctx, txTracker, tc.executorResponse, tc.oldStateRoot)
	   			if errWg != nil {
	   				errWg.Wait()
	   			}

	   			if tc.expectedError != nil {
	   				require.Equal(t, tc.expectedError, err)
	   			} else {
	   				require.Nil(t, err)
	   			}

	   			if tc.expectedStoredTx.batchResponse != nil {
	   				close(f.pendingL2BlocksToStore) // close the channel
	   				<-done                          // wait for the goroutine to finish
	   				f.pendingL2BlocksToStoreWG.Wait()
	   				require.Len(t, storedTxs, 1)
	   				actualTx := storedTxs[0] //nolint:gosec
	   				assertEqualTransactionToStore(t, tc.expectedStoredTx, actualTx)
	   			} else {
	   				require.Empty(t, storedTxs)
	   			}

	   			workerMock.AssertExpectations(t)
	   			stateMock.AssertExpectations(t)
	   		})
	   	}
}*/

/*func assertEqualTransactionToStore(t *testing.T, expectedTx, actualTx transactionToStore) {
	   require.Equal(t, expectedTx.from, actualTx.from)
	   require.Equal(t, expectedTx.hash, actualTx.hash)
	   require.Equal(t, expectedTx.response, actualTx.response)
	   require.Equal(t, expectedTx.batchNumber, actualTx.batchNumber)
	   require.Equal(t, expectedTx.timestamp, actualTx.timestamp)
	   require.Equal(t, expectedTx.coinbase, actualTx.coinbase)
	   require.Equal(t, expectedTx.oldStateRoot, actualTx.oldStateRoot)
	   require.Equal(t, expectedTx.isForcedBatch, actualTx.isForcedBatch)
	   require.Equal(t, expectedTx.flushId, actualTx.flushId)
}*/

/*func TestFinalizer_newWIPBatch(t *testing.T) {
	// arrange
	now = testNow
	defer func() {
		now = time.Now
	}()

	f = setupFinalizer(true)

	processRequest := state.ProcessRequest{
		Caller:       stateMetrics.SequencerCallerLabel,
		Timestamp_V1: now(),
		Transactions: decodedBatchL2Data,
	}
	stateRootErr := errors.New("state root must have value to close batch")
	txs := []types.Transaction{*tx}
	require.NoError(t, err)
	newBatchNum := f.wipBatch.batchNumber + 1
	expectedNewWipBatch := &Batch{
		batchNumber:        newBatchNum,
		coinbase:           f.sequencerAddress,
		initialStateRoot:   newHash,
		stateRoot:          newHash,
		timestamp:          now(),
		remainingResources: getMaxRemainingResources(f.batchConstraints),
	}
	closeBatchParams := ClosingBatchParameters{
		BatchNumber:          f.wipBatch.batchNumber,
	}

	batches := []*state.Batch{
		{
			BatchNumber:    f.wipBatch.batchNumber,
			StateRoot:      newHash,
			GlobalExitRoot: oldHash,
			Transactions:   txs,
			Timestamp:      now(),
			BatchL2Data:    decodedBatchL2Data,
		},
	}

	// For Empty Batch
	expectedNewWipEmptyBatch := *expectedNewWipBatch
	expectedNewWipEmptyBatch.initialStateRoot = oldHash
	expectedNewWipEmptyBatch.stateRoot = oldHash
	emptyBatch := *batches[0]
	emptyBatch.StateRoot = oldHash
	emptyBatch.Transactions = make([]types.Transaction, 0)
	emptyBatch.BatchL2Data = []byte{}
	emptyBatch.GlobalExitRoot = oldHash
	emptyBatchBatches := []*state.Batch{&emptyBatch}
	closeBatchParamsForEmptyBatch := closeBatchParams

	// For Forced Batch
	expectedForcedNewWipBatch := *expectedNewWipBatch
	expectedForcedNewWipBatch.batchNumber = expectedNewWipBatch.batchNumber + 1

	testCases := []struct {
		name                       string
		batches                    []*state.Batch
		closeBatchErr              error
		closeBatchParams           ClosingBatchParameters
		stateRootAndLERErr         error
		openBatchErr               error
		expectedWip                *Batch
		reprocessFullBatchResponse *state.ProcessBatchResponse
		expectedErr                error
		reprocessBatchErr          error
		forcedBatches              []state.ForcedBatch
	}{

		{
			name:               "Error StateRoot must have value",
			stateRootAndLERErr: stateRootErr,
			expectedErr:        stateRootErr,
		},
		{
			name:             "Error Close Batch",
			expectedWip:      expectedNewWipBatch,
			closeBatchParams: closeBatchParams,
			batches:          batches,
			closeBatchErr:    testErr,
			expectedErr:      fmt.Errorf("failed to close batch, err: %v", testErr),
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     f.wipBatch.stateRoot,
				NewLocalExitRoot: f.wipBatch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
		{
			name:             "Error Open Batch",
			expectedWip:      expectedNewWipBatch,
			closeBatchParams: closeBatchParams,
			batches:          batches,
			openBatchErr:     testErr,
			expectedErr:      fmt.Errorf("failed to open new batch, err: %v", testErr),
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     f.wipBatch.stateRoot,
				NewLocalExitRoot: f.wipBatch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
		{
			name:             "Success with closing non-empty batch",
			expectedWip:      expectedNewWipBatch,
			closeBatchParams: closeBatchParams,
			batches:          batches,
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     f.wipBatch.stateRoot,
				NewLocalExitRoot: f.wipBatch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
		{
			name:             "Success with closing empty batch",
			expectedWip:      &expectedNewWipEmptyBatch,
			closeBatchParams: closeBatchParamsForEmptyBatch,
			batches:          emptyBatchBatches,
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     oldHash,
				NewLocalExitRoot: f.wipBatch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
		{
			name: "Forced Batches",
			forcedBatches: []state.ForcedBatch{{
				BlockNumber:       1,
				ForcedBatchNumber: 1,
				Sequencer:         seqAddr,
				GlobalExitRoot:    oldHash,
				RawTxsData:        nil,
				ForcedAt:          now(),
			}},
			expectedWip:      &expectedForcedNewWipBatch,
			closeBatchParams: closeBatchParams,
			batches:          batches,
			reprocessFullBatchResponse: &state.ProcessBatchResponse{
				NewStateRoot:     f.wipBatch.stateRoot,
				NewLocalExitRoot: f.wipBatch.localExitRoot,
				IsRomOOCError:    false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			processRequest.GlobalExitRoot_V1 = oldHash
			processRequest.OldStateRoot = oldHash
			processRequest.BatchNumber = f.wipBatch.batchNumber
			f.nextForcedBatches = tc.forcedBatches

			currTxs := txs
			if tc.closeBatchParams.StateRoot == oldHash {
				currTxs = nil
				f.wipBatch.stateRoot = oldHash
				processRequest.Transactions = []byte{}
				defer func() {
					f.wipBatch.stateRoot = newHash
					processRequest.Transactions = decodedBatchL2Data
				}()

				executorMock.On("ProcessBatch", ctx, processRequest, true).Return(tc.reprocessFullBatchResponse, tc.reprocessBatchErr).Once()
			}

			if tc.stateRootAndLERErr == nil {
				stateMock.On("CloseBatch", ctx, tc.closeBatchParams).Return(tc.closeBatchErr).Once()
				stateMock.On("GetBatchByNumber", ctx, f.wipBatch.batchNumber, nil).Return(tc.batches[0], nilErr).Once()
				stateMock.On("GetForkIDByBatchNumber", f.wipBatch.batchNumber).Return(uint64(5))
				stateMock.On("GetTransactionsByBatchNumber", ctx, f.wipBatch.batchNumber).Return(currTxs, constants.EffectivePercentage, nilErr).Once()
				if tc.forcedBatches != nil && len(tc.forcedBatches) > 0 {
					fbProcessRequest := processRequest
					fbProcessRequest.BatchNumber = processRequest.BatchNumber + 1
					fbProcessRequest.OldStateRoot = newHash
					fbProcessRequest.Transactions = nil
					stateMock.On("GetLastTrustedForcedBatchNumber", ctx, nil).Return(tc.forcedBatches[0].ForcedBatchNumber-1, nilErr).Once()
					stateMock.On("ProcessForcedBatch", tc.forcedBatches[0].ForcedBatchNumber, fbProcessRequest).Return(tc.reprocessFullBatchResponse, nilErr).Once()
				}
				if tc.closeBatchErr == nil {
					stateMock.On("BeginStateTransaction", ctx).Return(dbTxMock, nilErr).Once()
					stateMock.On("OpenBatch", ctx, mock.Anything, dbTxMock).Return(tc.openBatchErr).Once()
					if tc.openBatchErr == nil {
						dbTxMock.On("Commit", ctx).Return(nilErr).Once()
					} else {
						dbTxMock.On("Rollback", ctx).Return(nilErr).Once()
					}
				}
				executorMock.On("ProcessBatch", ctx, processRequest, false).Return(tc.reprocessFullBatchResponse, tc.reprocessBatchErr).Once()
			}

			if tc.stateRootAndLERErr != nil {
				f.wipBatch.stateRoot = state.ZeroHash
				f.wipBatch.localExitRoot = state.ZeroHash
				defer func() {
					f.wipBatch.stateRoot = newHash
					f.wipBatch.localExitRoot = newHash
				}()
			}

			// act
			wipBatch, err := f.closeAndOpenNewWIPBatch(ctx)

			// assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Nil(t, wipBatch)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedWip, wipBatch)
			}
			stateMock.AssertExpectations(t)
			dbTxMock.AssertExpectations(t)
			executorMock.AssertExpectations(t)
		})
	}
}*/

/*func TestFinalizer_processForcedBatches(t *testing.T) {
	var err error
	f = setupFinalizer(false)
	now = testNow
	defer func() {
		now = time.Now
	}()
	ctx = context.Background()
	RawTxsData1 := make([]byte, 0, 2)
	RawTxsData1 = append(RawTxsData1, []byte(testBatchL2DataAsString)...)
	RawTxsData2 := make([]byte, 0, 2)
	RawTxsData2 = append(RawTxsData2, []byte(testBatchL2DataAsString)...)
	batchNumber := f.wipBatch.batchNumber
	decodedBatchL2Data, err = hex.DecodeHex(testBatchL2DataAsString)
	require.NoError(t, err)

	tx1 := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1), 100000, big.NewInt(1), RawTxsData1)
	tx2 := types.NewTransaction(1, common.HexToAddress("0x2"), big.NewInt(1), 100000, big.NewInt(1), RawTxsData2)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(pvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)
	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)

	txResp1 := &state.ProcessTransactionResponse{
		TxHash:    signedTx1.Hash(),
		StateRoot: stateRootHashes[0],
		Tx:        *signedTx1,
	}

	blockResp1 := &state.ProcessBlockResponse{
		TransactionResponses: []*state.ProcessTransactionResponse{txResp1},
	}

	txResp2 := &state.ProcessTransactionResponse{
		TxHash:    signedTx2.Hash(),
		StateRoot: stateRootHashes[1],
		Tx:        *signedTx2,
	}

	blockResp2 := &state.ProcessBlockResponse{
		TransactionResponses: []*state.ProcessTransactionResponse{txResp2},
	}

	batchResponse1 := &state.ProcessBatchResponse{
		NewBatchNumber: f.wipBatch.batchNumber + 1,
		BlockResponses: []*state.ProcessBlockResponse{blockResp1},
		NewStateRoot:   newHash,
	}

	batchResponse2 := &state.ProcessBatchResponse{
		NewBatchNumber: f.wipBatch.batchNumber + 2,
		BlockResponses: []*state.ProcessBlockResponse{blockResp2},
		NewStateRoot:   newHash2,
	}
	forcedBatch1 := state.ForcedBatch{
		ForcedBatchNumber: 2,
		GlobalExitRoot:    oldHash,
		RawTxsData:        RawTxsData1,
	}
	forcedBatch2 := state.ForcedBatch{
		ForcedBatchNumber: 3,
		GlobalExitRoot:    oldHash,
		RawTxsData:        RawTxsData2,
	}
	testCases := []struct {
		name                            string
		forcedBatches                   []state.ForcedBatch
		getLastTrustedForcedBatchNumErr error
		expectedErr                     error
		expectedStoredTx                []transactionToStore
		processInBetweenForcedBatch     bool
		getForcedBatchError             error
	}{
		{
			name:          "Success",
			forcedBatches: []state.ForcedBatch{forcedBatch1, forcedBatch2},
			expectedStoredTx: []transactionToStore{
				{
					hash:          signedTx1.Hash(),
					from:          auth.From,
					batchResponse: batchResponse1,
					batchNumber:   f.wipBatch.batchNumber + 1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  stateRootHashes[0],
					isForcedBatch: true,
					response:      txResp1,
				},
				{
					hash:          signedTx2.Hash(),
					from:          auth.From,
					batchResponse: batchResponse2,
					batchNumber:   f.wipBatch.batchNumber + 2,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  stateRootHashes[1],
					isForcedBatch: true,
					response:      txResp2,
				},
			},
		},
		{
			name:                            "GetLastTrustedForcedBatchNumber_Error",
			forcedBatches:                   []state.ForcedBatch{forcedBatch1},
			getLastTrustedForcedBatchNumErr: testErr,
			expectedErr:                     fmt.Errorf("failed to get last trusted forced batch number, err: %s", testErr),
		},
		{
			name:          "Skip Already Processed Forced Batches",
			forcedBatches: []state.ForcedBatch{{ForcedBatchNumber: 1}},
		},
		{
			name: "Process In-Between Unprocessed Forced Batches",
			forcedBatches: []state.ForcedBatch{
				forcedBatch2,
				forcedBatch1,
			},
			expectedStoredTx: []transactionToStore{
				{
					hash:          signedTx1.Hash(),
					from:          auth.From,
					batchResponse: batchResponse1,
					batchNumber:   f.wipBatch.batchNumber + 1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  stateRootHashes[0],
					isForcedBatch: true,
					response:      txResp1,
				},
				{
					hash:          signedTx2.Hash(),
					from:          auth.From,
					batchResponse: batchResponse2,
					batchNumber:   f.wipBatch.batchNumber + 2,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  stateRootHashes[1],
					isForcedBatch: true,
					response:      txResp2,
				},
			},
			processInBetweenForcedBatch: true,
		},
		{
			name: "Error GetForcedBatch",
			forcedBatches: []state.ForcedBatch{
				forcedBatch2,
				forcedBatch1,
			},
			expectedErr:                 fmt.Errorf("failed to get in-between forced batch %d, err: %s", 2, testErr),
			getForcedBatchError:         testErr, // This will be used to simulate the error
			processInBetweenForcedBatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			var newStateRoot common.Hash
			stateRoot := oldHash
			storedTxs := make([]transactionToStore, 0)
			f.pendingL2BlocksToStore = make(chan transactionToStore)
			if tc.expectedStoredTx != nil && len(tc.expectedStoredTx) > 0 {
				done = make(chan bool) // init a new done channel
				go func() {
					for tx := range f.pendingL2BlocksToStore {
						storedTxs = append(storedTxs, tx)
						f.pendingL2BlocksToStoreWG.Done()
					}
					done <- true // signal that the goroutine is done
				}()
			}
			f.nextForcedBatches = make([]state.ForcedBatch, len(tc.forcedBatches))
			copy(f.nextForcedBatches, tc.forcedBatches)
			internalBatchNumber := batchNumber
			stateMock.On("GetLastTrustedForcedBatchNumber", ctx, nil).Return(uint64(1), tc.getLastTrustedForcedBatchNumErr).Once()
			tc.forcedBatches = f.sortForcedBatches(tc.forcedBatches)

			if tc.getLastTrustedForcedBatchNumErr == nil {
				for i, forcedBatch := range tc.forcedBatches {
					// Skip already processed forced batches.
					if forcedBatch.ForcedBatchNumber == 1 {
						continue
					}

					internalBatchNumber += 1
					processRequest := state.ProcessRequest{
						BatchNumber:       internalBatchNumber,
						OldStateRoot:      stateRootHashes[i],
						GlobalExitRoot_V1: forcedBatch.GlobalExitRoot,
						Transactions:      forcedBatch.RawTxsData,
						Coinbase:          f.sequencerAddress,
						Timestamp_V1:      now(),
						Caller:            stateMetrics.SequencerCallerLabel,
					}
					var currResp *state.ProcessBatchResponse
					if tc.expectedStoredTx == nil {
						currResp = &state.ProcessBatchResponse{
							NewStateRoot:   stateRootHashes[i+1],
							NewBatchNumber: internalBatchNumber,
						}
					} else {
						for _, storedTx := range tc.expectedStoredTx {
							if storedTx.batchNumber == internalBatchNumber {
								currResp = storedTx.batchResponse
								break
							}
						}
					}
					stateMock.On("ProcessForcedBatch", forcedBatch.ForcedBatchNumber, processRequest).Return(currResp, nilErr).Once()
				}

				if tc.processInBetweenForcedBatch {
					stateMock.On("GetForcedBatch", ctx, uint64(2), nil).Return(&forcedBatch1, tc.getForcedBatchError).Once()
				}
			}

			workerMock.On("DeleteForcedTx", mock.Anything, mock.Anything).Return()
			workerMock.On("AddPendingTxToStore", mock.Anything, mock.Anything).Return()
			workerMock.On("AddForcedTx", mock.Anything, mock.Anything).Return()

			// act
			batchNumber, newStateRoot, err = f.processForcedBatches(ctx, batchNumber, stateRoot)

			// assert
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				if tc.expectedStoredTx != nil && len(tc.expectedStoredTx) > 0 {
					close(f.pendingL2BlocksToStore) // ensure the channel is closed
					<-done                          // wait for the goroutine to finish
					f.pendingL2BlocksToStoreWG.Wait()
					for i := range tc.expectedStoredTx {
						require.Equal(t, tc.expectedStoredTx[i], storedTxs[i])
					}
				}
				if len(tc.expectedStoredTx) > 0 {
					assert.Equal(t, stateRootHashes[len(stateRootHashes)-1], newStateRoot)
				}
				assert.Equal(t, batchNumber, internalBatchNumber)
				assert.NoError(t, tc.expectedErr)
				stateMock.AssertExpectations(t)
			}
		})
	}
}*/

/*func TestFinalizer_openWIPBatch(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	now = testNow
	defer func() {
		now = time.Now
	}()
	batchNum := f.wipBatch.batchNumber + 1
	expectedWipBatch := &Batch{
		batchNumber:        batchNum,
		coinbase:           f.sequencerAddress,
		initialStateRoot:   oldHash,
		imStateRoot:        oldHash,
		timestamp:          now(),
		remainingResources: getMaxRemainingResources(f.batchConstraints),
	}
	testCases := []struct {
		name         string
		openBatchErr error
		beginTxErr   error
		commitErr    error
		rollbackErr  error
		expectedWip  *Batch
		expectedErr  error
	}{
		{
			name:        "Success",
			expectedWip: expectedWipBatch,
		},
		{
			name:        "Error BeginTransaction",
			beginTxErr:  testErr,
			expectedErr: fmt.Errorf("failed to begin state transaction to open batch, err: %v", testErr),
		},
		{
			name:         "Error OpenBatch",
			openBatchErr: testErr,
			expectedErr:  fmt.Errorf("failed to open new batch, err: %v", testErr),
		},
		{
			name:        "Error Commit",
			commitErr:   testErr,
			expectedErr: fmt.Errorf("failed to commit database transaction for opening a batch, err: %v", testErr),
		},
		{
			name:         "Error Rollback",
			openBatchErr: testErr,
			rollbackErr:  testErr,
			expectedErr: fmt.Errorf(
				"failed to rollback dbTx: %s. Rollback err: %v",
				testErr.Error(), openBatchError,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			stateMock.On("BeginStateTransaction", ctx).Return(dbTxMock, tc.beginTxErr).Once()
			if tc.beginTxErr == nil {
				stateMock.On("OpenBatch", ctx, mock.Anything, dbTxMock).Return(tc.openBatchErr).Once()
			}

			if tc.expectedErr != nil && (tc.rollbackErr != nil || tc.openBatchErr != nil) {
				dbTxMock.On("Rollback", ctx).Return(tc.rollbackErr).Once()
			}

			if tc.expectedErr == nil || tc.commitErr != nil {
				dbTxMock.On("Commit", ctx).Return(tc.commitErr).Once()
			}

			// act
			wipBatch, err := f.openNewWIPBatch(ctx, batchNum, oldHash, oldHash, oldHash)

			// assert
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Nil(t, wipBatch)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedWip, wipBatch)
			}
			stateMock.AssertExpectations(t)
			dbTxMock.AssertExpectations(t)
		})
	}
}*/

// TestFinalizer_finalizeSIPBatch tests the finalizeSIPBatch method.
func TestFinalizer_finalizeSIPBatch(t *testing.T) {
	managerErr := fmt.Errorf("some err")

	testCases := []struct {
		name        string
		managerErr  error
		expectedErr error
	}{

		{
			name:        "Success",
			managerErr:  nil,
			expectedErr: nil,
		},
		{
			name:        "Error Manager",
			managerErr:  managerErr,
			expectedErr: managerErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f = setupFinalizer(true)
			// set wip batch has at least one L2 block as it can not be closed empty
			f.sipBatch.countOfL2Blocks++

			usedResources := getUsedBatchResources(f.batchConstraints, f.wipBatch.imRemainingResources)

			receipt := state.ProcessingReceipt{
				BatchNumber:    f.wipBatch.batchNumber,
				BatchResources: usedResources,
				ClosingReason:  f.wipBatch.closingReason,
			}

			// arrange
			stateMock.On("BeginStateTransaction", ctx).Return(dbTxMock, nilErr).Once()
			stateMock.On("CloseWIPBatch", ctx, receipt, mock.Anything).Return(tc.managerErr).Once()

			if tc.managerErr == nil {
				stateMock.On("GetBatchByNumber", ctx, f.sipBatch.batchNumber, nil).Return(&state.Batch{BatchNumber: f.sipBatch.batchNumber}, nilErr).Once()
				stateMock.On("GetForkIDByBatchNumber", f.wipBatch.batchNumber).Return(uint64(9)).Once()
				stateMock.On("GetL1InfoTreeDataFromBatchL2Data", ctx, mock.Anything, nil).Return(map[uint32]state.L1DataV2{}, state.ZeroHash, state.ZeroHash, nil)
				stateMock.On("ProcessBatchV2", ctx, mock.Anything, false).Return(&state.ProcessBatchResponse{}, "", nil)
				stateMock.On("UpdateBatchAsChecked", ctx, f.sipBatch.batchNumber, nil).Return(nil)
				dbTxMock.On("Commit", ctx).Return(nilErr).Once()
			} else {
				dbTxMock.On("Rollback", ctx).Return(nilErr).Once()
			}

			// act
			err := f.finalizeSIPBatch(ctx)

			// assert
			if tc.expectedErr != nil {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFinalizer_isDeadlineEncountered(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	now = testNow
	defer func() {
		now = time.Now
	}()
	testCases := []struct {
		name                        string
		nextForcedBatch             int64
		nextGER                     int64
		nextDelayedBatch            int64
		expected                    bool
		timestampResolutionDeadline bool
	}{
		{
			name:     "No deadlines",
			expected: false,
		},
		{
			name:            "Forced batch deadline",
			nextForcedBatch: now().Add(time.Second).Unix(),
			expected:        true,
		},
		{
			name:             "Delayed batch deadline",
			nextDelayedBatch: now().Add(time.Second).Unix(),
			expected:         false,
		},
		{
			name:                        "Timestamp resolution deadline",
			timestampResolutionDeadline: true,
			expected:                    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f.nextForcedBatchDeadline = tc.nextForcedBatch
			if tc.expected == true {
				now = func() time.Time {
					return testNow().Add(time.Second * 2)
				}
			}

			// specifically for "Timestamp resolution deadline" test case
			if tc.timestampResolutionDeadline == true {
				// ensure that the batch is not empty and the timestamp is in the past
				f.wipBatch.timestamp = now().Add(-f.cfg.BatchMaxDeltaTimestamp.Duration * 2)
				f.wipBatch.countOfL2Blocks = 1
			}

			// act
			actual, _ := f.checkIfFinalizeBatch()

			// assert
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestFinalizer_checkRemainingResources(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	ctx = context.Background()
	txResponse := &state.ProcessTransactionResponse{TxHash: oldHash}
	blockResponse := &state.ProcessBlockResponse{
		TransactionResponses: []*state.ProcessTransactionResponse{txResponse},
	}
	result := &state.ProcessBatchResponse{
		UsedZkCounters: state.ZKCounters{GasUsed: 1000},
		BlockResponses: []*state.ProcessBlockResponse{blockResponse},
	}
	remainingResources := state.BatchResources{
		ZKCounters: state.ZKCounters{GasUsed: 9000},
		Bytes:      10000,
	}
	f.wipBatch.imRemainingResources = remainingResources
	testCases := []struct {
		name                 string
		remaining            state.BatchResources
		overflow             bool
		overflowResource     string
		expectedWorkerUpdate bool
		expectedTxTracker    *TxTracker
	}{
		{
			name:                 "Success",
			remaining:            remainingResources,
			overflow:             false,
			expectedWorkerUpdate: false,
			expectedTxTracker:    &TxTracker{RawTx: []byte("test")},
		},
		{
			name: "Bytes Resource Exceeded",
			remaining: state.BatchResources{
				Bytes: 0,
			},
			overflow:             true,
			overflowResource:     "Bytes",
			expectedWorkerUpdate: true,
			expectedTxTracker:    &TxTracker{RawTx: []byte("test")},
		},
		{
			name: "ZkCounter Resource Exceeded",
			remaining: state.BatchResources{
				ZKCounters: state.ZKCounters{GasUsed: 0},
			},
			overflow:             true,
			overflowResource:     "CumulativeGas",
			expectedWorkerUpdate: true,
			expectedTxTracker:    &TxTracker{RawTx: make([]byte, 0)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f.wipBatch.imRemainingResources = tc.remaining
			stateMock.On("AddEvent", ctx, mock.Anything, nil).Return(nil)
			if tc.expectedWorkerUpdate {
				workerMock.On("UpdateTxZKCounters", txResponse.TxHash, tc.expectedTxTracker.From, result.UsedZkCounters).Return().Once()
			}

			// act
			overflow, overflowResource := f.wipBatch.imRemainingResources.Sub(state.BatchResources{ZKCounters: result.UsedZkCounters, Bytes: uint64(len(tc.expectedTxTracker.RawTx))})

			// assert
			assert.Equal(t, tc.overflow, overflow)
			assert.Equal(t, tc.overflowResource, overflowResource)
		})
	}
}

/*func TestFinalizer_handleTransactionError(t *testing.T) {
	// arrange
	f = setupFinalizer(true)
	nonce := uint64(0)
	txTracker := &TxTracker{Hash: txHash, From: senderAddr, Cost: big.NewInt(0)}
	testCases := []struct {
		name               string
		err                executor.RomError
		expectedDeleteCall bool
		updateTxStatus     pool.TxStatus
		expectedMoveCall   bool
		isRoomOOC          bool
	}{
		{
			name:               "Error OutOfCounters",
			err:                executor.RomError_ROM_ERROR_OUT_OF_COUNTERS_STEP,
			updateTxStatus:     pool.TxStatusInvalid,
			expectedDeleteCall: true,
			isRoomOOC:          true,
		},
		{
			name:             "Error IntrinsicInvalidNonce",
			err:              executor.RomError_ROM_ERROR_INTRINSIC_INVALID_NONCE,
			updateTxStatus:   pool.TxStatusFailed,
			expectedMoveCall: true,
		},
		{
			name:             "Error IntrinsicInvalidBalance",
			err:              executor.RomError_ROM_ERROR_INTRINSIC_INVALID_BALANCE,
			updateTxStatus:   pool.TxStatusFailed,
			expectedMoveCall: true,
		},
		{
			name:               "Error IntrinsicErrorChainId",
			err:                executor.RomError_ROM_ERROR_INTRINSIC_INVALID_CHAIN_ID,
			updateTxStatus:     pool.TxStatusFailed,
			expectedDeleteCall: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			if tc.expectedDeleteCall {
				workerMock.On("DeleteTx", txHash, senderAddr).Return()
				stateMock.On("UpdateTxStatus", ctx, txHash, tc.updateTxStatus, false, mock.Anything).Return(nil).Once()
			}
			if tc.expectedMoveCall {
				workerMock.On("MoveTxToNotReady", txHash, senderAddr, &nonce, big.NewInt(0)).Return([]*TxTracker{
					{
						Hash: txHash2,
					},
				}).Once()

				stateMock.On("UpdateTxStatus", ctx, txHash2, pool.TxStatusFailed, false, mock.Anything).Return(nil).Once()
			}

			result := &state.ProcessBatchResponse{
				IsRomOOCError: tc.isRoomOOC,
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					senderAddr: {Nonce: &nonce, Balance: big.NewInt(0)},
				},
				BlockResponses: []*state.ProcessBlockResponse{
					{
						TransactionResponses: []*state.ProcessTransactionResponse{
							{
								RomError: executor.RomErr(tc.err),
							},
						},
					},
				},
			}

			// act
			wg := f.handleProcessTransactionError(ctx, result, txTracker)
			if wg != nil {
				wg.Wait()
			}

			// assert
			workerMock.AssertExpectations(t)
		})
	}
}*/

/*func Test_processTransaction(t *testing.T) {
	f = setupFinalizer(true)
	gasUsed := uint64(100000)
	txTracker := &TxTracker{
		Hash:              txHash,
		From:              senderAddr,
		Nonce:             nonce1,
		GasPrice:          effectiveGasPrice,
		EffectiveGasPrice: effectiveGasPrice,
		L1GasPrice:        l1GasPrice,
		EGPLog: state.EffectiveGasPriceLog{
			ValueFinal:     new(big.Int).SetUint64(0),
			ValueFirst:     new(big.Int).SetUint64(0),
			ValueSecond:    new(big.Int).SetUint64(0),
			FinalDeviation: new(big.Int).SetUint64(0),
			MaxDeviation:   new(big.Int).SetUint64(0),
			GasPrice:       new(big.Int).SetUint64(0),
		},
		BatchResources: state.BatchResources{
			Bytes: 1000,
			ZKCounters: state.ZKCounters{
				GasUsed: 500,
			},
		},
		RawTx: []byte{0, 0, 1, 2, 3, 4, 5},
	}
	successfulTxResponse := &state.ProcessTransactionResponse{
		TxHash:    txHash,
		StateRoot: newHash,
		GasUsed:   gasUsed,
	}
	successfulBlockResponse := &state.ProcessBlockResponse{
		TransactionResponses: []*state.ProcessTransactionResponse{
			successfulTxResponse,
		},
	}

	successfulBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: newHash,
		BlockResponses: []*state.ProcessBlockResponse{
			successfulBlockResponse,
		},
		ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
			senderAddr: {
				Nonce: &nonce2,
			},
		},
	}
	outOfCountersErrBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: oldHash,
		BlockResponses: []*state.ProcessBlockResponse{
			{
				TransactionResponses: []*state.ProcessTransactionResponse{
					{
						StateRoot: oldHash,
						RomError:  runtime.ErrOutOfCountersKeccak,
						GasUsed:   gasUsed,
					},
				},
			},
		},
		IsRomOOCError: true,
	}
	outOfCountersExecutorErrBatchResp := *outOfCountersErrBatchResp
	outOfCountersExecutorErrBatchResp.IsRomOOCError = false
	testCases := []struct {
		name                   string
		ctx                    context.Context
		tx                     *TxTracker
		expectedResponse       *state.ProcessBatchResponse
		executorErr            error
		expectedErr            error
		expectedStoredTx       transactionToStore
		expectedUpdateTxStatus pool.TxStatus
	}{
		{
			name:             "Successful transaction processing",
			ctx:              context.Background(),
			tx:               txTracker,
			expectedResponse: successfulBatchResp,
			expectedStoredTx: transactionToStore{
				hash:          txHash,
				from:          senderAddr,
				batchNumber:   f.wipBatch.batchNumber,
				coinbase:      f.wipBatch.coinbase,
				timestamp:     f.wipBatch.timestamp,
				oldStateRoot:  newHash,
				batchResponse: successfulBatchResp,
				isForcedBatch: false,
				response:      successfulTxResponse,
			},
		},
		{
			name:                   "Out Of Counters err",
			ctx:                    context.Background(),
			tx:                     txTracker,
			expectedResponse:       outOfCountersErrBatchResp,
			expectedErr:            runtime.ErrOutOfCountersKeccak,
			expectedUpdateTxStatus: pool.TxStatusInvalid,
		},
		{
			name:                   "Executor err",
			ctx:                    context.Background(),
			tx:                     txTracker,
			expectedResponse:       &outOfCountersExecutorErrBatchResp,
			executorErr:            runtime.ErrOutOfCountersKeccak,
			expectedErr:            runtime.ErrOutOfCountersKeccak,
			expectedUpdateTxStatus: pool.TxStatusInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storedTxs := make([]transactionToStore, 0)
			f.pendingL2BlocksToStore = make(chan transactionToStore, 1)
			if tc.expectedStoredTx.batchResponse != nil {
				done = make(chan bool) // init a new done channel
				go func() {
					for tx := range f.pendingL2BlocksToStore {
						storedTxs = append(storedTxs, tx)
						f.pendingL2BlocksToStoreWG.Done()
					}
					done <- true // signal that the goroutine is done
				}()
			}

			stateMock.On("GetL1AndL2GasPrice").Return(uint64(1000000), uint64(100000)).Once()
			executorMock.On("ProcessBatch", tc.ctx, mock.Anything, true).Return(tc.expectedResponse, tc.executorErr).Once()
			if tc.executorErr == nil {
				workerMock.On("DeleteTx", tc.tx.Hash, tc.tx.From).Return().Once()
				stateMock.On("GetForkIDByBatchNumber", mock.Anything).Return(forkId5)
			}
			if tc.expectedErr == nil {
				workerMock.On("UpdateAfterSingleSuccessfulTxExecution", tc.tx.From, tc.expectedResponse.ReadWriteAddresses).Return([]*TxTracker{}).Once()
				workerMock.On("AddPendingTxToStore", tc.tx.Hash, tc.tx.From).Return().Once()
			}

			if tc.expectedUpdateTxStatus != "" {
				stateMock.On("UpdateTxStatus", tc.ctx, txHash, tc.expectedUpdateTxStatus, false, mock.Anything).Return(nil)
			}

			if errors.Is(tc.executorErr, runtime.ErrOutOfCountersKeccak) {
				workerMock.On("DeleteTx", tc.tx.Hash, tc.tx.From).Return().Once()
			}

			errWg, err := f.processTransaction(tc.ctx, tc.tx, true)

			if tc.expectedStoredTx.batchResponse != nil {
				close(f.pendingL2BlocksToStore) // ensure the channel is closed
				<-done                          // wait for the goroutine to finish
				f.pendingL2BlocksToStoreWG.Wait()
				// require.Equal(t, tc.expectedStoredTx, storedTxs[0])
			}
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			if errWg != nil {
				errWg.Wait()
			}

			workerMock.AssertExpectations(t)
			stateMock.AssertExpectations(t)
		})
	}
}*/

/*func Test_handleForcedTxsProcessResp(t *testing.T) {
	var chainID = new(big.Int).SetInt64(400)
	var pvtKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
	RawTxsData1 := make([]byte, 0, 2)
	RawTxsData2 := make([]byte, 0, 2)

	f = setupFinalizer(false)
	now = testNow
	defer func() {
		now = time.Now
	}()

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(pvtKey, "0x"))
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)

	tx1 := types.NewTransaction(0, common.HexToAddress("0x1"), big.NewInt(1), 100000, big.NewInt(1), RawTxsData1)
	tx2 := types.NewTransaction(1, common.HexToAddress("0x2"), big.NewInt(1), 100000, big.NewInt(1), RawTxsData2)

	signedTx1, err := auth.Signer(auth.From, tx1)
	require.NoError(t, err)

	signedTx2, err := auth.Signer(auth.From, tx2)
	require.NoError(t, err)

	tx1Plustx2, err := state.EncodeTransactions([]types.Transaction{*signedTx1, *signedTx2}, nil, 4)
	require.NoError(t, err)

	ctx = context.Background()
	txResponseOne := &state.ProcessTransactionResponse{
		TxHash:    signedTx1.Hash(),
		StateRoot: newHash,
		RomError:  nil,
		Tx:        *signedTx1,
	}
	txResponseTwo := &state.ProcessTransactionResponse{
		TxHash:    signedTx2.Hash(),
		StateRoot: newHash2,
		RomError:  nil,
		Tx:        *signedTx2,
	}
	blockResponseOne := &state.ProcessBlockResponse{
		TransactionResponses: []*state.ProcessTransactionResponse{
			txResponseOne,
		},
	}
	blockResponseTwo := &state.ProcessBlockResponse{
		TransactionResponses: []*state.ProcessTransactionResponse{
			txResponseTwo,
		},
	}
	successfulBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: newHash,
		BlockResponses: []*state.ProcessBlockResponse{
			blockResponseOne,
			blockResponseTwo,
		},
	}
	txResponseReverted := &state.ProcessTransactionResponse{
		Tx:        *signedTx1,
		TxHash:    signedTx1.Hash(),
		RomError:  runtime.ErrExecutionReverted,
		StateRoot: newHash,
	}
	blockResponseReverted := &state.ProcessBlockResponse{
		TransactionResponses: []*state.ProcessTransactionResponse{
			txResponseReverted,
		},
	}
	revertedBatchResp := &state.ProcessBatchResponse{
		BlockResponses: []*state.ProcessBlockResponse{
			blockResponseReverted,
		},
	}
	txResponseIntrinsicErr := &state.ProcessTransactionResponse{
		Tx:        *signedTx1,
		TxHash:    signedTx1.Hash(),
		RomError:  runtime.ErrIntrinsicInvalidChainID,
		StateRoot: newHash,
	}
	blockResponseIntrinsicErr := &state.ProcessBlockResponse{
		TransactionResponses: []*state.ProcessTransactionResponse{
			txResponseIntrinsicErr,
		},
	}

	intrinsicErrBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: newHash,
		BlockResponses: []*state.ProcessBlockResponse{
			blockResponseOne,
			blockResponseIntrinsicErr,
		},
	}

	testCases := []struct {
		name              string
		request           state.ProcessRequest
		result            *state.ProcessBatchResponse
		oldStateRoot      common.Hash
		expectedStoredTxs []transactionToStore
	}{
		{
			name: "Handle forced batch process response with successful transactions",
			request: state.ProcessRequest{
				Transactions: tx1Plustx2,
				BatchNumber:  1,
				Coinbase:     seqAddr,
				Timestamp_V1: now(),
				OldStateRoot: oldHash,
			},
			result:       successfulBatchResp,
			oldStateRoot: oldHash,
			expectedStoredTxs: []transactionToStore{
				{
					hash:          signedTx1.Hash(),
					from:          auth.From,
					batchNumber:   1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  oldHash,
					response:      txResponseOne,
					isForcedBatch: true,
					batchResponse: successfulBatchResp,
				},
				{
					hash:          signedTx2.Hash(),
					from:          auth.From,
					batchNumber:   1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  newHash,
					response:      txResponseTwo,
					isForcedBatch: true,
					batchResponse: successfulBatchResp,
				},
			},
		},
		{
			name: "Handle forced batch process response with reverted transactions",
			request: state.ProcessRequest{
				BatchNumber:  1,
				Coinbase:     seqAddr,
				Timestamp_V1: now(),
				OldStateRoot: oldHash,
			},
			result:       revertedBatchResp,
			oldStateRoot: oldHash,
			expectedStoredTxs: []transactionToStore{
				{
					hash:          signedTx1.Hash(),
					from:          auth.From,
					batchNumber:   1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  oldHash,
					response:      txResponseReverted,
					isForcedBatch: true,
					batchResponse: revertedBatchResp,
				}},
		},
		{
			name: "Handle forced batch process response with intrinsic ROM err",
			request: state.ProcessRequest{
				BatchNumber:  1,
				Coinbase:     seqAddr,
				Timestamp_V1: now(),
				OldStateRoot: oldHash,
			},

			result:       intrinsicErrBatchResp,
			oldStateRoot: oldHash,
			expectedStoredTxs: []transactionToStore{
				{
					hash:          signedTx1.Hash(),
					from:          auth.From,
					batchNumber:   1,
					coinbase:      seqAddr,
					timestamp:     now(),
					oldStateRoot:  oldHash,
					response:      txResponseOne,
					isForcedBatch: true,
					batchResponse: intrinsicErrBatchResp,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storedTxs := make([]transactionToStore, 0)
			f.pendingL2BlocksToStore = make(chan transactionToStore)

			// Mock storeProcessedTx to store txs into the storedTxs slice
			go func() {
				for tx := range f.pendingL2BlocksToStore {
					storedTxs = append(storedTxs, tx)
					f.pendingL2BlocksToStoreWG.Done()
				}
			}()

			workerMock.On("AddPendingTxToStore", mock.Anything, mock.Anything).Return()
			workerMock.On("DeleteForcedTx", mock.Anything, mock.Anything).Return()
			workerMock.On("AddForcedTx", mock.Anything, mock.Anything).Return()

			f.handleProcessForcedTxsResponse(ctx, tc.request, tc.result, tc.oldStateRoot)

			f.pendingL2BlocksToStoreWG.Wait()
			require.Nil(t, err)
			require.Equal(t, len(tc.expectedStoredTxs), len(storedTxs))
			for i := 0; i < len(tc.expectedStoredTxs); i++ {
				expectedTx := tc.expectedStoredTxs[i]
				actualTx := storedTxs[i]
				require.Equal(t, expectedTx, actualTx)
			}
		})
	}
}*/

/*func TestFinalizer_storeProcessedTx(t *testing.T) {
	f = setupFinalizer(false)
	testCases := []struct {
		name              string
		batchNum          uint64
		coinbase          common.Address
		timestamp         time.Time
		previousStateRoot common.Hash
		txResponse        *state.ProcessTransactionResponse
		isForcedBatch     bool
		expectedTxToStore transactionToStore
	}{
		{
			name:              "Normal transaction",
			batchNum:          1,
			coinbase:          seqAddr,
			timestamp:         time.Now(),
			previousStateRoot: oldHash,
			txResponse: &state.ProcessTransactionResponse{
				TxHash: txHash,
			},
			isForcedBatch: false,
			expectedTxToStore: transactionToStore{
				batchNumber:  1,
				coinbase:     seqAddr,
				timestamp:    now(),
				oldStateRoot: oldHash,
				response: &state.ProcessTransactionResponse{
					TxHash: txHash,
				},
				isForcedBatch: false,
			},
		},
		{
			name:              "Forced transaction",
			batchNum:          1,
			coinbase:          seqAddr,
			timestamp:         time.Now(),
			previousStateRoot: oldHash,
			txResponse: &state.ProcessTransactionResponse{
				TxHash: txHash2,
			},
			isForcedBatch: true,
			expectedTxToStore: transactionToStore{
				batchNumber:  1,
				coinbase:     seqAddr,
				timestamp:    now(),
				oldStateRoot: oldHash,
				response: &state.ProcessTransactionResponse{
					TxHash: txHash2,
				},
				isForcedBatch: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			stateMock.On("StoreProcessedTxAndDeleteFromPool", ctx, tc.expectedTxToStore).Return(nilErr)

			// act
			f.storeProcessedTx(ctx, tc.expectedTxToStore)

			// assert
			stateMock.AssertExpectations(t)
		})
	}
}*/

func TestFinalizer_updateWorkerAfterSuccessfulProcessing(t *testing.T) {
	testCases := []struct {
		name                  string
		txTracker             *TxTracker
		processBatchResponse  *state.ProcessBatchResponse
		expectedDeleteTxCount int
		expectedUpdateCount   int
	}{
		{
			name: "Successful update with one read-write address",
			txTracker: &TxTracker{
				Hash:  oldHash,
				From:  senderAddr,
				Nonce: nonce1,
			},
			processBatchResponse: &state.ProcessBatchResponse{
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					senderAddr: {
						Address: senderAddr,
						Nonce:   &nonce1,
					},
				},
			},
			expectedDeleteTxCount: 1,
			expectedUpdateCount:   1,
		},
		{
			name: "Successful update with multiple read-write addresses",
			txTracker: &TxTracker{
				Hash:  oldHash,
				From:  senderAddr,
				Nonce: 1,
			},
			processBatchResponse: &state.ProcessBatchResponse{
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{
					senderAddr: {
						Address: senderAddr,
						Nonce:   &nonce1,
					},
					receiverAddr: {
						Address: receiverAddr,
						Nonce:   &nonce2,
					},
				},
			},
			expectedDeleteTxCount: 1,
			expectedUpdateCount:   2,
		},
		{
			name: "No update when no read-write addresses provided",
			txTracker: &TxTracker{
				Hash:  oldHash,
				From:  senderAddr,
				Nonce: 1,
			},
			processBatchResponse: &state.ProcessBatchResponse{
				ReadWriteAddresses: map[common.Address]*state.InfoReadWrite{},
			},
			expectedDeleteTxCount: 1,
			expectedUpdateCount:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			finalizerInstance := setupFinalizer(false)
			workerMock.On("MoveTxPendingToStore", tc.txTracker.Hash, tc.txTracker.From).Times(tc.expectedDeleteTxCount)
			txsToDelete := make([]*TxTracker, 0, len(tc.processBatchResponse.ReadWriteAddresses))
			for _, infoReadWrite := range tc.processBatchResponse.ReadWriteAddresses {
				txsToDelete = append(txsToDelete, &TxTracker{
					Hash:         oldHash,
					From:         infoReadWrite.Address,
					FailedReason: &testErrStr,
				})
			}
			workerMock.On("UpdateAfterSingleSuccessfulTxExecution", tc.txTracker.From, tc.processBatchResponse.ReadWriteAddresses).
				Return(txsToDelete)
			if tc.expectedUpdateCount > 0 {
				poolMock.On("UpdateTxStatus", mock.Anything, mock.Anything, pool.TxStatusFailed, false, mock.Anything).Times(tc.expectedUpdateCount).Return(nil)
			}

			// act
			finalizerInstance.updateWorkerAfterSuccessfulProcessing(ctx, tc.txTracker.Hash, tc.txTracker.From, false, tc.processBatchResponse)

			// assert
			workerMock.AssertExpectations(t)
			stateMock.AssertExpectations(t)
		})
	}
}

/*func TestFinalizer_reprocessFullBatch(t *testing.T) {
	successfulResult := &state.ProcessBatchResponse{
		NewStateRoot: newHash,
	}
	roomOOCErrResult := &state.ProcessBatchResponse{
		NewStateRoot:  newHash,
		IsRomOOCError: true,
	}
	testCases := []struct {
		name                     string
		batchNum                 uint64
		oldStateRoot             common.Hash
		mockGetBatchByNumber     *state.Batch
		mockGetBatchByNumberErr  error
		expectedExecutorResponse *state.ProcessBatchResponse
		expectedResult           *state.ProcessBatchResponse
		expectedDecodeErr        error
		expectedExecutorErr      error
		expectedError            error
	}{
		{
			name:     "Success",
			batchNum: 1,
			mockGetBatchByNumber: &state.Batch{
				BatchNumber:    1,
				BatchL2Data:    decodedBatchL2Data,
				GlobalExitRoot: oldHash,
				Coinbase:       common.Address{},
				Timestamp:      time.Now(),
			},
			expectedExecutorResponse: successfulResult,
			expectedResult:           successfulResult,
		},
		{
			name:                    "Error while getting batch by number",
			batchNum:                1,
			mockGetBatchByNumberErr: errors.New("database err"),
			expectedError:           ErrGetBatchByNumber,
		},
		{
			name:     "Error processing batch",
			batchNum: 1,
			mockGetBatchByNumber: &state.Batch{
				BatchNumber:    1,
				BatchL2Data:    decodedBatchL2Data,
				GlobalExitRoot: oldHash,
				Coinbase:       common.Address{},
				Timestamp:      time.Now(),
			},
			expectedExecutorErr: ErrProcessBatch,
			expectedError:       ErrProcessBatch,
		},
		{
			name:     "RomOOCError",
			batchNum: 1,
			mockGetBatchByNumber: &state.Batch{
				BatchNumber:    1,
				BatchL2Data:    decodedBatchL2Data,
				GlobalExitRoot: oldHash,
				Coinbase:       common.Address{},
				Timestamp:      time.Now(),
			},
			expectedExecutorResponse: roomOOCErrResult,
			expectedError:            ErrProcessBatchOOC,
		},
		{
			name:     "Reprocessed batch has different state root",
			batchNum: 1,
			mockGetBatchByNumber: &state.Batch{
				BatchNumber:    1,
				BatchL2Data:    decodedBatchL2Data,
				GlobalExitRoot: oldHash,
				Coinbase:       common.Address{},
				Timestamp:      time.Now(),
			},
			expectedExecutorResponse: &state.ProcessBatchResponse{
				NewStateRoot: newHash2,
			},
			expectedError: ErrStateRootNoMatch,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f := setupFinalizer(true)
			stateMock.On("GetBatchByNumber", context.Background(), tc.batchNum, nil).Return(tc.mockGetBatchByNumber, tc.mockGetBatchByNumberErr).Once()
			//			if tc.name != "Error while getting batch by number" {
			//			stateMock.On("GetForkIDByBatchNumber", f.wipBatch.batchNumber).Return(uint64(7)).Once()
			//		}
			if tc.mockGetBatchByNumberErr == nil && tc.expectedDecodeErr == nil {
				stateMock.On("ProcessBatchV2", context.Background(), mock.Anything, false).Return(tc.expectedExecutorResponse, tc.expectedExecutorErr)
			}

			// act
			result, err := f.batchSanityCheck(context.Background(), tc.batchNum, f.wipBatch.initialStateRoot, newHash)

			// assert
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
			stateMock.AssertExpectations(t)
			stateMock.AssertExpectations(t)
		})
	}
}*/

func TestFinalizer_isBatchAlmostFull(t *testing.T) {
	// arrange
	testCases := []struct {
		name               string
		modifyResourceFunc func(resources state.BatchResources) state.BatchResources
		expectedResult     bool
	}{
		{
			name: "Is ready - MaxBatchBytesSize",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.Bytes = f.getConstraintThresholdUint64(bc.MaxBatchBytesSize) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxBatchBytesSize",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.Bytes = f.getConstraintThresholdUint64(bc.MaxBatchBytesSize) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxCumulativeGasUsed",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.GasUsed = f.getConstraintThresholdUint64(bc.MaxCumulativeGasUsed) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxCumulativeGasUsed",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.GasUsed = f.getConstraintThresholdUint64(bc.MaxCumulativeGasUsed) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxSteps",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.Steps = f.getConstraintThresholdUint32(bc.MaxSteps) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxSteps",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.Steps = f.getConstraintThresholdUint32(bc.MaxSteps) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxPoseidonPaddings",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.PoseidonPaddings = f.getConstraintThresholdUint32(bc.MaxPoseidonPaddings) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxPoseidonPaddings",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.PoseidonPaddings = f.getConstraintThresholdUint32(bc.MaxPoseidonPaddings) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxBinaries",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.Binaries = f.getConstraintThresholdUint32(bc.MaxBinaries) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxBinaries",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.Binaries = f.getConstraintThresholdUint32(bc.MaxBinaries) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxKeccakHashes",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.KeccakHashes = f.getConstraintThresholdUint32(bc.MaxKeccakHashes) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxKeccakHashes",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.KeccakHashes = f.getConstraintThresholdUint32(bc.MaxKeccakHashes) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxArithmetics",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.Arithmetics = f.getConstraintThresholdUint32(bc.MaxArithmetics) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxArithmetics",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.Arithmetics = f.getConstraintThresholdUint32(bc.MaxArithmetics) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxMemAligns",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.MemAligns = f.getConstraintThresholdUint32(bc.MaxMemAligns) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxMemAligns",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.MemAligns = f.getConstraintThresholdUint32(bc.MaxMemAligns) + 1
				return resources
			},
			expectedResult: false,
		},
		{
			name: "Is ready - MaxSHA256Hashes",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.Sha256Hashes_V2 = f.getConstraintThresholdUint32(bc.MaxSHA256Hashes) - 1
				return resources
			},
			expectedResult: true,
		},
		{
			name: "Is NOT ready - MaxSHA256Hashes",
			modifyResourceFunc: func(resources state.BatchResources) state.BatchResources {
				resources.ZKCounters.Sha256Hashes_V2 = f.getConstraintThresholdUint32(bc.MaxSHA256Hashes) + 1
				return resources
			},
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			f = setupFinalizer(true)
			maxRemainingResource := getMaxBatchResources(bc)
			f.wipBatch.imRemainingResources = tc.modifyResourceFunc(maxRemainingResource)

			// act
			result, closeReason := f.checkIfFinalizeBatch()

			// assert
			assert.Equal(t, tc.expectedResult, result)
			if tc.expectedResult {
				assert.Equal(t, state.ResourceMarginExhaustedClosingReason, closeReason)
			} else {
				assert.Equal(t, state.EmptyClosingReason, closeReason)
			}
		})
	}
}

func TestFinalizer_setNextForcedBatchDeadline(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	now = testNow
	defer func() {
		now = time.Now
	}()
	expected := now().Unix() + int64(f.cfg.ForcedBatchesTimeout.Duration.Seconds())

	// act
	f.setNextForcedBatchDeadline()

	// assert
	assert.Equal(t, expected, f.nextForcedBatchDeadline)
}

func TestFinalizer_getConstraintThresholdUint64(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	input := uint64(100)
	expect := input * uint64(f.cfg.ResourceExhaustedMarginPct) / 100

	// act
	result := f.getConstraintThresholdUint64(input)

	// assert
	assert.Equal(t, result, expect)
}

func TestFinalizer_getConstraintThresholdUint32(t *testing.T) {
	// arrange
	f = setupFinalizer(false)
	input := uint32(100)
	expect := input * f.cfg.ResourceExhaustedMarginPct / 100

	// act
	result := f.getConstraintThresholdUint32(input)

	// assert
	assert.Equal(t, result, expect)
}

func TestFinalizer_getRemainingResources(t *testing.T) {
	// act
	remainingResources := getMaxBatchResources(bc)

	// assert
	assert.Equal(t, remainingResources.ZKCounters.GasUsed, bc.MaxCumulativeGasUsed)
	assert.Equal(t, remainingResources.ZKCounters.KeccakHashes, bc.MaxKeccakHashes)
	assert.Equal(t, remainingResources.ZKCounters.PoseidonHashes, bc.MaxPoseidonHashes)
	assert.Equal(t, remainingResources.ZKCounters.PoseidonPaddings, bc.MaxPoseidonPaddings)
	assert.Equal(t, remainingResources.ZKCounters.MemAligns, bc.MaxMemAligns)
	assert.Equal(t, remainingResources.ZKCounters.Arithmetics, bc.MaxArithmetics)
	assert.Equal(t, remainingResources.ZKCounters.Binaries, bc.MaxBinaries)
	assert.Equal(t, remainingResources.ZKCounters.Steps, bc.MaxSteps)
	assert.Equal(t, remainingResources.ZKCounters.Sha256Hashes_V2, bc.MaxSHA256Hashes)
	assert.Equal(t, remainingResources.Bytes, bc.MaxBatchBytesSize)
}

func Test_isBatchFull(t *testing.T) {
	f = setupFinalizer(true)

	testCases := []struct {
		name            string
		batchCountOfTxs int
		maxTxsPerBatch  uint64
		expected        bool
	}{
		{
			name:            "Batch is not full",
			batchCountOfTxs: 5,
			maxTxsPerBatch:  10,
			expected:        false,
		},
		{
			name:            "Batch is full",
			batchCountOfTxs: 10,
			maxTxsPerBatch:  10,
			expected:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f.wipBatch.countOfTxs = tc.batchCountOfTxs
			f.batchConstraints.MaxTxsPerBatch = tc.maxTxsPerBatch

			assert.Equal(t, tc.expected, f.maxTxsPerBatchReached(f.wipBatch))
		})
	}
}

func Test_sortForcedBatches(t *testing.T) {
	f = setupFinalizer(false)

	testCases := []struct {
		name     string
		input    []state.ForcedBatch
		expected []state.ForcedBatch
	}{
		{
			name:     "Empty slice",
			input:    []state.ForcedBatch{},
			expected: []state.ForcedBatch{},
		},
		{
			name:     "Single item slice",
			input:    []state.ForcedBatch{{ForcedBatchNumber: 5}},
			expected: []state.ForcedBatch{{ForcedBatchNumber: 5}},
		},
		{
			name:     "Multiple items unsorted",
			input:    []state.ForcedBatch{{ForcedBatchNumber: 5}, {ForcedBatchNumber: 3}, {ForcedBatchNumber: 9}},
			expected: []state.ForcedBatch{{ForcedBatchNumber: 3}, {ForcedBatchNumber: 5}, {ForcedBatchNumber: 9}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := f.sortForcedBatches(testCase.input)
			assert.Equal(t, testCase.expected, result, "They should be equal")
		})
	}
}

func setupFinalizer(withWipBatch bool) *finalizer {
	wipBatch := new(Batch)
	poolMock = new(PoolMock)
	stateMock = new(StateMock)
	workerMock = new(WorkerMock)
	dbTxMock = new(DbTxMock)
	if withWipBatch {
		decodedBatchL2Data, err = hex.DecodeHex(testBatchL2DataAsString)
		if err != nil {
			panic(err)
		}
		wipBatch = &Batch{
			batchNumber:          1,
			coinbase:             seqAddr,
			initialStateRoot:     oldHash,
			imStateRoot:          newHash,
			timestamp:            now(),
			imRemainingResources: getMaxBatchResources(bc),
			closingReason:        state.EmptyClosingReason,
		}
	}
	eventStorage, err := nileventstorage.NewNilEventStorage()
	if err != nil {
		panic(err)
	}
	eventLog := event.NewEventLog(event.Config{}, eventStorage)
	return &finalizer{
		cfg:                        cfg,
		isSynced:                   isSynced,
		sequencerAddress:           seqAddr,
		workerIntf:                 workerMock,
		poolIntf:                   poolMock,
		stateIntf:                  stateMock,
		wipBatch:                   wipBatch,
		sipBatch:                   wipBatch,
		batchConstraints:           bc,
		nextForcedBatches:          make([]state.ForcedBatch, 0),
		nextForcedBatchDeadline:    0,
		nextForcedBatchesMux:       new(sync.Mutex),
		effectiveGasPrice:          pool.NewEffectiveGasPrice(poolCfg.EffectiveGasPrice),
		eventLog:                   eventLog,
		pendingL2BlocksToProcess:   make(chan *L2Block, pendingL2BlocksBufferSize),
		pendingL2BlocksToProcessWG: new(WaitGroupCount),
		pendingL2BlocksToStore:     make(chan *L2Block, pendingL2BlocksBufferSize),
		pendingL2BlocksToStoreWG:   new(WaitGroupCount),
		storedFlushID:              0,
		storedFlushIDCond:          sync.NewCond(new(sync.Mutex)),
		proverID:                   "",
		lastPendingFlushID:         0,
		pendingFlushIDCond:         sync.NewCond(new(sync.Mutex)),
	}
}
