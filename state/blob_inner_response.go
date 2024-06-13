package state

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ProcessBlobInnerResponse is the response of the process of a blob
// the fields are private, so you need the function to access the data
// To get the outcome of the execution you must use GetSuccesfulData() it will return a nil if the execution was not successful
// This is for forcing by interface don't access to results fields is it have an error
type ProcessBlobInnerResponse struct {
	succesfulData   ProcessBlobInnerResponseSuccesful // Here is the outcome of the execution
	isInvalid       bool                              // Is a variable of the ROM
	generalError    error
	romBlobError    error
	executorVersion string // Version of the executor e.g. "v7.0.0"
	errorDebugLog   string // This is debug.ErrorLog that is a debug string with context data of error
}

// ProcessBlobInnerResponseSuccesful is the data after a  successful call to ProcessBlobInner
type ProcessBlobInnerResponseSuccesful struct {
	newBlobStateRoot      common.Hash
	newBlobAccInputHash   common.Hash
	newNumBlob            uint64
	finalAccBatchHashData common.Hash
	batchData             [][]byte
}

func (p *ProcessBlobInnerResponseSuccesful) String() string {
	res := fmt.Sprintf("newBlobStateRoot: %s newBlobAccInputHash:%s newNumBlob:%d\n", p.newBlobStateRoot.String(), p.newBlobAccInputHash.String(), p.newNumBlob)
	res += fmt.Sprintf("finalAccBatchHashData: %s\n", p.finalAccBatchHashData.String())
	res += fmt.Sprintf("HowManyBatches: %d\n", p.HowManyBatches())
	for i := 0; i < p.HowManyBatches(); i++ {
		res += fmt.Sprintf("            Batch %d: Hash:%s\n", i, p.GetBatchHash(i).String())
	}
	return res
}

func (p *ProcessBlobInnerResponse) String() string {
	res := fmt.Sprintf("isInvalid: %t\n", p.isInvalid)
	if p.generalError != nil {
		res += fmt.Sprintf("generalError: %s\n", p.generalError.Error())
	}
	if p.romBlobError != nil {
		res += fmt.Sprintf("romBlobError: %s\n", p.romBlobError.Error())
	}
	if p.IsSuccessfulExecution() {
		res += p.succesfulData.String()
	}
	return res
}

// GetUnifiedError returns the combinations of errors of the execution
func (p *ProcessBlobInnerResponse) GetUnifiedError() error {
	if p.IsSuccessfulExecution() {
		return nil
	}
	return fmt.Errorf("ProcessBlobInnerV3 fails:version:%s  isInvalid: %t  general:%w romBlob:%w  errorLog:%s",
		p.executorVersion, p.isInvalid, p.generalError, p.romBlobError, p.errorDebugLog)
}

// IsSuccessfulExecution returns true if the execution was successful
func (p *ProcessBlobInnerResponse) IsSuccessfulExecution() bool {
	return !p.isInvalid && p.generalError == nil && p.romBlobError == nil
}

// GetSuccesfulData returns the outcome data of the execution
func (p *ProcessBlobInnerResponse) GetSuccesfulData() *ProcessBlobInnerResponseSuccesful {
	if !p.IsSuccessfulExecution() {
		log.Error("Trying to get successful data from a failed execution")
		return nil
	}
	return &p.succesfulData
}

// HowManyBatches returns the number of batches
func (p *ProcessBlobInnerResponseSuccesful) HowManyBatches() int {
	return len(p.batchData)
}

// GetBatchData returns the data of the batch
func (p *ProcessBlobInnerResponseSuccesful) GetBatchData(index int) []byte {
	return p.batchData[index]
}

// GetBatchHash returns the hash of the batch data
func (p *ProcessBlobInnerResponseSuccesful) GetBatchHash(index int) common.Hash {
	return crypto.Keccak256Hash(p.GetBatchData(index))
}

func newProcessBlobInnerProcessResponse(response *executor.ProcessBlobInnerResponseV3) *ProcessBlobInnerResponse {
	res := &ProcessBlobInnerResponse{
		succesfulData: ProcessBlobInnerResponseSuccesful{
			newBlobStateRoot:      common.BytesToHash(response.NewBlobStateRoot),
			newBlobAccInputHash:   common.BytesToHash(response.NewBlobAccInputHash),
			newNumBlob:            response.NewNumBlob,
			finalAccBatchHashData: common.BytesToHash(response.FinalAccBatchHashData),
			batchData:             response.BatchData,
		},
		isInvalid:    response.IsInvalid == cTrue,
		generalError: executor.ExecutorErr(response.Error),
		romBlobError: executor.RomBlobErr(response.ErrorRomBlob),
	}
	if response.Debug != nil {
		res.executorVersion = response.Debug.Version
		res.errorDebugLog = response.Debug.ErrorLog
	}
	return res
}
