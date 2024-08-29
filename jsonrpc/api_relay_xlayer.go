package jsonrpc

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/metrics"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// ApiRelayConfig is the api relay config
type ApiRelayConfig struct {
	Enabled  bool     `mapstructure:"Enabled"`
	DestURI  string   `mapstructure:"DestURI"`
	RPCs     []string `mapstructure:"RPCs"`
	Rerun    bool     `mapstructure:"Rerun"`
	RelayAll bool     `mapstructure:"RelayAll"`
}

func shouldRelayByMethod(localCfg ApiRelayConfig, name string) bool {
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

func shouldRerun(localCfg ApiRelayConfig) bool {
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		return getApolloConfig().ApiRelay.Rerun
	}
	return localCfg.Rerun
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
	if shouldRelay(localCfg, &request) {
		destURI := getRelayDestURI(localCfg.DestURI)
		ts := time.Now()
		res, err := client.JSONRPCRelay(destURI, request, shouldRerun(localCfg))
		if err != nil {
			log.Errorf("failed to relay %v to %s: %v", request.Method, destURI, err)
			metrics.RequestRelayFailCount(request.Method)
			return types.Response{}, false
		}
		log.Infof("relayed %v to %s in %v", request.Method, destURI, time.Since(ts))

		return res, true
	}

	return types.Response{}, false
}

func shouldRelay(localCfg ApiRelayConfig, request *types.Request) bool {
	enable := localCfg.Enabled &&
		localCfg.DestURI != ""
	relayAll := localCfg.RelayAll
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()
		enable = getApolloConfig().ApiRelay.Enabled &&
			getApolloConfig().ApiRelay.DestURI != ""
		relayAll = getApolloConfig().ApiRelay.RelayAll
	}
	if enable && relayAll {
		return true
	}

	return shouldRelayByMethod(localCfg, request.Method) && pass(request)
}
