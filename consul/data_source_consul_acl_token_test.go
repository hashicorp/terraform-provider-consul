package consul

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					testAccCheckDataSourceValue("data.consul_acl_token.read", "local", "true"),
					testAccCheckTokenExistsAndValidUUID("data.consul_acl_token.read", "secret_id"),
				),
			},
		},
	})
}

func testAccCheckTokenExistsAndValidUUID(n string, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		secretID := rs.Primary.Attributes[attr]
		r := regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}")
		if !r.MatchString(secretID) {
			return fmt.Errorf("No valid UUID format %q", secretID)
		}
		return nil
	}
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
