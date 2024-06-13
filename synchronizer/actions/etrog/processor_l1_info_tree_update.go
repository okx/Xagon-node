package etrog

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

// stateProcessorL1InfoTreeInterface interface required from state
type stateProcessorL1InfoTreeInterface interface {
	AddL1InfoTreeLeaf(ctx context.Context, L1InfoTreeLeaf *state.L1InfoTreeLeaf, dbTx pgx.Tx) (*state.L1InfoTreeExitRootStorageEntry, error)
}

// ProcessorL1InfoTreeUpdate implements L1EventProcessor for GlobalExitRootsOrder
type ProcessorL1InfoTreeUpdate struct {
	actions.ProcessorBase[ProcessorL1InfoTreeUpdate]
	state stateProcessorL1InfoTreeInterface
}

// NewProcessorL1InfoTreeUpdate new processor for GlobalExitRootsOrder
func NewProcessorL1InfoTreeUpdate(state stateProcessorL1InfoTreeInterface) *ProcessorL1InfoTreeUpdate {
	return &ProcessorL1InfoTreeUpdate{
		ProcessorBase: *actions.NewProcessorBase[ProcessorL1InfoTreeUpdate](
			[]etherman.EventOrder{etherman.L1InfoTreeOrder},
			actions.ForksIdToElderberry),
		state: state}
}

// Process process event
func (p *ProcessorL1InfoTreeUpdate) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	l1InfoTree := l1Block.L1InfoTree[order.Pos]
	ger := state.GlobalExitRoot{
		BlockNumber:     l1InfoTree.BlockNumber,
		MainnetExitRoot: l1InfoTree.MainnetExitRoot,
		RollupExitRoot:  l1InfoTree.RollupExitRoot,
		GlobalExitRoot:  l1InfoTree.GlobalExitRoot,
		Timestamp:       l1InfoTree.Timestamp,
	}
	l1IntoTreeLeaf := state.L1InfoTreeLeaf{
		GlobalExitRoot:    ger,
		PreviousBlockHash: l1InfoTree.PreviousBlockHash,
	}
	entry, err := p.state.AddL1InfoTreeLeaf(ctx, &l1IntoTreeLeaf, dbTx)
	if err != nil {
		log.Errorf("error storing the l1InfoTree(etrog). BlockNumber: %d, error: %v", l1Block.BlockNumber, err)
		return err
	}
	log.Infof("L1InfoTree(etrog) stored. BlockNumber: %d,GER:%s L1InfoTreeIndex: %d L1InfoRoot:%s", l1Block.BlockNumber, entry.GlobalExitRoot.GlobalExitRoot, entry.L1InfoTreeIndex, entry.L1InfoTreeRoot)
	return nil
}
