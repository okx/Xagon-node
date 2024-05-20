package jsonrpc

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/metrics"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// ApiAuthConfig is the api authentication config
type ApiAuthConfig struct {
	// Enabled defines if the api authentication is enabled
	Enabled bool `mapstructure:"Enabled"`
	// ApiKeys defines the api keys
	ApiKeys []KeyItem `mapstructure:"ApiKeys"`
}

// KeyItem is the api key item
type KeyItem struct {
	// Name defines the name of the key
	Project string `mapstructure:"Project"`
	// Key defines the key
	Key string `mapstructure:"Key"`
	// Timeout defines the timeout
	Timeout string `mapstructure:"Timeout"`
}

type apiAllow struct {
	allowKeys map[string]keyItem
	enable    bool
}

type keyItem struct {
	project string
	timeout time.Time
}

var al apiAllow

// InitApiAuth initializes the api authentication
func InitApiAuth(a ApiAuthConfig) {
	setApiAuth(a)
}

// setApiAuth sets the api authentication
func setApiAuth(a ApiAuthConfig) {
	al.enable = a.Enabled
	var tmp = make(map[string]keyItem)
	for _, k := range a.ApiKeys {
		k.Key = strings.ToLower(k.Key)
		parse, err := time.Parse("2006-01-02", k.Timeout)
		if err != nil {
			log.Warnf("parse key [%+v], error parsing timeout: %v", k, err)
			continue
		}
		if strings.ToLower(fmt.Sprintf("%x", md5.Sum([]byte(k.Project+k.Timeout)))) != k.Key {
			log.Warnf("project [%s], key [%s] is invalid, key = md5(Project+Timeout)", k.Project, k.Key)
			continue
		}
		tmp[k.Key] = keyItem{project: k.Project, timeout: parse}
	}
	al.allowKeys = tmp
}

func check(key string) error {
	key = strings.ToLower(key)
	if item, ok := al.allowKeys[key]; ok && time.Now().Before(item.timeout) {
		metrics.RequestAuthCount(al.allowKeys[key].project)
		return nil
	} else if ok && time.Now().After(item.timeout) {
		log.Warnf("project [%s], key [%s] has expired, ", item.project, key)
		metrics.RequestAuthErrorCount(metrics.RequestAuthErrorTypeKeyExpired)
		return errors.New("key has expired")
	}
	metrics.RequestAuthErrorCount(metrics.RequestAuthErrorTypeNoAuth)
	return errors.New("no authentication")
}

func apiAuthHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if al.enable {
			if er := check(path.Base(r.URL.Path)); er != nil {
				err := handleNoAuthErr(w, er)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
		handlerFunc(w, r)
	}
}

func apiAuthHandler(next http.Handler) http.Handler {
	return apiAuthHandlerFunc(next.ServeHTTP)
}

func handleNoAuthErr(w http.ResponseWriter, err error) error {
	respbytes, err := types.NewResponse(types.Request{JSONRPC: "2.0", ID: 0}, nil, types.NewRPCError(types.InvalidParamsErrorCode, err.Error())).Bytes()
	if err != nil {
		return err
	}
	_, err = w.Write(respbytes)
	if err != nil {
		return err
	}
	return nil
}
