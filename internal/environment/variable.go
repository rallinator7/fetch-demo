package environment

import (
	"fmt"
	"os"
)

type MissingEnvironmentVariableError struct {
	Variable string
}

func (e *MissingEnvironmentVariableError) Error() string {
	return fmt.Sprintf("environment variable %s is requried but isn't set", e.Variable)
}

type Logger interface {
	Info(...interface{})
	Fatal(...interface{})
}

// GetVariable attempts to retrieve an environment variable. It returns an error if the variable doesn't exist.
func GetVariable(key string) (string, error) {
	value, exist := os.LookupEnv(key)
	if !exist {
		return "", fmt.Errorf("GetVariable: %w", &MissingEnvironmentVariableError{
			Variable: key,
		})
	}

	return value, nil
}

// GetVariableFatal attempts to retrieve an environment variable. It throws a fatal log message if it doesn't exist.
func GetVariableFatal(key string, log Logger) string {
	value, err := GetVariable(key)
	if err != nil {
		log.Fatal(
			err.Error(),
		)
	}

	return value
}

// GetVariableInfo attempts to retrieve an environment variable. It throws an info log message if it doesn't exist and returns the empty string.
func GetVariableInfo(key string, log Logger) string {
	value, err := GetVariable(key)
	if err != nil {
		log.Info(
			fmt.Sprintf("%s is not set - setting to default", key),
		)
	}

	return value
}
