package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntry_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulConfigEntry_ServiceDefaults,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", "{\"MeshGateway\":{},\"Protocol\":\"https\"}"),
				),
			},
			{
				Config: testAccConsulConfigEntry_ProxyDefaults,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.foo", "name", "global"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "kind", "proxy-defaults"),
					resource.TestCheckResourceAttr("consul_config_entry.foo", "config_json", "{\"Config\":{\"foo\":\"bar\"},\"MeshGateway\":{}}"),
				),
			},
			{
				Config: TestAccConsulConfigEntry_ServiceRouter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_router", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_router", "kind", "service-router"),
					resource.TestCheckResourceAttr("consul_config_entry.service_router", "config_json", "{\"Routes\":[{\"Destination\":{\"Service\":\"admin\"},\"Match\":{\"HTTP\":{\"PathPrefix\":\"/admin\"}}}]}"),
				),
			},
			{
				Config: TestAccConsulConfigEntry_ServiceSplitter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "kind", "service-splitter"),
					resource.TestCheckResourceAttr("consul_config_entry.service_splitter", "config_json", "{\"Splits\":[{\"ServiceSubset\":\"v1\",\"Weight\":90},{\"ServiceSubset\":\"v2\",\"Weight\":10}]}"),
				),
			},
			{
				Config: TestAccConsulConfigEntry_ServiceResolver,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "kind", "service-resolver"),
					resource.TestCheckResourceAttr("consul_config_entry.service_resolver", "config_json", "{\"DefaultSubset\":\"v1\",\"Subsets\":{\"v1\":{\"Filter\":\"Service.Meta.version == v1\"},\"v2\":{\"Filter\":\"Service.Meta.version == v2\"}}}"),
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
				ExpectError: regexp.MustCompile("Provider produced inconsistent result after apply: When applying changes to consul_config_entry.foo, provider \"consul\" produced an unexpected new value for was present, but now absent."),
			},
		},
	})
}

const testAccConsulConfigEntry_ServiceDefaults = `
resource "consul_config_entry" "foo" {
	name = "foo"
	kind = "service-defaults"

	config_json = jsonencode({
		MeshGateway = {}
		Protocol    = "https"
	})
}
`

const testAccConsulConfigEntry_ProxyDefaults = `
resource "consul_config_entry" "foo" {
	name = "global"
	kind = "proxy-defaults"

	config_json = jsonencode({
		MeshGateway = {}
		Config      = {
			foo = "bar"
		}
	})
}
`

const TestAccConsulConfigEntry_ServiceRouter = `
resource "consul_config_entry" "web" {
	name = "web"
	kind = "service-defaults"

	config_json = jsonencode({
		MeshGateway = {}
		Protocol    = "http"
	})
}

resource "consul_config_entry" "admin_service_defaults" {
	name = "admin"
	kind = "service-defaults"

	config_json = jsonencode({
		MeshGateway = {}
		Protocol    = "http"
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
					Service = consul_config_entry.admin_service_defaults.name
				}
			}
			# NOTE: a default catch-all will send unmatched traffic to "web"
		]
	})
}
`

const TestAccConsulConfigEntry_ServiceSplitter = `
resource "consul_config_entry" "web" {
	name = "web"
	kind = "service-defaults"

	config_json = jsonencode({
		MeshGateway = {}
		Protocol    = "http"
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

	depends_on = [consul_config_entry.web]
}

resource "consul_config_entry" "service_splitter" {
	kind = "service-splitter"
	name = consul_config_entry.web.name

	config_json = jsonencode({
		Splits = [
			{
				Weight         = 90
				ServiceSubset = "v1"
			},
			{
				Weight         = 10
				ServiceSubset = "v2"
			},
		]
	})

	depends_on = [consul_config_entry.service_resolver]
}
`

const TestAccConsulConfigEntry_ServiceResolver = `
resource "consul_config_entry" "web" {
	name = "web"
	kind = "service-defaults"

	config_json = jsonencode({
		MeshGateway = {}
		Protocol    = "http"
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
