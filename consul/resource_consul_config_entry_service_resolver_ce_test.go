// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntryServiceResolverCETest(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testConsulConfigEntryServiceResolverCEWithRedirect,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "name", "consul-service-resolver-1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "meta.key", "value"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "connect_timeout", "10s"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "request_timeout", "10s"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "subsets.#", "2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "subsets.1420492792.name", "v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "subsets.1420492792.filter", "Service.Meta.version == v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "subsets.853348911.only_passing", "true"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "subsets.853348911.name", "v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "subsets.853348911.filter", "Service.Meta.version == v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "subsets.853348911.only_passing", "true"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "default_subset", "v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "redirect.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "redirect.1000671749.datacenter", "dc1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "load_balancer.1867531597.policy", "ring_hash"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "load_balancer.1867531597.ring_hash_config.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "load_balancer.1867531597.ring_hash_config.1100849243.minimum_ring_size", "3"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "load_balancer.1867531597.ring_hash_config.1100849243.maximum_ring_size", "10"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "load_balancer.1867531597.hash_policies.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "load_balancer.1867531597.hash_policies.0.field", "header"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.foo", "load_balancer.1867531597.hash_policies.0.field_value", "x-user-id"),
				),
			},
			{
				Config: testConsulConfigEntryServiceResolverCEWithFailover,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "name", "consul-service-resolver-2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "meta.key", "value"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "connect_timeout", "10s"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "request_timeout", "10s"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "subsets.#", "2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "subsets.1420492792.name", "v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "subsets.1420492792.filter", "Service.Meta.version == v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "subsets.1420492792.only_passing", "true"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "subsets.853348911.name", "v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "subsets.853348911.filter", "Service.Meta.version == v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "subsets.853348911.only_passing", "true"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "default_subset", "v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.#", "3"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1078808137.subset_name", "*"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1078808137.service", "backend"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1078808137.datacenters.0", "dc3"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1078808137.datacenters.1", "dc4"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1326363731.subset_name", "v2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1326363731.service", "frontend"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1326363731.datacenters.0", "dc2"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1420169321.subset_name", "v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1420169321.targets.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1420169321.targets.0.service_subset", "v1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "failover.1420169321.targets.0.datacenter", "dc1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "load_balancer.1867531597.policy", "ring_hash"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "load_balancer.1867531597.ring_hash_config.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "load_balancer.1867531597.ring_hash_config.1100849243.minimum_ring_size", "3"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "load_balancer.1867531597.ring_hash_config.1100849243.maximum_ring_size", "10"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "load_balancer.1867531597.hash_policies.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "load_balancer.1867531597.hash_policies.0.field", "header"),
					resource.TestCheckResourceAttr("consul_config_entry_service_resolver.bar", "load_balancer.1867531597.hash_policies.0.field_value", "x-user-id"),
				),
			},
		},
	})
}

const testConsulConfigEntryServiceResolverCEWithRedirect = `
resource "consul_config_entry_service_resolver" "foo" {
	name = "consul-service-resolver-1"
	meta = {
		key: "value"
	}
	connect_timeout = "10s"
	request_timeout = "10s"
	subsets {
		name = "v1"
		filter = "Service.Meta.version == v1"
		only_passing = true
	}
	subsets {
		name = "v2"
		filter = "Service.Meta.version == v2"
		only_passing = true
	}
	default_subset = "v1"
	redirect {
		datacenter = "dc1"
	}
	load_balancer {
		policy = "ring_hash"
		ring_hash_config {
			minimum_ring_size = 3
			maximum_ring_size = 10
		}
		hash_policies {
			field = "header"
			field_value = "x-user-id"
		}
	}
}
`

const testConsulConfigEntryServiceResolverCEWithFailover = `
resource "consul_config_entry_service_resolver" "bar" {
	name = "consul-service-resolver-2"
	meta = {
		key: "value"
	}
	connect_timeout = "10s"
	request_timeout = "10s"
	subsets {
		name = "v1"
		filter = "Service.Meta.version == v1"
		only_passing = true
	}
	subsets {
		name = "v2"
		filter = "Service.Meta.version == v2"
		only_passing = true
	}
	default_subset = "v1"
	failover {
		subset_name  = "*"
		service      = "backend"
		datacenters = ["dc3", "dc4"]
	}
	failover {
		subset_name = "v2"
		service = "frontend"
		datacenters = ["dc2"]
	}
	failover {
		subset_name = "v1"
		targets {
			service_subset = "v1"
			datacenter = "dc1"
		}
	}
	load_balancer {
		policy = "ring_hash"
		ring_hash_config {
			minimum_ring_size = 3
			maximum_ring_size = 10
		}
		hash_policies {
			field = "header"
			field_value = "x-user-id"
		}
	}
}
`
