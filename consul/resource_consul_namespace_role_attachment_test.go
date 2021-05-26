package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulNamespaceRoleAttachment(t *testing.T) {
	testRole := func(name string) func(*terraform.State) error {
		return func(s *terraform.State) error {
			client := getTestClient(testAccProvider.Meta())
			namespace, _, err := client.Namespaces().Read("testroleattachment", nil)
			if err != nil {
				return fmt.Errorf("failed to read namespace testroleattachment: %s", err)
			}
			if namespace == nil {
				return fmt.Errorf("namespace testroleattachment not found")
			}
			if len(namespace.ACLs.RoleDefaults) != 1 {
				return fmt.Errorf("wrong number of roles: %d", len(namespace.ACLs.RoleDefaults))
			}
			if namespace.ACLs.RoleDefaults[0].Name != name {
				return fmt.Errorf("wrong role, expected %q, found %q", name, namespace.ACLs.RoleDefaults[0].Name)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceNamespaceRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_namespace_role_attachment.test", "namespace", "testroleattachment"),
					resource.TestCheckResourceAttr("consul_namespace_role_attachment.test", "role", "role"),
					testRole("role"),
				),
			},
			{
				Config: testResourceNamespaceRoleConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_namespace_role_attachment.test", "namespace", "testroleattachment"),
					resource.TestCheckResourceAttr("consul_namespace_role_attachment.test", "role", "role2"),
					testRole("role2"),
				),
			},
			{
				Config: testResourceNamespaceRoleConfigUpdate,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "consul_namespace_role_attachment.test",
			},
		},
	})
}

const testResourceNamespaceRoleConfig = `
resource "consul_namespace" "test" {
	name = "testroleattachment"

	lifecycle {
		ignore_changes = [role_defaults]
	}
}

resource "consul_acl_role" "test" {
	name = "role"

	service_identities {
        service_name = "foo"
    }
}

resource "consul_namespace_role_attachment" "test" {
	namespace = consul_namespace.test.name
	role      = consul_acl_role.test.name
}
`

const testResourceNamespaceRoleConfigUpdate = `
resource "consul_namespace" "test" {
	name = "testroleattachment"

	lifecycle {
		ignore_changes = [role_defaults]
	}
}

resource "consul_acl_role" "test2" {
	name = "role2"

	service_identities {
        service_name = "foo"
    }
}

resource "consul_namespace_role_attachment" "test" {
	namespace = consul_namespace.test.name
	role      = consul_acl_role.test2.name
}
`
