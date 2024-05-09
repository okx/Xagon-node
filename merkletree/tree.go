package merkletree

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/hashdb"
	"github.com/ethereum/go-ethereum/common"
)

// StateTree provides methods to access and modify state in merkletree
type StateTree struct {
	grpcClient hashdb.HashDBServiceClient
}

// NewStateTree creates new StateTree.
func NewStateTree(client hashdb.HashDBServiceClient) *StateTree {
	return &StateTree{
		grpcClient: client,
	}
}

// GetBalance returns balance.
func (tree *StateTree) GetBalance(ctx context.Context, address common.Address, root []byte) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyEthAddrBalance(address)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key)
	proof, err := tree.get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}

// GetNonce returns nonce.
func (tree *StateTree) GetNonce(ctx context.Context, address common.Address, root []byte) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyEthAddrNonce(address)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key)
	proof, err := tree.get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}

// GetCodeHash returns code hash.
func (tree *StateTree) GetCodeHash(ctx context.Context, address common.Address, root []byte) ([]byte, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyContractCode(address)
	if err != nil {
		return nil, err
	}
	// this code gets only the hash of the smart contract code from the merkle tree
	k := new(big.Int).SetBytes(key)
	proof, err := tree.get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof.Value == nil {
		return nil, nil
	}

	valueBi := fea2scalar(proof.Value)
	return ScalarToFilledByteSlice(valueBi), nil
}

// GetCode returns code.
func (tree *StateTree) GetCode(ctx context.Context, address common.Address, root []byte) ([]byte, error) {
	scCodeHash, err := tree.GetCodeHash(ctx, address, root)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(scCodeHash)
	if k.Cmp(big.NewInt(0)) == 0 {
		return []byte{}, nil
	}

	// this code gets actual smart contract code from sc code storage
	scCode, err := tree.getProgram(ctx, scalarToh4(k))
	if err != nil {
		return nil, err
	}

	return scCode.Data, nil
}

// GetStorageAt returns Storage Value at specified position.
func (tree *StateTree) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root []byte) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyContractStorage(address, position.Bytes())
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key)
	proof, err := tree.get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}

// SetBalance sets balance.
func (tree *StateTree) SetBalance(ctx context.Context, address common.Address, balance *big.Int, root []byte, uuid string) (newRoot []byte, proof *UpdateProof, err error) {
	if balance.Cmp(big.NewInt(0)) == -1 {
		return nil, nil, fmt.Errorf("invalid balance")
	}

	r := new(big.Int).SetBytes(root)
	key, err := KeyEthAddrBalance(address)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key)
	balanceH8 := scalar2fea(balance)

	updateProof, err := tree.set(ctx, scalarToh4(r), scalarToh4(k), balanceH8, uuid)
	if err != nil {
		return nil, nil, err
	}

	return h4ToFilledByteSlice(updateProof.NewRoot), updateProof, nil
}

// SetNonce sets nonce.
func (tree *StateTree) SetNonce(ctx context.Context, address common.Address, nonce *big.Int, root []byte, uuid string) (newRoot []byte, proof *UpdateProof, err error) {
	if nonce.Cmp(big.NewInt(0)) == -1 {
		return nil, nil, fmt.Errorf("invalid nonce")
	}

	r := new(big.Int).SetBytes(root)
	key, err := KeyEthAddrNonce(address)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key)

	nonceH8 := scalar2fea(nonce)

	updateProof, err := tree.set(ctx, scalarToh4(r), scalarToh4(k), nonceH8, uuid)
	if err != nil {
		return nil, nil, err
	}

	return h4ToFilledByteSlice(updateProof.NewRoot), updateProof, nil
}

// SetCode sets smart contract code.
func (tree *StateTree) SetCode(ctx context.Context, address common.Address, code []byte, root []byte, uuid string) (newRoot []byte, proof *UpdateProof, err error) {
	// calculating smart contract code hash
	scCodeHash4, err := HashContractBytecode(code)
	if err != nil {
		return nil, nil, err
	}

	// store smart contract code by its hash
	err = tree.setProgram(ctx, scCodeHash4, code, true, uuid)
	if err != nil {
		return nil, nil, err
	}

	// set smart contract code hash as a leaf value in merkle tree
	r := new(big.Int).SetBytes(root)
	key, err := KeyContractCode(address)
	if err != nil {
		return nil, nil, err
	}
	k := new(big.Int).SetBytes(key)

	scCodeHash, err := hex.DecodeHex(H4ToString(scCodeHash4))
	if err != nil {
		return nil, nil, err
	}

	scCodeHashBI := new(big.Int).SetBytes(scCodeHash)
	scCodeHashH8 := scalar2fea(scCodeHashBI)

	updateProof, err := tree.set(ctx, scalarToh4(r), scalarToh4(k), scCodeHashH8, uuid)
	if err != nil {
		return nil, nil, err
	}

	// set code length as a leaf value in merkle tree
	key, err = KeyCodeLength(address)
	if err != nil {
		return nil, nil, err
	}
	k = new(big.Int).SetBytes(key)
	scCodeLengthBI := new(big.Int).SetInt64(int64(len(code)))
	scCodeLengthH8 := scalar2fea(scCodeLengthBI)

	updateProof, err = tree.set(ctx, updateProof.NewRoot, scalarToh4(k), scCodeLengthH8, uuid)
	if err != nil {
		return nil, nil, err
	}

	return h4ToFilledByteSlice(updateProof.NewRoot), updateProof, nil
}

