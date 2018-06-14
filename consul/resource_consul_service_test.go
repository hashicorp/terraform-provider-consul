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
					resource.TestCheckResourceAttr("consul_service.google", "id", "google"),
					resource.TestCheckResourceAttr("consul_service.google", "address", "www.google.com"),
					resource.TestCheckResourceAttr("consul_service.google", "node", "compute-google"),
					resource.TestCheckResourceAttr("consul_service.google", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.google", "tags.#", "1"),
					resource.TestCheckResourceAttr("consul_service.google", "tags.0", "tag0"),
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
					resource.TestCheckResourceAttr("consul_service.google", "id", "google"),
				),
			},
			resource.TestStep{
				Config: testAccConsulServiceConfigBasicNewTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.google", "tags.#", "2"),
					resource.TestCheckResourceAttr("consul_service.google", "tags.0", "tag0"),
					resource.TestCheckResourceAttr("consul_service.google", "tags.1", "tag1"),
				),
			},
			resource.TestStep{
				Config: testAccConsulServiceConfigBasicAddress,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.google", "id", "google"),
					resource.TestCheckResourceAttr("consul_service.google", "address", "lb.google.com"),
					resource.TestCheckResourceAttr("consul_service.google", "tags.#", "2"),
					resource.TestCheckResourceAttr("consul_service.google", "tags.0", "tag0"),
					resource.TestCheckResourceAttr("consul_service.google", "tags.1", "tag1"),
					resource.TestCheckResourceAttr("consul_service.google", "port", "80"),
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
					resource.TestCheckResourceAttr("consul_service.google", "id", "8ce84078-b32a-4039-bb68-17b13b7c2396"),
					resource.TestCheckResourceAttr("consul_service.google", "service_id", "8ce84078-b32a-4039-bb68-17b13b7c2396"),
					resource.TestCheckResourceAttr("consul_service.google", "address", "www.google.com"),
					resource.TestCheckResourceAttr("consul_service.google", "node", "compute-google"),
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
	name = "google"
	node = "external"
}
`

const testAccConsulServiceConfigBasic = `
resource "consul_service" "google" {
	name    = "google"
	node    = "${consul_node.compute.name}"
	port    = 80
	tags    = ["tag0"]
  }

  resource "consul_node" "compute" {
	name    = "compute-google"
	address = "www.google.com"
  }
`

const testAccConsulServiceConfigBasicNewTags = `
resource "consul_service" "google" {
	name    = "google"
	node    = "${consul_node.compute.name}"
	port    = 80
	tags    = ["tag0", "tag1"]
  }

  resource "consul_node" "compute" {
	name    = "compute-google"
	address = "www.google.com"
  }
`

const testAccConsulServiceConfigBasicAddress = `
resource "consul_service" "google" {
	name    = "google"
	address = "lb.google.com"
	node    = "${consul_node.compute.name}"
	port    = 80
	tags    = ["tag0", "tag1"]
  }

  resource "consul_node" "compute" {
	name    = "compute-google"
	address = "www.google.com"
  }
`
const testAccConsulServiceConfigServiceID = `
resource "consul_service" "google" {
	name       = "google"
	service_id = "8ce84078-b32a-4039-bb68-17b13b7c2396"
	node       = "${consul_node.compute.name}"
  }

  resource "consul_node" "compute" {
	name    = "compute-google"
	address = "www.google.com"
  }
`
