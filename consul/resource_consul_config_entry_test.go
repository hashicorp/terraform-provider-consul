package consul

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntry_basic(t *testing.T) {
	startTestServer(t)

	// Expected values for Consul Community Edition
	extraConf := ""
	configJSONServiceDefaults := "{\"Expose\":{},\"MeshGateway\":{},\"Protocol\":\"https\",\"TransparentProxy\":{}}"
	configJSONProxyDefaults := "{\"Config\":{\"foo\":\"bar\"},\"Expose\":{},\"MeshGateway\":{},\"TransparentProxy\":{}}"
	configJSONServiceRouter := "{\"Routes\":[{\"Destination\":{\"Namespace\":\"default\",\"Partition\":\"default\",\"Service\":\"admin\"},\"Match\":{\"HTTP\":{\"PathPrefix\":\"/admin\"}}}]}"
	configJSONServiceSplitter := "{\"Splits\":[{\"ServiceSubset\":\"v1\",\"Weight\":90},{\"ServiceSubset\":\"v2\",\"Weight\":10}]}"
	configJSONServiceResolver := "{\"DefaultSubset\":\"v1\",\"Subsets\":{\"v1\":{\"Filter\":\"Service.Meta.version == v1\"},\"v2\":{\"Filter\":\"Service.Meta.version == v2\"}}}"
	configJSONIngressGateway := "{\"Listeners\":[{\"Port\":8000,\"Protocol\":\"http\",\"Services\":[{\"Hosts\":null,\"Name\":\"*\"}]}],\"TLS\":{\"Enabled\":true}}"
	configJSONTerminatingGateway := "{\"Services\":[{\"Name\":\"billing\"}]}"

	if !serverIsConsulCommunityEdition(t) {
		extraConf = `Namespace: "default", Partition: "default"`
		configJSONServiceDefaults = "{\"Expose\":{},\"MeshGateway\":{},\"Partition\":\"default\",\"Protocol\":\"https\",\"TransparentProxy\":{}}"
		configJSONProxyDefaults = "{\"Config\":{\"foo\":\"bar\"},\"Expose\":{},\"MeshGateway\":{},\"Partition\":\"default\",\"TransparentProxy\":{}}"
		configJSONServiceRouter = "{\"Partition\":\"default\",\"Routes\":[{\"Destination\":{\"Namespace\":\"default\",\"Partition\":\"default\",\"Service\":\"admin\"},\"Match\":{\"HTTP\":{\"PathPrefix\":\"/admin\"}}}]}"
		configJSONServiceSplitter = "{\"Partition\":\"default\",\"Splits\":[{\"ServiceSubset\":\"v1\",\"Weight\":90},{\"ServiceSubset\":\"v2\",\"Weight\":10}]}"
		configJSONServiceResolver = "{\"DefaultSubset\":\"v1\",\"Partition\":\"default\",\"Subsets\":{\"v1\":{\"Filter\":\"Service.Meta.version == v1\"},\"v2\":{\"Filter\":\"Service.Meta.version == v2\"}}}"
		configJSONIngressGateway = "{\"Listeners\":[{\"Port\":8000,\"Protocol\":\"http\",\"Services\":[{\"Hosts\":null,\"Name\":\"*\",\"Namespace\":\"default\",\"Partition\":\"default\"}]}],\"Partition\":\"default\",\"TLS\":{\"Enabled\":true}}"
		configJSONTerminatingGateway = "{\"Partition\":\"default\",\"Services\":[{\"Name\":\"billing\",\"Namespace\":\"default\"}]}"
	}

	resource.Test(t, resource.TestCase{
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
	startTestServer(t)

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

func TestAccConsulConfigEntry_NamespaceEE(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulConfigEntry_DefaultNamespace,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "namespace", "default"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
				),
			},
			{
				Config: testAccConsulConfigEntry_Namespace,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.test_intentions", "name", "destination-service"),
					resource.TestCheckResourceAttr("consul_config_entry.test_intentions", "namespace", "example"),
					resource.TestCheckResourceAttr("consul_config_entry.test_intentions", "kind", "service-intentions"),
					resource.TestCheckResourceAttr("consul_config_entry.test_intentions", "config_json", "{\"Meta\":{\"foo\":\"bar\"},\"Partition\":\"default\",\"Sources\":[{\"Action\":\"allow\",\"Name\":\"source-service\",\"Namespace\":\"example\",\"Partition\":\"default\",\"Precedence\":9,\"Type\":\"consul\"}]}"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntry_ServicesExportedCE(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      TestAccConsulConfigEntry_exportedServicesCE,
				ExpectError: regexp.MustCompile(`Config entry kind "exported-services" requires Consul Enterprise`),
			},
		},
	})
}

