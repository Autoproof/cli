package autoproofapi

import (
	"net/http"
)

type APIKeyTransport struct {
	APIKey string

	Transport http.RoundTripper
}

func (t *APIKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.APIKey)
	return t.transport().RoundTrip(req)
}

func (t *APIKeyTransport) Client() *http.Client {
	return &http.Client{
		Transport: t,
	}
}

func (t *APIKeyTransport) transport() http.RoundTripper {
	if t.Transport == nil {
		return http.DefaultTransport
	}
	return t.Transport
}
