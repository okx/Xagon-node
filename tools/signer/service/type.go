package service

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Request is the request body
type Request struct {
	OperateType    int    `json:"operateType"`
	OperateAddress string `json:"operateAddress"`
	Symbol         int    `json:"symbol"`
	ProjectSymbol  int    `json:"projectSymbol"`
	RefOrderId     string `json:"refOrderId"`
	OperateSymbol  int    `json:"operateSymbol"`
	OperateAmount  int    `json:"operateAmount"`
	SysFrom        int    `json:"sysFrom"`
	OtherInfo      string `json:"otherInfo"`
}

func (args *Request) String() string {
	b, err := json.Marshal(*args)
	if err != nil {
		return fmt.Sprintf("%+v", *args)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *args)
	}
	return out.String()
}

// SeqData is the data for sequence operation
type SeqData struct {
	Batches            []Batch        `json:"batches"`
	SignaturesAndAddrs string         `json:"signaturesAndAddrs"`
	L2Coinbase         string         `json:"l2Coinbase"`
	ContractAddress    common.Address `json:"contractAddress"`
	GasLimit           uint64         `json:"gasLimit"`
	GasPrice           string         `json:"gasPrice"`
	Nonce              uint64         `json:"nonce"`
}

// Batch is the data for batch operation
type Batch struct {
	GlobalExitRoot     string `json:"globalExitRoot"`
	MinForcedTimestamp int64  `json:"minForcedTimestamp"`
	Timestamp          int64  `json:"timestamp"`
	Transactions       string `json:"transactions"`
	TransactionsHash   string `json:"transactionsHash"`
}

// AggData is the data for aggregate operation
type AggData struct {
	NewLocalExitRoot string         `json:"newLocalExitRoot"`
	NewStateRoot     string         `json:"newStateRoot"`
	FinalNewBatch    uint64         `json:"finalNewBatch"`
	Proof            []string       `json:"proof"`
	InitNumBatch     uint64         `json:"initNumBatch"`
	PendingStateNum  int            `json:"pendingStateNum"`
	ContractAddress  common.Address `json:"contractAddress"`
	GasLimit         uint64         `json:"gasLimit"`
	GasPrice         string         `json:"gasPrice"`
	Nonce            uint64         `json:"nonce"`
}

// Response is the response body
type Response struct {
	Code      int    `json:"code"`
	Data      string `json:"data"`
	DetailMsg string `json:"detailMsg"`
	Msg       string `json:"msg"`
	Status    int    `json:"status"`
	Success   bool   `json:"success"`
}
