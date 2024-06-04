package sequencesender

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

func (s *SequenceSender) tryToSendSequenceXLayer(ctx context.Context) {
	retry := false
	// process monitored sequences before starting a next cycle
	s.ethTxManager.ProcessPendingMonitoredTxs(ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
		if result.Status == ethtxmanager.MonitoredTxStatusConfirmed {
			if len(result.Txs) > 0 {
				var txL1BlockNumber uint64
				var txHash common.Hash
				receiptFound := false
				for _, tx := range result.Txs {
					if tx.Receipt != nil {
						txL1BlockNumber = tx.Receipt.BlockNumber.Uint64()
						txHash = tx.Tx.Hash()
						receiptFound = true
						break
					}
				}

				if !receiptFound {
					s.halt(ctx, fmt.Errorf("monitored tx %s for sequence [%d-%d] is confirmed but doesn't have a receipt", result.ID, s.lastSequenceInitialBatch, s.lastSequenceEndBatch))
				}

				// wait L1 confirmation blocks
				log.Infof("waiting %d L1 block confirmations for sequence [%d-%d], L1 block: %d, tx: %s",
					s.cfg.SequenceL1BlockConfirmations, s.lastSequenceInitialBatch, s.lastSequenceEndBatch, txL1BlockNumber, txHash)
				for {
					lastL1BlockHeader, err := s.etherman.GetLatestBlockHeader(ctx)
					if err != nil {
						log.Errorf("failed to get last L1 block number, err: %v", err)
					} else {
						lastL1BlockNumber := lastL1BlockHeader.Number.Uint64()

						if lastL1BlockNumber >= txL1BlockNumber+s.cfg.SequenceL1BlockConfirmations {
							log.Infof("continuing, last L1 block: %d", lastL1BlockNumber)
							break
						}
					}
					time.Sleep(waitRetryGetL1Block)
				}

				lastSCBatchNum, err := s.etherman.GetLatestBatchNumber()
				if err != nil {
					log.Warnf("failed to get from the SC last sequenced batch number, err: %v", err)
					return
				}

				// If it's the first time we call that function after the restart of the sequence-sender (lastSequenceBatch is 0) and we are having the
				// confirmation of a pending L1 tx sent before the sequence-sender was restarted, we don't know which batch was the last sequenced.
				// Therefore we cannot compare the last sequenced batch in the SC with the last sequenced from sequence-sender. We skip this check
				if s.lastSequenceEndBatch != 0 && (lastSCBatchNum != s.lastSequenceEndBatch) {
					s.halt(ctx, fmt.Errorf("last sequenced batch from SC %d doesn't match last sequenced batch sent %d", lastSCBatchNum, s.lastSequenceEndBatch))
				}
			} else {
				s.halt(ctx, fmt.Errorf("monitored tx %s for sequence [%d-%d] doesn't have transactions to be checked", result.ID, s.lastSequenceInitialBatch, s.lastSequenceEndBatch))
			}
		} else { // Monitored tx is failed
			retry = true
			mTxResultLogger := ethtxmanager.CreateMonitoredTxResultLogger(ethTxManagerOwner, result)
			mTxResultLogger.Error("failed to send sequence, TODO: review this fatal and define what to do in this case")
		}
	}, nil)

	if retry {
		return
	}

	sanityCheckOk, err := s.sanityCheck(ctx, retriesSanityCheck, waitRetrySanityCheck)
	if err != nil {
		s.halt(ctx, err)
	}
	if !sanityCheckOk {
		log.Info("sanity check failed, retrying...")
		time.Sleep(5 * time.Second) // nolint:gomnd
		return
	}

	// Check if should send sequence to L1
	log.Infof("getting sequences to send")
	sequences, err := s.getSequencesToSendXLayer(ctx)
	if err != nil || len(sequences) == 0 {
		if err != nil {
			log.Errorf("error getting sequences: %v", err)
		} else {
			log.Info("waiting for sequences to be worth sending to L1")
		}
		time.Sleep(s.cfg.WaitPeriodSendSequence.Duration)
		return
	}

	lastVirtualBatchNum, err := s.etherman.GetLatestBatchNumber()
	if err != nil {
		log.Errorf("failed to get last virtual batch num, err: %v", err)
		return
	}

	// Send sequences to L1
	sequenceCount := len(sequences)
	log.Infof("sending sequences to L1. From batch %d to batch %d", lastVirtualBatchNum+1, lastVirtualBatchNum+uint64(sequenceCount))

	// Check if we need to wait until last L1 block timestamp is L1BlockTimestampMargin seconds above the timestamp of the last L2 block in the sequence
	// Get last sequence
	lastSequence := sequences[sequenceCount-1]
	// Get timestamp of the last L2 block in the sequence
	lastL2BlockTimestamp := uint64(lastSequence.LastL2BLockTimestamp)

	timeMargin := int64(s.cfg.L1BlockTimestampMargin.Seconds())

	// Wait until last L1 block timestamp is timeMargin (L1BlockTimestampMargin) seconds above the timestamp of the last L2 block in the sequence
	for {
		// Get header of the last L1 block
		lastL1BlockHeader, err := s.etherman.GetLatestBlockHeader(ctx)
		if err != nil {
			log.Errorf("failed to get last L1 block timestamp, err: %v", err)
			return
		}

		elapsed, waitTime := s.marginTimeElapsed(lastL2BlockTimestamp, lastL1BlockHeader.Time, timeMargin)

		if !elapsed {
			log.Infof("waiting at least %d seconds to send sequences, time difference between last L1 block %d (ts: %d) and last L2 block %d (ts: %d) in the sequence is lower than %d seconds",
				waitTime, lastL1BlockHeader.Number, lastL1BlockHeader.Time, lastSequence.BatchNumber, lastL2BlockTimestamp, timeMargin)
			time.Sleep(time.Duration(waitTime) * time.Second)
		} else {
			log.Infof("continuing, time difference between last L1 block %d (ts: %d) and last L2 block %d (ts: %d) in the sequence is greater than %d seconds",
				lastL1BlockHeader.Number, lastL1BlockHeader.Time, lastSequence.BatchNumber, lastL2BlockTimestamp, timeMargin)
			break
		}
	}

	// Sanity check. Wait also until current time (now) is timeMargin (L1BlockTimestampMargin) seconds above the timestamp of the last L2 block in the sequence
	for {
		currentTime := uint64(time.Now().Unix())

		elapsed, waitTime := s.marginTimeElapsed(lastL2BlockTimestamp, currentTime, timeMargin)

		// Wait if the time difference is less than timeMargin (L1BlockTimestampMargin)
		if !elapsed {
			log.Infof("waiting at least %d seconds to send sequences, time difference between now (ts: %d) and last L2 block %d (ts: %d) in the sequence is lower than %d seconds",
				waitTime, currentTime, lastSequence.BatchNumber, lastL2BlockTimestamp, timeMargin)
			time.Sleep(time.Duration(waitTime) * time.Second)
		} else {
			log.Infof("sending sequences now, time difference between now (ts: %d) and last L2 block %d (ts: %d) in the sequence is also greater than %d seconds",
				currentTime, lastSequence.BatchNumber, lastL2BlockTimestamp, timeMargin)
			break
		}
	}

	// add sequence to be monitored
	firstSequence := sequences[0]
	dataAvailabilityMessage, err := s.da.PostSequence(ctx, sequences)
	if err != nil {
		log.Error("error posting sequences to the data availability protocol: ", err)
		return
	}
	to, data, err := s.etherman.BuildSequenceBatchesTxDataXLayer(s.cfg.SenderAddress, sequences, uint64(lastSequence.LastL2BLockTimestamp), firstSequence.BatchNumber-1, s.cfg.L2Coinbase, dataAvailabilityMessage)
	if err != nil {
		log.Error("error estimating new sequenceBatches to add to eth tx manager: ", err)
		return
	}

	monitoredTxID := fmt.Sprintf(monitoredIDFormat, firstSequence.BatchNumber, lastSequence.BatchNumber)
	err = s.ethTxManager.Add(ctx, ethTxManagerOwner, monitoredTxID, s.cfg.SenderAddress, to, nil, data, s.cfg.GasOffset, nil)
	if err != nil {
		mTxLogger := ethtxmanager.CreateLogger(ethTxManagerOwner, monitoredTxID, s.cfg.SenderAddress, to)
		mTxLogger.Errorf("error to add sequences tx to eth tx manager: ", err)
		return
	}

	s.lastSequenceInitialBatch = sequences[0].BatchNumber
	s.lastSequenceEndBatch = lastSequence.BatchNumber
}

