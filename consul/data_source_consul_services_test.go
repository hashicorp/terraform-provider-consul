package consul

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataConsulCatalogServices_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataConsulCatalogServicesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceValue("data.consul_services.read", "datacenter", "dc1"),
					testAccCheckDataSourceValue("data.consul_services.read", "services.%", "1"),
					testAccCheckDataSourceValue("data.consul_services.read", "services.consul", ""),
				),
			},
		},
	})
}

func TestAccDataConsulCatalogServices_alias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataConsulCatalogServicesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.consul_services.read", "services.%", "1"),
				),
			},
		},
	})
}

const testAccDataConsulCatalogServicesConfig = `
data "consul_services" "read" {
  query_options {
    allow_stale = true
    require_consistent = false
    token = ""
    wait_index = 0
    wait_time = "1m"
  }
}
`

const testAccDataConsulCatalogServicesAlias = `
data "consul_catalog_services" "read" {}
`
