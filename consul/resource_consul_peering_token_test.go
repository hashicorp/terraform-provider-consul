// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulPeeringToken_basic(t *testing.T) {
	ctx := context.Background()
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		CheckDestroy: func(s *terraform.State) error {
			peer, _, _ := client.Peerings().Read(ctx, "hello-world", nil)
			if peer != nil {
				return fmt.Errorf("the peer has not been removed")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConsulPeeringTokenBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("consul_peering_token.basic", "id", "hello-world"),
					resource.TestCheckResourceAttr("consul_peering_token.basic", "peer_name", "hello-world"),
					resource.TestCheckResourceAttrSet("consul_peering_token.basic", "peering_token"),
				),
			},
			{
				PreConfig: func() {
					_, err := client.Peerings().Delete(ctx, "hello-world", nil)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccConsulPeeringTokenBasic,
				Check: func(s *terraform.State) error {
					peer, _, err := client.Peerings().Read(ctx, "hello-world", nil)
					if err != nil {
						return err
					}
					if peer == nil {
						return fmt.Errorf("the peer does not exist")
					}
					return nil
				},
			},
		},
	})
}

const testAccConsulPeeringTokenBasic = `
resource "consul_peering_token" "basic" {
  peer_name = "hello-world"
}
`
