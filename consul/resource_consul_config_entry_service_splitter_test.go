// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntryServiceSplitterTest(t *testing.T) {
	providers, _ := startTestServer(t)

	var config string
	if serverIsConsulCommunityEdition(t) {
		config = testConsulConfigEntryServiceSplitter("", "")
	} else {
		config = testConsulConfigEntryServiceSplitter("default", "default")
	}

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "id", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "meta.%", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "meta.key", "value"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "name", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.#", "3"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.request_headers.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.request_headers.0.add.%", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.request_headers.0.remove.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.request_headers.0.set.%", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.request_headers.0.set.x-web-version", "from-v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.response_headers.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.response_headers.0.add.%", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.response_headers.0.remove.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.response_headers.0.set.%", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.response_headers.0.set.x-web-version", "to-v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.service", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.service_subset", "v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.0.weight", "80"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.request_headers.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.request_headers.0.add.%", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.request_headers.0.remove.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.request_headers.0.set.%", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.request_headers.0.set.x-web-version", "from-v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.response_headers.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.response_headers.0.add.%", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.response_headers.0.remove.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.response_headers.0.set.%", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.response_headers.0.set.x-web-version", "to-v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.service", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.service_subset", "v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.1.weight", "10"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.2.namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.2.partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.2.request_headers.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.2.response_headers.#", "0"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.2.service", "web"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.2.service_subset", "v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_splitter.foo", "splits.2.weight", "10"),
				),
			},
			{
				Config:            config,
				ResourceName:      "consul_config_entry_service_splitter.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"splits.2.request_headers.#",
					"splits.2.response_headers.#",
				},
			},
		},
	})
}

func testConsulConfigEntryServiceSplitter(namespace, partition string) string {
	return fmt.Sprintf(`
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
	name      = consul_config_entry.service_resolver.name
	namespace = "%s"
	partition = "%s"

	meta = {
		key = "value"
	}

	splits {
		weight         = 80
		service        = "web"
		service_subset = "v1"

		request_headers {
			set = {
				"x-web-version" = "from-v1"
			}
		}

		response_headers {
			set = {
				"x-web-version" = "to-v1"
			}
		}
	}

	splits {
		weight         = 10
		service        = "web"
		service_subset = "v2"

		request_headers {
			set = {
				"x-web-version" = "from-v2"
			}
		}

		response_headers {
			set = {
				"x-web-version" = "to-v2"
			}
		}
	}

	splits {
		weight         = 10
		service        = "web"
		service_subset = "v2"
	}
}
`, namespace, partition)
}
