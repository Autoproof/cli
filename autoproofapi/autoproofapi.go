package autoproofapi

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

const (
	defaultBaseURL = "https://app.autoproof.dev/api"
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
}

type ClientOption func(*Client)

func WithCustomClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithCustomBaseURL(baseURL *url.URL) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func NewClient(opts ...ClientOption) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := Client{
		httpClient: http.DefaultClient,
		baseURL:    baseURL,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

func (c *Client) newRequest(ctx context.Context, method string, urlStr string, body any) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
		q, _ := json.Marshal(body)
		fmt.Println(string(q))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL.JoinPath(urlStr).String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request, v any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		// If we got an error, and the context has been canceled, the context error is probably more useful.
		case <-req.Context().Done():
			return req.Context().Err()
		default:
		}

		return err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return err
	}

	// For 204 No content cases.
	if v == nil {
		return nil
	}

	decErr := json.NewDecoder(resp.Body).Decode(v)
	switch {
	case errors.Is(decErr, io.EOF):
		// Ignore EOF errors (caused by empty response body).
		err = nil
	case decErr != nil:
		err = decErr
	}

	return err
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	apiError := APIError{
		StatusCode: r.StatusCode,
	}
	if err := json.NewDecoder(r.Body).Decode(&apiError.Details); err != nil {
		return errors.Join(apiError, fmt.Errorf("failed to unmarshal error details: %v", err))
	}

	return apiError
}
