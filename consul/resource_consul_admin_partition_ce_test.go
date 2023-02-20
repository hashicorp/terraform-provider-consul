// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulAdminParition_CEBasic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck: func() {
			skipTestOnConsulEnterpriseEdition(t)
		},
		CheckDestroy: testAccCheckConsulACLTokenDestroy(client),
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulAdminPartitionBasic,
				ExpectError: regexp.MustCompile(`Unexpected response code: 404 \(Invalid URL path: not a recognized HTTP API endpoint\)`),
			},
		},
	})
}
