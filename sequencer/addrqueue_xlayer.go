package sequencer

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	"math/big"
)

func (a *addrQueue) checkQueueBalanceEnough(currentBalance *big.Int, cost *big.Int) bool {
	//// TODO apollo check config
	//if true {
	//	log.Infof("checkQueueBalanceEnough: currentBalance=%s, cost=%s", currentBalance.String(), cost.String())
	//	return true
	//}

	result := currentBalance.Cmp(cost) >= 0
	log.Infof("checkQueueBalanceEnough: currentBalance=%s, cost=%s, result:%v", currentBalance.String(), cost.String(), result)
	return result
}
