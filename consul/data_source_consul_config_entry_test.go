// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulConfigEntry_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsulConfigEntryMissing,
				ExpectError: regexp.MustCompile(`failed to read config entry service-defaults/foo: Unexpected response code: 404`),
			},
			{
				Config: testAccDataConsulConfigEntry,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_config_entry.read", "config_json", "{\"Expose\":{},\"MeshGateway\":{},\"Protocol\":\"https\",\"TransparentProxy\":{}}"),
					resource.TestCheckResourceAttr("data.consul_config_entry.read", "id", "service-defaults/foo"),
					resource.TestCheckResourceAttr("data.consul_config_entry.read", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("data.consul_config_entry.read", "name", "foo"),
				),
			},
		},
	})
}

const testAccDataConsulConfigEntry = `
resource "consul_config_entry" "test" {
	name = "foo"
	kind = "service-defaults"

	config_json = jsonencode({
		MeshGateway      = {}
		Protocol         = "http"
		TransparentProxy = {}
	})
}

data "consul_config_entry" "read" {
	name = consul_config_entry.test.name
	kind = consul_config_entry.test.kind
}
`

const testAccDataConsulConfigEntryMissing = `
data "consul_config_entry" "read" {
	name = "foo"
	kind = "service-defaults"
}
`
