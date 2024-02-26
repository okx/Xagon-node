package jsonrpc

import "github.com/0xPolygonHermez/zkevm-node/jsonrpc/nacos"

func (s *Server) registerNacos() {
	// start nacos client for registering restful service
	if s.config.Nacos.URLs != "" {
		nacos.StartNacosClient(s.config.Nacos.URLs, s.config.Nacos.NamespaceId, s.config.Nacos.ApplicationName, s.config.Nacos.ExternalListenAddr)
	}

	// start nacos client for registering restful service
	if s.config.NacosWs.URLs != "" {
		nacos.StartNacosClient(s.config.NacosWs.URLs, s.config.NacosWs.NamespaceId, s.config.NacosWs.ApplicationName, s.config.NacosWs.ExternalListenAddr)
	}
}

func (s *Server) getBatchReqLimit() (bool, uint) {
	var batchRequestEnable bool
	var batchRequestLimit uint
	// if apollo is enabled, get the config from apollo
	if getApolloConfig().Enable() {
		getApolloConfig().RLock()
		batchRequestEnable = getApolloConfig().BatchRequestsEnabled
		batchRequestLimit = getApolloConfig().BatchRequestsLimit
		getApolloConfig().RUnlock()
	} else {
		batchRequestEnable = s.config.BatchRequestsEnabled
		batchRequestLimit = s.config.BatchRequestsLimit
	}

	return batchRequestEnable, batchRequestLimit
}
