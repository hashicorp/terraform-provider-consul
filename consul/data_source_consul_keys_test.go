// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulKeys_basic(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeysConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.read", "read", "written"),
				),
			},
			{
				Config:      testAccDataConsulKeysNonExistantKeyConfig,
				ExpectError: regexp.MustCompile(`Key ".*" does not exist`),
			},
			{
				Config: testAccDataConsulKeysNonExistantKeyDefaultBehaviourConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.read", "read", ""),
				),
			},
			{
				Config: testAccDataConsulKeysNonExistantKeyWithDefaultConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.read", "read", "myvalue"),
				),
			},
			{
				Config: testAccDataConsulKeysExistantKeyWithDefaultAndEmptyValueConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.read", "read", "myvalue"),
				),
			},
		},
	})
}

func TestAccDataConsulKeys_namespaceCE(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
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
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { skipTestOnConsulCommunityEdition(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeysConfigNamespaceEE,
			},
		},
	})
}

func TestAccDataConsulKeys_datacenter(t *testing.T) {
	providers, _ := startRemoteDatacenterTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeysConfigDatacenter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.dc1", "read", ""),
					testAccCheckConsulKeysValue("data.consul_keys.dc2", "read", "dc2"),
				),
			},
		},
	})
}

const (
	testAccDataConsulKeysNonExistantKeyDefaultBehaviourConfig = `
data "consul_keys" "read" {
    key {
        path = "test/set"
        name = "read"
    }
}
`
	testAccDataConsulKeysNonExistantKeyConfig = `
data "consul_keys" "read" {
	error_on_missing_keys = true
    key {
        path = "test/set"
        name = "read"
    }
}
`

	testAccDataConsulKeysNonExistantKeyWithDefaultConfig = `
data "consul_keys" "read" {
    key {
        path = "test/set"
        name = "read"
		default = "myvalue"
    }
}

`

	testAccDataConsulKeysExistantKeyWithDefaultAndEmptyValueConfig = `
resource "consul_keys" "write" {
    datacenter = "dc1"

    key {
        path   = "test/set"
        value  = ""
		delete = true
    }
}

data "consul_keys" "read" {
    # Create a dependency on the resource so we're sure to
    # have the value in place before we try to read it.
    datacenter = consul_keys.write.datacenter

    key {
        path = "test/set"
        name = "read"
		default = "myvalue"
    }
}
`

	testAccDataConsulKeysConfig = `
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
    datacenter = consul_keys.write.datacenter

    key {
        path = "test/data_source"
        name = "read"
    }
}
`

	testAccDataConsulKeysConfigNamespaceCE = `
data "consul_keys" "read" {
  namespace  = "test-data-consul-keys"

  key {
    path = "test/data_source"
    name = "read"
  }
}`

	testAccDataConsulKeysConfigNamespaceEE = `
resource "consul_keys" "write" {
  datacenter = "dc1"

  key {
    path  = "test/data_source"
    value = "written"
  }
}

resource "consul_namespace" "test" {
  name = "test-data-consul-keys"
}

data "consul_keys" "read" {
  namespace  = consul_namespace.test.name
  datacenter = consul_keys.write.datacenter

  key {
    path = "test/data_source"
    name = "read"
  }
}`

	testAccDataConsulKeysConfigDatacenter = `
resource "consul_keys" "write" {
    datacenter = "dc2"

    key {
        path   = "test/dc"
        value  = "dc2"
        delete = true
    }
}

data "consul_keys" "dc1" {
    key {
        path = "test/dc"
        name = "read"
    }
}

data "consul_keys" "dc2" {
    datacenter = consul_keys.write.datacenter

    key {
        path = "test/dc"
        name = "read"
    }
}
`
)
