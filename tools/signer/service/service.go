package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/tools/signer/config"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

const (
	maxGas          = 5000000
	proofLen        = 24
	forkIDChunkSize = 20000
	statusSuccess   = 200
	operateTypeSeq  = 1
	operateTypeAgg  = 2
	codeSuccess     = 0
	codeFail        = 1
)

// Server is an API backend to handle RPC requests
type Server struct {
	ethCfg etherman.Config
	l1Cfg  etherman.L1Config
	ctx    context.Context

	seqPrivateKey *ecdsa.PrivateKey
	aggPrivateKey *ecdsa.PrivateKey
	ethClient     *etherman.Client

	seqAddress common.Address
	aggAddress common.Address

	result map[string]string
}

// NewServer creates a new server
func NewServer(cfg *config.Config, ctx context.Context) *Server {
	srv := &Server{
		ctx: ctx,
	}

	srv.ethCfg = etherman.Config{
		URL:              cfg.L1.RPC,
		ForkIDChunkSize:  forkIDChunkSize,
		MultiGasProvider: false,
	}

	srv.l1Cfg = etherman.L1Config{
		L1ChainID:                 cfg.L1.ChainId,
		ZkEVMAddr:                 cfg.L1.PolygonZkEVMAddress,
		MaticAddr:                 cfg.L1.PolygonMaticAddress,
		GlobalExitRootManagerAddr: cfg.L1.GlobalExitRootManagerAddr,
		DataCommitteeAddr:         cfg.L1.DataCommitteeAddr,
	}

	var err error
	srv.ethClient, err = etherman.NewClient(srv.ethCfg, srv.l1Cfg)
	if err != nil {
		log.Fatal("error creating etherman client. Error: %v", err)
	}

	_, srv.seqPrivateKey, err = srv.ethClient.LoadAuthFromKeyStoreX1(cfg.L1.SeqPrivateKey.Path, cfg.L1.SeqPrivateKey.Password)
	if err != nil {
		log.Fatal("error loading sequencer private key. Error: %v", err)
	}

	srv.seqAddress = crypto.PubkeyToAddress(srv.seqPrivateKey.PublicKey)
	log.Infof("Sequencer address: %s", srv.seqAddress.String())

	_, srv.aggPrivateKey, err = srv.ethClient.LoadAuthFromKeyStoreX1(cfg.L1.AggPrivateKey.Path, cfg.L1.AggPrivateKey.Password)
	if err != nil {
		log.Fatal("error loading aggregator private key. Error: %v", err)
	}

	srv.aggAddress = crypto.PubkeyToAddress(srv.aggPrivateKey.PublicKey)
	log.Infof("Agg address: %s", srv.aggAddress.String())

	srv.result = make(map[string]string)

	return srv
}

// Response is the response struct
func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data) //nolint:errcheck
}

// PostSignDataByOrderNo is the handler for the /priapi/v1/assetonchain/ecology/ecologyOperate endpoint
func (s *Server) PostSignDataByOrderNo(w http.ResponseWriter, r *http.Request) {
	log.Infof("PostSignDataByOrderNo start")
	response := Response{Code: codeFail, Data: "", DetailMsg: "", Msg: "", Status: statusSuccess, Success: false}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		response.DetailMsg = err.Error()
		sendJSONResponse(w, response)
		return
	}

	var requestData Request
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		response.DetailMsg = err.Error()
		sendJSONResponse(w, response)
		return
	}

	log.Infof("Request:%v", requestData.String())

	if value, ok := s.result[requestData.RefOrderId]; ok {
		response.DetailMsg = "already exist"
		log.Infof("already exist, key:%v, value:%v", requestData.RefOrderId, value)
		sendJSONResponse(w, response)
		return
	}

	if requestData.OperateType == operateTypeSeq {
		err, data := s.signSeq(requestData)
		if err != nil {
			response.DetailMsg = err.Error()
			log.Errorf("error signSeq: %v", err)
		} else {
			response.Code = codeSuccess
			response.Success = true
			s.result[requestData.RefOrderId] = data
		}
	} else if requestData.OperateType == operateTypeAgg {
		err, data := s.signAgg(requestData)
		if err != nil {
			response.DetailMsg = err.Error()
			log.Errorf("error signAgg: %v", err)
		} else {
			response.Code = codeSuccess
			response.Success = true
			s.result[requestData.RefOrderId] = data
		}
	} else {
		log.Error("error operateType")
		response.DetailMsg = "error operateType"
	}
	sendJSONResponse(w, response)
}

// signSeq is the handler for the /priapi/v1/assetonchain/ecology/ecologyOperate endpoint
func (s *Server) signSeq(requestData Request) (error, string) {
	var seqData SeqData
	err := json.Unmarshal([]byte(requestData.OtherInfo), &seqData)
	if err != nil {
		log.Errorf("error Unmarshal: %v", err)
		return err, ""
	}

	var sequences []types.Sequence
	var txHashs [][32]byte
	for _, batch := range seqData.Batches {
		var txsBytes []byte
		txsBytes, err := hex.DecodeHex(batch.Transactions)
		if err != nil {
			return err, ""
		}
		sequences = append(sequences, types.Sequence{
			BatchL2Data:          txsBytes,
			GlobalExitRoot:       common.HexToHash(batch.GlobalExitRoot),
			Timestamp:            batch.Timestamp,
			ForcedBatchTimestamp: batch.MinForcedTimestamp,
		})
		txHashs = append(txHashs, common.HexToHash(batch.TransactionsHash))
	}

	var signData []byte
	signData, err = hex.DecodeHex(seqData.SignaturesAndAddrs)
	if err != nil {
		signData = nil
	}

	_, data, err := s.ethClient.BuildMockSequenceBatchesTxData(s.seqAddress, sequences, common.HexToAddress(seqData.L2Coinbase), signData, txHashs)
	if err != nil {
		log.Errorf("error BuildSequenceBatchesTxData: %v", err)
		return err, ""
	}
	to := &seqData.ContractAddress

	// return s.getTxData(s.seqAddress, to, data)
	return s.getLegacyTxData(s.seqAddress, to, data, seqData.Nonce, seqData.GasLimit, seqData.GasPrice)
}

