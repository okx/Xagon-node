package jsonrpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/metrics"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/nacos"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

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

func (s *Server) handleWsMessage(httpRequest *http.Request, wsConn *concurrentWsConn, data []byte) ([]byte, error) {
	if validateWsRequest(data) != nil {
		return types.NewResponse(types.Request{}, nil, types.NewRPCError(types.InvalidRequestErrorCode, "Invalid json request")).Bytes()
	}
	single, err := s.isSingleRequest(data)
	if err != nil {
		return types.NewResponse(types.Request{}, nil, types.NewRPCError(types.InvalidRequestErrorCode, err.Error())).Bytes()
	}
	if single {
		return s.handler.HandleWs(data, wsConn, httpRequest)
	}
	return s.handleWsBatch(httpRequest, wsConn, data)
}

func (s *Server) handleWsBatch(httpRequest *http.Request, wsConn *concurrentWsConn, data []byte) ([]byte, error) {
	requests, err := s.parseRequests(data)
	if err != nil {
		return types.NewResponse(types.Request{}, nil, types.NewRPCError(types.InvalidRequestErrorCode, err.Error())).Bytes()
	}

	batchRequestEnable, batchRequestLimit := s.getBatchReqLimit()

	if !batchRequestEnable {
		return types.NewResponse(types.Request{}, nil, types.NewRPCError(types.InvalidRequestErrorCode, "Batch requests are disabled")).Bytes()
	}

	// Checking if batch requests limit is exceeded
	if batchRequestLimit > 0 {
		if len(requests) > int(batchRequestLimit) {
			errMsg := fmt.Sprintf("Batch requests limit exceeded: %d", batchRequestLimit)
			return types.NewResponse(types.Request{}, nil, types.NewRPCError(types.InvalidRequestErrorCode, errMsg)).Bytes()
		}
	}

	responses := make([]types.Response, 0, len(requests))

	for _, request := range requests {
		if !methodRateLimitAllow(request.Method) {
			responses = append(responses, types.NewResponse(request, nil, types.NewRPCError(types.InvalidParamsErrorCode, "server is too busy")))
			continue
		}
		st := time.Now()
		metrics.RequestMethodCount(request.Method)
		req := handleRequest{Request: request, wsConn: wsConn, HttpRequest: httpRequest}
		response := s.handler.Handle(req)
		responses = append(responses, response)
		metrics.RequestMethodDuration(request.Method, st)
	}

	return json.Marshal(responses)
}

func validateWsRequest(data []byte) error {
	if len(data) > maxRequestContentLength {
		return fmt.Errorf("content length too large (%d>%d)", len(data), maxRequestContentLength)
	}
	var valid bool
	var req types.Request
	if err := json.Unmarshal(data, &req); err != nil {
		valid = true
	}

	if !valid {
		var reqs []types.Request
		if err := json.Unmarshal(data, &reqs); err != nil {
			valid = true
		}
	}

	if !valid {
		return fmt.Errorf("invalid request")
	}

	return nil
}
