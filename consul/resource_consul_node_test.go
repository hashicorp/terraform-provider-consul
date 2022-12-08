package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulNode_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulNodeDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNodeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulNodeExists(client),
					testAccCheckConsulNodeValue("consul_node.foo", "address", "127.0.0.1"),
					testAccCheckConsulNodeValue("consul_node.foo", "name", "foo"),
				),
			},
			{
				Config:        testAccConsulNodeConfigBasic,
				ResourceName:  "consul_node.foo",
				ImportState:   true,
				ImportStateId: "foo",
			},
			{
				PreConfig: testAccRemoveConsulNode(t, client),
				Config:    testAccConsulNodeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulNodeExists(client),
					testAccCheckConsulNodeValue("consul_node.foo", "address", "127.0.0.1"),
					testAccCheckConsulNodeValue("consul_node.foo", "name", "foo"),
				),
			},
			{
				// consul_node must detect changes made to its address...
				PreConfig: testAccChangeConsulNodeAddress(t, client),
				Config:    testAccConsulNodeConfigBasic,
				Check:     testAccConsulNodeDetectAttributeChanges(client),
			},
			{
				// ... and to its meta information
				PreConfig: testAccChangeConsulNodeAddressMeta(t, client),
				Config:    testAccConsulNodeConfigBasic,
				Check:     testAccConsulNodeDetectAttributeChanges(client),
			},
		},
	})
}

func TestAccConsulNode_nodeMeta(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    providers,
		CheckDestroy: testAccCheckConsulNodeDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNodeConfigNodeMeta,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulNodeExists(client),
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
					testAccCheckConsulNodeExists(client),
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

func TestAccConsulNode_datacenter(t *testing.T) {
	providers, client := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulNodeDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNodeConfigDatacenter,
				Check: func(s *terraform.State) error {
					test := func(dc string) error {
						c := client.Catalog()
						opts := &consulapi.QueryOptions{
							Datacenter: dc,
						}
						nodes, _, err := c.Nodes(opts)
						if err != nil {
							return err
						}

						for _, n := range nodes {
							if n.Node == dc {
								return nil
							}
						}
						return fmt.Errorf("could not find node %q", dc)
					}
					if err := test("dc1"); err != nil {
						return err
					}
					return test("dc2")
				},
			},
		},
	})
}

func testAccCheckConsulNodeDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
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
}

func testAccCheckConsulNodeExists(client *consulapi.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

func testAccRemoveConsulNode(t *testing.T, client *consulapi.Client) func() {
	return func() {
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

func testAccChangeConsulNodeAddress(t *testing.T, client *consulapi.Client) func() {
	return func() {
		catalog := client.Catalog()
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

func testAccChangeConsulNodeAddressMeta(t *testing.T, client *consulapi.Client) func() {
	return func() {
		catalog := client.Catalog()
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

func testAccConsulNodeDetectAttributeChanges(client *consulapi.Client) func(*terraform.State) error {
	return func(s *terraform.State) error {
		catalog := client.Catalog()
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

const testAccConsulNodeConfigDatacenter = `
resource "consul_node" "dc1" {
	name 	= "dc1"
	address = "127.0.0.1"
}

resource "consul_node" "dc2" {
	datacenter = "dc2"
	name 	   = "dc2"
	address    = "127.0.0.1"
}
`
