package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccConsulPreparedQuery_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulPreparedQueryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulPreparedQueryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulPreparedQueryExists(),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "stored_token", "pq-token"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "service", "redis"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "near", "_agent"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "tags.#", "1"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "only_passing", "true"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "connect", "false"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "failover.0.nearest_n", "3"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "failover.0.datacenters.#", "2"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "template.0.type", "name_prefix_match"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "template.0.regexp", "hello"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "dns.0.ttl", "8m"),
				),
			},
			{
				Config: testAccConsulPreparedQueryConfigUpdate1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulPreparedQueryExists(),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "name", "baz"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "stored_token", "pq-token-updated"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "service", "memcached"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "near", "node1"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "tags.#", "2"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "only_passing", "false"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "connect", "true"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "failover.0.nearest_n", "2"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "failover.0.datacenters.#", "1"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "template.0.regexp", "goodbye"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "dns.0.ttl", "16m"),
				),
			},
			{
				PreConfig:          testAccConsulPreparedQueryNearestN(t),
				Config:             testAccConsulPreparedQueryConfigUpdate1,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConsulPreparedQueryConfigUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulPreparedQueryExists(),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "stored_token", ""),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "near", ""),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "tags.#", "0"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "failover.#", "0"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "template.#", "0"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "dns.#", "0"),
				),
			},
		},
	})
}

func TestAccConsulPreparedQuery_import(t *testing.T) {
	checkFn := func(s []*terraform.InstanceState) error {
		// Expect, 1 resource in state, and route count to be 1
		if len(s) != 1 {
			return fmt.Errorf("bad state: %s", s)
		}
		v, ok := s[0].Attributes["name"]
		if !ok || v != "foo" {
			return fmt.Errorf("bad name: %s", s)
		}
		v, ok = s[0].Attributes["stored_token"]
		if !ok || v != "pq-token" {
			return fmt.Errorf("bad stored_token: %s", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulPreparedQueryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulPreparedQueryConfig,
			},
			resource.TestStep{
				ResourceName:     "consul_prepared_query.foo",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

func checkPreparedQueryExists(s *terraform.State) bool {
	rn, ok := s.RootModule().Resources["consul_prepared_query.foo"]
	if !ok {
		return false
	}
	id := rn.Primary.ID

	client := testAccProvider.Meta().(*consulapi.Client).PreparedQuery()
	opts := &consulapi.QueryOptions{Datacenter: "dc1"}
	pq, _, err := client.Get(id, opts)
	return err == nil && pq != nil
}

func testAccCheckConsulPreparedQueryDestroy(s *terraform.State) error {
	if checkPreparedQueryExists(s) {
		return fmt.Errorf("Prepared query 'foo' still exists")
	}
	return nil
}

func testAccCheckConsulPreparedQueryExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !checkPreparedQueryExists(s) {
			return fmt.Errorf("Prepared query 'foo' does not exist")
		}
		return nil
	}
}

func testAccConsulPreparedQueryNearestN(t *testing.T) func() {
	return func() {
		client := testAccProvider.Meta().(*consulapi.Client)
		wOpts := &consulapi.WriteOptions{}
		qOpts := &consulapi.QueryOptions{}

		queries, _, err := client.PreparedQuery().List(qOpts)
		if err != nil {
			t.Fatalf("Failed to fetch prepared queries: %v", err)
		}
		if len(queries) != 1 {
			t.Fatal("Should have exactly one query")
		}

		pq := queries[0]

		// We change the value of nearest_n so the new plan should be non-empty
		pq.Service.Failover.NearestN = 1

		_, err = client.PreparedQuery().Update(pq, wOpts)
		if err != nil {
			t.Fatalf("Failed to update prepared query: %v", err)
		}
	}
}

const testAccConsulPreparedQueryConfig = `
resource "consul_prepared_query" "foo" {
	name = "foo"
	stored_token = "pq-token"
	service = "redis"
	tags = ["prod"]
	near = "_agent"
	only_passing = true

	failover {
		nearest_n = 3
		datacenters = ["dc1", "dc2"]
	}

	template {
		type = "name_prefix_match"
		regexp = "hello"
	}

	dns {
		ttl = "8m"
	}
}
`

const testAccConsulPreparedQueryConfigUpdate1 = `
resource "consul_prepared_query" "foo" {
	name = "baz"
	stored_token = "pq-token-updated"
	service = "memcached"
	tags = ["prod","sup"]
	near = "node1"
	only_passing = false
	connect = true

	failover {
		nearest_n = 2
		datacenters = ["dc2"]
	}

	template {
		type = "name_prefix_match"
		regexp = "goodbye"
	}

	dns {
		ttl = "16m"
	}
}
`

const testAccConsulPreparedQueryConfigUpdate2 = `
resource "consul_prepared_query" "foo" {
	name = "baz"
	service = "memcached"
}
`
