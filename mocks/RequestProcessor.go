// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	protocol "gitub.com/sriramr98/go_kvdb/core/protocol"
)

// RequestProcessor is an autogenerated mock type for the RequestProcessor type
type RequestProcessor struct {
	mock.Mock
}

// Process provides a mock function with given fields: request
func (_m *RequestProcessor) Process(request protocol.Request) (protocol.Response, error) {
	ret := _m.Called(request)

	var r0 protocol.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(protocol.Request) (protocol.Response, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(protocol.Request) protocol.Response); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Get(0).(protocol.Response)
	}

	if rf, ok := ret.Get(1).(func(protocol.Request) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRequestProcessor interface {
	mock.TestingT
	Cleanup(func())
}

// NewRequestProcessor creates a new instance of RequestProcessor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRequestProcessor(t mockConstructorTestingTNewRequestProcessor) *RequestProcessor {
	mock := &RequestProcessor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
