package ethtxmanager

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	authHandler  = "/inner/v3/okpool/vault/operate/get"
	referOrderId = "referOrderId"
)

// Response is the response structure for the rpc server
type Response struct {
	Code         int    `json:"code"`
	Data         string `json:"data"`
	DetailMsg    string `json:"detailMsg"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Message      string `json:"message"`
}

func (c *Client) startRPC() {
	if c == nil || !c.cfg.HTTP.Enable {
		log.Infof("rpc server is disabled")
		return
	}
	if c.srv != nil {
		log.Errorf("server already started")
		return
	}

	address := fmt.Sprintf("%s:%d", c.cfg.HTTP.Host, c.cfg.HTTP.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener: %v", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle(authHandler, c)

	c.srv = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: time.Minute,
		ReadTimeout:       time.Minute,
		WriteTimeout:      time.Minute,
	}
	log.Infof("http server started: %s", address)
	if err = c.srv.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Infof("http server stopped")
			return
		}
		log.Errorf("closed http connection: %v", err)
		return
	}

	return
}

func (c *Client) stopRPC() {
	if c == nil || c.srv == nil {
		return
	}

	if err := c.srv.Close(); err != nil {
		log.Errorf("failed to close http server: %v", err)
		return
	}

	return
}

// ServeHTTP handles the incoming HTTP requests
func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		var resp Response
		resp.Code = http.StatusMethodNotAllowed
		resp.ErrorMessage = "method not allowed"
		respBytes, _ := json.Marshal(resp)
		http.Error(w, string(respBytes), http.StatusMethodNotAllowed)

		return
	}
	queryParams := r.URL.Query()
	orderID := queryParams.Get(referOrderId)
	if orderID == "" {
		var resp Response
		resp.Code = http.StatusBadRequest
		resp.ErrorMessage = "order id not found"
		respBytes, _ := json.Marshal(resp)
		http.Error(w, string(respBytes), http.StatusBadRequest)
		return
	}

	data := getAuthInstance().get(orderID)
	if data == "" {
		var resp Response
		resp.Code = http.StatusBadRequest
		resp.ErrorMessage = "order id not found in cache"
		respBytes, _ := json.Marshal(resp)
		http.Error(w, string(respBytes), http.StatusBadRequest)
		return
	}
	var resp Response
	resp.Data = data
	respBytes, _ := json.Marshal(resp)
	w.Write(respBytes)
}
