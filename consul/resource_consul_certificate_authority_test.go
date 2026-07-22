// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulCertificateAuthority(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
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
				Config:      testAccConsulCertificateAuthorityConfigBoth,
				ExpectError: regexp.MustCompile(`"config": conflicts with config_json`),
			},
			{
				Config: testAccConsulCertificateAuthorityConfigJSON,
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

const testAccConsulCertificateAuthorityConfigBoth = `
resource "consul_certificate_authority" "test" {
	connect_provider = "consul"

	config = {
		LeafCertTTL         = "72h"
		RotationPeriod      = "1234h"
		IntermediateCertTTL = "5678h"
	}

	config_json = jsonencode({})
}
`

const testAccConsulCertificateAuthorityConfigJSON = `
resource "consul_certificate_authority" "test" {
	connect_provider = "consul"

	config_json = jsonencode({
		address = "http://localhost:8200"
		auth_method = {
			type       = "approle"
			mount_path = "approle"
			params = {
				role_id   = "role_id"
				secret_id = "secret_id"
			}
		}
		namespace             = "namespace"
		root_pki_path         = "root_pki_path"
		intermediate_pki_path = "intermediate_pki_path"
	})
}
`
