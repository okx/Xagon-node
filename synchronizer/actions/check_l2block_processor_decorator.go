package actions

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/jackc/pgx/v4"
)

// CheckL2BlockProcessorDecorator This class is just a decorator to call CheckL2Block
type CheckL2BlockProcessorDecorator struct {
	L1EventProcessor
	l2blockChecker *CheckL2BlockHash
}

// NewCheckL2BlockDecorator creates a new CheckL2BlockDecorator
func NewCheckL2BlockDecorator(l1EventProcessor L1EventProcessor, l2blockChecker *CheckL2BlockHash) *CheckL2BlockProcessorDecorator {
	return &CheckL2BlockProcessorDecorator{
		L1EventProcessor: l1EventProcessor,
		l2blockChecker:   l2blockChecker,
	}
}

// Process wraps the real Process and after check the L2Blocks
func (p *CheckL2BlockProcessorDecorator) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	res := p.L1EventProcessor.Process(ctx, order, l1Block, dbTx)
	if res != nil {
		return res
	}
	if p.l2blockChecker == nil {
		return nil
	}
	return p.l2blockChecker.CheckL2Block(ctx, dbTx)
}
