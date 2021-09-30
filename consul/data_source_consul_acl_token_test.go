package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataACLToken_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataACLTokenConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.consul_acl_token.read", "accessor_id"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "description", "test"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "expiration_time", ""),
					resource.TestCheckResourceAttrSet("data.consul_acl_token.read", "id"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "local", "false"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "node_identities.#", "1"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "node_identities.0.datacenter", "bar"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "node_identities.0.node_name", "foo"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "policies.#", "1"),
					resource.TestCheckResourceAttrSet("data.consul_acl_token.read", "policies.0.id"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "policies.0.name", "test"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "roles.#", "0"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "service_identities.#", "1"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "service_identities.0.datacenters.#", "1"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "service_identities.0.datacenters.0", "world"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "service_identities.0.service_name", "hello"),
				),
			},
		},
	})
}

func TestAccDataACLToken_namespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataACLTokenConfigNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccDataACLToken_namespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataACLTokenConfigNamespaceEE,
			},
		},
	})
}

const testAccDataACLTokenConfig = `
resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = ["${consul_acl_policy.test.name}"]
	local = false

	service_identities {
		service_name = "hello"
		datacenters = ["world"]
	}

	node_identities {
		node_name = "foo"
		datacenter = "bar"
	}
}

data "consul_acl_token" "read" {
    accessor_id = "${consul_acl_token.test.id}"
}
`

const testAccDataACLTokenConfigNamespaceCE = `
data "consul_acl_token" "read" {
  accessor_id = "foo"
  namespace   = "test-data-token"
}
`

const testAccDataACLTokenConfigNamespaceEE = `
resource "consul_namespace" "test" {
  name = "test-data-token"
}

resource "consul_acl_policy" "test" {
  name        = "test"
  rules       = "node \"\" { policy = \"read\" }"
  datacenters = [ "dc1" ]
  namespace   = consul_namespace.test.name
}

resource "consul_acl_token" "test" {
  description = "test"
  policies    = ["${consul_acl_policy.test.name}"]
  local       = true
  namespace   = consul_namespace.test.name
}

data "consul_acl_token" "read" {
  accessor_id = consul_acl_token.test.id
  namespace   = consul_namespace.test.name
}
`
