// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
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
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "namespace", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "partition", ""),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "meta.key", "value"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "sources.0.name", "frontend-webapp"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "sources.0.type", "consul"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "sources.0.action", "allow"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "sources.0.precedence", "9"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "sources.1.name", "nightly-cronjob"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "sources.1.type", "consul"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "sources.1.action", "deny"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "sources.1.precedence", "9"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "jwt.2394986741.providers.0.name", "okta"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "jwt.2394986741.providers.0.verify_claims.0.path.0", "/"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "jwt.2394986741.providers.0.verify_claims.0.path.1", "path1"),
					resource.TestCheckResourceAttr("consul_config_entry_service_intentions.foo", "jwt.2394986741.providers.0.verify_claims.1.path.0", "/path"),
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
		ClockSkewSeconds = 30
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
