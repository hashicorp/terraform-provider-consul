package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckConsulACLPolicyDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
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

}

func TestAccConsulACLPolicy_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulACLPolicyDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testResourceACLPolicyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_policy.test", "name", "test-policy"),
					resource.TestCheckResourceAttr("consul_acl_policy.test", "rules", "node_prefix \"\" { policy = \"read\" }"),
					resource.TestCheckResourceAttr("consul_acl_policy.test", "datacenters.#", "1"),
				),
			},
			{
				Config: testResourceACLPolicyConfigBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_policy.test", "name", "test-policy"),
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
	providers, _ := startTestServer(t)

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
		Providers: providers,
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
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
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
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
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
	name = "test-policy"
	rules = "node_prefix \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}`

const testResourceACLPolicyConfigBasicUpdate = `
resource "consul_acl_policy" "test" {
	name = "test-policy"
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
