package l2_shared

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// Implements PostClosedBatchChecker

type statePostClosedBatchCheckL2Block interface {
	GetLastL2BlockByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.L2Block, error)
}

// PostClosedBatchCheckL2Block is a struct that implements the PostClosedBatchChecker interface and check the las L2Block hash on close batch
type PostClosedBatchCheckL2Block struct {
	state statePostClosedBatchCheckL2Block
}

// NewPostClosedBatchCheckL2Block creates a new PostClosedBatchCheckL2Block
func NewPostClosedBatchCheckL2Block(state statePostClosedBatchCheckL2Block) *PostClosedBatchCheckL2Block {
	return &PostClosedBatchCheckL2Block{
		state: state,
	}
}

// CheckPostClosedBatch checks the last L2Block hash on close batch
func (p *PostClosedBatchCheckL2Block) CheckPostClosedBatch(ctx context.Context, processData ProcessData, dbTx pgx.Tx) error {
	if processData.TrustedBatch == nil {
		log.Warnf("%s trusted batch is nil", processData.DebugPrefix)
		return nil
	}
	if len(processData.TrustedBatch.Blocks) == 0 {
		log.Infof("%s trusted batch have no Blocks, so nothing to check", processData.DebugPrefix)
		return nil
	}

	// Get last L2Block from the database
	statelastL2Block, err := p.state.GetLastL2BlockByBatchNumber(ctx, processData.BatchNumber, dbTx)
	if err != nil {
		return err
	}
	if statelastL2Block == nil {
		return fmt.Errorf("last L2Block in the database is nil")
	}
	trustedLastL2Block := processData.TrustedBatch.Blocks[len(processData.TrustedBatch.Blocks)-1].Block
	log.Info(trustedLastL2Block)
	if statelastL2Block.Number().Cmp(big.NewInt(int64(trustedLastL2Block.Number))) != 0 {
		return fmt.Errorf("last L2Block in the database %s and the trusted batch %d are different", statelastL2Block.Number().String(), trustedLastL2Block.Number)
	}

	if statelastL2Block.Hash() != *trustedLastL2Block.Hash {
		return fmt.Errorf("last L2Block %s in the database %s and the trusted batch %s are different", statelastL2Block.Number().String(), statelastL2Block.Hash().String(), trustedLastL2Block.Hash.String())
	}
	log.Infof("%s last L2Block in the database %s and the trusted batch %s are the same", processData.DebugPrefix, statelastL2Block.Number().String(), trustedLastL2Block.Number)
	// Compare the two blocks

	return nil
}
