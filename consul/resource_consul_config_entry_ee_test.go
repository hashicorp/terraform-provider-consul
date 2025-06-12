// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConsulConfigEntryEE_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulConfigEntryEE_ServiceDefaults,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", "{\"Expose\":{},\"MeshGateway\":{},\"Protocol\":\"https\",\"TransparentProxy\":{}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_ServiceDefaultsPartition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "namespace", "ns1"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "partition", "part1"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", "{\"Expose\":{},\"MeshGateway\":{},\"Protocol\":\"https\",\"TransparentProxy\":{}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_ProxyDefaults,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "global"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "proxy-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", "{\"AccessLogs\":{},\"Config\":{\"foo\":\"bar\"},\"Expose\":{},\"MeshGateway\":{},\"TransparentProxy\":{}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_ServiceRouter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_router", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_router", "kind", "service-router"),
					resource.TestCheckResourceAttr("consul_config_entry.service_router", "config_json", "{\"Routes\":[{\"Destination\":{\"Namespace\":\"default\",\"Partition\":\"default\",\"Service\":\"admin\"},\"Match\":{\"HTTP\":{\"PathPrefix\":\"/admin\"}}}]}"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_ServiceSplitter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "kind", "service-splitter"),
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "config_json", "{\"Splits\":[{\"ServiceSubset\":\"v1\",\"Weight\":90},{\"ServiceSubset\":\"v2\",\"Weight\":10}]}"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_ServiceResolver,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "kind", "service-resolver"),
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "config_json", "{\"DefaultSubset\":\"v1\",\"Subsets\":{\"v1\":{\"Filter\":\"Service.Meta.version == v1\"},\"v2\":{\"Filter\":\"Service.Meta.version == v2\"}}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_IngressGateway,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.ingress_gateway", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.ingress_gateway", "kind", "ingress-gateway"),
					resource.TestCheckResourceAttr("consul_config_entry.ingress_gateway", "config_json", "{\"Listeners\":[{\"Port\":8000,\"Protocol\":\"http\",\"Services\":[{\"Hosts\":null,\"Name\":\"*\",\"Namespace\":\"default\",\"Partition\":\"default\"}]}],\"TLS\":{\"Enabled\":true}}"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_TerminatingGateway,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.terminating_gateway", "name", "foo-egress"),
					resource.TestCheckResourceAttr("consul_config_entry.terminating_gateway", "kind", "terminating-gateway"),
					resource.TestCheckResourceAttr("consul_config_entry.terminating_gateway", "config_json", "{\"Services\":[{\"Name\":\"billing\",\"Namespace\":\"default\"}]}"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_ServiceConfigAdminPartition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "example_server"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "partition", "part2"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "namespace", "ns2"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "config_json", "{\"Sources\":[{\"Action\":\"allow\",\"Name\":\"example_client\",\"Namespace\":\"ns1\",\"Partition\":\"part1\",\"Precedence\":9,\"Type\":\"consul\"}]}"),
				),
			},
			{
				Config:             testAccConsulConfigEntryEE_ServiceConfigL4,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "api-service"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config:             testAccConsulConfigEntryEE_ServiceConfigL7,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "fort-knox"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config:             testAccConsulConfigEntryEE_ServiceConfigL7b,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "api"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_ServiceConfigL7gRPC,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "billing"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_ServiceConfigL7Mixed,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "name", "api"),
					resource.TestCheckResourceAttr("consul_config_entry.service_intentions", "kind", "service-intentions"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntryEE_Namespace(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulConfigEntryEE_DefaultNamespace,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "namespace", "default"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
				),
			},
			{
				Config: testAccConsulConfigEntryEE_Namespace,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.test_intentions", "name", "destination-service"),
					resource.TestCheckResourceAttr("consul_config_entry.test_intentions", "namespace", "example"),
					resource.TestCheckResourceAttr("consul_config_entry.test_intentions", "kind", "service-intentions"),
					resource.TestCheckResourceAttr("consul_config_entry.test_intentions", "config_json", "{\"Meta\":{\"foo\":\"bar\"},\"Sources\":[{\"Action\":\"allow\",\"Name\":\"source-service\",\"Namespace\":\"example\",\"Partition\":\"default\",\"Precedence\":9,\"Type\":\"consul\"}]}"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntryEE_ServicesExported(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: TestAccConsulConfigEntryEE_exportedServicesEE,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.exported_services", "name", "test"),
					resource.TestCheckResourceAttr("consul_config_entry.exported_services", "kind", "exported-services"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntryEE_Mesh(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: TestAccConsulConfigEntryEE_meshEE,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.mesh", "name", "mesh"),
					resource.TestCheckResourceAttr("consul_config_entry.mesh", "kind", "mesh"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntryEE_JWTProvider_Remote(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             TestAccConsulConfigEntryEE_jwtRemote,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.jwt_provider", "id", "jwt-provider-okta"),
					resource.TestCheckResourceAttr("consul_config_entry.jwt_provider", "name", "okta"),
					resource.TestCheckResourceAttr("consul_config_entry.jwt_provider", "kind", "jwt-provider"),
					resource.TestCheckResourceAttr("consul_config_entry.jwt_provider", "config_json", "{\"ClockSkewSeconds\":30,\"Forwarding\":{\"HeaderName\":\"test-token\"},\"Issuer\":\"test-issuer\",\"JSONWebKeySet\":{\"Remote\":{\"FetchAsynchronously\":true,\"URI\":\"https://127.0.0.1:9091\"}}}"),
				),
			},
		},
	})
}

func TestAccConsulConfigEntryEE_JWTProvider_Local(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:             TestAccConsulConfigEntryEE_jwtLocal,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.jwt_provider", "id", "jwt-provider-auth0"),
					resource.TestCheckResourceAttr("consul_config_entry.jwt_provider", "name", "auth0"),
					resource.TestCheckResourceAttr("consul_config_entry.jwt_provider", "kind", "jwt-provider"),
					resource.TestCheckResourceAttr("consul_config_entry.jwt_provider", "config_json", "{\"ClockSkewSeconds\":30,\"Issuer\":\"auth0-issuer\",\"JSONWebKeySet\":{\"Local\":{\"JWKS\":\"eyJrZXlzIjogW3sKICAiY3J2IjogIlAtMjU2IiwKICAia2V5X29wcyI6IFsKICAgICJ2ZXJpZnkiCiAgXSwKICAia3R5IjogIkVDIiwKICAieCI6ICJXYzl1WnVQYUI3S2gyRk1jOXd0SmpSZThYRDR5VDJBWU5BQWtyWWJWanV3IiwKICAieSI6ICI2OGhSVEppSk5Pd3RyaDRFb1BYZVZuUnVIN2hpU0RKX2xtYmJqZkRmV3EwIiwKICAiYWxnIjogIkVTMjU2IiwKICAidXNlIjogInNpZyIsCiAgImtpZCI6ICJhYzFlOGY5MGVkZGY2MWM0MjljNjFjYTA1YjRmMmUwNyIKfV19\"}}}"),
				),
			},
		},
	})
}

