package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataConsulAutopilotHealth_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataAutopilotHealth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "healthy", "true"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "failure_tolerance", "0"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.#", "1"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.id", "<any>"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.name", "<any>"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.address", "<any>"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.serf_status", "alive"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.version", "<any>"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.leader", "true"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.last_contact", "<any>"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.last_term", "<any>"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.last_index", "<any>"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.healthy", "true"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.voter", "true"),
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.0.stable_since", "<any>"),
				),
			},
		},
	})
}

func TestAccDataConsulAutopilotHealth_config(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataAutopilotHealthDatacenter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_autopilot_health.read", "servers.#", "1"),
				),
			},
		},
	})
}

func TestAccDataConsulAutopilotHealth_wrongDatacenter(t *testing.T) {
	re, err := regexp.Compile("No path to datacenter")
	if err != nil {
		t.Fatalf("err: %#v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccDataAutopilotHealthWrongDatacenter,
				ExpectError: re,
			},
		},
	})
}

const testAccDataAutopilotHealth = `
data "consul_autopilot_health" "read" {}

output "health" {
  value = "${data.consul_autopilot_health.read.healthy}"
}
`

const testAccDataAutopilotHealthDatacenter = `
data "consul_autopilot_health" "read" {
	datacenter = "dc1"
}
`

const testAccDataAutopilotHealthWrongDatacenter = `
data "consul_autopilot_health" "read" {
	datacenter = "wrong_datacenter"
}
`
