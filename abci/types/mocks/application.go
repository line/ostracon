// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	abcitypes "github.com/Finschia/ostracon/abci/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/tendermint/tendermint/abci/types"
)

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

// ApplySnapshotChunk provides a mock function with given fields: _a0
func (_m *Application) ApplySnapshotChunk(_a0 types.RequestApplySnapshotChunk) types.ResponseApplySnapshotChunk {
	ret := _m.Called(_a0)

	var r0 types.ResponseApplySnapshotChunk
	if rf, ok := ret.Get(0).(func(types.RequestApplySnapshotChunk) types.ResponseApplySnapshotChunk); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseApplySnapshotChunk)
	}

	return r0
}

// BeginBlock provides a mock function with given fields: _a0
func (_m *Application) BeginBlock(_a0 abcitypes.RequestBeginBlock) types.ResponseBeginBlock {
	ret := _m.Called(_a0)

	var r0 types.ResponseBeginBlock
	if rf, ok := ret.Get(0).(func(abcitypes.RequestBeginBlock) types.ResponseBeginBlock); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseBeginBlock)
	}

	return r0
}

// BeginRecheckTx provides a mock function with given fields: _a0
func (_m *Application) BeginRecheckTx(_a0 abcitypes.RequestBeginRecheckTx) abcitypes.ResponseBeginRecheckTx {
	ret := _m.Called(_a0)

	var r0 abcitypes.ResponseBeginRecheckTx
	if rf, ok := ret.Get(0).(func(abcitypes.RequestBeginRecheckTx) abcitypes.ResponseBeginRecheckTx); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(abcitypes.ResponseBeginRecheckTx)
	}

	return r0
}

// CheckTxAsync provides a mock function with given fields: _a0, _a1
func (_m *Application) CheckTxAsync(_a0 types.RequestCheckTx, _a1 abcitypes.CheckTxCallback) {
	_m.Called(_a0, _a1)
}

// CheckTxSync provides a mock function with given fields: _a0
func (_m *Application) CheckTxSync(_a0 types.RequestCheckTx) abcitypes.ResponseCheckTx {
	ret := _m.Called(_a0)

	var r0 abcitypes.ResponseCheckTx
	if rf, ok := ret.Get(0).(func(types.RequestCheckTx) abcitypes.ResponseCheckTx); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(abcitypes.ResponseCheckTx)
	}

	return r0
}

// Commit provides a mock function with given fields:
func (_m *Application) Commit() types.ResponseCommit {
	ret := _m.Called()

	var r0 types.ResponseCommit
	if rf, ok := ret.Get(0).(func() types.ResponseCommit); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(types.ResponseCommit)
	}

	return r0
}

// DeliverTx provides a mock function with given fields: _a0
func (_m *Application) DeliverTx(_a0 types.RequestDeliverTx) types.ResponseDeliverTx {
	ret := _m.Called(_a0)

	var r0 types.ResponseDeliverTx
	if rf, ok := ret.Get(0).(func(types.RequestDeliverTx) types.ResponseDeliverTx); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseDeliverTx)
	}

	return r0
}

// EndBlock provides a mock function with given fields: _a0
func (_m *Application) EndBlock(_a0 types.RequestEndBlock) types.ResponseEndBlock {
	ret := _m.Called(_a0)

	var r0 types.ResponseEndBlock
	if rf, ok := ret.Get(0).(func(types.RequestEndBlock) types.ResponseEndBlock); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseEndBlock)
	}

	return r0
}

// EndRecheckTx provides a mock function with given fields: _a0
func (_m *Application) EndRecheckTx(_a0 abcitypes.RequestEndRecheckTx) abcitypes.ResponseEndRecheckTx {
	ret := _m.Called(_a0)

	var r0 abcitypes.ResponseEndRecheckTx
	if rf, ok := ret.Get(0).(func(abcitypes.RequestEndRecheckTx) abcitypes.ResponseEndRecheckTx); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(abcitypes.ResponseEndRecheckTx)
	}

	return r0
}

// Info provides a mock function with given fields: _a0
func (_m *Application) Info(_a0 types.RequestInfo) types.ResponseInfo {
	ret := _m.Called(_a0)

	var r0 types.ResponseInfo
	if rf, ok := ret.Get(0).(func(types.RequestInfo) types.ResponseInfo); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseInfo)
	}

	return r0
}

// InitChain provides a mock function with given fields: _a0
func (_m *Application) InitChain(_a0 types.RequestInitChain) types.ResponseInitChain {
	ret := _m.Called(_a0)

	var r0 types.ResponseInitChain
	if rf, ok := ret.Get(0).(func(types.RequestInitChain) types.ResponseInitChain); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseInitChain)
	}

	return r0
}

// ListSnapshots provides a mock function with given fields: _a0
func (_m *Application) ListSnapshots(_a0 types.RequestListSnapshots) types.ResponseListSnapshots {
	ret := _m.Called(_a0)

	var r0 types.ResponseListSnapshots
	if rf, ok := ret.Get(0).(func(types.RequestListSnapshots) types.ResponseListSnapshots); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseListSnapshots)
	}

	return r0
}

// LoadSnapshotChunk provides a mock function with given fields: _a0
func (_m *Application) LoadSnapshotChunk(_a0 types.RequestLoadSnapshotChunk) types.ResponseLoadSnapshotChunk {
	ret := _m.Called(_a0)

	var r0 types.ResponseLoadSnapshotChunk
	if rf, ok := ret.Get(0).(func(types.RequestLoadSnapshotChunk) types.ResponseLoadSnapshotChunk); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseLoadSnapshotChunk)
	}

	return r0
}

// OfferSnapshot provides a mock function with given fields: _a0
func (_m *Application) OfferSnapshot(_a0 types.RequestOfferSnapshot) types.ResponseOfferSnapshot {
	ret := _m.Called(_a0)

	var r0 types.ResponseOfferSnapshot
	if rf, ok := ret.Get(0).(func(types.RequestOfferSnapshot) types.ResponseOfferSnapshot); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseOfferSnapshot)
	}

	return r0
}

// Query provides a mock function with given fields: _a0
func (_m *Application) Query(_a0 types.RequestQuery) types.ResponseQuery {
	ret := _m.Called(_a0)

	var r0 types.ResponseQuery
	if rf, ok := ret.Get(0).(func(types.RequestQuery) types.ResponseQuery); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseQuery)
	}

	return r0
}

// SetOption provides a mock function with given fields: _a0
func (_m *Application) SetOption(_a0 types.RequestSetOption) types.ResponseSetOption {
	ret := _m.Called(_a0)

	var r0 types.ResponseSetOption
	if rf, ok := ret.Get(0).(func(types.RequestSetOption) types.ResponseSetOption); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(types.ResponseSetOption)
	}

	return r0
}

// NewApplication creates a new instance of Application. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplication(t interface {
	mock.TestingT
	Cleanup(func())
}) *Application {
	mock := &Application{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
