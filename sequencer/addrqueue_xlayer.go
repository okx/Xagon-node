package sequencer

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

// ForceQueueBalanceConfig is the global variable to force the queue balance enough
var ForceQueueBalanceConfig bool

func (a *addrQueue) checkQueueBalanceEnough(currentBalance *big.Int, cost *big.Int) bool {
	if getForceQueueBalanceEnoughConfig(ForceQueueBalanceConfig) {
		log.Debugf("checkQueueBalanceEnough, it always use true")
		return true
	}

	result := currentBalance.Cmp(cost) >= 0
	log.Debugf("checkQueueBalanceEnough: currentBalance=%s, cost=%s, result:%v", currentBalance.String(), cost.String(), result)
	return result
}

// SetForceQueueBalanceConfig sets the forceQueueBalanceEnough config
func SetForceQueueBalanceConfig(value bool) {
	log.Infof("SetForceQueueBalanceConfig form config: %v", value)
	ForceQueueBalanceConfig = value
}
