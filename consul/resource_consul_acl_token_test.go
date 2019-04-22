package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccCheckConsulACLTokenDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*consulapi.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "consul_acl" {
			continue
		}
		aclToken, _, err := client.ACL().TokenRead(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if aclToken != nil {
			return fmt.Errorf("ACL token %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccConsulACLToken_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_token.test", "description", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.1785148924", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "local", "true"),
				),
			},
			{
				Config: testResourceACLTokenConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_token.test", "description", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.111830242", "test2"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "local", "true"),
				),
			},
		},
	})
}

func TestAccConsulACLToken_import(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("bad state: %s", s)
		}
		v, ok := s[0].Attributes["description"]
		if !ok || v != "test" {
			return fmt.Errorf("bad description: %s", s)
		}
		v, ok = s[0].Attributes["policies.#"]
		if !ok || v != "1" {
			return fmt.Errorf("bad policies: %s", s)
		}
		v, ok = s[0].Attributes["local"]
		if !ok || v != "true" {
			return fmt.Errorf("bad local: %s", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenConfigBasic,
			},
			{
				ResourceName:     "consul_acl_token.test",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

const testResourceACLTokenConfigBasic = `
resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = ["${consul_acl_policy.test.name}"]
	local = true
}`

const testResourceACLTokenConfigUpdate = `
// Using another resource to force the update of consul_acl_token
resource "consul_acl_policy" "test2" {
	name = "test2"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = ["${consul_acl_policy.test2.name}"]
	local = true
}`
