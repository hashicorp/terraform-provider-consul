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
					resource.TestCheckResourceAttr("consul_acl_token.test", "local", "true"),
				),
			},
			{
				Config:  testResourceACLTokenConfigBasic,
				Destroy: false,
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
