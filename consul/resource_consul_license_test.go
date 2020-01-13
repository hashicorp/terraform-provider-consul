package consul

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulLicense_FailOnCommunityEdition(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulLicense,
				// Setting the Consul license will fail on the Community Edition
				ExpectError: regexp.MustCompile("failed to set license: Unexpected response code: 404"),
			},
		},
	})
}

func TestAccConsulLicense_BadLicense(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulLicense,
				ExpectError: regexp.MustCompile(`failed to set license: Unexpected response code: 400 \(Bad request: unknown version: .*\)`),
			},
		},
	})
}

func TestAccConsulLicense_CorrectLicense(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipTestOnConsulCommunityEdition(t)
			if _, err := os.Stat("../test_license.hclic"); os.IsNotExist(err) {
				t.Skip("This test needs a valid 'test_license.hclic' file to run.")
			}
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulLicense_CorrectLicense,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_license.license", "valid", "true"),
					resource.TestCheckResourceAttr("consul_license.license", "product", "consul"),
					resource.TestCheckResourceAttr("consul_license.license", "warnings.#", "0"),
				),
			},
		},
	})
}

func skipTestOnConsulCommunityEdition(t *testing.T) {
	testAccPreCheck(t)

	client := getClient(testAccProvider.Meta())
	self, err := client.Agent().Self()
	if err != nil {
		t.Fatalf("failed to get agent information: %v", err)
	}
	if !strings.HasSuffix(self["Config"]["Version"].(string), "+ent") {
		t.Skip("Test skipped on Consul Community Edition. Use a Consul Enterprise server to run this test.")
	}
}

func skipTestOnConsulEnterpriseEdition(t *testing.T) {
	testAccPreCheck(t)

	client := getClient(testAccProvider.Meta())
	self, err := client.Agent().Self()
	if err != nil {
		t.Fatalf("failed to get agent information: %v", err)
	}
	if strings.HasSuffix(self["Config"]["Version"].(string), "+ent") {
		t.Skip("Test skipped on Consul Enterprise Edition. Use a Consul Community server to run this test.")
	}
}

const testAccConsulLicense = `
resource "consul_license" "license" {
	license = "foobar"
}
`

const testAccConsulLicense_CorrectLicense = `
resource "consul_license" "license" {
	license = file("../test_license.hclic")
}
`
