package consul

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataConsulKeyPrefix_basic(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeyPrefixConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.read", "var.read1", "written1"),
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.read", "var.read2", "written2"),
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.read", "var.read3", "default3"),
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.read", "datacenter", "dc1"),
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.read", "path_prefix", "myapp/config/"),
					resource.TestCheckNoResourceAttr("data.consul_key_prefix.read", "subkeys.%"),
					resource.TestCheckNoResourceAttr("data.consul_key_prefix.read", "subkeys.key1"),
					resource.TestCheckNoResourceAttr("data.consul_key_prefix.read", "subkeys.key2/value"),
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.read2", "subkeys.%", "2"),
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.read2", "subkeys.key1", "written1"),
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.read2", "subkeys.key2/value", "written2"),
				),
			},
		},
	})
}

func TestAccDataConsulKeyPrefix_namespaceCE(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulEnterpriseEdition(t) },
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsulKeyPrefixConfigNamespaceCE,
				ExpectError: regexp.MustCompile("Unexpected response code: 400"),
			},
		},
	})
}

func TestAccDataConsulKeyPrefix_namespaceEE(t *testing.T) {
	startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeyPrefixConfigNamespaceEE,
			},
		},
	})
}

func TestAccDataConsulKeyPrefix_datacenter(t *testing.T) {
	startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeyPrefixConfigDatacenter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.dc1", "subkeys.%", "0"),
					testAccCheckConsulKeyPrefixAttribute("data.consul_key_prefix.dc2", "subkeys.%", "1"),
				),
			},
		},
	})
}

func testAccCheckConsulKeyPrefixAttribute(n, attr, val string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		out, ok := rn.Primary.Attributes[attr]
		if !ok {
			return fmt.Errorf("Attribute '%s' not found: %#v", attr, rn.Primary.Attributes)
		}
		if val != "<any>" && out != val {
			return fmt.Errorf("Attribute '%s' value '%s' != '%s'", attr, out, val)
		}
		if val == "<any>" && out == "" {
			return fmt.Errorf("Attribute '%s' value '%s'", attr, out)
		}
		return nil
	}
}

const testAccDataConsulKeyPrefixConfig = `
resource "consul_key_prefix" "write" {
    datacenter = "dc1"

    path_prefix = "myapp/config/"

    subkeys = {
        "key1" = "written1"
        "key2/value" = "written2"
    }
}

data "consul_key_prefix" "read" {
    # Create a dependency on the resource so we're sure to
    # have the value in place before we try to read it.
    datacenter = "${consul_key_prefix.write.datacenter}"

    path_prefix = "${consul_key_prefix.write.path_prefix}"

    subkey {
        path = "key1"
        name = "read1"
    }

    subkey {
        path = "key2/value"
        name = "read2"
    }

    subkey {
        path = "key3/foo/bar"
        name = "read3"
        default = "default3"
    }
}

data "consul_key_prefix" "read2" {
    # Create a dependency on the resource so we're sure to
    # have the value in place before we try to read it.
    datacenter = "${consul_key_prefix.write.datacenter}"

    path_prefix = "${consul_key_prefix.write.path_prefix}"
}
`

const testAccDataConsulKeyPrefixConfigNamespaceCE = `
data "consul_key_prefix" "read" {
  path_prefix = "foo/"
  namespace   = "test-key-prefix"
}
`

const testAccDataConsulKeyPrefixConfigNamespaceEE = `
resource "consul_key_prefix" "write" {
  path_prefix = "myapp/config/"

  subkeys = {
    "key1"       = "written1"
    "key2/value" = "written2"
  }
}

data "consul_key_prefix" "read" {
  path_prefix = "consul_key_prefix.write.path_prefix"
  namespace   = "test-key-prefix"
}
`

const testAccDataConsulKeyPrefixConfigDatacenter = `
resource "consul_key_prefix" "dc2" {
	datacenter  = "dc2"
	path_prefix = "test/"

	subkeys = {
		"dc"       = "dc2"
	}
}

data "consul_key_prefix" "dc1" {
	path_prefix = consul_key_prefix.dc2.path_prefix
}

data "consul_key_prefix" "dc2" {
	datacenter  = "dc2"
	path_prefix = consul_key_prefix.dc2.path_prefix
}
`
