// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulNodes_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulNodesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_nodes.read", "nodes.#", "1"),
					testAccCheckDataSourceValue("data.consul_nodes.read", "nodes.0.id", "<any>"),
					testAccCheckDataSourceValue("data.consul_nodes.read", "nodes.0.name", "<any>"),
					testAccCheckDataSourceValue("data.consul_nodes.read", "nodes.0.address", "<any>"),
				),
			},
		},
	})
}

func TestAccDataConsulNodes_alias(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulNodesAlias,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_catalog_nodes.read", "nodes.#", "1"),
				),
			},
		},
	})
}

func TestAccDataConsulNodes_datacenter(t *testing.T) {
	providers, _ := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulNodesDatacenter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_catalog_nodes.read", "nodes.#", "2"),
				),
			},
		},
	})
}

const testAccDataConsulNodesConfig = `
data "consul_nodes" "read" {
  query_options {
    allow_stale = true
    require_consistent = false
    token = ""
    wait_index = 0
    wait_time = "1m"
  }
}
`
const testAccDataConsulNodesAlias = `
data "consul_catalog_nodes" "read" {}
`

const testAccDataConsulNodesDatacenter = `
resource "consul_node" "dc2" {
	datacenter = "dc2"
	name 	   = "dc2"
	address    = "127.0.0.1"
}

data "consul_catalog_nodes" "read" {
	query_options {
		datacenter = consul_node.dc2.datacenter
	}
}
`
