package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckConsulACLTokenRoleAttachmentDestroy(s *terraform.State) error {
	client := getTestClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "consul_acl_token_role_attachment" {
			continue
		}
		tokenID, roleName, err := parseTwoPartID(rs.Primary.ID, "token", "role")
		if err != nil {
			return fmt.Errorf("Invalid ACL token role attachment id '%q'", rs.Primary.ID)
		}
		aclToken, _, _ := client.ACL().TokenRead(tokenID, nil)
		if aclToken != nil {
			for _, role := range aclToken.Roles {
				if role.Name == roleName {
					return fmt.Errorf("ACL token role attachment %q still exists", rs.Primary.ID)
				}
			}
		}
	}
	return nil
}

func testAccCheckTokenRoleName(s *terraform.State) error {
	rs, ok := s.RootModule().Resources["consul_acl_token.test"]
	if !ok {
		return fmt.Errorf("Not Found: consul_acl_token.test")
	}

	tokenID := rs.Primary.Attributes["id"]
	if tokenID == "" {
		return fmt.Errorf("No token ID is set")
	}

	client := getTestClient(testAccProvider.Meta())
	_, _, err := client.ACL().TokenRead(tokenID, nil)
	if err != nil {
		return fmt.Errorf("Unable to retrieve token %q", tokenID)
	}

	// Make sure the role has the same token_id
	rs, ok = s.RootModule().Resources["consul_acl_token_role_attachment.test"]
	if !ok {
		return fmt.Errorf("Not Found: consul_acl_token_role_attachment.test")
	}

	roleTokenID := rs.Primary.Attributes["token_id"]
	if roleTokenID == "" {
		return fmt.Errorf("No role token_id is set")
	}

	if roleTokenID != tokenID {
		return fmt.Errorf("%s != %s", roleTokenID, tokenID)
	}

	return nil
}

func TestAccConsulACLTokenRoleAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,

		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLTokenRoleAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenRoleAttachmentConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTokenRoleName,
					resource.TestCheckResourceAttr("consul_acl_token_role_attachment.test", "role", "test"),
				),
			},
			{
				Config: testResourceACLTokenRoleAttachmentConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTokenRoleName,
					resource.TestCheckResourceAttr("consul_acl_token_role_attachment.test", "role", "test2"),
				),
			},
			{
				Config: testResourceACLTokenRoleAttachmentConfigUpdate,
			},
		},
	})
}

func TestAccConsulACLTokenRoleAttachment_import(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("bad state: %s", s)
		}
		_, ok := s[0].Attributes["token_id"]
		if !ok {
			return fmt.Errorf("bad token_id: %s", s)
		}
		_, ok = s[0].Attributes["role"]
		if !ok {
			return fmt.Errorf("bad role: %s", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenRoleAttachmentConfigBasic,
			},
			{
				ResourceName:     "consul_acl_token_role_attachment.test",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

const testResourceACLTokenRoleAttachmentConfigBasic = `
resource "consul_acl_role" "test" {
	name = "test"

	service_identities {
        service_name = "foo"
    }
}

resource "consul_acl_token" "test" {
	description = "test"
	local = true

	lifecycle {
		ignore_changes = ["roles"]
	}
}

resource "consul_acl_token_role_attachment" "test" {
    token_id = "${consul_acl_token.test.id}"
    role  = "${consul_acl_role.test.name}"
}
`

const testResourceACLTokenRoleAttachmentConfigUpdate = `
// Using another resource to force the update of consul_acl_token
resource "consul_acl_role" "test2" {
	name      = "test2"

	service_identities {
        service_name = "bar"
    }
}

resource "consul_acl_token" "test" {
	description = "test"
	roles = []

	lifecycle {
		ignore_changes = ["roles"]
	}
}

resource "consul_acl_token_role_attachment" "test" {
    token_id = "${consul_acl_token.test.id}"
    role = "${consul_acl_role.test2.name}"
}`
