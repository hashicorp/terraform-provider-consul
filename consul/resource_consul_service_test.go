// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulService_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulServiceDestroy(client),
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
					testAccConsulExternalSource(client),
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
				PreConfig: testAccRemoveConsulService(t, client, "compute-example", "example"),
				Config:    testAccConsulServiceConfigBasicMeta,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "service_id", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "address", "www.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "node", "compute-example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "enable_tag_override", "true"),
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
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulServiceDestroy(client),
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
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulServiceDestroy(client),
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
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulServiceDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulServiceCheck,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "name", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "check.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.check_id", "service:redis1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.name", "Redis health check"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.notes", "Script based health check"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.status", "passing"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.http", "https://www.hashicorptest.com"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.interval", "5s"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.timeout", "1s"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.deregister_critical_service_after", "30s"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.header.#", "2"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.header.344754333.name", "bar"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.header.344754333.value.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.header.344754333.value.0", "test"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.header.2976766922.name", "foo"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.header.2976766922.value.#", "1"),
					resource.TestCheckResourceAttr("consul_service.example", "check.1210848548.header.2976766922.value.0", "test"),
					resource.TestCheckResourceAttr("consul_service.no-deregister", "check.1732085461.deregister_critical_service_after", "30s"),
				),
			},
			{
				Config: testAccConsulServiceCheckID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_service.example", "name", "example"),
					resource.TestCheckResourceAttr("consul_service.example", "port", "80"),
					resource.TestCheckResourceAttr("consul_service.example", "check.#", "1"),
				),
			},
			{
				// Changes in a registered check should be detected
				// See https://github.com/hashicorp/terraform-provider-consul/issues/237
				PreConfig: func() {
					c := client.Catalog()
					_, err := c.Register(&consulapi.CatalogRegistration{
						Node: "example",
						Service: &consulapi.AgentService{
							Service: "example",
							ID:      "service_id",
							Port:    80,
						},
						Checks: []*consulapi.HealthCheck{
							{
								ServiceID: "service_id",
								Node:      "example",
								CheckID:   "service:example",
								Name:      "Example health check",
								Status:    "passing",
								Type:      "http",
								Definition: consulapi.HealthCheckDefinition{
									HTTP:                                   "https://www.test.com",
									TLSSkipVerify:                          false,
									Method:                                 "PUT",
									IntervalDuration:                       5 * time.Second,
									TimeoutDuration:                        1 * time.Second,
									DeregisterCriticalServiceAfterDuration: 30 * time.Second,
								},
							},
						},
						SkipNodeUpdate: true,
					}, nil)

					if err != nil {
						t.Fatalf("failed to change the health-check: %s", err)
					}
				},
				Config:             testAccConsulServiceCheckID,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccConsulServiceCheckOrder(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
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
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulServiceSameServiceMultipleNodes,
			},
		},
	})
}

func TestAccConsulService_nodeDoesNotExist(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers:    providers,
		CheckDestroy: testAccCheckConsulServiceDestroy(client),
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulServiceConfigNoNode,
				ExpectError: regexp.MustCompile(`node does not exist: '*'`),
			},
		},
	})
}

func TestAccConsulService_dontOverrideNodeMeta(t *testing.T) {
	providers, _ := startTestServer(t)

	// This would raise an error if consul_service changed attributes of consul_node
	// since the next plan would not be empty
	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulServiceDontOverrideNodeMeta,
			},
		},
	})
}

func TestAccConsulService_multipleInstances(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulServiceMultipleInstances,
			},
			{
				PreConfig: testAccRemoveConsulService(t, client, "redis/redis1", "redis"),
				Config:    testAccConsulServiceMultipleInstances,
			},
		},
	})
}

func TestAccConsulService_NamespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulServiceNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccConsulService_NamespaceEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulServiceNamespaceEE,
			},
		},
	})
}

func TestAccConsulService_datacenter(t *testing.T) {
	providers, client := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulServiceDatacenter,
				Check: func(s *terraform.State) error {
					test := func(dc string, expected bool) error {
						c := client.Catalog()
						opts := &consulapi.QueryOptions{
							Datacenter: dc,
						}
						svc, _, err := c.Services(opts)
						if err != nil {
							return err
						}

						var found bool
						for s := range svc {
							if s == "test" {
								found = true
							}
						}
						if found != expected {
							return fmt.Errorf("unexpected result: %v, %v", found, expected)
						}
						return nil
					}

					if err := test("dc1", false); err != nil {
						return err
					}
					return test("dc2", true)
				},
			},
		},
	})
}

func testAccConsulExternalSource(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
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
}

func testAccCheckConsulServiceDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
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
}

func testAccRemoveConsulService(t *testing.T, client *consulapi.Client, node, serviceID string) func() {
	return func() {
		catalog := client.Catalog()
		wOpts := &consulapi.WriteOptions{}
		dereg := &consulapi.CatalogDeregistration{
			Node:      node,
			ServiceID: serviceID,
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
  name                = "example"
  node                = "${consul_node.compute.name}"
  port                = 80
  tags                = ["tag0"]
  enable_tag_override = true

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
// See https://github.com/hashicorp/terraform-provider-consul/issues/101
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

// Removing one instance of a service used to make Terraform complains when refreshing the plan
// See https://github.com/hashicorp/terraform-provider-consul/issues/146
const testAccConsulServiceMultipleInstances = `
resource "consul_node" "redis1" {
	name = "redis/redis1"
	address = "hostname1"
}
resource "consul_node" "redis2" {
	name = "redis/redis2"
	address = "hostname2"
}

resource "consul_service" "redis1" {
	name = "redis"
	node = consul_node.redis1.name
	port = 6379

	check {
		check_id                          = "service:redis1"
		name                              = "Redis health check"
		tcp                               = "127.0.0.1:6379"
		interval                          = "5s"
		timeout                           = "1s"
		deregister_critical_service_after = "30s"
	}
}

resource "consul_service" "redis2" {
	name = "redis"
	node = consul_node.redis2.name
	port = 6379

	check {
		check_id                          = "service:redis1"
		name                              = "Redis health check"
		tcp                               = "127.0.0.1:6379"
		interval                          = "5s"
		timeout                           = "1s"
		deregister_critical_service_after = "30s"
	}
}
`

const testAccConsulServiceNamespaceCE = `
resource "consul_node" "test" {
  name    = "test"
  address = "test.com"
}

resource "consul_service" "test" {
  name      = "test"
  namespace = "test"
  node      = consul_node.test.name
  port      = 80
}
`

const testAccConsulServiceNamespaceEE = `
resource "consul_namespace" "test" {
  name      = "test-service"
}

resource "consul_node" "test" {
  name    = "test"
  address = "test.com"
}

resource "consul_service" "test" {
  name      = "test"
  namespace = consul_namespace.test.name
  node      = consul_node.test.name
  port      = 80
}

data "consul_service" "test" {
  name = consul_service.test.name

  query_options {
    namespace = consul_namespace.test.name
  }
}
`

const testAccConsulServiceDatacenter = `
resource "consul_node" "test" {
	datacenter = "dc2"
	name       = "test"
	address    = "test.com"
}

resource "consul_service" "test" {
	datacenter = "dc2"
	name       = "test"
	node       = consul_node.test.name
	port       = 80
}
`
