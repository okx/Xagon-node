package nubit

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/dataavailability/nubit/proto"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GeneralDA struct {
	client     proto.GeneralDAClient
	commitTime int64
}

func NewGeneralDA(cfg *Config) (*GeneralDA, error) {
	conn, err := grpc.NewClient(cfg.NubitRpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	da := proto.NewGeneralDAClient(conn)

	return &GeneralDA{
		client: da,
	}, nil
}

func (s *GeneralDA) Init() error {
	return nil
}

func (s *GeneralDA) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	old := atomic.LoadInt64(&s.commitTime)
	lastCommitTime := time.Since(time.Unix(old, 0))
	if lastCommitTime < NubitMinCommitTime {
		time.Sleep(NubitMinCommitTime - lastCommitTime)
	}

	// Encode NubitDA blob data
	data := EncodeSequence(batchesData)

	if atomic.CompareAndSwapInt64(&s.commitTime, old, time.Now().Unix()) {
		reply, err := s.client.Store(ctx, &proto.StoreRequest{Data: data})
		if nil != err {
			return nil, err
		}
		return reply.Receipt, nil
	}
	return nil, fmt.Errorf("another batch data is being submitted")
}

func (s *GeneralDA) GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {

	result, err := s.client.Read(ctx, &proto.Receipt{Data: dataAvailabilityMessage})
	if err != nil {
		return nil, err
	}

	batchData, hashs := DecodeSequence(result.Result)
	if len(hashs) != len(batchHashes) {
		return nil, fmt.Errorf("invalid L2 batch data retrieval arguments, %d != %d", len(hashs), len(batchHashes))
	}

	for i, bh := range batchHashes {
		if bh != hashs[i] {
			return nil, fmt.Errorf("invalid L2 batch data hash retrieval arguments, %s != %s", bh.String(), hashs[i].String())
		}
	}

	return batchData, nil
}