const testAccConsulConfigEntryEE_ServiceDefaults = `
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

const testAccConsulConfigEntryEE_ServiceDefaultsPartition = `
resource "consul_admin_partition" "part1" {
  name = "part1"
}

resource "consul_namespace" "ns1" {
  name = "ns1"
  partition = consul_admin_partition.part1.name
}

resource "consul_config_entry" "foo" {
  name      = "foo"
  kind      = "service-defaults"
  namespace = consul_namespace.ns1.name
  partition = consul_admin_partition.part1.name

  config_json = jsonencode({
    Expose           = {}
    Protocol         = "https"
    TransparentProxy = {}
  })
}
`

const testAccConsulConfigEntryEE_ProxyDefaults = `
resource "consul_config_entry" "foo" {
	name = "global"
	kind = "proxy-defaults"

	config_json = jsonencode({
		AccessLogs = {}
		Config = {
			foo = "bar"
		}
		MeshGateway      = {}
		TransparentProxy = {}
	})
}
`

const testAccConsulConfigEntryEE_ServiceRouter = `
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

const testAccConsulConfigEntryEE_ServiceSplitter = `
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

const testAccConsulConfigEntryEE_ServiceResolver = `
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

const testAccConsulConfigEntryEE_IngressGateway = `
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
				Namespace        = "default"
				Partition        = "default"
			}]
		}]
	})
}
`

