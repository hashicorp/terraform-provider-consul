package consul

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccCheckConsulACLAgentTokenDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*consulapi.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "consul_acl_agent_token" {
			continue
		}
		aclEntry, _, err := client.ACL().Info(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if aclEntry != nil {
			return fmt.Errorf("ACL agent token %q still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccConsulACLAgentToken_basic(t *testing.T) {
	token, _ := uuid.GenerateUUID()
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckConsulACLAgentTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testResourceACLAgentTokenConfig_basic(token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_agent_token.test", "token", token),
				),
			},
		},
	})
}

func testResourceACLAgentTokenConfig_basic(token string) string {
	return fmt.Sprintf(`
resource "consul_acl" "test" {
	uuid = "%s"
	name = "Agent Token"
	type = "client"
	rules = "node \"\" { policy = \"write\" } service \"\" { policy = \"read\" }"
}

resource "consul_acl_agent_token" "test" {
	token = "${consul_acl.test.token}"
}`, token)
}
