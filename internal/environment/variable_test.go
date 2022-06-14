//go:build unit

package environment_test

import (
	"os"
	"testing"

	"github.com/rallinator7/fetch-demo/internal/environment"
	"github.com/rallinator7/fetch-demo/internal/environment/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetVariable(t *testing.T) {
	tests := map[string]struct {
		Key   string
		Value string
		Err   error
	}{
		"returns correct environment variable value": {
			Key:   "key",
			Value: "value",
			Err:   nil,
		},
		"returns error when variable isn't set": {
			Key:   "",
			Value: "",
			Err:   &environment.MissingEnvironmentVariableError{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			if test.Key != "" {
				os.Setenv(test.Key, test.Value)
			}

			env, err := environment.GetVariable(test.Key)
			if test.Err != nil {
				errorType := test.Err
				assert.ErrorAs(err, &errorType)
				return
			}

			assert.Equal(test.Value, env)
		})
	}
}

func TestGetVariableFatal(t *testing.T) {
	tests := map[string]struct {
		Key   string
		Value string
		Err   error
	}{
		"throws a fatal error when variable isn't set": {
			Key:   "test",
			Value: "test",
			Err:   &environment.MissingEnvironmentVariableError{},
		},
		"returns the value of the environment variable": {
			Key:   "test",
			Value: "test",
			Err:   nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			logger := &mocks.Logger{}

			if test.Err == nil {
				os.Setenv(test.Key, test.Value)
			} else {
				logger.On("Fatal", mock.AnythingOfType("string")).Return()
			}

			value := environment.GetVariableFatal(test.Key, logger)

			if test.Err == nil {
				assert.Equal(test.Value, value)
			} else {
				logger.AssertExpectations(t)
				logger.AssertNumberOfCalls(t, "Fatal", 1)
			}

			os.Unsetenv(test.Key)
		})
	}
}

func TestGetVariableInfo(t *testing.T) {
	tests := map[string]struct {
		Key   string
		Value string
		Err   error
	}{
		"throws an info message when variable isn't set": {
			Key:   "test",
			Value: "test",
			Err:   &environment.MissingEnvironmentVariableError{},
		},
		"returns the value of the environment variable": {
			Key:   "test",
			Value: "test",
			Err:   nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			logger := &mocks.Logger{}

			if test.Err == nil {
				os.Setenv(test.Key, test.Value)
			} else {
				logger.On("Info", mock.AnythingOfType("string")).Return()
			}

			value := environment.GetVariableInfo(test.Key, logger)

			if test.Err == nil {
				assert.Equal(test.Value, value)
			} else {
				logger.AssertExpectations(t)
				logger.AssertNumberOfCalls(t, "Info", 1)
			}

			os.Unsetenv(test.Key)
		})
	}
}
