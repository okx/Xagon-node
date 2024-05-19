package sequencer

// CountReadyTx provides a mock function with given fields:
func (_m *WorkerMock) CountReadyTx() uint64 {
	_m.Called()
	ret := _m.Called()
	return ret.Get(0).(uint64)
}
