/*

Helpful links to read up on go middlewares:

* https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
* https://justinas.org/writing-http-middleware-in-go/
* https://www.youtube.com/watch?v=xyDkyFjzFVc
*/

package middlewares

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Middleware is a type for decorating requests.
type Middleware func(http.Handler) http.Handler

// Apply wraps a list of middlewares around a handler and returns it
func Apply(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, adapter := range middlewares {
		h = adapter(h)
	}
	return h
}

// AccessControlAllowOrigin is a middleware for adding an access control header to requests
func AccessControlAllowOrigin(origin string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			h.ServeHTTP(w, r)
		})
	}
}

// AddHeaders is a middleware for adding arbitrary headers
func AddHeaders(headers map[string]string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range headers {
				w.Header().Set(k, v)
			}
			h.ServeHTTP(w, r)
		})
	}
}

type statusLoggingResponseWriter struct {
	http.ResponseWriter
	status    int
	bodyBytes int
}

func (w *statusLoggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
func (w *statusLoggingResponseWriter) Write(data []byte) (int, error) {
	length, err := w.ResponseWriter.Write(data)
	w.bodyBytes += length
	return length, err
}

func stripPort(remoteAddr string) string {
	splitIndex := strings.LastIndex(remoteAddr, ":")
	if splitIndex > 0 {
		strippedAddress := remoteAddr[0:splitIndex]
		return strippedAddress
	}
	return remoteAddr
}

func getRemoteAddr(r *http.Request) string {
	address := r.Header.Get("X-Real-IP")
	if len(address) > 0 {
		return stripPort(address)
	}
	address = r.Header.Get("X-Forwarded-For")
	if len(address) > 0 {
		return stripPort(address)
	}
	return stripPort(r.RemoteAddr)
}

// Logging is a middleware for adding a request log. Logs contains the following
// fields: level, timestamp, path, method, response_time, status, message, query,
// remote_addr, and user_agent.
//
// Logging accepts an optional list of closures that accept the inoming request
// and return a slice of zapcore.Field. Each closure is evaluated and its response
// fields are appended to the logged message after the request is handled
func Logging(closures ...func(*http.Request) []zapcore.Field) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := &statusLoggingResponseWriter{w, http.StatusOK, 0}

			defer func() {
				fields := []zapcore.Field{
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.Int("status", wrappedWriter.status),
					zap.String("query", r.Form.Encode()),
					zap.String("remote_addr", getRemoteAddr(r)),
					zap.String("user_agent", r.Header.Get("User-Agent")),
					zap.Int("body_bytes", wrappedWriter.bodyBytes),
				}

				for _, f := range closures {
					fields = append(fields, f(r)...)
				}
				zap.L().Info("", fields...)
			}()

			err := r.ParseForm()
			if err != nil {
				zap.L().Error("Error parsing form", zap.Error(err))
				http.Error(w, `{"error": "error parsing form"}`, http.StatusBadRequest)
				return
			}

			h.ServeHTTP(wrappedWriter, r)
		})
	}
}
