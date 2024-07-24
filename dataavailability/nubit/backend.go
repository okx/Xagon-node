package nubit

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"time"

	daTypes "github.com/0xPolygon/cdk-data-availability/types"
	polygondatacommittee "github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygondatacommittee_xlayer"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rollkit/go-da"
	"github.com/rollkit/go-da/proxy"
)

// NubitDABackend implements the DA integration with Nubit DA layer
type NubitDABackend struct {
	dataCommitteeContract *polygondatacommittee.PolygondatacommitteeXlayer
	client                da.DA
	config                *Config
	ns                    da.Namespace
	privKey               *ecdsa.PrivateKey
	commitTime            time.Time
	batchesDataCache      [][]byte
	batchesDataSize       uint64
}

// NewNubitDABackend is the factory method to create a new instance of NubitDABackend
func NewNubitDABackend(
	l1RPCURL string,
	dataCommitteeAddr common.Address,
	privKey *ecdsa.PrivateKey,
	cfg *Config,
) (*NubitDABackend, error) {
	ethClient, err := ethclient.Dial(l1RPCURL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", l1RPCURL, err)
		return nil, err
	}

	log.Infof("NubitDABackend config: %#v ", cfg)
	cn, err := proxy.NewClient(cfg.NubitRpcURL, cfg.NubitAuthKey)
	if err != nil {
		return nil, err
	}
	// TODO: Check if name byte array requires zero padding
	name, err := hex.DecodeString(cfg.NubitNamespace)
	if err != nil {
		log.Errorf("error decoding NubitDA namespace config: %+v", err)
		return nil, err
	}
	dataCommittee, err := polygondatacommittee.NewPolygondatacommitteeXlayer(dataCommitteeAddr, ethClient)
	if err != nil {
		return nil, err
	}
	log.Infof("NubitDABackend namespace: %s ", string(name))

	return &NubitDABackend{
		dataCommitteeContract: dataCommittee,
		config:                cfg,
		privKey:               privKey,
		ns:                    name,
		client:                cn,
		commitTime:            time.Now(),
		batchesDataCache:      [][]byte{},
		batchesDataSize:       0,
	}, nil
}

// Init initializes the NubitDA backend
func (backend *NubitDABackend) Init() error {
	return nil
}

// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
// as expected by the contract
func (backend *NubitDABackend) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	encodedData, err := MarshalBatchData(batchesData)
	if err != nil {
		log.Errorf("Marshal batch data failed: %s", err)
		return nil, err
	}

	// Add to batches data cache
	backend.batchesDataCache = append(backend.batchesDataCache, encodedData)
	backend.batchesDataSize += uint64(len(encodedData))
	if backend.batchesDataSize < backend.config.NubitMaxBatchesSize {
		log.Infof("Added batches data to NubitDABackend cache, current length: %+v", len(encodedData))
		return nil, nil
	}
	if time.Since(backend.commitTime) < 12*time.Second {
		time.Sleep(time.Since(backend.commitTime))
	}

	data, err := MarshalBatchData(backend.batchesDataCache)
	if err != nil {
		log.Errorf("Marshal batch data failed: %s", err)
		return nil, err
	}
	id, err := backend.client.Submit(ctx, [][]byte{data}, -1, backend.ns)
	if err != nil {
		log.Errorf("Submit batch data with NubitDA client failed: %s", err)
		return nil, err
	}
	log.Infof("Data submitted to Nubit DA: %d bytes against namespace %v sent with id %#x", len(backend.batchesDataCache), backend.ns, id)

	// Reset batches data cache and DA commit time
	backend.commitTime = time.Now()
	backend.batchesDataCache = [][]byte{}
	backend.batchesDataSize = 0

	// Get proof
	tries := uint64(0)
	posted := false
	for tries < backend.config.NubitGetProofMaxRetry {
		dataProof, err := backend.client.GetProofs(ctx, id, backend.ns)
		if err != nil {
			log.Infof("Proof not available: %s", err)
		}
		if len(dataProof) > 0 {
			log.Infof("Data proof from Nubit DA received: %+v", dataProof)
			posted = true
			break
		}

		tries += 1
		time.Sleep(backend.config.NubitGetProofWaitPeriod)
	}
	if !posted {
		log.Errorf("Get blob proof on Nubit DA failed: %s", err)
		return nil, err
	}

	// // TODO: use bridge API data
	// batchDAData := BatchDAData{ID: id}
	// log.Infof("Nubit DA data ID: %+v", batchDAData)
	// returnData, err := batchDAData.Encode()
	// if err != nil {
	// 	return nil, fmt.Errorf("Encode batch data failed: %w", err)
	// }

	// Sign sequence
	sequence := daTypes.Sequence{}
	for _, seq := range batchesData {
		sequence = append(sequence, seq)
	}
	signedSequence, err := sequence.Sign(backend.privKey)
	if err != nil {
		log.Errorf("Failed to sign sequence with pk: %v", err)
		return nil, err
	}
	signature := append(sequence.HashToSign(), signedSequence.Signature...)

	return signature, nil
}

// GetSequence gets the sequence data from NubitDA layer
func (backend *NubitDABackend) GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	batchDAData := BatchDAData{}
	err := batchDAData.Decode(dataAvailabilityMessage)
	if err != nil {
		log.Errorf("🏆    NubitDABackend.GetSequence.Decode:%s", err)
		return nil, err
	}
	log.Infof("🏆     Nubit GetSequence batchDAData:%+v", batchDAData)
	blob, err := backend.client.Get(ctx, batchDAData.ID, backend.ns)
	if err != nil {
		log.Errorf("🏆    NubitDABackend.GetSequence.Blob.Get:%s", err)
		return nil, err
	}
	log.Infof("🏆     Nubit GetSequence blob.data:%+v", len(blob))

	return blob, nil
}
