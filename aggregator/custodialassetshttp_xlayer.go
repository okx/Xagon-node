package aggregator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type signRequest struct {
	UserID         int            `json:"userId"`
	OperateType    int            `json:"operateType"` // 1 seq 2 agg
	OperateAddress common.Address `json:"operateAddress"`
	Symbol         int            `json:"symbol"`        // devnet 2882 mainnet 2
	ProjectSymbol  int            `json:"projectSymbol"` // 3011
	RefOrderID     string         `json:"refOrderId"`
	OperateSymbol  int            `json:"operateSymbol"`  // 2
	OperateAmount  int            `json:"operateAmount"`  // 0
	SysFrom        int            `json:"sysFrom"`        // 3
	OtherInfo      string         `json:"otherInfo"`      //
	DepositAddress string         `json:"depositAddress"` // "" for auth
	ToAddress      string         `json:"toAddress"`      // "" for auth
	BatchID        int            `json:"batchId"`        // 0
}

type signResponse struct {
	Code           int    `json:"code"`
	Data           string `json:"data"`
	DetailMessages string `json:"detailMsg"`
	Msg            string `json:"msg"`
	Status         int    `json:"status"`
	Success        bool   `json:"success"`
}

type signResultRequest struct {
	UserID        int    `json:"userId"`
	OrderID       string `json:"orderId"`
	ProjectSymbol int    `json:"projectSymbol"`
}

func (a *Aggregator) newSignRequest(operateType int, operateAddress common.Address, otherInfo string) *signRequest {
	refOrderID := uuid.New().String()
	return &signRequest{
		UserID:         a.cfg.CustodialAssets.UserID,
		OperateType:    operateType,
		OperateAddress: operateAddress,
		Symbol:         a.cfg.CustodialAssets.Symbol,
		ProjectSymbol:  a.cfg.CustodialAssets.ProjectSymbol,
		RefOrderID:     refOrderID,
		OperateSymbol:  a.cfg.CustodialAssets.OperateSymbol,
		OperateAmount:  a.cfg.CustodialAssets.OperateAmount,
		SysFrom:        a.cfg.CustodialAssets.SysFrom,
		OtherInfo:      otherInfo,
	}
}

func (a *Aggregator) newSignResultRequest(orderID string) *signResultRequest {
	return &signResultRequest{
		UserID:        a.cfg.CustodialAssets.UserID,
		OrderID:       orderID,
		ProjectSymbol: a.cfg.CustodialAssets.ProjectSymbol,
	}
}

