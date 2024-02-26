package jsonrpc

import (
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

// ApolloConfig is the apollo RPC dynamic config
type ApolloConfig struct {
	EnableApollo         bool            `json:"enable"`
	BatchRequestsEnabled bool            `json:"batchRequestsEnabled"`
	BatchRequestsLimit   uint            `json:"batchRequestsLimit"`
	GasLimitFactor       float64         `json:"gasLimitFactor"`
	DisableAPIs          []string        `json:"disableAPIs"`
	RateLimit            RateLimitConfig `json:"rateLimit"`

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

func (c *ApolloConfig) setDisableAPIs(disableAPIs []string) {
	if c == nil || !c.EnableApollo {
		return
	}
	c.DisableAPIs = make([]string, len(disableAPIs))
	copy(c.DisableAPIs, disableAPIs)
}

// UpdateConfig updates the apollo config
func UpdateConfig(apolloConfig Config) {
	getApolloConfig().Lock()
	getApolloConfig().EnableApollo = true
	getApolloConfig().BatchRequestsEnabled = apolloConfig.BatchRequestsEnabled
	getApolloConfig().BatchRequestsLimit = apolloConfig.BatchRequestsLimit
	getApolloConfig().GasLimitFactor = apolloConfig.GasLimitFactor
	getApolloConfig().setDisableAPIs(apolloConfig.DisableAPIs)
	setRateLimit(apolloConfig.RateLimit)
	getApolloConfig().Unlock()
}

func (e *EthEndpoints) isDisabled(rpc string) bool {
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return len(getApolloConfig().DisableAPIs) > 0 && types.Contains(getApolloConfig().DisableAPIs, rpc)
	}

	return len(e.cfg.DisableAPIs) > 0 && types.Contains(e.cfg.DisableAPIs, rpc)
}
