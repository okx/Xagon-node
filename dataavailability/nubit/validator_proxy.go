package nubit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: Replace the proxy with our validator sdk.

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ProveSharesParams represents the parameters for the ProveShares method
type ProveSharesParams struct {
	Height     int64  `json:"height"`
	StartShare uint64 `json:"startShare"`
	EndShare   uint64 `json:"endShare"`
}

// makeJSONRPCRequest sends a JSON-RPC request to the server
func makeJSONRPCRequest(url string, method string, params interface{}, id int) (*JSONRPCResponse, error) {
	// Create the JSON-RPC request payload
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      id,
	}

	// Encode the request into JSON
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send the request to the server
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Decode the JSON-RPC response
	var jsonResponse JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if jsonResponse.Error != nil {
		return &jsonResponse, fmt.Errorf("JSON-RPC error: %s", jsonResponse.Error.Message)
	}

	return &jsonResponse, nil
}

// Validator is the interface for the ProveShares method
type Validator interface {
	ProveShares(height int64, startShare, endShare uint64) (*ShareProof, error)
}

// JSONRPCValidator implements the Validator interface using JSON-RPC
type JSONRPCValidator struct {
	serverURL string
}

// NewJSONRPCValidator creates a new instance of JSONRPCValidator
func NewJSONRPCValidator(serverURL string) *JSONRPCValidator {
	return &JSONRPCValidator{serverURL: serverURL}
}

func (v *JSONRPCValidator) ProveShares(height int64, startShare, endShare uint64) (*ShareProof, error) {
	// Use makeJSONRPCRequest to send the request
	params := ProveSharesParams{
		Height:     height,
		StartShare: startShare,
		EndShare:   endShare,
	}

	response, err := makeJSONRPCRequest(v.serverURL, "ProveShares", params, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to make JSON-RPC request: %w", err)
	}

	// Decode the result into the expected type
	var result ShareProof
	if err := json.Unmarshal(response.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	return &result, nil
}
