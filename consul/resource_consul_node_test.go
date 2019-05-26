package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccConsulNode_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNodeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulNodeExists(),
					testAccCheckConsulNodeValue("consul_node.foo", "address", "127.0.0.1"),
					testAccCheckConsulNodeValue("consul_node.foo", "name", "foo"),
				),
			},
			{
				PreConfig: testAccRemoveConsulNode(t),
				Config:    testAccConsulNodeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulNodeExists(),
					testAccCheckConsulNodeValue("consul_node.foo", "address", "127.0.0.1"),
					testAccCheckConsulNodeValue("consul_node.foo", "name", "foo"),
				),
			},
			{
				// consul_node must detect changes made to its address...
				PreConfig: testAccChangeConsulNodeAddress(t),
				Config:    testAccConsulNodeConfigBasic,
				Check:     testAccConsulNodeDetectAttributeChanges,
			},
			{
				// ... and to its meta information
				PreConfig: testAccChangeConsulNodeAddressMeta(t),
				Config:    testAccConsulNodeConfigBasic,
				Check:     testAccConsulNodeDetectAttributeChanges,
			},
		},
	})
}

func TestAccConsulNode_nodeMeta(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNodeConfigNodeMeta,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulNodeExists(),
					testAccCheckConsulNodeValue("consul_node.foo", "address", "127.0.0.1"),
					testAccCheckConsulNodeValue("consul_node.foo", "name", "foo"),
					testAccCheckConsulNodeValue("consul_node.foo", "meta.%", "3"),
					testAccCheckConsulNodeValue("consul_node.foo", "meta.foo", "bar"),
					testAccCheckConsulNodeValue("consul_node.foo", "meta.update", "this"),
					testAccCheckConsulNodeValue("consul_node.foo", "meta.remove", "this"),
				),
			},
			{
				Config: testAccConsulNodeConfigNodeMeta_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulNodeExists(),
					testAccCheckConsulNodeValue("consul_node.foo", "address", "127.0.0.1"),
					testAccCheckConsulNodeValue("consul_node.foo", "name", "foo"),
					testAccCheckConsulNodeValue("consul_node.foo", "meta.%", "2"),
					testAccCheckConsulNodeValue("consul_node.foo", "meta.foo", "bar"),
					testAccCheckConsulNodeValue("consul_node.foo", "meta.update", "yes"),
					testAccCheckConsulNodeValueRemoved("consul_node.foo", "meta.remove"),
				),
			},
		},
	})
}

func testAccCheckConsulNodeDestroy(s *terraform.State) error {
	client := getClient(testAccProvider.Meta())
	catalog := client.Catalog()
	qOpts := consulapi.QueryOptions{}
	nodes, _, err := catalog.Nodes(&qOpts)
	if err != nil {
		return fmt.Errorf("Could not retrieve services: %#v", err)
	}
	for i := range nodes {
		if nodes[i].Node == "foo" {
			return fmt.Errorf("Node still exists: %#v", "foo")
		}
	}
	return nil
}

func testAccCheckConsulNodeExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := getClient(testAccProvider.Meta())
		catalog := client.Catalog()
		qOpts := consulapi.QueryOptions{}
		nodes, _, err := catalog.Nodes(&qOpts)
		if err != nil {
			return err
		}
		for i := range nodes {
			if nodes[i].Node == "foo" {
				return nil
			}
		}
		return fmt.Errorf("Service does not exist: %#v", "google")
	}
}

func testAccCheckConsulNodeValue(n, attr, val string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		out, ok := rn.Primary.Attributes[attr]
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

func testAccCheckConsulNodeValueRemoved(n, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		_, ok = rn.Primary.Attributes[attr]
		if ok {
			return fmt.Errorf("Attribute '%s' still present: %#v", attr, rn.Primary.Attributes)
		}
		return nil
	}
}

func testAccRemoveConsulNode(t *testing.T) func() {
	return func() {
		client := getClient(testAccProvider.Meta())
		catalog := client.Catalog()
		wOpts := &consulapi.WriteOptions{}
		dereg := &consulapi.CatalogDeregistration{
			Node: "foo",
		}
		_, err := catalog.Deregister(dereg, wOpts)
		if err != nil {
			t.Errorf("err: %v", err)
		}
	}
}

func testAccChangeConsulNodeAddress(t *testing.T) func() {
	return func() {
		catalog := testAccProvider.Meta().(*consulapi.Client).Catalog()
		wOpts := &consulapi.WriteOptions{}

		registration := &consulapi.CatalogRegistration{
			Address:    "wrong_address",
			Datacenter: "dc1",
			Node:       "foo",
			NodeMeta: map[string]string{
				"foo": "bar",
			},
		}
		_, err := catalog.Register(registration, wOpts)
		if err != nil {
			t.Errorf("err: %v", err)
		}
	}
}
func testAccChangeConsulNodeAddressMeta(t *testing.T) func() {
	return func() {
		catalog := testAccProvider.Meta().(*consulapi.Client).Catalog()
		wOpts := &consulapi.WriteOptions{}

		registration := &consulapi.CatalogRegistration{
			Address:    "127.0.0.1",
			Datacenter: "dc1",
			Node:       "foo",
		}
		_, err := catalog.Register(registration, wOpts)
		if err != nil {
			t.Errorf("err: %v", err)
		}
	}
}

func testAccConsulNodeDetectAttributeChanges(*terraform.State) error {
	catalog := testAccProvider.Meta().(*consulapi.Client).Catalog()
	n, _, err := catalog.Node("foo", &consulapi.QueryOptions{})
	if err != nil {
		return fmt.Errorf("Failed to read 'foo': %v", err)
	}
	if n == nil {
		return fmt.Errorf("No 'foo' node found")
	}
	if n.Node.Address != "127.0.0.1" {
		return fmt.Errorf("Wrong address: %s", n.Node.Address)
	}
	if len(n.Node.Meta) != 1 || n.Node.Meta["foo"] != "bar" {
		return fmt.Errorf("Wrong node meta: %v", n.Node.Meta)
	}
	return nil
}

const testAccConsulNodeConfigBasic = `
resource "consul_node" "foo" {
	name 	= "foo"
	address = "127.0.0.1"

	meta = {
		foo     = "bar"
	}
}
`

const testAccConsulNodeConfigNodeMeta = `
resource "consul_node" "foo" {
	name 	= "foo"
	address = "127.0.0.1"

	meta = {
		foo    = "bar"
		update = "this"
		remove = "this"
	}
}
`

const testAccConsulNodeConfigNodeMeta_Update = `
resource "consul_node" "foo" {
	name 	= "foo"
	address = "127.0.0.1"

	meta = {
		foo     = "bar"
		update  = "yes"
	}
}
`
