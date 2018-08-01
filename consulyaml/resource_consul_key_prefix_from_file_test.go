package consulyaml

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
			resource.TestStep{
				Config: testAccConsulKeyPrefixConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeyPrefixKeyValue("cheese", "chevre"),
					testAccCheckConsulKeyPrefixKeyValue("bread", "baguette"),
					testAccCheckConsulKeyPrefixKeyAbsent("species"),
					testAccCheckConsulKeyPrefixKeyAbsent("meat"),
				),
			},
		},
	})
}

func testAccCheckConsulKeyPrefixDestroy(s *terraform.State) error {
	kv := testAccProvider.Meta().(*consulapi.Client).KV()
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
		kv := testAccProvider.Meta().(*consulapi.Client).KV()
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
		kv := testAccProvider.Meta().(*consulapi.Client).KV()
		opts := &consulapi.WriteOptions{Datacenter: "dc1"}
		pair := &consulapi.KVPair{
			Key:   fullName,
			Value: []byte(value),
		}
		_, err := kv.Put(pair, opts)
		return err
	}
}

func testAccCheckConsulKeyPrefixKeyValue(name, value string) resource.TestCheckFunc {
	fullName := "prefix_test/" + name
	return func(s *terraform.State) error {
		kv := testAccProvider.Meta().(*consulapi.Client).KV()
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
		return nil
	}
}

const testAccConsulKeyPrefixConfig = `
resource "consul-yaml" "app" {
	datacenter = "dc1"

	path_prefix = "prefix_test/"
	subkeys_file = "test-fixtures/cheese.yaml"
}
`
