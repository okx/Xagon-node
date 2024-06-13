package l1_check_block_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block"
	"github.com/stretchr/testify/require"
)

var (
	errGenericToTestAsync       = fmt.Errorf("error_async")
	errReorgToTestAsync         = common.NewReorgError(uint64(1234), fmt.Errorf("fake reorg to test"))
	timeoutContextForAsyncTests = time.Second
)

type mockChecker struct {
	Wg             *sync.WaitGroup
	ErrorsToReturn []error
}

func (m *mockChecker) Step(ctx context.Context) error {
	defer m.Wg.Done()
	err := m.ErrorsToReturn[0]
	if len(m.ErrorsToReturn) > 0 {
		m.ErrorsToReturn = m.ErrorsToReturn[1:]
	}
	return err
}

// If checker.step() returns ok, the async object will relaunch the call
func TestAsyncRelaunchCheckerUntilReorgDetected(t *testing.T) {
	mockChecker := &mockChecker{ErrorsToReturn: []error{nil, nil, errGenericToTestAsync, errReorgToTestAsync}, Wg: &sync.WaitGroup{}}
	sut := l1_check_block.NewAsyncCheck(mockChecker)
	sut.SetPeriodTime(0)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContextForAsyncTests)
	defer cancel()
	mockChecker.Wg.Add(4)

	sut.Run(ctx, nil)

	mockChecker.Wg.Wait()
	result := sut.GetResultBlockingUntilAvailable(0)
	require.NotNil(t, result)
	require.Equal(t, uint64(1234), result.BlockNumber)
	require.Equal(t, true, result.ReorgDetected)
	require.Equal(t, nil, result.Err)
}

func TestAsyncGetResultIsNilUntilStops(t *testing.T) {
	mockChecker := &mockChecker{ErrorsToReturn: []error{nil, nil, errGenericToTestAsync, errReorgToTestAsync}, Wg: &sync.WaitGroup{}}
	sut := l1_check_block.NewAsyncCheck(mockChecker)
	sut.SetPeriodTime(0)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContextForAsyncTests)
	defer cancel()
	mockChecker.Wg.Add(4)
	require.Nil(t, sut.GetResult(), "before start result is Nil")

	sut.Run(ctx, nil)

	require.Nil(t, sut.GetResult(), "after start result is Nil")
	mockChecker.Wg.Wait()
	result := sut.GetResultBlockingUntilAvailable(0)
	require.NotNil(t, result)
}

// RunSynchronous it returns the first result, doesnt mind if a reorg or not
func TestAsyncGRunSynchronousReturnTheFirstResult(t *testing.T) {
	mockChecker := &mockChecker{ErrorsToReturn: []error{errGenericToTestAsync}, Wg: &sync.WaitGroup{}}
	sut := l1_check_block.NewAsyncCheck(mockChecker)
	sut.SetPeriodTime(0)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContextForAsyncTests)
	defer cancel()
	mockChecker.Wg.Add(1)

	result := sut.RunSynchronous(ctx)

	require.NotNil(t, result)
	require.Equal(t, uint64(0), result.BlockNumber)
	require.Equal(t, false, result.ReorgDetected)
	require.Equal(t, errGenericToTestAsync, result.Err)
}

func TestAsyncGRunSynchronousDontAffectGetResult(t *testing.T) {
	mockChecker := &mockChecker{ErrorsToReturn: []error{errGenericToTestAsync}, Wg: &sync.WaitGroup{}}
	sut := l1_check_block.NewAsyncCheck(mockChecker)
	sut.SetPeriodTime(0)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContextForAsyncTests)
	defer cancel()
	mockChecker.Wg.Add(1)

	result := sut.RunSynchronous(ctx)

	require.NotNil(t, result)
	require.Nil(t, sut.GetResult())
}

func TestAsyncStop(t *testing.T) {
	mockChecker := &mockChecker{ErrorsToReturn: []error{nil, nil, errGenericToTestAsync, errReorgToTestAsync}, Wg: &sync.WaitGroup{}}
	sut := l1_check_block.NewAsyncCheck(mockChecker)
	sut.SetPeriodTime(0)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContextForAsyncTests)
	defer cancel()
	require.Nil(t, sut.GetResult(), "before start result is Nil")
	mockChecker.Wg.Add(4)
	sut.Run(ctx, nil)
	sut.Stop()
	sut.Stop()

	result := sut.GetResultBlockingUntilAvailable(0)
	require.NotNil(t, result)
	mockChecker.Wg = &sync.WaitGroup{}
	mockChecker.Wg.Add(4)
	mockChecker.ErrorsToReturn = []error{nil, nil, errGenericToTestAsync, errReorgToTestAsync}
	sut.Run(ctx, nil)
	mockChecker.Wg.Wait()
	result = sut.GetResultBlockingUntilAvailable(0)
	require.NotNil(t, result)
}

func TestAsyncMultipleRun(t *testing.T) {
	mockChecker := &mockChecker{ErrorsToReturn: []error{nil, nil, errGenericToTestAsync, errReorgToTestAsync}, Wg: &sync.WaitGroup{}}
	sut := l1_check_block.NewAsyncCheck(mockChecker)
	sut.SetPeriodTime(0)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutContextForAsyncTests)
	defer cancel()
	require.Nil(t, sut.GetResult(), "before start result is Nil")
	mockChecker.Wg.Add(4)
	sut.Run(ctx, nil)
	sut.Run(ctx, nil)
	sut.Run(ctx, nil)
	result := sut.GetResultBlockingUntilAvailable(0)
	require.NotNil(t, result)
}
