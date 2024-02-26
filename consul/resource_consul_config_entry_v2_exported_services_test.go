// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccConsulExportedServicesV2_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulExportedServicesV2Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.test", "name", "test"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.test", "kind", "ExportedServices"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.test", "partition", "default"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.test", "services.#", "2"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.test", "services.0", "s1"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.test", "services.1", "s2"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.test", "sameness_group_consumers.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.test", "sameness_group_consumers.0", "sg1"),
				),
			},
			{
				Config: testAccConsulNamespaceExportedServicesV2Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.nstest", "name", "nstest"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.nstest", "kind", "NamespaceExportedServices"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.nstest", "partition", "default"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.nstest", "peer_consumers.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.nstest", "peer_consumers.0", "p1"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.nstest", "partition_consumers.#", "2"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.nstest", "partition_consumers.0", "ap1"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.nstest", "partition_consumers.1", "ap2"),
				),
			},
			{
				Config: testAccConsulPartitionExportedServicesV2Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.ptest", "name", "ptest"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.ptest", "kind", "PartitionExportedServices"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.ptest", "partition", "default"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.ptest", "peer_consumers.#", "1"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.ptest", "peer_consumers.0", "p1"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.ptest", "partition_consumers.#", "2"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.ptest", "partition_consumers.0", "ap1"),
					resource.TestCheckResourceAttr("consul_config_entry_v2_exported_services.ptest", "partition_consumers.1", "ap2"),
				),
			},
		},
	})
}

const testAccConsulExportedServicesV2Basic = `
resource "consul_config_entry_v2_exported_services" "test" {
	name = "test"
    kind = "ExportedServices"
    namespace = "default"
    partition = "default"
    services = ["s1", "s2"]
    sameness_group_consumers = ["sg1"]
}`

const testAccConsulNamespaceExportedServicesV2Basic = `
resource "consul_config_entry_v2_exported_services" "nstest" {
	name = "nstest"
    kind = "NamespaceExportedServices"
    namespace = "default"
    partition = "default"
    peer_consumers = ["p1"]
    partition_consumers = ["ap1", "ap2"]
}`

const testAccConsulPartitionExportedServicesV2Basic = `
resource "consul_config_entry_v2_exported_services" "ptest" {
	name = "ptest"
    kind = "PartitionExportedServices"
    partition = "default"
    peer_consumers = ["p1"]
    partition_consumers = ["ap1", "ap2"]
}`
