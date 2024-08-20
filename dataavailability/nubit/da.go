package nubit

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	share "github.com/RiemaLabs/nubit-node/da"
	"github.com/rollkit/go-da"

	client "github.com/RiemaLabs/nubit-node/rpc/rpc/client"
	nodeBlob "github.com/RiemaLabs/nubit-node/strucs/btx"
)

var ErrNamespaceMismatch = fmt.Errorf("namespace mismatch")

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

func (c *RpcNodeClient) GetProofs(ctx context.Context, ids []da.ID, namespace da.Namespace) ([]da.Proof, error) {

	pss := make([]da.Proof, len(ids))
	for i, id := range ids {
		hbyte, commitment := ParseID(id)
		height := binary.LittleEndian.Uint64(hbyte)
		np, err := share.NamespaceFromBytes(namespace)
		if err != nil {
			return nil, err
		}
		p, err := c.client.Blob.GetProof(ctx, height, np, nodeBlob.KzgCommitment(commitment))
		if err != nil {
			return nil, err
		}

		js, err := json.Marshal(p)
		if nil != err {
			return nil, err
		}

		pss[i] = js
	}
	return pss, nil
}

func (c *RpcNodeClient) Submit(ctx context.Context, blobs []da.Blob, gasPrice float64, namespace da.Namespace) ([]da.ID, error) {
	nsp, err := share.NamespaceFromBytes(namespace)
	if err != nil {
		panic(err)
	}

	bcs := make([]*nodeBlob.Blob, 0, len(blobs))
	commitments := make([]nodeBlob.KzgCommitment, 0, len(blobs))
	for _, b := range blobs {
		blob, err := nodeBlob.NewBlobV0(nsp, b)
		if err != nil {
			return nil, err
		}
		bcs = append(bcs, blob)
		commitments = append(commitments, blob.Commitment)
	}
	height, err := c.client.Blob.Submit(ctx, bcs, nodeBlob.GasPrice(gasPrice))
	if nil != err {
		return nil, err
	}

	hbytes := [8]byte{}
	binary.LittleEndian.PutUint64(hbytes[:], height)

	IDs := make([]da.ID, len(commitments))
	for i, c := range commitments {
		IDs[i] = MakeID(hbytes[:], c)
	}
	return IDs, err
}

// MaxBlobSize returns the max blob size
// TODO Temporarily set to 1M
func (c *RpcNodeClient) MaxBlobSize(ctx context.Context) (uint64, error) {
	return 1024 * 1024, nil
}

// Get returns Blob for each given ID, or an error.
//
// Error should be returned if ID is not formatted properly, there is no Blob for given ID or any other client-level
// error occurred (dropped connection, timeout, etc).
func (c *RpcNodeClient) Get(ctx context.Context, ids []da.ID, namespace da.Namespace) ([]da.Blob, error) {
	nsp, err := share.NamespaceFromBytes(namespace)
	if err != nil {
		return nil, err
	}
	blobs := make([]da.Blob, len(ids))
	for i, id := range ids {
		hbyte, commitment := ParseID(id)
		// if bytes.Equal(np, nsp) {
		// 	return nil, ErrNamespaceMismatch
		// }
		height := binary.LittleEndian.Uint64(hbyte)
		b, err := c.client.Blob.Get(ctx, height, nsp, nodeBlob.KzgCommitment(commitment))
		if err != nil {
			return nil, err
		}
		blobs[i] = b.Data
	}
	return blobs, nil

}

// GetIDs returns IDs of all Blobs located in DA at given height.
func (c *RpcNodeClient) GetIDs(ctx context.Context, height uint64, namespace da.Namespace) ([]da.ID, error) {
	shareNamespace, err := share.NamespaceFromBytes(namespace)
	if err != nil {
		return nil, err
	}

	blobs, err := c.client.Blob.GetAll(ctx, height, []share.Namespace{shareNamespace})
	if nil != err {
		return nil, err
	}
	ids := make([]da.ID, len(blobs))
	hbytes := [8]byte{}
	binary.LittleEndian.PutUint64(hbytes[:], height)
	for i, b := range blobs {
		ids[i] = MakeID(hbytes[:], b.Commitment)
	}

	return ids, nil
}

// Commit creates a Commitment for each given Blob.
func (c *RpcNodeClient) Commit(ctx context.Context, blobs []da.Blob, namespace da.Namespace) ([]da.Commitment, error) {
	nsp, err := share.NamespaceFromBytes(namespace)
	if err != nil {
		return nil, err
	}
	cs := make([]da.Commitment, len(blobs))
	for i, b := range blobs {
		bb, err := nodeBlob.NewBlobV0(nsp, b)
		if err != nil {
			return nil, err
		}
		cs[i] = bb.Commitment
	}
	return cs, nil
}

// Validate validates Commitments against the corresponding Proofs. This should be possible without retrieving the Blobs.
func (c *RpcNodeClient) Validate(ctx context.Context, ids []da.ID, proofs []da.Proof, namespace da.Namespace) ([]bool, error) {
	nsp, err := share.NamespaceFromBytes(namespace)
	if err != nil {
		return nil, err
	}
	valids := make([]bool, len(ids))
	for i, id := range ids {
		hbyte, commitment := ParseID(id)
		height := binary.LittleEndian.Uint64(hbyte)
		p, err := c.client.Blob.GetProof(ctx, height, nsp, nodeBlob.KzgCommitment(commitment))
		if err != nil {
			return nil, err
		}

		proof := new(nodeBlob.Proof)
		if err := json.Unmarshal(proofs[i], proof); err != nil {
			return nil, err
		}

		valids[i] = equal(proof, p)
	}
	return valids, nil
}

func equal(p, input *nodeBlob.Proof) bool {

	if p.Len() != input.Len() {
		return false
	}

	for i := 0; i < p.Len(); i++ {
		if (*p)[i].Equal((*input)[i]) {
			return false
		}
	}
	return true
}
