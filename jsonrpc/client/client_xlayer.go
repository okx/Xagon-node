package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

// JSONRPCRelay executes a 2.0 JSON RPC HTTP Post Request to the provided URL with
// types.Request, which is compatible with the Ethereum
// JSON RPC Server.
func JSONRPCRelay(url string, request types.Request) (types.Response, error) {
	httpRes, err := sendJSONRPC_HTTPRequest(url, request)
	if err != nil {
		return types.Response{}, err
	}

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return types.Response{}, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		return types.Response{}, fmt.Errorf("http error: %v - %v", httpRes.StatusCode, string(resBody))
	}

	var res types.Response
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return types.Response{}, err
	}
	if res.Error != nil {
		return types.Response{}, fmt.Errorf("response error: %v - %v", res.Error.Code, res.Error.Message)
	}
	return res, nil
}
