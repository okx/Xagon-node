package config

import (
	"flag"
	"os"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestCardona(t *testing.T) {
	cfg := Config{}
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.String(FlagNetwork, string(cardona), string(cardona))
	err := fs.Set(FlagNetwork, string(cardona))
	require.NoError(t, err)
	app := cli.NewApp()
	ctx := cli.NewContext(app, fs, nil)

	log.Info("flag=", ctx.String(FlagNetwork))
	cfg.loadNetworkConfig(ctx)
}
func TestLoadCustomNetworkConfig(t *testing.T) {
	tcs := []struct {
		description      string
		inputConfigStr   string
		expectedConfig   NetworkConfig
		expectedErrorMsg string
	}{
		{
			description: "happy path",
			inputConfigStr: `{
				"root": "0xBEEF",
				"rollupCreationBlockNumber": 69,
				"rollupManagerCreationBlockNumber": 60,
				"l1Config" : {
					"chainId": 420,
					"polygonZkEVMAddress": "0xc949254d682d8c9ad5682521675b8f43b102aec4",
					"polTokenAddress": "0xc949254d682d8c9ad5682521675b8f43b102aec4",
					"polygonZkEVMGlobalExitRootAddress": "0xc949254d682d8c9ad5682521675b8f43b102aec4"
				},
				"genesis": [
					{
						"balance": "0",
						"nonce": "2",
						"address": "0xc949254d682d8c9ad5682521675b8f43b102aec4"
					},
					{
						"balance": "0",
						"nonce": "1",
						"address": "0xae4bb80be56b819606589de61d5ec3b522eeb032",
						"bytecode": "0xbeef1",
						"storage": {
							"0x0000000000000000000000000000000000000000000000000000000000000002": "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"
						},
						"contractName": "PolygonZkEVMGlobalExitRootL2 proxy"
					},
					{
						"balance": "100000000000000000000000",
						"nonce": "2",
						"address": "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988",
						"bytecode": "0xbeef2",
						"storage": {
							"0x0000000000000000000000000000000000000000000000000000000000000000": "0xc949254d682d8c9ad5682521675b8f43b102aec4"
						},
						"contractName": "PolygonZkEVMBridge proxy"
					},
					{
						"balance": "0",
						"nonce": "1",
						"address": "0x61ba0248b0986c2480181c6e76b6adeeaa962483",
						"bytecode": "0xbeef3",
						"storage": {
							"0x0000000000000000000000000000000000000000000000000000000000000000": "0x01"
						}
					}
				]
			}`,
			expectedConfig: NetworkConfig{
				L1Config: etherman.L1Config{
					L1ChainID:                 420,
					ZkEVMAddr:                 common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"),
					PolAddr:                   common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"),
					GlobalExitRootManagerAddr: common.HexToAddress("0xc949254d682d8c9ad5682521675b8f43b102aec4"),
				},
				L2BridgeAddr: common.HexToAddress("0x9D98DeAbC42dd696Deb9e40b4f1CAB7dDBF55988"),
				Genesis: state.Genesis{
					Root:                     common.HexToHash("0xBEEF"),
					RollupBlockNumber:        69,
					RollupManagerBlockNumber: 60,
					Actions: []*state.GenesisAction{
						{
							Address: "0xc949254d682d8c9ad5682521675b8f43b102aec4",
							Type:    int(merkletree.LeafTypeNonce),
							Value:   "2",
						},
						{
							Address: "0xae4bb80be56b819606589de61d5ec3b522eeb032",
							Type:    int(merkletree.LeafTypeNonce),
							Value:   "1",
						},
						{
							Address:  "0xae4bb80be56b819606589de61d5ec3b522eeb032",
							Type:     int(merkletree.LeafTypeCode),
							Bytecode: "0xbeef1",
						},
						{
							Address:         "0xae4bb80be56b819606589de61d5ec3b522eeb032",
							Type:            int(merkletree.LeafTypeStorage),
							StoragePosition: "0x0000000000000000000000000000000000000000000000000000000000000002",
							Value:           "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988",
						},
						{
							Address: "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988",
							Type:    int(merkletree.LeafTypeBalance),
							Value:   "100000000000000000000000",
						},
						{
							Address: "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988",
							Type:    int(merkletree.LeafTypeNonce),
							Value:   "2",
						},
						{
							Address:  "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988",
							Type:     int(merkletree.LeafTypeCode),
							Bytecode: "0xbeef2",
						},
						{
							Address:         "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988",
							Type:            int(merkletree.LeafTypeStorage),
							StoragePosition: "0x0000000000000000000000000000000000000000000000000000000000000000",
							Value:           "0xc949254d682d8c9ad5682521675b8f43b102aec4",
						},
						{
							Address: "0x61ba0248b0986c2480181c6e76b6adeeaa962483",
							Type:    int(merkletree.LeafTypeNonce),
							Value:   "1",
						},
						{
							Address:  "0x61ba0248b0986c2480181c6e76b6adeeaa962483",
							Type:     int(merkletree.LeafTypeCode),
							Bytecode: "0xbeef3",
						},
						{
							Address:         "0x61ba0248b0986c2480181c6e76b6adeeaa962483",
							Type:            int(merkletree.LeafTypeStorage),
							StoragePosition: "0x0000000000000000000000000000000000000000000000000000000000000000",
							Value:           "0x01",
						},
					},
				},
			},
		},
		{
			description: "imported from network-config.example.json",
			inputConfigStr: `{
  "genesis": [
    {
      "address": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
      "balance": "1000000000000000000000"
    },
    {
      "address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
      "balance": "2000000000000000000000"
    },
    {
      "address": "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
      "balance": "3000000000000000000000"
    }
  ]
}`,
			expectedConfig: NetworkConfig{
				Genesis: state.Genesis{
					Actions: []*state.GenesisAction{
						{
							Address: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
							Type:    int(merkletree.LeafTypeBalance),
							Value:   "1000000000000000000000",
						},
						{
							Address: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
							Type:    int(merkletree.LeafTypeBalance),
							Value:   "2000000000000000000000",
						},
						{
							Address: "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
							Type:    int(merkletree.LeafTypeBalance),
							Value:   "3000000000000000000000",
						},
					},
				},
			},
		},
		{
			description:      "not valid JSON gives error",
			inputConfigStr:   "not a valid json",
			expectedErrorMsg: "failed to load genesis configuration from file. Error: invalid character 'o' in literal null (expecting 'u')",
		},
		{
			description:      "empty JSON gives error",
			expectedErrorMsg: "failed to load genesis configuration from file. Error: unexpected end of JSON input",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			file, err := os.CreateTemp("", "genesisConfig")
			require.NoError(t, err)
			defer func() {
				require.NoError(t, os.Remove(file.Name()))
			}()
			require.NoError(t, os.WriteFile(file.Name(), []byte(tc.inputConfigStr), 0600))

			flagSet := flag.NewFlagSet("test", flag.ExitOnError)
			flagSet.String(FlagNetwork, string(custom), "")
			flagSet.String(FlagCustomNetwork, file.Name(), "")
			ctx := cli.NewContext(nil, flagSet, nil)

			c := &Config{}
			if tc.expectedErrorMsg != "" {
				panicFunc := func() {
					c.loadNetworkConfig(ctx)
				}
				require.PanicsWithError(t, tc.expectedErrorMsg, panicFunc)
				// require.Panics(t, panicFunc)
			} else {
				c.loadNetworkConfig(ctx)
				require.Equal(t, tc.expectedConfig, c.NetworkConfig)
			}
		})
	}
}
