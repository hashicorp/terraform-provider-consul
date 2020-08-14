package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataACLPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceACLPolicyConfigNotFound,
				ExpectError: regexp.MustCompile("Could not find policy 'not-found'"),
			},
			{
				Config: testAccDataSourceACLPolicyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_acl_policy.test", "name", "test"),
					testAccCheckDataSourceValue("data.consul_acl_policy.test", "description", "foo"),
					testAccCheckDataSourceValue("data.consul_acl_policy.test", "rules", "node_prefix \"\" { policy = \"read\" }"),
					testAccCheckDataSourceValue("data.consul_acl_policy.test", "datacenters.#", "1"),
					testAccCheckDataSourceValue("data.consul_acl_policy.test", "datacenters.0", "dc1"),
				),
			},
		},
	})
}

func TestAccDataACLPolicy_namespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceACLPolicyNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccDataACLPolicy_namespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceACLPolicyNamespaceEE,
			},
		},
	})
}

const testAccDataSourceACLPolicyConfigNotFound = `
data "consul_acl_policy" "test" {
	name = "not-found"
}
`

const testAccDataSourceACLPolicyConfigBasic = `
resource "consul_acl_policy" "test" {
	name = "test"
	description = "foo"
	rules = "node_prefix \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

data "consul_acl_policy" "test" {
	name = consul_acl_policy.test.name
}
`

const testAccDataSourceACLPolicyNamespaceCE = `
data "consul_acl_policy" "test" {
  name      = "test"
  namespace = "test-policy"
}
`

const testAccDataSourceACLPolicyNamespaceEE = `
resource "consul_acl_policy" "test" {
  name      = "test"
  rules     = "node_prefix \"\" { policy = \"read\" }"
  namespace = consul_namespace.test.name
}

resource "consul_namespace" "test" {
  name = "test-data-policy"
}

data "consul_acl_policy" "test" {
  name      = consul_acl_policy.test.name
  namespace = consul_namespace.test.name
}
`
