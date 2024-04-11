package state

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

const (
	double       = 2
	ether155V    = 27
	etherPre155V = 35
	// Decoding constants
	headerByteLength uint64 = 1
	sLength          uint64 = 32
	rLength          uint64 = 32
	vLength          uint64 = 1
	c0               uint64 = 192 // 192 is c0. This value is defined by the rlp protocol
	ff               uint64 = 255 // max value of rlp header
	shortRlp         uint64 = 55  // length of the short rlp codification
	f7               uint64 = 247 // 192 + 55 = c0 + shortRlp

	// EfficiencyPercentageByteLength is the length of the effective percentage in bytes
	EfficiencyPercentageByteLength uint64 = 1
)

// EncodeTransactions RLP encodes the given transactions
func EncodeTransactions(txs []types.Transaction, effectivePercentages []uint8, forkID uint64) ([]byte, error) {
	var batchL2Data []byte

	for i, tx := range txs {
		txData, err := prepareRPLTxData(tx)
		if err != nil {
			return nil, err
		}
		batchL2Data = append(batchL2Data, txData...)

		if forkID >= FORKID_DRAGONFRUIT {
			effectivePercentageAsHex, err := hex.DecodeHex(fmt.Sprintf("%x", effectivePercentages[i]))
			if err != nil {
				return nil, err
			}
			batchL2Data = append(batchL2Data, effectivePercentageAsHex...)
		}
	}

	return batchL2Data, nil
}

func prepareRPLTxData(tx types.Transaction) ([]byte, error) {
	v, r, s := tx.RawSignatureValues()
	sign := 1 - (v.Uint64() & 1)

	nonce, gasPrice, gas, to, value, data, chainID := tx.Nonce(), tx.GasPrice(), tx.Gas(), tx.To(), tx.Value(), tx.Data(), tx.ChainId()

	rlpFieldsToEncode := []interface{}{
		nonce,
		gasPrice,
		gas,
		to,
		value,
		data,
	}

	if !IsPreEIP155Tx(tx) {
		rlpFieldsToEncode = append(rlpFieldsToEncode, chainID)
		rlpFieldsToEncode = append(rlpFieldsToEncode, uint(0))
		rlpFieldsToEncode = append(rlpFieldsToEncode, uint(0))
	}

	txCodedRlp, err := rlp.EncodeToBytes(rlpFieldsToEncode)
	if err != nil {
		return nil, err
	}

	newV := new(big.Int).Add(big.NewInt(ether155V), big.NewInt(int64(sign)))
	newRPadded := fmt.Sprintf("%064s", r.Text(hex.Base))
	newSPadded := fmt.Sprintf("%064s", s.Text(hex.Base))
	newVPadded := fmt.Sprintf("%02s", newV.Text(hex.Base))
	txData, err := hex.DecodeString(hex.EncodeToString(txCodedRlp) + newRPadded + newSPadded + newVPadded)
	if err != nil {
		return nil, err
	}
	return txData, nil
}

// EncodeTransactionsWithoutEffectivePercentage RLP encodes the given transactions without the effective percentage
func EncodeTransactionsWithoutEffectivePercentage(txs []types.Transaction) ([]byte, error) {
	var batchL2Data []byte

	for _, tx := range txs {
		txData, err := prepareRPLTxData(tx)
		if err != nil {
			return nil, err
		}
		batchL2Data = append(batchL2Data, txData...)
	}

	return batchL2Data, nil
}

// EncodeTransactionWithoutEffectivePercentage RLP encodes the given transaction without the effective percentage
func EncodeTransactionWithoutEffectivePercentage(tx types.Transaction) ([]byte, error) {
	return EncodeTransactionsWithoutEffectivePercentage([]types.Transaction{tx})
}

// EncodeTransaction RLP encodes the given transaction
func EncodeTransaction(tx types.Transaction, effectivePercentage uint8, forkID uint64) ([]byte, error) {
	return EncodeTransactions([]types.Transaction{tx}, []uint8{effectivePercentage}, forkID)
}

