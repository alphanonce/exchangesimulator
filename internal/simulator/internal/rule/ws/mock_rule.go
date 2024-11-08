// Code generated by mockery v2.45.1. DO NOT EDIT.

package ws

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockRule is an autogenerated mock type for the Rule type
type MockRule struct {
	mock.Mock
}

// Handle provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *MockRule) Handle(_a0 context.Context, _a1 Message, _a2 Connection, _a3 Connection) error {
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

// MatchMessage provides a mock function with given fields: _a0
func (_m *MockRule) MatchMessage(_a0 Message) bool {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for MatchMessage")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(Message) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewMockRule creates a new instance of MockRule. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRule(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRule {
	mock := &MockRule{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
