// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	p2p "github.com/Finschia/ostracon/p2p"
	mock "github.com/stretchr/testify/mock"
)

// NodeInfo is an autogenerated mock type for the NodeInfo type
type NodeInfo struct {
	mock.Mock
}

// CompatibleWith provides a mock function with given fields: other
func (_m *NodeInfo) CompatibleWith(other p2p.NodeInfo) error {
	ret := _m.Called(other)

	var r0 error
	if rf, ok := ret.Get(0).(func(p2p.NodeInfo) error); ok {
		r0 = rf(other)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *NodeInfo) ID() p2p.ID {
	ret := _m.Called()

	var r0 p2p.ID
	if rf, ok := ret.Get(0).(func() p2p.ID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(p2p.ID)
	}

	return r0
}

// NetAddress provides a mock function with given fields:
func (_m *NodeInfo) NetAddress() (*p2p.NetAddress, error) {
	ret := _m.Called()

	var r0 *p2p.NetAddress
	if rf, ok := ret.Get(0).(func() *p2p.NetAddress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*p2p.NetAddress)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Validate provides a mock function with given fields:
func (_m *NodeInfo) Validate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewNodeInfo interface {
	mock.TestingT
	Cleanup(func())
}

// NewNodeInfo creates a new instance of NodeInfo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNodeInfo(t mockConstructorTestingTNewNodeInfo) *NodeInfo {
	mock := &NodeInfo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
