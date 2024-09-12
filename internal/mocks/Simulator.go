// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	time "time"

	types "alphanonce.com/exchangesimulator/internal/types"
)

// Simulator is an autogenerated mock type for the Simulator type
type Simulator struct {
	mock.Mock
}

// Process provides a mock function with given fields: request, startTime
func (_m *Simulator) Process(request types.Request, startTime time.Time) (types.Response, time.Time) {
	ret := _m.Called(request, startTime)

	if len(ret) == 0 {
		panic("no return value specified for Process")
	}

	var r0 types.Response
	var r1 time.Time
	if rf, ok := ret.Get(0).(func(types.Request, time.Time) (types.Response, time.Time)); ok {
		return rf(request, startTime)
	}
	if rf, ok := ret.Get(0).(func(types.Request, time.Time) types.Response); ok {
		r0 = rf(request, startTime)
	} else {
		r0 = ret.Get(0).(types.Response)
	}

	if rf, ok := ret.Get(1).(func(types.Request, time.Time) time.Time); ok {
		r1 = rf(request, startTime)
	} else {
		r1 = ret.Get(1).(time.Time)
	}

	return r0, r1
}

// NewSimulator creates a new instance of Simulator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSimulator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Simulator {
	mock := &Simulator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
