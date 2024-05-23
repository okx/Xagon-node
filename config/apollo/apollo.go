package apollo

import (
	"errors"
	"strings"

	nodeconfig "github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
)

// Client is the apollo client
type Client struct {
	*agollo.Client
	config *nodeconfig.Config
}

// NewClient creates a new apollo client
func NewClient(conf *nodeconfig.Config) *Client {
	if conf == nil || !conf.Apollo.Enable || conf.Apollo.IP == "" || conf.Apollo.AppID == "" || len(conf.Apollo.NamespaceName) == 0 {
		log.Infof("apollo is not enabled, config: %+v", conf.Apollo)
		return nil
	}
	unitNamespace, err := validateNamespace(conf.Apollo.NamespaceName)
	if err != nil {
		log.Fatalf("failed init apollo: %v", err)
		return nil
	}
	unitNamespaceMap = unitNamespace
	var namespaces []string
	for _, namespace := range unitNamespace {
		namespaces = append(namespaces, namespace)
	}
	c := &config.AppConfig{
		IP:             conf.Apollo.IP,
		AppID:          conf.Apollo.AppID,
		NamespaceName:  strings.Join(namespaces, ","),
		Cluster:        "default",
		IsBackupConfig: false,
	}

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		log.Fatalf("failed init apollo: %v", err)
	}

	apc := &Client{
		Client: client,
		config: conf,
	}
	client.AddChangeListener(&CustomChangeListener{apc})

	return apc
}

var unitNamespaceMap = make(map[string]string)

func validateNamespace(namespaceItem string) (map[string]string, error) {
	namespaceitems := strings.Split(namespaceItem, ",")
	var unitNamespaceMapi = make(map[string]string)
	for _, item := range namespaceitems {
		split := strings.SplitN(item, ":", 2)
		if len(split) != 2 {
			log.Errorf("invalid namespace: %s, please configure in correct format \"unit:namespace\"", item)
			return nil, errors.New("invalid namespace")
		}
		unitNamespaceMapi[split[0]] = split[1]
	}
	return unitNamespaceMapi, nil
}

// LoadConfig loads the config
func (c *Client) LoadConfig() (loaded bool) {
	if c == nil {
		return false
	}
	for unit, namesapce := range unitNamespaceMap {
		cache := c.GetConfigCache(namesapce)
		var tempval interface{}
		cache.Range(func(key, value interface{}) bool {
			tempval = value
			return true
		})
		loaded = true
		switch unit {
		case RPC:
			c.loadJsonRPC(tempval)
		case Sequencer:
			c.loadSequencer(tempval)
		case L2GasPriceSuggester:
			c.loadL2GasPricer(tempval)
		case Pool:
			c.loadPool(tempval)
		}
	}
	return
}

// CustomChangeListener is the custom change listener
type CustomChangeListener struct {
	*Client
}

// OnChange is the change listener
func (c *CustomChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	var namespaceunitmap = make(map[string]string)
	for unit, namespace := range unitNamespaceMap {
		namespaceunitmap[namespace] = unit
	}
	for key, value := range changeEvent.Changes {
		if unit, ok := namespaceunitmap[changeEvent.Namespace]; ok {
			switch unit {
			case Halt:
				c.fireHalt(key, value)
			case RPC:
				c.fireJsonRPC(key, value)
			case Sequencer:
				c.fireSequencer(key, value)
			case L2GasPriceSuggester:
				c.fireL2GasPricer(key, value)
			case Pool:
				c.firePool(key, value)
			}
		}
	}
}

// OnNewestChange is the newest change listener
func (c *CustomChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
}
