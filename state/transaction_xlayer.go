package state

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

// EstimateGasOpt for a transaction
func (s *State) EstimateGasOpt(transaction *types.Transaction, senderAddress common.Address, isGasFreeSender bool, l2BlockNumber *uint64, dbTx pgx.Tx) (uint64, []byte, error) {
	const ethTransferGas = 21000

	ctx := context.Background()

	t0 := time.Now()

	var l2Block *L2Block
	var err error
	if l2BlockNumber == nil {
		l2Block, err = s.GetLastL2Block(ctx, dbTx)
	} else {
		l2Block, err = s.GetL2BlockByNumber(ctx, *l2BlockNumber, dbTx)
	}
	if err != nil {
		return 0, nil, err
	}

	t1 := time.Now()
	getBlockTime := t1.Sub(t0)

	batch, err := s.GetBatchByL2BlockNumber(ctx, l2Block.NumberU64(), dbTx)
	if err != nil {
		return 0, nil, err
	}

	t2 := time.Now()
	getBatchTime := t2.Sub(t1)

	forkID := s.GetForkIDByBatchNumber(batch.BatchNumber)
	latestL2BlockNumber, err := s.GetLastL2BlockNumber(ctx, dbTx)
	if err != nil {
		return 0, nil, err
	}

	t3 := time.Now()
	getForkIDTime := t3.Sub(t2)

	loadedNonce, err := s.tree.GetNonce(ctx, senderAddress, l2Block.Root().Bytes())
	if err != nil {
		return 0, nil, err
	}
	nonce := loadedNonce.Uint64()

	t4 := time.Now()
	getNonceTime := t4.Sub(t3)

	highEnd := MaxTxGasLimit

	// if gas price is set, set the highEnd to the max amount
	// of the account afford
	isGasPriceSet := !isGasFreeSender && transaction.GasPrice().BitLen() != 0
	if isGasPriceSet {
		senderBalance, err := s.tree.GetBalance(ctx, senderAddress, l2Block.Root().Bytes())
		if errors.Is(err, ErrNotFound) {
			senderBalance = big.NewInt(0)
		} else if err != nil {
			return 0, nil, err
		}

		availableBalance := new(big.Int).Set(senderBalance)
		// check if the account has funds to pay the transfer value
		if transaction.Value() != nil {
			if transaction.Value().Cmp(availableBalance) > 0 {
				return 0, nil, ErrInsufficientFundsForTransfer
			}

			// deduct the value from the available balance
			availableBalance.Sub(availableBalance, transaction.Value())
		}

		// Check the gas allowance for this account, make sure high end is capped to it
		gasAllowance := new(big.Int).Div(availableBalance, transaction.GasPrice())
		if gasAllowance.IsUint64() && highEnd > gasAllowance.Uint64() {
			log.Debugf("Gas estimation high-end capped by allowance [%d]", gasAllowance.Uint64())
			highEnd = gasAllowance.Uint64()
		}
	}

	// if the tx gas is set and it is smaller than the highEnd,
	// limit the highEnd to the maximum allowed by the tx gas
	if transaction.Gas() != 0 && transaction.Gas() < highEnd {
		highEnd = transaction.Gas()
	}

	// set start values for lowEnd and highEnd:
	lowEnd, err := core.IntrinsicGas(transaction.Data(), transaction.AccessList(), s.isContractCreation(transaction), true, false, false)
	if err != nil {
		return 0, nil, err
	}

	// if the intrinsic gas is the same as the constant value for eth transfer
	// and the transaction has a receiver address
	if lowEnd == ethTransferGas && transaction.To() != nil {
		receiver := *transaction.To()
		// check if the receiver address is not a smart contract
		code, err := s.tree.GetCode(ctx, receiver, l2Block.Root().Bytes())
		if err != nil {
			log.Warnf("error while getting code for address %v: %v", receiver.String(), err)
		} else if len(code) == 0 {
			// in case it is just an account, we can avoid the execution and return
			// the transfer constant amount
			return lowEnd, nil, nil
		}
	}

	t5 := time.Now()
	getEndTime := t5.Sub(t4)

	// testTransaction runs the transaction with the specified gas value.
	// it returns a status indicating if the transaction has failed, if it
	// was reverted and the accompanying error

	// Check if the highEnd is a good value to make the transaction pass, if it fails we
	// can return immediately.
	log.Debugf("Estimate gas. Trying to execute TX with %v gas", highEnd)
	var estimationResult *testGasEstimationResult
	if forkID < FORKID_ETROG {
		estimationResult, err = s.internalTestGasEstimationTransactionV1(ctx, batch, l2Block, latestL2BlockNumber, transaction, forkID, senderAddress, isGasFreeSender, highEnd, nonce, false)
	} else {
		estimationResult, err = s.internalTestGasEstimationTransactionV2(ctx, batch, l2Block, latestL2BlockNumber, transaction, forkID, senderAddress, isGasFreeSender, highEnd, nonce, false)
	}
	if err != nil {
		return 0, nil, err
	}
	if estimationResult.failed {
		if estimationResult.reverted {
			return 0, estimationResult.returnValue, estimationResult.executionError
		}

		if estimationResult.ooc {
			return 0, nil, estimationResult.executionError
		}

		// The transaction shouldn't fail, for whatever reason, at highEnd
		return 0, nil, fmt.Errorf(
			"gas required exceeds allowance (%d)",
			highEnd,
		)
	}

	t6 := time.Now()
	internalGasTime := t6.Sub(t5)

	log.Infof("state-EstimateGas time. getBlock:%vms, getBatch:%vms, getForkID:%vms, getNonce:%vms, getEnd:%vms, internalGas:%vms",
		getBlockTime.Milliseconds(), getBatchTime.Milliseconds(), getForkIDTime.Milliseconds(), getNonceTime.Milliseconds(), getEndTime.Milliseconds(), internalGasTime.Milliseconds())

	return estimationResult.gasUsed, nil, nil
}
