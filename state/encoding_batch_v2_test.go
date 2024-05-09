package state

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/stretchr/testify/require"
)

const (
	// changeL2Block + deltaTimeStamp + indexL1InfoTree
	codedL2BlockHeader = "0b73e6af6f00000000"
	// 2 x [ tx coded in RLP + r,s,v,efficiencyPercentage]
	codedRLP2Txs1 = "ee02843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e88080bff0e780ba7db409339fd3f71969fa2cbf1b8535f6c725a1499d3318d3ef9c2b6340ddfab84add2c188f9efddb99771db1fe621c981846394ea4f035c85bcdd51bffee03843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880805b346aa02230b22e62f73608de9ff39a162a6c24be9822209c770e3685b92d0756d5316ef954eefc58b068231ccea001fb7ac763ebe03afd009ad71cab36861e1bff"
	// 2 x [ tx coded in RLP + r,s,v,efficiencyPercentage]
	codedRLP2Txs2 = "ee80843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880801cee7e01dc62f69a12c3510c6d64de04ee6346d84b6a017f3e786c7d87f963e75d8cc91fa983cd6d9cf55fff80d73bd26cd333b0f098acc1e58edb1fd484ad731bffee01843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880803ee20a0764440b016c4a2ee4e7e4eb3a5a97f1e6a6c9f40bf5ecf50f95ff636d63878ddb3e997e519826c7bb26fb7c5950a208e1ec722a9f1c568c4e479b40341cff"
	codedL2Block1 = codedL2BlockHeader + codedRLP2Txs1
	codedL2Block2 = codedL2BlockHeader + codedRLP2Txs2
	// Batch 420.000 (Incaberry) from Testnet
	realBatchIncaberry      = "ee8307c4848402faf08082520894417a7ba2d8d0060ae6c54fd098590db854b9c1d58609184e72a000808205a28080e8c76f8b8ec579362a4ef92dc1c8c372ad4ef6372a20903b3997408743e86239394ad6decc3bc080960b6c62ad78bc09913cba88fd98d595457b3462ed1494b91cffee8307c4858402faf08082520894417a7ba2d8d0060ae6c54fd098590db854b9c1d58609184e72a000808205a28080ed0de9758ff75ae777821e45178da0163c719341188220050cc4ad33048cd9cb272951662ae72269cf611528d591fcf682c8bad4402d98dbac4abc1b2be1ca431cffee8307c4868402faf08082520894417a7ba2d8d0060ae6c54fd098590db854b9c1d58609184e72a000808205a280807c94882ecf48d65b6240e7355c32e7d1a56366fd9571471cb664463ad2afecdd564d24abbea5b38b74dda029cdac3109f199f5e3e683acfbe43e7f27fe23b60b1cffee8307c4878402faf08082520894417a7ba2d8d0060ae6c54fd098590db854b9c1d58609184e72a000808205a280801b5e85cc1b402403a625610d4319558632cffd2b14a15bc031b9ba644ecc48a332bcc608e894b9ede61220767558e1d9e02780b53dbdd9bcc01de0ab2b1742951bffee8307c4888402faf08082520894417a7ba2d8d0060ae6c54fd098590db854b9c1d58609184e72a000808205a2808089eee14afeead54c815953a328ec52d441128e71d08ff75b4e5cd23db6fa67e774ca24e8878368eee5ad4562340edebcfb595395d40f8a5b0301e19ced92af5f1cffee8307c4898402faf08082520894417a7ba2d8d0060ae6c54fd098590db854b9c1d58609184e72a000808205a280807b672107c41caf91cff9061241686dd37e8d1e013d81f7f383b76afa93b7ff85413d4fc4c7e9613340b8fc29aefd0c42a3db6d75340b1bec0b895d324bcfa02e1cffee8307c48a8402faf08082520894417a7ba2d8d0060ae6c54fd098590db854b9c1d58609184e72a000808205a28080efadeca94da405cf44881670bc8b2464d006af41f20517e82339c72d73543c5c4e1e546eea07b4b751e3e2f909bd4026f742684c923bf666985f9a5a1cd91cde1bffee8307c48b8402faf08082520894417a7ba2d8d0060ae6c54fd098590db854b9c1d58609184e72a000808205a2808092ac34e2d6a38c7df5df96c78f9d837daaa7f74352d8c42fe671ef8ba6565ae350648c7e736a0017bf90370e766720c410441f6506765c70fad91ce046c1fad61bfff86c8206838402faf08082803194828f7ceca102de66a6ed4f4b6abee0bd1bd4f9dc80b844095ea7b3000000000000000000000000e907ec70b4efbb28efbf6f4ffb3ae0d34012eaa00000000000000000000000000000000000000000000000011a8297a4dca080008205a28080579cfefee3fa664c8b59190de80454da9642b7647a46b929c9fcc89105b2d5575d28665bef2bb1052db0d36ec1e92bc7503efaa74798fe3630b8867318c20d4e1cff"
	realBatchConvertedEtrog = codedL2BlockHeader + realBatchIncaberry
)