// EncodeUnsignedTransaction RLP encodes the given unsigned transaction
func EncodeUnsignedTransaction(tx types.Transaction, chainID uint64, forcedNonce *uint64, forkID uint64) ([]byte, error) {
	v, _ := new(big.Int).SetString("0x1c", 0)
	r, _ := new(big.Int).SetString("0xa54492cfacf71aef702421b7fbc70636537a7b2fbe5718c5ed970a001bb7756b", 0)
	s, _ := new(big.Int).SetString("0x2e9fb27acc75955b898f0b12ec52aa34bf08f01db654374484b80bf12f0d841e", 0)

	sign := 1 - (v.Uint64() & 1)

	nonce, gasPrice, gas, to, value, data, chainID := tx.Nonce(), tx.GasPrice(), tx.Gas(), tx.To(), tx.Value(), tx.Data(), chainID //nolint:gomnd
	log.Debug(nonce, " ", gasPrice, " ", gas, " ", to, " ", value, " ", len(data), " ", chainID)

	if forcedNonce != nil {
		nonce = *forcedNonce
		log.Debug("Forced nonce: ", nonce)
	}

	txCodedRlp, err := rlp.EncodeToBytes([]interface{}{
		nonce,
		gasPrice,
		gas,
		to,
		value,
		data,
		big.NewInt(0).SetUint64(chainID), uint(0), uint(0),
	})

	if err != nil {
		return nil, err
	}

	newV := new(big.Int).Add(big.NewInt(ether155V), big.NewInt(int64(sign)))
	newRPadded := fmt.Sprintf("%064s", r.Text(hex.Base))
	newSPadded := fmt.Sprintf("%064s", s.Text(hex.Base))
	newVPadded := fmt.Sprintf("%02s", newV.Text(hex.Base))
	effectivePercentageAsHex := fmt.Sprintf("%x", MaxEffectivePercentage)
	// Only add EffectiveGasprice if forkID is equal or higher than DRAGONFRUIT_FORKID
	if forkID < FORKID_DRAGONFRUIT {
		effectivePercentageAsHex = ""
	}
	txData, err := hex.DecodeString(hex.EncodeToString(txCodedRlp) + newRPadded + newSPadded + newVPadded + effectivePercentageAsHex)
	if err != nil {
		return nil, err
	}

	return txData, nil
}

// DecodeTxs extracts Transactions for its encoded form
func DecodeTxs(txsData []byte, forkID uint64) ([]types.Transaction, []byte, []uint8, error) {
	// Process coded txs
	var pos uint64
	var txs []types.Transaction
	var efficiencyPercentages []uint8
	txDataLength := uint64(len(txsData))
	if txDataLength == 0 {
		return txs, txsData, nil, nil
	}
	for pos < txDataLength {
		num, err := strconv.ParseUint(hex.EncodeToString(txsData[pos:pos+1]), hex.Base, hex.BitSize64)
		if err != nil {
			log.Debug("error parsing header length: ", err)
			return []types.Transaction{}, txsData, []uint8{}, err
		}
		// First byte is the length and must be ignored
		if num < c0 {
			log.Debugf("error num < c0 : %d, %d", num, c0)
			return []types.Transaction{}, txsData, []uint8{}, ErrInvalidData
		}
		length := num - c0
		if length > shortRlp { // If rlp is bigger than length 55
			// n is the length of the rlp data without the header (1 byte) for example "0xf7"
			if (pos + 1 + num - f7) > txDataLength {
				log.Debug("error parsing length: ", err)
				return []types.Transaction{}, txsData, []uint8{}, err
			}
			n, err := strconv.ParseUint(hex.EncodeToString(txsData[pos+1:pos+1+num-f7]), hex.Base, hex.BitSize64) // +1 is the header. For example 0xf7
			if err != nil {
				log.Debug("error parsing length: ", err)
				return []types.Transaction{}, txsData, []uint8{}, err
			}
			if n+num < f7 {
				log.Debug("error n + num < f7: ", err)
				return []types.Transaction{}, txsData, []uint8{}, ErrInvalidData
			}
			length = n + num - f7 // num - f7 is the header. For example 0xf7
		}

		endPos := pos + length + rLength + sLength + vLength + headerByteLength

		if forkID >= FORKID_DRAGONFRUIT {
			endPos += EfficiencyPercentageByteLength
		}

		if endPos > txDataLength {
			err := fmt.Errorf("endPos %d is bigger than txDataLength %d", endPos, txDataLength)
			log.Debug("error parsing header: ", err)
			return []types.Transaction{}, txsData, []uint8{}, ErrInvalidData
		}

		if endPos < pos {
			err := fmt.Errorf("endPos %d is smaller than pos %d", endPos, pos)
			log.Debug("error parsing header: ", err)
			return []types.Transaction{}, txsData, []uint8{}, ErrInvalidData
		}

		if endPos < pos {
			err := fmt.Errorf("endPos %d is smaller than pos %d", endPos, pos)
			log.Debug("error parsing header: ", err)
			return []types.Transaction{}, txsData, []uint8{}, ErrInvalidData
		}

		fullDataTx := txsData[pos:endPos]
		dataStart := pos + length + headerByteLength
		txInfo := txsData[pos:dataStart]
		rData := txsData[dataStart : dataStart+rLength]
		sData := txsData[dataStart+rLength : dataStart+rLength+sLength]
		vData := txsData[dataStart+rLength+sLength : dataStart+rLength+sLength+vLength]

		if forkID >= FORKID_DRAGONFRUIT {
			efficiencyPercentage := txsData[dataStart+rLength+sLength+vLength : endPos]
			efficiencyPercentages = append(efficiencyPercentages, efficiencyPercentage[0])
		}

		pos = endPos

		// Decode rlpFields
		var rlpFields [][]byte
		err = rlp.DecodeBytes(txInfo, &rlpFields)
		if err != nil {
			log.Error("error decoding tx Bytes: ", err, ". fullDataTx: ", hex.EncodeToString(fullDataTx), "\n tx: ", hex.EncodeToString(txInfo), "\n Txs received: ", hex.EncodeToString(txsData))
			return []types.Transaction{}, txsData, []uint8{}, ErrInvalidData
		}

		legacyTx, err := RlpFieldsToLegacyTx(rlpFields, vData, rData, sData)
		if err != nil {
			log.Debug("error creating tx from rlp fields: ", err, ". fullDataTx: ", hex.EncodeToString(fullDataTx), "\n tx: ", hex.EncodeToString(txInfo), "\n Txs received: ", hex.EncodeToString(txsData))
			return []types.Transaction{}, txsData, []uint8{}, err
		}

		tx := types.NewTx(legacyTx)
		txs = append(txs, *tx)
	}
	return txs, txsData, efficiencyPercentages, nil
}

