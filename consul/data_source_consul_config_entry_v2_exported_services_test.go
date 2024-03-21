package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataExportedServicesV2_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceExportedServicesV2ConfigNotFound,
				SkipFunc:    skipIfConsulVersionLT(client, "1.18.0"),
				ExpectError: regexp.MustCompile(`exported services config not found: not-found`),
			},
			{
				Config:   testAccDataSourceExportedServicesV2ConfigBasic,
				SkipFunc: skipIfConsulVersionLT(client, "1.18.0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "name", "test"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "kind", "ExportedServices"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "namespace", "default"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "partition", "default"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "services.0", "s1"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "sameness_group_consumers.0", "sg1"),
				),
			},
			{
				Config:   testAccDataSourceNamespaceExportedServicesV2ConfigBasic,
				SkipFunc: skipIfConsulVersionLT(client, "1.18.0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "name", "test"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "kind", "NamespaceExportedServices"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "partition", "default"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "partition_consumers.0", "default"),
				),
			},
			{
				Config:   testAccDataSourcePartitionExportedServicesV2ConfigBasic,
				SkipFunc: skipIfConsulVersionLT(client, "1.18.0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "name", "test"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "kind", "PartitionExportedServices"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "partition", "default"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "peer_consumers.0", "peer1"),
				),
			},
		},
	})
}

const testAccDataSourcePartitionExportedServicesV2ConfigBasic = `
resource "consul_config_entry_v2_exported_services" "test" {
	name = "test"
    kind = "PartitionExportedServices"
    partition = "default"
    peer_consumers = ["peer1"]
}

data "consul_config_entry_v2_exported_services" "read" {
	name = consul_config_entry_v2_exported_services.test.name
    kind = consul_config_entry_v2_exported_services.test.kind
    namespace = consul_config_entry_v2_exported_services.test.namespace
    partition = consul_config_entry_v2_exported_services.test.partition
}
`

const testAccDataSourceNamespaceExportedServicesV2ConfigBasic = `
resource "consul_config_entry_v2_exported_services" "test" {
	name = "test"
    kind = "NamespaceExportedServices"
    namespace = "default"
    partition = "default"
    partition_consumers = ["default"]
}

data "consul_config_entry_v2_exported_services" "read" {
	name = consul_config_entry_v2_exported_services.test.name
    kind = consul_config_entry_v2_exported_services.test.kind
    namespace = consul_config_entry_v2_exported_services.test.namespace
    partition = consul_config_entry_v2_exported_services.test.partition
}
`

const testAccDataSourceExportedServicesV2ConfigBasic = `
resource "consul_config_entry_v2_exported_services" "test" {
	name = "test"
    kind = "ExportedServices"
    namespace = "default"
    partition = "default"
    services = ["s1"]
    sameness_group_consumers = ["sg1"]
}

data "consul_config_entry_v2_exported_services" "read" {
	name = consul_config_entry_v2_exported_services.test.name
    kind = consul_config_entry_v2_exported_services.test.kind
    namespace = consul_config_entry_v2_exported_services.test.namespace
    partition = consul_config_entry_v2_exported_services.test.partition
}
`

const testAccDataSourceExportedServicesV2ConfigNotFound = `
data "consul_config_entry_v2_exported_services" "test" {
	name = "not-found"
    kind = "ExportedServices"
    namespace = "default"
    partition = "default"
}
`
