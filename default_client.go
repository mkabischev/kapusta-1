package kapusta

import (
	"net/http"

	"golang.org/x/net/context"
)

// DefaultClient is adapter for default http.Client
type DefaultClient struct {
	client http.Client
}

// Do Client interface implementation
func (c *DefaultClient) Do(ctx context.Context, r *http.Request) (*http.Response, error) {
	return c.client.Do(r)
}

// NewDefaultClient returns new instance of DefaultClient
func NewDefaultClient(c http.Client) Client {
	return &DefaultClient{c}
}
