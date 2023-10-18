// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntrySamenessGroupTest(t *testing.T) {
	providers, _ := startTestServer(t)

	t.Run("community-edition", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
			Providers: providers,
			Steps: []resource.TestStep{
				{
					Config:      testConsulConfigEntrySamenessGroup,
					ExpectError: regexp.MustCompile("enterprise-only feature"),
				},
			},
		})
	})

	t.Run("enterprise-edition", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
			Providers: providers,
			Steps: []resource.TestStep{
				{
					Config: testConsulConfigEntrySamenessGroup,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "default_for_failover", "true"),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "id", "test"),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "include_local", "true"),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "members.#", "4"),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "members.0.partition", "store-east"),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "members.0.peer", ""),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "members.1.partition", "inventory-east"),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "members.1.peer", ""),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "members.2.partition", ""),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "members.2.peer", "dc2-store-west"),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "members.3.partition", ""),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "members.3.peer", "dc2-inventory-west"),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "name", "test"),
						resource.TestCheckResourceAttr("consul_config_entry_sameness_group.foo", "partition", ""),
					),
				},
				{
					Config:            testConsulConfigEntrySamenessGroup,
					ResourceName:      "consul_config_entry_sameness_group.foo",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

}

const testConsulConfigEntrySamenessGroup = `
resource "consul_config_entry_sameness_group" "foo" {
	name                 = "test"
	default_for_failover = true
	include_local        = true

	members { partition = "store-east" }
	members { partition = "inventory-east" }
	members { peer      = "dc2-store-west" }
	members { peer      = "dc2-inventory-west" }

}
`
