package go_middlewares

import (
	"testing"

	"github.com/skuid/go-middlewares/authn/google"
)

func TestInterface(t *testing.T) {
	cases := []struct {
		a Authorizer
	}{
		{google.New()},
	}
	for _, c := range cases {
		c.a.Authorize()
	}
}
