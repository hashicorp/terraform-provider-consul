package consul

import (
	"fmt"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulKeyPrefix_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckConsulKeyPrefixKeyAbsent("species"),
			testAccCheckConsulKeyPrefixKeyAbsent("meat"),
			testAccCheckConsulKeyPrefixKeyAbsent("cheese"),
			testAccCheckConsulKeyPrefixKeyAbsent("bread"),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeyPrefixNoKeys,
			},
			{
				Config: testAccConsulKeyPrefixConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeyPrefixKeyValue("cheese", "chevre", 0),
					testAccCheckConsulKeyPrefixKeyValue("bread", "baguette", 0),
					testAccCheckConsulKeyPrefixKeyValue("condiment/first", "tomato", 2),
					testAccCheckConsulKeyPrefixKeyValue("condiment/second", "salad", 4),
					testAccCheckConsulKeyPrefixKeyAbsent("species"),
					testAccCheckConsulKeyPrefixKeyAbsent("meat"),
				),
			},
			{
				Config:             testAccConsulKeyPrefixConfig,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					// This will add a rogue key that Terraform isn't
					// expecting, causing a non-empty plan that wants
					// to remove it.
					testAccAddConsulKeyPrefixRogue("species", "gorilla"),
				),
			},
			{
				Config: testAccConsulKeyPrefixConfig_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeyPrefixKeyValue("meat", "ham", 0),
					testAccCheckConsulKeyPrefixKeyValue("bread", "batard", 0),
					testAccCheckConsulKeyPrefixKeyValue("condiment/second", "mayonnaise", 4),
					testAccCheckConsulKeyPrefixKeyValue("condiment/third", "onion", 0),
					testAccCheckConsulKeyPrefixKeyAbsent("condiment/first"),
					testAccCheckConsulKeyPrefixKeyAbsent("cheese"),
					testAccCheckConsulKeyPrefixKeyAbsent("species"),
				),
			},
			{
				Config:             testAccConsulKeyPrefixConfig_Update,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccAddConsulKeyPrefixRogue("species", "gorilla"),
				),
			},
			{
				PreConfig: func() {
					kv := getTestClient(testAccProvider.Meta()).KV()
					kv.DeleteTree("", nil)
				},
				Config: testAccConsulKeyPrefixConfig_root,
			},
		},
	})
}

func TestAccCheckConsulKeyPrefix_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeyPrefixConfig_Import,
			},
			{
				Config:            testAccConsulKeyPrefixConfig_Import,
				ResourceName:      "consul_key_prefix.app",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsulKeyPrefix_namespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulKeyPrefixConfig_namespaceCE,
				ExpectError: regexp.MustCompile("Unexpected response code: 400"),
			},
		},
	})
}

func TestAccConsulKeyPrefix_namespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeyPrefixConfig_namespaceEE,
			},
		},
	})
}

// TestAccConsulKeyPrefix_deleted checks that resource will recreate keys
// the consul_key_prefix resource if all the keys has been deleted on Consul
func TestAccConsulKeyPrefix_deleted(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Apply the config and remove the prefix in Consul
				Config: testAccConsulKeyPrefixConfig_deleted,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeyPrefixKeyValue("first", "plop", 0),
					testAccCheckConsulKeyPrefixKeyValue("second", "plip", 0),
				),
			},
			{
				// This will remove all the key_prefix in Consul
				// causing a non-empty plan that wants to recreate it.
				PreConfig: testAccDeleteConsulKeyPrefix(t, "prefix_test/"),
				// This step should recreate the missing keys
				Config: testAccConsulKeyPrefixConfig_deleted,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeyPrefixKeyValue("first", "plop", 0),
					testAccCheckConsulKeyPrefixKeyValue("second", "plip", 0),
				),
			},
			{
				// Apply again and remove one key under the prefix
				Config: testAccConsulKeyPrefixConfig_deleted,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeyPrefixKeyValue("first", "plop", 0),
					testAccCheckConsulKeyPrefixKeyValue("second", "plip", 0),
				),
			},
			{
				// Remove the first key, this should cause a non-empty plan
				// to recreate it
				PreConfig: testAccDeleteConsulKey(t, "prefix_test/first"),
				// This step should recreate the missing key
				Config: testAccConsulKeyPrefixConfig_deleted,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeyPrefixKeyValue("first", "plop", 0),
					testAccCheckConsulKeyPrefixKeyValue("second", "plip", 0),
				),
			},
		},
	})
}

func TestAccConsulKeyPrefix_datacenter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccRemoteDatacenterPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConsulKeyPrefixConfig_datacenter,
				Check:  testAccCheckConsulKeysDatacenter,
			},
		},
	})
}

