// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConsulNetworkArea_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { skipTestOnConsulCommunityEdition(t) },
		Providers:    providers,
		CheckDestroy: testAccConsulNetworkAreaCheckDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNetworkAreaBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_network_area.test", "peer_datacenter", "foo"),
					resource.TestCheckResourceAttr("consul_network_area.test", "use_tls", "false"),
					resource.TestCheckResourceAttr("consul_network_area.test", "retry_join.#", "1"),
					resource.TestCheckResourceAttr("consul_network_area.test", "retry_join.0", "1.2.3.4"),
				),
			},
			{
				Config: testAccConsulNetworkAreaBasic_update1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_network_area.test", "peer_datacenter", "foo"),
					resource.TestCheckResourceAttr("consul_network_area.test", "use_tls", "true"),
					resource.TestCheckResourceAttr("consul_network_area.test", "retry_join.#", "1"),
					resource.TestCheckResourceAttr("consul_network_area.test", "retry_join.0", "1.2.3.4"),
				),
			},
			{
				Config: testAccConsulNetworkAreaBasic_update2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_network_area.test", "peer_datacenter", "bar"),
					resource.TestCheckResourceAttr("consul_network_area.test", "use_tls", "true"),
					resource.TestCheckResourceAttr("consul_network_area.test", "retry_join.#", "0"),
				),
			},
		},
	})
}

func TestAccConsulNetworkArea_CommunityEdition(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulNetworkAreaBasic,
				ExpectError: regexp.MustCompile("failed to create network area: Unexpected response code: 404"),
			},
		},
	})
}

func TestAccConsulNetworkArea_datacenter(t *testing.T) {
	providers, client := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipTestOnConsulCommunityEdition(t)
		},
		Providers:    providers,
		CheckDestroy: testAccConsulNetworkAreaCheckDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNetworkAreaDatacenter,
				Check: func(s *terraform.State) error {
					test := func(dc, peer string) error {
						c := client.Operator()
						opts := &consulapi.QueryOptions{
							Datacenter: dc,
						}
						area, _, err := c.AreaList(opts)
						if err != nil {
							return err
						}
						if len(area) != 1 {
							return fmt.Errorf("wrong number of network area: %#v", area)
						}
						if area[0].PeerDatacenter != peer {
							return fmt.Errorf("unexpected peer: %s", area[0].PeerDatacenter)
						}
						return nil
					}
					if err := test("dc1", "dc2"); err != nil {
						return err
					}
					return test("dc2", "dc1")
				},
			},
			{
				Config: testAccConsulNetworkAreaBasic_update1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_network_area.test", "peer_datacenter", "foo"),
					resource.TestCheckResourceAttr("consul_network_area.test", "use_tls", "true"),
					resource.TestCheckResourceAttr("consul_network_area.test", "retry_join.#", "1"),
					resource.TestCheckResourceAttr("consul_network_area.test", "retry_join.0", "1.2.3.4"),
				),
			},
			{
				Config: testAccConsulNetworkAreaBasic_update2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_network_area.test", "peer_datacenter", "bar"),
					resource.TestCheckResourceAttr("consul_network_area.test", "use_tls", "true"),
					resource.TestCheckResourceAttr("consul_network_area.test", "retry_join.#", "0"),
				),
			},
		},
	})
}

func testAccConsulNetworkAreaCheckDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		operator := client.Operator()

		qOpts := &consulapi.QueryOptions{}
		areas, _, err := operator.AreaList(qOpts)
		if err != nil {
			return fmt.Errorf("Failed to fetch network areas: %v", err)
		}

		if len(areas) != 0 {
			return fmt.Errorf("Some areas have not been destroyed: %v", areas)
		}

		return nil
	}
}

const testAccConsulNetworkAreaBasic = `
resource "consul_network_area" "test" {
	peer_datacenter = "foo"
	retry_join = ["1.2.3.4"]
}
`

const testAccConsulNetworkAreaBasic_update1 = `
resource "consul_network_area" "test" {
	peer_datacenter = "foo"
	retry_join = ["1.2.3.4"]

	use_tls = true
}
`

const testAccConsulNetworkAreaBasic_update2 = `
resource "consul_network_area" "test" {
	peer_datacenter = "bar"
	retry_join = []

	use_tls = true
}
`

const testAccConsulNetworkAreaDatacenter = `
resource "consul_network_area" "dc1" {
	peer_datacenter = "dc2"
	retry_join = []
}

resource "consul_network_area" "dc2" {
	datacenter      = "dc2"
	peer_datacenter = "dc1"
	retry_join = []
}
`
