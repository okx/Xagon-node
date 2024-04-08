package pgstatestorage

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// AddForkID adds a new forkID to the storage
func (p *PostgresStorage) AddForkID(ctx context.Context, forkID state.ForkIDInterval, dbTx pgx.Tx) error {
	const addForkIDSQL = "INSERT INTO state.fork_id (from_batch_num, to_batch_num, fork_id, version, block_num) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (fork_id) DO UPDATE SET block_num = $5 WHERE state.fork_id.fork_id = $3;"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addForkIDSQL, forkID.FromBatchNumber, forkID.ToBatchNumber, forkID.ForkId, forkID.Version, forkID.BlockNumber)
	return err
}

// GetForkIDs get all the forkIDs stored
func (p *PostgresStorage) GetForkIDs(ctx context.Context, dbTx pgx.Tx) ([]state.ForkIDInterval, error) {
	const getForkIDsSQL = "SELECT from_batch_num, to_batch_num, fork_id, version, block_num FROM state.fork_id ORDER BY from_batch_num ASC"
	q := p.getExecQuerier(dbTx)

	rows, err := q.Query(ctx, getForkIDsSQL)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	forkIDs := make([]state.ForkIDInterval, 0, len(rows.RawValues()))

	for rows.Next() {
		var forkID state.ForkIDInterval
		if err := rows.Scan(
			&forkID.FromBatchNumber,
			&forkID.ToBatchNumber,
			&forkID.ForkId,
			&forkID.Version,
			&forkID.BlockNumber,
		); err != nil {
			return forkIDs, err
		}
		forkIDs = append(forkIDs, forkID)
	}
	return forkIDs, err
}

// UpdateForkID updates the forkID stored in db
func (p *PostgresStorage) UpdateForkIDToBatchNumber(ctx context.Context, forkID state.ForkIDInterval, dbTx pgx.Tx) error {
	const updateForkIDSQL = "UPDATE state.fork_id SET to_batch_num = $1 WHERE fork_id = $2"
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, updateForkIDSQL, forkID.ToBatchNumber, forkID.ForkId); err != nil {
		return err
	}
	return nil
}

// UpdateForkID updates the forkID stored in db
func (p *PostgresStorage) UpdateForkIDBlockNumber(ctx context.Context, forkdID uint64, newBlockNumber uint64, updateMemCache bool, dbTx pgx.Tx) error {
	const sql = "UPDATE state.fork_id SET block_num = $1 WHERE fork_id = $2"
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, sql, forkdID, newBlockNumber); err != nil {
		return err
	}
	if updateMemCache {
		log.Debugf("Updating forkID %d in memory", forkdID)
		forkIDs, err := p.GetForkIDs(ctx, dbTx)
		if err != nil {
			log.Error("error getting oldForkIDs. Error: ", err)
			return err
		}
		p.UpdateForkIDIntervalsInMemory(forkIDs)
	}
	return nil
}

// UpdateForkIDIntervalsInMemory updates the forkID intervals in memory
func (p *PostgresStorage) UpdateForkIDIntervalsInMemory(intervals []state.ForkIDInterval) {
	log.Infof("Updating forkIDs. Setting %d forkIDs", len(intervals))
	log.Infof("intervals: %#v", intervals)
	p.cfg.ForkIDIntervals = intervals
}

// AddForkIDInterval updates the forkID intervals
func (p *PostgresStorage) AddForkIDInterval(ctx context.Context, newForkID state.ForkIDInterval, dbTx pgx.Tx) error {
	// Add forkId to db and memori variable
	oldForkIDs, err := p.GetForkIDs(ctx, dbTx)
	if err != nil {
		log.Error("error getting oldForkIDs. Error: ", err)
		return err
	}
	if len(oldForkIDs) == 0 {
		p.UpdateForkIDIntervalsInMemory([]state.ForkIDInterval{newForkID})
	} else {
		var forkIDs []state.ForkIDInterval
		forkIDs = oldForkIDs
		// Check to detect forkID inconsistencies
		if forkIDs[len(forkIDs)-1].ForkId >= newForkID.ForkId {
			errMsg := "error checking forkID sequence. Last ForkID stored: %d. New ForkID received: %d"
			err := fmt.Errorf(errMsg, forkIDs[len(forkIDs)-1].ForkId, newForkID.ForkId)
			log.Errorf(err.Error())
			return err
		}
		forkIDs[len(forkIDs)-1].ToBatchNumber = newForkID.FromBatchNumber - 1
		err := p.UpdateForkIDToBatchNumber(ctx, forkIDs[len(forkIDs)-1], dbTx)
		if err != nil {
			log.Errorf("error updating forkID: %d. Error: %v", forkIDs[len(forkIDs)-1].ForkId, err)
			return err
		}
		forkIDs = append(forkIDs, newForkID)

		p.UpdateForkIDIntervalsInMemory(forkIDs)
	}
	err = p.AddForkID(ctx, newForkID, dbTx)
	if err != nil {
		log.Errorf("error adding forkID %d. Error: %v", newForkID.ForkId, err)
		return err
	}
	return nil
}