// DecodeTx decodes a string rlp tx representation into a types.Transaction instance
func DecodeTx(encodedTx string) (*types.Transaction, error) {
	b, err := hex.DecodeHex(encodedTx)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return tx, nil
}

// GenerateReceipt generates a receipt from a processed transaction
func GenerateReceipt(blockNumber *big.Int, processedTx *ProcessTransactionResponse, txIndex uint, forkID uint64) *types.Receipt {
	receipt := &types.Receipt{
		Type:             uint8(processedTx.Type),
		BlockNumber:      blockNumber,
		GasUsed:          processedTx.GasUsed,
		TxHash:           processedTx.Tx.Hash(),
		TransactionIndex: txIndex,
		ContractAddress:  processedTx.CreateAddress,
		Logs:             processedTx.Logs,
	}
	if forkID <= FORKID_ETROG {
		receipt.PostState = processedTx.StateRoot.Bytes()
		receipt.CumulativeGasUsed = processedTx.GasUsed
	} else {
		receipt.PostState = []byte{}
		receipt.CumulativeGasUsed = processedTx.CumulativeGasUsed
	}
	if processedTx.EffectiveGasPrice != "" {
		effectiveGasPrice, ok := big.NewInt(0).SetString(processedTx.EffectiveGasPrice, 0)
		if !ok {
			log.Errorf("error converting effective gas price %s to big.Int", processedTx.EffectiveGasPrice)
		}
		receipt.EffectiveGasPrice = effectiveGasPrice
	}

	// TODO: this fix is temporary while the Executor is returning a
	// different Tx hash for the TxHash, Log.TxHash and Tx.Hash().
	// At the moment, the processedTx.TxHash and Log[n].TxHash are
	// returning a different hash than the Hash of the transaction
	// sent to be processed by the Executor.
	// The processedTx.Tx.Hash() is correct.
	for i := 0; i < len(receipt.Logs); i++ {
		receipt.Logs[i].TxHash = processedTx.Tx.Hash()
	}
	if forkID <= FORKID_ETROG {
		if processedTx.RomError == nil {
			receipt.Status = types.ReceiptStatusSuccessful
		} else {
			receipt.Status = types.ReceiptStatusFailed
		}
	} else {
		receipt.Status = uint64(processedTx.Status)
	}

	return receipt
}

// IsPreEIP155Tx checks if the tx is a tx that has a chainID as zero and
// V field is either 27 or 28
func IsPreEIP155Tx(tx types.Transaction) bool {
	v, _, _ := tx.RawSignatureValues()
	return tx.ChainId().Uint64() == 0 && (v.Uint64() == 27 || v.Uint64() == 28)
}

// CheckLogOrder checks the order of the logs. The order should be incremental
func CheckLogOrder(logs []*types.Log) bool {
	logsAux := make([]*types.Log, len(logs))
	copy(logsAux, logs)
	sort.Slice(logsAux, func(i, j int) bool {
		return logsAux[i].Index < logsAux[j].Index
	})
	if len(logs) != len(logsAux) {
		return false
	}
	for i := range logs {
		if logsAux[i].Index != logs[i].Index {
			log.Debug("Array index: ", i, ". Index of log on each array: ", logsAux[i].Index, logs[i].Index)
			return false
		}
	}
	return true
}

// Ptr returns a pointer for any instance
func Ptr[T any](v T) *T {
	return &v
}

// HashByteArray returns the hash of the given byte array
func HashByteArray(data []byte) common.Hash {
	sha := sha3.NewLegacyKeccak256()
	sha.Write(data)
	return common.BytesToHash(sha.Sum(nil))
}
