package jsonrpc

import (
	"net/http"
	"path"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/metrics"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
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
}

type apiAllow struct {
	allowKeys map[string]string
	enable    bool
}

var al apiAllow

// InitApiAuth initializes the api authentication
func InitApiAuth(a ApiAuthConfig) {
	setApiAuth(a)
}

// setApiAuth sets the api authentication
func setApiAuth(a ApiAuthConfig) {
	al.enable = a.Enabled
	var tmp = make(map[string]string)
	for _, k := range a.ApiKeys {
		tmp[k.Key] = k.Project
	}
	al.allowKeys = tmp
}

func check(key string) bool {
	if _, ok := al.allowKeys[key]; ok {
		metrics.RequestAuthCount(al.allowKeys[key])
		return true
	}
	metrics.RequestAuthErrorCount()
	return false
}

func apiAuthHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if al.enable {
			base := path.Base(r.URL.Path)
			if base == "" || !check(base) {
				err := handleNoAuthErr(w)
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

func handleNoAuthErr(w http.ResponseWriter) error {
	respbytes, err := types.NewResponse(types.Request{JSONRPC: "2.0", ID: 0}, nil, types.NewRPCError(types.InvalidParamsErrorCode, "no authentication")).Bytes()
	if err != nil {
		return err
	}
	_, err = w.Write(respbytes)
	if err != nil {
		return err
	}
	return nil
}
