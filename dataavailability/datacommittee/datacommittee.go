package datacommittee

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"strings"

	"github.com/0xPolygon/cdk-data-availability/client"
	daTypes "github.com/0xPolygon/cdk-data-availability/types"
	polygondatacommittee "github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygondatacommittee_xlayer"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/net/context"
)

const unexpectedHashTemplate = "missmatch on transaction data. Expected hash %s, actual hash: %s"

// DataCommitteeMember represents a member of the Data Committee
type DataCommitteeMember struct {
	Addr common.Address
	URL  string
}

// DataCommittee represents a specific committee
type DataCommittee struct {
	AddressesHash      common.Hash
	Members            []DataCommitteeMember
	RequiredSignatures uint64
}

// DataCommitteeBackend implements the DAC integration
type DataCommitteeBackend struct {
	dataCommitteeContract      *polygondatacommittee.PolygondatacommitteeXlayer
	privKey                    *ecdsa.PrivateKey
	dataCommitteeClientFactory client.Factory

	committeeMembers        []DataCommitteeMember
	selectedCommitteeMember int
	ctx                     context.Context
}

// New creates an instance of DataCommitteeBackend
func New(
	l1RPCURL string,
	dataCommitteeAddr common.Address,
	privKey *ecdsa.PrivateKey,
	dataCommitteeClientFactory client.Factory,
) (*DataCommitteeBackend, error) {
	ethClient, err := ethclient.Dial(l1RPCURL)
	if err != nil {
		log.Errorf("error connecting to %s: %+v", l1RPCURL, err)
		return nil, err
	}
	dataCommittee, err := polygondatacommittee.NewPolygondatacommitteeXlayer(dataCommitteeAddr, ethClient)
	if err != nil {
		return nil, err
	}
	return &DataCommitteeBackend{
		dataCommitteeContract:      dataCommittee,
		privKey:                    privKey,
		dataCommitteeClientFactory: dataCommitteeClientFactory,
		ctx:                        context.Background(),
	}, nil
}

// Init loads the DAC to be cached when needed
func (d *DataCommitteeBackend) Init() error {
	committee, err := d.getCurrentDataCommittee()
	if err != nil {
		return err
	}
	selectedCommitteeMember := -1
	if committee != nil {
		d.committeeMembers = committee.Members
		if len(committee.Members) > 0 {
			selectedCommitteeMember = rand.Intn(len(committee.Members)) //nolint:gosec
		}
	}
	d.selectedCommitteeMember = selectedCommitteeMember
	return nil
}

// GetSequence gets backend data one hash at a time. This should be optimized on the DAC side to get them all at once.
func (d *DataCommitteeBackend) GetSequence(ctx context.Context, hashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	// TODO: optimize this on the DAC side by implementing a multi batch retrieve api
	var batchData [][]byte
	for _, h := range hashes {
		data, err := d.GetBatchL2Data(h)
		if err != nil {
			return nil, err
		}
		batchData = append(batchData, data)
	}
	return batchData, nil
}

// GetBatchL2Data returns the data from the DAC. It checks that it matches with the expected hash
func (d *DataCommitteeBackend) GetBatchL2Data(hash common.Hash) ([]byte, error) {
	intialMember := d.selectedCommitteeMember
	found := false
	for !found && intialMember != -1 {
		member := d.committeeMembers[d.selectedCommitteeMember]
		log.Infof("trying to get data from %s at %s", member.Addr.Hex(), member.URL)
		c := d.dataCommitteeClientFactory.New(member.URL)
		data, err := c.GetOffChainData(d.ctx, hash)
		if err != nil {
			log.Warnf(
				"error getting data from DAC node %s at %s: %s",
				member.Addr.Hex(), member.URL, err,
			)
			d.selectedCommitteeMember = (d.selectedCommitteeMember + 1) % len(d.committeeMembers)
			if d.selectedCommitteeMember == intialMember {
				break
			}
			continue
		}
		actualTransactionsHash := crypto.Keccak256Hash(data)
		if actualTransactionsHash != hash {
			unexpectedHash := fmt.Errorf(
				unexpectedHashTemplate, hash, actualTransactionsHash,
			)
			log.Warnf(
				"error getting data from DAC node %s at %s: %s",
				member.Addr.Hex(), member.URL, unexpectedHash,
			)
			d.selectedCommitteeMember = (d.selectedCommitteeMember + 1) % len(d.committeeMembers)
			if d.selectedCommitteeMember == intialMember {
				break
			}
			continue
		}
		return data, nil
	}
	if err := d.Init(); err != nil {
		return nil, fmt.Errorf("error loading data committee: %s", err)
	}
	return nil, fmt.Errorf("couldn't get the data from any committee member")
}

type signatureMsg struct {
	addr      common.Address
	signature []byte
	err       error
}

// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
// as expected by the contract
func (s *DataCommitteeBackend) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	// Get current committee
	committee, err := s.getCurrentDataCommittee()
	if err != nil {
		return nil, err
	}

	// Authenticate as trusted sequencer by signing the sequences
	sequence := daTypes.Sequence{}
	for _, seq := range batchesData {
		sequence = append(sequence, seq)
	}
	signedSequence, err := sequence.Sign(s.privKey)
	if err != nil {
		return nil, err
	}

	// Request signatures to all members in parallel
	ch := make(chan signatureMsg, len(committee.Members))
	signatureCtx, cancelSignatureCollection := context.WithCancel(ctx)
	for _, member := range committee.Members {
		go requestSignatureFromMember(signatureCtx, *signedSequence, member, ch)
	}

	// Collect signatures
	msgs := []signatureMsg{}
	var (
		collectedSignatures uint64
		failedToCollect     uint64
	)
	for collectedSignatures < committee.RequiredSignatures {
		msg := <-ch
		if msg.err != nil {
			log.Errorf("error when trying to get signature from %s: %s", msg.addr, msg.err)
			failedToCollect++
			if len(committee.Members)-int(failedToCollect) < int(committee.RequiredSignatures) {
				cancelSignatureCollection()
				return nil, errors.New("too many members failed to send their signature")
			}
		} else {
			log.Infof("received signature from %s", msg.addr)
			collectedSignatures++
		}
		msgs = append(msgs, msg)
	}

	// Stop requesting as soon as we have N valid signatures
	cancelSignatureCollection()

	return buildSignaturesAndAddrs(signatureMsgs(msgs), committee.Members), nil
}

func requestSignatureFromMember(ctx context.Context, signedSequence daTypes.SignedSequence, member DataCommitteeMember, ch chan signatureMsg) {
	// request
	c := client.New(member.URL)
	log.Infof("sending request to sign the sequence to %s at %s", member.Addr.Hex(), member.URL)
	signature, err := c.SignSequence(signedSequence)
	if err != nil {
		ch <- signatureMsg{
			addr: member.Addr,
			err:  err,
		}
		return
	}
	// verify returned signature
	signedSequence.Signature = signature
	signer, err := signedSequence.Signer()
	if err != nil {
		ch <- signatureMsg{
			addr: member.Addr,
			err:  err,
		}
		return
	}
	if signer != member.Addr {
		ch <- signatureMsg{
			addr: member.Addr,
			err:  fmt.Errorf("invalid signer. Expected %s, actual %s", member.Addr.Hex(), signer.Hex()),
		}
		return
	}
	ch <- signatureMsg{
		addr:      member.Addr,
		signature: signature,
	}
}

func buildSignaturesAndAddrs(sigs signatureMsgs, members []DataCommitteeMember) []byte {
	const (
		sigLen  = 65
		addrLen = 20
	)
	res := make([]byte, 0, len(sigs)*sigLen+len(members)*addrLen)
	sort.Sort(sigs)
	for _, msg := range sigs {
		log.Debugf("adding signature %s from %s", common.Bytes2Hex(msg.signature), msg.addr.Hex())
		res = append(res, msg.signature...)
	}
	for _, member := range members {
		log.Debugf("adding addr %s", common.Bytes2Hex(member.Addr.Bytes()))
		res = append(res, member.Addr.Bytes()...)
	}
	log.Debugf("full res %s", common.Bytes2Hex(res))
	return res
}

type signatureMsgs []signatureMsg

func (s signatureMsgs) Len() int { return len(s) }
func (s signatureMsgs) Less(i, j int) bool {
	return strings.ToUpper(s[i].addr.Hex()) < strings.ToUpper(s[j].addr.Hex())
}
func (s signatureMsgs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// getCurrentDataCommittee return the currently registered data committee
func (d *DataCommitteeBackend) getCurrentDataCommittee() (*DataCommittee, error) {
	addrsHash, err := d.dataCommitteeContract.CommitteeHash(&bind.CallOpts{Pending: false})
	if err != nil {
		return nil, fmt.Errorf("error getting CommitteeHash from L1 SC: %w", err)
	}
	reqSign, err := d.dataCommitteeContract.RequiredAmountOfSignatures(&bind.CallOpts{Pending: false})
	if err != nil {
		return nil, fmt.Errorf("error getting RequiredAmountOfSignatures from L1 SC: %w", err)
	}
	members, err := d.getCurrentDataCommitteeMembers()
	if err != nil {
		return nil, err
	}

	return &DataCommittee{
		AddressesHash:      common.Hash(addrsHash),
		RequiredSignatures: reqSign.Uint64(),
		Members:            members,
	}, nil
}

// getCurrentDataCommitteeMembers return the currently registered data committee members
func (d *DataCommitteeBackend) getCurrentDataCommitteeMembers() ([]DataCommitteeMember, error) {
	nMembers, err := d.dataCommitteeContract.GetAmountOfMembers(&bind.CallOpts{Pending: false})
	if err != nil {
		return nil, fmt.Errorf("error getting GetAmountOfMembers from L1 SC: %w", err)
	}
	members := make([]DataCommitteeMember, 0, nMembers.Int64())
	for i := int64(0); i < nMembers.Int64(); i++ {
		member, err := d.dataCommitteeContract.Members(&bind.CallOpts{Pending: false}, big.NewInt(i))
		if err != nil {
			return nil, fmt.Errorf("error getting Members %d from L1 SC: %w", i, err)
		}
		members = append(members, DataCommitteeMember{
			Addr: member.Addr,
			URL:  member.Url,
		})
	}
	return members, nil
}
