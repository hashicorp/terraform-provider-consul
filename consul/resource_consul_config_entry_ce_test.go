package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntryCE_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulConfigEntryCE_ServiceDefaults,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", "{\"Expose\":{},\"MeshGateway\":{},\"Protocol\":\"https\",\"TransparentProxy\":{}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_ServiceDefaultsOptionalField,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", "{\"Expose\":{},\"MeshGateway\":{},\"Protocol\":\"https\",\"TransparentProxy\":{}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_ProxyDefaults,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "global"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "proxy-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", "{\"Config\":{\"foo\":\"bar\"},\"Expose\":{},\"MeshGateway\":{},\"TransparentProxy\":{}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_ServiceRouter,
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccConsulConfigEntryCE_ServiceSplitter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "kind", "service-splitter"),
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "config_json", "{\"Splits\":[{\"ServiceSubset\":\"v1\",\"Weight\":90},{\"ServiceSubset\":\"v2\",\"Weight\":10}]}"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_ServiceResolver,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "kind", "service-resolver"),
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "config_json", "{\"DefaultSubset\":\"v1\",\"Subsets\":{\"v1\":{\"Filter\":\"Service.Meta.version == v1\"},\"v2\":{\"Filter\":\"Service.Meta.version == v2\"}}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_IngressGateway,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.ingress_gateway", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.ingress_gateway", "kind", "ingress-gateway"),
					resource.TestCheckResourceAttr("consul_config_entry.ingress_gateway", "config_json", "{\"Listeners\":[{\"Port\":8000,\"Protocol\":\"http\",\"Services\":[{\"Hosts\":null,\"Name\":\"*\"}]}],\"TLS\":{\"Enabled\":true}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_TerminatingGateway,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.terminating_gateway", "name", "foo-egress"),
					resource.TestCheckResourceAttr("consul_config_entry.terminating_gateway", "kind", "terminating-gateway"),
					resource.TestCheckResourceAttr("consul_config_entry.terminating_gateway", "config_json", "{\"Services\":[{\"Name\":\"billing\"}]}"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_ServiceConfigL4,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "api-service"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_ServiceConfigL7,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "fort-knox"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_ServiceConfigL7b,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "api"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_ServiceConfigL7gRPC,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "billing"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntryCE_ServiceConfigL7Mixed,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "api"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config:       testAccConsulConfigEntryCE_ServiceConfigL7Mixed,
				ImportState:  true,
				ResourceName: "consul_config_entry.service_intentions",
				ExpectError:  regexp.MustCompile(`expected path of the form "<kind>/<name>" or "<partition>/<namespace>/<kind>/<name>"`),
			},
			{
				Config:        testAccConsulConfigEntryCE_ServiceConfigL7Mixed,
				ImportState:   true,
				ResourceName:  "consul_config_entry.service_intentions",
				ImportStateId: "service-defaults/api",
			},
			{
				Config:        testAccConsulConfigEntryCE_ServiceConfigL7Mixed,
				ImportState:   true,
				ResourceName:  "consul_config_entry.service_intentions",
				ImportStateId: "default/default/service-defaults/api",
			},
		},
	})
}

func TestAccConsulConfigEntryCE_ServicesExported(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: TestAccConsulConfigEntryCE_exportedServicesCE,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.exported_services", "config_json", "{\"Services\":[{\"Consumers\":[{\"Peer\":\"us-east-2\"}],\"Name\":\"test\"}]}"),
					resource.TestCheckResourceAttr("consul_config_entry.exported_services", "id", "exported-services-default"),
					resource.TestCheckResourceAttr("consul_config_entry.exported_services", "kind", "exported-services"),
					resource.TestCheckResourceAttr("consul_config_entry.exported_services", "name", "default"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntryCE_Mesh(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: TestAccConsulConfigEntryCE_meshCE,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.mesh", "name", "mesh"),
					resource.TestCheckResourceAttr("consul_config_entry.mesh", "kind", "mesh"),
				),
			},
		},
	})
}

const testAccConsulConfigEntryCE_ServiceDefaults = `
resource "consul_config_entry" "foo" {
	name = "foo"
	kind = "service-defaults"

	config_json = jsonencode({
		MeshGateway      = {}
		Protocol         = "https"
		TransparentProxy = {}
	})
}
`

