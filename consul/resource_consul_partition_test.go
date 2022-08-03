package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var partitionEnterpriseFeature = regexp.MustCompile("(?i)Consul Enterprise feature")

func TestAccConsulPartition_FailOnCommunityEdition(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulPartition,
				ExpectError: regexp.MustCompile("failed to create partition: Unexpected response code: 404"),
			},
		},
	})
}

func TestAccConsulPartition(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulPartition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_partition.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_partition.test", "description", "test partition"),
				),
			},
			{
				Config: testAccConsulPartition_Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_partition.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_partition.test", "description", "updated description"),
				),
			},
			{
				ResourceName:      "consul_partition.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccConsulPartition = `
resource "consul_partition" "test" {
	name        = "test"
	description = "test partition"
}
`

const testAccConsulPartition_Update = `
resource "consul_partition" "test" {
  name        = "test"
  description = "updated description"
}`
