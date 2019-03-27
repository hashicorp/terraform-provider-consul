package consul

import (
	"fmt"
	"os"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

// The providers and configuration used to test permissions
var testAccMasterProviders map[string]terraform.ResourceProvider
var testAccMasterProvider *schema.Provider

const masterToken = "master-token"
const testAccMasterProviderConfiguration = `
provider "consul" {
	token = "` + masterToken + `"
}`

func getMasterClient() (*consulapi.Client, error) {
	rp, err := testAccMasterProviderFactory()
	client := rp.Meta().(*consulapi.Client)
	return client, err
}

func testAccMasterProviderFactory() (*schema.Provider, error) {
	testAccMasterProvider = Provider().(*schema.Provider)
	raw := map[string]interface{}{
		"token": masterToken,
	}
	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		return nil, err
	}

	err = testAccMasterProvider.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		return nil, err
	}
	return testAccMasterProvider, nil
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"consul": testAccProvider,
	}

	testAccMasterProvider, err := testAccMasterProviderFactory()
	if err != nil {
		panic(fmt.Sprintf("err: %s", err))
	}
	testAccMasterProviders = map[string]terraform.ResourceProvider{
		"consul": testAccMasterProvider,
	}
}

func TestResourceProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func TestResourceProvider_Configure(t *testing.T) {
	rp := Provider()

	raw := map[string]interface{}{
		"address":    "demo.consul.io:80",
		"datacenter": "nyc3",
		"scheme":     "https",
	}

	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = rp.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_ConfigureTLS(t *testing.T) {
	rp := Provider()

	raw := map[string]interface{}{
		"address":    "demo.consul.io:80",
		"ca_file":    "test-fixtures/cacert.pem",
		"cert_file":  "test-fixtures/usercert.pem",
		"datacenter": "nyc3",
		"key_file":   "test-fixtures/userkey.pem",
		"scheme":     "https",
	}

	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = rp.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_CAPath(t *testing.T) {
	rp := Provider()

	raw := map[string]interface{}{
		"address": "demo.consul.io:90",
		"ca_path": "test-fixtures/capath",
		"scheme":  "https",
	}

	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = rp.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_ConfigureTLSInsecureHttps(t *testing.T) {
	rp := Provider()

	raw := map[string]interface{}{
		"address":        "demo.consul.io:80",
		"datacenter":     "nyc3",
		"scheme":         "https",
		"insecure_https": true,
	}

	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = rp.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_ConfigureTLSInsecureHttpsMismatch(t *testing.T) {
	rp := Provider()

	raw := map[string]interface{}{
		"address":        "demo.consul.io:80",
		"datacenter":     "nyc3",
		"scheme":         "http",
		"insecure_https": true,
	}

	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = rp.Configure(terraform.NewResourceConfig(rawConfig))
	if err == nil {
		t.Fatal("Provider should error if insecure_https is set but scheme is not https")
	}
}

func TestResourceProvider_tokenIsSensitive(t *testing.T) {
	rp := Provider()

	for _, resource := range rp.Resources() {
		schema, err := rp.GetSchema(&terraform.ProviderSchemaRequest{
			ResourceTypes: []string{resource.Name},
		})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if token, ok := schema.ResourceTypes[resource.Name].Attributes["token"]; ok {
			if !token.Sensitive {
				t.Fatalf("token should be marked as sensitive for %v", resource.Name)
			}
		}
	}

	for _, datasource := range rp.DataSources() {
		schema, err := rp.GetSchema(&terraform.ProviderSchemaRequest{
			DataSources: []string{datasource.Name},
		})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if token, ok := schema.DataSources[datasource.Name].Attributes["token"]; ok {
			if !token.Sensitive {
				t.Fatalf("token should be marked as sensitive for %v", datasource.Name)
			}
		}
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CONSUL_HTTP_ADDR"); v != "" {
		return
	}
	if v := os.Getenv("CONSUL_ADDRESS"); v != "" {
		return
	}
	t.Fatal("Either CONSUL_ADDRESS or CONSUL_HTTP_ADDR must be set for acceptance tests")
}
