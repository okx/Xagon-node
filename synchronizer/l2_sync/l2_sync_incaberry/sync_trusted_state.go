package l2_sync_incaberry

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/metrics"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type zkEVMClientInterface interface {
	BatchNumber(ctx context.Context) (uint64, error)
	BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error)
}

// TrustedState contains the last trusted batches and stateRoot (cache)
type TrustedState struct {
	LastTrustedBatches []*state.Batch
	LastStateRoot      *common.Hash
}

type syncTrustedBatchesStateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ProcessBatch(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, egpLog *state.EffectiveGasPriceLog, globalExitRoot, blockInfoRoot common.Hash, dbTx pgx.Tx) (*state.L2Header, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error
}
type syncTrustedBatchesSynchronizerInterface interface {
	PendingFlushID(flushID uint64, proverID string)
	CheckFlushID(dbTx pgx.Tx) error
}

// SyncTrustedBatchesAction is the action that synchronizes the trusted state
type SyncTrustedBatchesAction struct {
	zkEVMClient  zkEVMClientInterface
	state        syncTrustedBatchesStateInterface
	sync         syncTrustedBatchesSynchronizerInterface
	TrustedState TrustedState
}

// CleanTrustedState Clean cache of TrustedBatches and StateRoot
func (s *SyncTrustedBatchesAction) CleanTrustedState() {
	s.TrustedState.LastTrustedBatches = nil
	s.TrustedState.LastStateRoot = nil
}

// NewSyncTrustedStateExecutor creates a new syncTrustedBatchesAction for incaberry
func NewSyncTrustedStateExecutor(zkEVMClient zkEVMClientInterface, state syncTrustedBatchesStateInterface, sync syncTrustedBatchesSynchronizerInterface) *SyncTrustedBatchesAction {
	return &SyncTrustedBatchesAction{
		zkEVMClient:  zkEVMClient,
		state:        state,
		sync:         sync,
		TrustedState: TrustedState{},
	}
}

// GetCachedBatch implements syncinterfaces.SyncTrustedStateExecutor. Returns a cached batch
func (s *SyncTrustedBatchesAction) GetCachedBatch(batchNumber uint64) *state.Batch {
	if s.TrustedState.LastTrustedBatches == nil {
		return nil
	}
	for _, batch := range s.TrustedState.LastTrustedBatches {
		if batch.BatchNumber == batchNumber {
			return batch
		}
	}
	return nil
}

