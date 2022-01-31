package consul

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulDatacenters_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if os.Getenv("TEST_REMOTE_DATACENTER") != "" {
				t.Skip("Test skipped. Unset TEST_REMOTE_DATACENTER to run this test.")
			}
		},
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

func TestAccDataConsulDatacenters_multipleDatacenters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccRemoteDatacenterPreCheck(t)
		},
		Providers: testAccProviders,
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
