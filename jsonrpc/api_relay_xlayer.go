package jsonrpc

import (
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

type ApiRelayConfig struct {
	Enabled bool     `mapstructure:"Enabled"`
	DestURI string   `mapstructure:"DestURI"`
	RPCs    []string `mapstructure:"RPCs"`
}

func shouldRelay(localCfg ApiRelayConfig, name string) bool {
	enable := localCfg.Enabled && localCfg.DestURI != ""
	contained := types.Contains(localCfg.RPCs, name)
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		enable = getApolloConfig().ApiRelay.Enabled && getApolloConfig().ApiRelay.DestURI != ""
		contained = types.Contains(getApolloConfig().ApiRelay.RPCs, name)
	}

	return enable && contained
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