func testAccCheckConsulKeyPrefixDestroy(s *terraform.State) error {
	client := getTestClient(testAccProvider.Meta())
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

func testAccCheckConsulKeyPrefixKeyAbsent(name string) resource.TestCheckFunc {
	fullName := "prefix_test/" + name
	return func(s *terraform.State) error {
		client := getTestClient(testAccProvider.Meta())
		kv := client.KV()
		opts := &consulapi.QueryOptions{Datacenter: "dc1"}
		pair, _, err := kv.Get(fullName, opts)
		if err != nil {
			return err
		}
		if pair != nil {
			return fmt.Errorf("key '%s' exists, but shouldn't", fullName)
		}
		return nil
	}
}

// This one is actually not a check, but rather a mutation step. It writes
// a value directly into Consul, bypassing our Terraform resource.
func testAccAddConsulKeyPrefixRogue(name, value string) resource.TestCheckFunc {
	fullName := "prefix_test/" + name
	return func(s *terraform.State) error {
		client := getTestClient(testAccProvider.Meta())
		kv := client.KV()
		opts := &consulapi.WriteOptions{Datacenter: "dc1"}
		pair := &consulapi.KVPair{
			Key:   fullName,
			Value: []byte(value),
		}
		_, err := kv.Put(pair, opts)
		return err
	}
}

// This one is actually not a check, but rather a mutation step.
// It removes the prefix_test "folder" (all keys under this prefix)
func testAccDeleteConsulKeyPrefix(t *testing.T, prefix string) func() {
	return func() {
		kv := getTestClient(testAccProvider.Meta()).KV()
		_, err := kv.DeleteTree(prefix, &consulapi.WriteOptions{Datacenter: "dc1"})
		if err != nil {
			t.Fatalf("failed to delete tree: %v", err)
		}
	}
}

// This one is actually not a check, but rather a mutation step.
// It removes one key in Consul
func testAccDeleteConsulKey(t *testing.T, key string) func() {
	return func() {
		kv := getTestClient(testAccProvider.Meta()).KV()
		_, err := kv.Delete(key, &consulapi.WriteOptions{Datacenter: "dc1"})
		if err != nil {
			t.Fatalf("failed to delete key %q: %v", key, err)
		}
	}
}

func testAccCheckConsulKeyPrefixKeyValue(name, value string, flags uint64) resource.TestCheckFunc {
	fullName := "prefix_test/" + name
	return func(s *terraform.State) error {
		client := getTestClient(testAccProvider.Meta())
		kv := client.KV()
		opts := &consulapi.QueryOptions{Datacenter: "dc1"}
		pair, _, err := kv.Get(fullName, opts)
		if err != nil {
			return err
		}
		if pair == nil {
			return fmt.Errorf("key %v doesn't exist, but should", fullName)
		}
		if string(pair.Value) != value {
			return fmt.Errorf("key %v has value %v; want %v", fullName, pair.Value, value)
		}
		if pair.Flags != flags {
			return fmt.Errorf("key %v has flags %v; want %v", fullName, pair.Flags, flags)
		}
		return nil
	}
}

const testAccConsulKeyPrefixNoKeys = `
resource "consul_key_prefix" "app" {
	datacenter = "dc1"
	path_prefix = "prefix_test/"
}`

const testAccConsulKeyPrefixConfig = `
resource "consul_key_prefix" "app" {
	datacenter = "dc1"

    path_prefix = "prefix_test/"

    subkeys = {
        cheese = "chevre"
        bread = "baguette"
	}

	subkey {
		path  = "condiment/first"
		value = "tomato"
		flags = 2
	}

	subkey {
		path  = "condiment/second"
		value = "salad"
		flags = 4
	}
}
`

const testAccConsulKeyPrefixConfig_Update = `
resource "consul_key_prefix" "app" {
	datacenter = "dc1"

    path_prefix = "prefix_test/"

    subkeys = {
        bread = "batard"
        meat = "ham"
    }

	subkey {
		path  = "condiment/second"
		value = "mayonnaise"
		flags = 4
	}

	subkey {
		path  = "condiment/third"
		value = "onion"
	}
}
`

const testAccConsulKeyPrefixConfig_Import = `
resource "consul_key_prefix" "app" {
	datacenter = "dc1"

    path_prefix = "prefix_test/"

    subkeys = {
        bread = "batard"
        meat = "ham"
    }
}
`

const testAccConsulKeyPrefixConfig_namespaceCE = `
resource "consul_key_prefix" "test" {
  path_prefix = "prefix_test/"
  namespace   = "test-key-prefix"

  subkeys = {
    bread = "batard"
    meat = "ham"
  }
}`

const testAccConsulKeyPrefixConfig_namespaceEE = `
resource "consul_namespace" "test" {
  name = "test-key-prefix"
}

resource "consul_key_prefix" "test" {
  path_prefix = "prefix_test/"
  namespace   = consul_namespace.test.name

  subkeys = {
    bread = "batard"
    meat = "ham"
  }
}`

const testAccConsulKeyPrefixConfig_deleted = `
resource "consul_key_prefix" "app" {
	datacenter = "dc1"
    path_prefix = "prefix_test/"

	subkey {
		path  = "first"
		value = "plop"
	}

	subkey {
		path  = "second"
		value = "plip"
	}
}
`

const testAccConsulKeyPrefixConfig_datacenter = `
resource "consul_key_prefix" "dc1" {
    path_prefix = "foo/"

	subkey {
		path  = "dc"
		value = "dc1"
	}
}

resource "consul_key_prefix" "dc2" {
	datacenter = "dc2"
    path_prefix = "foo/"

	subkey {
		path  = "dc"
		value = "dc2"
	}
}
`

const testAccConsulKeyPrefixConfig_root = `
resource "consul_key_prefix" "root" {
    path_prefix = ""

	subkey {
		path  = "foo"
		value = "bar"
	}
}
`
