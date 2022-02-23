package consul

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Config is configuration defined in the provider block
type Config struct {
	Datacenter    string `mapstructure:"datacenter"`
	Address       string `mapstructure:"address"`
	Scheme        string `mapstructure:"scheme"`
	HttpAuth      string `mapstructure:"http_auth"`
	Token         string `mapstructure:"token"`
	CAFile        string `mapstructure:"ca_file"`
	CAPem         string `mapstructure:"ca_pem"`
	CertFile      string `mapstructure:"cert_file"`
	CertPEM       string `mapstructure:"cert_pem"`
	KeyFile       string `mapstructure:"key_file"`
	KeyPEM        string `mapstructure:"key_pem"`
	CAPath        string `mapstructure:"ca_path"`
	InsecureHttps bool   `mapstructure:"insecure_https"`
	Namespace     string `mapstructure:"namespace"`
	client        *consulapi.Client
	resourceData  *schema.ResourceData
}

// Client returns a new client for accessing consul.
func (c *Config) Client() (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()
	if c.Datacenter != "" {
		config.Datacenter = c.Datacenter
	}
	if c.Address != "" {
		config.Address = c.Address
	}
	if c.Scheme != "" {
		config.Scheme = c.Scheme
	}

	if c.CAFile != "" {
		config.TLSConfig.CAFile = c.CAFile
	}
	if c.CAPem != "" {
		config.TLSConfig.CAPem = []byte(c.CAPem)
	}
	if c.CertFile != "" {
		config.TLSConfig.CertFile = c.CertFile
	}
	if c.CertPEM != "" {
		config.TLSConfig.CertPEM = []byte(c.CertPEM)
	}
	if c.KeyFile != "" {
		config.TLSConfig.KeyFile = c.KeyFile
	}
	if c.KeyPEM != "" {
		config.TLSConfig.KeyPEM = []byte(c.KeyPEM)
	}
	if c.CAPath != "" {
		config.TLSConfig.CAPath = c.CAPath
	}
	if c.InsecureHttps {
		if config.Scheme != "https" {
			return nil, fmt.Errorf("insecure_https is meant to be used when scheme is https")
		}
		config.TLSConfig.InsecureSkipVerify = c.InsecureHttps
	}

	// This is a temporary workaround to add the Content-Type header when
	// needed until the fix is released in the Consul api client.
	config.HttpClient = &http.Client{
		Transport: transport{config.Transport},
	}

	if config.Transport.TLSClientConfig == nil {
		tlsClientConfig, err := consulapi.SetupTLSConfig(&config.TLSConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create http client: %s", err)
		}

		config.Transport.TLSClientConfig = tlsClientConfig
	}

	if c.HttpAuth != "" {
		var username, password string
		if strings.Contains(c.HttpAuth, ":") {
			split := strings.SplitN(c.HttpAuth, ":", 2)
			username = split[0]
			password = split[1]
		} else {
			username = c.HttpAuth
		}
		config.HttpAuth = &consulapi.HttpBasicAuth{Username: username, Password: password}
	}

	if c.Token != "" {
		config.Token = c.Token
	}

	client, err := consulapi.NewClient(config)

	log.Printf("[INFO] Consul Client configured with address: '%s', scheme: '%s', datacenter: '%s'"+
		", insecure_https: '%t'", config.Address, config.Scheme, config.Datacenter, config.TLSConfig.InsecureSkipVerify)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// transport adds the Content-Type header to all requests that might need it
// until we update the API client to a version with
// https://github.com/hashicorp/consul/pull/10204 at which time we will be able
// to remove this hack.
type transport struct {
	http.RoundTripper
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		// This will not be the appropriate Content-Type for the license and
		// snapshot endpoints but this is only temporary and Consul does not
		// actually use the header anyway.
		req.Header.Add("Content-Type", "application/json")
	}
	return t.RoundTripper.RoundTrip(req)
}
