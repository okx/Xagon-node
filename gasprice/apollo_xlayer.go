package gasprice

import "sync"

// ApolloConfig is the apollo RPC dynamic config
type ApolloConfig struct {
	EnableApollo bool
	conf         Config
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
	getApolloConfig().conf = apolloConfig
	getApolloConfig().Unlock()
}

func (c *ApolloConfig) get() Config {
	c.RLock()
	defer c.RUnlock()
	return c.conf
}
