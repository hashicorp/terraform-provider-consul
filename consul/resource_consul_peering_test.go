package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulPeering_basic(t *testing.T) {
	providers, _ := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulPeeringBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_peering.basic", "deleted_at", ""),
					resource.TestCheckResourceAttr("consul_peering.basic", "exported_service_count", "0"),
					resource.TestCheckResourceAttr("consul_peering.basic", "id", "test"),
					resource.TestCheckResourceAttr("consul_peering.basic", "imported_service_count", "0"),
					resource.TestCheckResourceAttr("consul_peering.basic", "meta.%", "1"),
					resource.TestCheckResourceAttr("consul_peering.basic", "meta.foo", "bar"),
					resource.TestCheckResourceAttr("consul_peering.basic", "peer_ca_pems.#", "0"),
					resource.TestCheckResourceAttrSet("consul_peering.basic", "peer_id"),
					resource.TestCheckResourceAttr("consul_peering.basic", "peer_name", "test"),
					resource.TestCheckResourceAttr("consul_peering.basic", "peer_server_addresses.#", "1"),
					resource.TestCheckResourceAttr("consul_peering.basic", "peer_server_addresses.0", "127.0.0.1:8508"),
					resource.TestCheckResourceAttr("consul_peering.basic", "peer_server_name", "server.dc2.consul"),
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
  address = "http://localhost:8501"
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
