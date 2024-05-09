package merkletree

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/hashdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// NewMTDBServiceClient creates a new MTDB client.
func NewMTDBServiceClient(ctx context.Context, c Config) (hashdb.HashDBServiceClient, *grpc.ClientConn, context.CancelFunc) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	mtDBConn, err := grpc.NewClient(c.URI, opts...)
	if err != nil {
		log.Fatalf("fail to create grpc connection to merkletree: %v", err)
	}

	log.Infof("trying to connect to merkletree: %v", c.URI)
	const maxWaitSeconds = 120
	ctx, cancel := context.WithTimeout(ctx, maxWaitSeconds*time.Second)
	mtDBConn.Connect()
	err = waitForConnection(ctx, mtDBConn)
	if err != nil {
		log.Fatalf("fail to connect to merkletree: %v", err)
	}
	log.Infof("connected to merkletree")

	mtDBClient := hashdb.NewHashDBServiceClient(mtDBConn)
	return mtDBClient, mtDBConn, cancel
}

func waitForConnection(ctx context.Context, conn *grpc.ClientConn) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			s := conn.GetState()
			if s == connectivity.Ready {
				return nil
			}
		}
	}
}
