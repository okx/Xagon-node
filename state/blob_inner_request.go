package state

import (
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
)

// ProcessBlobInnerProcessRequest is the request to process a blob
// you must use the builder to create the request
type ProcessBlobInnerProcessRequest struct {
	oldBlobStateRoot    common.Hash
	oldBlobAccInputHash common.Hash
	oldNumBlob          uint64
	oldStateRoot        common.Hash
	forkId              uint64
	lastL1InfoTreeIndex uint32
	lastL1InfoTreeRoot  common.Hash
	timestampLimit      uint64
	coinbase            common.Address
	zkGasLimit          uint64
	blobType            BlobType
}

// NewProcessBlobInnerProcessRequest creates a new ProcessBlobInnerProcessRequest
func NewProcessBlobInnerProcessRequest(forkid uint64, blob *BlobInner,
	previousSequence *BlobSequence,
	currentSequence BlobSequence) (*ProcessBlobInnerProcessRequest, error) {
	res := &ProcessBlobInnerProcessRequest{
		forkId:           forkid,
		blobType:         blob.Type,
		oldBlobStateRoot: ZeroHash, // Is always zero!
	}
	if previousSequence == nil {
		res.setAsFirstBlob()
	} else {
		res.setPreviousSequence(*previousSequence)
	}
	res.setBlob(blob)
	res.setCurrentSequence(currentSequence)
	return res, nil
}

func (p *ProcessBlobInnerProcessRequest) setAsFirstBlob() {
	p.oldBlobStateRoot = ZeroHash
	p.oldBlobAccInputHash = ZeroHash
	p.oldNumBlob = 0
	p.oldStateRoot = ZeroHash
}

func (p *ProcessBlobInnerProcessRequest) setCurrentSequence(seq BlobSequence) {
	p.coinbase = seq.L2Coinbase
}

func (p *ProcessBlobInnerProcessRequest) setPreviousSequence(previousSequence BlobSequence) {
	p.oldBlobAccInputHash = previousSequence.FinalAccInputHash
	p.oldNumBlob = previousSequence.LastBlobSequenced
}

func (p *ProcessBlobInnerProcessRequest) setBlob(blob *BlobInner) {
	p.lastL1InfoTreeIndex = blob.L1InfoLeafIndex
	p.lastL1InfoTreeRoot = blob.L1InfoTreeRoot
	p.timestampLimit = uint64(blob.MaxSequenceTimestamp.Unix()) // Convert time.Time to uint64
	p.zkGasLimit = blob.ZkGasLimit
}

func convertBlobInnerProcessRequestToExecutor(request ProcessBlobInnerProcessRequest, data []byte) *executor.ProcessBlobInnerRequestV3 {
	return &executor.ProcessBlobInnerRequestV3{
		OldBlobStateRoot:    request.oldBlobStateRoot.Bytes(),
		OldBlobAccInputHash: request.oldBlobAccInputHash.Bytes(),
		OldNumBlob:          request.oldNumBlob,
		OldStateRoot:        request.oldStateRoot.Bytes(),
		ForkId:              request.forkId,
		LastL1InfoTreeIndex: request.lastL1InfoTreeIndex,
		LastL1InfoTreeRoot:  request.lastL1InfoTreeRoot.Bytes(),
		TimestampLimit:      request.timestampLimit,
		Coinbase:            request.coinbase.String(),
		ZkGasLimit:          request.zkGasLimit,
		BlobType:            uint32(request.blobType),
		BlobData:            data,
	}
}
