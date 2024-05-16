package celestia

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
	openrpc "github.com/celestiaorg/celestia-openrpc"
	"github.com/celestiaorg/celestia-openrpc/types/blob"
	"github.com/celestiaorg/celestia-openrpc/types/share"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const unexpectedHashTemplate = "mismatch on transaction data. Expected hash %s, actual hash: %s"

// Config Configuration for the CelestiaDA node
type Config struct {
	// L1RpcUrl           string
	GasPrice float64
	Rpc      string
	// TendermintRPC      string
	NamespaceId string
	AuthToken   string
	// BlobstreamXAddress string
	// EventChannelSize   uint64
}

// Batch Types to help with Batch Data Aggregation
//type Batch struct {
//	ChangeL2BlockMarker uint8
//	DeltaTimestamp      uint32
//	IndexL1InfoTree     uint32
//	Transactions        []Transaction
//}
//
//type Transaction struct {
//	RLP               []byte
//	R, S, V           []byte
//	EfficiencyPercent byte
//}

// CelestiaBackend implements the Celestia integration
type CelestiaBackend struct {
	Cfg    Config
	Client *openrpc.Client
	// Trpc        *http.HTTP
	Namespace share.Namespace
	// BlobstreamX *blobstreamx.BlobstreamX
}

func New(cfg Config) (*CelestiaBackend, error) {
	if cfg.NamespaceId == "" {
		return nil, errors.New("namespace id cannot be blank")
	}
	nsBytes, err := hex.DecodeString(cfg.NamespaceId)
	if err != nil {
		return nil, err
	}

	namespace, err := share.NewBlobNamespaceV0(nsBytes)
	if err != nil {
		return nil, err
	}

	// var trpc *http.HTTP
	// trpc, err = http.New(cfg.TendermintRPC, "/websocket")

	// if cfg.EventChannelSize == 0 {
	// 	cfg.EventChannelSize = 100
	// }

	return &CelestiaBackend{
		Cfg:       cfg,
		Client:    nil,
		Namespace: namespace,
	}, nil
}

func (c *CelestiaBackend) Init() error {
	daClient, err := openrpc.NewClient(context.Background(), c.Cfg.Rpc, c.Cfg.AuthToken)
	if err != nil {
		return err
	}
	c.Client = daClient
	return nil
}

func (c *CelestiaBackend) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	var aggregatedBatch []byte
	for _, batchData := range batchesData {
		aggregatedBatch = append(aggregatedBatch, batchData...)
	}
	dataBlob, err := blob.NewBlobV0(c.Namespace, aggregatedBatch)
	if err != nil {
		log.Warn("Error creating blob", "err", err)
		return nil, err
	}

	commitment, err := blob.CreateCommitment(dataBlob)
	if err != nil {
		log.Warn("Error creating commitment", "err", err)
		return nil, err
	}

	height, err := c.Client.Blob.Submit(ctx, []*blob.Blob{dataBlob}, openrpc.GasPrice(c.Cfg.GasPrice))
	if err != nil {
		log.Warn("Blob Submission error", "err", err)
		return nil, err
	}

	if height == 0 {
		log.Warn("Unexpected height from blob response", "height", height)
		return nil, errors.New("unexpected response code")
	}

	proofs, err := c.Client.Blob.GetProof(ctx, height, c.Namespace, commitment)
	if err != nil {
		log.Warn("Error retrieving proof", "err", err)
		return nil, err
	}

	included, err := c.Client.Blob.Included(ctx, height, c.Namespace, proofs, commitment)
	if err != nil || !included {
		log.Warn("Error checking for inclusion", "err", err, "proof", proofs)
		return nil, err
	}
	log.Info("Successfully posted blob", "height", height, "commitment", hex.EncodeToString(commitment))

	txCommitment := [32]byte{}
	copy(txCommitment[:], commitment)

	blobPointer := BlobPointer{
		BlockHeight:  height,
		TxCommitment: txCommitment,
	}

	blobPointerData, err := blobPointer.MarshalBinary()
	if err != nil {
		log.Warn("BlobPointer MashalBinary error", "err", err)
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, blobPointerData)
	if err != nil {
		log.Warn("blob pointer data serialization failed", "err", err)
		return nil, err
	}

	// will need to wait for Blobstream or Equivalency service
	return buf.Bytes(), nil
}

func (c *CelestiaBackend) GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	blobPointer := BlobPointer{}
	err := blobPointer.UnmarshalBinary(dataAvailabilityMessage)
	if err != nil {
		return nil, err
	}
	dataBlob, err := c.Client.Blob.Get(ctx, blobPointer.BlockHeight, c.Namespace, blobPointer.TxCommitment[:])
	if err != nil {
		return nil, err
	}

	aggregatedBatchData := dataBlob.Data

	// need to read data using the dataAvailabilityMessage
	decodedBatches, err := decodeBatches(aggregatedBatchData)
	if err != nil {
		fmt.Printf("Error decoding batches: %v\n", err)
		return nil, err
	}

	if len(decodedBatches) != len(batchHashes) {
		return nil, fmt.Errorf("error decoded batches length and batch hashes length mismatch")
	}

	for i, batchData := range decodedBatches {
		actualTransactionsHash := crypto.Keccak256Hash(batchData)
		if actualTransactionsHash != batchHashes[i] {
			unexpectedHash := fmt.Errorf(
				unexpectedHashTemplate, batchHashes[i], actualTransactionsHash,
			)
			log.Warnf(
				"error getting data from Celestia node at hegith %s with commitment %s: %s",
				blobPointer.BlockHeight, blobPointer.TxCommitment, unexpectedHash,
			)
		}
	}
	return decodedBatches, nil
}

func decodeBatches(data []byte) ([][]byte, error) {
	var batches [][]byte
	offset := 0

	for offset < len(data) {
		// Skip the changeL2BlockMarker
		offset++

		// Read the DeltaTimestamp and IndexL1InfoTree
		_ = binary.BigEndian.Uint32(data[offset : offset+4]) // DeltaTimestamp
		offset += 4
		_ = binary.BigEndian.Uint32(data[offset : offset+4]) // IndexL1InfoTree
		offset += 4

		var batchBytes []byte
		for offset < len(data) {
			rlpLen := binary.BigEndian.Uint32(data[offset : offset+4])
			offset += 4
			rlpEnd := offset + int(rlpLen)
			batchBytes = append(batchBytes, data[offset:rlpEnd]...)
			offset = rlpEnd

			rLen := binary.BigEndian.Uint32(data[offset : offset+4])
			offset += 4
			rEnd := offset + int(rLen)
			batchBytes = append(batchBytes, data[offset:rEnd]...)
			offset = rEnd

			sLen := binary.BigEndian.Uint32(data[offset : offset+4])
			offset += 4
			sEnd := offset + int(sLen)
			batchBytes = append(batchBytes, data[offset:sEnd]...)
			offset = sEnd

			vLen := binary.BigEndian.Uint32(data[offset : offset+4])
			offset += 4
			vEnd := offset + int(vLen)
			batchBytes = append(batchBytes, data[offset:vEnd]...)
			offset = vEnd

			batchBytes = append(batchBytes, data[offset]) // EfficiencyPercent
			offset++

			if offset >= len(data) {
				break
			}
		}

		batches = append(batches, batchBytes)
	}

	return batches, nil
}
