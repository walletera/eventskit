// Code generated by mockery v2.42.2. DO NOT EDIT.

package fake

import (
	context "context"

	errors "github.com/walletera/eventskit/errors"
	fake "github.com/walletera/eventskit/fake"

	mock "github.com/stretchr/testify/mock"
)

// MockEventHandler is an autogenerated mock type for the EventHandler type
type MockEventHandler struct {
	mock.Mock
}

type MockEventHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEventHandler) EXPECT() *MockEventHandler_Expecter {
	return &MockEventHandler_Expecter{mock: &_m.Mock}
}

// HandleFakeEvent provides a mock function with given fields: ctx, e
func (_m *MockEventHandler) HandleFakeEvent(ctx context.Context, e fake.Event) errors.ProcessingError {
	ret := _m.Called(ctx, e)

	if len(ret) == 0 {
		panic("no return value specified for HandleFakeEvent")
	}

	var r0 errors.ProcessingError
	if rf, ok := ret.Get(0).(func(context.Context, fake.Event) errors.ProcessingError); ok {
		r0 = rf(ctx, e)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(errors.ProcessingError)
		}
	}

	return r0
}

// MockEventHandler_HandleFakeEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HandleFakeEvent'
type MockEventHandler_HandleFakeEvent_Call struct {
	*mock.Call
}

// HandleFakeEvent is a helper method to define mock.On call
//   - ctx context.Context
//   - e fake.Event
func (_e *MockEventHandler_Expecter) HandleFakeEvent(ctx interface{}, e interface{}) *MockEventHandler_HandleFakeEvent_Call {
	return &MockEventHandler_HandleFakeEvent_Call{Call: _e.mock.On("HandleFakeEvent", ctx, e)}
}

func (_c *MockEventHandler_HandleFakeEvent_Call) Run(run func(ctx context.Context, e fake.Event)) *MockEventHandler_HandleFakeEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(fake.Event))
	})
	return _c
}

func (_c *MockEventHandler_HandleFakeEvent_Call) Return(_a0 errors.ProcessingError) *MockEventHandler_HandleFakeEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEventHandler_HandleFakeEvent_Call) RunAndReturn(run func(context.Context, fake.Event) errors.ProcessingError) *MockEventHandler_HandleFakeEvent_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEventHandler creates a new instance of MockEventHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventHandler {
	mock := &MockEventHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}