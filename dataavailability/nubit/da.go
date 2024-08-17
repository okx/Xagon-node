package nubit

import (
	"context"
	"encoding/binary"

	share "github.com/RiemaLabs/nubit-node/da"

	client "github.com/RiemaLabs/nubit-node/rpc/rpc/client"
	nodeBlob "github.com/RiemaLabs/nubit-node/strucs/btx"
)

type ID []byte
type Namespace []byte
type Blob []byte

type RpcNodeClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	client *client.Client
}

func NewNodeRpc(url, token string) *RpcNodeClient {
	ctx, cancel := context.WithCancel(context.TODO())
	cn, err := client.NewClient(ctx, url, token)
	if err != nil {
		panic(err)
	}
	return &RpcNodeClient{
		ctx:    ctx,
		cancel: cancel,
		client: cn,
	}
}

func (c *RpcNodeClient) Context() context.Context {
	return c.ctx
}

func (c *RpcNodeClient) Client() *client.Client {
	return c.client
}

func (c *RpcNodeClient) Close() {
	c.client.Close()
	c.cancel()
}

// func (c *RpcNodeClient) RandSendblobFornamespace(ns namespace.Namespace) (uint64, error) {

// 	ctx, cancle := context.WithTimeout(c.ctx, 40*time.Second)
// 	defer cancle()

// 	nsp, err := share.NamespaceFromBytes(ns.Bytes())
// 	if err != nil {
// 		panic(err)
// 	}

// 	body, err := nodeBlob.NewBlobV0(nsp, bytes.Repeat([]byte{0x01}, 500*1024))
// 	if err != nil {
// 		return 0, err
// 	}

// 	return c.client.Blob.Submit(ctx, []*nodeBlob.Blob{body}, 0.01)
// }

// func (c *RpcNodeClient) NodeInfo() string {
// 	info, err := c.client.Node.Info(c.ctx)
// 	if err != nil {
// 		return ""
// 	}
// 	return fmt.Sprintf("type:%s, version:%s", info.Type, info.APIVersion)
// }

func (c *RpcNodeClient) GetProofs(ctx context.Context, ids []ID, namespace Namespace) ([]*nodeBlob.Proof, error) {

	pss := make([]*nodeBlob.Proof, len(ids))
	for i, id := range ids {
		height := binary.BigEndian.Uint64(id[:8])
		np, err := share.NamespaceFromBytes(namespace)
		if err != nil {
			return nil, err
		}
		p, err := c.client.Blob.GetProof(ctx, height, np, nodeBlob.KzgCommitment(id[8:]))
		if err != nil {
			return nil, err
		}
		pss[i] = p
	}
	return pss, nil
}

func (c *RpcNodeClient) Submit(ctx context.Context, blobs []Blob, gasPrice float64, namespace Namespace) ([]ID, error) {
	nsp, err := share.NamespaceFromBytes(namespace)
	if err != nil {
		panic(err)
	}

	bcs := make([]*nodeBlob.Blob, 0, len(blobs))
	for _, b := range blobs {
		blob, err := nodeBlob.NewBlobV0(nsp, b)
		if err != nil {
			return nil, err
		}
		bcs = append(bcs, blob)
	}
	height, err := c.client.Blob.Submit(ctx, bcs, nodeBlob.GasPrice(gasPrice))
	_ = height
	// TODO
	return []ID{}, err
}

/*
	// MaxBlobSize returns the max blob size
	MaxBlobSize(ctx context.Context) (uint64, error)

	// Get returns Blob for each given ID, or an error.
	//
	// Error should be returned if ID is not formatted properly, there is no Blob for given ID or any other client-level
	// error occurred (dropped connection, timeout, etc).
	Get(ctx context.Context, ids []ID, namespace Namespace) ([]Blob, error)

	// GetIDs returns IDs of all Blobs located in DA at given height.
	GetIDs(ctx context.Context, height uint64, namespace Namespace) ([]ID, error)

	// GetProofs returns inclusion Proofs for Blobs specified by their IDs.
	GetProofs(ctx context.Context, ids []ID, namespace Namespace) ([]Proof, error)

	// Commit creates a Commitment for each given Blob.
	Commit(ctx context.Context, blobs []Blob, namespace Namespace) ([]Commitment, error)

	// Submit submits the Blobs to Data Availability layer.
	//
	// This method is synchronous. Upon successful submission to Data Availability layer, it returns the IDs identifying blobs
	// in DA.
	Submit(ctx context.Context, blobs []Blob, gasPrice float64, namespace Namespace) ([]ID, error)

	// Validate validates Commitments against the corresponding Proofs. This should be possible without retrieving the Blobs.
	Validate(ctx context.Context, ids []ID, proofs []Proof, namespace Namespace) ([]bool, error)
*/
