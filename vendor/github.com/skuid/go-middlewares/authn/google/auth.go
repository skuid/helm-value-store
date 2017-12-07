package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/skuid/spec/middlewares"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	plus "google.golang.org/api/plus/v1"
)

// Authorizer authorizes a google user against a whitelist of domains
type Authorizer struct {
	authorizedDomains map[string]bool
	tokenMap          sync.Map
}

func (a *Authorizer) domains() []string {
	response := []string{}
	for k := range a.authorizedDomains {
		response = append(response, k)
	}
	return response
}

func (a *Authorizer) containsDomain(domain string) bool {
	_, ok := a.authorizedDomains[domain]
	return ok
}

// WithAuthorizedDomains adds the specified domains to the whitelist
func WithAuthorizedDomains(domains ...string) func(*Authorizer) {
	return func(a *Authorizer) {
		for _, domain := range domains {
			a.authorizedDomains[domain] = true
		}
	}
}

// New returns a new Authorizer with default options and applies any supplied
// option functions
func New(opts ...func(*Authorizer)) *Authorizer {
	a := &Authorizer{
		authorizedDomains: map[string]bool{},
		tokenMap:          sync.Map{},
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func getPerson(ctx context.Context, authorizationValue string) (*plus.Person, error) {
	googleURL := "https://www.googleapis.com/plus/v1/people/me"
	req, err := http.NewRequest(http.MethodGet, googleURL, nil)
	if err != nil {
		zap.L().Error("Error creating request", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Authorization", authorizationValue)

	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		zap.L().Error("Error making request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			zap.L().Error("Error decoding body", zap.Error(err))
			return nil, err
		}

		zap.L().Info(
			"Not authorized against Google",
			zap.Int("code", resp.StatusCode),
			zap.ByteString("body", body),
		)
		return nil, fmt.Errorf("Not authorized")
	}

	person := &plus.Person{}
	err = json.NewDecoder(resp.Body).Decode(person)
	if err != nil {
		zap.L().Error("Error decoding response", zap.Error(err))
		return nil, err
	}
	return person, nil
}

// Check each email and get the first that ends with a valid domain
func (a *Authorizer) extractEmail(person *plus.Person) string {
	for _, email := range person.Emails {
		if idx := strings.LastIndex(email.Value, "@"); idx > 0 {
			domain := email.Value[idx+1:]
			if a.containsDomain(domain) {
				return email.Value
			}
		}
	}
	return ""
}

func (a *Authorizer) authorize(ctx context.Context, token string) bool {
	if token == "" {
		zap.L().Debug("Empty token")
		return false
	}

	var username string

	usernameIface, ok := a.tokenMap.Load(token)

	// User is not in cache
	if !ok {

		// Get the person's info
		person, err := getPerson(ctx, token)
		if err != nil {
			zap.L().Info("Couldn't get person", zap.String("token", token))
			return false
		}

		// Validate the person's domain
		validated := a.containsDomain(person.Domain)
		if !validated {
			zap.L().Info("Invalid domain", zap.String("person", fmt.Sprintf("%#v", person)), zap.Strings("domains", a.domains()))
			return false
		}

		// Get the first matching email
		username = a.extractEmail(person)
		if username == "" {
			zap.L().Info("Couldn't get person's email", zap.String("person", fmt.Sprintf("%#v", person)))
			return false
		}
		zap.L().Debug("Storing valid user", zap.String("user", username))
		a.tokenMap.Store(token, username)
	} else {
		username = usernameIface.(string)
	}
	zap.L().Debug("Successfully authorized", zap.String("user", username))

	return true
}

// LoggingClosure adds a "user" field for an authorized user
func (a *Authorizer) LoggingClosure(r *http.Request) []zapcore.Field {

	token := r.Header.Get("Authorization")
	if token == "" {
		zap.L().Debug("No authorization header on request")
		return []zapcore.Field{}
	}

	usernameIface, ok := a.tokenMap.Load(token)

	if !ok {
		zap.L().Debug("Couldn't exract user from request")
		return []zapcore.Field{}
	}
	username := usernameIface.(string)

	return []zapcore.Field{zap.String("user", username)}
}

// Authorize is a middleware for adding AuthZ to routes
func (a *Authorizer) Authorize() middlewares.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok := a.authorize(r.Context(), r.Header.Get("Authorization"))
			if !ok {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
