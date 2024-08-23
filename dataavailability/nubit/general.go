package nubit

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	"github.com/RiemaLabs/nuport-offchain/das"
	"github.com/ethereum/go-ethereum/common"
)

type DaClient interface {
	das.NubitWriter
	das.NubitReader
}

type GeneralDA struct {
	client DaClient
}

func NewGeneralDA(cfg *Config, ethRpc string) (*GeneralDA, error) {
	c := &das.DAConfig{
		Rpc:         cfg.NubitRpcURL,
		AuthToken:   cfg.NubitAuthKey,
		GasPrice:    0.01,
		NamespaceId: cfg.NubitNamespace,
		ValidatorConfig: &das.ValidatorConfig{
			TendermintRPC:  cfg.NubitValidatorURL,
			BlobstreamAddr: "", // TODO: add blobstream address
			VerifyAddr:     "", // TODO: add verify address
			EthClient:      ethRpc,
		},
	}
	da, err := das.NewNubitDA(c, nil)
	if err != nil {
		return nil, err
	}
	return &GeneralDA{
		client: da,
	}, nil
}

func (s *GeneralDA) Init() error {
	return nil
}

func (s *GeneralDA) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	// lastCommitTime := time.Since(backend.commitTime)
	// if lastCommitTime < NubitMinCommitTime {
	// 	time.Sleep(NubitMinCommitTime - lastCommitTime)
	// }

	// Encode NubitDA blob data
	data := EncodeSequence(batchesData)

	return s.client.Store(ctx, data)
}

func (s *GeneralDA) GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {

	var flag byte

	buf := bytes.NewReader(dataAvailabilityMessage)
	binary.Read(buf, binary.BigEndian, &flag)
	if flag != das.NubitMessageHeaderFlag {
		return nil, fmt.Errorf("invalid data availability message flag, %d", flag)
	}

	bpBinary := make([]byte, len(dataAvailabilityMessage)-1)

	bp := &das.BlobPointer{}
	if err := bp.UnmarshalBinary(bpBinary); err != nil {
		return nil, err
	}

	result, err := s.client.Read(ctx, bp)
	if err != nil {
		return nil, err
	}

	batchData, hashs := DecodeSequence(result.Message)
	if len(hashs) != len(batchHashes) {
		return nil, fmt.Errorf("invalid L2 batch data retrieval arguments, %d != %d", len(hashs), len(batchHashes))
	}

	for i, bh := range batchHashes {
		if bh != hashs[i] {
			return nil, fmt.Errorf("invalid L2 batch data hash retrieval arguments, %s != %s", bh.String(), hashs[i].String())
		}
	}

	return batchData, nil
}
