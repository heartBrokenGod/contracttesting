package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pacttest/model"
)

var (
	ErrUnavailable = errors.New("api unavailable")
)

type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
}

// Gets single user from the api
func (c *Client) GetUser(id int) (*model.User, error) {
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("/users/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var user model.User
	_, err = c.do(req, &user)

	if err != nil {
		return nil, ErrUnavailable
	}

	return &user, err
}

func (c *Client) newRequest(method, path string, body any) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Admin Service")

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
