package l1_check_block

// This make a pre-check of blocks but don't mark them as checked
// It checks blocks between a segment: example:
// real check point SAFE:
//  pre check:  (SAFE+1) -> (LATEST-32)
// It gets all pending blocks
//  - Start cheking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/jackc/pgx/v4"
)

var (
	// ErrDeSync is an error that indicates that from the starting of verification to end something have been changed on state
	ErrDeSync = errors.New("DeSync: a block hash is different from the state block hash")
)

// StatePreCheckInterfacer is an interface for the state
type StatePreCheckInterfacer interface {
	GetUncheckedBlocks(ctx context.Context, fromBlockNumber uint64, toBlockNumber uint64, dbTx pgx.Tx) ([]*state.Block, error)
}

// PreCheckL1BlockHash is a struct that implements a checker of L1Block hash
type PreCheckL1BlockHash struct {
	L1Client                  L1Requester
	State                     StatePreCheckInterfacer
	InitialSegmentBlockNumber SafeL1BlockNumberFetcher
	EndSegmentBlockNumber     SafeL1BlockNumberFetcher
}

// NewPreCheckL1BlockHash creates a new CheckL1BlockHash
func NewPreCheckL1BlockHash(l1Client L1Requester, state StatePreCheckInterfacer,
	initial, end SafeL1BlockNumberFetcher) *PreCheckL1BlockHash {
	return &PreCheckL1BlockHash{
		L1Client:                  l1Client,
		State:                     state,
		InitialSegmentBlockNumber: initial,
		EndSegmentBlockNumber:     end,
	}
}

// Name is a method that returns the name of the checker
func (p *PreCheckL1BlockHash) Name() string {
	return logPrefix + ":memory_check: "
}

// Step is a method that checks the L1 block hash, run until all blocks are checked and returns
func (p *PreCheckL1BlockHash) Step(ctx context.Context) error {
	from, err := p.InitialSegmentBlockNumber.GetSafeBlockNumber(ctx, p.L1Client)
	if err != nil {
		return err
	}
	to, err := p.EndSegmentBlockNumber.GetSafeBlockNumber(ctx, p.L1Client)
	if err != nil {
		return err
	}
	if from > to {
		log.Warnf("%s: fromBlockNumber(%s) %d is greater than toBlockNumber(%s) %d, Check configuration", p.Name(), p.InitialSegmentBlockNumber.Description(), from, p.EndSegmentBlockNumber.Description(), to)
		return nil
	}

	blocksToCheck, err := p.State.GetUncheckedBlocks(ctx, from, to, nil)
	if err != nil {
		log.Warnf("%s can't get unchecked blocks, so it discard the reorg error", p.Name())
		return err
	}
	msg := fmt.Sprintf("%s: Checking blocks from (%s) %d to (%s) %d -> len(blocks)=%d", p.Name(), p.InitialSegmentBlockNumber.Description(), from, p.EndSegmentBlockNumber.Description(), to, len(blocksToCheck))
	if len(blocksToCheck) == 0 {
		log.Debugf(msg)
		return nil
	}
	log.Infof(msg)
	startTime := time.Now()
	for _, block := range blocksToCheck {
		// check block
		err = CheckBlockHash(ctx, block, p.L1Client, p.Name())
		if common.IsReorgError(err) {
			// Double-check the state block that still is the same
			log.Debugf("%s: Reorg detected at blockNumber: %d, checking that the block on State doesn't have change", p.Name(), block.BlockNumber)
			isTheSame, errBlockIsTheSame := p.checkThatStateBlockIsTheSame(ctx, block)
			if errBlockIsTheSame != nil {
				log.Warnf("%s can't double-check that blockNumber %d haven't changed, so it discard the reorg error", p.Name(), block.BlockNumber)
				return err
			}
			if !isTheSame {
				log.Infof("%s: DeSync detected, blockNumber: %d is different now that when we started the check", p.Name(), block.BlockNumber)
				return ErrDeSync
			}
			log.Infof("%s: Reorg detected and verified the state block, blockNumber: %d", p.Name(), block.BlockNumber)
			return err
		}
		if err != nil {
			return err
		}
	}
	elapsed := time.Since(startTime)
	log.Infof("%s: Checked blocks from (%s) %d to (%s) %d -> len(blocks):%d elapsed: %s", p.Name(), p.InitialSegmentBlockNumber.Description(), from, p.EndSegmentBlockNumber.Description(), to, len(blocksToCheck), elapsed.String())

	return nil
}

// CheckBlockHash is a method that checks the L1 block hash
// returns true if is the same
func (p *PreCheckL1BlockHash) checkThatStateBlockIsTheSame(ctx context.Context, block *state.Block) (bool, error) {
	blocks, err := p.State.GetUncheckedBlocks(ctx, block.BlockNumber, block.BlockNumber, nil)
	if err != nil {
		log.Warnf("%s: Fails to get blockNumber %d in state .Err:%s", p.Name(), block.BlockNumber, err.Error())
		return false, err
	}
	if len(blocks) == 0 {
		// The block is checked or deleted, so it is not the same
		log.Debugf("%s: The blockNumber %d is no longer in the state (or checked or deleted)", p.Name(), block.BlockNumber)
		return false, nil
	}
	stateBlock := blocks[0]
	if stateBlock.BlockNumber != block.BlockNumber {
		msg := fmt.Sprintf("%s: The blockNumber returned by state %d is different from the state blockNumber %d",
			p.Name(), block.BlockNumber, stateBlock.BlockNumber)
		log.Warn(msg)
		return false, fmt.Errorf(msg)
	}
	if stateBlock.BlockHash != block.BlockHash {
		msg := fmt.Sprintf("%s: The blockNumber %d differs the hash checked %s from current in state %s",
			p.Name(), block.BlockNumber, block.BlockHash.String(), stateBlock.BlockHash.String())
		log.Warn(msg)
		return false, nil
	}
	// The block is the same
	return true, nil
}
