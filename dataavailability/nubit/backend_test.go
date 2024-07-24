package nubit

import (
	"crypto/ecdsa"
	"encoding/hex"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/rollkit/go-da/proxy"
)

func NewNubitDABackendTest(url string, authKey string, pk *ecdsa.PrivateKey) (*NubitDABackend, error) {
	cn, err := proxy.NewClient(url, authKey)
	if err != nil || cn == nil {
		return nil, err
	}

	name, err := hex.DecodeString("00000000000000000000000000000000000000000000706f6c79676f6e")
	if err != nil {
		return nil, err
	}

	log.Infof("Nubit Namespace: %s ", string(name))
	return &NubitDABackend{
		ns:         name,
		client:     cn,
		privKey:    pk,
		commitTime: time.Now(),
	}, nil
}
