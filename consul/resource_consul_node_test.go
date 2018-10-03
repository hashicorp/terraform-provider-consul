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
			resource.TestStep{
				Config: testAccConsulNodeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulNodeExists(),
					testAccCheckConsulNodeValue("consul_node.foo", "address", "127.0.0.1"),
					testAccCheckConsulNodeValue("consul_node.foo", "name", "foo"),
				),
			},
			resource.TestStep{
				PreConfig: func() {
					catalog := testAccProvider.Meta().(*consulapi.Client).Catalog()
					wOpts := &consulapi.WriteOptions{}
					dereg := &consulapi.CatalogDeregistration{
						Node: "foo",
					}
					_, err := catalog.Deregister(dereg, wOpts)
					if err != nil {
						t.Errorf("err: %v", err)
					}
				},
				Config: testAccConsulNodeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulNodeExists(),
					testAccCheckConsulNodeValue("consul_node.foo", "address", "127.0.0.1"),
					testAccCheckConsulNodeValue("consul_node.foo", "name", "foo"),
				),
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
			resource.TestStep{
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
			resource.TestStep{
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
	catalog := testAccProvider.Meta().(*consulapi.Client).Catalog()
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
		catalog := testAccProvider.Meta().(*consulapi.Client).Catalog()
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

const testAccConsulNodeConfigBasic = `
resource "consul_node" "foo" {
	name 	= "foo"
	address = "127.0.0.1"
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
