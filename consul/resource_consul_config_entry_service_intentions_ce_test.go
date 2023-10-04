// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulConfigEntryServiceIntentionsCETest(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testConsulConfigEntryServiceIntentionsCE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "name", "service-intention"),
				),
			},
		},
	})
}

const testConsulConfigEntryServiceIntentionsCE = `
resource "consul_config_entry" "jwt_provider" {
	name = "okta"
	kind = "jwt-provider"

	config_json = jsonencode({
		Issuer = "test-issuer"
		JSONWebKeySet = {
			Remote = {
				URI = "https://127.0.0.1:9091"
				FetchAsynchronously = true
			}
		}
	})
}
resource "consul_config_entry_service_intentions" "foo" {
	name = "service-intention"
	meta = {
		key = "value"
	}
	jwt {
		providers {
			name = consul_config_entry.jwt_provider.name
			verify_claims {
				path = ["/", "path1"]
				value = ""
			}
			verify_claims {
				path = ["/path"]
				value = "value"
			}
		}
	}
	sources {
		action     = "allow"
		name       = "frontend-webapp"
		precedence = 9
		type       = "consul"
	}
	sources {
		name       = "nightly-cronjob"
		precedence = 9
		type       = "consul"
		action = "deny"
	}
}
`
