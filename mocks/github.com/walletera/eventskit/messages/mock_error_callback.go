// Code generated by mockery v2.42.2. DO NOT EDIT.

package messages

import (
	mock "github.com/stretchr/testify/mock"
	werrors "github.com/walletera/werrors"
)

// MockErrorCallback is an autogenerated mock type for the ErrorCallback type
type MockErrorCallback struct {
	mock.Mock
}

type MockErrorCallback_Expecter struct {
	mock *mock.Mock
}

func (_m *MockErrorCallback) EXPECT() *MockErrorCallback_Expecter {
	return &MockErrorCallback_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: processingError
func (_m *MockErrorCallback) Execute(processingError werrors.WError) {
	_m.Called(processingError)
}

// MockErrorCallback_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockErrorCallback_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - processingError werrors.WError
func (_e *MockErrorCallback_Expecter) Execute(processingError interface{}) *MockErrorCallback_Execute_Call {
	return &MockErrorCallback_Execute_Call{Call: _e.mock.On("Execute", processingError)}
}

func (_c *MockErrorCallback_Execute_Call) Run(run func(processingError werrors.WError)) *MockErrorCallback_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(werrors.WError))
	})
	return _c
}

func (_c *MockErrorCallback_Execute_Call) Return() *MockErrorCallback_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockErrorCallback_Execute_Call) RunAndReturn(run func(werrors.WError)) *MockErrorCallback_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockErrorCallback creates a new instance of MockErrorCallback. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockErrorCallback(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockErrorCallback {
	mock := &MockErrorCallback{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
