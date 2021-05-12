package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulACLRole_basic(t *testing.T) {
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

func TestAccConsulACLRole_import(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("bad state: %s", s)
		}
		v, ok := s[0].Attributes["name"]
		if !ok || v != "foo" {
			return fmt.Errorf("bad name: %s", s)
		}
		v, ok = s[0].Attributes["description"]
		if !ok || v != "bar" {
			return fmt.Errorf("bad description: %s", s)
		}
		v, ok = s[0].Attributes["policies.#"]
		if !ok || v != "1" {
			return fmt.Errorf("bad policies: %s", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLRoleConfigBasic,
			},
			{
				ResourceName:     "consul_acl_role.test",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

func TestAccConsulACLRole_NamespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testResourceACLRoleNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccConsulACLRole_NamespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLRoleNamespaceEE,
			},
		},
	})
}

func testRoleDestroy(s *terraform.State) error {
	ACL := getTestClient(testAccProvider.Meta()).ACL()
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
	name        = "test"
	rules       = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role" "test" {
	name        = "foo"
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

const testResourceACLRoleNamespaceCE = `
resource "consul_acl_role" "test" {
  name      = "test"
  namespace = "test-role"
}
`

const testResourceACLRoleNamespaceEE = `
resource "consul_namespace" "test" {
  name = "test-role"
}

resource "consul_acl_role" "test" {
  name      = "test-role"
  namespace = consul_namespace.test.name
}
`
