package mocks

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// AddInnerTx provides a mock function with given fields: ctx, txHash, innerTx
func (_m *PoolMock) AddInnerTx(ctx context.Context, txHash common.Hash, innerTx []byte) error {
	ret := _m.Called(ctx, txHash, innerTx)

	if len(ret) == 0 {
		panic("no return value specified for AddTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, []byte) error); ok {
		r0 = rf(ctx, txHash, innerTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetInnerTx provides a mock function with given fields: ctx, txHash
func (_m *PoolMock) GetInnerTx(ctx context.Context, txHash common.Hash) (string, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for GetPendingTxs")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (string, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) string); ok {
		r0 = rf(ctx, txHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMinSuggestedGasPriceWithDelta provides a mock function with given fields: ctx, delta
func (_m *PoolMock) GetMinSuggestedGasPriceWithDelta(ctx context.Context, delta time.Duration) (uint64, error) {
	ret := _m.Called(ctx, delta)

	if len(ret) == 0 {
		panic("no return value specified for GetMinSuggestedGasPriceWithDelta")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration) (uint64, error)); ok {
		return rf(ctx, delta)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration) uint64); ok {
		r0 = rf(ctx, delta)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uint64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Duration) error); ok {
		r1 = rf(ctx, delta)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsPendingStatEnabled provides a mock function with given fields: ctx
func (_m *PoolMock) IsPendingStatEnabled(ctx context.Context) bool {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for IsPendingStatEnabled")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(bool)
		}
	}

	return r0
}
