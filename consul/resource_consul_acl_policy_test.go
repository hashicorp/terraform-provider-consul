package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckConsulACLPolicyDestroy(s *terraform.State) error {
	client := getTestClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "consul_acl_policy" {
			continue
		}
		secret, _, err := client.ACL().Info(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if secret != nil {
			return fmt.Errorf("ACL %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccConsulACLPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLPolicyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_policy.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_acl_policy.test", "rules", "node_prefix \"\" { policy = \"read\" }"),
					resource.TestCheckResourceAttr("consul_acl_policy.test", "datacenters.#", "1"),
				),
			},
			{
				Config: testResourceACLPolicyConfigBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_policy.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_acl_policy.test", "rules", "node_prefix \"\" { policy = \"write\" }"),
					resource.TestCheckResourceAttr("consul_acl_policy.test", "datacenters.#", "1"),
				),
			},
			{
				Config:  testResourceACLPolicyConfigBasicUpdate,
				Destroy: true,
			},
		},
	})
}

func TestAccConsulACLPolicy_import(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("bad state: %s", s)
		}
		v, ok := s[0].Attributes["rules"]
		if !ok || v != `node_prefix "" { policy = "read" }` {
			return fmt.Errorf("bad rules: %s", s)
		}
		v, ok = s[0].Attributes["description"]
		if !ok || v != "" {
			return fmt.Errorf("bad description: %s", s)
		}
		v, ok = s[0].Attributes["datacenters.#"]
		if !ok || v != "1" {
			return fmt.Errorf("bad datacenters: %s", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLPolicyConfigBasic,
			},
			{
				ResourceName:     "consul_acl_policy.test",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

func TestAccConsulACLPolicy_NamespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testResourceACLPolicyNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccConsulACLPolicy_NamespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLPolicyNamespaceEE,
			},
		},
	})
}

const testResourceACLPolicyConfigBasic = `
resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node_prefix \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}`

const testResourceACLPolicyConfigBasicUpdate = `
resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node_prefix \"\" { policy = \"write\" }"
	datacenters = [ "dc1" ]
}`

const testResourceACLPolicyNamespaceCE = `
resource "consul_acl_policy" "test" {
  name      = "test"
  rules     = "service \"app\" { policy = \"write\"}"
  namespace = "test-policy"
}
`

const testResourceACLPolicyNamespaceEE = `
resource "consul_namespace" "test" {
  name = "test-policy"
}

resource "consul_acl_policy" "test" {
  name      = "test"
  rules     = "service \"app\" { policy = \"write\"}"
  namespace = consul_namespace.test.name
}
`
