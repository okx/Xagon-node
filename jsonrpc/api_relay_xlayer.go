package jsonrpc

import (
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/metrics"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// ApiRelayConfig is the api relay config
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

func tryRelay(localCfg ApiRelayConfig, request types.Request) (types.Response, bool) {
	if shouldRelay(localCfg, request.Method) && pass(&request) {
		destURI := getRelayDestURI(localCfg.DestURI)
		res, err := client.JSONRPCRelay(destURI, request)
		if err != nil {
			log.Errorf("failed to relay %v to %s: %v", request.Method, destURI, err)
			metrics.RequestRelayFailCount(request.Method)
			return types.Response{}, false
		}

		return res, true
	}

	return types.Response{}, false
}
