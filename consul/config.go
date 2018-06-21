package consul

import (
	"fmt"
	"log"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

// Config is configuration defined in the provider block
type Config struct {
	Datacenter    string `mapstructure:"datacenter"`
	Address       string `mapstructure:"address"`
	Scheme        string `mapstructure:"scheme"`
	HttpAuth      string `mapstructure:"http_auth"`
	Token         string `mapstructure:"token"`
	CAFile        string `mapstructure:"ca_file"`
	CertFile      string `mapstructure:"cert_file"`
	KeyFile       string `mapstructure:"key_file"`
	InsecureHttps bool   `mapstructure:"insecure_https"`
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

	tlsConfig := consulapi.TLSConfig{}
	tlsConfig.CAFile = c.CAFile
	tlsConfig.CertFile = c.CertFile
	tlsConfig.KeyFile = c.KeyFile
	if c.InsecureHttps {
		if config.Scheme != "https" {
			return nil, fmt.Errorf("insecure_https is meant to be used when scheme is https")
		}
		tlsConfig.InsecureSkipVerify = c.InsecureHttps
	}

	var err error
	config.HttpClient, err = consulapi.NewHttpClient(config.Transport, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create http client: %s", err)
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
		", insecure_https: '%t'", config.Address, config.Scheme, config.Datacenter, tlsConfig.InsecureSkipVerify)
	if err != nil {
		return nil, err
	}
	return client, nil
}
