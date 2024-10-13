package ethtxmanager

import (
	"net/http"
	"fmt"
	"net"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"time"
	"encoding/json"
)

const (
	authHandler  = "/inner/v3/okpool/vault/operate/get"
	referOrderId = "referOrderId"
)

type Response struct {
	Code         int    `json:"code"`
	Data         string `json:"data"`
	DetailMsg    string `json:"detailMsg"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Message      string `json:"message"`
}

func (c *Client) startRPC() error {
	if c == nil || !c.cfg.RPC.Enable {
		return nil
	}
	if c.srv != nil {
		return fmt.Errorf("server already started")
	}

	address := fmt.Sprintf("%s:%d", c.cfg.RPC.Host, c.cfg.RPC.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener: %v", err)
		return err
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
			return nil
		}
		log.Errorf("closed http connection: %v", err)
		return err
	}

	return nil
}

func (c *Client) stopRPC() error {
	if c == nil || c.srv == nil {
		return nil
	}

	if err := c.srv.Close(); err != nil {
		log.Errorf("failed to close http server: %v", err)
		return err
	}

	return nil
}

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
