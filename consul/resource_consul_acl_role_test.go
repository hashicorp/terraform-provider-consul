package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulACLRole_basic(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLRoleConfigBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_role.test", "description", "bar"),
					resource.TestCheckResourceAttrSet("consul_acl_role.test", "id"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "name", "foo"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "namespace", ""),
					resource.TestCheckResourceAttr("consul_acl_role.test", "node_identities.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "policies.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "service_identities.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "service_identities.3690720679.datacenters.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "service_identities.3690720679.service_name", "foo"),
				),
			},
			{
				Config: testResourceACLRoleConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_role.test", "description", ""),
					resource.TestCheckResourceAttrSet("consul_acl_role.test", "id"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "name", "baz"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "namespace", ""),
					resource.TestCheckResourceAttr("consul_acl_role.test", "node_identities.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "node_identities.0.datacenter", "world"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "node_identities.0.node_name", "hello"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "policies.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "service_identities.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "service_identities.2708159462.datacenters.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_role.test", "service_identities.2708159462.service_name", "bar"),
				),
			},
			{
				Config:            testResourceACLRoleConfigUpdate,
				ResourceName:      "consul_acl_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsulACLRole_NamespaceCE(t *testing.T) {
	startTestServer(t)

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
	startTestServer(t)

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
	name        = "test-role"
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
	name      = "baz"

	service_identities {
		service_name = "bar"
	}

	node_identities {
		node_name = "hello"
		datacenter = "world"
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
