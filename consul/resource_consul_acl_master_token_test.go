package consul

import (
	"io/ioutil"
	"regexp"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/resource"
)

var resetRegexp = regexp.MustCompile(`Unexpected response code: 403 \(Permission denied: ACL bootstrap no longer allowed \(reset index: (\d+)\)\)`)

func resetACL(t *testing.T) {
	// Developmemt server has already ACL enabled so we must first disabled them
	// to do this test
	// https://learn.hashicorp.com/consul/advanced/day-1-operations/acl-guide#ensure-the-acl-system-is-configured-properly
	config := &consulapi.Config{}
	client, err := consulapi.NewClient(config)
	if err != nil {
		t.Fatalf("Error while creating Consul Client: %s", err)
	}
	// Get the reset index from the API error message
	_, _, err = client.ACL().Bootstrap()
	if err == nil {
		t.Fatalf("ACL should already be enabled.")
	}
	if !resetRegexp.MatchString(err.Error()) {
		t.Fatalf("Error while bootstraping ACL: %s", err)
	}
	// Write the reset index in <data-directory>/acl-bootstrap-reset
	resetIndex := resetRegexp.FindStringSubmatch(err.Error())[1]
	// Consul data-dir is one level higher than the tests
	err = ioutil.WriteFile("../acl-bootstrap-reset", []byte(resetIndex), 0644)
	if err != nil {
		t.Fatalf("Failed to write acl-bootstrap-reset: %s", err)
	}
	// Bootstraping should now work again
}

func TestAccConsulACLMasterToken_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t); resetACL(t) },
		Steps: []resource.TestStep{
			{
				Config: testResourceACLMasterTokenConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_acl_master_token.test", "description", "Bootstrap Token (Global Management)"),
					resource.TestCheckResourceAttr("consul_acl_master_token.test", "policies.#", "1"),
					resource.TestCheckResourceAttr("consul_acl_master_token.test", "local", "false"),
					resource.TestCheckResourceAttrSet("consul_acl_master_token.test", "token"),
				),
			},
			{
				// Destroy will leave the configuration as is
				Config:  testResourceACLMasterTokenConfigBasic,
				Destroy: true,
			},
			{
				// Bootstraping when the cluster already has ACL enabled should
				// raise an error
				Config:      testResourceACLMasterTokenConfigBasic,
				ExpectError: regexp.MustCompile("Permission denied: ACL bootstrap no longer allowed"),
			},
		},
	})
}

const testResourceACLMasterTokenConfigBasic = `
resource "consul_acl_master_token" "test" {}
`
