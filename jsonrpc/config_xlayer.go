package jsonrpc

// NacosConfig has parameters to config the nacos client
type NacosConfig struct {
	// URLs nacos server urls for discovery service of rest api, url is separated by ","
	URLs string `mapstructure:"URLs"`

	// NamespaceId nacos namepace id for discovery service of rest api
	NamespaceId string `mapstructure:"NamespaceId"`

	// ApplicationName rest application name in  nacos
	ApplicationName string `mapstructure:"ApplicationName"`

	// ExternalListenAddr Set the rest-server external ip and port, when it is launched by Docker
	ExternalListenAddr string `mapstructure:"ExternalListenAddr"`
}
