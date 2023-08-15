// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	network "gitub.com/sriramr98/go_kvdb/core/network"
)

// Listener is an autogenerated mock type for the Listener type
type Listener struct {
	mock.Mock
}

// Listen provides a mock function with given fields: _a0, addr
func (_m *Listener) Listen(_a0 string, addr string) (network.Processor, error) {
	ret := _m.Called(_a0, addr)

	var r0 network.Processor
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (network.Processor, error)); ok {
		return rf(_a0, addr)
	}
	if rf, ok := ret.Get(0).(func(string, string) network.Processor); ok {
		r0 = rf(_a0, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.Processor)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, addr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewListener interface {
	mock.TestingT
	Cleanup(func())
}

// NewListener creates a new instance of Listener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewListener(t mockConstructorTestingTNewListener) *Listener {
	mock := &Listener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
