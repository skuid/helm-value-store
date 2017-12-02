package middlewares

import (
	"net/http"
	"strings"

	"google.golang.org/grpc"
)

// HandleGrpc passes gRPC requests to the given *grpc.Server and passes http
// requests on to the http handler
func HandleGrpc(server *grpc.Server) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.HasPrefix(
				r.Header.Get("Content-Type"), "application/grpc") {
				server.ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}
		})

	}
}
