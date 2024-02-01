package jsonrpc

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jackc/pgx/v4"
)

var debugEndPoints *DebugEndpoints
var once sync.Once

// GetInternalTransactions returns a transaction by his hash
func (e *EthEndpoints) GetInternalTransactions(hash types.ArgHash) (interface{}, types.Error) {
	if e.isDisabled("eth_getInternalTransactions") {
		return RPCErrorResponse(types.DefaultErrorCode, "not supported yet", nil, true)
	}
	once.Do(func() {
		debugEndPoints = &DebugEndpoints{
			state: e.state,
		}
	})
	return debugEndPoints.txMan.NewDbTxScope(debugEndPoints.state, func(ctx context.Context, dbTx pgx.Tx) (interface{}, types.Error) {
		ret, err := debugEndPoints.buildInnerTransaction(ctx, hash.Hash(), dbTx)
		if err != nil {
			return ret, err
		}

		jr, ok := ret.(json.RawMessage)
		if !ok {
			return nil, types.NewRPCError(types.ParserErrorCode, "cant transfer to json raw message")
		}
		r, stderr := jr.MarshalJSON()
		if stderr != nil {
			return nil, types.NewRPCError(types.ParserErrorCode, stderr.Error())
		}
		var of okFrame
		stderr = json.Unmarshal(r, &of)
		if stderr != nil {
			return nil, types.NewRPCError(types.ParserErrorCode, stderr.Error())
		}
		result := internalTxTraceToInnerTxs(of)

		return result, nil
	})
}

type okLog struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    hexutil.Bytes  `json:"data"`
}

type okFrame struct {
	Type         string          `json:"type"`
	From         common.Address  `json:"from"`
	Gas          string          `json:"gas"`
	GasUsed      string          `json:"gasUsed"`
	To           *common.Address `json:"to,omitempty" rlp:"optional"`
	Input        string          `json:"input" rlp:"optional"`
	Output       string          `json:"output,omitempty" rlp:"optional"`
	Error        string          `json:"error,omitempty" rlp:"optional"`
	RevertReason string          `json:"revertReason,omitempty"`
	Calls        []okFrame       `json:"calls,omitempty" rlp:"optional"`
	Logs         []okLog         `json:"logs,omitempty" rlp:"optional"`
	// Placed at end on purpose. The RLP will be decoded to 0 instead of
	// nil if there are non-empty elements after in the struct.
	Value string `json:"value,omitempty" rlp:"optional"`
}

func internalTxTraceToInnerTxs(tx okFrame) []*InnerTx {
	dfs := dfs{}
	indexMap := make(map[int]int)
	indexMap[0] = 1
	var level = 0
	var index = 1
	isError := tx.Error != ""
	dfs.dfs(tx, level, index, indexMap, isError)
	return dfs.innerTxs
}

type dfs struct {
	innerTxs []*InnerTx
}

func inArray(dst string, src []string) bool {
	for _, v := range src {
		if v == dst {
			return true
		}
	}
	return false
}

func (d *dfs) dfs(tx okFrame, level int, index int, indexMap map[int]int, isError bool) {
	if !inArray(strings.ToLower(tx.Type), []string{"call", "create", "create2",
		"callcode", "delegatecall", "staticcall", "selfdestruct"}) {
		return
	}
	name := strings.ToLower(tx.Type)
	for i := 0; i < level; i++ {
		if indexMap[i] == 0 {
			continue
		}
		name = name + "_" + strconv.Itoa(indexMap[i])
	}
	innerTx := internalTxTraceToInnerTx(tx, name, level, index)
	if !isError {
		isError = innerTx.IsError
	} else {
		innerTx.IsError = isError
	}
	d.innerTxs = append(d.innerTxs, innerTx)
	index = 0
	for _, call := range tx.Calls {
		index++
		indexMap[level] = index
		d.dfs(call, level+1, index+1, indexMap, isError)
	}
	if len(tx.Calls) == 0 {
		return
	}
}

// InnerTx represents a struct type for internal transactions.
type InnerTx struct {
	Dept          big.Int `json:"dept"`
	InternalIndex big.Int `json:"internal_index"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	Input         string  `json:"input"`
	Output        string  `json:"output"`
	IsError       bool    `json:"is_error"`
	GasUsed       uint64  `json:"gas_used"`
	Value         string  `json:"value"`
	ValueWei      string  `json:"value_wei"`
	CallValueWei  string  `json:"call_value_wei"`
	Error         string  `json:"error"`
	Gas           uint64  `json:"gas"`
	//ReturnGas    uint64 `json:"return_gas"`
	CallType     string `json:"call_type"`
	Name         string `json:"name"`
	TraceAddress string `json:"trace_address"`
	CodeAddress  string `json:"code_address"`
}

func internalTxTraceToInnerTx(currentTx okFrame, name string, depth int, index int) *InnerTx {
	value := currentTx.Value
	if value == "" {
		value = "0x0"
	}
	var toAddress string
	if currentTx.To != nil {
		toAddress = currentTx.To.String()
	}
	gas, _ := strconv.ParseUint(currentTx.Gas, 0, 64)
	gasUsed, _ := strconv.ParseUint(currentTx.GasUsed, 0, 64)
	valueWei, _ := hexutil.DecodeBig(value)
	callTx := &InnerTx{
		Dept:         *big.NewInt(int64(depth)),
		From:         currentTx.From.String(),
		To:           toAddress,
		ValueWei:     valueWei.String(),
		CallValueWei: value,
		CallType:     strings.ToLower(currentTx.Type),
		Name:         name,
		Input:        currentTx.Input,
		Output:       currentTx.Output,
		Gas:          gas,
		GasUsed:      gasUsed,
		IsError:      false, // TODO Nested errors
		//ReturnGas:    currentTx.Gas - currentTx.GasUsed,
	}
	callTx.InternalIndex = *big.NewInt(int64(index - 1))
	if strings.ToLower(currentTx.Type) == "callcode" {
		callTx.CodeAddress = currentTx.To.String()
	}
	if strings.ToLower(currentTx.Type) == "delegatecall" {
		callTx.ValueWei = ""
	}
	if currentTx.Error != "" {
		callTx.Error = currentTx.Error
		callTx.IsError = true
	}
	return callTx
}