// signAgg is the handler for the /priapi/v1/assetonchain/ecology/ecologyOperate endpoint
func (s *Server) signAgg(requestData Request) (error, string) {
	var aggData AggData
	err := json.Unmarshal([]byte(requestData.OtherInfo), &aggData)
	if err != nil {
		log.Errorf("error Unmarshal: %v", err)
		return err, ""
	}

	newLocal, err := hex.DecodeHex(aggData.NewLocalExitRoot)
	if err != nil {
		log.Errorf("error DecodeHex: %v", err)
		return err, ""
	}

	newStateRoot, err := hex.DecodeHex(aggData.NewStateRoot)
	if err != nil {
		log.Errorf("error DecodeHex: %v", err)
		return err, ""
	}

	if len(aggData.Proof) != proofLen {
		log.Errorf("agg data len is not 24")
		return fmt.Errorf("agg proof len is not 24"), ""
	}
	proofStr := "0x"
	for _, v := range aggData.Proof {
		proofStr += v
	}

	log.Infof("proofStr: %v", proofStr)

	proof := &prover.FinalProof{
		Proof: proofStr,
	}

	var inputs = &ethmanTypes.FinalProofInputs{
		NewLocalExitRoot: newLocal,
		NewStateRoot:     newStateRoot,
		FinalProof:       proof,
	}

	_, data, err := s.ethClient.BuildTrustedVerifyBatchesTxData(aggData.InitNumBatch, aggData.FinalNewBatch, inputs)
	if err != nil {
		log.Errorf("error BuildTrustedVerifyBatchesTxData: %v", err)
		return err, ""
	}

	to := &aggData.ContractAddress
	// return s.getTxData(s.aggAddress, to, data)
	return s.getLegacyTxData(s.aggAddress, to, data, aggData.Nonce, aggData.GasLimit, aggData.GasPrice)
}

func (s *Server) getLegacyTxData(from common.Address, to *common.Address, data []byte, nonce, gasLimit uint64, gasPrice string) (error, string) {
	bigFloatGasPrice := new(big.Float)
	bigFloatGasPrice, _ = bigFloatGasPrice.SetString(gasPrice)
	result := new(big.Float).Mul(bigFloatGasPrice, new(big.Float).SetInt(big.NewInt(params.Ether)))
	gp := new(big.Int)
	result.Int(gp)

	tx := ethTypes.NewTx(&ethTypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gp,
		Gas:      gasLimit,
		To:       to,
		Data:     data,
	})

	signedTx, err := s.ethClient.SignTx(s.ctx, from, tx)
	if err != nil {
		log.Errorf("error SignTx: %v", err)
		return err, ""
	}

	txBin, err := signedTx.MarshalBinary()
	if err != nil {
		log.Errorf("error MarshalBinary: %v", err)
		return err, ""
	}

	log.Infof("TxHash: %v", signedTx.Hash().String())
	return nil, hex.EncodeToString(txBin)
}

// nolint:unused
func (s *Server) getTxData(from common.Address, to *common.Address, data []byte) (error, string) {
	nonce, err := s.ethClient.CurrentNonce(s.ctx, from)
	if err != nil {
		log.Errorf("error CurrentNonce: %v", err)
		return err, ""
	}

	tx := ethTypes.NewTx(&ethTypes.DynamicFeeTx{
		To:   to,
		Data: data,
	})
	signedTx, err := s.ethClient.SignTx(s.ctx, from, tx) //nolint:staticcheck
	if err != nil {
		log.Errorf("error SignTx: %v", err)
		return err, ""
	}

	// get gas price
	gasPrice, err := s.ethClient.SuggestedGasPrice(s.ctx)
	if err != nil {
		err := fmt.Errorf("failed to get suggested gas price: %w", err)
		log.Error(err.Error())
		return err, ""
	}

	tx = ethTypes.NewTx(&ethTypes.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: gasPrice,
		GasFeeCap: gasPrice,
		Gas:       maxGas,
		To:        to,
		Data:      data,
	})
	signedTx, err = s.ethClient.SignTx(s.ctx, from, tx)
	if err != nil {
		log.Errorf("error SignTx: %v", err)
		return err, ""
	}

	txBin, err := signedTx.MarshalBinary()
	if err != nil {
		log.Errorf("error MarshalBinary: %v", err)
		return err, ""
	}

	log.Infof("TxHash: %v", signedTx.Hash().String())
	return nil, hex.EncodeToString(txBin)
}

// GetSignDataByOrderNo is the handler for the /priapi/v1/assetonchain/ecology/ecologyOperate endpoint
func (s *Server) GetSignDataByOrderNo(w http.ResponseWriter, r *http.Request) {
	response := Response{Code: codeFail, Data: "", DetailMsg: "", Msg: "", Status: statusSuccess, Success: false}

	orderID := r.URL.Query().Get("orderId")
	projectSymbol := r.URL.Query().Get("projectSymbol")
	log.Infof("GetSignDataByOrderNo: %v,%v", orderID, projectSymbol)
	if value, ok := s.result[orderID]; ok {
		response.Code = codeSuccess
		response.Success = true
		response.Data = value
	} else {
		response.DetailMsg = "not exist"
	}

	sendJSONResponse(w, response)
}
