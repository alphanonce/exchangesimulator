// Code generated by mockery v2.45.1. DO NOT EDIT.

package ws

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockMessageHandler is an autogenerated mock type for the MessageHandler type
type MockMessageHandler struct {
	mock.Mock
}

// Handle provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *MockMessageHandler) Handle(_a0 context.Context, _a1 Message, _a2 Connection, _a3 Connection) error {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	if len(ret) == 0 {
		panic("no return value specified for Handle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, Message, Connection, Connection) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockMessageHandler creates a new instance of MockMessageHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMessageHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMessageHandler {
	mock := &MockMessageHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
