package pgstatestorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// CheckProofContainsCompleteSequences checks if a recursive proof contains complete sequences
func (p *PostgresStorage) CheckProofContainsCompleteSequences(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) (bool, error) {
	const getProofContainsCompleteSequencesSQL = `
		SELECT EXISTS (SELECT 1 FROM state.sequences s1 WHERE s1.from_batch_num = $1) AND
			   EXISTS (SELECT 1 FROM state.sequences s2 WHERE s2.to_batch_num = $2)
		`
	e := p.getExecQuerier(dbTx)
	var exists bool
	err := e.QueryRow(ctx, getProofContainsCompleteSequencesSQL, proof.BatchNumber, proof.BatchNumberFinal).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return exists, err
	}
	return exists, nil
}

// GetProofReadyForFinal return the proof that is ready to generate the final proof
func (p *PostgresStorage) GetProofReadyForFinal(ctx context.Context, lastVerfiedBatchNumber uint64, dbTx pgx.Tx) (*state.Proof, error) {
	const getProofReadyForFinalSQL = `
		SELECT 
			p.batch_num, 
			p.batch_num_final,
			p.proof,
			p.proof_id,
			p.input_prover,
			p.prover,
			p.prover_id,
			p.generating_since,
			p.created_at,
			p.updated_at
		FROM state.batch_proof p
		WHERE batch_num = $1 AND generating_since IS NULL AND
			EXISTS (SELECT 1 FROM state.sequences s1 WHERE s1.from_batch_num = p.batch_num) AND
			EXISTS (SELECT 1 FROM state.sequences s2 WHERE s2.to_batch_num = p.batch_num_final)		
		`

	var proof *state.Proof = &state.Proof{}

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getProofReadyForFinalSQL, lastVerfiedBatchNumber+1)
	err := row.Scan(&proof.BatchNumber, &proof.BatchNumberFinal, &proof.Proof, &proof.ProofID, &proof.InputProver, &proof.Prover, &proof.ProverID, &proof.GeneratingSince, &proof.CreatedAt, &proof.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return proof, err
}

// GetBatchProofsToAggregate return the next 2 batch proofs that it are possible to aggregate
func (p *PostgresStorage) GetBatchProofsToAggregate(ctx context.Context, dbTx pgx.Tx) (*state.Proof, *state.Proof, error) {
	var (
		proof1 *state.Proof = &state.Proof{}
		proof2 *state.Proof = &state.Proof{}
	)

	// TODO: add comments to explain the query
	const getBatchProofsToAggregateSQL = `
		SELECT 
			p1.batch_num as p1_batch_num, 
			p1.batch_num_final as p1_batch_num_final, 
			p1.proof as p1_proof,	
			p1.proof_id as p1_proof_id, 
			p1.input_prover as p1_input_prover, 
			p1.prover as p1_prover,
			p1.prover_id as p1_prover_id,
			p1.generating_since as p1_generating_since,
			p1.created_at as p1_created_at,
			p1.updated_at as p1_updated_at,
			p2.batch_num as p2_batch_num, 
			p2.batch_num_final as p2_batch_num_final, 
			p2.proof as p2_proof,	
			p2.proof_id as p2_proof_id, 
			p2.input_prover as p2_input_prover, 
			p2.prover as p2_prover,
			p2.prover_id as p2_prover_id,
			p2.generating_since as p2_generating_since,
			p2.created_at as p2_created_at,
			p2.updated_at as p2_updated_at
		FROM state.batch_proof p1 INNER JOIN state.batch_proof p2 ON p1.batch_num_final = p2.batch_num - 1
		WHERE p1.blob_inner_num = p2.blob_inner_num AND
			  p1.generating_since IS NULL AND p2.generating_since IS NULL AND 
		 	  p1.proof IS NOT NULL AND p2.proof IS NOT NULL
		ORDER BY p1.batch_num ASC
		LIMIT 1
		`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getBatchProofsToAggregateSQL)
	err := row.Scan(
		&proof1.BatchNumber, &proof1.BatchNumberFinal, &proof1.Proof, &proof1.ProofID, &proof1.InputProver, &proof1.Prover, &proof1.ProverID, &proof1.GeneratingSince, &proof1.CreatedAt, &proof1.UpdatedAt,
		&proof2.BatchNumber, &proof2.BatchNumberFinal, &proof2.Proof, &proof2.ProofID, &proof2.InputProver, &proof2.Prover, &proof2.ProverID, &proof2.GeneratingSince, &proof2.CreatedAt, &proof2.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, state.ErrNotFound
	} else if err != nil {
		return nil, nil, err
	}

	return proof1, proof2, err
}

