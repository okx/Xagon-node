package pgstatestorage

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// AddBlobSequence adds a new blob sequence to the state.
// TODO: Add support to ReceivedAt
func (p *PostgresStorage) AddBlobSequence(ctx context.Context, blobSequence *state.BlobSequence, dbTx pgx.Tx) error {
	const addBlobSequenceSQL = "INSERT INTO state.blob_sequence (index, block_num, coinbase, final_acc_input_hash, first_blob_sequenced, last_blob_sequenced, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addBlobSequenceSQL, blobSequence.BlobSequenceIndex, blobSequence.BlockNumber, blobSequence.L2Coinbase.String(), blobSequence.FinalAccInputHash.String(), blobSequence.FirstBlobSequenced, blobSequence.LastBlobSequenced, blobSequence.CreateAt)
	return err
}

// GetLastBlobSequence returns the last blob sequence stored in the state.
// TODO: Add support to ReceivedAt
func (p *PostgresStorage) GetLastBlobSequence(ctx context.Context, dbTx pgx.Tx) (*state.BlobSequence, error) {
	var (
		coinbase          string
		finalAccInputHash string
		createAt          time.Time
		blobSequence      state.BlobSequence
	)
	const getLastBlobSequenceSQL = "SELECT index, coinbase, final_acc_input_hash, first_blob_sequenced, last_blob_sequenced, created_at FROM state.blob_sequence ORDER BY index DESC LIMIT 1"

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getLastBlobSequenceSQL).Scan(&blobSequence.BlobSequenceIndex, &coinbase, &finalAccInputHash, &blobSequence.FirstBlobSequenced, &blobSequence.LastBlobSequenced, &createAt)
	if errors.Is(err, pgx.ErrNoRows) {
		// If none on database return a nil object
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	blobSequence.L2Coinbase = common.HexToAddress(coinbase)
	blobSequence.FinalAccInputHash = common.HexToHash(finalAccInputHash)
	blobSequence.CreateAt = createAt
	return &blobSequence, nil
}
