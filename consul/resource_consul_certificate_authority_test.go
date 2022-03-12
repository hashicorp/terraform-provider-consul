package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulCertificateAuthority(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulCertificateAuthorityConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_certificate_authority.test", "connect_provider", "consul"),
					resource.TestCheckResourceAttr("consul_certificate_authority.test", "config.%", "3"),
					resource.TestCheckResourceAttr("consul_certificate_authority.test", "config.LeafCertTTL", "72h"),
					resource.TestCheckResourceAttr("consul_certificate_authority.test", "config.RotationPeriod", "1234h"),
					resource.TestCheckResourceAttr("consul_certificate_authority.test", "config.IntermediateCertTTL", "5678h"),
				),
			},
			{
				Config:            testAccConsulCertificateAuthorityConfig,
				ResourceName:      "consul_certificate_authority.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccConsulCertificateAuthorityConfig = `
resource "consul_certificate_authority" "test" {
	connect_provider = "consul"

	config = {
		LeafCertTTL         = "72h"
		RotationPeriod      = "1234h"
		IntermediateCertTTL = "5678h"
	}
}
`
