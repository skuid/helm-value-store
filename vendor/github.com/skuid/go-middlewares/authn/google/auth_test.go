package google

import (
	"net/http"
	"testing"

	"go.uber.org/zap/zapcore"
	plus "google.golang.org/api/plus/v1"
)

func TestAuthorizedDomains(t *testing.T) {

	cases := []struct {
		includedDomains []string
	}{
		{
			[]string{"skuid.com", "skuidify.com"},
		},
		{
			[]string{"skuid.com"},
		},
	}

	for _, c := range cases {
		authorizer := New(WithAuthorizedDomains(c.includedDomains...))

		for _, domain := range c.includedDomains {
			if !authorizer.containsDomain(domain) {
				t.Errorf("Expected domain %s to be included! Got false", domain)
			}
		}
	}
}

func TestExtractEmail(t *testing.T) {

	includedDomain := "skuid.com"
	cases := []struct {
		person *plus.Person
		want   string
	}{
		{
			&plus.Person{Emails: []*plus.PersonEmails{&plus.PersonEmails{Value: "micah@skuid.com"}}},
			"micah@skuid.com",
		},
		{
			&plus.Person{Emails: []*plus.PersonEmails{
				&plus.PersonEmails{Type: "home", Value: "micah@example.com"},
				&plus.PersonEmails{Type: "work", Value: "micah@skuid.com"},
			}},
			"micah@skuid.com",
		},
		{
			&plus.Person{Emails: []*plus.PersonEmails{&plus.PersonEmails{Type: "home", Value: "micah@example.com"}}},
			"",
		},
	}

	authorizer := New(WithAuthorizedDomains(includedDomain))
	for _, c := range cases {
		if got := authorizer.extractEmail(c.person); c.want != got {
			t.Errorf("Expected email '%s', got '%s'", c.want, got)
		}
	}
}

func TestLoggingClosure(t *testing.T) {

	cases := []struct {
		tokenMapSeed       map[string]string
		authorizationValue string
		want               *zapcore.Field
	}{
		{
			map[string]string{"Bearer abc123": "micah@skuid.com"},
			"Bearer abc123",
			&zapcore.Field{Key: "user", String: "micah@skuid.com"},
		},
		{
			map[string]string{"Bearer abc123": "micah@skuid.com"},
			"Bearer 111111",
			nil,
		},
		{
			map[string]string{"Bearer abc123": "micah@skuid.com"},
			"",
			nil,
		},
	}

	for _, c := range cases {
		authorizer := New()
		for k, v := range c.tokenMapSeed {
			authorizer.tokenMap.Store(k, v)
		}

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		if c.authorizationValue != "" {
			req.Header.Set("Authorization", c.authorizationValue)
		}

		got := authorizer.LoggingClosure(req)
		if len(got) == 0 {
			if c.want != nil {
				t.Errorf("Got 0 loggable fields, expected %#v", c.want)
			}
			continue
		}
		if firstField := got[0]; firstField.Key != c.want.Key && firstField.String != c.want.String {
			t.Errorf("Expected fields '%#v', got '%#v'", c.want, got)
		}
	}
}
