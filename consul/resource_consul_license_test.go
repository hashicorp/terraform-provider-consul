// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConsulLicense_CommunityEdition(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipTestOnConsulEnterpriseEdition(t)
		},
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulLicense,
				ExpectError: regexp.MustCompile("failed to set license: Unexpected response code: 404"),
			},
		},
	})
}

func TestAccConsulLicense_EnterpriseEdition(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipTestOnConsulCommunityEdition(t)
		},
		Providers: providers,
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