// SetStorageAt sets storage value at specified position.
func (tree *StateTree) SetStorageAt(ctx context.Context, address common.Address, position *big.Int, value *big.Int, root []byte, uuid string) (newRoot []byte, proof *UpdateProof, err error) {
	r := new(big.Int).SetBytes(root)
	key, err := KeyContractStorage(address, position.Bytes())
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key)
	valueH8 := scalar2fea(value)
	updateProof, err := tree.set(ctx, scalarToh4(r), scalarToh4(k), valueH8, uuid)
	if err != nil {
		return nil, nil, err
	}

	return h4ToFilledByteSlice(updateProof.NewRoot), updateProof, nil
}

func (tree *StateTree) get(ctx context.Context, root, key []uint64) (*Proof, error) {
	result, err := tree.grpcClient.Get(ctx, &hashdb.GetRequest{
		Root: &hashdb.Fea{Fe0: root[0], Fe1: root[1], Fe2: root[2], Fe3: root[3]},
		Key:  &hashdb.Fea{Fe0: key[0], Fe1: key[1], Fe2: key[2], Fe3: key[3]},
	})
	if err != nil {
		return nil, err
	}

	value, err := string2fea(result.Value)
	if err != nil {
		return nil, err
	}
	return &Proof{
		Root:  []uint64{root[0], root[1], root[2], root[3]},
		Key:   key,
		Value: value,
	}, nil
}

func (tree *StateTree) getProgram(ctx context.Context, key []uint64) (*ProgramProof, error) {
	result, err := tree.grpcClient.GetProgram(ctx, &hashdb.GetProgramRequest{
		Key: &hashdb.Fea{Fe0: key[0], Fe1: key[1], Fe2: key[2], Fe3: key[3]},
	})
	if err != nil {
		return nil, err
	}

	return &ProgramProof{
		Data: result.Data,
	}, nil
}

func (tree *StateTree) set(ctx context.Context, oldRoot, key, value []uint64, uuid string) (*UpdateProof, error) {
	feaValue := fea2string(value)
	if strings.HasPrefix(feaValue, "0x") { // nolint
		feaValue = feaValue[2:]
	}
	result, err := tree.grpcClient.Set(ctx, &hashdb.SetRequest{
		OldRoot:     &hashdb.Fea{Fe0: oldRoot[0], Fe1: oldRoot[1], Fe2: oldRoot[2], Fe3: oldRoot[3]},
		Key:         &hashdb.Fea{Fe0: key[0], Fe1: key[1], Fe2: key[2], Fe3: key[3]},
		Value:       feaValue,
		Details:     false,
		Persistence: hashdb.Persistence_PERSISTENCE_DATABASE,
		BatchUuid:   uuid,
		TxIndex:     0,
		BlockIndex:  0,
	})
	if err != nil {
		return nil, err
	}

	var newValue []uint64
	if result.NewValue != "" {
		newValue, err = string2fea(result.NewValue)
		if err != nil {
			return nil, err
		}
	}

	return &UpdateProof{
		OldRoot:  oldRoot,
		NewRoot:  []uint64{result.NewRoot.Fe0, result.NewRoot.Fe1, result.NewRoot.Fe2, result.NewRoot.Fe3},
		Key:      key,
		NewValue: newValue,
	}, nil
}

func (tree *StateTree) setProgram(ctx context.Context, key []uint64, data []byte, persistent bool, uuid string) error {
	persistence := hashdb.Persistence_PERSISTENCE_TEMPORARY
	if persistent {
		persistence = hashdb.Persistence_PERSISTENCE_DATABASE
	}

	_, err := tree.grpcClient.SetProgram(ctx, &hashdb.SetProgramRequest{
		Key:         &hashdb.Fea{Fe0: key[0], Fe1: key[1], Fe2: key[2], Fe3: key[3]},
		Data:        data,
		Persistence: persistence,
		BatchUuid:   uuid,
		TxIndex:     0,
		BlockIndex:  0,
	})
	return err
}

// Flush flushes all changes to the persistent storage.
func (tree *StateTree) Flush(ctx context.Context, newStateRoot common.Hash, uuid string) error {
	flushRequest := &hashdb.FlushRequest{BatchUuid: uuid, NewStateRoot: newStateRoot.String(), Persistence: hashdb.Persistence_PERSISTENCE_DATABASE}
	_, err := tree.grpcClient.Flush(ctx, flushRequest)
	return err
}

// StartBlock starts a new block.
func (tree *StateTree) StartBlock(ctx context.Context, oldRoot common.Hash, uuid string) error {
	startBlockRequest := &hashdb.StartBlockRequest{
		BatchUuid:    uuid,
		OldStateRoot: oldRoot.String(),
		Persistence:  hashdb.Persistence_PERSISTENCE_DATABASE}
	_, err := tree.grpcClient.StartBlock(ctx, startBlockRequest)
	return err
}

// FinishBlock finishes a block.
func (tree *StateTree) FinishBlock(ctx context.Context, newRoot common.Hash, uuid string) error {
	finishBlockRequest := &hashdb.FinishBlockRequest{
		BatchUuid:    uuid,
		NewStateRoot: newRoot.String(),
		Persistence:  hashdb.Persistence_PERSISTENCE_DATABASE}
	_, err := tree.grpcClient.FinishBlock(ctx, finishBlockRequest)

	return err
}
