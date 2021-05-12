package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulNetworkSegments_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulNetworkSegmentsBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_network_segments.test", "segments.#", "1"),
					resource.TestCheckResourceAttr("data.consul_network_segments.test", "segments.0", ""),
				),
			},
		},
	})
}

func TestAccConsulNetworkSegments_CommunityEdition(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulNetworkSegmentsBasic,
				ExpectError: regexp.MustCompile("Failed to get segment list"),
			},
		},
	})
}

func TestAccConsulNetworkSegments_datacenter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccRemoteDatacenterPreCheck(t)
			skipTestOnConsulCommunityEdition(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulNetworkSegmentsDatacenter,
				ExpectError: regexp.MustCompile("Failed to get segment list"),
			},
		},
	})
}

const testAccConsulNetworkSegmentsBasic = `
data "consul_network_segments" "test" {}
`

const testAccConsulNetworkSegmentsDatacenter = `
data "consul_network_segments" "test" {
	datacenter = "dc3"
}
`
