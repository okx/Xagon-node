package config

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-data-streamer/log"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

const (
	// FlagCfg is the flag for cfg
	FlagCfg = "cfg"
)

// L1 is the default configuration values
type L1 struct {
	ChainId       uint64                   `mapstructure:"ChainId"`
	RPC           string                   `mapstructure:"RPC"`
	SeqPrivateKey types.KeystoreFileConfig `mapstructure:"SeqPrivateKey"`
	AggPrivateKey types.KeystoreFileConfig `mapstructure:"AggPrivateKey"`

	PolygonMaticAddress       common.Address `mapstructure:"PolygonMaticAddress"`
	GlobalExitRootManagerAddr common.Address `mapstructure:"GlobalExitRootManagerAddress"`
	DataCommitteeAddr         common.Address `mapstructure:"DataCommitteeAddress"`
	PolygonZkEVMAddress       common.Address `mapstructure:"PolygonZkEVMAddress"`
}

// Config is the configuration for the tool
type Config struct {
	Port int        `mapstructure:"Port"`
	L1   L1         `mapstructure:"L1"`
	Log  log.Config `mapstructure:"Log"`
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
