// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccConsulServiceSplitterConfigCEEntryTest(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testConsulServiceSplitterConfigEntryCE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service_splitter_config_entry.foo", "name", "service-splitter"),
					resource.TestCheckResourceAttr("consul_service_splitter_config_entry.foo", "meta.key", "value"),
					resource.TestCheckResourceAttr("consul_service_splitter_config_entry.foo", "splits.#", "2"),
					resource.TestCheckResourceAttr("consul_service_splitter_config_entry.foo", "splits.0.weight", "90"),
					resource.TestCheckResourceAttr("consul_service_splitter_config_entry.foo", "splits.0.service_subset", "v1"),
					resource.TestCheckResourceAttr("consul_service_splitter_config_entry.foo", "splits.1.weight", "10"),
					resource.TestCheckResourceAttr("consul_service_splitter_config_entry.foo", "splits.1.service_subset", "v2"),
				),
			},
		},
	})
}

const testConsulServiceSplitterConfigEntryCE = `
resource "consul_service_splitter_config_entry" "foo" {
	name      = "web" 
	meta      = {
		key = "value"
	}
	splits {
		weight  = 90                   
		service_subset  = "v1"                
		service = "web"
	}
	splits {
		weight  = 10
		service = "web"
		service_subset  = "v2"
	}
}
`
