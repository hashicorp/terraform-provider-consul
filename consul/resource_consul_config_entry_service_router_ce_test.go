// Copyright IBM Corp. 2014, 2025
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
				Config: testConsulConfigEntryServiceRouterCE_Empty,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "id", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.#", "0"),
				),
			},
			{
				Config: testConsulConfigEntryServiceRouterCE,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "id", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "meta.%", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.idle_timeout", "0s"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.num_retries", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.prefix_rewrite", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.request_headers.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.request_timeout", "0s"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.response_headers.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.retry_on.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.retry_on_connect_failure", "false"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.retry_on_status_codes.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.service", "admin"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.service_subset", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.header.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.methods.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.path_exact", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.path_prefix", "/admin"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.path_regex", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.query_param.#", "0"),
				),
			},
			{
				Config: testConsulConfigEntryServiceRouterCE_noMatch,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "id", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "meta.%", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.idle_timeout", "0s"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.num_retries", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.prefix_rewrite", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.request_headers.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.request_timeout", "0s"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.response_headers.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.retry_on.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.retry_on_connect_failure", "false"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.retry_on_status_codes.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.service", "admin"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.0.service_subset", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.#", "0"),
				),
			},
			{
				Config: testConsulConfigEntryServiceRouterCE_noDestination,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "id", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "meta.%", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.destination.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.header.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.methods.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.path_exact", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.path_prefix", "/admin"),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.path_regex", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_router.foo", "routes.0.match.0.http.0.query_param.#", "0"),
				),
			},
		},
	})
}

const (
	testConsulConfigEntryServiceRouterCE_Empty = `
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
	}
`

	testConsulConfigEntryServiceRouterCE = `
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
			service = consul_config_entry.admin_service_defaults.name
		}
	}
}`

	testConsulConfigEntryServiceRouterCE_noMatch = `
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
		destination {
			service = consul_config_entry.admin_service_defaults.name
		}
	}
}`

	testConsulConfigEntryServiceRouterCE_noDestination = `
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
	}
}`
)
