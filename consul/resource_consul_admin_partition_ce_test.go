package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulAdminParition_CEBasic(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck: func() {
			skipTestOnConsulEnterpriseEdition(t)
		},
		CheckDestroy: testAccCheckConsulACLTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulAdminPartitionBasic,
				ExpectError: regexp.MustCompile(`Unexpected response code: 404 \(Invalid URL path: not a recognized HTTP API endpoint\)`),
			},
		},
	})
}
