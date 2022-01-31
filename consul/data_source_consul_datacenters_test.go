package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulDatacenters_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulDatacentersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_datacenters.read", "datacenters.#", "1"),
					testAccCheckDataSourceValue("data.consul_datacenters.read", "datacenters.0", "dc1"),
				),
			},
		},
	})
}

const testAccDataConsulDatacentersConfig = `
data "consul_datacenters" "read" {}
`
