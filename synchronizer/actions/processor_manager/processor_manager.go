package processor_manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

var (
	// ErrCantProcessThisEvent is used when the object is not found
	ErrCantProcessThisEvent = errors.New("not a processor for this event/forkid")
)

// L1EventProcessors is a manager of L1EventProcessor, it have processor for each forkId and event
//
//	  and it could:
//		 	- Returns specific processor for a forkId and event (Get function)
//	  	- Execute a event for a forkId and event (Process function)
//
// To build the object use L1EventProcessorsBuilder
type L1EventProcessors struct {
	// forkId -> event -> processor
	processors map[actions.ForkIdType]map[etherman.EventOrder]actions.L1EventProcessor
}

// NewL1EventProcessors returns a empty new L1EventProcessors
func NewL1EventProcessors() *L1EventProcessors {
	return &L1EventProcessors{
		processors: make(map[actions.ForkIdType]map[etherman.EventOrder]actions.L1EventProcessor),
	}
}

// Get returns the processor, first try specific, if not wildcard and if not found returns nil
func (p *L1EventProcessors) Get(forkId actions.ForkIdType, event etherman.EventOrder) actions.L1EventProcessor {
	if _, ok := p.processors[forkId]; !ok {
		if forkId == actions.WildcardForkId {
			return nil
		}
		return p.Get(actions.WildcardForkId, event)
	}
	if _, ok := p.processors[forkId][event]; !ok {
		if forkId == actions.WildcardForkId {
			return nil
		}
		return p.Get(actions.WildcardForkId, event)
	}
	return p.processors[forkId][event]
}

// Process execute the event for the forkId and event
func (p *L1EventProcessors) Process(ctx context.Context, forkId actions.ForkIdType, order etherman.Order, block *etherman.Block, dbTx pgx.Tx) error {
	processor := p.Get(forkId, order.Name)
	if processor == nil {
		var strBlockNumber string
		if block != nil {
			strBlockNumber = fmt.Sprintf("%d", block.BlockNumber)
		} else {
			strBlockNumber = "nil"
		}
		return fmt.Errorf("can't process blocknumber:%s event:%s, forkid:%d because: %w", strBlockNumber, order.Name, forkId, ErrCantProcessThisEvent)
	}
	if block.BlockNumber == 5733221 {
		log.Infof("Processing block %d, block time : %d", block.BlockNumber, block.ReceivedAt.Unix())
	}
	return processor.Process(ctx, order, block, dbTx)
}
