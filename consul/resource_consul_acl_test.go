package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccCheckConsulACLDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*consulapi.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "consul_acl" {
			continue
		}
		secret, _, err := client.ACL().Info(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if secret != nil {
			return fmt.Errorf("ACL %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccConsulACL_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceTokenConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_acl.test", "type", "client"),
					resource.TestCheckResourceAttr("consul_acl.test", "rules", "node \"\" { policy = \"read\" }"),
				),
			},
		},
	})
}

func testResourceTokenConfig_basic() string {
	return `
resource "consul_acl" "test" {
	name = "test"
	type = "client"
	rules = "node \"\" { policy = \"read\" }"
}`
}

func TestAccConsulACL_uuid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceTokenConfig_uuid(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_acl.test", "uuid", "a49b6f0a-a939-4966-a7e7-c7177a103653"),
					resource.TestCheckResourceAttr("consul_acl.test", "type", "client"),
					resource.TestCheckResourceAttr("consul_acl.test", "rules", "node \"\" { policy = \"read\" }"),
				),
			},
		},
	})
}

func testResourceTokenConfig_uuid() string {
	return `
resource "consul_acl" "test" {
	uuid = "a49b6f0a-a939-4966-a7e7-c7177a103653"
	name = "test"
	type = "client"
	rules = "node \"\" { policy = \"read\" }"
}`
}

func TestAccConsulACL_management(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceTokenConfig_management(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_acl.test", "type", "management"),
					resource.TestCheckResourceAttr("consul_acl.test", "rules", ""),
				),
			},
		},
	})
}

func testResourceTokenConfig_management() string {
	return `
resource "consul_acl" "test" {
	name = "test"
	type = "management"
}`
}
