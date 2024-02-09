package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataExportedServicesV2_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceExportedServicesV2ConfigNotFound,
				ExpectError: regexp.MustCompile(`exported services config not found: not-found`),
			},
			{
				Config: testAccDataSourceExportedServicesV2ConfigBasic,
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
				Config: testAccDataSourceComputedExportedServicesV2ConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "name", "default"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "kind", "ComputedExportedServices"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "partition", "default"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "services.0", "s1"),
					resource.TestCheckResourceAttr("data.consul_config_entry_v2_exported_services.read", "sameness_group_consumers.0", "sg1"),
				),
			},
		},
	})
}

const testAccDataSourceComputedExportedServicesV2ConfigBasic = `
resource "consul_config_entry_v2_exported_services" "test" {
	name = "test"
    kind = "ExportedServices"
    namespace = "default"
    partition = "default"
    services = ["s1"]
    partition_consumers = ["default"]
}

data "consul_config_entry_v2_exported_services" "read" {
	name = "default"
    kind = "ComputedExportedServices"
    partition = "default"
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
