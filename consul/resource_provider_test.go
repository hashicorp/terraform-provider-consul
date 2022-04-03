package consul

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

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

	err := rp.Configure(terraform.NewResourceConfigRaw(raw))
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

	err := rp.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestResourceProvider_ConfigureTLSPem(t *testing.T) {
	rp := Provider()

	caPem, err := ioutil.ReadFile("test-fixtures/cacert.pem")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	certPem, err := ioutil.ReadFile("test-fixtures/usercert.pem")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	keyPem, err := ioutil.ReadFile("test-fixtures/userkey.pem")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	raw := map[string]interface{}{
		"address":    "demo.consul.io:80",
		"ca_pem":     string(caPem),
		"cert_pem":   string(certPem),
		"datacenter": "nyc3",
		"key_pem":    string(keyPem),
		"scheme":     "https",
	}

	err = rp.Configure(terraform.NewResourceConfigRaw(raw))
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

	err := rp.Configure(terraform.NewResourceConfigRaw(raw))
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

	err := rp.Configure(terraform.NewResourceConfigRaw(raw))
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

	err := rp.Configure(terraform.NewResourceConfigRaw(raw))
	if err == nil {
		t.Fatal("Provider should error if insecure_https is set but scheme is not https")
	}
}

// token is sometime nested inside the object
// func checkToken(name string, resource *configschema.Block) error {
// 	for key, value := range resource.BlockTypes {
// 		if err := checkToken(fmt.Sprintf("%s.%s", name, key), &value.Block); err != nil {
// 			return err
// 		}
// 	}

// 	for key, value := range resource.Attributes {
// 		if (key == "token" || strings.HasSuffix(key, ".token")) && !value.Sensitive {
// 			return fmt.Errorf("token should be marked as sensitive for %s.%s", name, key)
// 		}
// 	}
// 	return nil
// }

// func TestResourceProvider_tokenIsSensitive(t *testing.T) {
// 	rp := Provider()

// 	for _, resource := range rp.Resources() {
// 		schema, err := rp.GetSchema(&terraform.ProviderSchemaRequest{
// 			ResourceTypes: []string{resource.Name},
// 		})
// 		if err != nil {
// 			t.Fatalf("err: %v", err)
// 		}
// 		if err = checkToken(resource.Name, schema.ResourceTypes[resource.Name]); err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	for _, datasource := range rp.DataSources() {
// 		schema, err := rp.GetSchema(&terraform.ProviderSchemaRequest{
// 			DataSources: []string{datasource.Name},
// 		})
// 		if err != nil {
// 			t.Fatalf("err: %v", err)
// 		}

// 		if err = checkToken(datasource.Name, schema.DataSources[datasource.Name]); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// }

func TestAccTokenReadProviderConfigureWithHeaders(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testHeaderConfig,
			},
		},
	})

	rootProvider := Provider().(*schema.Provider)

	rootProviderResource := &schema.Resource{
		Schema: rootProvider.Schema,
	}
	rootProviderData := rootProviderResource.TestResourceData()
	if _, err := providerConfigure(rootProviderData); err != nil {
		t.Fatal(err)
	}
}

func startServerWithConfig(t *testing.T, config string) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	f, err := os.CreateTemp("", "consul_*.hcl")
	if err != nil {
		t.Fatalf("fail to create Consul config file: %s", err)
	}
	if _, err := f.WriteString(config); err != nil {
		t.Fatalf("fail to write Consul config: %s", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("fail to close Consul config file: %s", err)
	}

	path := os.Getenv("CONSUL_TEST_BINARY")
	if path == "" {
		path = "consul"
	}
	cmd := exec.Command(path, "agent", "-dev", "-config-file", f.Name())

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start Consul: %s", err)
	}
	t.Cleanup(func() {
		cmd.Process.Kill()
		cmd.Process.Wait()
	})
}

func waitForService(t *testing.T) (map[string]terraform.ResourceProvider, *consulapi.Client) {
	os.Setenv("CONSUL_HTTP_ADDR", "http://localhost:8500")
	os.Setenv("CONSUL_HTTP_TOKEN", "master-token")

	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		t.Fatalf("failed to instantiate client: %v", err)
	}

	var services []*consulapi.ServiceEntry
	for i := 0; i < 20; i++ {
		services, _, err = client.Health().Service("consul", "", true, nil)
		if err == nil && len(services) == 1 && len(services[0].Node.Meta) == 1 {
			return map[string]terraform.ResourceProvider{
				"consul": Provider(),
			}, client
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatalf("timeout while waiting for Consul to start, last error: %v, %d services", err, len(services))
	return nil, nil
}

func startTestServer(t *testing.T) (map[string]terraform.ResourceProvider, *consulapi.Client) {
	startServerWithConfig(
		t,
		`
			ui = true

			limits = {
				http_max_conns_per_client = -1
			}

			acl = {
				enabled = true
				default_policy = "allow"
				down_policy = "extend-cache"

				tokens = {
					master = "master-token"
				}
			}
		`,
	)

	return waitForService(t)
}

func startRemoteDatacenterTestServer(t *testing.T) (map[string]terraform.ResourceProvider, *consulapi.Client) {
	startServerWithConfig(
		t,
		`
			ui = true
			datacenter = "dc2"
			primary_datacenter = "dc1"

			limits = {
				http_max_conns_per_client = -1
			}

			acl = {
				enabled = true
				default_policy = "allow"
				down_policy = "extend-cache"

				tokens = {
					replication = "master-token"
				}
			}

			ports = {
				dns = -1
				grpc = -1
				http = 8501
				server = 8305
				serf_lan = 8306
				serf_wan = 8307
			}
		`,
	)
	startServerWithConfig(
		t,
		`
			ui = true
			primary_datacenter = "dc1"

			limits = {
				http_max_conns_per_client = -1
			}

			acl = {
				enabled = true
				default_policy = "allow"
				down_policy = "extend-cache"

				tokens = {
					master = "master-token"
				}
			}

			retry_join_wan = ["127.0.0.1:8307"]
		`,
	)

	providers, client := waitForService(t)
	for i := 0; i < 10; i++ {
		datacenters, err := client.Catalog().Datacenters()
		if err == nil && len(datacenters) == 2 {
			return providers, client
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatal("wait for the two datacenters to get synced")
	return nil, nil
}

func serverIsConsulCommunityEdition(t *testing.T) bool {
	path := os.Getenv("CONSUL_TEST_BINARY")
	if path == "" {
		path = "consul"
	}
	cmd := exec.Command(path, "version", "-format=json")

	data, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to get `consul version` output: %v", err)
	}

	type Output struct {
		Version string
	}
	var output Output
	if err := json.Unmarshal(data, &output); err != nil {
		t.Fatalf("failed to unmarshal Consul version: %v", err)
	}

	return !strings.HasSuffix(output.Version, "+ent")
}

func skipTestOnConsulCommunityEdition(t *testing.T) {
	if serverIsConsulCommunityEdition(t) {
		t.Skip("Test skipped on Consul Community Edition. Use a Consul Enterprise server to run this test.")
	}
}

func skipTestOnConsulEnterpriseEdition(t *testing.T) {
	if !serverIsConsulCommunityEdition(t) {
		t.Skip("Test skipped on Consul Enterprise Edition. Use a Consul Community server to run this test.")
	}
}

// lintignore: AT004
var testHeaderConfig = `
provider "consul" {
	header {
		name  = "auth"
		value = "123"
	}
}

data "consul_key_prefix" "read" {
	path_prefix = "foo/"
}
`