func (a *Aggregator) postCustodialAssets(ctx context.Context, request *signRequest) error {
	if a == nil || !a.cfg.CustodialAssets.Enable {
		return errCustodialAssetsNotEnabled
	}
	mLog := log.WithFields(getTraceID(ctx))

	payload, err := a.sortSignRequestMarshal(request)
	if err != nil {
		return fmt.Errorf("error marshal request: %w", err)
	}

	reqSignURL, err := url.JoinPath(a.cfg.CustodialAssets.URL, a.cfg.CustodialAssets.RequestSignURI)
	if err != nil {
		return fmt.Errorf("error join url: %w", err)
	}

	req, err := http.NewRequest("POST", reqSignURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	mLog.Infof("post custodial assets request: %v", string(payload))

	req.Header.Set("Content-Type", "application/json")
	err = a.auth(ctx, req)
	if err != nil {
		return fmt.Errorf("error auth: %w", err)
	}
	mLog.Infof("post custodial assets request: %+v", req)

	client := &http.Client{
		Timeout: httpTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	mLog.Infof("post custodial assets response: %v", string(body))
	var signResp signResponse
	err = json.Unmarshal(body, &signResp)
	if err != nil {
		return fmt.Errorf("error unmarshal %v response body: %w", resp, err)
	}
	if signResp.Status != 200 || !signResp.Success {
		return fmt.Errorf("error response %v status: %v", signResp, signResp.Status)
	}

	return nil
}

func (a *Aggregator) sortSignRequestMarshal(request *signRequest) ([]byte, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshal request: %w", err)
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(payload, &jsonMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal body error: %v", err)
	}
	sortedKeys := make([]string, 0, len(jsonMap))
	for k := range jsonMap {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	// Creating a new map with sorted keys
	sortedJSON := make(map[string]interface{})
	for _, key := range sortedKeys {
		sortedJSON[key] = jsonMap[key]
	}

	// Converting the sorted map to JSON
	return json.Marshal(sortedJSON)
}

func (a *Aggregator) querySignResult(ctx context.Context, request *signResultRequest) (*signResponse, error) {
	if a == nil || !a.cfg.CustodialAssets.Enable {
		return nil, errCustodialAssetsNotEnabled
	}
	querySignURL, err := url.JoinPath(a.cfg.CustodialAssets.URL, a.cfg.CustodialAssets.QuerySignURI)
	if err != nil {
		return nil, fmt.Errorf("error join url: %w", err)
	}
	params := url.Values{}
	params.Add("orderId", request.OrderID)
	params.Add("projectSymbol", fmt.Sprintf("%d", request.ProjectSymbol))
	fullQuerySignURL := fmt.Sprintf("%s?%s", querySignURL, params.Encode())

	req, err := http.NewRequest("GET", fullQuerySignURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	err = a.auth(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error auth: %w", err)
	}
	mLog := log.WithFields(getTraceID(ctx))
	mLog.Infof("get sign result request: %+v", req)

	client := &http.Client{
		Timeout: httpTimeout,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	mLog.Infof("query sign result response: %v", string(body))
	var signResp signResponse
	err = json.Unmarshal(body, &signResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal %v response body: %w", response, err)
	}

	if signResp.Status != 200 || len(signResp.Data) == 0 {
		return nil, fmt.Errorf("error response %v status: %v", signResp, signResp.Status)
	}

	return &signResp, nil
}

func (a *Aggregator) waitResult(parentCtx context.Context, request *signResultRequest) (*signResponse, error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()
	ctx, cancel := context.WithTimeout(parentCtx, a.cfg.CustodialAssets.WaitResultTimeout.Duration)
	defer cancel()

	mLog := log.WithFields(getTraceID(ctx))
	for {
		result, err := a.querySignResult(ctx, request)
		if err == nil {
			mLog.Infof("query sign result success: %v", result)
			return result, nil
		}
		mLog.Infof("query sign result failed: %v", err)

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

func (a *Aggregator) postSignRequestAndWaitResult(ctx context.Context, request *signRequest) ([]byte, error) {
	if a == nil || !a.cfg.CustodialAssets.Enable {
		return nil, errCustodialAssetsNotEnabled
	}
	mLog := log.WithFields(getTraceID(ctx))
	if err := a.postCustodialAssets(ctx, request); err != nil {
		return nil, fmt.Errorf("error post custodial assets: %w", err)
	}
	mLog.Infof("post custodial assets success")
	result, err := a.waitResult(ctx, a.newSignResultRequest(request.RefOrderID))
	if err != nil {
		return nil, fmt.Errorf("error wait result: %w", err)
	}
	mLog.Infof("wait result success: %v", result)

	return hex.DecodeHex(result.Data)
}

func (a *Aggregator) postApproveAndWaitResult(ctx context.Context, request *signRequest) ([]byte, error) {
	if a == nil || !a.cfg.CustodialAssets.Enable {
		return nil, errCustodialAssetsNotEnabled
	}
	mLog := log.WithFields(getTraceID(ctx))
	if err := a.postCustodialAssets(ctx, request); err != nil {
		return nil, fmt.Errorf("error post custodial assets: %w", err)
	}
	mLog.Infof("post custodial assets success")
	result, err := a.waitResult(ctx, a.newSignResultRequest(request.RefOrderID))
	if err != nil {
		return nil, fmt.Errorf("error wait result: %w", err)
	}
	mLog.Infof("wait result success: %v", result)

	return hex.DecodeHex(result.Data)
}
