// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulKeys_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulKeysDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeysConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysExists(client),
					testAccCheckConsulKeysValue("consul_keys.app", "enabled", "true"),
					testAccCheckConsulKeysValue("consul_keys.app", "set", "acceptance"),
					testAccCheckConsulKeysValue("consul_keys.app", "remove_one", "hello"),
					resource.TestCheckResourceAttr("consul_keys.app", "key.4258512057.flags", "0"),
				),
			},
			{
				Config: testAccConsulKeysConfig_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysExists(client),
					testAccCheckConsulKeysValue("consul_keys.app", "enabled", "true"),
					testAccCheckConsulKeysValue("consul_keys.app", "set", "acceptanceUpdated"),
					testAccCheckConsulKeysRemoved("consul_keys.app", "remove_one"),
				),
			},
		},
	})
}

func TestAccConsulKeys_EmptyValue(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeysEmptyValue,
				Check:  testAccCheckConsulKeysExists(client),
			},
		},
	})
}

func TestAccConsulKeys_NamespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulKeysNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccConsulKeys_NamespaceEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeysNamespaceEE,
			},
		},
	})
}

func TestAccConsulKeys_Datacenter(t *testing.T) {
	providers, client := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeysDatacenter,
				Check:  testAccCheckConsulKeysDatacenter(client),
			},
		},
	})
}

func testAccCheckConsulKeysDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		kv := client.KV()
		opts := &consulapi.QueryOptions{Datacenter: "dc1"}
		pair, _, err := kv.Get("test/set", opts)
		if err != nil {
			return err
		}
		if pair != nil {
			return fmt.Errorf("Key still exists: %#v", pair)
		}
		return nil
	}
}

func testAccCheckConsulKeysExists(client *consulapi.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		kv := client.KV()
		opts := &consulapi.QueryOptions{Datacenter: "dc1"}
		pair, _, err := kv.Get("test/set", opts)
		if err != nil {
			return err
		}
		if pair == nil {
			return fmt.Errorf("Key 'test/set' does not exist")
		}
		return nil
	}
}

func testAccCheckConsulKeysValue(n, attr, val string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		out, ok := rn.Primary.Attributes["var."+attr]
		if !ok {
			return fmt.Errorf("Attribute '%s' not found: %#v", attr, rn.Primary.Attributes)
		}
		if val != "<any>" && out != val {
			return fmt.Errorf("Attribute '%s' value '%s' != '%s'", attr, out, val)
		}
		if val == "<any>" && out == "" {
			return fmt.Errorf("Attribute '%s' value '%s'", attr, out)
		}
		return nil
	}
}

func testAccCheckConsulKeysRemoved(n, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		_, ok = rn.Primary.Attributes["var."+attr]
		if ok {
			return fmt.Errorf("Attribute '%s' still present: %#v", attr, rn.Primary.Attributes)
		}
		return nil
	}
}

func testAccCheckConsulKeysDatacenter(client *consulapi.Client) func(s *terraform.State) error {
	test := func(dc string) error {
		kv := client.KV()
		opts := &consulapi.QueryOptions{
			Datacenter: dc,
		}
		pair, _, err := kv.Get("foo/dc", opts)
		if err != nil {
			return err
		}
		if kv == nil {
			return fmt.Errorf("key 'dc' does not exist")
		}
		value := string(pair.Value)
		if value != dc {
			return fmt.Errorf("wrong value: %q", value)
		}
		return nil
	}

	return func(s *terraform.State) error {
		if err := test("dc1"); err != nil {
			return err
		}
		return test("dc2")
	}
}

const testAccConsulKeysConfig = `
resource "consul_keys" "app" {
	datacenter = "dc1"
	key {
		name = "enabled"
		path = "test/enabled"
		default = "true"
	}
	key {
		name = "set"
		path = "test/set"
		value = "acceptance"
		delete = true
	}
	key {
		name = "remove_one"
		path = "test/remove_one"
		value = "hello"
		delete = true
	}
}
`

const testAccConsulKeysConfig_Update = `
resource "consul_keys" "app" {
	datacenter = "dc1"
	key {
		name = "enabled"
		path = "test/enabled"
		default = "true"
	}
	key {
		name = "set"
		path = "test/set"
		value = "acceptanceUpdated"
		flags = 64
		delete = true
	}
}
`

const testAccConsulKeysEmptyValue = `
resource "consul_keys" "consul" {
	key {
	  path  = "test/set"
	  value = ""
	  delete = true
	}
}`

const testAccConsulKeysNamespaceCE = `
resource "consul_keys" "consul" {
  namespace = "test-keys"

  key {
    path  = "test/set"
    value = ""
    delete = true
  }
}`

const testAccConsulKeysNamespaceEE = `
resource "consul_namespace" "test" {
  name = "test-keys"
}

resource "consul_keys" "consul" {
  namespace = consul_namespace.test.name

  key {
    path   = "test/set"
    value  = ""
    delete = true
  }
}`

const testAccConsulKeysDatacenter = `
resource "consul_keys" "dc1" {
	key {
		path   = "foo/dc"
		value  = "dc1"
		delete = true
	}
}

resource "consul_keys" "dc2" {
	datacenter = "dc2"

	key {
		path   = "foo/dc"
		value  = "dc2"
		delete = true
	}
}
`
