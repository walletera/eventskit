// Code generated by mockery v2.42.2. DO NOT EDIT.

package messages

import (
	mock "github.com/stretchr/testify/mock"
	messages "github.com/walletera/eventskit/messages"
)

// MockProcessorOpt is an autogenerated mock type for the ProcessorOpt type
type MockProcessorOpt struct {
	mock.Mock
}

type MockProcessorOpt_Expecter struct {
	mock *mock.Mock
}

func (_m *MockProcessorOpt) EXPECT() *MockProcessorOpt_Expecter {
	return &MockProcessorOpt_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: opts
func (_m *MockProcessorOpt) Execute(opts *messages.ProcessorOpts) {
	_m.Called(opts)
}

// MockProcessorOpt_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockProcessorOpt_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - opts *messages.ProcessorOpts
func (_e *MockProcessorOpt_Expecter) Execute(opts interface{}) *MockProcessorOpt_Execute_Call {
	return &MockProcessorOpt_Execute_Call{Call: _e.mock.On("Execute", opts)}
}

func (_c *MockProcessorOpt_Execute_Call) Run(run func(opts *messages.ProcessorOpts)) *MockProcessorOpt_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*messages.ProcessorOpts))
	})
	return _c
}

func (_c *MockProcessorOpt_Execute_Call) Return() *MockProcessorOpt_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockProcessorOpt_Execute_Call) RunAndReturn(run func(*messages.ProcessorOpts)) *MockProcessorOpt_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockProcessorOpt creates a new instance of MockProcessorOpt. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockProcessorOpt(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockProcessorOpt {
	mock := &MockProcessorOpt{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