func TestDecodeEmptyBatchV2(t *testing.T) {
	batchL2Data, err := hex.DecodeString("")
	require.NoError(t, err)

	batch, err := DecodeBatchV2(batchL2Data)
	require.NoError(t, err)
	require.Equal(t, 0, len(batch.Blocks))
}

func TestDecodeBatches(t *testing.T) {
	type testCase struct {
		name          string
		batchL2Data   string
		expectedError error
	}
	testCases := []testCase{
		{
			name:          "batch dont start with 0x0b (changeL2Block)",
			batchL2Data:   "0c",
			expectedError: ErrInvalidBatchV2,
		},
		{
			name:          "batch no enough  data to decode block.deltaTimestamp",
			batchL2Data:   "0b010203",
			expectedError: ErrInvalidBatchV2,
		},
		{
			name:          "batch no enough  data to decode block.index",
			batchL2Data:   "0b01020304010203",
			expectedError: ErrInvalidBatchV2,
		},
		{
			name:          "batch no enough  data to decode block.index",
			batchL2Data:   "0b01020304010203",
			expectedError: ErrInvalidBatchV2,
		},
		{
			name:          "valid batch no trx, just L2Block",
			batchL2Data:   "0b0102030401020304",
			expectedError: nil,
		},
		{
			name:          "invalid batch bad RLP codification",
			batchL2Data:   "0b" + "01020304" + "01020304" + "7f",
			expectedError: ErrInvalidRLP,
		},
		{
			name:          "1 block + 2 txs",
			batchL2Data:   "0b" + "73e6af6f" + "00000000" + codedRLP2Txs1 + codedRLP2Txs2,
			expectedError: nil,
		},
		{
			name:          "1 block + 1 txs",
			batchL2Data:   "0b" + "73e6af6f" + "00000000" + codedRLP2Txs1,
			expectedError: nil,
		},
		{
			name:          "1 block + 1 txs, missiging efficiencyPercentage",
			batchL2Data:   "0b" + "73e6af6f" + "00000000" + codedRLP2Txs1[0:len(codedRLP2Txs1)-2],
			expectedError: ErrInvalidBatchV2,
		},
		{
			name:          "real batch converted to etrog",
			batchL2Data:   realBatchConvertedEtrog,
			expectedError: nil,
		},
		{
			name:          "pass a V1 batch(incaberry) must fail",
			batchL2Data:   realBatchIncaberry,
			expectedError: ErrBatchV2DontStartWithChangeL2Block,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Debug("************************ ", tc.name, " ************************")
			data, err := hex.DecodeString(tc.batchL2Data)
			require.NoError(t, err)
			_, err = DecodeBatchV2(data)
			if err != nil {
				log.Debugf("[%s] %v", tc.name, err)
			}
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDecodeBatchV2(t *testing.T) {
	batchL2Data, err := hex.DecodeString(codedL2Block1)
	require.NoError(t, err)
	batchL2Data2, err := hex.DecodeString(codedL2Block2)
	require.NoError(t, err)
	batch := append(batchL2Data, batchL2Data2...)
	decodedBatch, err := DecodeBatchV2(batch)
	require.NoError(t, err)
	require.Equal(t, 2, len(decodedBatch.Blocks))
	require.Equal(t, uint32(0x73e6af6f), decodedBatch.Blocks[0].DeltaTimestamp)
	require.Equal(t, uint32(0x00000000), decodedBatch.Blocks[0].IndexL1InfoTree)
}

func TestDecodeRLPLength(t *testing.T) {
	type testCase struct {
		name           string
		data           string
		expectedError  error
		expectedResult uint64
	}
	testCases := []testCase{
		{
			name:          "must start >= 0xc0",
			data:          "bf",
			expectedError: ErrInvalidRLP,
		},
		{
			name:           "shortRLP: c0 -> len=0",
			data:           "c0",
			expectedResult: 1,
		},
		{
			name:           "shortRLP: c1 -> len=1",
			data:           "c1",
			expectedResult: 2, // 1 byte header + 1 byte of data
		},
		{
			name:           "shortRLP: byte>0xf7",
			data:           "f7",
			expectedResult: 56, // 1 byte header + 55 bytes of data
		},
		{
			name:          "longRLP: f8: 1 extra byte, missing data",
			data:          "f8",
			expectedError: ErrInvalidRLP,
		},
		{
			name:           "longRLP: f8:size is stored in next byte ->0x01 (code add the length of bytes of the size??)",
			data:           "f8" + "01",
			expectedResult: 3, // 2 bytes of header + 1 byte of data
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.Debug("************************ ", tc.name, " ************************")
			data, err := hex.DecodeString(tc.data)
			require.NoError(t, err)
			length, err := decodeRLPListLengthFromOffset(data, 0)
			if err != nil {
				log.Debugf("[%s] %v", tc.name, err)
			}
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, length)
			}
		})
	}
}

