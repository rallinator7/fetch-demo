// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

type Logger_Expecter struct {
	mock *mock.Mock
}

func (_m *Logger) EXPECT() *Logger_Expecter {
	return &Logger_Expecter{mock: &_m.Mock}
}

// Fatal provides a mock function with given fields: _a0
func (_m *Logger) Fatal(_a0 ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, _a0...)
	_m.Called(_ca...)
}

// Logger_Fatal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Fatal'
type Logger_Fatal_Call struct {
	*mock.Call
}

// Fatal is a helper method to define mock.On call
//  - _a0 ...interface{}
func (_e *Logger_Expecter) Fatal(_a0 ...interface{}) *Logger_Fatal_Call {
	return &Logger_Fatal_Call{Call: _e.mock.On("Fatal",
		append([]interface{}{}, _a0...)...)}
}

func (_c *Logger_Fatal_Call) Run(run func(_a0 ...interface{})) *Logger_Fatal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Logger_Fatal_Call) Return() *Logger_Fatal_Call {
	_c.Call.Return()
	return _c
}

// Info provides a mock function with given fields: _a0
func (_m *Logger) Info(_a0 ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, _a0...)
	_m.Called(_ca...)
}

// Logger_Info_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Info'
type Logger_Info_Call struct {
	*mock.Call
}

// Info is a helper method to define mock.On call
//  - _a0 ...interface{}
func (_e *Logger_Expecter) Info(_a0 ...interface{}) *Logger_Info_Call {
	return &Logger_Info_Call{Call: _e.mock.On("Info",
		append([]interface{}{}, _a0...)...)}
}

func (_c *Logger_Info_Call) Run(run func(_a0 ...interface{})) *Logger_Info_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Logger_Info_Call) Return() *Logger_Info_Call {
	_c.Call.Return()
	return _c
}

// NewLogger creates a new instance of Logger. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewLogger(t testing.TB) *Logger {
	mock := &Logger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
