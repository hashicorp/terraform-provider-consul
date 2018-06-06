package consul

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataConsulCatalogNodes_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataConsulCatalogNodesConfig,
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
func TestAccDataConsulCatalogNodes_alias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataConsulCatalogNodesAlias,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_catalog_nodes.read", "nodes.#", "1"),
				),
			},
		},
	})
}

const testAccDataConsulCatalogNodesConfig = `
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
const testAccDataConsulCatalogNodesAlias = `
data "consul_catalog_nodes" "read" {}
`
