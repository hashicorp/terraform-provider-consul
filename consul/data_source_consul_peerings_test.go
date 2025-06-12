// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataConsulPeerings_basic(t *testing.T) {
	providers, _ := startPeeringTestServers(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulPeeringsNone,
			},
			{
				// The peering takes a bit of time to be established so we
				// expect the state to change between the apply and the refresh
				ExpectNonEmptyPlan: true,
				Config:             testAccDataConsulPeeringsBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "id", "peers"),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.#"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.deleted_at", ""),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.0.id"),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.0.meta.%"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.meta.foo", "bar"),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.0.name"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.peer_ca_pems.#", "1"),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.0.peer_id"),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.0.peer_server_addresses.#"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.peer_server_addresses.0", "127.0.0.1:9503"),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.0.peer_server_name"),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.0.state"),
				),
			},
		},
	})
}

const testAccDataConsulPeeringsNone = `
data "consul_peerings" "basic" {}
`

const testAccDataConsulPeeringsBasic = `
provider "consul" {}

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

data "consul_peerings" "basic" {
  depends_on = [consul_peering.basic]
}
`
