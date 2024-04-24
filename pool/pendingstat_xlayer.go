package pool

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool/trace"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

// PendingStatCfg is the configuration for the pending stat
type PendingStatCfg struct {
	Enable        bool           `mapstructure:"Enable"`
	Interval      types.Duration `mapstructure:"Interval"`
	StaleInterval types.Duration `mapstructure:"StaleInterval"`
	CacheInternal types.Duration `mapstructure:"CacheInternal"`
}

func (p *Pool) startPendingStat() {
	if p.cfg.PendingStat.Enable {
		go state.InfiniteSafeRun(p.updatePendingStat, "error updating pending stat", p.cfg.PendingStat.Interval.Duration)
		go state.InfiniteSafeRun(p.updatePendingStatCache, "error updating pending stat cache", p.cfg.PendingStat.CacheInternal.Duration)
	}
}

// updatePendingStat updates the pending statistics
// 1. find all pending transactions count
// 2. find all pending address and min nonce received before the stale interval
// 3. find all pending address that skip nonce
// 4. find all pending address that have balance issue
func (p *Pool) updatePendingStat() {
	if !getEnablePendingStat(p.cfg.PendingStat.Enable) {
		return
	}
	ctx := context.WithValue(context.Background(), trace.ID, uuid.New().String())
	locked, err := p.storage.LockStat(ctx, p.cfg.PendingStat.Interval.Duration)
	if err != nil {
		return
	}
	if !locked {
		return
	}
	defer func() {
		err = p.storage.UnLockStat(ctx)
		if err != nil {
			log.WithFields(trace.GetID(ctx)).Error("error unlocking stat", "err", err)
		}
	}()

	mLog := log.WithFields(trace.GetID(ctx))
	mLog.Infof("updating pending stat")

	totalCount, err := p.CountTransactionsByStatus(ctx, TxStatusPending)
	if err != nil {
		mLog.Error("error getting pending transactions count", "err", err)
		return
	}
	mLog.Infof("total pending transactions %v", totalCount)

	address, nonces, err := p.storage.GetPendingFromAndMinNonceBefore(ctx, p.cfg.PendingStat.StaleInterval.Duration)
	if err != nil {
		mLog.Error("error getting pending address and min nonce before ", err)
		return
	}
	mLog.Infof("pending address count %v before %v", len(address), p.cfg.PendingStat.StaleInterval.Duration)

	skipNonceAddress, continueNonceAddress, continueNonces, errNonceAddress, noncesNormal, err := p.filterAddress(ctx, address, nonces)
	if err != nil {
		mLog.Error("error filtering skip nonce address", "err", err)
		return
	}
	mLog.Infof("skip nonce address count %v, continue nonce address count %v, err nonce address count %v", len(skipNonceAddress), len(continueNonceAddress), len(errNonceAddress))

	totalSkipNonceTransactions, err := p.countSkipNonceTransactions(ctx, skipNonceAddress)
	if err != nil {
		mLog.Error("error counting skip nonce transactions", "err", err)
		return
	}
	mLog.Infof("total skip nonce transactions %v", totalSkipNonceTransactions)

	totalBalanceIssueTransactions, err := p.countBalanceIssueTransactions(ctx, continueNonceAddress, continueNonces)
	if err != nil {
		mLog.Error("error counting balance issue transactions", "err", err)
		return
	}
	mLog.Infof("total balance issue transactions %v", totalBalanceIssueTransactions)

	totalErrNonceTransactions, err := p.countErrNonceTransactions(ctx, errNonceAddress, noncesNormal)
	if err != nil {
		mLog.Error("error counting nonce issue transactions", "err", err)
		return
	}
	mLog.Infof("total nonce issue transactions %v", totalErrNonceTransactions)

	err = p.storage.UpdateStatAndUnlock(ctx, totalCount, totalSkipNonceTransactions, totalBalanceIssueTransactions, totalErrNonceTransactions)
	if err != nil {
		mLog.Error("error updating stat and unlock", "err", err)
	} else {
		mLog.Infof("total %v, skip nonce %v, balance issue %v, nonce issue %v", totalCount, totalSkipNonceTransactions, totalBalanceIssueTransactions, totalErrNonceTransactions)
	}
}

func (p *Pool) filterAddress(ctx context.Context, addresses []common.Address, nonces []uint64) ([]common.Address, []common.Address, []uint64, []common.Address, []uint64, error) {
	var skipNonceAddresses []common.Address
	var continueNonceAddresses []common.Address
	var continueNonces []uint64
	var errNonceAddresses []common.Address
	var noncesNormal []uint64

	lastL2Block, err := p.state.GetLastL2Block(ctx, nil)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to load last l2 block while adding tx to the pool: %w", err)
	}
	for i, addr := range addresses {
		nonce, err := p.state.GetNonce(ctx, addr, lastL2Block.Root())
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("failed to load nonce while adding tx to the pool: %w", err)
		}
		if nonces[i]-nonce > 0 {
			skipNonceAddresses = append(skipNonceAddresses, addr)
		} else if (nonces[i] - nonce) == 0 {
			continueNonceAddresses = append(continueNonceAddresses, addr)
			continueNonces = append(continueNonces, nonces[i])
		} else {
			errNonceAddresses = append(errNonceAddresses, addr)
			noncesNormal = append(noncesNormal, nonce)
		}
	}
	return skipNonceAddresses, continueNonceAddresses, continueNonces, errNonceAddresses, noncesNormal, nil
}

func (p *Pool) countSkipNonceTransactions(ctx context.Context, addresses []common.Address) (uint64, error) {
	var totalSkipNonceTransactions uint64
	for _, addr := range addresses {
		count, err := p.storage.CountTransactionsByFromAndStatus(ctx, addr, TxStatusPending)
		if err != nil {
			return 0, err
		}
		totalSkipNonceTransactions += count
	}
	return totalSkipNonceTransactions, nil
}

func (p *Pool) countBalanceIssueTransactions(ctx context.Context, addresses []common.Address, nonces []uint64) (uint64, error) {
	var totalBalanceIssueTransactions uint64
	mLog := log.WithFields(trace.GetID(ctx))

	lastL2Block, err := p.state.GetLastL2Block(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to load last l2 block while adding tx to the pool: %w", err)
	}

	for i, addr := range addresses {
		txs, err := p.storage.GetTxsByFromAndNonce(ctx, addr, nonces[i])
		if err != nil || len(txs) > 1 {
			mLog.Warnf("error getting transactions by from %v and nonce %v tx count %v, err %v", addr, nonces[i], len(txs), err)
			continue
		}
		balance, err := p.state.GetBalance(ctx, addr, lastL2Block.Root())
		if err != nil {
			mLog.Warnf("error getting balance for address %v, l2block %v, err %v", addr, err, lastL2Block.Root())
			continue
		}
		if balance.Cmp(txs[0].Cost()) < 0 {
			count, err := p.storage.CountTransactionsByFromAndStatus(ctx, addr, TxStatusPending)
			if err != nil {
				mLog.Warnf("error getting transactions count by from %v, err %v", addr, err)
			}
			totalBalanceIssueTransactions += count
		}
	}

	return totalBalanceIssueTransactions, nil
}

func (p *Pool) countErrNonceTransactions(ctx context.Context, addresses []common.Address, nonces []uint64) (uint64, error) {
	var totalErrNonceTransactions uint64
	for i, addr := range addresses {
		count, err := p.storage.CountTransactionsByFromStatusAndNonce(ctx, addr, nonces[i], TxStatusPending)
		if err != nil {
			return 0, err
		}
		totalErrNonceTransactions += count
	}
	return totalErrNonceTransactions, nil
}
