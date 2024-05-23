package jsonrpc

type apiRelayConfig struct {
	Enabled bool   `mapstructure:"Enabled"`
	DestURI string `mapstructure:"DestURI"`
}
