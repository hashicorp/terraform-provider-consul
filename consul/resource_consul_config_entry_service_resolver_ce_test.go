// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntryServiceResolverCETest(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testConsulConfigEntryServiceResolverCE,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testConsulConfigEntryServiceResolverCE = `
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
