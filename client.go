package mystrom

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ErrStatus respresents a non success status code error.
var ErrStatus = errors.New("status code error")

type Client struct {
	userAgent string
	apiKey    string

	httpClient *http.Client
}

func NewClient(opts ...Option) *Client {
	client := Client{
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		opt(&client)
	}

	return &client
}

type Option func(*Client)

// WithAPIKey defines the API key to be used.
func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithUserAgent allows to change the user agent.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithHTTPClient allows to replace the http client.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		c.httpClient = h
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	baseURL *url.URL,
	method, path string,
	params url.Values,
	body interface{}, //nolint:unparam
) (*http.Request, error) {
	if params == nil {
		params = url.Values{}
	}

	rel := &url.URL{Path: path}
	u := baseURL.ResolveReference(rel)
	u.RawQuery = params.Encode()

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)

		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if c.apiKey != "" {
		req.Header.Set("Token", c.apiKey)
	}

	return req, nil
}

func (c *Client) doJSON(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&v)
		if err != nil {
			return nil, err
		}
	}

	return resp, err
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.Body != nil {
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("error reading body: %w", err)
			}
			body = bytes.TrimSpace(body)

			return resp, fmt.Errorf("%s: %d, %w '%s'", http.StatusText(resp.StatusCode), resp.StatusCode, ErrStatus, body)
		}
		return resp, fmt.Errorf("%s: %d, %w", http.StatusText(resp.StatusCode), resp.StatusCode, ErrStatus)
	}

	return resp, err
}
