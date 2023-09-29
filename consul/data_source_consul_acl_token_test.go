// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataACLToken_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
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
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "policies.0.name", "test-token"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "roles.#", "0"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "service_identities.#", "1"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "service_identities.0.datacenters.#", "1"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "service_identities.0.datacenters.0", "world"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "service_identities.0.service_name", "hello"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "templated_policies.#", "2"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "templated_policies.0.datacenters.#", "1"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "templated_policies.0.datacenters.0", "world"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "templated_policies.0.template_variables.#", "1"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "templated_policies.0.template_variables.0.name", "web"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "templated_policies.0.template_name", "builtin/service"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "templated_policies.1.template_variables.#", "0"),
					resource.TestCheckResourceAttr("data.consul_acl_token.read", "templated_policies.1.template_name", "builtin/dns"),
				),
			},
		},
	})
}

func TestAccDataACLToken_namespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
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
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
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
	name = "test-token"
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

	templated_policies {
		template_name = "builtin/service"
		datacenters = ["world"]
		template_variables {
			name = "web"
		}
	}

	templated_policies {
		template_name = "builtin/dns"
		datacenters = ["world"]
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
  name        = "test-token"
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