// AddBatchProof adds a batch proof to the storage
func (p *PostgresStorage) AddBatchProof(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) error {
	const addBatchProofSQL = "INSERT INTO state.batch_proof (batch_num, batch_num_final, proof, proof_id, input_prover, prover, prover_id, generating_since, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	e := p.getExecQuerier(dbTx)
	now := time.Now().UTC().Round(time.Microsecond)
	_, err := e.Exec(ctx, addBatchProofSQL, proof.BatchNumber, proof.BatchNumberFinal, proof.Proof, proof.ProofID, proof.InputProver, proof.Prover, proof.ProverID, proof.GeneratingSince, now, now)
	return err
}

// UpdateBatchProof updates a batch proof in the storage
func (p *PostgresStorage) UpdateBatchProof(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) error {
	const addBatchProofSQL = "UPDATE state.batch_proof SET proof = $3, proof_id = $4, input_prover = $5, prover = $6, prover_id = $7, generating_since = $8, updated_at = $9 WHERE batch_num = $1 AND batch_num_final = $2"
	e := p.getExecQuerier(dbTx)
	now := time.Now().UTC().Round(time.Microsecond)
	_, err := e.Exec(ctx, addBatchProofSQL, proof.BatchNumber, proof.BatchNumberFinal, proof.Proof, proof.ProofID, proof.InputProver, proof.Prover, proof.ProverID, proof.GeneratingSince, now)
	return err
}

// DeleteBatchProofs deletes from the storage the batch proofs falling inside the batch numbers range.
func (p *PostgresStorage) DeleteBatchProofs(ctx context.Context, batchNumber uint64, batchNumberFinal uint64, dbTx pgx.Tx) error {
	const deleteBatchProofSQL = "DELETE FROM state.batch_proof WHERE batch_num >= $1 AND batch_num_final <= $2"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, deleteBatchProofSQL, batchNumber, batchNumberFinal)
	return err
}

// CleanupBatchProofs deletes from the storage the batch proofs up to the specified batch number included.
func (p *PostgresStorage) CleanupBatchProofs(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	const deleteBatchProofSQL = "DELETE FROM state.batch_proof WHERE batch_num_final <= $1"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, deleteBatchProofSQL, batchNumber)
	return err
}

// CleanupLockedBatchProofs deletes from the storage the proofs locked in generating state for more than the provided threshold.
func (p *PostgresStorage) CleanupLockedBatchProofs(ctx context.Context, duration string, dbTx pgx.Tx) (int64, error) {
	interval, err := toPostgresInterval(duration)
	if err != nil {
		return 0, err
	}
	sql := fmt.Sprintf("DELETE FROM state.batch_proof WHERE generating_since < (NOW() - interval '%s')", interval)
	e := p.getExecQuerier(dbTx)
	ct, err := e.Exec(ctx, sql)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}

// DeleteUngeneratedBatchProofs deletes ungenerated proofs. This method is meant to be use during aggregator boot-up sequence
func (p *PostgresStorage) DeleteUngeneratedBatchProofs(ctx context.Context, dbTx pgx.Tx) error {
	const deleteUngeneratedProofsSQL = "DELETE FROM state.batch_proof WHERE generating_since IS NOT NULL"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, deleteUngeneratedProofsSQL)
	return err
}

func toPostgresInterval(duration string) (string, error) {
	unit := duration[len(duration)-1]
	var pgUnit string

	switch unit {
	case 's':
		pgUnit = "second"
	case 'm':
		pgUnit = "minute"
	case 'h':
		pgUnit = "hour"
	default:
		return "", state.ErrUnsupportedDuration
	}

	isMoreThanOne := duration[0] != '1' || len(duration) > 2 //nolint:gomnd
	if isMoreThanOne {
		pgUnit = pgUnit + "s"
	}

	return fmt.Sprintf("%s %s", duration[:len(duration)-1], pgUnit), nil
}
