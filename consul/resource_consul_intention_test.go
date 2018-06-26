package consul

import (
	"fmt"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccConsulIntention_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulIntentionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulIntentionConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_intention.example", "source_name", "api"),
					resource.TestCheckResourceAttr("consul_intention.example", "destination_name", "db"),
					resource.TestCheckResourceAttr("consul_intention.example", "action", "allow"),
					resource.TestCheckResourceAttr("consul_intention.example", "description", "something about example"),
					resource.TestCheckResourceAttr("consul_intention.example", "meta.foo", "bar"),
					resource.TestCheckResourceAttr("consul_intention.example", "meta.baz", "bat"),
				),
			},
		},
	})
}

func TestAccConsulIntention_badAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulIntentionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccConsulIntentionConfigBadAction,
				ExpectError: regexp.MustCompile("expected action to be one of"),
			},
		},
	})
}

func testAccCheckConsulIntentionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*consulapi.Client)

	qOpts := consulapi.QueryOptions{}
	intentions, _, err := client.Connect().Intentions(&qOpts)
	if err != nil {
		return fmt.Errorf("Failed to retrieve intentions: %v", err)
	}

	if len(intentions) > 0 {
		return fmt.Errorf("Intentions still exsist: %v", intentions)
	}

	return nil
}

const testAccConsulIntentionConfigBasic = `
resource "consul_intention" "example" {
	source_name      = "api"
	destination_name = "db"
	action           = "allow"

	description = "something about example"
	meta {
		foo = "bar"
		baz = "bat"
	}
}
`

const testAccConsulIntentionConfigBadAction = `
resource "consul_intention" "example" {
	source_name      = "api"
	destination_name = "db"
	action           = "foobar"
}
`
