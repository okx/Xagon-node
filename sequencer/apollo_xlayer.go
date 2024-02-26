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
