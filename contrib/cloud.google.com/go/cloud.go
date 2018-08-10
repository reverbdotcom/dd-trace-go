// Package cloud provides functions to trace the cloud.google.com/go package.
package cloud // import "gopkg.in/DataDog/dd-trace-go.v1/contrib/cloud.google.com/go"

import (
	"net/http"

	htransport "google.golang.org/api/transport/http"
	apitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/api"
)

// NewClient creates a new http.Client that's configured to trace API requests
// for Google APIs.
func NewClient(options ...Option) (*http.Client, error) {
	cfg := newConfig(options...)
	client, _, err := htransport.NewClient(cfg.ctx)
	if err != nil {
		return nil, err
	}
	client.Transport = apitrace.WrapRoundTripper(client.Transport,
		apitrace.WithContext(cfg.ctx),
		apitrace.WithServiceName(cfg.serviceName))
	return client, nil
}
