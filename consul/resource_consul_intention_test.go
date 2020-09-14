package consul

import (
	"fmt"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulIntention_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulIntentionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulIntentionConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_intention.example", "source_name", "api"),
					resource.TestCheckResourceAttr("consul_intention.example", "datacenter", "dc1"),
					resource.TestCheckResourceAttr("consul_intention.example", "destination_name", "db"),
					resource.TestCheckResourceAttr("consul_intention.example", "action", "allow"),
					resource.TestCheckResourceAttr("consul_intention.example", "description", "something about example"),
					resource.TestCheckResourceAttr("consul_intention.example", "meta.foo", "bar"),
					resource.TestCheckResourceAttr("consul_intention.example", "meta.baz", "bat"),
				),
			},
			{
				PreConfig: testAccRemoveConsulIntention(t),
				Config:    testAccConsulIntentionConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_intention.example", "source_name", "api"),
					resource.TestCheckResourceAttr("consul_intention.example", "datacenter", "dc1"),
					resource.TestCheckResourceAttr("consul_intention.example", "destination_name", "db"),
					resource.TestCheckResourceAttr("consul_intention.example", "action", "allow"),
					resource.TestCheckResourceAttr("consul_intention.example", "description", "something about example"),
					resource.TestCheckResourceAttr("consul_intention.example", "meta.foo", "bar"),
					resource.TestCheckResourceAttr("consul_intention.example", "meta.baz", "bat"),
				),
			},
			{
				Config:            testAccConsulIntentionConfigBasic,
				ResourceName:      "consul_intention.example",
				ImportState:       true,
				ImportStateVerify: true,
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
			{
				Config:      testAccConsulIntentionConfigBadAction,
				ExpectError: regexp.MustCompile("expected action to be one of"),
			},
		},
	})
}

func TestAccConsulIntention_namespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConsulIntentionConfigNamespace,
			},
		},
	})
}

func TestAccConsulIntention_namespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConsulIntentionConfigNamespace,
			},
		},
	})
}

func testAccCheckConsulIntentionDestroy(s *terraform.State) error {
	client := getClient(testAccProvider.Meta())

	qOpts := consulapi.QueryOptions{}
	intentions, _, err := client.Connect().Intentions(&qOpts)
	if err != nil {
		return fmt.Errorf("Failed to retrieve intentions: %v", err)
	}

	if len(intentions) > 0 {
		return fmt.Errorf("Intentions still exist: %v", intentions)
	}

	return nil
}

func testAccRemoveConsulIntention(t *testing.T) func() {
	return func() {
		client := getClient(testAccProvider.Meta())
		connect := client.Connect()
		qOpts := &consulapi.QueryOptions{}
		iM := &consulapi.IntentionMatch{
			By:    consulapi.IntentionMatchSource,
			Names: []string{"api"},
		}

		resp, _, err := connect.IntentionMatch(iM, qOpts)
		if err != nil {
			t.Fatalf("Failed to retrieve intentions by match err: %v", err)
		}

		intentions, hasMatch := resp["api"]
		if !hasMatch {
			t.Fatalf("No intention with source api was found")
		}

		var iid string
		for _, i := range intentions {
			if _, ok := i.Meta["is_tf_acc_test"]; ok {
				iid = i.ID
				break
			}
		}

		if iid == "" {
			t.Fatalf("Failed to find the intention created by Terraform")
		}

		_, err = connect.IntentionDelete(iid, &consulapi.WriteOptions{})
		if err != nil {
			t.Errorf("Failed to delete the intention. err: %s", err)
		}
	}
}

func TestAccConsulIntention_dc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulIntentionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulIntentionDc,
				ExpectError: regexp.MustCompile("No path to datacenter"),
			},
		},
	})
}

const testAccConsulIntentionConfigBasic = `
resource "consul_intention" "example" {
	source_name      = "api"
	destination_name = "db"
	action           = "allow"

	description = "something about example"
	meta = {
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

const testAccConsulIntentionConfigNamespace = `
resource "consul_intention" "example" {
	source_name           = "api"
	source_namespace      = "ns"
	destination_name      = "db"
	destination_namespace = "ns"

	action = "allow"
}
`

const testAccConsulIntentionDc = `
resource "consul_intention" "example" {
	datacenter       = "ny3"
	source_name      = "api"
	destination_name = "db"
	action           = "allow"

	meta = {
		is_tf_acc_test = "yes"
	}
}
`
