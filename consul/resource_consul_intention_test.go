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
			resource.TestStep{
				PreConfig: testAccRemoveConsulIntention(t),
				Config:    testAccConsulIntentionConfigBasic,
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

func testAccRemoveConsulIntention(t *testing.T) func() {
	return func() {
		connect := testAccProvider.Meta().(*consulapi.Client).Connect()
		qOpts := &consulapi.QueryOptions{}
		iM := &consulapi.IntentionMatch{
			By:    consulapi.IntentionMatchSource,
			Names: []string{"api"},
		}

		resp, _, err := connect.IntentionMatch(iM, qOpts)
		if err != nil {
			t.Errorf("Failed to retrieve intentions by match err: %v", err)
		}

		intentions, hasMatch := resp["api"]
		if !hasMatch {
			t.Errorf("No intention with source api was found")
		}

		var iid string
		for _, i := range intentions {
			if _, ok := i.Meta["is_tf_acc_test"]; ok {
				iid = i.ID
				break
			}
		}

		if iid == "" {
			t.Errorf("Failed to find the intention created by Terraform")
		}

		_, err = connect.IntentionDelete(iid, &consulapi.WriteOptions{})
		if err != nil {
			t.Errorf("Failed to delete the intention. err: %s", err)
		}
	}
}

const testAccConsulIntentionConfigBasic = `
resource "consul_intention" "example" {
	source_name      = "api"
	destination_name = "db"
	action           = "allow"

	description = "something about example"
	meta {
		foo            = "bar"
		baz            = "bat"
		is_tf_acc_test = "yes"
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
