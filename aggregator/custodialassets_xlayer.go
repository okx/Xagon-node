package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	agglayertx "github.com/0xPolygon/agglayer/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type contextKey string

const (
	sigLen                 = 4
	hashLen                = 32
	proofLen               = 24
	traceID     contextKey = "traceID"
	httpTimeout            = 2 * time.Minute
)

func getTraceID(ctx context.Context) (string, string) {
	if ctx == nil || ctx.Value(traceID) == nil {
		return "", ""
	}
	return string(traceID), ctx.Value(traceID).(string)
}

var (
	errCustodialAssetsNotEnabled = errors.New("custodial assets not enabled")
)

type httpAggregatorZKP struct {
	NewStateRoot     string `json:"newStateRoot"`
	NewLocalExitRoot string `json:"newLocalExitRoot"`
	Proof            string `json:"proof"`
}
type httpAggregator struct {
	LastVerifiedBatch string            `json:"lastVerifiedBatch"`
	NewVerifiedBatch  string            `json:"newVerifiedBatch"`
	ZKP               httpAggregatorZKP `json:"ZKP"`
	ContractAddress   common.Address    `json:"contractAddress"`
}

func (a *Aggregator) signTx(ctx context.Context, tx agglayertx.Tx) (*agglayertx.SignedTx, error) {
	lastVerifiedBatch := tx.LastVerifiedBatch.Hex()
	newVerifiedBatch := tx.NewVerifiedBatch.Hex()
	newStateRoot := tx.ZKP.NewStateRoot.Hex()
	newLocalExitRoot := tx.ZKP.NewLocalExitRoot.Hex()
	proof := tx.ZKP.Proof.Hex()

	httpPayload := httpAggregator{
		LastVerifiedBatch: lastVerifiedBatch,
		NewVerifiedBatch:  newVerifiedBatch,
		ZKP: httpAggregatorZKP{
			NewStateRoot:     newStateRoot,
			NewLocalExitRoot: newLocalExitRoot,
			Proof:            proof,
		},
	}
	otherInfo, err := json.Marshal(httpPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal http payload: %v", err)
	}

	signature, err := a.postSignRequestAndWaitResult(ctx, a.newSignRequest(a.cfg.CustodialAssets.OperateTypeAgg, a.cfg.CustodialAssets.AggregatorAddr, string(otherInfo)))
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %v", err)
	}
	return &agglayertx.SignedTx{
		Tx:        tx,
		Signature: signature,
	}, nil
}

type httpApprove struct {
	ApproveToAddress string `json:"approveToAddress"`
	ApproveAmount    string `json:"approveAmount"`
	ContractAddress  string `json:"contractAddress"`
	GasLimit         uint64 `json:"gasLimit"`
	GasPrice         string `json:"gasPrice"`
	Nonce            uint64 `json:"nonce"`
}

func (a *Aggregator) approve(approveToAddress, approveAmount, contractAddress string, GasLimit, Nonce uint64, GasPrice string) ([]byte, error) {
	httpPayload := httpApprove{
		ApproveToAddress: approveToAddress,
		ApproveAmount:    approveAmount,
		ContractAddress:  contractAddress,
		GasLimit:         GasLimit,
		Nonce:            Nonce,
		GasPrice:         GasPrice,
	}
	otherInfo, err := json.Marshal(httpPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal http payload: %v", err)
	}

	ctx := context.WithValue(context.Background(), traceID, uuid.New().String())
	return a.postApproveAndWaitResult(ctx, a.newSignRequest(a.cfg.CustodialAssets.OperateTypeAgg, a.cfg.CustodialAssets.AggregatorAddr, string(otherInfo)))
}
