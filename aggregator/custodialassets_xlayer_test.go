package aggregator

import (
	"context"
	"math/big"
	"testing"
	"time"

	agglayerTypes "github.com/0xPolygon/agglayer/rpc/types"
	"github.com/0xPolygon/agglayer/tx"
	zktypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

const (
	domain                 = "http://asset-onchain.base-defi.svc.test.local:7001"
	seqAddr                = "1a13bddcc02d363366e04d4aa588d3c125b0ff6f"
	aggAddr                = "66e39a1e507af777e8c385e2d91559e20e306303"
	contractAddr           = "8947dc90862f386968966b22dfe5edf96435bc2f"
	contractAddrAgg        = "1d5298ee11f7cd56fb842b7894346bfb2e47a95f"
	l1ChainID       uint64 = 11155111
	AccessKey              = ""
	SecretKey              = ""
	// domain       = "http://127.0.0.1:7001"
	// seqAddr      = "f39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	// aggAddr      = "70997970C51812dc3A010C7d01b50e0d17dc79C8"
	// contractAddr = "812cB73e48841a6736bB94c65c56341817cE6304"
)

func TestClientPostSignRequestAndWaitResultTxSigner(t *testing.T) {
	newStateRoot := common.BigToHash(big.NewInt(100))
	t.Log(newStateRoot.Hex())
	newLocalExitRoot := common.BigToHash(big.NewInt(1000))
	t.Log(newLocalExitRoot.Hex())
	proof := agglayerTypes.ArgBytes([]byte("sampleProof"))

	t.Log(proof.Hex())

	tnx := tx.Tx{
		LastVerifiedBatch: agglayerTypes.ArgUint64(36),
		NewVerifiedBatch:  *agglayerTypes.ArgUint64Ptr(72),
		ZKP: tx.ZKP{
			NewStateRoot:     newStateRoot,
			NewLocalExitRoot: newLocalExitRoot,
			Proof:            proof,
		},
		RollupID: 2,
	}
	t.Log(tnx.Hash())

	agg := &Aggregator{
		cfg: Config{
			CustodialAssets: CustodialAssetsConfig{
				Enable:            false,
				URL:               domain,
				Symbol:            2882,
				SequencerAddr:     common.HexToAddress(seqAddr),
				AggregatorAddr:    common.HexToAddress(aggAddr),
				WaitResultTimeout: zktypes.NewDuration(4 * time.Minute),
				OperateTypeSeq:    5,
				OperateTypeAgg:    6,
				ProjectSymbol:     3011,
				OperateSymbol:     2,
				SysFrom:           3,
				UserID:            0,
				OperateAmount:     0,
				RequestSignURI:    "/priapi/v1/assetonchain/ecology/ecologyOperate",
				QuerySignURI:      "/priapi/v1/assetonchain/ecology/querySignDataByOrderNo",
				AccessKey:         AccessKey,
				SecretKey:         SecretKey,
			},
		},
	}
	ctx := context.WithValue(context.Background(), traceID, uuid.New().String())

	myTx, err := agg.signTx(ctx, tnx)
	t.Log(myTx.Tx)
	t.Log(myTx.Signature)
	t.Log(err)
}

func TestTxHash(t *testing.T) {
	newStateRoot := common.BigToHash(big.NewInt(10000000))
	t.Log(newStateRoot.Hex())
	newLocalExitRoot := common.BigToHash(big.NewInt(500000000))
	t.Log(newLocalExitRoot.Hex())
	proof := agglayerTypes.ArgBytes([]byte("sampleProof"))

	t.Log(proof.Hex())

	tnx := tx.Tx{
		LastVerifiedBatch: agglayerTypes.ArgUint64(20000000),
		NewVerifiedBatch:  *agglayerTypes.ArgUint64Ptr(300000000),
		ZKP: tx.ZKP{
			NewStateRoot:     newStateRoot,
			NewLocalExitRoot: newLocalExitRoot,
			Proof:            proof,
		},
		RollupID: 2,
	}
	t.Log(tnx.Hash())
}
