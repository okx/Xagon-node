package config

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-data-streamer/log"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	// FlagCfg is the flag for cfg
	FlagCfg = "cfg"
	// FlagGenesis is the flag for genesis file
	FlagGenesis = "genesis"
)

// OnlineConfig is the configuration for the online data streamer
type OnlineConfig struct {
	URI        string                  `mapstructure:"URI"`
	StreamType datastreamer.StreamType `mapstructure:"StreamType"`
}

// MTConfig is the configuration for the merkle tree
type MTConfig struct {
	URI        string `mapstructure:"URI"`
	MaxThreads int    `mapstructure:"MaxThreads"`
	CacheFile  string `mapstructure:"CacheFile"`
}

// StreamServerCfg is the configuration for the offline data streamer
type StreamServerCfg struct {
	// Port to listen on
	Port uint16 `mapstructure:"Port"`
	// Filename of the binary data file
	Filename string `mapstructure:"Filename"`
	// Version of the binary data file
	Version uint8 `mapstructure:"Version"`
	// ChainID is the chain ID
	ChainID uint64 `mapstructure:"ChainID"`
	// Log is the log configuration
	Log log.Config `mapstructure:"Log"`
	// UpgradeEtrogBatchNumber is the batch number of the upgrade etrog
	UpgradeEtrogBatchNumber uint64 `mapstructure:"UpgradeEtrogBatchNumber"`
	// WriteTimeout is the TCP write timeout when sending data to a datastream client
	WriteTimeout types.Duration `mapstructure:"WriteTimeout"`
}

// Config is the configuration for the tool
type Config struct {
	Online     OnlineConfig    `mapstructure:"Online"`
	Offline    StreamServerCfg `mapstructure:"Offline"`
	StateDB    db.Config       `mapstructure:"StateDB"`
	Executor   executor.Config `mapstructure:"Executor"`
	MerkleTree MTConfig        `mapstructure:"MerkleTree"`
	Log        log.Config      `mapstructure:"Log"`
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

// Load parses the configuration values from the config file and environment variables
func Load(ctx *cli.Context) (*Config, error) {
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
	viper.SetEnvPrefix("ZKEVM_DATA_STREAMER")
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

	return cfg, nil
}
