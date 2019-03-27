package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccCheckConsulACLAgentTokenDestroy(s *terraform.State) error {
	client, err := getMasterClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "consul_acl_agent_token" {
			continue
		}
		aclEntry, _, err := client.ACL().Info(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if aclEntry != nil {
			return fmt.Errorf("ACL agent token %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccConsulACLAgentToken_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLAgentTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLAgentTokenConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_agent_token.test", "description", "test"),
					resource.TestCheckResourceAttr("consul_acl_agent_token.test", "policies.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_agent_token.test", "local", "true"),
					resource.TestCheckResourceAttrSet("consul_acl_agent_token.test", "token"),
				),
			},
		},
	})
}

const testResourceACLAgentTokenConfigBasic = `
resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node_prefix \"\" { policy = \"write\" } service_prefix \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_agent_token" "test" {
	description = "test"
	policies = ["${consul_acl_policy.test.name}"]
	local = true
}`
