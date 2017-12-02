/*
Package example contains an example web application that uses packages
contained in this project.
*/
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/skuid/spec"
	"github.com/skuid/spec/lifecycle"
	"github.com/skuid/spec/middlewares"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	rand.Seed(int64(time.Now().Second()))
}

// A default handler
func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello!"}`))
}

// A function that returns an error
func barf(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"message": "stop it"}`))
}

// A function that sleeps a variable amount of time
func random(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	x := rand.Int() % 6000
	time.Sleep(time.Duration(int64(x)) * time.Millisecond)

	w.Write([]byte(fmt.Sprintf(`{"slept": "%dms"}`, x)))
}

// flip sets lifecycle.Ready to the inverse of it's current state
func flip(w http.ResponseWriter, r *http.Request) {
	if lifecycle.Ready {
		lifecycle.Ready = false
	} else {
		lifecycle.Ready = true
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"Ready": "%t"}`, lifecycle.Ready)))
}

func main() {

	level := spec.LevelPflagPCommandLine("level", "l", zapcore.InfoLevel, "Log level")
	flag.Parse()

	l, _ := spec.NewStandardLevelLogger(*level) // handle error
	zap.ReplaceGlobals(l)

	zap.L().Debug("Some debug log")
	zap.L().Info("Some info log")
	zap.L().Error("Some error log")
	//zap.L().Fatal("Some fatal log") // Terminates application

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/barf", barf)
	mux.HandleFunc("/flip", flip)
	mux.HandleFunc("/random", random)

	handler := middlewares.Apply(
		mux,
		middlewares.InstrumentRoute(),
		middlewares.Logging(),
		middlewares.AccessControlAllowOrigin("*"),
		middlewares.AddHeaders(map[string]string{"X-Frame-Options": "DENY"}),
	)

	internalMux := http.NewServeMux()
	internalMux.Handle("/", handler)

	// Listen on 3001 for metrics and healthchecks
	go spec.MetricsServer(3001)

	hostPort := ":3000"

	zap.L().Info("Server is starting", zap.String("listen", hostPort))

	server := &http.Server{Addr: hostPort, Handler: internalMux}
	lifecycle.ShutdownOnTerm(server)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		zap.L().Fatal("Error listening", zap.Error(err))
	}
	zap.L().Info("Server gracefully stopped")
}
