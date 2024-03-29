package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	agglayertx "github.com/0xPolygon/agglayer/tx"
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
	LastVirtualBatch string            `json:"lastVirtualBatch"`
	NewVirtualBatch  string            `json:"newVirtualBatch"`
	ZKP              httpAggregatorZKP `json:"ZKP"`
}

func (a *Aggregator) signTx(ctx context.Context, tx agglayertx.Tx) (*agglayertx.SignedTx, error) {
	lastVirtualBatch := tx.LastVerifiedBatch.Hex()
	newVirtualBatch := tx.NewVerifiedBatch.Hex()
	newStateRoot := tx.ZKP.NewStateRoot.Hex()
	newLocalExitRoot := tx.ZKP.NewLocalExitRoot.Hex()
	proof := tx.ZKP.Proof.Hex()
	httpPayload := httpAggregator{
		LastVirtualBatch: lastVirtualBatch,
		NewVirtualBatch:  newVirtualBatch,
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
