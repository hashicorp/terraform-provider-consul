// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckConsulACLTokenDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
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
}

func TestAccConsulACLToken_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulACLTokenDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenConfigBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("consul_acl_token.test", "accessor_id"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "description", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "expiration_time", ""),
					resource.TestCheckResourceAttrSet("consul_acl_token.test", "id"),
					resource.TestCheckResourceAttrSet("consul_acl_token.test", "local"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "node_identities.#", "0"),
					resource.TestCheckResourceAttrSet("consul_acl_token.test", "policies.#"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "service_identities.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.#", "0"),
				),
			},
			{
				Config: testResourceACLTokenConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("consul_acl_token.test", "accessor_id"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "description", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "expiration_time", ""),
					resource.TestCheckResourceAttrSet("consul_acl_token.test", "id"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "local", "false"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "node_identities.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "node_identities.0.datacenter", "bar"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "node_identities.0.node_name", "foo"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.111830242", "test2"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "service_identities.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "service_identities.0.datacenters.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "service_identities.0.datacenters.0", "world"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "service_identities.0.service_name", "hello"),
				),
			},
			{
				Config: testResourceACLTokenConfigUpdate,
			},
			{
				Config:   testResourceACLTokenConfigUpdateTemplatedPolicies,
				SkipFunc: skipIfConsulVersionLT(client, "1.17.0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.#", "2"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.0.datacenters.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.0.datacenters.0", "world"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.0.template_variables.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.0.template_variables.0.name", "web"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.0.template_name", "builtin/service"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.1.template_variables.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.1.template_name", "builtin/dns"),
				),
			},
			{
				Config: testResourceACLTokenConfigRole,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("consul_acl_token.test", "accessor_id"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "description", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "expiration_time", ""),
					resource.TestCheckResourceAttrSet("consul_acl_token.test", "id"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "local", "false"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "node_identities.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "policies.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "roles.1785148924", "test"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "service_identities.#", "0"),
					resource.TestCheckResourceAttr("consul_acl_token.test", "templated_policies.#", "0"),
				),
			},
		},
	})
}

func TestAccConsulACLToken_import(t *testing.T) {
	providers, _ := startTestServer(t)

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
		Providers: providers,
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
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testResourceACLTokenConfigNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccConsulACLToken_namespaceEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLTokenConfigNamespaceEE,
			},
		},
	})
}

const (
	testResourceACLTokenConfigBasic = `
resource "consul_acl_policy" "test" {
	name = "test-token-basic"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = [consul_acl_policy.test.name]
	local = true
}`

	testResourceACLTokenConfigUpdate = `
// Using another resource to force the update of consul_acl_token
resource "consul_acl_policy" "test2" {
	name = "test2"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = ["${consul_acl_policy.test2.name}"]

	service_identities {
		service_name = "hello"
		datacenters = ["world"]
	}

	node_identities {
		node_name = "foo"
		datacenter = "bar"
	}
}`

	testResourceACLTokenConfigUpdateTemplatedPolicies = `
// Using another resource to force the update of consul_acl_token
resource "consul_acl_policy" "test2" {
	name = "test2"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = ["${consul_acl_policy.test2.name}"]

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
}`

	testResourceACLTokenConfigRole = `
resource "consul_acl_role" "test" {
    name      = "test"
}

resource "consul_acl_token" "test" {
	description = "test"
	roles = [consul_acl_role.test.name]
}`

	testResourceACLTokenConfigNamespaceCE = `
resource "consul_acl_token" "test" {
  description = "test"
  namespace   = "test"
}`

	testResourceACLTokenConfigNamespaceEE = `
resource "consul_namespace" "test" {
  name = "test-token"
}
resource "consul_acl_token" "test" {
  description = "test"
  namespace   = consul_namespace.test.name
}`
)
