package go_middlewares

import (
	"net/http"

	"github.com/skuid/spec/middlewares"
	"go.uber.org/zap/zapcore"
)

// Authorizer is an interface for authorizing requests
type Authorizer interface {
	Authorize() middlewares.Middleware
	LoggingClosure(r *http.Request) []zapcore.Field
}
