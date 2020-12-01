package consul

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntry_basic(t *testing.T) {
	// This needs to be called before serverIsConsulCommunityEdition() as the
	// test provider won't be initialized for unit tests.
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			resource.TestEnvVar))
		return
	}

	// Expected values for Consul Community Edition
	extraConf := ""
	configJSONServiceDefaults := "{\"Expose\":{},\"MeshGateway\":{},\"Protocol\":\"https\"}"
	configJSONProxyDefaults := "{\"Config\":{\"foo\":\"bar\"},\"Expose\":{},\"MeshGateway\":{}}"
	configJSONServiceRouter := "{\"Routes\":[{\"Destination\":{\"Namespace\":\"default\",\"Service\":\"admin\"},\"Match\":{\"HTTP\":{\"PathPrefix\":\"/admin\"}}}]}"
	configJSONServiceSplitter := "{\"Splits\":[{\"ServiceSubset\":\"v1\",\"Weight\":90},{\"ServiceSubset\":\"v2\",\"Weight\":10}]}"
	configJSONServiceResolver := "{\"DefaultSubset\":\"v1\",\"Subsets\":{\"v1\":{\"Filter\":\"Service.Meta.version == v1\"},\"v2\":{\"Filter\":\"Service.Meta.version == v2\"}}}"
	configJSONIngressGateway := "{\"Listeners\":[{\"Port\":8000,\"Protocol\":\"http\",\"Services\":[{\"Hosts\":null,\"Name\":\"*\"}]}],\"TLS\":{\"Enabled\":true}}"
	configJSONTerminatingGateway := "{\"Services\":[{\"Name\":\"billing\"}]}"

	if !serverIsConsulCommunityEdition(t) {
		extraConf = `Namespace: "default"`
		configJSONServiceDefaults = "{\"Expose\":{},\"MeshGateway\":{},\"Namespace\":\"default\",\"Protocol\":\"https\"}"
		configJSONProxyDefaults = "{\"Config\":{\"foo\":\"bar\"},\"Expose\":{},\"MeshGateway\":{},\"Namespace\":\"default\"}"
		configJSONServiceRouter = "{\"Namespace\":\"default\",\"Routes\":[{\"Destination\":{\"Namespace\":\"default\",\"Service\":\"admin\"},\"Match\":{\"HTTP\":{\"PathPrefix\":\"/admin\"}}}]}"
		configJSONServiceSplitter = "{\"Namespace\":\"default\",\"Splits\":[{\"ServiceSubset\":\"v1\",\"Weight\":90},{\"ServiceSubset\":\"v2\",\"Weight\":10}]}"
		configJSONServiceResolver = "{\"DefaultSubset\":\"v1\",\"Namespace\":\"default\",\"Subsets\":{\"v1\":{\"Filter\":\"Service.Meta.version == v1\"},\"v2\":{\"Filter\":\"Service.Meta.version == v2\"}}}"
		configJSONIngressGateway = "{\"Listeners\":[{\"Port\":8000,\"Protocol\":\"http\",\"Services\":[{\"Hosts\":null,\"Name\":\"*\",\"Namespace\":\"default\"}]}],\"Namespace\":\"default\",\"TLS\":{\"Enabled\":true}}"
		configJSONTerminatingGateway = "{\"Namespace\":\"default\",\"Services\":[{\"Name\":\"billing\",\"Namespace\":\"default\"}]}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulConfigEntry_ServiceDefaults(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", configJSONServiceDefaults),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceDefaultsOptionalField(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", configJSONServiceDefaults),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceDefaultsExtraField,
				ExpectError: regexp.MustCompile(`errors during apply: Failed to decode config entry: 1 error\(s\) decoding:

\* '' has invalid keys: ThisFieldDoesNotExists`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", configJSONServiceDefaults),
				),
			},
			{
				Config: testAccConsulConfigEntry_ProxyDefaults(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "global"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "proxy-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", configJSONProxyDefaults),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceRouter(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_router", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_router", "kind", "service-router"),
					resource.TestCheckResourceAttr("consul_config_entry.service_router", "config_json", configJSONServiceRouter),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceSplitter(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "kind", "service-splitter"),
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "config_json", configJSONServiceSplitter),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceResolver(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "kind", "service-resolver"),
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "config_json", configJSONServiceResolver),
				),
			},
			{
				Config: testAccConsulConfigEntry_IngressGateway(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.ingress_gateway", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.ingress_gateway", "kind", "ingress-gateway"),
					resource.TestCheckResourceAttr("consul_config_entry.ingress_gateway", "config_json", configJSONIngressGateway),
				),
			},
			{
				Config: testAccConsulConfigEntry_TerminatingGateway(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.terminating_gateway", "name", "foo-egress"),
					resource.TestCheckResourceAttr("consul_config_entry.terminating_gateway", "kind", "terminating-gateway"),
					resource.TestCheckResourceAttr("consul_config_entry.terminating_gateway", "config_json", configJSONTerminatingGateway),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceConfigL4(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "api-service"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceConfigL7(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "fort-knox"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceConfigL7b(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "api"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceConfigL7gRPC(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "billing"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntry_ServiceConfigL7Mixed(extraConf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "api"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntry_Errors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulConfigEntry_ProxyDefaultsWrongName,
				ExpectError: regexp.MustCompile("failed to read config entry after setting it.\nThis may happen when some attributes have an unexpected value.\nRead the documentation at https://www.consul.io/docs/agent/config-entries/proxy-defaults.html\nto see what values are expected."),
			},
		},
	})
}

func testAccConsulConfigEntry_ServiceDefaults(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "foo" {
	name = "foo"
	kind = "service-defaults"

	config_json = jsonencode({
		MeshGateway = {}
		Protocol    = "https"
		%s
	})
}
`, extraConf)
}

func testAccConsulConfigEntry_ServiceDefaultsOptionalField(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "foo" {
	name = "foo"
	kind = "service-defaults"

	config_json = jsonencode({
		Protocol    = "https"
		%s
	})
}
`, extraConf)
}

const testAccConsulConfigEntry_ServiceDefaultsExtraField = `
resource "consul_config_entry" "foo" {
	name = "foo"
	kind = "service-defaults"

	config_json = jsonencode({
		ThisFieldDoesNotExists = true
		Protocol               = "https"
	})
}
`

func testAccConsulConfigEntry_ProxyDefaults(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "foo" {
	name = "global"
	kind = "proxy-defaults"

	config_json = jsonencode({
		Config = {
			foo = "bar"
		}
		%s
	})
}
`, extraConf)
}

func testAccConsulConfigEntry_ServiceRouter(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "web" {
	name = "web"
	kind = "service-defaults"

	config_json = jsonencode({
		Protocol = "http"
		%s
	})
}

resource "consul_config_entry" "admin_service_defaults" {
	name = "admin"
	kind = "service-defaults"

	config_json = jsonencode({
		Protocol = "http"
		%s
	})
}

resource "consul_config_entry" "service_router" {
	name = consul_config_entry.web.name
	kind = "service-router"

	config_json = jsonencode({
		%s
		Routes = [
			{
				Match = {
					HTTP = {
						PathPrefix = "/admin"
					}
				}

				Destination = {
					Namespace = "default"
					Service   = consul_config_entry.admin_service_defaults.name
				}
			}
			# NOTE: a default catch-all will send unmatched traffic to "web"
		]
	})
}
`, extraConf, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceSplitter(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "web" {
	name = "web"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol = "http"
	})
}

resource "consul_config_entry" "service_resolver" {
	kind = "service-resolver"
	name = consul_config_entry.web.name

	config_json = jsonencode({
		%s
		DefaultSubset = "v1"

		Subsets = {
			"v1" = {
				Filter = "Service.Meta.version == v1"
			}
			"v2" = {
				Filter = "Service.Meta.version == v2"
			}
		}
	})
}

resource "consul_config_entry" "service_splitter" {
	kind = "service-splitter"
	name = consul_config_entry.service_resolver.name

	config_json = jsonencode({
		%s
		Splits = [
			{
				Weight         = 90
				ServiceSubset = "v1"
			},
			{
				Weight        = 10
				ServiceSubset = "v2"
			},
		]
	})
}
`, extraConf, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceResolver(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "web" {
	name = "web"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol = "http"
	})
}

resource "consul_config_entry" "service_resolver" {
	kind = "service-resolver"
	name = consul_config_entry.web.name

	config_json = jsonencode({
		%s
		DefaultSubset = "v1"

		Subsets = {
			"v1" = {
				Filter = "Service.Meta.version == v1"
			}
			"v2" = {
				Filter = "Service.Meta.version == v2"
			}
		}

	})
}
`, extraConf, extraConf)
}

