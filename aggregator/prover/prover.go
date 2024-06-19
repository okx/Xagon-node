package prover

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/metrics"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

var (
	ErrBadProverResponse    = errors.New("Prover returned wrong type for response")  //nolint:revive
	ErrProverInternalError  = errors.New("Prover returned INTERNAL_ERROR response")  //nolint:revive
	ErrProverCompletedError = errors.New("Prover returned COMPLETED_ERROR response") //nolint:revive
	ErrBadRequest           = errors.New("Prover returned ERROR for a bad request")  //nolint:revive
	ErrUnspecified          = errors.New("Prover returned an UNSPECIFIED response")  //nolint:revive
	ErrUnknown              = errors.New("Prover returned an unknown response")      //nolint:revive
	ErrProofCanceled        = errors.New("Proof has been canceled")                  //nolint:revive
)

// Prover abstraction of the grpc prover client.
type Prover struct {
	sync.Mutex

	name                      string
	id                        string
	address                   net.Addr
	proofStatePollingInterval types.Duration
	stream                    AggregatorService_ChannelServer
}

// New returns a new Prover instance.
func New(stream AggregatorService_ChannelServer, addr net.Addr, proofStatePollingInterval types.Duration) (*Prover, error) {
	p := &Prover{
		stream:                    stream,
		address:                   addr,
		proofStatePollingInterval: proofStatePollingInterval,
	}
	status, err := p.Status()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve prover id %w", err)
	}
	p.name = status.ProverName
	p.id = status.ProverId
	return p, nil
}

// Name returns the Prover name.
func (p *Prover) Name() string { return p.name }

// ID returns the Prover ID.
func (p *Prover) ID() string { return p.id }

// Addr returns the prover IP address.
func (p *Prover) Addr() string {
	if p.address == nil {
		return ""
	}
	return p.address.String()
}

// Status gets the prover status.
func (p *Prover) Status() (*GetStatusResponse, error) {
	req := &AggregatorMessage{
		Request: &AggregatorMessage_GetStatusRequest{
			GetStatusRequest: &GetStatusRequest{},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return nil, err
	}
	if msg, ok := res.Response.(*ProverMessage_GetStatusResponse); ok {
		return msg.GetStatusResponse, nil
	}
	return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &ProverMessage_GetStatusResponse{}, res.Response)
}

// IsIdle returns true if the prover is idling.
func (p *Prover) IsIdle() (bool, error) {
	status, err := p.Status()
	if err != nil {
		return false, err
	}
	return status.Status == GetStatusResponse_STATUS_IDLE, nil
}

// SupportsForkID returns true if the prover supports the given fork id.
func (p *Prover) SupportsForkID(forkID uint64) bool {
	status, err := p.Status()
	if err != nil {
		log.Warnf("Error asking status for prover ID %s: %v", p.ID(), err)
		return false
	}

	log.Debugf("Prover %s supports fork ID %d", p.ID(), status.ForkId)

	return status.ForkId == forkID
}

// BatchProof instructs the prover to generate a batch proof for the provided
// input. It returns the ID of the proof being computed.
func (p *Prover) BatchProof(input *InputProver) (*string, error) {
	metrics.WorkingProver()

	req := &AggregatorMessage{
		Request: &AggregatorMessage_GenBatchProofRequest{
			GenBatchProofRequest: &GenBatchProofRequest{Input: input},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return nil, err
	}

	if msg, ok := res.Response.(*ProverMessage_GenBatchProofResponse); ok {
		switch msg.GenBatchProofResponse.Result {
		case Result_RESULT_UNSPECIFIED:
			return nil, fmt.Errorf("failed to generate proof %s, %w, input %v", msg.GenBatchProofResponse.String(), ErrUnspecified, input)
		case Result_RESULT_OK:
			return &msg.GenBatchProofResponse.Id, nil
		case Result_RESULT_ERROR:
			return nil, fmt.Errorf("failed to generate proof %s, %w, input %v", msg.GenBatchProofResponse.String(), ErrBadRequest, input)
		case Result_RESULT_INTERNAL_ERROR:
			return nil, fmt.Errorf("failed to generate proof %s, %w, input %v", msg.GenBatchProofResponse.String(), ErrProverInternalError, input)
		default:
			return nil, fmt.Errorf("failed to generate proof %s, %w,input %v", msg.GenBatchProofResponse.String(), ErrUnknown, input)
		}
	}

	return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &ProverMessage_GenBatchProofResponse{}, res.Response)
}

