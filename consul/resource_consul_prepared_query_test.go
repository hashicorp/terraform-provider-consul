package consul

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "ignore_check_ids.#", "3"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "ignore_check_ids.0", "1"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "ignore_check_ids.1", "2"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "ignore_check_ids.2", "3"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "node_meta.%", "1"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "node_meta.foo", "bar"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "service_meta.%", "1"),
					resource.TestCheckResourceAttr("consul_prepared_query.foo", "service_meta.spam", "ham"),
					testAccCheckConsulPreparedQueryAttributes,
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
			{
				Config: testAccConsulPreparedQueryConfig,
			},
			{
				ResourceName:     "consul_prepared_query.foo",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

func TestAccConsulPreparedQuery_blocks(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulPreparedQueryBlocks,
			},
			{
				Config: testAccConsulPreparedQueryBlocks2,
			},
			{
				Config: testAccConsulPreparedQueryBlocks3,
			},
			{
				Config: testAccConsulPreparedQueryBlocks4,
			},
		},
	})
}

func TestAccConsulPreparedQuery_datacenter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccRemoteDatacenterPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulPreparedQueryDatacenter,
				Check: func(s *terraform.State) error {
					test := func(dc string) error {
						c := getTestClient(testAccProvider.Meta()).PreparedQuery()
						opts := &api.QueryOptions{
							Datacenter: dc,
						}
						pq, _, err := c.List(opts)
						if err != nil {
							return err
						}

						if len(pq) != 1 {
							return fmt.Errorf("wrong number of prepared queries: %#v", pq)
						}
						if pq[0].Name != dc {
							return fmt.Errorf("unknown prepared query %q in datacenter %q", pq[0].Name, dc)
						}
						return nil
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

func getPreparedQuery(s *terraform.State) (*api.PreparedQueryDefinition, error) {
	rn, ok := s.RootModule().Resources["consul_prepared_query.foo"]
	if !ok {
		return nil, fmt.Errorf("Counld not find resource in state")
	}
	id := rn.Primary.ID

	c := getTestClient(testAccProvider.Meta())
	client := c.PreparedQuery()
	opts := &api.QueryOptions{Datacenter: "dc1"}
	pq, _, err := client.Get(id, opts)
	if len(pq) != 1 {
		return nil, fmt.Errorf("Wrong number of prepared queries")
	}
	return pq[0], err
}

func checkPreparedQueryExists(s *terraform.State) bool {
	pq, err := getPreparedQuery(s)
	return err == nil && pq != nil
}

func testAccCheckConsulPreparedQueryAttributes(s *terraform.State) error {
	pq, err := getPreparedQuery(s)

	if err != nil {
		return err
	}

	if pq.Token != "" {
		return fmt.Errorf("Wrong value for 'stored_token': %v", pq.Token)
	}
	if pq.Service.Near != "" {
		return fmt.Errorf("Wrong value for 'near': %v", pq.Service.Near)
	}
	if len(pq.Service.Tags) != 0 {
		return fmt.Errorf("Wrong value for 'tags': %v", pq.Service.Tags)
	}
	if !reflect.DeepEqual(pq.Service.IgnoreCheckIDs, []string{"1", "2", "3"}) {
		return fmt.Errorf("Wrong value for 'ignore_check_ids': %v", pq.Service.IgnoreCheckIDs)
	}
	if !reflect.DeepEqual(pq.Service.ServiceMeta, map[string]string{"spam": "ham"}) {
		return fmt.Errorf("Wrong value for 'service_meta': %v", pq.Service.ServiceMeta)
	}
	if !reflect.DeepEqual(pq.Service.NodeMeta, map[string]string{"foo": "bar"}) {
		return fmt.Errorf("Wrong value for 'node_meta': %v", pq.Service.NodeMeta)
	}
	return nil
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
		client := getTestClient(testAccProvider.Meta())
		wOpts := &api.WriteOptions{}
		qOpts := &api.QueryOptions{}

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
	name             = "baz"
	service          = "memcached"
	ignore_check_ids = ["1", "2", "3"]

	node_meta = {
		foo = "bar"
	}

	service_meta = {
		spam = "ham"
	}
}
`

const testAccConsulPreparedQueryBlocks = `
resource "consul_prepared_query" "foo" {
	name = "foo"
	stored_token = "pq-token"
	service = "redis"
	tags = ["prod"]
	near = "_agent"
	only_passing = true

	failover {
		nearest_n = 0
		datacenters = ["dc1", "dc2"]
	}
}
`

const testAccConsulPreparedQueryBlocks2 = `
resource "consul_prepared_query" "foo" {
	name = "foo"
	stored_token = "pq-token"
	service = "redis"
	tags = ["prod"]
	near = "_agent"
	only_passing = true

	failover {
		nearest_n = 0
		datacenters = []
	}
}
`

const testAccConsulPreparedQueryBlocks3 = `
resource "consul_prepared_query" "foo" {
	name = "foo"
	stored_token = "pq-token"
	service = "redis"
	tags = ["prod"]
	near = "_agent"
	only_passing = true

	dns {
		ttl = ""
	}
}
`

const testAccConsulPreparedQueryBlocks4 = `
resource "consul_prepared_query" "foo" {
	name = "foo"
	stored_token = "pq-token"
	service = "redis"
	tags = ["prod"]
	near = "_agent"
	only_passing = true

	template {
		type   = ""
		regexp = ""
	}
}
`

const testAccConsulPreparedQueryDatacenter = `
resource "consul_prepared_query" "dc1" {
	name = "dc1"
	service = "redis"
}

resource "consul_prepared_query" "dc2" {
	datacenter = "dc2"
	name       = "dc2"
	service    = "redis"
}
`
