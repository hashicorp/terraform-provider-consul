package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataACLRole_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceACLRoleConfigNotFound,
				ExpectError: regexp.MustCompile("could not find role 'not-found'"),
			},
			{
				Config: testAccDataSourceACLRoleConfigBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "description", "bar"),
					resource.TestCheckResourceAttrSet("data.consul_acl_role.test", "id"),
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "name", "foo"),
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "node_identities.#", "1"),
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "node_identities.0.datacenter", "world"),
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "node_identities.0.node_name", "hello"),
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "policies.#", "1"),
					resource.TestCheckResourceAttrSet("data.consul_acl_role.test", "policies.0.id"),
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "policies.0.name", "test-role"),
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "service_identities.#", "1"),
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "service_identities.0.datacenters.#", "0"),
					resource.TestCheckResourceAttr("data.consul_acl_role.test", "service_identities.0.service_name", "foo"),
				),
			},
		},
	})
}

func TestAccDataACLRole_namespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceACLRoleConfigNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccDataACLRole_namespaceEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceACLRoleConfigNamespaceEE,
			},
		},
	})
}

const testAccDataSourceACLRoleConfigNotFound = `
data "consul_acl_role" "test" {
	name = "not-found"
}
`

const testAccDataSourceACLRoleConfigBasic = `
resource "consul_acl_policy" "test-read" {
	name = "test-role"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role" "test" {
	name      = "foo"
	description = "bar"

	policies = [
		consul_acl_policy.test-read.id
	]

	service_identities {
		service_name = "foo"
	}

	node_identities {
		node_name = "hello"
		datacenter = "world"
	}
}

data "consul_acl_role" "test" {
	name = consul_acl_role.test.name
}
`
const testAccDataSourceACLRoleConfigNamespaceCE = `
data "consul_acl_role" "test" {
  name      = "test"
  namespace = "test-data-role"
}
`

const testAccDataSourceACLRoleConfigNamespaceEE = `
resource "consul_namespace" "test" {
  name = "test-data-role"
}

resource "consul_acl_policy" "test-read" {
  name        = "test"
  rules       = "node \"\" { policy = \"read\" }"
  namespace   = consul_namespace.test.name
}

resource "consul_acl_role" "test" {
  name        = "foo"
  description = "bar"
  namespace   = consul_namespace.test.name

  policies = [
    consul_acl_policy.test-read.id
  ]

  service_identities {
    service_name = "foo"
  }
}

data "consul_acl_role" "test" {
  name      = consul_acl_role.test.name
  namespace = consul_namespace.test.name
}
`
