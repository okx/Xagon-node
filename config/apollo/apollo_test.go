package apollo

import (
	"testing"

	nodeConfig "github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

func TestApolloClient_LoadConfig(t *testing.T) {
	nc := &nodeConfig.Config{
		Apollo: types.ApolloConfig{
			IP:            "",
			AppID:         "xlayer-devnet",
			NamespaceName: "jsonrpc-ro.txt,jsonrpc-roHalt.properties",
			Enable:        true,
		},
	}
	client := NewClient(nc)

	client.LoadConfig()
	t.Log(nc.RPC)
	// time.Sleep(2 * time.Minute)
	t.Log(nc.RPC)
}
