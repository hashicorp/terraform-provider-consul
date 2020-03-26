package consul

import (
	"fmt"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulNetworkArea_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { skipTestOnConsulCommunityEdition(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccConsulNetworkAreaCheckDestroy,
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
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulNetworkAreaBasic,
				ExpectError: regexp.MustCompile("Failed to create network area: Unexpected response code: 404"),
			},
		},
	})
}

func testAccConsulNetworkAreaCheckDestroy(s *terraform.State) error {
	client := getClient(testAccProvider.Meta())
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
