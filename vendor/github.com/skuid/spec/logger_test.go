package spec_test

import (
	"testing"

	"github.com/skuid/spec"
	"go.uber.org/zap"
)

func TestNewStandardLogger(t *testing.T) {
	l, err := spec.NewStandardLogger()
	if err != nil {
		t.Errorf("Unexpected error calling NewStandardLogger(): %v", err)
	}
	restore := zap.ReplaceGlobals(l)
	defer restore()
}

func ExampleNewStandardLogger() {
	l, _ := spec.NewStandardLogger() // handle error
	zap.ReplaceGlobals(l)
}
