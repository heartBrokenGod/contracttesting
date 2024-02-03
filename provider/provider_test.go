package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	l "log"

	"github.com/pact-foundation/pact-go/v2/log"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/pact-foundation/pact-go/v2/version"
	"github.com/stretchr/testify/assert"
)

// var dir, _ = os.Getwd()

var pactDir = "/home/neel/Learning/go_learning/pacttest/consumer/client/pacts"

var requestFilterCalled = false
var stateHandlerCalled = false

func TestV4HTTPProvider(t *testing.T) {
	log.SetLogLevel("Error")
	version.CheckVersion()

	// Start provider API in the background
	go startServer()

	verifier := provider.NewVerifier()

	// Authorization middleware
	// This is your chance to modify the request before it hits your provider
	// NOTE: this should be used very carefully, as it has the potential to
	// _change_ the contract
	f := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Println("[DEBUG] HOOK request filter")
			requestFilterCalled = true
			r.Header.Add("Authorization", "Bearer 1234-dynamic-value")
			next.ServeHTTP(w, r)
		})
	}

	// Verify the Provider with local Pact Files
	err := verifier.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL:        "http://127.0.0.1:8111",
		Provider:               "ThisProvider",
		ProviderVersion:        "1.0.2",
		BrokerURL:              "http://localhost:9292", // os.Getenv("PACT_BROKER_BASE_URL"),
		DisableSSLVerification: true,
		ProviderBranch:         "main",
		BrokerUsername:         "pactbrokeruser",
		BrokerPassword:         "pactbrokerpassword",
		ConsumerVersionSelectors: []provider.Selector{
			&provider.ConsumerVersionSelector{
				Consumer: "ThisConsumer",
				Branch:   "main",
			},
		},
		PublishVerificationResults: true,
		RequestFilter:              f,
		BeforeEach: func() error {
			l.Println("[DEBUG] HOOK before each")
			return nil
		},
		AfterEach: func() error {
			l.Println("[DEBUG] HOOK after each")
			return nil
		},
		StateHandlers: models.StateHandlers{
			"User with given id exist": func(setup bool, s models.ProviderState) (models.ProviderStateResponse, error) {
				stateHandlerCalled = true

				if setup {
					l.Println("[DEBUG] HOOK calling user foo exists state handler", s)
				} else {
					l.Println("[DEBUG] HOOK teardown the 'User foo exists' state")
				}

				// ... do something, such as create "foo" in the database

				// Optionally (if there are generators in the pact) return provider state values to be used in the verification
				return models.ProviderStateResponse{"uuid": "1234"}, nil
			},
		},
	})

	// t.Log(err)
	assert.NoError(t, err)
	assert.True(t, requestFilterCalled)
	assert.True(t, stateHandlerCalled)
}

func startServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/users/10", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		user := &User{
			FirstName: "Some",
			LastName:  "User",
			ID:        10,
			Type:      "someType",
			Username:  "someUser",
		}
		userJson, _ := json.Marshal(user)
		fmt.Fprint(w, string(userJson))
	})

	l.Fatal(http.ListenAndServe("127.0.0.1:8111", mux))
}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Username  string `json:"Username"`
}
