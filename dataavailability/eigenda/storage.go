package eigenda

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// ErrNotFound is used when DA information is not found
var ErrNotFound = errors.New("not found")

// DAStorage is an in-memory data availability storage for the mock implementation.
// Will need a thorough design to think through the backend storage for DA information on
// the node. For now, this mock node PoC quickly implements this by indexing block hashes
// to the index of the sequence data that is submitted on the EigenDA layer.
type DAStorage struct {
	inner map[common.Hash][]byte
	mutex *sync.RWMutex
}

// Get tries to get the DA information indexed to the block hash.
func (s *DAStorage) Get(hash common.Hash) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	info, ok := s.inner[hash]
	if ok {
		return info, nil
	} else {
		return nil, ErrNotFound
	}
}

// Add indexes the DA information with the block hash into the in-memory storage.
func (s *DAStorage) Add(hash common.Hash, message []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.inner[hash] = message
	return nil
}

// Update updates the DA information with the block hash into the in-memory storage.
func (s *DAStorage) Update(hash common.Hash, message []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.inner[hash] = message
	return nil
}