const testAccConsulConfigEntryEE_TerminatingGateway = `
resource "consul_config_entry" "terminating_gateway" {
	name = "foo-egress"
	kind = "terminating-gateway"

	config_json = jsonencode({
		Services = [{
			Name = "billing"
			Namespace: "default"
		}]
	})
}
`

const testAccConsulConfigEntryEE_ServiceConfigAdminPartition = `
resource "consul_admin_partition" "part1" {
  name = "part1"
}

resource "consul_admin_partition" "part2" {
  name = "part2"
}

resource "consul_namespace" "ns1" {
  name = "ns1"
  partition = consul_admin_partition.part1.name
}

resource "consul_namespace" "ns2" {
  name = "ns2"
  partition = consul_admin_partition.part2.name
}

resource "consul_config_entry" "service_intentions" {
  kind      = "service-intentions"
  name      = "example_server"
  namespace = consul_namespace.ns2.name
  partition = consul_admin_partition.part2.name

  config_json = jsonencode({
    Sources = [{
      Action     = "allow"
      Name       = "example_client"
      Namespace  = consul_namespace.ns1.name
      Partition  = consul_admin_partition.part1.name
      Precedence = 9
      Type       = "consul"
    }]
  })
}
`

const testAccConsulConfigEntryEE_ServiceConfigL4 = `
resource "consul_config_entry" "jwt_provider" {
	name = "okta"
	kind = "jwt-provider"

	config_json = jsonencode({
		Issuer = "test-issuer"
		JSONWebKeySet = {
			Remote = {
				URI = "https://127.0.0.1:9091"
				FetchAsynchronously = true
			}
		}
	})
}

resource "consul_config_entry" "service_intentions" {
	name = "api-service"
	kind = "service-intentions"

	config_json = jsonencode({
		JWT = {
			Providers = [
				{ name = "okta" }
			]
		}
		Sources = [
			{
				Namespace  = "default"
				Partition  = "default"
				Action     = "allow"
				Name       = "frontend-webapp"
				Precedence = 9
				Type       = "consul"
			},
            {
				Namespace  = "default"
				Partition  = "default"
				Action     = "allow"
				Name       = "nightly-cronjob"
				Precedence = 9
				Type       = "consul"
			}
		]
	})
}
`

const testAccConsulConfigEntryEE_ServiceConfigL7 = `
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

resource "consul_config_entry" "jwt_provider" {
	name = "okta"
	kind = "jwt-provider"

	config_json = jsonencode({
		Issuer = "test-issuer"
		JSONWebKeySet = {
			Remote = {
				URI = "https://127.0.0.1:9091"
				FetchAsynchronously = true
			}
		}
	})
}

resource "consul_config_entry" "service_intentions" {
	name = consul_config_entry.sd.name
	kind = "service-intentions"

	config_json = jsonencode({
		Sources = [
			{
				Namespace   = "default"
				Partition   = "default"
				Name        = "contractor-webapp"
				Permissions = [
					{
						Action = "allow"
						HTTP   = {
							Methods   = ["GET", "HEAD"]
							PathExact = "/healtz"
						}
						JWT = {
							Providers = [
								{ name = consul_config_entry.jwt_provider.name }
							]
						}
					}
				]
				Precedence = 9
				Type       = "consul"
			},
			{
				Namespace  = "default"
				Partition  = "default"
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

const testAccConsulConfigEntryEE_ServiceConfigL7b = `
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

resource "consul_config_entry" "jwt_provider" {
	name = "okta"
	kind = "jwt-provider"

	config_json = jsonencode({
		Issuer = "test-issuer"
		JSONWebKeySet = {
			Remote = {
				URI = "https://127.0.0.1:9091"
				FetchAsynchronously = true
			}
		}
	})
}

