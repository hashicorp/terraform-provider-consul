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
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testConsulConfigEntryServiceIntentionsCE = `

	name = "service-intention-3"
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
		action     = "allow"
		name       = "nightly-cronjob"
		precedence = 9
		type       = "consul"
	}
`
