package consul

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataConsulServiceHealth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataConsulServiceHealth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_service_health.consul", "service", "consul"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "datacenter", "dc1"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "near", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "tag", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "node_meta.%", "1"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "node_meta.consul-network-segment", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "passing", "false"),

					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.#", "1"),

					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.node_id", "<any>"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.node_name", "<any>"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.node_address", "<any>"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.node_datacenter", "dc1"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.node_tagged_addresses.%", "2"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.node_meta.%", "1"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.node_meta.consul-network-segment", ""),

					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.service_id", "consul"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.service_name", "consul"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.service_tags.#", "0"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.service_address", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.service_port", "8300"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.service_meta.%", "0"),

					testAccCheckDataSourceValue("data.consul_service_health.consul", "nodes.0.checks.#", "1"),
				),
			},
		},
	})
}

const testAccDataConsulServiceHealth = `
data "consul_service_health" "consul" {
	service = "consul"

	node_meta {
		// Consul development server has this node meta information
		consul-network-segment = ""
	}
}
`
