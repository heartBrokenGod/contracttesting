package client

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"pacttest/model"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestClientUnit_GetUser(t *testing.T) {
// 	userID := 10

// 	// setup a mock server
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		assert.Equal(t, r.URL.String(), fmt.Sprintf("/users/%d", userID))
// 		user, _ := json.Marshal(model.User{
// 			FirstName: "Sally",
// 			LastName:  "McDougall",
// 			ID:        userID,
// 			Type:      "admin",
// 			Username:  "smcdougall",
// 		})
// 		w.Write([]byte(user))
// 	}))
// 	defer server.Close()

// 	// setup client
// 	u, _ := url.Parse(server.URL)
// 	client := &Client{
// 		BaseURL: u,
// 	}

// 	user, err := client.GetUser(userID)
// 	assert.NoError(t, err)

// 	// Assert basic fact
// 	assert.Equal(t, user.ID, userID)

// }
