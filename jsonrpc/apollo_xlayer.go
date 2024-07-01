package jsonrpc

import (
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/lru_xlayer"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

// ApolloConfig is the apollo RPC dynamic config
type ApolloConfig struct {
	EnableApollo         bool
	BatchRequestsEnabled bool
	BatchRequestsLimit   uint
	GasLimitFactor       float64
	EnableEstimateGasOpt bool
	DisableAPIs          []string
	RateLimit            RateLimitConfig
	DynamicGP            DynamicGPConfig
	ApiAuthentication    ApiAuthConfig
	ApiRelay             ApiRelayConfig

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

func (c *ApolloConfig) setApiRelayCfg(apiRelayCfg ApiRelayConfig) {
	if c == nil || !c.EnableApollo {
		return
	}
	c.ApiRelay.Enabled = apiRelayCfg.Enabled
	c.ApiRelay.DestURI = apiRelayCfg.DestURI
	c.ApiRelay.RPCs = make([]string, len(apiRelayCfg.RPCs))
	copy(c.ApiRelay.RPCs, apiRelayCfg.RPCs)
}

// UpdateConfig updates the apollo config
func UpdateConfig(apolloConfig Config) {
	getApolloConfig().Lock()
	getApolloConfig().EnableApollo = true
	getApolloConfig().BatchRequestsEnabled = apolloConfig.BatchRequestsEnabled
	getApolloConfig().BatchRequestsLimit = apolloConfig.BatchRequestsLimit
	getApolloConfig().GasLimitFactor = apolloConfig.GasLimitFactor
	getApolloConfig().EnableEstimateGasOpt = apolloConfig.EnableEstimateGasOpt
	getApolloConfig().setDisableAPIs(apolloConfig.DisableAPIs)
	setRateLimit(apolloConfig.RateLimit)
	setApiAuth(apolloConfig.ApiAuthentication)
	getApolloConfig().DynamicGP = apolloConfig.DynamicGP
	getApolloConfig().setApiRelayCfg(apolloConfig.ApiRelay)
	lru_xlayer.SetConfig(apolloConfig.LRUConfig)
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
