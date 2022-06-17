package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulPeerings_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulPeeringsNone,
			},
			{
				Config: testAccDataConsulPeeringsBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "id", "peers"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.#", "2"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.deleted_at", ""),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.0.id"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.meta.%", "0"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.name", "hello-world"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.partition", ""),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.peer_ca_pems.#", "0"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.peer_id", ""),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.peer_server_addresses.#", "0"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.peer_server_name", ""),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.0.state", "INITIAL"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.deleted_at", ""),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.1.id"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.meta.%", "1"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.meta.foo", "bar"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.name", "test"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.partition", ""),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.peer_ca_pems.#", "0"),
					resource.TestCheckResourceAttrSet("data.consul_peerings.basic", "peers.1.peer_id"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.peer_server_addresses.#", "1"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.peer_server_addresses.0", "127.0.0.1:8300"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.peer_server_name", "server.dc1.consul"),
					resource.TestCheckResourceAttr("data.consul_peerings.basic", "peers.1.state", "INITIAL"),
				),
			},
		},
	})
}

const testAccDataConsulPeeringsNone = `
data "consul_peerings" "basic" {}
`

const testAccDataConsulPeeringsBasic = `
resource "consul_peering_token" "basic" {
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
  datacenter = consul_peering.basic.datacenter
}
`
