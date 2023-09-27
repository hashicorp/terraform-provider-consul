// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntryServiceSplitterEETest(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testConsulConfigEntryServiceSplitterEE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "namespace", "ns"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "partition", "pt"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "meta.key", "value"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.#", "2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.weight", "90"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.service", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.service_subset", "v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.request_headers.0.set.x-web-version", "from-v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.response_headers.0.set.x-web-version", "to-v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.weight", "10"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.service_subset", "v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.service", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.request_headers.0.set.x-web-version", "from-v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.response_headers.0.set.x-web-version", "to-v2"),
				),
			},
		},
	})
}

const testConsulConfigEntryServiceSplitterEE = `
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

resource "consul_config_entry_service_splitter" "foo" {
	name = consul_config_entry.service_resolver.name
	namespace = "ns"
	partition = "pt"
	meta = {
		key = "value"
	}
	splits {
		weight  = 90
		service_subset  = "v1"
		service = "web"
		request_headers {
			set = {
				"x-web-version": "from-v1"
			}
		}
		response_headers {
			set = {
				"x-web-version": "to-v1"
			}
		}
	}
	splits {
		weight  = 10
		service = "web"
		service_subset  = "v2"
		request_headers {
			set = {
				"x-web-version": "from-v2"
			}
		}
		response_headers {
			set = {
				"x-web-version": "to-v2"
			}
		}
	}
}
`
