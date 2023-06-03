// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	network "gitub.com/sriramr98/go_kvdb/core/network"
)

// Processor is an autogenerated mock type for the Processor type
type Processor struct {
	mock.Mock
}

// Accept provides a mock function with given fields:
func (_m *Processor) Accept() (network.Conn, error) {
	ret := _m.Called()

	var r0 network.Conn
	var r1 error
	if rf, ok := ret.Get(0).(func() (network.Conn, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() network.Conn); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(network.Conn)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *Processor) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewProcessor interface {
	mock.TestingT
	Cleanup(func())
}

// NewProcessor creates a new instance of Processor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProcessor(t mockConstructorTestingTNewProcessor) *Processor {
	mock := &Processor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
