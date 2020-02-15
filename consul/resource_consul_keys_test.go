package consul

import (
	"fmt"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulKeys_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulKeysDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeysConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysExists(),
					testAccCheckConsulKeysValue("consul_keys.app", "enabled", "true"),
					testAccCheckConsulKeysValue("consul_keys.app", "set", "acceptance"),
					testAccCheckConsulKeysValue("consul_keys.app", "remove_one", "hello"),
					resource.TestCheckResourceAttr("consul_keys.app", "key.4258512057.flags", "0"),
				),
			},
			{
				Config: testAccConsulKeysConfig_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysExists(),
					testAccCheckConsulKeysValue("consul_keys.app", "enabled", "true"),
					testAccCheckConsulKeysValue("consul_keys.app", "set", "acceptanceUpdated"),
					testAccCheckConsulKeysRemoved("consul_keys.app", "remove_one"),
				),
			},
		},
	})
}

func TestAccConsulKeys_EmptyValue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeysEmptyValue,
				Check:  testAccCheckConsulKeysExists(),
			},
		},
	})
}

func TestAccConsulKeys_NamespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulKeysNamespaceCE,
				ExpectError: regexp.MustCompile("Namespaces is a Consul Enterprise feature"),
			},
		},
	})
}

func TestAccConsulKeys_NamespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeysNamespaceEE,
			},
		},
	})
}

func testAccCheckConsulKeysFlags(path string, flags int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		kv := testAccProvider.Meta().(*consulapi.Client).KV()
		opts := &consulapi.QueryOptions{Datacenter: "dc1"}
		pair, _, err := kv.Get(path, opts)
		if err != nil {
			return err
		}
		if pair == nil {
			return fmt.Errorf("Key '%v' does not exist", path)
		}
		if int(pair.Flags) != flags {
			return fmt.Errorf("Wrong flags for '%v': %v != %v", path, int(pair.Flags), flags)
		}
		return nil
	}
}

func testAccCheckConsulKeysDestroy(s *terraform.State) error {
	client := getClient(testAccProvider.Meta())
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

func testAccCheckConsulKeysExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := getClient(testAccProvider.Meta())
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
    path  = "test/set"
    value = ""
    delete = true
  }
}`
