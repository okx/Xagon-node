package actions

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Implements PostClosedBatchChecker

type stateGetL2Block interface {
	GetL2BlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*state.L2Block, error)
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
}

type trustedRPCGetL2Block interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

// CheckL2BlockHash is a struct that implements a checker of L2Block hash
type CheckL2BlockHash struct {
	state              stateGetL2Block
	trustedClient      trustedRPCGetL2Block
	lastL2BlockChecked uint64
	// Is a modulus used to choose the l2block to check
	modulusL2BlockToCheck uint64
}

// NewCheckL2BlockHash creates a new CheckL2BlockHash
func NewCheckL2BlockHash(state stateGetL2Block,
	trustedClient trustedRPCGetL2Block,
	initialL2BlockNumber uint64,
	modulusBlockNumber uint64) (*CheckL2BlockHash, error) {
	if modulusBlockNumber == 0 {
		return nil, fmt.Errorf("error: modulusBlockNumber is zero")
	}
	return &CheckL2BlockHash{
		state:                 state,
		trustedClient:         trustedClient,
		lastL2BlockChecked:    initialL2BlockNumber,
		modulusL2BlockToCheck: modulusBlockNumber,
	}, nil
}

// CheckL2Block checks the  L2Block hash between the local and the trusted
func (p *CheckL2BlockHash) CheckL2Block(ctx context.Context, dbTx pgx.Tx) error {
	lastLocalL2BlockNumber, err := p.state.GetLastL2BlockNumber(ctx, dbTx)
	if errors.Is(err, state.ErrNotFound) || errors.Is(err, state.ErrStateNotSynchronized) {
		log.Debugf("checkL2block:No L2Block  in database. err: %s", err.Error())
		return nil
	}
	if err != nil {
		log.Errorf("checkL2block: Error getting last L2Block from the database. err: %s", err.Error())
		return err
	}
	shouldCheck, l2BlockNumber := p.GetNextL2BlockToCheck(lastLocalL2BlockNumber, p.GetMinimumL2BlockToCheck())
	if !shouldCheck {
		return nil
	}
	err = p.iterationCheckL2Block(ctx, l2BlockNumber, dbTx)
	if err != nil {
		return err
	}
	return nil
}

// GetNextL2BlockToCheck returns true is need to check and the blocknumber
func (p *CheckL2BlockHash) GetNextL2BlockToCheck(lastLocalL2BlockNumber, minL2BlockNumberToCheck uint64) (bool, uint64) {
	l2BlockNumber := max(minL2BlockNumberToCheck, lastLocalL2BlockNumber)
	if l2BlockNumber > lastLocalL2BlockNumber {
		log.Infof("checkL2block: skip check L2block (next to check: %d) currently LastL2BlockNumber: %d", minL2BlockNumberToCheck, lastLocalL2BlockNumber)
		return false, 0
	}
	if l2BlockNumber%p.modulusL2BlockToCheck != 0 {
		return false, 0
	}
	return true, l2BlockNumber
}

// GetMinimumL2BlockToCheck returns the minimum L2Block to check
func (p *CheckL2BlockHash) GetMinimumL2BlockToCheck() uint64 {
	if p.modulusL2BlockToCheck == 0 {
		return p.lastL2BlockChecked + 1
	}
	return ((p.lastL2BlockChecked / p.modulusL2BlockToCheck) + 1) * p.modulusL2BlockToCheck
}

// GetL2Blocks returns localL2Block and trustedL2Block
func (p *CheckL2BlockHash) GetL2Blocks(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*state.L2Block, *types.Block, error) {
	localL2Block, err := p.state.GetL2BlockByNumber(ctx, blockNumber, dbTx)
	if err != nil {
		log.Debugf("checkL2block: Error getting L2Block %d from the database. err: %s", blockNumber, err.Error())
		return nil, nil, err
	}
	trustedL2Block, err := p.trustedClient.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Errorf("checkL2block: Error getting L2Block %d from the Trusted RPC. err:%s", blockNumber, err.Error())
		return nil, nil, nil
	}
	return localL2Block, trustedL2Block, nil
}

// CheckPostClosedBatch checks the last L2Block hash on close batch
func (p *CheckL2BlockHash) iterationCheckL2Block(ctx context.Context, l2BlockNumber uint64, dbTx pgx.Tx) error {
	prefixLogs := fmt.Sprintf("checkL2block: L2BlockNumber: %d ", l2BlockNumber)
	localL2Block, trustedL2Block, err := p.GetL2Blocks(ctx, l2BlockNumber, dbTx)
	if errors.Is(err, state.ErrNotFound) || errors.Is(err, state.ErrStateNotSynchronized) {
		log.Debugf("%s not found in the database", prefixLogs, l2BlockNumber)
		return nil
	}
	if err != nil {
		log.Errorf("%s Error getting  from the database and trusted. err: %s", prefixLogs, err.Error())
		return err
	}
	if localL2Block == nil || trustedL2Block == nil {
		log.Errorf("%s localL2Block or trustedL2Block is nil", prefixLogs, l2BlockNumber)
		return nil
	}

	if err := compareL2Blocks(prefixLogs, localL2Block, trustedL2Block); err != nil {
		log.Errorf("%s Error comparing L2Blocks from the database and trusted. err: %s", prefixLogs, err.Error())
		return err
	}

	log.Infof("%s checked L2Block in the database  and the trusted batch are the same %s", prefixLogs, localL2Block.Hash().String())
	// Compare the two blocks
	p.lastL2BlockChecked = l2BlockNumber
	return nil
}

func compareL2Blocks(prefixLogs string, localL2Block *state.L2Block, trustedL2Block *types.Block) error {
	if localL2Block == nil || trustedL2Block == nil {
		return fmt.Errorf("%s localL2Block or trustedL2Block are nil", prefixLogs)
	}
	if localL2Block.Hash() != trustedL2Block.Hash() {
		return fmt.Errorf("%s localL2Block.Hash %s and trustedL2Block.Hash %s are different", prefixLogs, localL2Block.Hash().String(), trustedL2Block.Hash().String())
	}
	if localL2Block.ParentHash() != trustedL2Block.ParentHash() {
		return fmt.Errorf("%s localL2Block.ParentHash %s and trustedL2Block.ParentHash %s are different", prefixLogs, localL2Block.ParentHash().String(), trustedL2Block.ParentHash().String())
	}
	return nil
}
