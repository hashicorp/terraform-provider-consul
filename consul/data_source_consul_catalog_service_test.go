package consul

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataConsulCatalogService_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataConsulCatalogServiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "datacenter", "dc1"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.#", "1"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.address", "<all>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.create_index", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.enable_tag_override", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.id", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.modify_index", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.name", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.node_address", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.node_id", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.node_meta.%", "1"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.node_meta.consul-network-segment", ""),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.node_name", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.port", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.tagged_addresses.%", "2"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read", "service.0.tags.#", "0"),
				),
			},
		},
	})
}

func TestAccDataConsulCatalogService_filtered(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataConsulCatalogServiceFilteredConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "datacenter", "dc1"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.#", "1"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.0.address", "192.168.10.10"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.0.create_index", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.0.enable_tag_override", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.0.id", "redis1"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.0.modify_index", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.0.name", "redis"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.0.node_name", "foobar"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.0.port", "<any>"),
					testAccCheckDataSourceValue("data.consul_catalog_service.read_f", "service.0.tags.#", "2"),
				),
			},
		},
	})
}

const testAccDataConsulCatalogServiceConfig = `
data "consul_catalog_service" "read" {
  query_options {
    allow_stale = true
    require_consistent = false
    token = ""
    wait_index = 0
    wait_time = "1m"
  }

  name = "consul"
}
`

const testAccDataConsulCatalogServiceFilteredConfig = `
resource "consul_catalog_entry" "service1" {
  address = "192.168.10.11"
  node    = "foobar_dummy"
  datacenter = "dc1"

  service = {
    id      = "redis2"
    name    = "redis"
    port    = 8000
    tags    = ["v1"]
  }
}

resource "consul_catalog_entry" "service2" {
  address = "192.168.10.10"
  node    = "foobar"
  datacenter = "${consul_catalog_entry.service1.datacenter}"

  service = {
    id      = "redis1"
    name    = "redis"
    port    = 8000
    tags    = ["master", "v1"]
  }
}

data "consul_catalog_service" "read_f" {
  query_options {
    allow_stale = true
    require_consistent = false
    token = ""
    wait_index = 0
    wait_time = "1m"
  }

  name = "redis"
  tag = "master"
  datacenter = "${consul_catalog_entry.service2.datacenter}"
}
`
