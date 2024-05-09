package state

import (
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
)

// TestConvertToProcessBatchResponseV3 for test purposes
func (s *State) TestConvertToProcessBatchResponseV3(batchResponse *executor.ProcessBatchResponseV3) (*ProcessBatchResponse, error) {
	return s.convertToProcessBatchResponseV3(batchResponse)
}

func (s *State) convertToProcessBatchResponseV3(batchResponse *executor.ProcessBatchResponseV3) (*ProcessBatchResponse, error) {
	blockResponses, isRomLevelError, isRomOOCError, err := s.convertToProcessBlockResponseV2(batchResponse.BlockResponses)
	if err != nil {
		return nil, err
	}
	isRomOOCError = isRomOOCError || executor.IsROMOutOfCountersError(batchResponse.ErrorRom)
	readWriteAddresses, err := convertToReadWriteAddressesV2(batchResponse.ReadWriteAddresses)
	if err != nil {
		return nil, err
	}

	return &ProcessBatchResponse{
		NewStateRoot:              common.BytesToHash(batchResponse.NewStateRoot),
		NewAccInputHash:           common.BytesToHash(batchResponse.NewAccInputHash),
		NewLocalExitRoot:          common.BytesToHash(batchResponse.NewLocalExitRoot),
		UsedZkCounters:            convertToUsedZKCountersV3(batchResponse),
		ReservedZkCounters:        convertToReservedZKCountersV3(batchResponse),
		BlockResponses:            blockResponses,
		ExecutorError:             executor.ExecutorErr(batchResponse.Error),
		ReadWriteAddresses:        readWriteAddresses,
		FlushID:                   batchResponse.FlushId,
		StoredFlushID:             batchResponse.StoredFlushId,
		ProverID:                  batchResponse.ProverId,
		IsExecutorLevelError:      batchResponse.Error != executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR,
		IsRomLevelError:           isRomLevelError,
		IsRomOOCError:             isRomOOCError,
		GasUsed_V2:                batchResponse.GasUsed,
		SMTKeys_V2:                convertToKeys(batchResponse.SmtKeys),
		ProgramKeys_V2:            convertToKeys(batchResponse.ProgramKeys),
		ForkID:                    batchResponse.ForkId,
		InvalidBatch_V2:           batchResponse.InvalidBatch != 0,
		RomError_V2:               executor.RomErr(batchResponse.ErrorRom),
		OldStateRoot_V2:           common.BytesToHash(batchResponse.OldStateRoot),
		NewLastTimestamp_V3:       batchResponse.NewLastTimestamp,
		CurrentL1InfoTreeRoot_V3:  common.BytesToHash(batchResponse.CurrentL1InfoTreeRoot),
		CurrentL1InfoTreeIndex_V3: batchResponse.CurrentL1InfoTreeIndex,
	}, nil
}

func convertToUsedZKCountersV3(resp *executor.ProcessBatchResponseV3) ZKCounters {
	return ZKCounters{
		GasUsed:          resp.GasUsed,
		KeccakHashes:     resp.CntKeccakHashes,
		PoseidonHashes:   resp.CntPoseidonHashes,
		PoseidonPaddings: resp.CntPoseidonPaddings,
		MemAligns:        resp.CntMemAligns,
		Arithmetics:      resp.CntArithmetics,
		Binaries:         resp.CntBinaries,
		Steps:            resp.CntSteps,
		Sha256Hashes_V2:  resp.CntSha256Hashes,
	}
}

func convertToReservedZKCountersV3(resp *executor.ProcessBatchResponseV3) ZKCounters {
	return ZKCounters{
		// There is no "ReserveGasUsed" in the response, so we use "GasUsed" as it will make calculations easier
		GasUsed:          resp.GasUsed,
		KeccakHashes:     resp.CntReserveKeccakHashes,
		PoseidonHashes:   resp.CntReservePoseidonHashes,
		PoseidonPaddings: resp.CntReservePoseidonPaddings,
		MemAligns:        resp.CntReserveMemAligns,
		Arithmetics:      resp.CntReserveArithmetics,
		Binaries:         resp.CntReserveBinaries,
		Steps:            resp.CntReserveSteps,
		Sha256Hashes_V2:  resp.CntReserveSha256Hashes,
	}
}
