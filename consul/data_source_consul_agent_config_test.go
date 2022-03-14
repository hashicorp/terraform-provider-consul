package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulAgentConfig_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulAgentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_agent_config.example", "datacenter", "dc1"),
					resource.TestCheckResourceAttr("data.consul_agent_config.example", "server", "true"),
					resource.TestCheckResourceAttrSet("data.consul_agent_config.example", "node_name"),
					resource.TestCheckResourceAttrSet("data.consul_agent_config.example", "node_id"),
					resource.TestCheckResourceAttrSet("data.consul_agent_config.example", "version"),
				),
			},
		},
	})
}

const testAccDataConsulAgentConfig = `
data "consul_agent_config" "example" {}
`
