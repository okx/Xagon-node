package jsonrpc

import "strings"

const (
	noSubCall = ""
)

type rpcMethod struct {
	Name    string   `mapstructure:"Name"`
	SubCall []string `mapstructure:"SubCall"`
}

type apiRelayConfig struct {
	Enabled bool        `mapstructure:"Enabled"`
	DestURI string      `mapstructure:"DestURI"`
	RPCs    []rpcMethod `mapstructure:"RPCs"`
}

func (c *apiRelayConfig) shouldRelay(name, subCall string) bool {
	if c == nil || !c.Enabled || c.DestURI == "" {
		return false
	}
	for _, rpc := range c.RPCs {
		if strings.ToLower(rpc.Name) == strings.ToLower(name) {
			if len(rpc.SubCall) == 0 {
				return true
			}
			for _, sc := range rpc.SubCall {
				if strings.ToLower(sc) == strings.ToLower(subCall) {
					return true
				}
			}
		}
	}

	return false
}

func (e *EthEndpoints) shouldRelay(name string, subCall string) bool {
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()

		return getApolloConfig().ApiRelayCfg.shouldRelay(name, subCall)
	}

	return e.cfg.ApiRelayCfg.shouldRelay(name, subCall)
}

func getRelayDestURI(localDestURI string) string {
	ret := localDestURI
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		defer getApolloConfig().RUnlock()

		ret = getApolloConfig().ApiRelayCfg.DestURI
	}

	return ret
}
