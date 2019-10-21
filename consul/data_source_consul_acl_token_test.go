package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataACLToken_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataACLTokenConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_acl_token.read", "description", "test"),
					testAccCheckDataSourceValue("data.consul_acl_token.read", "policies.#", "1"),
					testAccCheckDataSourceValue("data.consul_acl_token.read", "policies.0.name", "test"),
					testAccCheckDataSourceValue("data.consul_acl_token.read", "policies.0.id", "<any>"),
					testAccCheckDataSourceValue("data.consul_acl_token.read", "local", "true"),
				),
			},
		},
	})
}

const testAccDataACLTokenConfig = `
resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = ["${consul_acl_policy.test.name}"]
	local = true
}

data "consul_acl_token" "read" {
    accessor_id = "${consul_acl_token.test.id}"
}
`
