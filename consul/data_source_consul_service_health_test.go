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
					testAccCheckDataSourceValue("data.consul_service_health.consul", "name", "consul"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "datacenter", "dc1"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "near", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "tag", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "node_meta.%", "1"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "node_meta.consul-network-segment", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "passing", "true"),

					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.#", "1"),

					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.node.0.id", "<any>"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.node.0.name", "<any>"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.node.0.address", "<any>"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.node.0.datacenter", "dc1"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.node.0.tagged_addresses.%", "2"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.node.0.meta.%", "1"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.node.0.meta.consul-network-segment", ""),

					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.service.0.id", "consul"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.service.0.name", "consul"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.service.0.tags.#", "0"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.service.0.address", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.service.0.port", "8300"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.service.0.meta.%", "0"),

					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.#", "1"),

					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.0.node", "<any>"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.0.id", "serfHealth"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.0.name", "Serf Health Status"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.0.status", "passing"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.0.notes", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.0.output", "Agent alive and reachable"),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.0.service_id", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.0.service_name", ""),
					testAccCheckDataSourceValue("data.consul_service_health.consul", "results.0.checks.0.service_tags.#", "0"),
				),
			},
		},
	})
}

const testAccDataConsulServiceHealth = `
data "consul_service_health" "consul" {
	name = "consul"

	node_meta {
		// Consul development server has this node meta information
		consul-network-segment = ""
	}
}
`
