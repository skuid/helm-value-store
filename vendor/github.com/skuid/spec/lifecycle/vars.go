/*
Package lifecycle provides variables for controlling an application's
lifecycle, and a function for gracefully shutting down an http.Server.

It exposes these variables with two HTTP handlers:

	/live
	/ready

The package is sometimes only imported for the side effect of registering its
HTTP handlers. To use it this way, link this package into your program:

	import _ "github.com/skuid/spec/lifecycle"

When not using the default multiplexer in the http.ListenAndServe function
call, the handlers are available for adding separately.
*/
package lifecycle

import (
	"net/http"
	"os"
	"os/signal"
)

// Shutdown is a boolean that represents whether the application has received
// a SIGTERM.
var Shutdown = false

// Ready is a boolean that represents whether the application is ready.
var Ready = true

// ShutdownTimer is a configuration option for this package that sets the
// amount of time in seconds an application should wait before exiting
// after receiving a SIGTERM.
var ShutdownTimer int64 = 15

// LivenessHandler reports on the status of Shutdown
func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if Shutdown {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "shutdown"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}

// ReadinessHandler reports on the status of Ready
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if !Ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "not ready"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ready"}`))
}

func init() {

	termChan := make(chan os.Signal)
	signal.Notify(termChan, term)

	go func() {
		for range termChan {
			Ready = false
			Shutdown = true
		}
	}()

	http.Handle("/live", http.HandlerFunc(LivenessHandler))
	http.Handle("/ready", http.HandlerFunc(ReadinessHandler))
}
