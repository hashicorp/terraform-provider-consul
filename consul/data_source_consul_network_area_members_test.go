package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulNetworkAreaMembers_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNetworkAreaMembersBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_network_area_members.test", "members.#", "1"),
					resource.TestCheckResourceAttr("data.consul_network_area_members.test", "members.0.address", "127.0.0.1"),
					resource.TestCheckResourceAttr("data.consul_network_area_members.test", "members.0.port", "8300"),
					resource.TestCheckResourceAttr("data.consul_network_area_members.test", "members.0.datacenter", "dc1"),
					resource.TestCheckResourceAttr("data.consul_network_area_members.test", "members.0.role", "server"),
					resource.TestCheckResourceAttr("data.consul_network_area_members.test", "members.0.protocol", "2"),
					resource.TestCheckResourceAttr("data.consul_network_area_members.test", "members.0.status", "alive"),
					resource.TestCheckResourceAttr("data.consul_network_area_members.test", "members.0.rtt", "0"),
				),
			},
		},
	})
}

func TestAccConsulNetworkAreaMembers_CommunityEdition(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulNetworkAreaMembers_CommunityEdition,
				ExpectError: regexp.MustCompile("Failed to fetch the list of members"),
			},
		},
	})
}

func TestAccConsulNetworkAreaMembers_datacenter(t *testing.T) {
	providers, _ := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipTestOnConsulCommunityEdition(t)
		},
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNetworkAreaMembers_datacenter,
				Check:  resource.TestCheckResourceAttr("data.consul_network_area_members.test", "members.0.datacenter", "dc2"),
			},
		},
	})
}

const testAccConsulNetworkAreaMembersBasic = `
resource "consul_network_area" "test" {
	peer_datacenter = "foo"
	retry_join = ["1.2.3.4"]
}

data "consul_network_area_members" "test" {
	uuid = consul_network_area.test.id
}`

const testAccConsulNetworkAreaMembers_CommunityEdition = `
data "consul_network_area_members" "test" {
	uuid = "1.2.3.4"
}`

const testAccConsulNetworkAreaMembers_datacenter = `
resource "consul_network_area" "test" {
	datacenter = "dc2"
	peer_datacenter = "foo"
	retry_join = ["1.2.3.4"]
}

data "consul_network_area_members" "test" {
	datacenter = "dc2"
	uuid = consul_network_area.test.id
}`
