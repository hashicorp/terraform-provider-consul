// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulService_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulServiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_service.read", "datacenter", "dc1"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.#", "1"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.address", "<all>"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.create_index", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.enable_tag_override", "false"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.id", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.modify_index", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.name", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.node_address", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.node_id", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.node_meta.%", "2"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.node_meta.consul-network-segment", ""),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.node_name", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.port", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.tagged_addresses.%", "4"),
					testAccCheckDataSourceValue("data.consul_service.read", "service.0.tags.#", "0"),
				),
			},
			{
				Config: testAccDataConsulServiceConfig_wrongFilter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_service.read", "service.#", "0"),
				),
			},
		},
	})
}

func TestAccDataConsulService_filtered(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulServiceFilteredConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_service.read_f", "datacenter", "dc1"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.#", "1"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.address", "192.168.10.10"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.create_index", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.enable_tag_override", "true"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.id", "redis2"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.modify_index", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.name", "redis"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.node_name", "foobar_dummy"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.port", "<any>"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.tags.#", "2"),
					testAccCheckDataSourceValue("data.consul_service.read_f", "service.0.meta.test", "test"),
				),
			},
		},
	})
}

func TestAccDataConsulService_alias(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulServiceAlias,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_catalog_service.read", "service.#", "1"),
				),
			},
		},
	})
}

func TestAccDataConsulService_namespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsulServiceNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccDataConsulService_namespaceEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulServiceNamespaceEE,
			},
		},
	})
}

const testAccDataConsulServiceConfig = `
data "consul_service" "read" {
  query_options {
    allow_stale = true
    require_consistent = false
    token = ""
    wait_index = 0
    wait_time = "1m"
  }

  name   = "consul"
  filter = "ServiceName == consul"
}
`

const testAccDataConsulServiceConfig_wrongFilter = `
data "consul_service" "read" {
  query_options {
    allow_stale = true
    require_consistent = false
    token = ""
    wait_index = 0
    wait_time = "1m"
  }

  name   = "consul"
  filter = "ServiceName != consul"
}
`

const testAccDataConsulServiceFilteredConfig = `
resource "consul_node" "node" {
  address = "192.168.10.10"
  name    = "foobar_dummy"
}

resource "consul_service" "service1" {
	node = "${consul_node.node.name}"
	datacenter = "dc1"

	service_id = "redis1"
	name       = "redis"
	port       = 8000
	tags       = ["v1"]
}

resource "consul_service" "service2" {
	node = "${consul_node.node.name}"
	datacenter = "dc1"

	service_id          = "redis2"
	name                = "redis"
	port                = 8000
	tags                = ["master", "v1"]
	enable_tag_override = true

	meta       = {
		test  = "test"
	}
}

data "consul_service" "read_f" {
  query_options {
    allow_stale = true
    require_consistent = false
    token = ""
    wait_index = 0
	wait_time = "1m"
  }

  name = "redis"
  tag = "master"
  datacenter = "${consul_service.service2.datacenter}"
}
`

const testAccDataConsulServiceAlias = `
data "consul_catalog_service" "read" {
  name = "consul"
}
`

const testAccDataConsulServiceNamespaceCE = `
data "consul_catalog_service" "read" {
  name = "consul"

  query_options {
    namespace = "test-data-service"
  }
}
`

const testAccDataConsulServiceNamespaceEE = `
data "consul_catalog_service" "read" {
  name = "consul"

  query_options {
    namespace = "test-data-service"
  }
}
`
