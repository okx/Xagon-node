package sequencer

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	pmetric "github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

var countinterval = 10

func (s *Sequencer) countPendingTx() {
	for {
		<-time.After(time.Second * time.Duration(countinterval))
		transactions, err := s.pool.CountPendingTransactions(context.Background())
		if err != nil {
			log.Errorf("load pending tx from pool: %v", err)
			continue
		}
		pmetric.PendingTxCount(int(transactions))
	}
}

func (s *Sequencer) updateReadyTxCount() {
	err := s.pool.UpdateReadyTxCount(context.Background(), getPoolReadyTxCounter().getReadyTxCount())
	if err != nil {
		log.Errorf("error adding ready tx count: %v", err)
	}
}

func (s *Sequencer) countReadyTx() {
	state.InfiniteSafeRun(s.updateReadyTxCount, "error counting ready tx", time.Second)
}