func TestAccConsulConfigEntry_ServicesExportedEE(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccConsulConfigEntry_exportedServicesEE,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.exported_services", "name", "test"),
					resource.TestCheckResourceAttr("consul_config_entry.exported_services", "kind", "exported-services"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntry_MeshCE(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccConsulConfigEntry_meshCE,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.mesh", "name", "mesh"),
					resource.TestCheckResourceAttr("consul_config_entry.mesh", "kind", "mesh"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntry_MeshEE(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccConsulConfigEntry_meshEE,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.mesh", "name", "mesh"),
					resource.TestCheckResourceAttr("consul_config_entry.mesh", "kind", "mesh"),
				),
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
		MeshGateway      = {}
		Protocol         = "https"
		TransparentProxy = {}
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
		Expose           = {}
		Protocol         = "https"
		TransparentProxy = {}
		%s
	})
}
`, extraConf)
}

func testAccConsulConfigEntry_ProxyDefaults(extraConf string) string {
	return fmt.Sprintf(`
resource "consul_config_entry" "foo" {
	name = "global"
	kind = "proxy-defaults"

	config_json = jsonencode({
		Config = {
			foo = "bar"
		}
		MeshGateway      = {}
		TransparentProxy = {}
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
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
		Protocol         = "http"
		%s
	})
}

resource "consul_config_entry" "admin_service_defaults" {
	name = "admin"
	kind = "service-defaults"

	config_json = jsonencode({
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
		Protocol         = "http"
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
					Partition = "default"
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
		Protocol         = "http"
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
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
		Protocol         = "http"
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
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
	var partition string
	if extraConf != "" {
		partition = `Partition = "default"`
	}

	return fmt.Sprintf(`
resource "consul_config_entry" "ingress_gateway" {
	name = "foo"
	kind = "ingress-gateway"

	config_json = jsonencode({
		TLS = {
			Enabled = true
		}
		Listeners = [{
			Port      = 8000
			Protocol  = "http"
			Services = [{
				Hosts = null
				Name  = "*"
				%s
			}]
		}]
		%s
	})
}
`, extraConf, partition)
}

func testAccConsulConfigEntry_TerminatingGateway(extraConf string) string {
	var partition string
	if extraConf != "" {
		extraConf = `Namespace: "default"`
		partition = `Partition = "default"`
	}

	return fmt.Sprintf(`
resource "consul_config_entry" "terminating_gateway" {
	name = "foo-egress"
	kind = "terminating-gateway"

	config_json = jsonencode({
		Services = [{
			Name = "billing"
			%s
		}]
		%s
	})
}
`, extraConf, partition)
}

func testAccConsulConfigEntry_ServiceConfigL4(extraConf string) string {
	var partition string
	if extraConf != "" {
		partition = `Partition = "default"`
	}

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
`, partition, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceConfigL7(extraConf string) string {
	var partition string
	if extraConf != "" {
		partition = `Partition = "default"`
	}

	return fmt.Sprintf(`
resource "consul_config_entry" "sd" {
	name = "fort-knox"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol         = "http"
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
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
`, partition, partition, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceConfigL7b(extraConf string) string {
	var partition string
	if extraConf != "" {
		partition = `Partition = "default"`
	}

	return fmt.Sprintf(`
resource "consul_config_entry" "sd" {
	name = "api"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol         = "http"
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
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
				Name      = "report-generator"
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
`, partition, partition, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceConfigL7gRPC(extraConf string) string {
	var partition string
	if extraConf != "" {
		partition = `Partition = "default"`
	}

	return fmt.Sprintf(`
resource "consul_config_entry" "sd" {
	name = "billing"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol         = "grpc"
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
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
`, partition, partition, extraConf, extraConf)
}

func testAccConsulConfigEntry_ServiceConfigL7Mixed(extraConf string) string {
	var partition string
	if extraConf != "" {
		partition = `Partition = "default"`
	}

	return fmt.Sprintf(`
resource "consul_config_entry" "sd" {
	name = "api"
	kind = "service-defaults"

	config_json = jsonencode({
		%s
		Protocol         = "grpc"
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
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
`, extraConf, partition, extraConf, extraConf, extraConf)
}

const testAccConsulConfigEntry_DefaultNamespace = `
resource "consul_config_entry" "foo" {
	name      = "foo"
	kind      = "service-defaults"
	namespace = "default"

	config_json = jsonencode({
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
		Protocol         = "https"
		Namespace        = "default"
		Partition        = "default"
	})
}
`

const testAccConsulConfigEntry_Namespace = `
resource "consul_namespace" "example_namespace" {
	name = "example"
	description = "Example namespace"
}

resource "consul_config_entry" "test_intentions" {
	name = "destination-service"
	kind = "service-intentions"
	namespace = consul_namespace.example_namespace.name

	config_json = jsonencode({
		Sources = [
		  {
			Action     = "allow"
			Name       = "source-service"
			Namespace  = "example"
			Partition  = "default"
			Precedence = 9
			Type       = "consul"
		  }
		]
		Meta = {
			foo = "bar"
		}
		Partition = "default"
	  })
}
`

const TestAccConsulConfigEntry_mesh = `
`

const TestAccConsulConfigEntry_exportedServicesCE = `
resource "consul_config_entry" "exported_services" {
	name = "test"
	kind = "exported-services"

	config_json = jsonencode({
		Services = [{
			Name = "test"
			Namespace = "default"
			Consumers = [{
				Partition = "default"
			}]
		}]
	})
}
`

const TestAccConsulConfigEntry_exportedServicesEE = `
resource "consul_admin_partition" "test" {
	name = "test"
}

resource "consul_config_entry" "exported_services" {
	name = "test"
	kind = "exported-services"

	config_json = jsonencode({
		Partition = consul_admin_partition.test.name
		Services = [{
			Name = "test"
			Namespace = "default"
			Consumers = [{
				Partition = "default"
			}]
		}]
	})
}
`

const TestAccConsulConfigEntry_meshCE = `
resource "consul_config_entry" "mesh" {
	name = "mesh"
	kind = "mesh"

	config_json = jsonencode({
		TransparentProxy = {
			MeshDestinationsOnly = true
		}
	})
}
`

const TestAccConsulConfigEntry_meshEE = `
resource "consul_config_entry" "mesh" {
	name = "mesh"
	kind = "mesh"

	config_json = jsonencode({
		Partition = "default"

		TransparentProxy = {
			MeshDestinationsOnly = true
		}
	})
}
`
