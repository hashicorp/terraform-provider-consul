// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulAdminParition_EntBasic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck: func() {
			skipTestOnConsulCommunityEdition(t)
		},
		CheckDestroy: testAccCheckConsulACLTokenDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulAdminPartitionBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_admin_partition.test", "name", "hello"),
					resource.TestCheckResourceAttr("consul_admin_partition.test", "description", "world"),
					resource.TestCheckResourceAttr("consul_admin_partition.test", "disable_gossip", "true"),
				),
			},
			{
				PreConfig: func() {
					partitions := client.Partitions()
					if _, err := partitions.Delete(context.TODO(), "hello", nil); err != nil {
						t.Fatalf("failed to remove partition: %s", err)
					}
				},
				Config: testAccConsulAdminPartitionBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_admin_partition.test", "name", "hello"),
					resource.TestCheckResourceAttr("consul_admin_partition.test", "description", "world"),
					resource.TestCheckResourceAttr("consul_admin_partition.test", "disable_gossip", "true"),
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
        disable_gossip = true
}
