package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var namespaceEnterpriseFeature = regexp.MustCompile("(?i)Namespaces .* Consul Enterprise feature")

func TestAccConsulNamespace_FailOnCommunityEdition(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulNamespace,
				ExpectError: regexp.MustCompile("failed to create namespace: Unexpected response code: 404"),
			},
		},
	})
}

func TestAccConsulNamespace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNamespace,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_namespace.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_namespace.test", "description", "test namespace"),
					resource.TestCheckResourceAttr("consul_namespace.test", "meta.%", "1"),
					resource.TestCheckResourceAttr("consul_namespace.test", "policy_defaults.#", "0"),
					resource.TestCheckResourceAttr("consul_namespace.test", "policy_defaults.#", "0"),
				),
			},
			{
				Config: testAccConsulNamespace_Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_namespace.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_namespace.test", "description", "updated description"),
					resource.TestCheckResourceAttr("consul_namespace.test", "meta.%", "0"),
					resource.TestCheckResourceAttr("consul_namespace.test", "role_defaults.#", "1"),
					resource.TestCheckResourceAttr("consul_namespace.test", "role_defaults.0", "foo"),
					resource.TestCheckResourceAttr("consul_namespace.test", "policy_defaults.#", "1"),
					resource.TestCheckResourceAttr("consul_namespace.test", "policy_defaults.0", "bar"),
				),
			},
		},
	})
}

const testAccConsulNamespace = `
resource "consul_namespace" "test" {
	name        = "test"
	description = "test namespace"

	meta = {
		foo = "bar"
	}
}
`

const testAccConsulNamespace_Update = `
resource "consul_acl_role" "test" {
  name = "foo"
}

resource "consul_acl_policy" "test" {
  name  = "bar"
  rules = "node_prefix \"\" { policy = \"read\" }"
}

resource "consul_namespace" "test" {
  name        = "test"
  description = "updated description"

  policy_defaults = [
    consul_acl_policy.test.name
  ]

  role_defaults = [
    consul_acl_role.test.name
  ]
}`
