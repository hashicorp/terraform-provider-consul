package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulDatacenters_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
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

func TestAccDataConsulDatacenters_multipleDatacenters(t *testing.T) {
	providers, _ := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulDatacentersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_datacenters.read", "datacenters.#", "2"),
					testAccCheckDataSourceValue("data.consul_datacenters.read", "datacenters.0", "dc1"),
					testAccCheckDataSourceValue("data.consul_datacenters.read", "datacenters.1", "dc2"),
				),
			},
		},
	})
}

const testAccDataConsulDatacentersConfig = `
data "consul_datacenters" "read" {}
`