const testAccConsulConfigEntry_ProxyDefaultsWrongName = `
resource "consul_config_entry" "foo" {
	name = "foo"
	kind = "proxy-defaults"

	config_json = jsonencode({
		Config = {
			foo = "bar"
		}
	})
}
`

func testAccConsulConfigEntry_IngressGateway(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "ingress_gateway" {
	name = "foo"
	kind = "ingress-gateway"

	config_json = jsonencode({
		%s
		TLS = {
			Enabled = true
		}
		Listeners = [{
			Port     = 8000
			Protocol = "http"
			Services = [{
				Hosts = null
				Name  = "*"
				%s
			}]
		}]
	})
}
`, extraConf, extraConf)
}

func testAccConsulConfigEntry_TerminatingGateway(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "terminating_gateway" {
	name = "foo-egress"
	kind = "terminating-gateway"

	config_json = jsonencode({
		%s
		Services = [{
			Name = "billing"
			%s
		}]
	})
}
`, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceConfigL4(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "service_intentions" {
	name = "api-service"
	kind = "service-intentions"

	config_json = jsonencode({
		%s
		Sources = [
			{
				%s
				Action     = "allow"
				Name       = "frontend-webapp"
				Precedence = 9
				Type       = "consul"
			},
            {
				%s
				Action     = "allow"
				Name       = "nightly-cronjob"
				Precedence = 9
				Type       = "consul"
			}
		]
	})
}
`, extraConf, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceConfigL7(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "sd" {
	name = "fort-knox"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol = "http"
	})
}

resource "consul_config_entry" "service_intentions" {
	name = consul_config_entry.sd.name
	kind = "service-intentions"

	config_json = jsonencode({
		%s
		Sources = [
			{
				%s
				Name        = "contractor-webapp"
				Permissions = [
					{
						Action = "allow"
						HTTP   = {
							Methods   = ["GET", "HEAD"]
							PathExact = "/healtz"
						}
					}
				]
				Precedence = 9
				Type       = "consul"
			},
			{
				%s
				Name        = "admin-dashboard-webapp",
				Permissions = [
					{
						Action = "deny",
						HTTP = {
							PathPrefix= "/debugz"
						}
					},
					{
						Action= "allow"
						HTTP = {
							PathPrefix= "/"
						}
					}
				],
				Precedence = 9
				Type       = "consul"
			}
		]
	})
}
`, extraConf, extraConf, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceConfigL7b(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "sd" {
	name = "api"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol = "http"
	})
}

resource "consul_config_entry" "service_intentions" {
	name = consul_config_entry.sd.name
	kind = "service-intentions"

	config_json = jsonencode({
		%s
		Sources = [
			{
				%s
				Name        = "admin-dashboard"
				Permissions = [
					{
						Action = "allow"
						HTTP = {
							Methods    = ["GET", "PUT", "POST", "DELETE", "HEAD"]
							PathPrefix = "/v2"
						}
					}
				],
				Precedence = 9
				Type = "consul"
			},
			{
				%s
				Name = "report-generator"
				Permissions = [
					{
						Action = "allow"
						HTTP = {
							Methods = ["GET"]
							PathPrefix = "/v2/widgets"
						}
					}
				],
				Precedence = 9,
				Type = "consul"
			}
		]
	})
}
`, extraConf, extraConf, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceConfigL7gRPC(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "sd" {
	name = "billing"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol = "grpc"
	})
}

resource "consul_config_entry" "service_intentions" {
	name = consul_config_entry.sd.name
	kind = "service-intentions"

	config_json = jsonencode({
		%s
		Sources = [
			{
				%s
				Name = "frontend-web"
				Permissions = [
					{
						Action = "deny"
						HTTP = {
							PathExact = "/mycompany.BillingService/IssueRefund"
						}
					},
					{
						Action = "allow"
						HTTP = {
							PathPrefix = "/mycompany.BillingService/"
						}
					}
				],
				Precedence = 9
				Type = "consul"
			},
			{
				%s
				Name = "support-portal"
				Permissions = [
					{
						Action = "allow"
						HTTP = {
							PathPrefix = "/mycompany.BillingService/"
						}
					}
				],
				Precedence = 9
				Type = "consul"
			}
		]
	})
}
`, extraConf, extraConf, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceConfigL7Mixed(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "sd" {
	name = "api"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol = "grpc"
	})
}

resource "consul_config_entry" "service_intentions" {
	name = consul_config_entry.sd.name
	kind = "service-intentions"

	config_json = jsonencode({
		%s
		Sources = [
			{
				%s
				Action     = "deny"
				Name       = "hackathon-project"
				Precedence = 9
				Type       = "consul"
			},
			{
				%s
				Action     = "allow"
				Name       = "web"
				Precedence = 9
				Type       = "consul"
			},
			{
				%s
				Name = "nightly-reconciler"
				Permissions = [
					{
						Action = "allow"
						HTTP = {
							Methods   = ["POST"]
							PathExact = "/v1/reconcile-data"
						}
					}
				]
				Precedence = 9
				Type       = "consul"
			}
		]
	})
}
`, extraConf, extraConf, extraConf, extraConf, extraConf)
}
