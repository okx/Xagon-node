package ethtxmanager

// RPCConfig is the configuration for the rpc server
type RPCConfig struct {
	// Enable is the flag to enable the rpc server
	Enable bool `mapstructure:"Enable"`

	// Host is the host of the rpc server
	Host string `mapstructure:"Host"`

	// RPCPort is the port of the rpc server
	Port int `mapstructure:"Port"`
}
