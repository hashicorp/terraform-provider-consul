package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulLicense_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulLicense,
				// Setting the Consul license will fail on the Community Edition
				ExpectError: regexp.MustCompile("failed to set license"),
			},
		},
	})
}

const testAccConsulLicense = `
resource "consul_license" "license" {
	license = "foobar"
}
`