// GetForkIDByBlockNumber returns the fork id for a given block number
func (p *PostgresStorage) GetForkIDByBlockNumber(blockNumber uint64) uint64 {
	if p.cfg.AvoidForkIDInMemory {
		const query = `
			SELECT fork_id
			  FROM state.fork_id
			 WHERE block_num <= $1
			 ORDER BY fork_id DESC
			 LIMIT 1`
		q := p.getExecQuerier(nil)

		var forkID uint64
		err := q.QueryRow(context.Background(), query, blockNumber).Scan(&forkID)
		if errors.Is(err, pgx.ErrNoRows) {
			return 1
		} else if err != nil {
			log.Warnf("failed to get forkID by blockNumber from db, falling back to in memory information, err: %v", err)
			return p.GetForkIDByBlockNumberInMemory(blockNumber)
		}

		return forkID
	} else {
		return p.GetForkIDByBlockNumberInMemory(blockNumber)
	}
}

// GetForkIDByBlockNumber returns the fork id for a given block number in memory
func (p *PostgresStorage) GetForkIDByBlockNumberInMemory(blockNumber uint64) uint64 {
	for _, index := range sortIndexForForkdIDSortedByBlockNumber(p.cfg.ForkIDIntervals) {
		// reverse travesal
		interval := p.cfg.ForkIDIntervals[len(p.cfg.ForkIDIntervals)-1-index]
		if blockNumber >= interval.BlockNumber {
			return interval.ForkId
		}
	}
	// If not found return the  fork id 1
	return 1
}

func sortIndexForForkdIDSortedByBlockNumber(forkIDs []state.ForkIDInterval) []int {
	sortedIndex := make([]int, len(forkIDs))
	for i := range sortedIndex {
		sortedIndex[i] = i
	}
	cmpFunc := func(i, j int) bool {
		return forkIDs[sortedIndex[i]].BlockNumber < forkIDs[sortedIndex[j]].BlockNumber
	}
	sort.Slice(sortedIndex, cmpFunc)
	return sortedIndex
}

// GetForkIDByBatchNumber returns the fork id for a given batch number
func (p *PostgresStorage) GetForkIDByBatchNumber(batchNumber uint64) uint64 {
	if batchNumber == 0 {
		batchNumber = 1
	}

	if p.cfg.AvoidForkIDInMemory {
		const query = `
			SELECT fork_id FROM state.fork_id
			 WHERE from_batch_num <= $1 AND to_batch_num >= $1
			 ORDER BY fork_id DESC
			 LIMIT 1`
		q := p.getExecQuerier(nil)

		var forkID uint64
		err := q.QueryRow(context.Background(), query, batchNumber).Scan(&forkID)
		if errors.Is(err, pgx.ErrNoRows) {
			const query = `
				SELECT fork_id
				FROM state.fork_id
				ORDER BY fork_id DESC
				LIMIT 1`
			q := p.getExecQuerier(nil)
			err := q.QueryRow(context.Background(), query).Scan(&forkID)
			if errors.Is(err, pgx.ErrNoRows) {
				log.Warnf("can't find forkID by batchNumber in the db, falling back to in memory information, err: %v", err)
				return p.GetForkIDByBatchNumberInMemory(batchNumber)
			} else if err != nil {
				log.Warnf("failed to get forkID by batchNumber from db, falling back to in memory information, err: %v", err)
				return p.GetForkIDByBatchNumberInMemory(batchNumber)
			}
		} else if err != nil {
			log.Warnf("failed to get forkID by batchNumber from db, falling back to in memory information, err: %v", err)
			return p.GetForkIDByBatchNumberInMemory(batchNumber)
		}

		return forkID
	} else {
		return p.GetForkIDByBatchNumberInMemory(batchNumber)
	}
}

// GetForkIDByBatchNumberInMemory returns the fork id for a given batch number
func (p *PostgresStorage) GetForkIDByBatchNumberInMemory(batchNumber uint64) uint64 {
	if batchNumber == 0 {
		batchNumber = 1
	}

	// If NumBatchForkIdUpgrade is defined (!=0) we are performing forkid upgrade process
	// In this case, if the batchNumber is the next to the NumBatchForkIdUpgrade, we need to return the
	// new "future" forkId (ForkUpgradeNewForkId)
	if (p.cfg.ForkUpgradeBatchNumber) != 0 && (batchNumber > p.cfg.ForkUpgradeBatchNumber) {
		return p.cfg.ForkUpgradeNewForkId
	}

	for _, interval := range p.cfg.ForkIDIntervals {
		if batchNumber >= interval.FromBatchNumber && batchNumber <= interval.ToBatchNumber {
			return interval.ForkId
		}
	}

	// If not found return the last fork id
	return p.cfg.ForkIDIntervals[len(p.cfg.ForkIDIntervals)-1].ForkId
}

// GetForkIDInMemory get the forkIDs stored in cache, or nil if not found
func (p *PostgresStorage) GetForkIDInMemory(forkId uint64) *state.ForkIDInterval {
	for _, interval := range p.cfg.ForkIDIntervals {
		if interval.ForkId == forkId {
			return &interval
		}
	}
	return nil
}
