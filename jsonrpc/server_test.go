package jsonrpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/mocks"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	maxRequestsPerIPAndSecond        = 1000
	chainID                   uint64 = 1000
)

type mockedServer struct {
	Config              Config
	Server              *Server
	ServerURL           string
	ServerWebSocketsURL string
}

type mocksWrapper struct {
	Pool     *mocks.PoolMock
	State    *mocks.StateMock
	Etherman *mocks.EthermanMock
	Storage  *storageMock
}

func newMockedServer(t *testing.T, cfg Config) (*mockedServer, *mocksWrapper, *ethclient.Client) {
	pool := mocks.NewPoolMock(t)
	st := mocks.NewStateMock(t)
	etherman := mocks.NewEthermanMock(t)
	storage := newStorageMock(t)
	apis := map[string]bool{
		APIEth:    true,
		APINet:    true,
		APIDebug:  true,
		APIZKEVM:  true,
		APITxPool: true,
		APIWeb3:   true,
	}

	var newL2BlockEventHandler state.NewL2BlockEventHandler = func(e state.NewL2BlockEvent) {}
	st.On("RegisterNewL2BlockEventHandler", mock.IsType(newL2BlockEventHandler)).Once()
	st.On("StartToMonitorNewL2Blocks").Once()

	services := []Service{}
	if _, ok := apis[APIEth]; ok {
		services = append(services, Service{
			Name:    APIEth,
			Service: NewEthEndpoints(cfg, chainID, pool, st, etherman, storage),
		})
	}

	if _, ok := apis[APINet]; ok {
		services = append(services, Service{
			Name:    APINet,
			Service: NewNetEndpoints(cfg, chainID),
		})
	}

	if _, ok := apis[APIZKEVM]; ok {
		services = append(services, Service{
			Name:    APIZKEVM,
			Service: NewZKEVMEndpoints(cfg, pool, st, etherman),
		})
	}

	if _, ok := apis[APITxPool]; ok {
		services = append(services, Service{
			Name:    APITxPool,
			Service: &TxPoolEndpoints{},
		})
	}

	if _, ok := apis[APIDebug]; ok {
		services = append(services, Service{
			Name:    APIDebug,
			Service: NewDebugEndpoints(cfg, st, etherman),
		})
	}

	if _, ok := apis[APIWeb3]; ok {
		services = append(services, Service{
			Name:    APIWeb3,
			Service: &Web3Endpoints{},
		})
	}
	server := NewServer(cfg, chainID, pool, st, storage, services)

	go func() {
		err := server.Start()
		if err != nil {
			panic(err)
		}
	}()

	serverURL := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	for {
		fmt.Println("waiting server to get ready...") // fmt is used here to avoid race condition with logs
		res, err := http.Get(serverURL)               //nolint:gosec
		if err == nil && res.StatusCode == http.StatusOK {
			fmt.Println("server ready!") // fmt is used here to avoid race condition with logs
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	ethClient, err := ethclient.Dial(serverURL)
	require.NoError(t, err)

	serverWebSocketsURL := fmt.Sprintf("ws://%s:%d", cfg.WebSockets.Host, cfg.WebSockets.Port)

	msv := &mockedServer{
		Config:              cfg,
		Server:              server,
		ServerURL:           serverURL,
		ServerWebSocketsURL: serverWebSocketsURL,
	}

	mks := &mocksWrapper{
		Pool:     pool,
		State:    st,
		Etherman: etherman,
		Storage:  storage,
	}

	return msv, mks, ethClient
}

func getSequencerDefaultConfig() Config {
	cfg := Config{
		Host:                         "0.0.0.0",
		Port:                         9123,
		MaxRequestsPerIPAndSecond:    maxRequestsPerIPAndSecond,
		MaxCumulativeGasUsed:         300000,
		BatchRequestsEnabled:         true,
		MaxLogsCount:                 10000,
		MaxLogsBlockRange:            10000,
		MaxNativeBlockHashBlockRange: 60000,
		WebSockets: WebSocketsConfig{
			Enabled:   true,
			Host:      "0.0.0.0",
			Port:      9133,
			ReadLimit: 0,
		},
	}
	return cfg
}

func getNonSequencerDefaultConfig(sequencerNodeURI string) Config {
	cfg := getSequencerDefaultConfig()
	cfg.Port = 9124
	cfg.SequencerNodeURI = sequencerNodeURI
	return cfg
}

func newSequencerMockedServer(t *testing.T) (*mockedServer, *mocksWrapper, *ethclient.Client) {
	cfg := getSequencerDefaultConfig()
	return newMockedServer(t, cfg)
}

func newMockedServerWithCustomConfig(t *testing.T, cfg Config) (*mockedServer, *mocksWrapper, *ethclient.Client) {
	return newMockedServer(t, cfg)
}

func newNonSequencerMockedServer(t *testing.T, sequencerNodeURI string) (*mockedServer, *mocksWrapper, *ethclient.Client) {
	cfg := getNonSequencerDefaultConfig(sequencerNodeURI)
	return newMockedServer(t, cfg)
}

func (s *mockedServer) GetWSClient() *ethclient.Client {
	ethClient, err := ethclient.Dial(s.ServerWebSocketsURL)
	if err != nil {
		panic(err)
	}

	return ethClient
}

func (s *mockedServer) Stop() {
	err := s.Server.Stop()
	if err != nil {
		panic(err)
	}
}

func (s *mockedServer) JSONRPCCall(method string, parameters ...interface{}) (types.Response, error) {
	return client.JSONRPCCall(s.ServerURL, method, parameters...)
}

func (s *mockedServer) JSONRPCBatchCall(calls ...client.BatchCall) ([]types.Response, error) {
	return client.JSONRPCBatchCall(s.ServerURL, calls...)
}

func (s *mockedServer) ChainID() uint64 {
	return chainID
}

func TestBatchRequests(t *testing.T) {
	type testCase struct {
		Name                 string
		BatchRequestsEnabled bool
		BatchRequestsLimit   uint
		NumberOfRequests     int
		ExpectedError        error
		SetupMocks           func(m *mocksWrapper, tc testCase)
	}

	st := trie.NewStackTrie(nil)
	block := state.NewL2Block(
		state.NewL2Header(&ethTypes.Header{Number: big.NewInt(2), UncleHash: ethTypes.EmptyUncleHash, Root: ethTypes.EmptyRootHash}),
		[]*ethTypes.Transaction{ethTypes.NewTransaction(1, common.Address{}, big.NewInt(1), 1, big.NewInt(1), []byte{})},
		nil,
		[]*ethTypes.Receipt{ethTypes.NewReceipt([]byte{}, false, uint64(0))},
		st,
	)

	testCases := []testCase{
		{
			Name:                 "batch requests disabled",
			BatchRequestsEnabled: false,
			BatchRequestsLimit:   0,
			NumberOfRequests:     10,
			ExpectedError:        fmt.Errorf("400 - " + types.ErrBatchRequestsDisabled.Error() + "\n"),
			SetupMocks:           func(m *mocksWrapper, tc testCase) {},
		},
		{
			Name:                 "batch requests over the limit",
			BatchRequestsEnabled: true,
			BatchRequestsLimit:   5,
			NumberOfRequests:     6,
			ExpectedError:        fmt.Errorf("413 - " + types.ErrBatchRequestsLimitExceeded.Error() + "\n"),
			SetupMocks: func(m *mocksWrapper, tc testCase) {
			},
		},
		{
			Name:                 "batch requests unlimited",
			BatchRequestsEnabled: true,
			BatchRequestsLimit:   0,
			NumberOfRequests:     100,
			ExpectedError:        nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.On("GetLastL2BlockNumber", context.Background(), nil).Return(block.Number().Uint64(), nil).Times(tc.NumberOfRequests)
				m.State.On("GetL2BlockByNumber", context.Background(), block.Number().Uint64(), nil).Return(block, nil).Times(tc.NumberOfRequests)
				m.State.On("GetTransactionReceipt", context.Background(), mock.Anything, nil).Return(ethTypes.NewReceipt([]byte{}, false, uint64(0)), nil)
			},
		},
		{
			Name:                 "batch requests equal the limit",
			BatchRequestsEnabled: true,
			BatchRequestsLimit:   5,
			NumberOfRequests:     5,
			ExpectedError:        nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.On("GetLastL2BlockNumber", context.Background(), nil).Return(block.Number().Uint64(), nil).Times(tc.NumberOfRequests)
				m.State.On("GetL2BlockByNumber", context.Background(), block.Number().Uint64(), nil).Return(block, nil).Times(tc.NumberOfRequests)
				m.State.On("GetTransactionReceipt", context.Background(), mock.Anything, nil).Return(ethTypes.NewReceipt([]byte{}, false, uint64(0)), nil)
			},
		},
		{
			Name:                 "batch requests under the limit",
			BatchRequestsEnabled: true,
			BatchRequestsLimit:   5,
			NumberOfRequests:     4,
			ExpectedError:        nil,
			SetupMocks: func(m *mocksWrapper, tc testCase) {
				m.State.On("GetLastL2BlockNumber", context.Background(), nil).Return(block.Number().Uint64(), nil).Times(tc.NumberOfRequests)
				m.State.On("GetL2BlockByNumber", context.Background(), block.Number().Uint64(), nil).Return(block, nil).Times(tc.NumberOfRequests)
				m.State.On("GetTransactionReceipt", context.Background(), mock.Anything, nil).Return(ethTypes.NewReceipt([]byte{}, false, uint64(0)), nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase

			cfg := getSequencerDefaultConfig()
			cfg.BatchRequestsEnabled = tc.BatchRequestsEnabled
			cfg.BatchRequestsLimit = tc.BatchRequestsLimit
			s, m, _ := newMockedServerWithCustomConfig(t, cfg)

			tc.SetupMocks(m, tc)

			calls := []client.BatchCall{}

			for i := 0; i < tc.NumberOfRequests; i++ {
				calls = append(calls, client.BatchCall{
					Method:     "eth_getBlockByNumber",
					Parameters: []interface{}{"latest"},
				})
			}

			result, err := s.JSONRPCBatchCall(calls...)
			if testCase.ExpectedError == nil {
				assert.Equal(t, testCase.NumberOfRequests, len(result))
			} else {
				assert.Equal(t, 0, len(result))
				assert.Equal(t, testCase.ExpectedError.Error(), err.Error())
			}

			s.Stop()
		})
	}
}

func TestRequestValidation(t *testing.T) {
	type testCase struct {
		Name                    string
		Method                  string
		Content                 []byte
		ContentType             string
		ExpectedStatusCode      int
		ExpectedResponseHeaders map[string][]string
		ExpectedMessage         string
	}

	testCases := []testCase{
		{
			Name:               "OPTION request",
			Method:             http.MethodOptions,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"application/json"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "",
		},
		{
			Name:               "GET request",
			Method:             http.MethodGet,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"application/json"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "zkEVM JSON RPC Server",
		},
		{
			Name:               "HEAD request",
			Method:             http.MethodHead,
			ExpectedStatusCode: http.StatusMethodNotAllowed,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "",
		},
		{
			Name:               "PUT request",
			Method:             http.MethodPut,
			ExpectedStatusCode: http.StatusMethodNotAllowed,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "method PUT not allowed\n",
		},
		{
			Name:               "PATCH request",
			Method:             http.MethodPatch,
			ExpectedStatusCode: http.StatusMethodNotAllowed,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "method PATCH not allowed\n",
		},
		{
			Name:               "DELETE request",
			Method:             http.MethodDelete,
			ExpectedStatusCode: http.StatusMethodNotAllowed,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "method DELETE not allowed\n",
		},
		{
			Name:               "CONNECT request",
			Method:             http.MethodConnect,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type": {"text/plain; charset=utf-8"},
			},
			ExpectedMessage: "404 page not found\n",
		},
		{
			Name:               "TRACE request",
			Method:             http.MethodTrace,
			ExpectedStatusCode: http.StatusMethodNotAllowed,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "method TRACE not allowed\n",
		},
		{
			Name:               "Request content bigger than limit",
			Method:             http.MethodPost,
			Content:            make([]byte, maxRequestContentLength+1),
			ExpectedStatusCode: http.StatusRequestEntityTooLarge,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "content length too large (5242881>5242880)\n",
		},
		{
			Name:               "Invalid content type",
			Method:             http.MethodPost,
			ContentType:        "text/html",
			ExpectedStatusCode: http.StatusUnsupportedMediaType,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "invalid content type, only application/json is supported\n",
		},
		{
			Name:               "Empty request body",
			Method:             http.MethodPost,
			ContentType:        contentType,
			Content:            []byte(""),
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "empty request body\n",
		},
		{
			Name:               "Invalid json",
			Method:             http.MethodPost,
			ContentType:        contentType,
			Content:            []byte("this is not a json format string"),
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "invalid json object request body\n",
		},
		{
			Name:               "Incomplete json object",
			Method:             http.MethodPost,
			ContentType:        contentType,
			Content:            []byte("{ \"field\":"),
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "invalid json object request body\n",
		},
		{
			Name:               "Incomplete json array",
			Method:             http.MethodPost,
			ContentType:        contentType,
			Content:            []byte("[ { \"field\":"),
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponseHeaders: map[string][]string{
				"Content-Type":                 {"text/plain; charset=utf-8"},
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"POST, OPTIONS"},
				"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
			},
			ExpectedMessage: "invalid json array request body\n",
		},
	}

	s, _, _ := newSequencerMockedServer(t)
	defer s.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			reqBodyReader := bytes.NewReader(tc.Content)
			httpReq, err := http.NewRequest(tc.Method, s.ServerURL, reqBodyReader)
			require.NoError(t, err)

			httpReq.Header.Add("Content-type", tc.ContentType)

			httpRes, err := http.DefaultClient.Do(httpReq)
			require.NoError(t, err)

			resBody, err := io.ReadAll(httpRes.Body)
			require.NoError(t, err)
			defer httpRes.Body.Close()

			message := string(resBody)
			assert.Equal(t, tc.ExpectedStatusCode, httpRes.StatusCode)
			assert.Equal(t, tc.ExpectedMessage, message)

			for responseHeaderKey, responseHeaderValue := range tc.ExpectedResponseHeaders {
				assert.ElementsMatch(t, httpRes.Header[responseHeaderKey], responseHeaderValue)
			}
		})
	}
}

