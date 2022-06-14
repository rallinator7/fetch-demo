package logger

import (
	"net/http"
	"runtime/debug"
)

type Logger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

// Middleware logs the incoming HTTP request & its duration.
func Middleware(l Logger) func(http.Handler) http.Handler {
	wrapperFn := func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			// in case we panic in our go routine we can still capture the log
			defer func() {
				err := recover()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					l.Errorw(
						"panic recovered",
						"err", err,
						"stack", debug.Stack(),
					)
				}
			}()

			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			l.Infow(
				"request",
				"code-status", wrapped.status,
				"method", r.Method,
				"remote-address", r.RemoteAddr,
				"path", r.URL.EscapedPath(),
				"user-agent", r.UserAgent(),
			)
		}

		return http.HandlerFunc(fn)
	}

	return wrapperFn
}
