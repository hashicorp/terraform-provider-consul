package consul

import (
	"fmt"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccConsulService_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "service_id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "address", "www.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "node", "compute-example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.0", "tag0"),
					testAccConsulExternalSource,
				),
			},
			resource.TestStep{
				PreConfig: func() {
					catalog := testAccProvider.Meta().(*consulapi.Client).Catalog()
					wOpts := &consulapi.WriteOptions{}
					dereg := &consulapi.CatalogDeregistration{
						Node:      "comput-example",
						ServiceID: "example",
					}
					_, err := catalog.Deregister(dereg, wOpts)
					if err != nil {
						t.Errorf("err: %v", err)
					}
				},
				Config: testAccConsulServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "service_id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "address", "www.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "node", "compute-example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.0", "tag0"),
					testAccConsulExternalSource,
				),
			},
		},
	})
}

func TestAccConsulService_basicModify(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "id", "example"),
				),
			},
			resource.TestStep{
				Config: testAccConsulServiceConfigBasicNewTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "tags.#", "2"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.0", "tag0"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.1", "tag1"),
				),
			},
			resource.TestStep{
				Config: testAccConsulServiceConfigBasicAddress,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "address", "lb.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.#", "2"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.0", "tag0"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.1", "tag1"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
				),
			},
		},
	})
}

func TestAccConsulService_serviceID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulServiceConfigServiceID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "id", "8ce84078-b32a-4039-bb68-17b13b7c2396"),
					resource.TestCheckResourceAttr("consul_service.example", "service_id", "8ce84078-b32a-4039-bb68-17b13b7c2396"),
					resource.TestCheckResourceAttr("consul_service.example", "address", "www.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "node", "compute-example"),
				),
			},
		},
	})
}

func TestAccConsulService_nodeDoesNotExist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccConsulServiceConfigNoNode,
				ExpectError: regexp.MustCompile(`Node does not exist: '*'`),
			},
		},
	})
}

func testAccConsulExternalSource(s *terraform.State) error {
	client := testAccProvider.Meta().(*consulapi.Client)
	qOpts := consulapi.QueryOptions{}

	service, _, err := client.Catalog().Service("example", "", &qOpts)
	if err != nil {
		return fmt.Errorf("Failed to retrieve service: %v", err)
	}

	for _, s := range service {
		source, ok := s.ServiceMeta["external-source"]
		if !ok || source != "terraform" {
			return fmt.Errorf("external-source not set")
		}
	}
	return nil
}

func testAccCheckConsulServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*consulapi.Client)

	qOpts := consulapi.QueryOptions{}
	services, _, err := client.Catalog().Services(&qOpts)
	if err != nil {
		return fmt.Errorf("Failed to retrieve services: %v", err)
	}

	if len(services) > 1 {
		return fmt.Errorf("Matching services still exsist: %v", services)
	}

	return nil
}

const testAccConsulServiceConfigNoNode = `
resource "consul_service" "example" {
	name = "example"
	node = "external"
}
`

const testAccConsulServiceConfigBasic = `
resource "consul_service" "example" {
	name    = "example"
	node    = "${consul_node.compute.name}"
	port    = 80
	tags    = ["tag0"]
  }

  resource "consul_node" "compute" {
	name    = "compute-example"
	address = "www.hashicorptest.com"
  }
`

const testAccConsulServiceConfigBasicNewTags = `
resource "consul_service" "example" {
	name    = "example"
	node    = "${consul_node.compute.name}"
	port    = 80
	tags    = ["tag0", "tag1"]
  }

  resource "consul_node" "compute" {
	name    = "compute-example"
	address = "www.hashicorptest.com"
  }
`

const testAccConsulServiceConfigBasicAddress = `
resource "consul_service" "example" {
	name    = "example"
	address = "lb.hashicorptest.com"
	node    = "${consul_node.compute.name}"
	port    = 80
	tags    = ["tag0", "tag1"]
  }

  resource "consul_node" "compute" {
	name    = "compute-example"
	address = "www.hashicorptest.com"
  }
`
const testAccConsulServiceConfigServiceID = `
resource "consul_service" "example" {
	name       = "example"
	service_id = "8ce84078-b32a-4039-bb68-17b13b7c2396"
	node       = "${consul_node.compute.name}"
  }

  resource "consul_node" "compute" {
	name    = "compute-example"
	address = "www.hashicorptest.com"
  }
`