// getSequencesToSend generates an array of sequences to be send to L1.
// If the array is empty, it doesn't necessarily mean that there are no sequences to be sent,
// it could be that it's not worth it to do so yet.
func (s *SequenceSender) getSequencesToSendXLayer(ctx context.Context) ([]types.Sequence, error) {
	lastVirtualBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get last virtual batch num, err: %v", err)
	}
	log.Debugf("last virtual batch number: %d", lastVirtualBatchNum)

	currentBatchNumToSequence := lastVirtualBatchNum + 1
	log.Debugf("current batch number to sequence: %d", currentBatchNumToSequence)

	sequences := []types.Sequence{}

	// Add sequences until too big for a single L1 tx or last batch is reached
	for {
		//Check if the next batch belongs to a new forkid, in this case we need to stop sequencing as we need to
		//wait the upgrade of forkid is completed and s.cfg.NumBatchForkIdUpgrade is disabled (=0) again
		if (s.cfg.ForkUpgradeBatchNumber != 0) && (currentBatchNumToSequence == (s.cfg.ForkUpgradeBatchNumber + 1)) {
			return nil, fmt.Errorf("aborting sequencing process as we reached the batch %d where a new forkid is applied (upgrade)", s.cfg.ForkUpgradeBatchNumber+1)
		}

		// Add new sequence
		batch, err := s.state.GetBatchByNumber(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			if err == state.ErrNotFound {
				break
			}
			log.Debugf("failed to get batch by number %d, err: %v", currentBatchNumToSequence, err)
			return nil, err
		}

		// Check if batch is closed and checked (sequencer sanity check was successful)
		isChecked, err := s.state.IsBatchChecked(ctx, currentBatchNumToSequence, nil)
		if err != nil {
			log.Debugf("failed to check if batch %d is closed and checked, err: %v", currentBatchNumToSequence, err)
			return nil, err
		}

		if !isChecked {
			// Batch is not closed and checked
			break
		}

		seq := types.Sequence{
			BatchL2Data: batch.BatchL2Data,
			BatchNumber: batch.BatchNumber,
		}

		if batch.ForcedBatchNum != nil {
			forcedBatch, err := s.state.GetForcedBatch(ctx, *batch.ForcedBatchNum, nil)
			if err != nil {
				return nil, err
			}

			// Get L1 block for the forced batch
			fbL1Block, err := s.state.GetBlockByNumber(ctx, forcedBatch.BlockNumber, nil)
			if err != nil {
				return nil, err
			}

			seq.GlobalExitRoot = forcedBatch.GlobalExitRoot
			seq.ForcedBatchTimestamp = forcedBatch.ForcedAt.Unix()
			seq.PrevBlockHash = fbL1Block.ParentHash
			// Set sequence timestamps as the forced batch timestamp
			seq.LastL2BLockTimestamp = seq.ForcedBatchTimestamp
		} else {
			// Set sequence timestamps as the latest l2 block timestamp
			lastL2Block, err := s.state.GetLastL2BlockByBatchNumber(ctx, currentBatchNumToSequence, nil)
			if err != nil {
				return nil, err
			}
			if lastL2Block == nil {
				return nil, fmt.Errorf("no last L2 block returned from the state for batch %d", currentBatchNumToSequence)
			}

			// Get timestamp of the last L2 block in the sequence
			seq.LastL2BLockTimestamp = lastL2Block.ReceivedAt.Unix()
		}

		sequences = append(sequences, seq)
		if len(sequences) == int(s.cfg.MaxBatchesForL1) {
			log.Infof(
				"sequence should be sent to L1, because MaxBatchesForL1 (%d) has been reached",
				s.cfg.MaxBatchesForL1,
			)
			return sequences, nil
		}

		//Check if the current batch is the last before a change to a new forkid, in this case we need to close and send the sequence to L1
		if (s.cfg.ForkUpgradeBatchNumber != 0) && (currentBatchNumToSequence == (s.cfg.ForkUpgradeBatchNumber)) {
			log.Infof("sequence should be sent to L1, as we have reached the batch %d from which a new forkid is applied (upgrade)", s.cfg.ForkUpgradeBatchNumber)
			return sequences, nil
		}

		// Increase batch num for next iteration
		currentBatchNumToSequence++
	}

	// Reached latest batch. Decide if it's worth to send the sequence, or wait for new batches
	if len(sequences) == 0 {
		log.Info("no batches to be sequenced")
		return nil, nil
	}

	lastBatchVirtualizationTime, err := s.state.GetTimeForLatestBatchVirtualization(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		log.Warnf("failed to get last l1 interaction time, err: %v. Sending sequences as a conservative approach", err)
		return sequences, nil
	}
	if lastBatchVirtualizationTime.Before(time.Now().Add(-s.cfg.LastBatchVirtualizationTimeMaxWaitPeriod.Duration)) {
		// TODO: implement check profitability
		// if s.checker.IsSendSequencesProfitable(new(big.Int).SetUint64(estimatedGas), sequences) {
		log.Info("sequence should be sent to L1, because too long since didn't send anything to L1")
		return sequences, nil
		//}
	}

	log.Info("not enough time has passed since last batch was virtualized, and the sequence could be bigger")
	return nil, nil
}

// SetDataProvider sets the data provider
func (s *SequenceSender) SetDataProvider(da dataAbilitier) {
	s.da = da
}
