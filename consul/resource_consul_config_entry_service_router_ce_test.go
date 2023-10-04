// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntryServiceRouterCETest(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testConsulConfigEntryServiceRouterCE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.365255188.http.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.365255188.http.3609927257.path_prefix", "/admin"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.1670225453.retry_on_connect_failure", "false"),
				),
			},
		},
	})
}

const testConsulConfigEntryServiceRouterCE = `
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
resource "consul_config_entry_service_router" "foo" {
	name = consul_config_entry.web.name

	routes {
		 match {
			 http {
				  path_prefix = "/admin"
			 }
		 }

		 destination {
			 service   = consul_config_entry.admin_service_defaults.name
		 }
	 }
}
`
