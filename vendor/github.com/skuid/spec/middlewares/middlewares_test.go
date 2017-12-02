package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestGetRemoteAddr(t *testing.T) {

	cases := []struct {
		name       string
		header     map[string]string
		remoteAddr string
		want       string
	}{
		{
			"getRemoteAddr consumes X-Real-IP",
			map[string]string{
				"X-Real-IP":       "11.22.33.44:1234",
				"X-Forwarded-For": "11.22.33.44:1234",
			},
			"11.22.33.44:1234",
			"11.22.33.44",
		},
		{
			"getRemoteAddr consumes X-Forwarded-For",
			map[string]string{
				"X-Forwarded-For": "11.22.33.44:1234",
			},
			"11.22.33.44:1234",
			"11.22.33.44",
		},
		{
			"getRemoteAddr strips port",
			map[string]string{
				"X-Real-IP":       "11.22.33.44:1234",
				"X-Forwarded-For": "11.22.33.44:1234",
			},
			"11.22.33.44:12312",
			"11.22.33.44",
		},
		{
			"getRemoteAddr handles no header",
			map[string]string{},
			"11.22.33.44:2000",
			"11.22.33.44",
		},
		{
			"getRemoteAddr handles no port",
			map[string]string{},
			"11.22.33.44",
			"11.22.33.44",
		},
		{
			"getRemoteAddr handles ipv6",
			map[string]string{},
			"[::]:1234",
			"[::]",
		},
	}
	for _, c := range cases {
		request, err := http.NewRequest("GET", "http://localhost", nil)
		if err != nil {
			t.Errorf("Failed %s: Unable to create a new http.Request{}", c.name)
		}
		request.RemoteAddr = c.remoteAddr
		for header, value := range c.header {
			request.Header.Set(header, value)
		}
		if got := getRemoteAddr(request); got != c.want {
			t.Errorf("Failed %s: getRemoteAddr() Expected: %v, got: %v", c.name, c.want, got)
		}
	}
}

func TestLoggingClosures(t *testing.T) {

	// Set up the observer and inject it into the logger
	core, observed := observer.New(zapcore.DebugLevel)
	opt := zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return core
	})
	logger := zap.NewExample(opt)
	reset := zap.ReplaceGlobals(logger)
	defer reset()

	handleRequest := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, client")
	}

	closure := func(r *http.Request) []zapcore.Field {
		return []zapcore.Field{zap.String("user", r.Header.Get("user"))}
	}

	server := httptest.NewServer(Logging(closure)(http.HandlerFunc(handleRequest)))
	defer server.Close()

	request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	request.Header.Add("user", "alfanzo")
	client := http.Client{}
	_, err := client.Do(request)
	if err != nil {
		t.Fatalf("Error making request to test server: %s", err.Error())
	}

	logger.Sync()
	if observed.Len() == 0 {
		t.Fatal("Expected log! Got no logs")
	}
	loggedMessage := observed.All()[0]
	if loggedMessage.Context[7].String != "alfanzo" {
		t.Errorf("Didn't find alfanzo")
	}
}
