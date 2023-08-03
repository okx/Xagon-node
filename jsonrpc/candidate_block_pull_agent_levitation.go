// LEVITATION_BEGIN

package jsonrpc

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	"time"
)

// CandidateBlockPullAgent is pull candidate blocks from decentralized sequencer
type CandidateBlockPullAgent struct {
	config        Config
	chainID       uint64
	exitRequested bool
}

// NewCandidateBlockPullAgent returns the CandidateBlockPullAgent
func NewCandidateBlockPullAgent(
	cfg Config,
	chainID uint64,
) *CandidateBlockPullAgent {
	agent := &CandidateBlockPullAgent{
		config:        cfg,
		chainID:       chainID,
		exitRequested: false,
	}
	return agent
}

// Start initializes the JSON RPC server to listen for request
func (agent *CandidateBlockPullAgent) Start() error {
	log.Infof("LEVITATION:Starting NewCandidateBlockPullAgent")
	go agent.doWork()
	log.Infof("LEVITATION:Started NewCandidateBlockPullAgent")
	return nil
}

// Stop shutdown the rpc server
func (agent *CandidateBlockPullAgent) Stop() error {
	agent.exitRequested = true
	return nil
}

func (agent *CandidateBlockPullAgent) doWork() {
	log.Infof("LEVITATION:Starting NewCandidateBlockPullAgent work loop")
	for !agent.exitRequested {
		log.Infof("Trying to pull candidate block")
		time.Sleep(1 * time.Second)
	}
}

// LEVITATION_END