// AggregatedProof instructs the prover to generate an aggregated proof from
// the two inputs provided. It returns the ID of the proof being computed.
func (p *Prover) AggregatedProof(inputProof1, inputProof2 string) (*string, error) {
	metrics.WorkingProver()

	req := &AggregatorMessage{
		Request: &AggregatorMessage_GenAggregatedProofRequest{
			GenAggregatedProofRequest: &GenAggregatedProofRequest{
				RecursiveProof_1: inputProof1,
				RecursiveProof_2: inputProof2,
			},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return nil, err
	}

	if msg, ok := res.Response.(*ProverMessage_GenAggregatedProofResponse); ok {
		switch msg.GenAggregatedProofResponse.Result {
		case Result_RESULT_UNSPECIFIED:
			return nil, fmt.Errorf("failed to aggregate proofs %s, %w, input 1 %s, input 2 %s",
				msg.GenAggregatedProofResponse.String(), ErrUnspecified, inputProof1, inputProof2)
		case Result_RESULT_OK:
			return &msg.GenAggregatedProofResponse.Id, nil
		case Result_RESULT_ERROR:
			return nil, fmt.Errorf("failed to aggregate proofs %s, %w, input 1 %s, input 2 %s",
				msg.GenAggregatedProofResponse.String(), ErrBadRequest, inputProof1, inputProof2)
		case Result_RESULT_INTERNAL_ERROR:
			return nil, fmt.Errorf("failed to aggregate proofs %s, %w, input 1 %s, input 2 %s",
				msg.GenAggregatedProofResponse.String(), ErrProverInternalError, inputProof1, inputProof2)
		default:
			return nil, fmt.Errorf("failed to aggregate proofs %s, %w, input 1 %s, input 2 %s",
				msg.GenAggregatedProofResponse.String(), ErrUnknown, inputProof1, inputProof2)
		}
	}

	return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &ProverMessage_GenAggregatedProofResponse{}, res.Response)
}

// FinalProof instructs the prover to generate a final proof for the given
// input. It returns the ID of the proof being computed.
func (p *Prover) FinalProof(inputProof string, aggregatorAddr string) (*string, error) {
	metrics.WorkingProver()

	req := &AggregatorMessage{
		Request: &AggregatorMessage_GenFinalProofRequest{
			GenFinalProofRequest: &GenFinalProofRequest{
				RecursiveProof: inputProof,
				AggregatorAddr: aggregatorAddr,
			},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return nil, err
	}

	if msg, ok := res.Response.(*ProverMessage_GenFinalProofResponse); ok {
		switch msg.GenFinalProofResponse.Result {
		case Result_RESULT_UNSPECIFIED:
			return nil, fmt.Errorf("failed to generate final proof %s, %w, input %s",
				msg.GenFinalProofResponse.String(), ErrUnspecified, inputProof)
		case Result_RESULT_OK:
			return &msg.GenFinalProofResponse.Id, nil
		case Result_RESULT_ERROR:
			return nil, fmt.Errorf("failed to generate final proof %s, %w, input %s",
				msg.GenFinalProofResponse.String(), ErrBadRequest, inputProof)
		case Result_RESULT_INTERNAL_ERROR:
			return nil, fmt.Errorf("failed to generate final proof %s, %w, input %s",
				msg.GenFinalProofResponse.String(), ErrProverInternalError, inputProof)
		default:
			return nil, fmt.Errorf("failed to generate final proof %s, %w, input %s",
				msg.GenFinalProofResponse.String(), ErrUnknown, inputProof)
		}
	}
	return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &ProverMessage_GenFinalProofResponse{}, res.Response)
}

