// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckConsulACLRolePolicyAttachmentDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "consul_acl_role_policy_attachment" {
				continue
			}
			roleID, policyName, err := parseTwoPartID(rs.Primary.ID, "role", "policy")
			if err != nil {
				return fmt.Errorf("Invalid role policy attachment id '%q'", rs.Primary.ID)
			}
			role, _, _ := client.ACL().RoleRead(roleID, nil)
			if role != nil {
				for _, iPolicy := range role.Policies {
					if iPolicy.Name == policyName {
						return fmt.Errorf("role policy attachment %q still exists", rs.Primary.ID)
					}
				}
			}
		}
		return nil
	}

}

func testAccCheckRolePolicyID(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["consul_acl_role.test_role"]
		if !ok {
			return fmt.Errorf("Not Found: consul_acl_role.test_role")
		}

		roleID := rs.Primary.Attributes["id"]
		if roleID == "" {
			return fmt.Errorf("No token ID is set")
		}

		_, _, err := client.ACL().RoleRead(roleID, nil)
		if err != nil {
			return fmt.Errorf("Unable to retrieve role %q", roleID)
		}

		// Make sure the policy has then same role_id
		rs, ok = s.RootModule().Resources["consul_acl_role_policy_attachment.test"]
		if !ok {
			return fmt.Errorf("Not Found: consul_acl_role_policy_attachment.test")
		}

		policyTokenID := rs.Primary.Attributes["role_id"]
		if policyTokenID == "" {
			return fmt.Errorf("No policy role_id is set")
		}

		if policyTokenID != roleID {
			return fmt.Errorf("%s != %s", policyTokenID, roleID)
		}

		return nil
	}
}

func TestAccConsulACLRolePolicyAttachment_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulACLRolePolicyAttachmentDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testResourceACLRolePolicyAttachmentConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRolePolicyID(client),
					resource.TestCheckResourceAttr("consul_acl_role_policy_attachment.test", "policy", "test-attachment"),
				),
			},
			{
				Config: testResourceACLRolePolicyAttachmentConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRolePolicyID(client),
					resource.TestCheckResourceAttr("consul_acl_role_policy_attachment.test", "policy", "test2"),
				),
			},
			{
				Config: testResourceACLRolePolicyAttachmentConfigUpdate,
			},
			{
				Config:            testResourceACLRolePolicyAttachmentConfigUpdate,
				ResourceName:      "consul_acl_role_policy_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testResourceACLRolePolicyAttachmentConfigBasic = `
resource "consul_acl_policy" "test_policy" {
	name = "test-attachment"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role" "test_role" {
    name = "test"

    lifecycle {
		ignore_changes = ["policies"]
	}
}

resource "consul_acl_role_policy_attachment" "test" {
    role_id = consul_acl_role.test_role.id
    policy  = consul_acl_policy.test_policy.name
}
`

const testResourceACLRolePolicyAttachmentConfigUpdate = `
// Using another resource to force the update of consul_acl_role
resource "consul_acl_policy" "test2" {
	name = "test2"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role" "test_role" {
    name = "test"

    lifecycle {
		ignore_changes = ["policies"]
	}
}

resource "consul_acl_role_policy_attachment" "test" {
    role_id = consul_acl_role.test_role.id
    policy  = consul_acl_policy.test2.name
}`
