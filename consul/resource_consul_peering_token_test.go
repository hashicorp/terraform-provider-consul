package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulPeeringToken_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulPeeringTokenBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_peering_token.basic", "id", "hello-world"),
					resource.TestCheckResourceAttr("consul_peering_token.basic", "peer_name", "hello-world"),
					resource.TestCheckResourceAttrSet("consul_peering_token.basic", "peering_token"),
				),
			},
		},
	})
}

const testAccConsulPeeringTokenBasic = `
resource "consul_peering_token" "basic" {
  peer_name = "hello-world"
}
`