// CancelProofRequest asks the prover to stop the generation of the proof
// matching the provided proofID.
func (p *Prover) CancelProofRequest(proofID string) error {
	req := &AggregatorMessage{
		Request: &AggregatorMessage_CancelRequest{
			CancelRequest: &CancelRequest{Id: proofID},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return err
	}
	if msg, ok := res.Response.(*ProverMessage_CancelResponse); ok {
		switch msg.CancelResponse.Result {
		case Result_RESULT_UNSPECIFIED:
			return fmt.Errorf("failed to cancel proof id [%s], %w, %s",
				proofID, ErrUnspecified, msg.CancelResponse.String())
		case Result_RESULT_OK:
			return nil
		case Result_RESULT_ERROR:
			return fmt.Errorf("failed to cancel proof id [%s], %w, %s",
				proofID, ErrBadRequest, msg.CancelResponse.String())
		case Result_RESULT_INTERNAL_ERROR:
			return fmt.Errorf("failed to cancel proof id [%s], %w, %s",
				proofID, ErrProverInternalError, msg.CancelResponse.String())
		default:
			return fmt.Errorf("failed to cancel proof id [%s], %w, %s",
				proofID, ErrUnknown, msg.CancelResponse.String())
		}
	}
	return fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &ProverMessage_CancelResponse{}, res.Response)
}

// WaitRecursiveProof waits for a recursive proof to be generated by the prover
// and returns it.
func (p *Prover) WaitRecursiveProof(ctx context.Context, proofID string) (string, error) {
	res, err := p.waitProof(ctx, proofID)
	if err != nil {
		return "", err
	}
	resProof := res.Proof.(*GetProofResponse_RecursiveProof)
	return resProof.RecursiveProof, nil
}

// WaitFinalProof waits for the final proof to be generated by the prover and
// returns it.
func (p *Prover) WaitFinalProof(ctx context.Context, proofID string) (*FinalProof, error) {
	res, err := p.waitProof(ctx, proofID)
	if err != nil {
		return nil, err
	}
	resProof := res.Proof.(*GetProofResponse_FinalProof)
	return resProof.FinalProof, nil
}

// waitProof waits for a proof to be generated by the prover and returns the
// prover response.
func (p *Prover) waitProof(ctx context.Context, proofID string) (*GetProofResponse, error) {
	defer metrics.IdlingProver()

	req := &AggregatorMessage{
		Request: &AggregatorMessage_GetProofRequest{
			GetProofRequest: &GetProofRequest{
				// TODO(pg): set Timeout field?
				Id: proofID,
			},
		},
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			res, err := p.call(req)
			if err != nil {
				return nil, err
			}
			if msg, ok := res.Response.(*ProverMessage_GetProofResponse); ok {
				switch msg.GetProofResponse.Result {
				case GetProofResponse_RESULT_PENDING:
					time.Sleep(p.proofStatePollingInterval.Duration)
					continue
				case GetProofResponse_RESULT_UNSPECIFIED:
					return nil, fmt.Errorf("failed to get proof ID: %s, %w, prover response: %s",
						proofID, ErrUnspecified, msg.GetProofResponse.String())
				case GetProofResponse_RESULT_COMPLETED_OK:
					return msg.GetProofResponse, nil
				case GetProofResponse_RESULT_ERROR:
					return nil, fmt.Errorf("failed to get proof with ID %s, %w, prover response: %s",
						proofID, ErrBadRequest, msg.GetProofResponse.String())
				case GetProofResponse_RESULT_COMPLETED_ERROR:
					return nil, fmt.Errorf("failed to get proof with ID %s, %w, prover response: %s",
						proofID, ErrProverCompletedError, msg.GetProofResponse.String())
				case GetProofResponse_RESULT_INTERNAL_ERROR:
					return nil, fmt.Errorf("failed to get proof ID: %s, %w, prover response: %s",
						proofID, ErrProverInternalError, msg.GetProofResponse.String())
				case GetProofResponse_RESULT_CANCEL:
					return nil, fmt.Errorf("proof generation was cancelled for proof ID %s, %w, prover response: %s",
						proofID, ErrProofCanceled, msg.GetProofResponse.String())
				default:
					return nil, fmt.Errorf("failed to get proof ID: %s, %w, prover response: %s",
						proofID, ErrUnknown, msg.GetProofResponse.String())
				}
			}
			return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &ProverMessage_GetProofResponse{}, res.Response)
		}
	}
}

// call sends a message to the prover and waits to receive the response over
// the connection stream.
func (p *Prover) call(req *AggregatorMessage) (*ProverMessage, error) {
	p.Lock()
	defer p.Unlock()
	if err := p.stream.Send(req); err != nil {
		return nil, err
	}
	res, err := p.stream.Recv()
	if err != nil {
		return nil, err
	}
	return res, nil
}
