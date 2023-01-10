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
			roleID, policyID, err := parseTwoPartID(rs.Primary.ID, "role_id", "policy_id")
			if err != nil {
				return fmt.Errorf("Invalid ACL role attachment id! '%q'", rs.Primary.ID)
			}
			aclRole, _, _ := client.ACL().RoleRead(roleID, nil)
			if aclRole != nil {
				for _, policy := range aclRole.Policies {
					if policy.ID == policyID {
						return fmt.Errorf("ACL role policy attachment %q still exists", rs.Primary.ID)
					}
				}
			}
		}
		return nil
	}
}

func testAccCheckRolePolicyName(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["consul_acl_role_policy_attachment.test"]
		if !ok {
			return fmt.Errorf("Not Found: consul_acl_role_policy_attachment.test")
		}

		roleIDFromState, policyIDFromState, err := parseTwoPartID(rs.Primary.ID, "role", "policy")
		if err != nil {
			return err
		}

		_, _, err = client.ACL().RoleRead(roleIDFromState, nil)
		if err != nil {
			return fmt.Errorf("Unable to retrieve role %q", roleIDFromState)
		}

		rs, ok = s.RootModule().Resources["consul_acl_role_policy_attachment.test"]
		if !ok {
			return fmt.Errorf("Not Found: consul_acl_role_policy_attachment.test")
		}

		roleIDFromResource := rs.Primary.Attributes["role_id"]
		if roleIDFromResource == "" {
			return fmt.Errorf("No role is set in attachment")
		}

		if roleIDFromResource != roleIDFromState {
			return fmt.Errorf("%s != %s", roleIDFromResource, roleIDFromState)
		}

		policyIDFromResource := rs.Primary.Attributes["policy_id"]
		if policyIDFromResource == "" {
			return fmt.Errorf("No policy is set in attachment")
		}

		if policyIDFromResource != policyIDFromState {
			return fmt.Errorf("%s != %s", policyIDFromState, policyIDFromResource)
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
					testAccCheckRolePolicyName(client),
				),
			},
			{
				Config: testResourceACLRolePolicyAttachmentConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRolePolicyName(client),
				),
			},
			{
				Config: testResourceACLRolePolicyAttachmentConfigUpdate,
			},
		},
	})
}

func TestAccConsulACLRolePolicyAttachment_import(t *testing.T) {
	providers, _ := startTestServer(t)

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("bad state: %s", s)
		}
		_, ok := s[0].Attributes["role_id"]
		if !ok {
			return fmt.Errorf("bad role: %s", s)
		}
		_, ok = s[0].Attributes["policy_id"]
		if !ok {
			return fmt.Errorf("bad policy: %s", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLRolePolicyAttachmentConfigBasic,
			},
			{
				ResourceName:     "consul_acl_role_policy_attachment.test",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

const testResourceACLRolePolicyAttachmentConfigBasic = `
resource "consul_acl_policy" "default" {
	name = "default"
	rules = "service \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role" "test" {
	name = "test"

	service_identities {
        service_name = "foo"
    }

	policies = ["${consul_acl_policy.default.id}"]

	lifecycle {
		ignore_changes = ["policies"]
	}	
}

resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role_policy_attachment" "test" {
    role_id   = "${consul_acl_role.test.id}"
    policy_id = "${consul_acl_policy.test.id}"
}
`

const testResourceACLRolePolicyAttachmentConfigUpdate = `
resource "consul_acl_policy" "default" {
	name = "default"
	rules = "service \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role" "test" {
	name = "test"

	service_identities {
        service_name = "foo"
    }

	policies = ["${consul_acl_policy.default.id}"]

	lifecycle {
		ignore_changes = ["policies"]
	}	
}

resource "consul_acl_policy" "test2" {
	name = "test2"
	rules = "node \"\" { policy = \"write\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role_policy_attachment" "test" {
    role_id   = "${consul_acl_role.test.id}"
    policy_id = "${consul_acl_policy.test2.id}"
}
`
