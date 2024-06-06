package eigenda

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/dataavailability"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrDisperseFailed         = errors.New("disperse blob on EigenDA layer failed")
	ErrInsufficientSignatures = errors.New("insufficient signatures, confirmation threshold for blob not met")
)

// DataAvailabilityProvider is the EigenDA backend manager that holds the DA implementation.
// It contains the implementation of SequenceSender and SequenceRetriever of the zkevm's
// dataavailability package.
type DataAvailabilityProvider struct {
	cfg    dataavailability.Config
	state  DAStorage
	client *DisperserClient
}

// Factory method for a new DataAvailibilityProvider instance
func NewDataAvailabilityProvider(cfg dataavailability.Config) *DataAvailabilityProvider {
	// Initialize in-memory DA storage
	s := DAStorage{
		inner: map[common.Hash][]byte{},
		mutex: &sync.RWMutex{},
	}
	signer := MockBlobRequestSigner{}

	p := &DataAvailabilityProvider{
		cfg:    cfg,
		state:  s,
		client: NewDisperserClient(&cfg, signer),
	}

	return p
}

func (d *DataAvailabilityProvider) Init() error {
	return nil
}

func (d *DataAvailabilityProvider) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	// Send blob to EigenDA disperser
	blobData := EncodeSequence(batchesData)
	_, idBytes, err := d.client.DisperseBlob(ctx, blobData, []uint8{})
	if err != nil {
		fmt.Println("failed to send blob to EigenDA disperser: ", err)
		return nil, err
	}
	fmt.Println("sent blob to EigenDA disperser, request id: ", base64.StdEncoding.EncodeToString(idBytes))

	var blobStatusReply *disperser_rpc.BlobStatusReply
	for {
		blobStatusReply, err = d.client.GetBlobStatus(ctx, idBytes)
		if err != nil {
			fmt.Printf("error getting blob status: %v\n", err)
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
		}

		// Wait period before retrieving blob status
		time.Sleep(d.cfg.RetrieveBlobStatusPeriod.Duration)
	}

	if blobStatusReply == nil {
		err = fmt.Errorf("empty blob status reply returned")
		return nil, err
	}

	// Get abi-encoded data availability message
	info := blobStatusReply.GetInfo()
	data, err := GetBlobData(info)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return TryEncodeToDataAvailabilityMessage(data)
}

func (d *DataAvailabilityProvider) GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	// Get blob from EigenDA layer
	var batchesData [][]byte
	for _, hash := range batchHashes {
		batchData, err := d.GetBatchL2Data(ctx, hash)
		if err != nil {
			return nil, err
		}
		batchesData = append(batchesData, batchData)
	}

	// Try decoding data availability message
	blobData, err := TryDecodeFromDataAvailabilityMessage(dataAvailabilityMessage)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	batchHeaderHash, err := blobData.BlobVerificationProof.BatchMetadata.GetBatchHeaderHash()
	if err != nil {
		return nil, err
	}

	reply, err := d.client.RetrieveBlob(ctx, batchHeaderHash, blobData.BlobVerificationProof.BlobIndex)
	if err != nil {
		fmt.Printf("failed to retrieve blob: %v\n", err)
		return nil, err
	}
	data, _ := DecodeSequence(reply.GetData())
	batchesData = append(batchesData, data...)

	return batchesData, nil
}

func (d *DataAvailabilityProvider) StoreDataAvailabilityMessage(ctx context.Context, batchHash common.Hash, dataAvailabilityMessage []byte) error {
	// Store blob information inside in-memory DA storage
	err := d.state.Add(batchHash, dataAvailabilityMessage)
	if err != nil {
		fmt.Printf("error adding data availability message into storage: %v\n", err)
		// Should not come here, but we will panic the mock node if indexing fails
		panic(err)
	}

	return nil
}

// GetBatchL2Data returns the data from the EigenDA layer operators. It checks the DA storage to get the
// requestID used when submitting the batch data to the DA.
func (d *DataAvailabilityProvider) GetBatchL2Data(ctx context.Context, hash common.Hash) ([]byte, error) {
	msg, err := d.state.Get(hash)
	if err != nil {
		fmt.Println("failed to get blob info from DA storage")
		return nil, err
	}
	blobData, err := TryDecodeFromDataAvailabilityMessage(msg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	batchHeaderHash, err := blobData.BlobVerificationProof.BatchMetadata.GetBatchHeaderHash()
	if err != nil {
		return nil, err
	}

	reply, err := d.client.RetrieveBlob(ctx, batchHeaderHash, blobData.BlobVerificationProof.BlobIndex)
	if err != nil {
		fmt.Printf("failed to retrieve blob: %v\n", err)
		return nil, err
	}
	data := reply.GetData()
	batchesData, batchesHash := DecodeSequence(data)

	// Get batch data from batch hash
	for idx, h := range batchesHash {
		if h == hash {
			return batchesData[idx], nil
		}
	}
	return nil, fmt.Errorf("failed to get batch data from hash, corrupted DA storage")
}

// Get data availability message from request ID
func (d *DataAvailabilityProvider) GetDataAvailabilityMessageFromId(ctx context.Context, requestId []byte) ([]byte, error) {
	blobStatusReply, err := d.client.GetBlobStatus(ctx, requestId)
	if err != nil {
		fmt.Printf("error getting blob status: %v\n", err)
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
		fmt.Println(err)
		return nil, err
	}

	return TryEncodeToDataAvailabilityMessage(blobData)
}
