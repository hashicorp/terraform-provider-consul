package consul

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulAdminParition_EntBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck: func() {
			testAccPreCheck(t)
			skipTestOnConsulCommunityEdition(t)
		},
		CheckDestroy: testAccCheckConsulACLTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulAdminPartitionBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_admin_partition.test", "name", "hello"),
					resource.TestCheckResourceAttr("consul_admin_partition.test", "description", "world"),
				),
			},
			{
				PreConfig: func() {
					client := getTestClient(testAccProvider.Meta())
					partitions := client.Partitions()
					if _, err := partitions.Delete(context.TODO(), "hello", nil); err != nil {
						t.Fatalf("failed to remove partition: %s", err)
					}
				},
				Config: testAccConsulAdminPartitionBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_admin_partition.test", "name", "hello"),
					resource.TestCheckResourceAttr("consul_admin_partition.test", "description", "world"),
				),
			},
			{
				ImportState:       true,
				ResourceName:      "consul_admin_partition.test",
				ImportStateVerify: true,
			},
		},
	})
}

const testAccConsulAdminPartitionBasic = `
resource "consul_admin_partition" "test" {
	name        = "hello"
	description = "world"
}
`
