package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulPeering_basic(t *testing.T) {
	providers, _ := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsulPeeringMissing,
				ExpectError: regexp.MustCompile(`no peer name "hello" found`),
			},
			{
				Config: testAccDataConsulPeeringBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_peering.basic", "deleted_at", ""),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "exported_service_count", "0"),
					resource.TestCheckResourceAttrSet("data.consul_peering.basic", "id"),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "imported_service_count", "0"),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "meta.%", "1"),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "meta.foo", "bar"),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "peer_ca_pems.#", "0"),
					resource.TestCheckResourceAttrSet("data.consul_peering.basic", "peer_id"),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "peer_name", "test"),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "peer_server_addresses.#", "1"),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "peer_server_addresses.0", "127.0.0.1:8508"),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "peer_server_name", "server.dc2.consul"),
					resource.TestCheckResourceAttr("data.consul_peering.basic", "state", "ESTABLISHING"),
				),
			},
		},
	})
}

const testAccDataConsulPeeringMissing = `
data "consul_peering" "basic" {
  peer_name = "hello"
}
`

const testAccDataConsulPeeringBasic = `
provider "consul" {}

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

data "consul_peering" "basic" {
  peer_name = consul_peering.basic.peer_name
}
`
