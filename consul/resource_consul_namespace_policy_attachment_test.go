package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulNamespacePolicyAttachment(t *testing.T) {
	providers, client := startTestServer(t)

	testPolicy := func(name string) func(*terraform.State) error {
		return func(s *terraform.State) error {
			namespace, _, err := client.Namespaces().Read("testattachment", nil)
			if err != nil {
				return fmt.Errorf("failed to read namespace testattachment: %s", err)
			}
			if namespace == nil {
				return fmt.Errorf("namespace testattachment not found")
			}
			if len(namespace.ACLs.PolicyDefaults) != 1 {
				return fmt.Errorf("wrong number of policies: %d", len(namespace.ACLs.PolicyDefaults))
			}
			if namespace.ACLs.PolicyDefaults[0].Name != name {
				return fmt.Errorf("wrong policy, expected %q, found %q", name, namespace.ACLs.PolicyDefaults[0].Name)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceNamespacePolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_namespace_policy_attachment.test", "namespace", "testattachment"),
					resource.TestCheckResourceAttr("consul_namespace_policy_attachment.test", "policy", "policy"),
					testPolicy("policy"),
				),
			},
			{
				Config: testResourceNamespacePolicyConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_namespace_policy_attachment.test", "namespace", "testattachment"),
					resource.TestCheckResourceAttr("consul_namespace_policy_attachment.test", "policy", "policy2"),
					testPolicy("policy2"),
				),
			},
			{
				Config: testResourceNamespacePolicyConfigUpdate,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "consul_namespace_policy_attachment.test",
			},
		},
	})
}

const testResourceNamespacePolicyConfig = `
resource "consul_namespace" "test" {
	name = "testattachment"

	lifecycle {
		ignore_changes = [policy_defaults]
	}
}

resource "consul_acl_policy" "test" {
	name        = "policy"
	rules       = <<-RULE
	  node_prefix "" {
		policy = "read"
	  }
	RULE
}

resource "consul_namespace_policy_attachment" "test" {
	namespace = consul_namespace.test.name
	policy    = consul_acl_policy.test.name
}
`

const testResourceNamespacePolicyConfigUpdate = `
resource "consul_namespace" "test" {
	name = "testattachment"

	lifecycle {
		ignore_changes = [policy_defaults]
	}
}

resource "consul_acl_policy" "test2" {
	name        = "policy2"
	rules       = <<-RULE
	  node_prefix "" {
		policy = "read"
	  }
	RULE
}

resource "consul_namespace_policy_attachment" "test" {
	namespace = consul_namespace.test.name
	policy    = consul_acl_policy.test2.name
}
`
