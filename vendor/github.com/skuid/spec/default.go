package spec

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/skuid/spec/lifecycle"
	_ "github.com/skuid/spec/metrics" // import spec metrics
	"go.uber.org/zap"
)

// MetricsServer is a convenience function for handling metrics and healthchecks.
// It listens on the provided port without authentication.
//
//   go spec.MetricsServer(3001)
//
func MetricsServer(port int) {
	internalMux := http.NewServeMux()
	internalMux.Handle("/metrics", promhttp.Handler())
	internalMux.HandleFunc("/live", lifecycle.LivenessHandler)
	internalMux.HandleFunc("/ready", lifecycle.ReadinessHandler)
	hostPort := fmt.Sprintf(":%d", port)

	zap.L().Info("Metrics server is starting", zap.String("listen", hostPort))
	httpServer := &http.Server{Addr: hostPort, Handler: internalMux}
	lifecycle.ShutdownOnTerm(httpServer)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		zap.L().Fatal("Error listening", zap.Error(err))
	}
	zap.L().Info("Metrics server gracefully stopped")
}
