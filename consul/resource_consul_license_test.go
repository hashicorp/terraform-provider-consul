package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulLicense(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulLicense,
				ExpectError: regexp.MustCompile(`failed to set license: Unexpected response code: 400 \(Bad request: unknown version: .*\)`),
			},
		},
	})
}

const testAccConsulLicense = `
resource "consul_license" "license" {
	license = "foobar"
}
`
