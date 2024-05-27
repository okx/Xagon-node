package jsonrpc

import (
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

type apiRelayConfig struct {
	Enabled bool     `mapstructure:"Enabled"`
	DestURI string   `mapstructure:"DestURI"`
	RPCs    []string `mapstructure:"RPCs"`
}

func (e *EthEndpoints) shouldRelay(name string) bool {
	if !e.cfg.ApiRelay.Enabled || e.cfg.ApiRelay.DestURI == "" {
		return false
	}

	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()

		return types.Contains(getApolloConfig().ApiRelay.RPCs, name)
	}

	return types.Contains(e.cfg.ApiRelay.RPCs, name)
}

func getRelayDestURI(localDestURI string) string {
	ret := localDestURI
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()

		ret = getApolloConfig().ApiRelay.DestURI
	}

	return ret
}
