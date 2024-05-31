package sequencer

import (
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// ApolloConfig is the apollo RPC dynamic config
type ApolloConfig struct {
	EnableApollo           bool
	FullBatchSleepDuration types.Duration
	PackBatchSpacialList   []string
	GasPriceMultiple       float64
	InitGasPriceMultiple   float64
	QueryPendingTxsLimit   uint64

	sync.RWMutex
}

var apolloConfig = &ApolloConfig{}

// getApolloConfig returns the singleton instance
func getApolloConfig() *ApolloConfig {
	return apolloConfig
}

// Enable returns true if apollo is enabled
func (c *ApolloConfig) Enable() bool {
	if c == nil || !c.EnableApollo {
		return false
	}
	c.RLock()
	defer c.RUnlock()
	return c.EnableApollo
}

// UpdateConfig updates the apollo config
func UpdateConfig(apolloConfig Config) {
	getApolloConfig().Lock()
	getApolloConfig().EnableApollo = true
	getApolloConfig().FullBatchSleepDuration = apolloConfig.Finalizer.FullBatchSleepDuration
	getApolloConfig().PackBatchSpacialList = apolloConfig.PackBatchSpacialList
	getApolloConfig().GasPriceMultiple = apolloConfig.GasPriceMultiple
	getApolloConfig().InitGasPriceMultiple = apolloConfig.InitGasPriceMultiple
	getApolloConfig().QueryPendingTxsLimit = apolloConfig.QueryPendingTxsLimit
	getApolloConfig().Unlock()
}

func getFullBatchSleepDuration(localDuration time.Duration) time.Duration {
	var ret time.Duration
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		ret = getApolloConfig().FullBatchSleepDuration.Duration
	} else {
		ret = localDuration
	}

	return ret
}

func getPackBatchSpacialList(addrs []string) map[string]bool {
	ret := make(map[string]bool, len(addrs))
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		addrs = getApolloConfig().PackBatchSpacialList
	}

	for _, addr := range addrs {
		ret[addr] = true
	}

	return ret
}

func getGasPriceMultiple(gpMul float64) float64 {
	ret := gpMul
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		ret = getApolloConfig().GasPriceMultiple
	}

	return ret
}

func getInitGasPriceMultiple(gpMul float64) float64 {
	ret := gpMul
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		ret = getApolloConfig().InitGasPriceMultiple
	}

	return ret
}

func getQueryPendingTxsLimit(limit uint64) uint64 {
	ret := limit
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		ret = getApolloConfig().QueryPendingTxsLimit
	}

	return ret
}