const testAccConsulConfigEntryCE_ServiceDefaultsOptionalField = `
resource "consul_config_entry" "foo" {
	name = "foo"
	kind = "service-defaults"

	config_json = jsonencode({
		Expose           = {}
		Protocol         = "https"
		TransparentProxy = {}
	})
}
`

const testAccConsulConfigEntryCE_ProxyDefaults = `
resource "consul_config_entry" "foo" {
	name = "global"
	kind = "proxy-defaults"

	config_json = jsonencode({
		Config = {
			foo = "bar"
		}
		MeshGateway      = {}
		TransparentProxy = {}
	})
}
`

const testAccConsulConfigEntryCE_ServiceRouter = `
resource "consul_config_entry" "web" {
	name = "web"
	kind = "service-defaults"

	config_json = jsonencode({
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
		Protocol         = "http"
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
	})
}

resource "consul_config_entry" "service_router" {
	name = consul_config_entry.web.name
	kind = "service-router"

	config_json = jsonencode({
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
`

const testAccConsulConfigEntryCE_ServiceSplitter = `
resource "consul_config_entry" "web" {
	name = "web"
	kind = "service-defaults"

	config_json = jsonencode({
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
`

const testAccConsulConfigEntryCE_ServiceResolver = `
resource "consul_config_entry" "web" {
	name = "web"
	kind = "service-defaults"

	config_json = jsonencode({
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
`

const testAccConsulConfigEntryCE_ProxyDefaultsWrongName = `
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

const testAccConsulConfigEntryCE_IngressGateway = `
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
			}]
		}]
	})
}
`

const testAccConsulConfigEntryCE_TerminatingGateway = `
resource "consul_config_entry" "terminating_gateway" {
	name = "foo-egress"
	kind = "terminating-gateway"

	config_json = jsonencode({
		Services = [{
			Name = "billing"
		}]
	})
}
`

const testAccConsulConfigEntryCE_ServiceConfigL4 = `
resource "consul_config_entry" "service_intentions" {
	name = "api-service"
	kind = "service-intentions"

	config_json = jsonencode({
		Sources = [
			{
				Action     = "allow"
				Name       = "frontend-webapp"
				Precedence = 9
				Type       = "consul"
			},
            {
				Action     = "allow"
				Name       = "nightly-cronjob"
				Precedence = 9
				Type       = "consul"
			}
		]
	})
}
`

const testAccConsulConfigEntryCE_ServiceConfigL7 = `
resource "consul_config_entry" "sd" {
	name = "fort-knox"
	kind = "service-defaults"

	config_json = jsonencode({
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
		Sources = [
			{
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
`

const testAccConsulConfigEntryCE_ServiceConfigL7b = `
resource "consul_config_entry" "sd" {
	name = "api"
	kind = "service-defaults"

	config_json = jsonencode({
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
		Sources = [
			{
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
`

const testAccConsulConfigEntryCE_ServiceConfigL7gRPC = `
resource "consul_config_entry" "sd" {
	name = "billing"
	kind = "service-defaults"

	config_json = jsonencode({
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
		Sources = [
			{
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
`

const testAccConsulConfigEntryCE_ServiceConfigL7Mixed = `
resource "consul_config_entry" "sd" {
	name = "api"
	kind = "service-defaults"

	config_json = jsonencode({
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
		Sources = [
			{
				Action     = "deny"
				Name       = "hackathon-project"
				Precedence = 9
				Type       = "consul"
			},
			{
				Action     = "allow"
				Name       = "web"
				Precedence = 9
				Type       = "consul"
			},
			{
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
`

const TestAccConsulConfigEntryCE_mesh = `
`

const TestAccConsulConfigEntryCE_exportedServicesCE = `
resource "consul_config_entry" "exported_services" {
	name = "default"
	kind = "exported-services"

	config_json = jsonencode({
		Services = [{
			Name = "test"
			Consumers = [{
				Peer = "us-east-2"
			}]
		}]
	})
}
`

const TestAccConsulConfigEntryCE_exportedServicesEE = `
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

const TestAccConsulConfigEntryCE_meshCE = `
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

const TestAccConsulConfigEntryCE_meshEE = `
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
