package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulACLAuthMethod_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,

		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLRoleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_role.test", "name", "foo"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "description", "bar"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "policies.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "service_identities.#", "1"),
				),
			},
			{
				Config: testResourceACLRoleConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_role.test", "name", "baz"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "description", ""),
					resource.TestCheckResourceAttr("consul_acl_role.test", "policies.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "service_identities.#", "1"),
				),
			},
		},
	})
}

func testRoleDestroy(s *terraform.State) error {
	ACL := getClient(testAccProvider.Meta()).ACL()
	qOpts := &consulapi.QueryOptions{}

	role, _, err := ACL.RoleReadByName("baz", qOpts)
	if err != nil {
		return err
	}

	if role != nil {
		return fmt.Errorf("Role 'baz' still exists")
	}

	return nil
}

const testResourceACLRoleConfigBasic = `
resource "consul_acl_policy" "test-read" {
	name = "test"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role" "test" {
	name = "foo"
	description = "bar"

	policies = [
		"${consul_acl_policy.test-read.id}"
	]

	service_identities {
		service_name = "foo"
	}
}`

const testResourceACLRoleConfigUpdate = `
resource "consul_acl_role" "test" {
	name = "baz"

	service_identities {
		service_name = "bar"
	}
}`
