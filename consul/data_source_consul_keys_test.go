package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulKeys_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeysConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.read", "read", "written"),
				),
			},
		},
	})
}

func TestAccDataConsulKeys_namespaceCE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsulKeysConfigNamespaceCE,
				ExpectError: regexp.MustCompile("Unexpected response code: 400"),
			},
		},
	})
}

func TestAccDataConsulKeys_namespaceEE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeysConfigNamespaceEE,
			},
		},
	})
}

const testAccDataConsulKeysConfig = `
resource "consul_keys" "write" {
    datacenter = "dc1"

    key {
        path = "test/data_source"
        value = "written"
    }
}

data "consul_keys" "read" {
    # Create a dependency on the resource so we're sure to
    # have the value in place before we try to read it.
    datacenter = "${consul_keys.write.datacenter}"

    key {
        path = "test/data_source"
        name = "read"
    }
}
`

const testAccDataConsulKeysConfigNamespaceCE = `
data "consul_keys" "read" {
  namespace  = "test-data-consul-keys"

  key {
    path = "test/data_source"
    name = "read"
  }
}`

const testAccDataConsulKeysConfigNamespaceEE = `
resource "consul_keys" "write" {
  datacenter = "dc1"

  key {
    path = "test/data_source"
    value = "written"
  }
}

resource "consul_namespace" "test" {
  name = "test-data-consul-keys"
}

data "consul_keys" "read" {
  namespace = consul_namespace.test.name
  datacenter = consul_keys.write.datacenter

  key {
    path = "test/data_source"
    name = "read"
  }
}`
