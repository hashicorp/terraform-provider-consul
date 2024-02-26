package tools

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hashicorp/consul/api"
)

// NewHttpClient returns an http client configured with the given Transport and TLS
// config.
func NewHttpClient(transport *http.Transport, tlsConf api.TLSConfig) (*http.Client, error) {
	client := &http.Client{
		Transport: transport,
	}

	// TODO (slackpad) - Once we get some run time on the HTTP/2 support we
	// should turn it on by default if TLS is enabled. We would basically
	// just need to call http2.ConfigureTransport(transport) here. We also
	// don't want to introduce another external dependency on
	// golang.org/x/net/http2 at this time. For a complete recipe for how
	// to enable HTTP/2 support on a transport suitable for the API client
	// library see agent/http_test.go:TestHTTPServer_H2.

	if transport.TLSClientConfig == nil {
		tlsClientConfig, err := api.SetupTLSConfig(&tlsConf)

		if err != nil {
			return nil, err
		}

		transport.TLSClientConfig = tlsClientConfig
	}

	return client, nil
}

func CloseResponseBody(resp *http.Response) error {
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.Body.Close()
}

// decodeBody is used to JSON decode a body
func DecodeBody(resp *http.Response, out interface{}) error {
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}
