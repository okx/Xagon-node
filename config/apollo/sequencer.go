package apollo

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/apolloconfig/agollo/v4/storage"
)

func (c *Client) loadSequencer(value interface{}) {
	dstConf, err := c.unmarshal(value)
	if err != nil {
		log.Fatalf("failed to unmarshal sequncer config: %v", err)
	}
	c.config.Sequencer = dstConf.Sequencer
	c.config.Sequencer.StreamServer.Log.Outputs = make([]string, len(dstConf.Sequencer.StreamServer.Log.Outputs))
	copy(c.config.Sequencer.StreamServer.Log.Outputs, dstConf.Sequencer.StreamServer.Log.Outputs)

	log.Infof("loaded sequencer from apollo config: %+v", value.(string))
}

// fireSequencer fires the sequencer config change
// FullBatchSleepDuration
func (c *Client) fireSequencer(key string, value *storage.ConfigChange) {
	newConf, err := c.unmarshal(value.NewValue)
	if err != nil {
		log.Errorf("failed to unmarshal sequencer config: %v error: %v", value.NewValue, err)
		return
	}
	log.Infof("apollo sequencer old config : %+v", value.OldValue.(string))
	log.Infof("apollo sequencer config changed: %+v", value.NewValue.(string))
	sequencer.UpdateConfig(newConf.Sequencer)
}
