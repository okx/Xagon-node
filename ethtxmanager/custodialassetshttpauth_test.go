package ethtxmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientGenerateSignature(t *testing.T) {
	client := &Client{
		cfg: Config{
			CustodialAssets: CustodialAssetsConfig{
				SecretKey: "12doxpwjkengkjna",
			},
		},
	}
	ctx := context.Background()
	treeMap := make(map[string][]string)
	treeMap["0"] = []string{"0x07f67d4195bc9940f07eb901ef18f1e9e4af12d7"}
	treeMap["1"] = []string{"127"}
	treeMap["2"] = []string{"true"}
	auth, err := client.generateSignature(ctx, treeMap, "{\"testBOdy\":45251}", algoSha256)
	require.NoError(t, err)
	require.Equal(t, "si/fTWlDg6+V9OFOM3CictCuqtGfUjKZ3keGLwxM/walrXtQaN8K/PnGTvFvc4q6pb/80HtZIy+hjeugAx8VPLmIXmKpJ5H4mbGEVQe7bk4=", auth)

	auth, err = client.generateSignature(ctx, treeMap, "{\"testBOdy\":45251}", algoMd5)
	require.NoError(t, err)
	require.Equal(t, "r1HPNFEpExR61IOWf8mP5oOxbz4RnyU0LF67Dv3hlBWuXy6bJA5x467NbVTldTVOdoAB/16w1cc+7/Y5fDEcC7mIXmKpJ5H4mbGEVQe7bk4=", auth)
}
