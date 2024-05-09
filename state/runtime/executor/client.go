package executor

import (
	"context"
	"os/exec"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// NewExecutorClient is the executor client constructor.
func NewExecutorClient(ctx context.Context, c Config) (ExecutorServiceClient, *grpc.ClientConn, context.CancelFunc) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(c.MaxGRPCMessageSize)),
	}
	const maxWaitSeconds = 20
	const maxRetries = 5
	var innerCtx context.Context
	var cancel context.CancelFunc

	connectionRetries := 0

	var executorConn *grpc.ClientConn
	var err error
	delay := 2
	for connectionRetries < maxRetries {
		innerCtx, cancel = context.WithTimeout(ctx, maxWaitSeconds*time.Second)
		executorConn, err = grpc.NewClient(c.URI, opts...)
		if err != nil {
			log.Fatalf("fail to create grpc connection to merkletree: %v", err)
		}

		log.Infof("trying to connect to executor: %v", c.URI)
		executorConn.Connect()
		err = waitForConnection(innerCtx, executorConn)
		if err != nil {
			log.Infof("Retrying connection to executor #%d", connectionRetries)
			time.Sleep(time.Duration(delay) * time.Second)
			connectionRetries = connectionRetries + 1
			out, err := exec.Command("docker", []string{"logs", "zkevm-prover"}...).Output()
			if err == nil {
				log.Infof("Prover logs:\n%s\n", out)
			}
		} else {
			log.Infof("connected to executor")
			break
		}
	}

	if connectionRetries == maxRetries {
		log.Fatalf("fail to dial: %v", err)
	}
	executorClient := NewExecutorServiceClient(executorConn)
	return executorClient, executorConn, cancel
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