resource "consul_config_entry" "service_intentions" {
	name = consul_config_entry.sd.name
	kind = "service-intentions"

	config_json = jsonencode({
		Sources = [
			{
				Namespace   = "default"
				Partition   = "default"
				Name        = "admin-dashboard"
				Permissions = [
					{
						Action = "allow"
						HTTP = {
							Methods    = ["GET", "PUT", "POST", "DELETE", "HEAD"]
							PathPrefix = "/v2"
						}
						JWT = {
							Providers = [
								{ name = consul_config_entry.jwt_provider.name }
							]
						}
					}
				],
				Precedence = 9
				Type = "consul"
			},
			{
				Namespace = "default"
				Partition = "default"
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

const testAccConsulConfigEntryEE_ServiceConfigL7gRPC = `
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
				Namespace  = "default"
				Partition  = "default"
				Name       = "frontend-web"
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
				Namespace  = "default"
				Partition  = "default"
				Name       = "support-portal"
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

const testAccConsulConfigEntryEE_ServiceConfigL7Mixed = `
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
				Namespace  = "default"
				Partition  = "default"
				Action     = "deny"
				Name       = "hackathon-project"
				Precedence = 9
				Type       = "consul"
			},
			{
				Namespace  = "default"
				Partition  = "default"
				Action     = "allow"
				Name       = "web"
				Precedence = 9
				Type       = "consul"
			},
			{
				Namespace  = "default"
				Partition  = "default"
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

const testAccConsulConfigEntryEE_DefaultNamespace = `
resource "consul_config_entry" "foo" {
	name      = "foo"
	kind      = "service-defaults"
	namespace = "default"

	config_json = jsonencode({
		Expose           = {}
		MeshGateway      = {}
		TransparentProxy = {}
		Protocol         = "https"
	})
}
`

const testAccConsulConfigEntryEE_Namespace = `
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
	  })
}
`

const TestAccConsulConfigEntryEE_mesh = `
`

const TestAccConsulConfigEntryEE_exportedServicesCE = `
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

const TestAccConsulConfigEntryEE_exportedServicesEE = `
resource "consul_admin_partition" "test" {
	name = "test"
}

resource "consul_config_entry" "exported_services" {
	name = "test"
	kind = "exported-services"
	partition = consul_admin_partition.test.name

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

const TestAccConsulConfigEntryEE_meshCE = `
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

const TestAccConsulConfigEntryEE_meshEE = `
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

const TestAccConsulConfigEntryEE_jwtRemote = `
resource "consul_config_entry" "jwt_provider" {
	name = "okta"
	kind = "jwt-provider"

	config_json = jsonencode({
		Issuer = "test-issuer"
		JSONWebKeySet = {
			Remote = {
				URI = "https://127.0.0.1:9091"
				FetchAsynchronously = true
			}
		}
		Forwarding = {
			HeaderName = "test-token"
		}
	})
}
`

const TestAccConsulConfigEntryEE_jwtLocal = `
resource "consul_config_entry" "jwt_provider" {
	name = "auth0"
	kind = "jwt-provider"

	config_json = jsonencode({
		Issuer = "auth0-issuer"
		JSONWebKeySet = {
			Local = {
        JWKS = "eyJrZXlzIjogW3sKICAiY3J2IjogIlAtMjU2IiwKICAia2V5X29wcyI6IFsKICAgICJ2ZXJpZnkiCiAgXSwKICAia3R5IjogIkVDIiwKICAieCI6ICJXYzl1WnVQYUI3S2gyRk1jOXd0SmpSZThYRDR5VDJBWU5BQWtyWWJWanV3IiwKICAieSI6ICI2OGhSVEppSk5Pd3RyaDRFb1BYZVZuUnVIN2hpU0RKX2xtYmJqZkRmV3EwIiwKICAiYWxnIjogIkVTMjU2IiwKICAidXNlIjogInNpZyIsCiAgImtpZCI6ICJhYzFlOGY5MGVkZGY2MWM0MjljNjFjYTA1YjRmMmUwNyIKfV19"
    	}
		}
	})
}
`
