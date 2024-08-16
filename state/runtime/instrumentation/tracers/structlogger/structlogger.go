package structlogger

import (
	"encoding/json"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/fakevm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Config are the configuration options for structured logger the EVM
type Config struct {
	EnableMemory     bool // enable memory capture
	DisableStack     bool // disable stack capture
	DisableStorage   bool // disable storage capture
	EnableReturnData bool // enable return data capture
}

// StructLogRes represents the debug trace information for each opcode
type StructLogRes struct {
	Pc            uint64             `json:"pc"`
	Op            string             `json:"op"`
	Gas           uint64             `json:"gas"`
	GasCost       uint64             `json:"gasCost"`
	Depth         int                `json:"depth"`
	Error         string             `json:"error,omitempty"`
	Stack         *[]string          `json:"stack,omitempty"`
	ReturnData    *string            `json:"returnData,omitempty"`
	Memory        *[]string          `json:"memory,omitempty"`
	Storage       *map[string]string `json:"storage,omitempty"`
	RefundCounter uint64             `json:"refund,omitempty"`
}

type TraceResponse struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue interface{}    `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

type JSONLogger struct {
	cfg Config
}

func NewStructLogger(cfg Config) *JSONLogger {
	return &JSONLogger{cfg}
}

func (l *JSONLogger) ParseTrace(result *runtime.ExecutionResult, receipt types.Receipt) (json.RawMessage, error) {
	structLogs := make([]StructLogRes, 0, len(result.FullTrace.Steps))
	memory := fakevm.NewMemory()
	for _, step := range result.FullTrace.Steps {
		errRes := ""
		if step.Error != nil {
			errRes = step.Error.Error()
		}

		op := step.OpCode
		if op == "SHA3" {
			op = "KECCAK256"
		} else if op == "STOP" && step.Pc == 0 {
			// this stop is generated for calls with single
			// step(no depth increase) and must be ignored
			continue
		}

		structLogRes := StructLogRes{
			Pc:            step.Pc,
			Op:            op,
			Gas:           step.Gas,
			GasCost:       step.GasCost,
			Depth:         step.Depth,
			Error:         errRes,
			RefundCounter: step.Refund,
		}

		if !l.cfg.DisableStack {
			stack := make([]string, 0, len(step.Stack))
			for _, stackItem := range step.Stack {
				if stackItem != nil {
					stack = append(stack, hex.EncodeBig(stackItem))
				}
			}
			structLogRes.Stack = &stack
		}

		if l.cfg.EnableMemory {
			memory.Resize(uint64(step.MemorySize))
			if len(step.Memory) > 0 {
				memory.Set(uint64(step.MemoryOffset), uint64(len(step.Memory)), step.Memory)
			}

			if step.MemorySize > 0 {
				// Populate the structLog memory
				step.Memory = memory.Data()

				// Convert memory to string array
				const memoryChunkSize = 32
				memoryArray := make([]string, 0, len(step.Memory))

				for i := 0; i < len(step.Memory); i = i + memoryChunkSize {
					slice32Bytes := make([]byte, memoryChunkSize)
					copy(slice32Bytes, step.Memory[i:i+memoryChunkSize])
					memoryStringItem := hex.EncodeToString(slice32Bytes)
					memoryArray = append(memoryArray, memoryStringItem)
				}

				structLogRes.Memory = &memoryArray
			} else {
				memory = fakevm.NewMemory()
				structLogRes.Memory = &[]string{}
			}
		}

		if !l.cfg.DisableStorage && len(step.Storage) > 0 {
			storage := make(map[string]string, len(step.Storage))
			for storageKey, storageValue := range step.Storage {
				k := hex.EncodeToString(storageKey.Bytes())
				v := hex.EncodeToString(storageValue.Bytes())
				storage[k] = v
			}
			structLogRes.Storage = &storage
		}

		var returnData *string
		if l.cfg.EnableReturnData && len(step.ReturnData) > 0 {
			rd := hex.EncodeToHex(step.ReturnData)
			returnData = &rd
		}

		structLogRes.ReturnData = returnData

		structLogs = append(structLogs, structLogRes)
	}

	var rv interface{}
	if l.cfg.EnableReturnData {
		rv = common.Bytes2Hex(result.ReturnValue)
	}

	failed := receipt.Status == types.ReceiptStatusFailed

	resp := TraceResponse{
		Gas:         receipt.GasUsed,
		Failed:      failed,
		ReturnValue: rv,
		StructLogs:  structLogs,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return b, nil
}
