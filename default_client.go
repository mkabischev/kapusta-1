package kapusta

import (
	"net/http"

	"golang.org/x/net/context"
)

// DefaultClient is adapter for default http.Client
type HTTPClientAdapter struct {
	client *http.Client
}

// Do Client interface implementation
func (c *HTTPClientAdapter) Do(ctx context.Context, r *http.Request) (*http.Response, error) {
	return c.client.Do(r)
}

// NewDefaultClient returns new instance of DefaultClient
func NewHTTPClientAdapter(c *http.Client) Client {
	return &HTTPClientAdapter{c}
}
