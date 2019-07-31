package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccCheckConsulACLTokenPolicyAttachmentDestroy(s *terraform.State) error {
	client := getClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "consul_acl_token_policy_attachment" {
			continue
		}
		tokenID, policyName, err := parseTwoPartID(rs.Primary.ID)
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

func testAccCheckTokenExists(n string, tokenID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		tokenID := rs.Primary.Attributes["id"]
		if tokenID == "" {
			return fmt.Errorf("No token ID is set")
		}

		client := getClient(testAccProvider.Meta())
		_, _, err := client.ACL().TokenRead(tokenID, nil)
		if err != nil {
			return fmt.Errorf("Unable to retrieve token %q", tokenID)
		}

		return nil
	}
}

func TestAccConsulACLTokenPolicyAttachment_basic(t *testing.T) {
	var tokenID string

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,

		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLTokenPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenPolicyAttachmentConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTokenExists("consul_acl_token.test", &tokenID),
					resource.TestCheckResourceAttr("consul_acl_token_policy_attachment.test", "token_id", tokenID),
					resource.TestCheckResourceAttr("consul_acl_token_policy_attachment.test", "policy", "test"),
				),
			},
			{
				Config: testResourceACLTokenPolicyAttachmentConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTokenExists("consul_acl_token.test", &tokenID),
					resource.TestCheckResourceAttr("consul_acl_token_policy_attachment.test", "token_id", tokenID),
					resource.TestCheckResourceAttr("consul_acl_token_policy_attachment.test", "policy", "test"),
				),
			},
			{
				Config: testResourceACLTokenPolicyAttachmentConfigUpdate,
			},
		},
	})
}

func TestAccConsulACLTokenPolicyAttachment_import(t *testing.T) {
	var tokenID string

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("bad state: %s", s)
		}
		v, ok := s[0].Attributes["token_id"]
		if !ok || v != tokenID {
			return fmt.Errorf("bad token_id: %s", s)
		}
		v, ok = s[0].Attributes["policy"]
		if !ok || v != "test" {
			return fmt.Errorf("bad policy: %s", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenPolicyAttachmentConfigBasic,
				Check:  testAccCheckTokenExists("", &tokenID),
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
	name = "test"
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
