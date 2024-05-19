package config

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/gasprice"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/sequencesender"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	// FlagYes is the flag for yes.
	FlagYes = "yes"
	// FlagCfg is the flag for cfg.
	FlagCfg = "cfg"
	// FlagNetwork is the flag for the network name. Valid values: ["testnet", "mainnet", "cardona", "custom"].
	FlagNetwork = "network"
	// FlagCustomNetwork is the flag for the custom network file. This is required if --network=custom
	FlagCustomNetwork = "custom-network-file"
	// FlagAmount is the flag for amount.
	FlagAmount = "amount"
	// FlagRemoteMT is the flag for remote-merkletree.
	FlagRemoteMT = "remote-merkletree"
	// FlagComponents is the flag for components.
	FlagComponents = "components"
	// FlagHTTPAPI is the flag for http.api.
	FlagHTTPAPI = "http.api"
	// FlagKeyStorePath is the path of the key store file containing the private key of the account going to sing and approve the tokens
	FlagKeyStorePath = "key-store-path"
	// FlagPassword is the password needed to decrypt the key store
	FlagPassword = "password"
	// FlagMigrations is the flag for migrations.
	FlagMigrations = "migrations"
	// FlagOutputFile is the flag for the output file
	FlagOutputFile = "output"
	// FlagMaxAmount is the flag to avoid to use the flag FlagAmount
	FlagMaxAmount = "max-amount"
	// FlagDocumentationFileType is the flag for the choose which file generate json-schema
	FlagDocumentationFileType = "config-file"
)

/*
Config represents the configuration of the entire Hermez Node
The file is [TOML format]
You could find some examples:
  - `config/environments/local/local.node.config.toml`: running a permisionless node
  - `config/environments/mainnet/node.config.toml`
  - `config/environments/public/node.config.toml`
  - `test/config/test.node.config.toml`: configuration for a trusted node used in CI

[TOML format]: https://en.wikipedia.org/wiki/TOML
*/
type Config struct {
	// This define is a trusted node (`true`) or a permission less (`false`). If you don't known
	// set to `false`
	IsTrustedSequencer bool `mapstructure:"IsTrustedSequencer"`
	// Last batch number before  a forkid change (fork upgrade). That implies that
	// greater batch numbers are going to be trusted but no virtualized neither verified.
	// So after the batch number `ForkUpgradeBatchNumber` is virtualized and verified you could update
	// the system (SC,...) to new forkId and remove this value to allow the system to keep
	// Virtualizing and verifying the new batchs.
	// Check issue [#2236](https://github.com/0xPolygonHermez/zkevm-node/issues/2236) to known more
	// This value overwrite `SequenceSender.ForkUpgradeBatchNumber`
	ForkUpgradeBatchNumber uint64 `mapstructure:"ForkUpgradeBatchNumber"`
	// Which is the new forkId
	ForkUpgradeNewForkId uint64 `mapstructure:"ForkUpgradeNewForkId"`
	// Configure Log level for all the services, allow also to store the logs in a file
	Log log.Config
	// Configuration of the etherman (client for access L1)
	Etherman etherman.Config
	// Configuration for ethereum transaction manager
	EthTxManager ethtxmanager.Config
	// Pool service configuration
	Pool pool.Config
	// Configuration for RPC service. THis one offers a extended Ethereum JSON-RPC API interface to interact with the node
	RPC jsonrpc.Config
	// Configuration of service `Syncrhonizer`. For this service is also really important the value of `IsTrustedSequencer`
	// because depending of this values is going to ask to a trusted node for trusted transactions or not
	Synchronizer synchronizer.Config
	// Configuration of the sequencer service
	Sequencer sequencer.Config
	// Configuration of the sequence sender service
	SequenceSender sequencesender.Config
	// Configuration of the aggregator service
	Aggregator aggregator.Config
	// Configuration of the genesis of the network. This is used to known the initial state of the network
	NetworkConfig NetworkConfig
	// Configuration of the gas price suggester service
	L2GasPriceSuggester gasprice.Config
	// Configuration of the executor service
	Executor executor.Config
	// Configuration of the merkle tree client service. Not use in the node, only for testing
	MTClient merkletree.Config
	// Configuration of the metrics service, basically is where is going to publish the metrics
	Metrics metrics.Config
	// Configuration of the event database connection
	EventLog event.Config
	// Configuration of the hash database connection
	HashDB db.Config
	// State service configuration
	State state.Config
	// Apollo configuration
	Apollo types.ApolloConfig
	// ForceBatchAddress Address of the L1 ForceBatch contract
	Fork9UpgradeBatch uint64 `mapstructure:"Fork9UpgradeBatch"`
}

// Default parses the default configuration values.
func Default() (*Config, error) {
	var cfg Config
	viper.SetConfigType("toml")

	err := viper.ReadConfig(bytes.NewBuffer([]byte(DefaultValues)))
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Load loads the configuration
func Load(ctx *cli.Context, loadNetworkConfig bool) (*Config, error) {
	cfg, err := Default()
	if err != nil {
		return nil, err
	}
	configFilePath := ctx.String(FlagCfg)
	if configFilePath != "" {
		dirName, fileName := filepath.Split(configFilePath)

		fileExtension := strings.TrimPrefix(filepath.Ext(fileName), ".")
		fileNameWithoutExtension := strings.TrimSuffix(fileName, "."+fileExtension)

		viper.AddConfigPath(dirName)
		viper.SetConfigName(fileNameWithoutExtension)
		viper.SetConfigType(fileExtension)
	}
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("ZKEVM_NODE")
	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			log.Infof("config file not found")
		} else {
			log.Infof("error reading config file: ", err)
			return nil, err
		}
	}

	decodeHooks := []viper.DecoderConfigOption{
		// this allows arrays to be decoded from env var separated by ",", example: MY_VAR="value1,value2,value3"
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(mapstructure.TextUnmarshallerHookFunc(), mapstructure.StringToSliceHookFunc(","))),
	}

	err = viper.Unmarshal(&cfg, decodeHooks...)
	if err != nil {
		return nil, err
	}

	if loadNetworkConfig {
		// Load genesis parameters
		cfg.loadNetworkConfig(ctx)
	}
	return cfg, nil
}
