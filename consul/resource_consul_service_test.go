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
		PreCheck:     func() { testAccPreCheck(t) },
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
					resource.TestCheckResourceAttr("consul_service.example", "meta.%", "0"),
					testAccConsulExternalSource,
				),
			},
			{
				Config: testAccConsulServiceConfigBasicMeta,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "service_id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "address", "www.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "node", "compute-example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.0", "tag0"),
					resource.TestCheckResourceAttr("consul_service.example", "meta.%", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "meta.test", "test"),
				),
			},
			{
				PreConfig: testAccRemoveConsulService(t),
				Config:    testAccConsulServiceConfigBasicMeta,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "service_id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "address", "www.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "node", "compute-example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "tags.0", "tag0"),
					resource.TestCheckResourceAttr("consul_service.example", "meta.%", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "meta.test", "test"),
				),
			},
		},
	})
}

func TestAccConsulService_basicModify(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
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
					resource.TestCheckResourceAttr("consul_service.example", "meta.test", "test"),
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
		PreCheck:     func() { testAccPreCheck(t) },
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulServiceCheck,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "name", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "check.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.check_id", "service:redis1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.name", "Redis health check"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.notes", "Script based health check"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.status", "passing"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.http", "https://www.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.interval", "5s"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.timeout", "1s"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.deregister_critical_service_after", "30s"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.header.#", "2"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.header.344754333.name", "bar"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.header.344754333.value.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.header.344754333.value.0", "test"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.header.2976766922.name", "foo"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.header.2976766922.value.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1674222516.header.2976766922.value.0", "test"),
					resource.TestCheckResourceAttr("consul_service.no-deregister", "check.3059104818.deregister_critical_service_after", "30s"),
				),
			},
			resource.TestStep{
				Config: testAccConsulServiceCheckID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "name", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "check.#", "1"),
				),
			},
		},
	})
}

func TestAccConsulServiceCheckOrder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulServiceCheckOrder,
			},
		},
	})
}

// When the same service is defined on multiple nodes, the health-checks must
// be associated to the correct instance.
func TestAccDataConsulServiceSameServiceMultipleNodes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataConsulServiceSameServiceMultipleNodes,
			},
		},
	})
}

func TestAccConsulService_nodeDoesNotExist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
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

func TestAccConsulService_dontOverrideNodeMeta(t *testing.T) {
	// This would raise an error if consul_service changed attributes of consul_node
	// since the next plan would not be empty
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulServiceDontOverrideNodeMeta,
			},
		},
	})
}

func testAccConsulExternalSource(s *terraform.State) error {
	client := getClient(testAccProvider.Meta())
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
	client := getClient(testAccProvider.Meta())
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
		client := getClient(testAccProvider.Meta())
		catalog := client.Catalog()
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

const testAccConsulServiceCheckOrder = `
resource "consul_node" "external" {
	name    = "external-example"
	address = "www.hashicorptest.com"
}

resource "consul_service" "no-deregister" {
	name     = "example-external"
	node     = "${consul_node.external.name}"
	port     = 80

	check {
		// Consul seems to order checks alphabetically by check_id
		check_id = "service:redis2"
		name = "Redis health check"
		http = "https://www.google.com"
		interval = "5s"
		timeout = "1s"
	}

	check {
		check_id = "service:redis1"
		name = "Redis health check"
		http = "https://www.google.com"
		interval = "5s"
		timeout = "1s"
	}
}
`

const testAccConsulServiceCheckID = `
resource "consul_service" "example" {
  name       = "example"
  service_id = "service_id"
  node       = consul_node.example.name
  port       = 80

  check {
    check_id                          = "service:example"
    name                              = "Example health check"
    status                            = "passing"
    http                              = "https://www.hashicorptest.com"
    tls_skip_verify                   = false
    method                            = "PUT"
    interval                          = "5s"
    timeout                           = "1s"
    deregister_critical_service_after = "30s"
  }
}

resource "consul_node" "example" {
name    = "example"
address = "www.example.com"
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

const testAccConsulServiceConfigBasicMeta = `
resource "consul_service" "example" {
  name    = "example"
  node    = "${consul_node.compute.name}"
  port    = 80
  tags    = ["tag0"]
  meta    = {
	test  = "test"
  }
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
	meta    = {
		test  = "test"
	}
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

const testAccDataConsulServiceSameServiceMultipleNodes = `
resource "consul_node" "compute1" {
  name    = "compute-google1"
  address = "www.google.com"
}

resource "consul_service" "google1" {
  name    = "google"
  node    = "${consul_node.compute1.name}"
  port    = 80
  tags    = ["tag0"]

  check {
    check_id                          = "service:redis1"
    name                              = "Redis health check"
    status                            = "passing"
    http                              = "https://www.hashicorptest.com"
    tls_skip_verify                   = false
    method                            = "PUT"
    interval                          = "5s"
    timeout                           = "1s"
    deregister_critical_service_after = "30s"

    header {
      name  = "foo"
      value = ["test"]
    }

    header {
      name  = "bar"
      value = ["test"]
    }
  }
}

resource "consul_node" "compute2" {
  name    = "compute-google2"
  address = "www.google.com"
}

resource "consul_service" "google2" {
  name    = "google"
  node    = "${consul_node.compute2.name}"
  port    = 80
  tags    = ["tag0"]

  check {
    check_id                          = "service:redis1"
    name                              = "Redis health check"
    status                            = "critical"
    http                              = "https://www.hashicorptest.com"
    tls_skip_verify                   = false
    method                            = "PUT"
    interval                          = "5s"
    timeout                           = "1s"
    deregister_critical_service_after = "30s"

    header {
      name  = "foo"
      value = ["test"]
    }

    header {
      name  = "bar"
      value = ["test"]
    }
  }
}
`

// Regression test, creating a service used to make changes to the associated node
// See https://github.com/terraform-providers/terraform-provider-consul/issues/101
const testAccConsulServiceDontOverrideNodeMeta = `
resource "consul_node" "compute" {
	name    = "compute-example"
	address = "www.hashicorptest.com"

	meta = {
	  foo = "bar"
	}
}

resource "consul_service" "example" {
	name = "example"
	node = "${consul_node.compute.name}"
	port = 80
}
`
