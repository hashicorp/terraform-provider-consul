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
			{
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
			{
				PreConfig: testAccRemoveConsulService(t),
				Config:    testAccConsulServiceConfigBasic,
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
			{
				Config: testAccConsulServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "id", "example"),
				),
			},
			{
				Config: testAccConsulServiceConfigBasicNewTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "tags.#", "2"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.0", "tag0"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.1", "tag1"),
				),
			},
			{
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
			{
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

func TestAccConsulServiceCheck(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulServiceCheck,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "name", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "check.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.check_id", "service:redis1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.name", "Redis health check"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.notes", "Script based health check"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.status", "passing"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.http", "https://www.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.interval", "5s"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.timeout", "1s"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.deregister_critical_service_after", "30s"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.header.#", "2"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.header.344754333.name", "bar"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.header.344754333.value.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.header.344754333.value.0", "test"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.header.2976766922.name", "foo"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.header.2976766922.value.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.0.header.2976766922.value.0", "test"),
					resource.TestCheckResourceAttr("consul_service.no-deregister", "check.0.deregister_critical_service_after", "30s"),
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
			{
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
		return fmt.Errorf("Matching services still exist: %v", services)
	}

	return nil
}

func testAccRemoveConsulService(t *testing.T) func() {
	return func() {
		catalog := testAccProvider.Meta().(*consulapi.Client).Catalog()
		wOpts := &consulapi.WriteOptions{}
		dereg := &consulapi.CatalogDeregistration{
			Node:      "compute-example",
			ServiceID: "example",
		}
		_, err := catalog.Deregister(dereg, wOpts)
		if err != nil {
			t.Errorf("err: %v", err)
		}
	}
}

const testAccConsulServiceConfigNoNode = `
resource "consul_service" "example" {
	name = "example"
	node = "external"
}
`

const testAccConsulServiceCheck = `
resource "consul_node" "compute" {
	name    = "compute-example"
	address = "www.hashicorptest.com"
}

resource "consul_service" "example" {
	name = "example"
	node = "${consul_node.compute.name}"
	port = 80

	check {
		check_id = "service:redis1"
		name = "Redis health check"
		notes = "Script based health check"
		status = "passing"
		http = "https://www.hashicorptest.com"
		tls_skip_verify = false
		method = "PUT"
		interval = "5s"
		timeout = "1s"
		deregister_critical_service_after = "30s"

		header {
		  name = "foo"
		  value = ["test"]
		}

		header {
		  name = "bar"
		  value = ["test"]
		}
	}
}

resource "consul_node" "external" {
	name    = "external-example"
	address = "www.hashicorptest.com"
}

resource "consul_service" "external" {
	name     = "example-external"
	node     = "${consul_node.external.name}"
	external = true
	port     = 80

	check {
		check_id = "service:redis1"
		name = "Redis health check"
		notes = "Script based health check"
		http = "https://www.google.com"
		interval = "5s"
		timeout = "1s"
		deregister_critical_service_after = "30s"
	}
}

resource "consul_service" "no-deregister" {
	name     = "example-external"
	node     = "${consul_node.external.name}"
	external = true
	port     = 80

	check {
		check_id = "service:redis1"
		name = "Redis health check"
		notes = "Script based health check"
		http = "https://www.google.com"
		interval = "5s"
		timeout = "1s"
	}
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
