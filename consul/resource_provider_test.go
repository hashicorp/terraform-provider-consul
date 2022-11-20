package consul

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/sdk/testutil/retry"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	initialManagementToken = "12345678-1234-1234-1234-1234567890ab"
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

	caPem, err := os.ReadFile("test-fixtures/cacert.pem")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	certPem, err := os.ReadFile("test-fixtures/usercert.pem")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	keyPem, err := os.ReadFile("test-fixtures/userkey.pem")
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

func startServerWithConfig(t *testing.T, configFile string) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	os.Setenv("CONSUL_HTTP_TOKEN", initialManagementToken)

	path := os.Getenv("CONSUL_TEST_BINARY")
	if path == "" {
		path = "consul"
	}
	cmd := exec.Command(path, "agent", "-dev", "-config-file", "test-fixtures/"+configFile)

	if os.Getenv("TF_ACC_CONSUL_LOG") != "" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start Consul: %s", err)
	}
	t.Cleanup(func() {
		cmd.Process.Kill()
		cmd.Process.Wait()
	})
}

func waitForService(t *testing.T, address string) (terraform.ResourceProvider, *consulapi.Client) {
	config := consulapi.DefaultConfig()
	config.Address = address
	config.Token = initialManagementToken
	client, err := consulapi.NewClient(config)
	if err != nil {
		t.Fatalf("failed to instantiate client: %v", err)
	}

	var services []*consulapi.ServiceEntry
	for i := 0; i < 20; i++ {
		services, _, err = client.Health().Service("consul", "", true, nil)
		if err == nil && len(services) == 1 && len(services[0].Node.Meta) == 1 {
			return Provider(), client
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatalf("timeout while waiting for %s to start, last error: %v, %d services", address, err, len(services))
	return nil, nil
}

func startTestServer(t *testing.T) (map[string]terraform.ResourceProvider, *consulapi.Client) {
	startServerWithConfig(t, "consul.hcl")

	provider, client := waitForService(t, "http://localhost:8500")

	return map[string]terraform.ResourceProvider{
		"consul": provider,
	}, client
}

func startRemoteDatacenterTestServer(t *testing.T) (map[string]terraform.ResourceProvider, *consulapi.Client) {
	if os.Getenv("SKIP_REMOTE_DATACENTER_TESTS") != "" {
		t.Skip("Remote datacenter skipped because SKIP_REMOTE_DATACENTER_TESTS is set")
	}

	startServerWithConfig(t, "consul.hcl")
	startServerWithConfig(t, "consul-secondary.hcl")

	provider, client := waitForService(t, "http://localhost:8500")
	remoteProvider, _ := waitForService(t, "http://localhost:9500")

	for i := 0; i < 20; i++ {
		datacenters, err := client.Catalog().Datacenters()
		if err == nil && len(datacenters) == 2 {
			return map[string]terraform.ResourceProvider{
				"consul":       provider,
				"consulremote": remoteProvider,
			}, client
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatal("wait for the two datacenters to get synced")
	return nil, nil
}

func waitForActiveCARoot(t testing.TB, address string) {
	// don't need to fully decode the response
	type rootsResponse struct {
		ActiveRootID string
		TrustDomain  string
		Roots        []interface{}
	}

	retry.Run(t, func(r *retry.R) {
		// Query the API and check the status code.
		url := address + "/v1/agent/connect/ca/roots"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("x-consul-token", initialManagementToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			r.Fatalf("failed http get '%s': %v", url, err)
		}
		defer resp.Body.Close()
		// Roots will return an error status until it's been bootstrapped. We could
		// parse the body and sanity check but that causes either import cycles
		// since this is used in both `api` and consul test or duplication. The 200
		// is all we really need to wait for.
		if resp.StatusCode != 200 {
			r.Fatalf("failed OK response: Bad status code: %d", resp.StatusCode)
		}

		var roots rootsResponse

		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&roots); err != nil {
			r.Fatal(err)
		}

		if roots.ActiveRootID == "" || len(roots.Roots) < 1 {
			r.Fatalf("/v1/agent/connect/ca/roots returned 200 but without roots: %+v", roots)
		}
	})
}

func startPeeringTestServers(t *testing.T) (map[string]terraform.ResourceProvider, *consulapi.Client) {
	if os.Getenv("SKIP_REMOTE_DATACENTER_TESTS") != "" {
		t.Skip("Remote datacenter skipped because SKIP_REMOTE_DATACENTER_TESTS is set")
	}

	startServerWithConfig(t, "consul-peering-blue.hcl")
	startServerWithConfig(t, "consul-peering-green.hcl")

	provider, client := waitForService(t, "http://localhost:8500")
	remoteProvider, _ := waitForService(t, "http://localhost:9500")

	waitForActiveCARoot(t, "http://localhost:8500")
	waitForActiveCARoot(t, "http://localhost:9500")

	return map[string]terraform.ResourceProvider{
		"consul":       provider,
		"consulremote": remoteProvider,
	}, client
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
