package mystrom

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(
		WithAPIKey("abcd"),
	)

	if c.apiKey != "abcd" {
		t.Error("expected abcd, got", c.apiKey)
	}
}

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/test" {
			t.Error("expected /test, got", req.URL.Path)
		}

		if req.Header.Get("Token") != "abcd" {
			t.Error("expected Authorization: abcd, got", req.Header.Get("Authorization"))
		}

		if req.Header.Get("User-Agent") != "test/0.1" {
			t.Error("expected User-Agent: test/0.1, got", req.Header.Get("User-Agent"))
		}

		_, err := rw.Write([]byte(`OK`))
		if err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()

	c := client(server)
	baseURL, _ := url.Parse(server.URL)

	req, err := c.newRequest(context.Background(), baseURL, http.MethodGet, "test", nil, nil)
	if err != nil {
		t.Error(err)
	}

	res, err := c.do(req)
	if err != nil {
		t.Error(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error("expected", http.StatusOK, "got", res.StatusCode)
	}

	if string(body) != "OK" {
		t.Error("expected OK, got", string(body))
	}
}

func TestDoJSON(t *testing.T) {
	type result struct {
		Result string `json:"result"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		_, err := rw.Write([]byte(`{"result": "success"}`))
		if err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()

	c := client(server)
	baseURL, _ := url.Parse(server.URL)

	req, err := c.newRequest(context.Background(), baseURL, http.MethodGet, "test", nil, nil)
	if err != nil {
		t.Error(err)
	}

	var r result

	_, err = c.doJSON(req, &r)
	if err != nil {
		t.Error(err)
	}

	if r.Result != "success" {
		t.Error("expected success, got", r.Result)
	}
}

func client(server *httptest.Server) Client {
	return Client{
		httpClient: server.Client(),
		userAgent:  "test/0.1",
		apiKey:     "abcd",
	}
}
