// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccConsulServiceSplitterConfigEEEntryTest(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testConsulServiceSplitterConfigEntryEE,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testConsulServiceSplitterConfigEntryEE = `
resource "consul_service_splitter_config_entry" "service-splitter-config-entry" {
	name      = "service-splitter" 
	meta      = {
		key = "value"
	}
	namespace = "namespace"
	partition = "partition"
	splits {
		weight         = 90                   
		service        = "frontend"            
		service_subset  = "v1"                
	}
	splits {
		weight         = 10
		service        = "frontend"
		service_subset  = "v2"
	}
}
`
