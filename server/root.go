package server

import (
	"github.com/skuid/go-middlewares"
	"github.com/skuid/helm-value-store/store"
)

// ApiController stores metadata for the API
type ApiController struct {
	releaseStore store.ReleaseStore
	authorizers  []go_middlewares.Authorizer
	timeout      int64
}

// ControllerOpt is a func that modifies an ApiController
type ControllerOpt func(*ApiController)

// WithAutorizer sets the authorizer for an ApiController
func WithAuthorizers(azs ...go_middlewares.Authorizer) ControllerOpt {
	return func(a *ApiController) {
		a.authorizers = azs
	}
}

// WithTimeout sets the timeout in seconds on an ApiController
func WithTimeout(timeout int64) ControllerOpt {
	return func(a *ApiController) {
		a.timeout = timeout
	}
}

// NewApiController returns a new API controller with a default timeout of 300 seconds
func NewApiController(s store.ReleaseStore, opts ...ControllerOpt) *ApiController {
	response := &ApiController{
		releaseStore: s,
		timeout:      300,
	}
	for _, opt := range opts {
		opt(response)
	}

	return response
}