func TestEncodeBatchV2(t *testing.T) {
	block1 := L2BlockRaw{
		ChangeL2BlockHeader: ChangeL2BlockHeader{
			DeltaTimestamp:  123,
			IndexL1InfoTree: 456,
		},
		Transactions: []L2TxRaw{},
	}

	block2 := L2BlockRaw{
		ChangeL2BlockHeader: ChangeL2BlockHeader{
			DeltaTimestamp:  789,
			IndexL1InfoTree: 101112,
		},
		Transactions: []L2TxRaw{},
	}
	blocks := []L2BlockRaw{block1, block2}

	expectedBatchData := []byte{
		0xb, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x1, 0xc8, 0xb, 0x0, 0x0, 0x3, 0x15, 0x0, 0x1, 0x8a, 0xf8,
	}

	batchData, err := EncodeBatchV2(&BatchRawV2{Blocks: blocks})
	require.NoError(t, err)
	require.Equal(t, expectedBatchData, batchData)
}

func TestDecodeEncodeBatchV2(t *testing.T) {
	batchL2Data, err := hex.DecodeString(codedL2Block1 + codedL2Block2)
	require.NoError(t, err)
	decodedBatch, err := DecodeBatchV2(batchL2Data)
	require.NoError(t, err)
	require.Equal(t, 2, len(decodedBatch.Blocks))
	encoded, err := EncodeBatchV2(decodedBatch)
	require.NoError(t, err)
	require.Equal(t, batchL2Data, encoded)
}

func TestEncodeEmptyBatchV2Fails(t *testing.T) {
	l2Batch := BatchRawV2{}
	_, err := EncodeBatchV2(&l2Batch)
	require.ErrorIs(t, err, ErrInvalidBatchV2)
	_, err = EncodeBatchV2(nil)
	require.ErrorIs(t, err, ErrInvalidBatchV2)
}

func TestDecodeForcedBatchV2(t *testing.T) {
	batchL2Data, err := hex.DecodeString(codedRLP2Txs1)
	require.NoError(t, err)
	decodedBatch, err := DecodeForcedBatchV2(batchL2Data)
	require.NoError(t, err)
	require.Equal(t, 2, len(decodedBatch.Transactions))
}

func TestDecodeForcedBatchV2WithRegularBatch(t *testing.T) {
	batchL2Data, err := hex.DecodeString(codedL2Block1)
	require.NoError(t, err)
	_, err = DecodeForcedBatchV2(batchL2Data)
	require.Error(t, err)
}

func TestEncodeBatchV2WithTxInBinary(t *testing.T) {
	block1 := L2BlockRaw{
		ChangeL2BlockHeader: ChangeL2BlockHeader{
			DeltaTimestamp:  123,
			IndexL1InfoTree: 456,
		},
		Transactions: []L2TxRaw{
			{
				EfficiencyPercentage: 255,
				TxAlreadyEncoded:     true,
				Data:                 []byte{0x01, 0x02, 0x03},
			},
		},
	}

	block2 := L2BlockRaw{
		ChangeL2BlockHeader: ChangeL2BlockHeader{
			DeltaTimestamp:  789,
			IndexL1InfoTree: 101112,
		},
		Transactions: []L2TxRaw{},
	}
	blocks := []L2BlockRaw{block1, block2}

	expectedBatchData := []byte{
		0xb, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x1, 0xc8, 0x1, 0x2, 0x3, 0xff, 0xb, 0x0, 0x0, 0x3, 0x15, 0x0, 0x1, 0x8a, 0xf8,
	}

	batchData, err := EncodeBatchV2(&BatchRawV2{Blocks: blocks})
	require.NoError(t, err)
	require.Equal(t, expectedBatchData, batchData)
}
