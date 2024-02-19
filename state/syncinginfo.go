package state

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
)

// SyncInfoDataOnStorage stores information regarding the syncing status of the node in the database
type SyncInfoDataOnStorage struct {
	InitialSyncingBatch         uint64
	LastBatchNumberSeen         uint64
	LastBatchNumberConsolidated uint64
}

// SyncingInfo stores information regarding the syncing status of the node
type SyncingInfo struct {
	InitialSyncingBlock   uint64 // L2Block corresponding to InitialSyncingBatch
	CurrentBlockNumber    uint64 // last L2Block in state
	EstimatedHighestBlock uint64 // estimated highest L2Block in state

	// IsSynchronizing indicates if the node is syncing (true -> syncing, false -> fully synced)
	IsSynchronizing bool
}

// GetSyncingInfo returns information regarding the syncing status of the node
func (p *State) GetSyncingInfo(ctx context.Context, dbTx pgx.Tx) (SyncingInfo, error) {
	var info SyncingInfo

	syncData, err := p.GetSyncInfoData(ctx, dbTx)
	if errors.Is(err, ErrNotFound) {
		return SyncingInfo{}, ErrStateNotSynchronized
	} else if err != nil {
		return SyncingInfo{}, err
	}

	info.InitialSyncingBlock, err = p.GetFirstL2BlockNumberForBatchNumber(ctx, syncData.InitialSyncingBatch, dbTx)
	if errors.Is(err, ErrNotFound) {
		return SyncingInfo{}, ErrStateNotSynchronized
	} else if err != nil {
		return SyncingInfo{}, err
	}

	lastBlockNumber, err := p.GetLastL2BlockNumber(ctx, dbTx)
	if errors.Is(err, ErrNotFound) {
		return SyncingInfo{}, ErrStateNotSynchronized
	} else if err != nil {
		return SyncingInfo{}, err
	}
	info.CurrentBlockNumber = lastBlockNumber

	lastBatchNumber, err := p.GetLastBatchNumber(ctx, dbTx)
	if errors.Is(err, ErrNotFound) {
		return SyncingInfo{}, ErrStateNotSynchronized
	} else if err != nil {
		return SyncingInfo{}, err
	}

	info.IsSynchronizing = syncData.LastBatchNumberSeen > lastBatchNumber
	if info.IsSynchronizing {
		// Estimation of block counting 1 l2block per missing batch
		info.EstimatedHighestBlock = lastBlockNumber + (syncData.LastBatchNumberConsolidated - lastBatchNumber)
	} else {
		info.EstimatedHighestBlock = lastBlockNumber
	}

	return info, nil
}
