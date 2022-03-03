package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulLicense_CommunityEdition(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipTestOnConsulEnterpriseEdition(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulLicense,
				ExpectError: regexp.MustCompile("failed to set license: Unexpected response code: 404"),
			},
		},
	})
}

func TestAccConsulLicense_EnterpriseEdition(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipTestOnConsulCommunityEdition(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulLicense,
				ExpectError: regexp.MustCompile("failed to set license: Unexpected response code: 405"),
			},
		},
	})
}

const testAccConsulLicense = `
resource "consul_license" "license" {
	license = "foobar"
}
`
