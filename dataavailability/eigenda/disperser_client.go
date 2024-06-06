package eigenda

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/dataavailability"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// type DisperserClient interface {
// 	DisperseBlob(ctx context.Context, data []byte, customQuorums []uint8) (*disperser.BlobStatus, []byte, error)
// 	DisperseBlobAuthenticated(ctx context.Context, data []byte, customQuorums []uint8) (*disperser.BlobStatus, []byte, error)
// 	GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error)
// }

type DisperserClient struct {
	cfg    *dataavailability.Config
	signer core.BlobRequestSigner
}

func NewDisperserClient(cfg *dataavailability.Config, signer core.BlobRequestSigner) *DisperserClient {
	return &DisperserClient{
		cfg:    cfg,
		signer: signer,
	}
}

func (c *DisperserClient) getDialOptions() []grpc.DialOption {
	if c.cfg.UseSecureGrpcFlag {
		cfg := &tls.Config{MinVersion: tls.VersionTLS13}
		credential := credentials.NewTLS(cfg)
		return []grpc.DialOption{grpc.WithTransportCredentials(credential)}
	} else {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}
}

func (c *DisperserClient) DisperseBlob(ctx context.Context, data []byte, quorums []uint8) (*disperser.BlobStatus, []byte, error) {
	addr := fmt.Sprintf("%v:%v", c.cfg.Hostname, c.cfg.Port)

	dialOptions := c.getDialOptions()
	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = conn.Close() }()

	DisperserClient := disperser_rpc.NewDisperserClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, c.cfg.Timeout.Duration)
	defer cancel()

	quorumNumbers := make([]uint32, len(quorums))
	for i, q := range quorums {
		quorumNumbers[i] = uint32(q)
	}

	// check every 32 bytes of data are within the valid range for a bn254 field element
	_, err = rs.ToFrArray(data)
	if err != nil {
		return nil, nil, fmt.Errorf("encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617 %w", err)
	}
	request := &disperser_rpc.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: quorumNumbers,
	}

	reply, err := DisperserClient.DisperseBlob(ctxTimeout, request)
	if err != nil {
		return nil, nil, err
	}

	blobStatus, err := disperser.FromBlobStatusProto(reply.GetResult())
	if err != nil {
		return nil, nil, err
	}

	return blobStatus, reply.GetRequestId(), nil
}

func (c *DisperserClient) DisperseBlobAuthenticated(ctx context.Context, data []byte, quorums []uint8) (*disperser.BlobStatus, []byte, error) {

	addr := fmt.Sprintf("%v:%v", c.cfg.Hostname, c.cfg.Port)

	dialOptions := c.getDialOptions()
	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = conn.Close() }()

	DisperserClient := disperser_rpc.NewDisperserClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, c.cfg.Timeout.Duration)

	defer cancel()

	stream, err := DisperserClient.DisperseBlobAuthenticated(ctxTimeout)
	if err != nil {
		return nil, nil, fmt.Errorf("error while calling DisperseBlobAuthenticated: %w", err)
	}

	quorumNumbers := make([]uint32, len(quorums))
	for i, q := range quorums {
		quorumNumbers[i] = uint32(q)
	}

	// check every 32 bytes of data are within the valid range for a bn254 field element
	_, err = rs.ToFrArray(data)
	if err != nil {
		return nil, nil, fmt.Errorf("encountered an error to convert a 32-bytes into a valid field element, please use the correct format where every 32bytes(big-endian) is less than 21888242871839275222246405745257275088548364400416034343698204186575808495617, %w", err)
	}
	request := &disperser_rpc.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: quorumNumbers,
		AccountId:           c.signer.GetAccountID(),
	}

	// Send the initial request
	err = stream.Send(&disperser_rpc.AuthenticatedRequest{Payload: &disperser_rpc.AuthenticatedRequest_DisperseRequest{
		DisperseRequest: request,
	}})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Get the Challenge
	reply, err := stream.Recv()
	if err != nil {
		return nil, nil, fmt.Errorf("error while receiving: %w", err)
	}
	authHeaderReply, ok := reply.Payload.(*disperser_rpc.AuthenticatedReply_BlobAuthHeader)
	if !ok {
		return nil, nil, errors.New("expected challenge")
	}

	authHeader := core.BlobAuthHeader{
		BlobCommitments: encoding.BlobCommitments{},
		AccountID:       "",
		Nonce:           authHeaderReply.BlobAuthHeader.ChallengeParameter,
	}

	authData, err := c.signer.SignBlobRequest(authHeader)
	if err != nil {
		return nil, nil, errors.New("error signing blob request")
	}

	// Process challenge and send back challenge_reply
	err = stream.Send(&disperser_rpc.AuthenticatedRequest{Payload: &disperser_rpc.AuthenticatedRequest_AuthenticationData{
		AuthenticationData: &disperser_rpc.AuthenticationData{
			AuthenticationData: authData,
		},
	}})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send challenge reply: %w", err)
	}

	reply, err = stream.Recv()
	if err != nil {
		return nil, nil, fmt.Errorf("error while receiving final reply: %w", err)
	}
	disperseReply, ok := reply.Payload.(*disperser_rpc.AuthenticatedReply_DisperseReply) // Process the final disperse_reply
	if !ok {
		return nil, nil, errors.New("expected DisperseReply")
	}

	blobStatus, err := disperser.FromBlobStatusProto(disperseReply.DisperseReply.GetResult())
	if err != nil {
		return nil, nil, err
	}

	return blobStatus, disperseReply.DisperseReply.GetRequestId(), nil
}

func (c *DisperserClient) GetBlobStatus(ctx context.Context, requestID []byte) (*disperser_rpc.BlobStatusReply, error) {
	addr := fmt.Sprintf("%v:%v", c.cfg.Hostname, c.cfg.Port)
	dialOptions := c.getDialOptions()
	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		return nil, err
	}

	DisperserClient := disperser_rpc.NewDisperserClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	request := &disperser_rpc.BlobStatusRequest{
		RequestId: requestID,
	}

	reply, err := DisperserClient.GetBlobStatus(ctxTimeout, request)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *DisperserClient) RetrieveBlob(ctx context.Context, batchHeaderHash []byte, blobIndex uint32) (*disperser_rpc.RetrieveBlobReply, error) {
	addr := fmt.Sprintf("%v:%v", c.cfg.Hostname, c.cfg.Port)
	dialOptions := c.getDialOptions()
	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		return nil, err
	}

	DisperserClient := disperser_rpc.NewDisperserClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	request := &disperser_rpc.RetrieveBlobRequest{
		BatchHeaderHash: batchHeaderHash,
		BlobIndex:       blobIndex,
	}
	reply, err := DisperserClient.RetrieveBlob(ctxTimeout, request)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
