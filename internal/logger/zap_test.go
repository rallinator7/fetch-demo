//go:build unit

package logger_test

import (
	"testing"

	"github.com/rallinator7/fetch-demo/internal/logger"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	tests := map[string]struct {
		Message string
		Content []interface{}
	}{
		"creates a new logger that can log": {
			Message: "it works!",
			Content: []interface{}{"testKey", "testValue"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)

			sugar, err := logger.New()
			require.NoError(err)

			sugar.Infow(test.Message, test.Content...)
		})
	}
}
