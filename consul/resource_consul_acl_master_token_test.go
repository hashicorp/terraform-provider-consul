package consul

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccConsulACLMasterToken_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLMasterTokenConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("consul_acl_master_token.test", "token"),
				),
			},
		},
	})
}

func testResourceACLMasterTokenConfig_basic() string {
	return `
resource "consul_acl_master_token" "test" {
}`
}
