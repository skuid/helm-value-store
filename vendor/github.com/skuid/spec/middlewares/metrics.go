package middlewares

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "Counter of requests broken out for each verb, path, and response code.",
		},
		[]string{"verb", "path", "code"},
	)
	requestLatencies = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_latencies",
			Help: "Response latency distribution in microseconds for each verb and path",
			// Use buckets ranging from 125 ms to 8 seconds.
			Buckets: prometheus.ExponentialBuckets(125000, 2.0, 7),
		},
		[]string{"verb", "path"},
	)
	requestLatenciesSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration_microseconds",
			Help: "Response latency summary in microseconds for each verb and path.",
			// Make the sliding window of 1h.
			MaxAge: time.Hour,
		},
		[]string{"verb", "path"},
	)
)

func register() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestLatencies)
	prometheus.MustRegister(requestLatenciesSummary)
}

func monitor(verb, path string, httpCode int, reqStart time.Time) {
	elapsed := float64((time.Since(reqStart)) / time.Microsecond)

	requestCounter.WithLabelValues(verb, path, codeToString(httpCode)).Inc()
	requestLatencies.WithLabelValues(verb, path).Observe(elapsed)
	requestLatenciesSummary.WithLabelValues(verb, path).Observe(elapsed)
}

func init() {
	register()
}

// InstrumentRoute is a middleware for adding the following metrics for each
// route:
//
//     # Counter
//     http_request_count{"verb", "path", "code}
//     # Histogram
//     http_request_latencies{"verb", "path"}
//     # Summary
//     http_request_duration_microseconds{"verb", "path", "code}
//
func InstrumentRoute() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			wrappedWriter := &statusLoggingResponseWriter{w, http.StatusOK, 0}

			defer func() {
				monitor(r.Method, r.URL.Path, wrappedWriter.status, now)
			}()
			h.ServeHTTP(wrappedWriter, r)

		})
	}
}
