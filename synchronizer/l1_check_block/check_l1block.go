package l1_check_block

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// This object check old L1block to double-check that the L1block hash is correct
// - Get first not checked block
// - Get last block on L1 (safe/finalized/ or minus -n)

// L1Requester is an interface for GETH client
type L1Requester interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// StateInterfacer is an interface for the state
type StateInterfacer interface {
	GetFirstUncheckedBlock(ctx context.Context, fromBlockNumber uint64, dbTx pgx.Tx) (*state.Block, error)
	UpdateCheckedBlockByNumber(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx pgx.Tx) error
}

// SafeL1BlockNumberFetcher is an interface for fetching the  L1 block number reference point (safe, finalized,...)
type SafeL1BlockNumberFetcher interface {
	GetSafeBlockNumber(ctx context.Context, l1Client L1Requester) (uint64, error)
	Description() string
}

// CheckL1BlockHash is a struct that implements a checker of L1Block hash
type CheckL1BlockHash struct {
	L1Client               L1Requester
	State                  StateInterfacer
	SafeBlockNumberFetcher SafeL1BlockNumberFetcher
}

// NewCheckL1BlockHash creates a new CheckL1BlockHash
func NewCheckL1BlockHash(l1Client L1Requester, state StateInterfacer, safeBlockNumberFetcher SafeL1BlockNumberFetcher) *CheckL1BlockHash {
	return &CheckL1BlockHash{
		L1Client:               l1Client,
		State:                  state,
		SafeBlockNumberFetcher: safeBlockNumberFetcher,
	}
}

// Name is a method that returns the name of the checker
func (p *CheckL1BlockHash) Name() string {
	return logPrefix + " main_checker: "
}

// Step is a method that checks the L1 block hash, run until all blocks are checked and returns
func (p *CheckL1BlockHash) Step(ctx context.Context) error {
	stateBlock, err := p.State.GetFirstUncheckedBlock(ctx, uint64(0), nil)
	if errors.Is(err, state.ErrNotFound) {
		log.Debugf("%s: No unchecked blocks to check", p.Name())
		return nil
	}
	if err != nil {
		return err
	}
	if stateBlock == nil {
		log.Warnf("%s: function CheckL1Block receive a nil pointer", p.Name())
		return nil
	}
	safeBlockNumber, err := p.SafeBlockNumberFetcher.GetSafeBlockNumber(ctx, p.L1Client)
	if err != nil {
		return err
	}
	log.Debugf("%s: checking from block (%s) %d first block to check: %d....", p.Name(), p.SafeBlockNumberFetcher.Description(), safeBlockNumber, stateBlock.BlockNumber)
	return p.doAllBlocks(ctx, *stateBlock, safeBlockNumber)
}

func (p *CheckL1BlockHash) doAllBlocks(ctx context.Context, firstStateBlock state.Block, safeBlockNumber uint64) error {
	var err error
	startTime := time.Now()
	stateBlock := &firstStateBlock
	numBlocksChecked := 0
	for {
		lastStateBlockNumber := stateBlock.BlockNumber
		if stateBlock.BlockNumber > safeBlockNumber {
			log.Debugf("%s: block %d to check is not still safe enough (%s) %d ", p.Name(), stateBlock.BlockNumber, p.SafeBlockNumberFetcher.Description(), safeBlockNumber, logPrefix)
			return nil
		}
		err = p.doBlock(ctx, stateBlock)
		if err != nil {
			return err
		}
		numBlocksChecked++
		stateBlock, err = p.State.GetFirstUncheckedBlock(ctx, lastStateBlockNumber, nil)
		if errors.Is(err, state.ErrNotFound) {
			diff := time.Since(startTime)
			log.Infof("%s: checked all blocks (%d) (using as safe Block Point(%s): %d) time:%s", p.Name(), numBlocksChecked, p.SafeBlockNumberFetcher.Description(), safeBlockNumber, diff)
			return nil
		}
	}
}

func (p *CheckL1BlockHash) doBlock(ctx context.Context, stateBlock *state.Block) error {
	err := CheckBlockHash(ctx, stateBlock, p.L1Client, p.Name())
	if err != nil {
		return err
	}
	log.Infof("%s: L1Block: %d hash: %s is correct marking as checked", p.Name(), stateBlock.BlockNumber,
		stateBlock.BlockHash.String())
	err = p.State.UpdateCheckedBlockByNumber(ctx, stateBlock.BlockNumber, true, nil)
	if err != nil {
		log.Errorf("%s: Error updating block %d as checked. err: %s", p.Name(), stateBlock.BlockNumber, err.Error())
		return err
	}
	return nil
}

// CheckBlockHash is a method that checks the L1 block hash
func CheckBlockHash(ctx context.Context, stateBlock *state.Block, L1Client L1Requester, checkerName string) error {
	if stateBlock == nil {
		log.Warn("%s function CheckL1Block receive a nil pointer", checkerName)
		return nil
	}
	l1Block, err := L1Client.HeaderByNumber(ctx, big.NewInt(int64(stateBlock.BlockNumber)))
	if err != nil {
		return err
	}
	if l1Block == nil {
		err = fmt.Errorf("%s request of block: %d to L1 returns a nil", checkerName, stateBlock.BlockNumber)
		log.Error(err.Error())
		return err
	}
	if l1Block.Hash() != stateBlock.BlockHash {
		msg := fmt.Sprintf("%s Reorg detected at block %d l1Block.Hash=%s != stateBlock.Hash=%s. ", checkerName, stateBlock.BlockNumber,
			l1Block.Hash().String(), stateBlock.BlockHash.String())
		if l1Block.ParentHash != stateBlock.ParentHash {
			msg += fmt.Sprintf(" ParentHash are also different. l1Block.ParentHash=%s !=  stateBlock.ParentHash=%s", l1Block.ParentHash.String(), stateBlock.ParentHash.String())
		}
		log.Errorf(msg)
		return common.NewReorgError(stateBlock.BlockNumber, fmt.Errorf(msg))
	}
	return nil
}
