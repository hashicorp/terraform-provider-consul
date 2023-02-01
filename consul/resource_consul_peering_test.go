// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulPeering_basic(t *testing.T) {
	providers, _ := startPeeringTestServers(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulPeeringBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_peering.basic", "deleted_at", ""),
					resource.TestCheckResourceAttr("consul_peering.basic", "id", "test"),
					resource.TestCheckResourceAttr("consul_peering.basic", "meta.%", "1"),
					resource.TestCheckResourceAttr("consul_peering.basic", "meta.foo", "bar"),
					resource.TestCheckResourceAttr("consul_peering.basic", "peer_ca_pems.#", "1"),
					resource.TestCheckResourceAttrSet("consul_peering.basic", "peer_id"),
					resource.TestCheckResourceAttr("consul_peering.basic", "peer_name", "test"),
					resource.TestCheckResourceAttr("consul_peering.basic", "peer_server_addresses.#", "1"),
					resource.TestCheckResourceAttr("consul_peering.basic", "peer_server_addresses.0", "127.0.0.1:9503"),
					resource.TestCheckResourceAttrSet("consul_peering.basic", "peer_server_name"),
					resource.TestCheckResourceAttrSet("consul_peering.basic", "peering_token"),
					resource.TestCheckResourceAttr("consul_peering.basic", "state", "ESTABLISHING"),
				),
			},
			{
				Config:                  testAccConsulPeeringBasic,
				ImportState:             true,
				ResourceName:            "consul_peering.basic",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"peering_token"},
			},
		},
	})
}

const testAccConsulPeeringBasic = `
provider "consulremote" {
  address = "http://localhost:9500"
}

resource "consul_peering_token" "basic" {
  provider  = consulremote
  peer_name = "hello-world"
}

resource "consul_peering" "basic" {
  peer_name     = "test"
  peering_token = consul_peering_token.basic.peering_token

  meta = {
    foo = "bar"
  }
}
`
