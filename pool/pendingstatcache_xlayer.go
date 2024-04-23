package pool

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool/trace"
	"github.com/google/uuid"
)

// PendingStat is the pending stat
type PendingStat struct {
	// Total is the total number of pending transactions
	Total uint64 `json:"total"`

	// SkipNonce is the number of transactions that skipped nonce
	SkipNonce uint64 `json:"skipNonce"`

	// BalanceIssue is the number of transactions that have balance issue
	BalanceIssue uint64 `json:"balanceIssue"`

	// ErrorNonce is the number of transactions that have nonce issue
	ErrorNonce uint64 `json:"errorNonce"`
}

var pendingStatInst *PendingStat
var pendingStatOnce sync.Once

// GetPendingStat returns the singleton instance
func GetPendingStat() *PendingStat {
	pendingStatOnce.Do(func() {
		pendingStatInst = &PendingStat{}
	})
	return pendingStatInst
}

func (ps *PendingStat) setStat(stat PendingStat) {
	atomic.StoreUint64(&ps.Total, stat.Total)
	atomic.StoreUint64(&ps.SkipNonce, stat.SkipNonce)
	atomic.StoreUint64(&ps.BalanceIssue, stat.BalanceIssue)
}

// GetStat returns the pending stat
func (ps *PendingStat) GetStat() PendingStat {
	return PendingStat{
		Total:        atomic.LoadUint64(&ps.Total),
		SkipNonce:    atomic.LoadUint64(&ps.SkipNonce),
		BalanceIssue: atomic.LoadUint64(&ps.BalanceIssue),
	}
}

func (p *Pool) updatePendingStatCache() {
	if !getEnablePendingStat(p.cfg.PendingStat.Enable) {
		return
	}
	ctx := context.WithValue(context.Background(), trace.ID, uuid.New().String())
	mLog := log.WithFields(trace.GetID(ctx))
	total, skip, balanceIssue, nonceIssue, err := p.storage.GetStat(ctx)
	if err != nil {
		mLog.Errorf("error getting stat: %v", err)
		return
	}
	GetPendingStat().setStat(PendingStat{
		Total:        total,
		SkipNonce:    skip,
		BalanceIssue: balanceIssue,
		ErrorNonce:   nonceIssue,
	})
}
