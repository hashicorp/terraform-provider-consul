package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulServices_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
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
			{
				Config: testAccDataConsulCatalogServicesConfig,
				Check:  resource.TestCheckResourceAttr("data.consul_services.read", "services.%", "1"),
			},
		},
	})
}

func TestAccDataConsulCatalogServices_badToken(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsulCatalogServicesBadTokenConfig,
				ExpectError: regexp.MustCompile(`Unexpected response code: 403 \(ACL not found\)`),
			},
		},
	})
}

func TestAccDataConsulServices_namespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsulServicesNamespaceCE,
				ExpectError: namespaceEnterpriseFeature,
			},
		},
	})
}

func TestAccDataConsulServices_namespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulServicesNamespaceEE,
			},
		},
	})
}

func TestAccDataConsulServices_datacenter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccRemoteDatacenterPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulCatalogServicesDatacenter,
				Check:  testAccCheckDataSourceValue("data.consul_services.read", "services.%", "2"),
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

const testAccDataConsulCatalogServicesBadTokenConfig = `
data "consul_services" "read" {
  query_options {
    token = "foobar"
  }
}
`

const testAccDataConsulServicesNamespaceCE = `
data "consul_services" "read" {
  query_options {
    namespace = "test-data-services"
  }
}
`

const testAccDataConsulServicesNamespaceEE = `
data "consul_services" "read" {
  query_options {
    namespace = "test-data-services"
  }
}
`

const testAccDataConsulCatalogServicesDatacenter = `
resource "consul_node" "test" {
	datacenter = "dc2"
	name       = "test"
	address    = "test.com"
}

resource "consul_service" "test" {
	datacenter = "dc2"
	name       = "test"
	node       = consul_node.test.name
	port       = 80
}

data "consul_services" "read" {
  	query_options {
	  	datacenter = consul_service.test.datacenter
  	}
}
`
