package client

import (
	"fmt"
	"net/http"
	"net/url"
	"pacttest/model"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/stretchr/testify/assert"
)

func TestClientContract(t *testing.T) {
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "ThisConsumer",
		Provider: "ThisProvider",
	})
	assert.NoError(t, err)
	err = mockProvider.
		AddInteraction().Given("User with given id exist").
		UponReceiving("A request for user 10").
		WithRequest(http.MethodGet, "/users/10").
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			// b.JSONBody()
			b.JSONBody(&model.User{
				FirstName: "Some",
				LastName:  "User",
				ID:        10,
				Type:      "someType",
				Username:  "someUser",
			})
		}).ExecuteTest(
		t,
		func(config consumer.MockServerConfig) error {

			uri, err := url.Parse(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			assert.NoError(t, err)
			t.Log("config host : ", config.Port)
			t.Log(uri)
			// Act: test our API client behaves correctly
			// Initialise the API client and point it at the Pact mock server
			client := &Client{
				BaseURL: uri,
			}

			// Execute the API client
			product, err := client.GetUser(10)

			// Assert: check the result
			assert.NoError(t, err)
			assert.Equal(t, 10, product.ID)

			return err
		},
	)

	assert.NoError(t, err)

}
