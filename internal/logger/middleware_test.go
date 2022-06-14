//go:build unit

package logger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rallinator7/fetch-demo/internal/logger"
	"github.com/rallinator7/fetch-demo/internal/logger/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	tests := map[string]struct {
		Status int
		Path   string
		Method string
		Panic  bool
	}{
		"correctly logs to output on successful call": {
			Status: 200,
			Path:   "/get",
			Method: "GET",
			Panic:  false,
		},
		"correctly logs to output on panic": {
			Status: 200,
			Path:   "/get",
			Method: "GET",
			Panic:  true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)
			logs := &mocks.Logger{}

			if test.Panic {
				logs.On(
					"Errorw",
					"panic recovered",
					"err", mock.Anything,
					"stack", mock.Anything,
				).Return()
			} else {
				logs.On(
					"Infow",
					"request",
					"code-status", test.Status,
					"method", test.Method,
					"remote-address", mock.Anything,
					"path", test.Path,
					"user-agent", mock.Anything,
				).Return()
			}

			fakeHandle := func(w http.ResponseWriter, r *http.Request) {
				if test.Panic {
					panic("ouch")
				} else {
					w.WriteHeader(test.Status)
				}
			}

			request, err := http.NewRequest("GET", "/get", nil)
			require.NoError(err)
			recorder := httptest.NewRecorder()

			logMiddleware := logger.Middleware(logs)
			hfunc := logMiddleware(http.HandlerFunc(fakeHandle))
			hfunc.ServeHTTP(recorder, request)

			if test.Panic {
				assert.Equal(500, recorder.Code)
			} else {
				assert.Equal(test.Status, recorder.Code)
			}

			logs.AssertExpectations(t)
		})
	}
}
