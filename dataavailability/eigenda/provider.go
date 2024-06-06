package eigenda

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/dataavailability"
	"github.com/0xPolygonHermez/zkevm-node/log"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// ErrDisperseFailed is used when blob dispersion failed
	ErrDisperseFailed = errors.New("disperse blob on EigenDA layer failed")
	// ErrInsufficientSignatures is there are insufficient signatures
	ErrInsufficientSignatures = errors.New("insufficient signatures, confirmation threshold for blob not met")
)

// DataAvailabilityProvider is the EigenDA backend manager that holds the DA implementation.
// It contains the implementation of SequenceSender and SequenceRetriever of the zkevm's
// dataavailability package.
type DataAvailabilityProvider struct {
	cfg    dataavailability.Config
	client *DisperserClient
}

// NewDataAvailabilityProvider creates a new data availability provider
func NewDataAvailabilityProvider(cfg dataavailability.Config) *DataAvailabilityProvider {
	// TODO: Switch to authenticated blob dispersion pipeline
	signer := MockBlobRequestSigner{}
	p := &DataAvailabilityProvider{
		cfg:    cfg,
		client: NewDisperserClient(&cfg, signer),
	}

	return p
}

// Init initializes the data availability provider
func (d *DataAvailabilityProvider) Init() error {
	return nil
}

// PostSequence encodes and posts the sequence to the EigenDA layer and verifies that
// the blob gets confirmed or finalized. If the sequence is successfully posted, the
// pipeline encodes the blob data and returns the data availability message
func (d *DataAvailabilityProvider) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	// Send blob to EigenDA disperser
	blobData := EncodeSequence(batchesData)
	_, idBytes, err := d.client.DisperseBlob(ctx, blobData, []uint8{})
	if err != nil {
		log.Error("failed to send blob to EigenDA disperser: ", err)
		return nil, err
	}
	log.Debug("sent blob to EigenDA disperser, request id: ", base64.StdEncoding.EncodeToString(idBytes))

	startTime := time.Now()
	var blobStatusReply *disperser_rpc.BlobStatusReply
	for {
		blobStatusReply, err = d.client.GetBlobStatus(ctx, idBytes)
		if err != nil {
			log.Error("error getting blob status: %v\n", err)
			return nil, err
		}

		// Get blob status
		currStatus := blobStatusReply.GetStatus()
		if currStatus == disperser_rpc.BlobStatus_CONFIRMED || currStatus == disperser_rpc.BlobStatus_FINALIZED {
			break
		} else if currStatus == disperser_rpc.BlobStatus_FAILED {
			return nil, ErrDisperseFailed
		} else if currStatus == disperser_rpc.BlobStatus_INSUFFICIENT_SIGNATURES {
			return nil, ErrInsufficientSignatures
		} else {
			// status == BlobStatus_PROCESSING || BlobStatus_DISPERSING || BlobStatus_UNKNOWN
			if time.Since(startTime) > d.cfg.BlobStatusConfirmedTimeout.Duration {
				err = fmt.Errorf("blob status confirmation timeout")
				log.Error("Error: ", err)
				return nil, err
			}
		}
		time.Sleep(d.cfg.RetrieveBlobStatusPeriod.Duration)
	}

	if blobStatusReply == nil {
		return nil, fmt.Errorf("empty blob status reply returned")
	}

	// Get abi-encoded data availability message
	info := blobStatusReply.GetInfo()
	data, err := GetBlobData(info)
	if err != nil {
		log.Error("Error getting blob data: ", err)
		return nil, err
	}

	return TryEncodeToDataAvailabilityMessage(data)
}

// GetSequence gets blob data from the EigenDA layer and decodes the blob data into
// L2 batches data. The EigenDA provider does not use batchHashes to retrieve the L2
// batches data.
func (d *DataAvailabilityProvider) GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	var batchesData [][]byte
	blobData, err := TryDecodeFromDataAvailabilityMessage(dataAvailabilityMessage)
	if err != nil {
		log.Error("Error decoding from da message: ", err)
		return nil, err
	}
	batchHeaderHash, err := blobData.BlobVerificationProof.BatchMetadata.GetBatchHeaderHash()
	if err != nil {
		return nil, err
	}

	reply, err := d.client.RetrieveBlob(ctx, batchHeaderHash, blobData.BlobVerificationProof.BlobIndex)
	if err != nil {
		log.Error("Error retrieving blob from EigenDA disperser: ", err)
		return nil, err
	}
	data, _ := DecodeSequence(reply.GetData())
	batchesData = append(batchesData, data...)

	return batchesData, nil
}

// GetDataAvailabilityMessageFromId gets the data availability message from request ID
func (d *DataAvailabilityProvider) GetDataAvailabilityMessageFromId(ctx context.Context, requestId []byte) ([]byte, error) {
	blobStatusReply, err := d.client.GetBlobStatus(ctx, requestId)
	if err != nil {
		log.Error("Error getting blob status from EigenDA disperser: ", err)
		return nil, err
	}

	// Get blob status
	status := blobStatusReply.GetStatus()
	confirmedFlag := status == disperser_rpc.BlobStatus_CONFIRMED || status == disperser_rpc.BlobStatus_FINALIZED
	if !confirmedFlag {
		return nil, fmt.Errorf("EigenDA blob not confirmed, unable to retrieve blob information")
	}

	info := blobStatusReply.GetInfo()
	blobData, err := GetBlobData(info)
	if err != nil {
		return nil, err
	}

	return TryEncodeToDataAvailabilityMessage(blobData)
}
