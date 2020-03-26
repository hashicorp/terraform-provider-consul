package consul

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckConsulACLTokenDestroy(s *terraform.State) error {
	client := getClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "consul_acl_token" {
			continue
		}
		aclToken, _, _ := client.ACL().TokenRead(rs.Primary.ID, nil)
		if aclToken != nil {
			return fmt.Errorf("ACL token %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccConsulACLToken_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,

		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_token.test", "description", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "roles.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.1785148924", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "local", "true"),
				),
			},
			{
				Config: testResourceACLTokenConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_token.test", "description", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "roles.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.111830242", "test2"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "local", "false"),
				),
			},
			{
				Config: testResourceACLTokenConfigUpdate,
			},
			{
				Config: testResourceACLTokenConfigRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_token.test", "description", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "roles.1785148924", "test"),
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

func TestAccConsulACLToken_namespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testResourceACLTokenConfigNamespaceCE,
				ExpectError: regexp.MustCompile("Namespaces is a Consul Enterprise feature"),
			},
		},
	})
}

func TestAccConsulACLToken_namespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenConfigNamespaceEE,
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
}`

const testResourceACLTokenConfigRole = `
resource "consul_acl_role" "test" {
    name = "test"
}

resource "consul_acl_token" "test" {
	description = "test"
	roles = [consul_acl_role.test.name]
}`

const testResourceACLTokenConfigNamespaceCE = `
resource "consul_acl_token" "test" {
  description = "test"
  namespace   = "test"
}`

const testResourceACLTokenConfigNamespaceEE = `
resource "consul_namespace" "test" {
  name = "test-token"
}
resource "consul_acl_token" "test" {
  description = "test"
  namespace   = consul_namespace.test.name
}`
