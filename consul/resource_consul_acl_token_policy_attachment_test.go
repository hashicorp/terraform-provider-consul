// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckConsulACLTokenPolicyAttachmentDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "consul_acl_token_policy_attachment" {
				continue
			}
			tokenID, policyName, err := parseTwoPartID(rs.Primary.ID, "token", "policy")
			if err != nil {
				return fmt.Errorf("Invalid ACL token policy attachment id '%q'", rs.Primary.ID)
			}
			aclToken, _, _ := client.ACL().TokenRead(tokenID, nil)
			if aclToken != nil {
				for _, iPolicy := range aclToken.Policies {
					if iPolicy.Name == policyName {
						return fmt.Errorf("ACL token policy attachment %q still exists", rs.Primary.ID)
					}
				}
			}
		}
		return nil
	}

}

func testAccCheckTokenPolicyID(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["consul_acl_token.test"]
		if !ok {
			return fmt.Errorf("Not Found: consul_acl_token.test")
		}

		tokenID := rs.Primary.Attributes["id"]
		if tokenID == "" {
			return fmt.Errorf("No token ID is set")
		}

		_, _, err := client.ACL().TokenRead(tokenID, nil)
		if err != nil {
			return fmt.Errorf("Unable to retrieve token %q", tokenID)
		}

		// Make sure the policy has then same token_id
		rs, ok = s.RootModule().Resources["consul_acl_token_policy_attachment.test"]
		if !ok {
			return fmt.Errorf("Not Found: consul_acl_token_policy_attachment.test")
		}

		policyTokenID := rs.Primary.Attributes["token_id"]
		if policyTokenID == "" {
			return fmt.Errorf("No policy token_id is set")
		}

		if policyTokenID != tokenID {
			return fmt.Errorf("%s != %s", policyTokenID, tokenID)
		}

		return nil
	}
}

func TestAccConsulACLTokenPolicyAttachment_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulACLTokenPolicyAttachmentDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenPolicyAttachmentConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTokenPolicyID(client),
					resource.TestCheckResourceAttr("consul_acl_token_policy_attachment.test", "policy", "test-attachment"),
				),
			},
			{
				Config: testResourceACLTokenPolicyAttachmentConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTokenPolicyID(client),
					resource.TestCheckResourceAttr("consul_acl_token_policy_attachment.test", "policy", "test2"),
				),
			},
			{
				Config: testResourceACLTokenPolicyAttachmentConfigUpdate,
			},
		},
	})
}

func TestAccConsulACLTokenPolicyAttachment_import(t *testing.T) {
	providers, _ := startTestServer(t)

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("bad state: %s", s)
		}
		_, ok := s[0].Attributes["token_id"]
		if !ok {
			return fmt.Errorf("bad token_id: %s", s)
		}
		v, ok := s[0].Attributes["policy"]
		if !ok || v != "test-attachment" {
			return fmt.Errorf("bad policy: %s", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenPolicyAttachmentConfigBasic,
			},
			{
				ResourceName:     "consul_acl_token_policy_attachment.test",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

const testResourceACLTokenPolicyAttachmentConfigBasic = `
resource "consul_acl_policy" "test" {
	name = "test-attachment"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	local = true

	lifecycle {
		ignore_changes = ["policies"]
	}
}

resource "consul_acl_token_policy_attachment" "test" {
    token_id = "${consul_acl_token.test.id}"
    policy = "${consul_acl_policy.test.name}"
}
`

const testResourceACLTokenPolicyAttachmentConfigUpdate = `
// Using another resource to force the update of consul_acl_token
resource "consul_acl_policy" "test2" {
	name = "test2"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = []

	lifecycle {
		ignore_changes = ["policies"]
	}
}

resource "consul_acl_token_policy_attachment" "test" {
    token_id = "${consul_acl_token.test.id}"
    policy = "${consul_acl_policy.test2.name}"
}`
