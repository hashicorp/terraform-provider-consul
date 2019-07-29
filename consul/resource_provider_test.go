package consul

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"consul": testAccProvider,
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

// token is sometime nested inside the object
func checkToken(name string, resource *configschema.Block) error {
	for key, value := range resource.BlockTypes {
		if err := checkToken(fmt.Sprintf("%s.%s", name, key), &value.Block); err != nil {
			return err
		}
	}

	for key, value := range resource.Attributes {
		if (key == "token" || strings.HasSuffix(key, ".token")) && !value.Sensitive {
			return fmt.Errorf("token should be marked as sensitive for %s.%s", name, key)
		}
	}
	return nil
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
		if err = checkToken(resource.Name, schema.ResourceTypes[resource.Name]); err != nil {
			t.Fatal(err)
		}
	}

	for _, datasource := range rp.DataSources() {
		schema, err := rp.GetSchema(&terraform.ProviderSchemaRequest{
			DataSources: []string{datasource.Name},
		})
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if err = checkToken(datasource.Name, schema.DataSources[datasource.Name]); err != nil {
			t.Fatal(err)
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
