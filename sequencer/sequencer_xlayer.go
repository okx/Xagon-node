package sequencer

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	pmetric "github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

var countinterval = 10

func contains(s []string, ele common.Address) bool {
	for _, e := range s {
		if common.HexToAddress(e) == ele {
			return true
		}
	}
	return false
}

func containsMethod(data string, methods []string) bool {
	for _, m := range methods {
		if strings.HasPrefix(data, m) {
			return true
		}
	}
	return false
}

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
	}

	// check if tx is from special project
	if pool.GetEnableSpecialFreeGasList(s.poolCfg.EnableFreeGasList) {
		fromToName, freeGpList := pool.GetSpecialFreeGasList(s.poolCfg.FreeGasList)
		info := freeGpList[fromToName[txTracker.FromStr]]
		if info != nil &&
			contains(info.ToList, *tx.To()) &&
			containsMethod("0x"+common.Bytes2Hex(tx.Data()), info.MethodSigs) {
			gpMul = info.GasPriceMultiple
			return
		}
	}

	return
}
