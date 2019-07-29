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

func TestAccDataConsulServiceHealthPassing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataConsulServiceHealthPassingSetup,
				Check:  resource.ComposeTestCheckFunc(),
			},
			resource.TestStep{
				Config: testAccDataConsulServiceHealthPassingFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_service_health.google", "name", "google"),
					testAccCheckDataSourceValue("data.consul_service_health.google", "passing", "false"),
					testAccCheckDataSourceValue("data.consul_service_health.google", "results.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccDataConsulServiceHealthPassingTrue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_service_health.google", "name", "google"),
					testAccCheckDataSourceValue("data.consul_service_health.google", "passing", "true"),
					testAccCheckDataSourceValue("data.consul_service_health.google", "results.#", "1"),
				),
			},
		},
	})
}

const testAccDataConsulServiceHealth = `
data "consul_service_health" "consul" {
	name = "consul"

	node_meta = {
		// Consul development server has this node meta information
		consul-network-segment = ""
	}
}
`

const testAccDataConsulServiceHealthPassingSetup = `
resource "consul_node" "compute1" {
  name    = "compute-google1"
  address = "www.google.com"
}

resource "consul_service" "google1" {
  name    = "google"
  node    = "${consul_node.compute1.name}"
  port    = 80
  tags    = ["tag0"]

  check {
    check_id                          = "service:redis1"
    name                              = "Redis health check"
    status                            = "passing"
    http                              = "https://www.hashicorptest.com"
    tls_skip_verify                   = false
    method                            = "PUT"
    interval                          = "5s"
    timeout                           = "1s"
    deregister_critical_service_after = "30s"

    header {
      name  = "foo"
      value = ["test"]
    }

    header {
      name  = "bar"
      value = ["test"]
    }
  }
}

resource "consul_node" "compute2" {
  name    = "compute-google2"
  address = "www.google.com"
}

resource "consul_service" "google2" {
  name    = "google"
  node    = "${consul_node.compute2.name}"
  port    = 80
  tags    = ["tag1"]

  check {
    check_id                          = "service:redis1"
    name                              = "Redis health check"
    status                            = "critical"
    http                              = "https://www.hashicorptest.com"
    tls_skip_verify                   = false
    method                            = "PUT"
    interval                          = "5s"
    timeout                           = "1s"
    deregister_critical_service_after = "30s"

    header {
      name  = "foo"
      value = ["test"]
    }

    header {
      name  = "bar"
      value = ["test"]
    }
  }
}`

const testAccDataConsulServiceHealthPassingFalse = testAccDataConsulServiceHealthPassingSetup + `
data "consul_service_health" "google" {
  name    = "google"
	passing = false
	// wait_for = "10s"
}
`
const testAccDataConsulServiceHealthPassingTrue = testAccDataConsulServiceHealthPassingSetup + `
data "consul_service_health" "google" {
	name     = "google"
	// wait_for = "10s"
}
`