func TestMaxRequestPerIPPerSec(t *testing.T) {
	// this is the number of requests the test will execute
	// it's important to keep this number with an amount of
	// requests that the machine running this test is able
	// to execute in a single second
	const numberOfRequests = 100
	// the number of workers are the amount of go routines
	// the machine is able to run at the same time without
	// consuming all the resources and making the go routines
	// to affect each other performance, this number may vary
	// depending on the machine spec running the test.
	// a good number to this generally is a number close to
	// the number of cores or threads provided by the CPU.
	const workers = 12
	// it's important to keep this limit smaller than the
	// number of requests the test is going to perform, so
	// the test can have some requests rejected.
	const maxRequestsPerIPAndSecond = 20

	cfg := getSequencerDefaultConfig()
	cfg.MaxRequestsPerIPAndSecond = maxRequestsPerIPAndSecond
	s, m, _ := newMockedServerWithCustomConfig(t, cfg)
	defer s.Stop()

	// since the limitation is made by second,
	// the test waits 1 sec before starting because request are made during the
	// server creation to check its availability. Waiting this second means
	// we have a fresh second without any other request made.
	time.Sleep(time.Second)

	// create a wait group to wait for all the requests to return
	wg := sync.WaitGroup{}
	wg.Add(numberOfRequests)

	// prepare mocks with specific amount of times it can be called
	// this makes us sure the code is calling these methods only for
	// allowed requests
	times := int(cfg.MaxRequestsPerIPAndSecond)
	m.State.On("GetLastL2BlockNumber", context.Background(), nil).Return(uint64(1), nil).Times(times)

	// prepare the workers to process the requests as long as a job is available
	requestsLimitedCount := uint64(0)
	jobs := make(chan int, numberOfRequests)
	// put each worker to work
	for i := 0; i < workers; i++ {
		// each worker works in a go routine to be able to have many
		// workers working concurrently
		go func() {
			// a worker keeps working indefinitely looking for new jobs
			for {
				// waits until a job is available
				<-jobs
				// send the request
				_, err := s.JSONRPCCall("eth_blockNumber")
				// if the request works well or gets rejected due to max requests per sec, it's ok
				// otherwise we stop the test and log the error.
				if err != nil {
					if err.Error() == "429 - You have reached maximum request limit." {
						atomic.AddUint64(&requestsLimitedCount, 1)
					} else {
						require.NoError(t, err)
					}
				}

				// registers in the wait group a request was executed and has returned
				wg.Done()
			}
		}()
	}

	// add jobs to notify workers accordingly to the number
	// of requests the test wants to send to the server
	for i := 0; i < numberOfRequests; i++ {
		jobs <- i
	}

	// wait for all the requests to return
	wg.Wait()

	// checks if all the exceeded requests were limited
	assert.Equal(t, uint64(numberOfRequests-maxRequestsPerIPAndSecond), requestsLimitedCount)

	// wait the server to process the last requests without breaking the
	// connection abruptly
	time.Sleep(time.Second)
}
