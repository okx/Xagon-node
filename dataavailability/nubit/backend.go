package nubit

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rollkit/go-da"
	"github.com/rollkit/go-da/proxy"
)

// NubitDABackend implements the DA integration with Nubit DA layer
type NubitDABackend struct {
	client     da.DA
	validator  Validator
	config     *Config
	namespace  da.Namespace
	privKey    *ecdsa.PrivateKey
	commitTime time.Time
}

const heightLen = 8

func MakeID(height []byte, commitment da.Commitment) da.ID {
	id := make([]byte, heightLen+len(commitment))
	copy(id[:heightLen], height)
	copy(id[heightLen:], commitment)
	return id
}

// NewNubitDABackend is the factory method to create a new instance of NubitDABackend
func NewNubitDABackend(
	cfg *Config,
	privKey *ecdsa.PrivateKey,
) (*NubitDABackend, error) {
	log.Infof("NubitDABackend config: %#v", cfg)
	cn, err := proxy.NewClient(cfg.NubitRpcURL, cfg.NubitAuthKey)
	if err != nil {
		return nil, err
	}

	va := NewJSONRPCValidator(cfg.NubitValidatorURL)

	hexStr := hex.EncodeToString([]byte(cfg.NubitNamespace))
	name, err := hex.DecodeString(strings.Repeat("0", NubitNamespaceBytesLength-len(hexStr)) + hexStr)
	if err != nil {
		log.Errorf("error decoding NubitDA namespace config: %+v", err)
		return nil, err
	}
	log.Infof("NubitDABackend namespace: %s", string(name))

	return &NubitDABackend{
		config:     cfg,
		privKey:    privKey,
		namespace:  name,
		client:     cn,
		validator:  va,
		commitTime: time.Now(),
	}, nil
}

// Init initializes the NubitDA backend
func (backend *NubitDABackend) Init() error {
	return nil
}

// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
// as expected by the contract
func (backend *NubitDABackend) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	// Check commit time interval validation
	lastCommitTime := time.Since(backend.commitTime)
	if lastCommitTime < NubitMinCommitTime {
		time.Sleep(NubitMinCommitTime - lastCommitTime)
	}

	// Encode NubitDA blob data
	data := EncodeSequence(batchesData)

	// Submit the blob.
	// TODO: Add the gas estimation function.
	blobIDs, err := backend.client.Submit(ctx, [][]byte{data}, 0.01, backend.namespace)

	// Ensure only a single blobID returned
	if err != nil || len(blobIDs) != 1 {
		log.Errorf("Submit batch data with NubitDA client failed: %s", err)
		return nil, err
	}

	blobID := blobIDs[0]
	height := blobID[:heightLen]
	log.Infof("Current Nubit height: %d", binary.LittleEndian.Uint64(height[:]))
	commitment := blobID[heightLen:]

	// blobID implementation:
	// height: 8 bytes
	// commitment: 32 bytes

	backend.commitTime = time.Now()
	log.Infof("Data submitted to Nubit DA: %d bytes against namespace %v sent with id %#x", len(data), backend.namespace, blobID)

	// Get proof of batches data on NubitDA layer
	tries := uint64(0)
	posted := false
	proofs := make([]*NamespaceRangeProof, 0)
	for tries < backend.config.NubitGetProofMaxRetry {
		time.Sleep(backend.config.NubitGetProofWaitPeriod.Duration)
		inclusionProofBytes, err := backend.client.GetProofs(ctx, [][]byte{blobID}, backend.namespace)
		if err != nil {
			log.Errorf("Inclusion proof is not available: %s", err)
			// Retries
			tries += 1
		} else {
			// Never retry if the interface is invalid/incompatible.
			proofs, err = UnmarshalInclusionProofs(inclusionProofBytes)
			if err != nil {
				log.Errorf("Inclusion proof is not valid: %s", err)
				return nil, err
			} else {
				posted = true
				break
			}
		}
	}
	if !posted {
		log.Errorf("Get inclusion proof on Nubit DA failed: %s", err)
		return nil, err
	}

	// Get the row proof via the NamespaceRangeProof.

	sharesProofs := make([]ShareProof, 0)
	for _, p := range proofs {
		height_i := int64(binary.LittleEndian.Uint64(height[:]))
		start := uint64(p.Start)
		end := uint64(p.End)
		sharesProof, err := backend.validator.ProveShares(height_i, start, end)
		if err != nil {
			log.Errorf("Get shares proof on Nubit Validator failed: %s", err)
			return nil, err
		}
		sharesProofs = append(sharesProofs, *sharesProof)
	}

	// TODO: Finish the proof convertion in the polygon context.

	daMessage := BlobData{
		NubitHeight: height,
		Commitment:  commitment,
		SharesProof: make([]byte, 0),
	}

	return TryEncodeToDataAvailabilityMessage(daMessage)
}

// GetSequence gets the sequence data from NubitDA layer
func (backend *NubitDABackend) GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	blobData, err := TryDecodeFromDataAvailabilityMessage(dataAvailabilityMessage)
	if err != nil {
		log.Error("Error decoding from da message: ", err)
		return nil, err
	}

	blobID := MakeID(blobData.NubitHeight, blobData.NubitHeight)

	reply, err := backend.client.Get(ctx, [][]byte{blobID}, backend.namespace)
	if err != nil || len(reply) != 1 {
		log.Error("Error retrieving blob from NubitDA client: ", err)
		return nil, err
	}

	batchesData, _ := DecodeSequence(reply[0])
	return batchesData, nil
}