// SyncTrustedState synchronizes information from the trusted sequencer
// related to the trusted state when the node has all the information from
// l1 synchronized
func (s *SyncTrustedBatchesAction) SyncTrustedState(ctx context.Context, latestSyncedBatch uint64, maximumBatchNumberToProcess uint64) error {
	log.Info("syncTrustedState: Getting trusted state info")
	start := time.Now()
	lastTrustedStateBatchNumberSeen, err := s.zkEVMClient.BatchNumber(ctx)
	metrics.GetTrustedBatchNumberTime(time.Since(start))
	if err != nil {
		log.Warn("syncTrustedState: error syncing trusted state. Error: ", err)
		return err
	}
	lastTrustedStateBatchNumber := min(lastTrustedStateBatchNumberSeen, maximumBatchNumberToProcess)
	log.Debug("syncTrustedState: lastTrustedStateBatchNumber ", lastTrustedStateBatchNumber)
	log.Debug("syncTrustedState: latestSyncedBatch ", latestSyncedBatch)
	log.Debug("syncTrustedState: lastTrustedStateBatchNumberSeen ", lastTrustedStateBatchNumberSeen)
	if lastTrustedStateBatchNumber < latestSyncedBatch {
		return nil
	}

	batchNumberToSync := latestSyncedBatch
	for batchNumberToSync <= lastTrustedStateBatchNumber {
		if batchNumberToSync == 0 {
			batchNumberToSync++
			continue
		}
		start = time.Now()
		batchToSync, err := s.zkEVMClient.BatchByNumber(ctx, big.NewInt(0).SetUint64(batchNumberToSync))
		metrics.GetTrustedBatchInfoTime(time.Since(start))
		if err != nil {
			log.Warnf("syncTrustedState: failed to get batch %d from trusted state. Error: %v", batchNumberToSync, err)
			return err
		}

		dbTx, err := s.state.BeginStateTransaction(ctx)
		if err != nil {
			log.Errorf("syncTrustedState: error creating db transaction to sync trusted batch %d: %v", batchNumberToSync, err)
			return err
		}
		start = time.Now()
		cbatches, lastStateRoot, err := s.processTrustedBatch(ctx, batchToSync, dbTx)
		metrics.ProcessTrustedBatchTime(time.Since(start))
		if err != nil {
			log.Errorf("syncTrustedState: error processing trusted batch %d: %v", batchNumberToSync, err)
			rollbackErr := dbTx.Rollback(ctx)
			if rollbackErr != nil {
				log.Errorf("syncTrustedState: error rolling back db transaction to sync trusted batch %d: %v", batchNumberToSync, rollbackErr)
				return rollbackErr
			}
			return err
		}
		log.Debug("syncTrustedState: Checking FlushID to commit trustedState data to db")
		err = s.sync.CheckFlushID(dbTx)
		if err != nil {
			log.Errorf("syncTrustedState: error checking flushID. Error: %v", err)
			rollbackErr := dbTx.Rollback(ctx)
			if rollbackErr != nil {
				log.Errorf("syncTrustedState: error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
		}

		if err := dbTx.Commit(ctx); err != nil {
			log.Errorf("syncTrustedState: error committing db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
			return err
		}
		s.TrustedState.LastTrustedBatches = cbatches
		s.TrustedState.LastStateRoot = lastStateRoot
		batchNumberToSync++
	}

	log.Info("syncTrustedState: Trusted state fully synchronized")
	return nil
}

func (s *SyncTrustedBatchesAction) processTrustedBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, *common.Hash, error) {
	log.Debugf("Processing trusted batch: %d", uint64(trustedBatch.Number))
	trustedBatchL2Data := trustedBatch.BatchL2Data
	batches := s.TrustedState.LastTrustedBatches
	log.Debug("len(batches): ", len(batches))
	batches, err := s.getCurrentBatches(ctx, batches, trustedBatch, dbTx)
	if err != nil {
		log.Error("error getting currentBatches. Error: ", err)
		return nil, nil, err
	}

	if batches[0] != nil && (((trustedBatch.StateRoot == common.Hash{}) && (batches[0].StateRoot != common.Hash{})) ||
		len(batches[0].BatchL2Data) > len(trustedBatchL2Data)) {
		log.Error("error: inconsistency in data received from trustedNode")
		log.Infof("BatchNumber. stored: %d. synced: %d", batches[0].BatchNumber, uint64(trustedBatch.Number))
		log.Infof("GlobalExitRoot. stored:  %s. synced: %s", batches[0].GlobalExitRoot.String(), trustedBatch.GlobalExitRoot.String())
		log.Infof("LocalExitRoot. stored:  %s. synced: %s", batches[0].LocalExitRoot.String(), trustedBatch.LocalExitRoot.String())
		log.Infof("StateRoot. stored:  %s. synced: %s", batches[0].StateRoot.String(), trustedBatch.StateRoot.String())
		log.Infof("Coinbase. stored:  %s. synced: %s", batches[0].Coinbase.String(), trustedBatch.Coinbase.String())
		log.Infof("Timestamp. stored:  %d. synced: %d", uint64(batches[0].Timestamp.Unix()), uint64(trustedBatch.Timestamp))
		log.Infof("BatchL2Data. stored: %s. synced: %s", common.Bytes2Hex(batches[0].BatchL2Data), common.Bytes2Hex(trustedBatchL2Data))
		return nil, nil, fmt.Errorf("error: inconsistency in data received from trustedNode")
	}

	if s.TrustedState.LastStateRoot == nil && (batches[0] == nil || (batches[0].StateRoot == common.Hash{})) {
		log.Debug("Setting stateRoot of previous batch. StateRoot: ", batches[1].StateRoot)
		// Previous synchronization incomplete. Needs to reprocess all txs again
		s.TrustedState.LastStateRoot = &batches[1].StateRoot
	} else if batches[0] != nil && (batches[0].StateRoot != common.Hash{}) {
		// Previous synchronization completed
		s.TrustedState.LastStateRoot = &batches[0].StateRoot
	}

	request := state.ProcessRequest{
		BatchNumber:     uint64(trustedBatch.Number),
		OldStateRoot:    *s.TrustedState.LastStateRoot,
		OldAccInputHash: batches[1].AccInputHash,
		Coinbase:        common.HexToAddress(trustedBatch.Coinbase.String()),
		Timestamp_V1:    time.Unix(int64(trustedBatch.Timestamp), 0),
	}
	// check if batch needs to be synchronized
	if batches[0] != nil {
		if checkIfSynced(batches, trustedBatch) {
			log.Debugf("Batch %d already synchronized", uint64(trustedBatch.Number))
			return batches, s.TrustedState.LastStateRoot, nil
		}
		log.Infof("Batch %d needs to be updated", uint64(trustedBatch.Number))

		// Find txs to be processed and included in the trusted state
		if *s.TrustedState.LastStateRoot == batches[1].StateRoot {
			prevBatch := uint64(trustedBatch.Number) - 1
			log.Infof("ResetTrustedState: processTrustedBatch: %d Cleaning state until batch:%d  ", trustedBatch.Number, prevBatch)
			// Delete txs that were stored before restart. We need to reprocess all txs because the intermediary stateRoot is only stored in memory
			err := s.state.ResetTrustedState(ctx, prevBatch, dbTx)
			if err != nil {
				log.Error("error resetting trusted state. Error: ", err)
				return nil, nil, err
			}
			// All txs need to be processed
			request.Transactions = trustedBatchL2Data
			// Reopen batch
			err = s.openBatch(ctx, trustedBatch, dbTx)
			if err != nil {
				log.Error("error openning batch. Error: ", err)
				return nil, nil, err
			}
			request.GlobalExitRoot_V1 = trustedBatch.GlobalExitRoot
			request.Transactions = trustedBatchL2Data
		} else {
			// Only new txs need to be processed
			storedTxs, syncedTxs, _, syncedEfficiencyPercentages, err := s.decodeTxs(trustedBatchL2Data, batches)
			if err != nil {
				return nil, nil, err
			}
			if len(storedTxs) < len(syncedTxs) {
				forkID := s.state.GetForkIDByBatchNumber(batches[0].BatchNumber)
				txsToBeAdded := syncedTxs[len(storedTxs):]
				if forkID >= state.FORKID_DRAGONFRUIT {
					syncedEfficiencyPercentages = syncedEfficiencyPercentages[len(storedTxs):]
				}
				log.Infof("Processing %d new txs with forkID: %d", len(txsToBeAdded), forkID)

				request.Transactions, err = state.EncodeTransactions(txsToBeAdded, syncedEfficiencyPercentages, forkID)
				if err != nil {
					log.Error("error encoding txs (%d) to be added to the state. Error: %v", len(txsToBeAdded), err)
					return nil, nil, err
				}
				log.Debug("request.Transactions: ", common.Bytes2Hex(request.Transactions))
			} else {
				log.Info("Nothing to sync. Node updated. Checking if it is closed")
				isBatchClosed := trustedBatch.StateRoot.String() != state.ZeroHash.String()
				if isBatchClosed {
					//Sanity check
					if s.TrustedState.LastStateRoot != nil && trustedBatch.StateRoot != *s.TrustedState.LastStateRoot {
						log.Errorf("batch %d, different batchL2Datas (trustedBatchL2Data: %s, batches[0].BatchL2Data: %s). Decoded txs are len(storedTxs): %d, len(syncedTxs): %d", uint64(trustedBatch.Number), trustedBatchL2Data.Hex(), "0x"+common.Bytes2Hex(batches[0].BatchL2Data), len(storedTxs), len(syncedTxs))
						for _, tx := range storedTxs {
							log.Error("stored txHash : ", tx.Hash())
						}
						for _, tx := range syncedTxs {
							log.Error("synced txHash : ", tx.Hash())
						}
						log.Errorf("batch: %d, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), *s.TrustedState.LastStateRoot, trustedBatch.StateRoot)
						return nil, nil, fmt.Errorf("batch: %d, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), *s.TrustedState.LastStateRoot, trustedBatch.StateRoot)
					}
					receipt := state.ProcessingReceipt{
						BatchNumber:   uint64(trustedBatch.Number),
						StateRoot:     trustedBatch.StateRoot,
						LocalExitRoot: trustedBatch.LocalExitRoot,
						BatchL2Data:   trustedBatchL2Data,
						AccInputHash:  trustedBatch.AccInputHash,
					}
					log.Debugf("closing batch %d", uint64(trustedBatch.Number))
					if err := s.state.CloseBatch(ctx, receipt, dbTx); err != nil {
						// This is a workaround to avoid closing a batch that was already closed
						if err.Error() != state.ErrBatchAlreadyClosed.Error() {
							log.Errorf("error closing batch %d", uint64(trustedBatch.Number))
							return nil, nil, err
						} else {
							log.Warnf("CASE 02: the batch [%d] was already closed", uint64(trustedBatch.Number))
							log.Info("batches[0].BatchNumber: ", batches[0].BatchNumber)
							log.Info("batches[0].AccInputHash: ", batches[0].AccInputHash)
							log.Info("batches[0].StateRoot: ", batches[0].StateRoot)
							log.Info("batches[0].LocalExitRoot: ", batches[0].LocalExitRoot)
							log.Info("batches[0].GlobalExitRoot: ", batches[0].GlobalExitRoot)
							log.Info("batches[0].Coinbase: ", batches[0].Coinbase)
							log.Info("batches[0].ForcedBatchNum: ", batches[0].ForcedBatchNum)
							log.Info("####################################")
							log.Info("batches[1].BatchNumber: ", batches[1].BatchNumber)
							log.Info("batches[1].AccInputHash: ", batches[1].AccInputHash)
							log.Info("batches[1].StateRoot: ", batches[1].StateRoot)
							log.Info("batches[1].LocalExitRoot: ", batches[1].LocalExitRoot)
							log.Info("batches[1].GlobalExitRoot: ", batches[1].GlobalExitRoot)
							log.Info("batches[1].Coinbase: ", batches[1].Coinbase)
							log.Info("batches[1].ForcedBatchNum: ", batches[1].ForcedBatchNum)
							log.Info("###############################")
							log.Info("trustedBatch.BatchNumber: ", trustedBatch.Number)
							log.Info("trustedBatch.AccInputHash: ", trustedBatch.AccInputHash)
							log.Info("trustedBatch.StateRoot: ", trustedBatch.StateRoot)
							log.Info("trustedBatch.LocalExitRoot: ", trustedBatch.LocalExitRoot)
							log.Info("trustedBatch.GlobalExitRoot: ", trustedBatch.GlobalExitRoot)
							log.Info("trustedBatch.Coinbase: ", trustedBatch.Coinbase)
							log.Info("trustedBatch.ForcedBatchNum: ", trustedBatch.ForcedBatchNumber)
						}
					}
					batches[0].AccInputHash = trustedBatch.AccInputHash
					batches[0].StateRoot = trustedBatch.StateRoot
					batches[0].LocalExitRoot = trustedBatch.LocalExitRoot
				}
				return batches, &trustedBatch.StateRoot, nil
			}
		}
		// Update batchL2Data
		err := s.state.UpdateBatchL2Data(ctx, batches[0].BatchNumber, trustedBatchL2Data, dbTx)
		if err != nil {
			log.Errorf("error opening batch %d", uint64(trustedBatch.Number))
			return nil, nil, err
		}
		batches[0].BatchL2Data = trustedBatchL2Data
		log.Debug("BatchL2Data updated for batch: ", batches[0].BatchNumber)
	} else {
		log.Infof("Batch %d needs to be synchronized", uint64(trustedBatch.Number))
		err := s.openBatch(ctx, trustedBatch, dbTx)
		if err != nil {
			log.Error("error openning batch. Error: ", err)
			return nil, nil, err
		}
		request.GlobalExitRoot_V1 = trustedBatch.GlobalExitRoot
		request.Transactions = trustedBatchL2Data
	}

	log.Debugf("Processing sequencer for batch %d", uint64(trustedBatch.Number))

	processBatchResp, err := s.processAndStoreTxs(ctx, trustedBatch, request, dbTx)
	if err != nil {
		log.Error("error procesingAndStoringTxs. Error: ", err)
		return nil, nil, err
	}

	log.Debug("TrustedBatch.StateRoot ", trustedBatch.StateRoot)
	isBatchClosed := trustedBatch.StateRoot.String() != state.ZeroHash.String()
	if isBatchClosed {
		//Sanity check
		if trustedBatch.StateRoot != processBatchResp.NewStateRoot {
			log.Error("trustedBatchL2Data: ", trustedBatchL2Data)
			log.Error("request.Transactions: ", request.Transactions)
			log.Errorf("batch: %d after processing some txs, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), processBatchResp.NewStateRoot.String(), trustedBatch.StateRoot.String())
			return nil, nil, fmt.Errorf("batch: %d, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), processBatchResp.NewStateRoot.String(), trustedBatch.StateRoot.String())
		}
		receipt := state.ProcessingReceipt{
			BatchNumber:   uint64(trustedBatch.Number),
			StateRoot:     processBatchResp.NewStateRoot,
			LocalExitRoot: processBatchResp.NewLocalExitRoot,
			BatchL2Data:   trustedBatchL2Data,
			AccInputHash:  trustedBatch.AccInputHash,
		}

		log.Debugf("closing batch %d", uint64(trustedBatch.Number))
		if err := s.state.CloseBatch(ctx, receipt, dbTx); err != nil {
			// This is a workarround to avoid closing a batch that was already closed
			if err.Error() != state.ErrBatchAlreadyClosed.Error() {
				log.Errorf("error closing batch %d", uint64(trustedBatch.Number))
				return nil, nil, err
			} else {
				log.Warnf("CASE 01: batch [%d] was already closed", uint64(trustedBatch.Number))
			}
		}
		log.Info("Batch closed right after processing some tx")
		if batches[0] != nil {
			log.Debug("Updating batches[0] values...")
			batches[0].AccInputHash = trustedBatch.AccInputHash
			batches[0].StateRoot = trustedBatch.StateRoot
			batches[0].LocalExitRoot = trustedBatch.LocalExitRoot
			batches[0].BatchL2Data = trustedBatchL2Data
		}
	}

	log.Infof("Batch %d synchronized", uint64(trustedBatch.Number))
	return batches, &processBatchResp.NewStateRoot, nil
}

func (s *SyncTrustedBatchesAction) openBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx) error {
	log.Debugf("Opening batch %d", trustedBatch.Number)
	var batchL2Data []byte = trustedBatch.BatchL2Data
	processCtx := state.ProcessingContext{
		BatchNumber:    uint64(trustedBatch.Number),
		Coinbase:       common.HexToAddress(trustedBatch.Coinbase.String()),
		Timestamp:      time.Unix(int64(trustedBatch.Timestamp), 0),
		GlobalExitRoot: trustedBatch.GlobalExitRoot,
		BatchL2Data:    &batchL2Data,
	}
	if trustedBatch.ForcedBatchNumber != nil {
		fb := uint64(*trustedBatch.ForcedBatchNumber)
		processCtx.ForcedBatchNum = &fb
	}
	err := s.state.OpenBatch(ctx, processCtx, dbTx)
	if err != nil {
		log.Error("error opening batch: ", trustedBatch.Number)
		return err
	}
	return nil
}

func (s *SyncTrustedBatchesAction) decodeTxs(trustedBatchL2Data types.ArgBytes, batches []*state.Batch) ([]ethTypes.Transaction, []ethTypes.Transaction, []uint8, []uint8, error) {
	forkID := s.state.GetForkIDByBatchNumber(batches[0].BatchNumber)
	syncedTxs, _, syncedEfficiencyPercentages, err := state.DecodeTxs(trustedBatchL2Data, forkID)
	if err != nil {
		log.Errorf("error decoding synced txs from trustedstate. Error: %v, TrustedBatchL2Data: %s", err, trustedBatchL2Data.Hex())
		return nil, nil, nil, nil, err
	}
	storedTxs, _, storedEfficiencyPercentages, err := state.DecodeTxs(batches[0].BatchL2Data, forkID)
	if err != nil {
		log.Errorf("error decoding stored txs from trustedstate. Error: %v, batch.BatchL2Data: %s", err, common.Bytes2Hex(batches[0].BatchL2Data))
		return nil, nil, nil, nil, err
	}
	log.Debug("len(storedTxs): ", len(storedTxs))
	log.Debug("len(syncedTxs): ", len(syncedTxs))
	return storedTxs, syncedTxs, storedEfficiencyPercentages, syncedEfficiencyPercentages, nil
}

func (s *SyncTrustedBatchesAction) getCurrentBatches(ctx context.Context, batches []*state.Batch, trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, error) {
	if len(batches) == 0 || batches[0] == nil || (batches[0] != nil && uint64(trustedBatch.Number) != batches[0].BatchNumber) {
		log.Debug("Updating batch[0] value!")
		batch, err := s.state.GetBatchByNumber(ctx, uint64(trustedBatch.Number), dbTx)
		if err != nil && err != state.ErrNotFound {
			log.Warnf("failed to get batch %v from local trusted state. Error: %v", trustedBatch.Number, err)
			return nil, err
		}
		var prevBatch *state.Batch
		if len(batches) == 0 || batches[0] == nil || (batches[0] != nil && uint64(trustedBatch.Number-1) != batches[0].BatchNumber) {
			log.Debug("Updating batch[1] value!")
			prevBatch, err = s.state.GetBatchByNumber(ctx, uint64(trustedBatch.Number-1), dbTx)
			if err != nil && err != state.ErrNotFound {
				log.Warnf("failed to get prevBatch %v from local trusted state. Error: %v", trustedBatch.Number-1, err)
				return nil, err
			}
		} else {
			prevBatch = batches[0]
		}
		log.Debug("batch: ", batch)
		log.Debug("prevBatch: ", prevBatch)
		batches = []*state.Batch{batch, prevBatch}
	}
	return batches, nil
}

func (s *SyncTrustedBatchesAction) processAndStoreTxs(ctx context.Context, trustedBatch *types.Batch, request state.ProcessRequest, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	processBatchResp, err := s.state.ProcessBatch(ctx, request, true)
	if err != nil {
		log.Errorf("error processing sequencer batch for batch: %v", trustedBatch.Number)
		return nil, err
	}
	s.sync.PendingFlushID(processBatchResp.FlushID, processBatchResp.ProverID)

	log.Debugf("Storing %d blocks for batch %v", len(processBatchResp.BlockResponses), trustedBatch.Number)
	if processBatchResp.IsExecutorLevelError {
		log.Warn("executorLevelError detected. Avoid store txs...")
		return processBatchResp, nil
	} else if processBatchResp.IsRomOOCError {
		log.Warn("romOOCError detected. Avoid store txs...")
		return processBatchResp, nil
	}
	for _, block := range processBatchResp.BlockResponses {
		for _, tx := range block.TransactionResponses {
			if state.IsStateRootChanged(executor.RomErrorCode(tx.RomError)) {
				log.Infof("TrustedBatch info: %+v", processBatchResp)
				log.Infof("Storing trusted tx %+v", tx)
				if _, err = s.state.StoreTransaction(ctx, uint64(trustedBatch.Number), tx, trustedBatch.Coinbase, uint64(trustedBatch.Timestamp), nil, block.GlobalExitRoot, block.BlockInfoRoot, dbTx); err != nil {
					log.Errorf("failed to store transactions for batch: %v. Tx: %s", trustedBatch.Number, tx.TxHash.String())
					return nil, err
				}
			}
		}
	}
	return processBatchResp, nil
}

func checkIfSynced(batches []*state.Batch, trustedBatch *types.Batch) bool {
	matchNumber := batches[0].BatchNumber == uint64(trustedBatch.Number)
	matchGER := batches[0].GlobalExitRoot.String() == trustedBatch.GlobalExitRoot.String()
	matchLER := batches[0].LocalExitRoot.String() == trustedBatch.LocalExitRoot.String()
	matchSR := batches[0].StateRoot.String() == trustedBatch.StateRoot.String()
	matchCoinbase := batches[0].Coinbase.String() == trustedBatch.Coinbase.String()
	matchTimestamp := uint64(batches[0].Timestamp.Unix()) == uint64(trustedBatch.Timestamp)
	matchL2Data := hex.EncodeToString(batches[0].BatchL2Data) == hex.EncodeToString(trustedBatch.BatchL2Data)

	if matchNumber && matchGER && matchLER && matchSR &&
		matchCoinbase && matchTimestamp && matchL2Data {
		return true
	}
	log.Infof("matchNumber %v %d %d", matchNumber, batches[0].BatchNumber, uint64(trustedBatch.Number))
	log.Infof("matchGER %v %s %s", matchGER, batches[0].GlobalExitRoot.String(), trustedBatch.GlobalExitRoot.String())
	log.Infof("matchLER %v %s %s", matchLER, batches[0].LocalExitRoot.String(), trustedBatch.LocalExitRoot.String())
	log.Infof("matchSR %v %s %s", matchSR, batches[0].StateRoot.String(), trustedBatch.StateRoot.String())
	log.Infof("matchCoinbase %v %s %s", matchCoinbase, batches[0].Coinbase.String(), trustedBatch.Coinbase.String())
	log.Infof("matchTimestamp %v %d %d", matchTimestamp, uint64(batches[0].Timestamp.Unix()), uint64(trustedBatch.Timestamp))
	log.Infof("matchL2Data %v", matchL2Data)
	return false
}
