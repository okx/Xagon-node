package sequencer

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	pmetric "github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
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

func (s *Sequencer) checkFreeGas(tx pool.Transaction, txTracker *TxTracker) (freeGp, claimTx bool, gpMul float64) {
	// check if tx is init-free-gas and if it can be prior pack
	freeGp = tx.GasPrice().Cmp(big.NewInt(0)) == 0
	gpMul = getInitGasPriceMultiple(s.cfg.InitGasPriceMultiple)

	// check if tx is bridge-claim
	addrs := getPackBatchSpacialList(s.cfg.PackBatchSpacialList)
	if addrs[txTracker.FromStr] {
		claimTx = true
		gpMul = getGasPriceMultiple(s.cfg.GasPriceMultiple)
		return
	}

	// check if tx is from special project
	if pool.GetEnableSpecialFreeGasList(s.poolCfg.EnableFreeGasList) {
		fromToName, freeGpList := pool.GetSpecialFreeGasList(s.poolCfg.FreeGasList)
		info := freeGpList[fromToName[txTracker.FromStr]]
		if info != nil &&
			tx.To() != nil &&
			pool.Contains(info.ToList, *tx.To()) &&
			pool.ContainsMethod("0x"+common.Bytes2Hex(tx.Data()), info.MethodSigs) {
			gpMul = info.GasPriceMultiple
			return
		}
	}

	return
}
