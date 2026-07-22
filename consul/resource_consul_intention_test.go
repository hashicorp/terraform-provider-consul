// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    providers,
		CheckDestroy: testAccCheckConsulIntentionDestroy(client),
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
				PreConfig: testAccRemoveConsulIntention(t, client),
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
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    providers,
		CheckDestroy: testAccCheckConsulIntentionDestroy(client),
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulIntentionConfigBadAction,
				ExpectError: regexp.MustCompile("expected action to be one of"),
			},
		},
	})
}

func TestAccConsulIntention_namespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulIntentionConfigNamespaceCE,
				ExpectError: regexp.MustCompile("Namespaces is a Consul Enterprise feature"),
			},
		},
	})
}

func TestAccConsulIntention_namespaceEE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConsulIntentionConfigNamespaceEE,
			},
		},
	})
}

func testAccCheckConsulIntentionDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
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
}

func testAccRemoveConsulIntention(t *testing.T, client *consulapi.Client) func() {
	return func() {
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

const testAccConsulIntentionConfigNamespaceCE = `
resource "consul_intention" "example" {
	source_name           = "api"
	source_namespace      = "ns"
	destination_name      = "db"
	destination_namespace = "ns"

	action = "allow"
}
`

const testAccConsulIntentionConfigNamespaceEE = `
resource "consul_namespace" "ns" {
	name = "ns"
}

resource "consul_intention" "example" {
	source_name           = "api"
	source_namespace      = consul_namespace.ns.name
	destination_name      = "db"
	destination_namespace = consul_namespace.ns.name

	action = "allow"
}
`
