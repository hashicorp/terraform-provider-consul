// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataConsulKeysNonExistentKeys(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsulKeysNonExistantKeyConfig,
				ExpectError: regexp.MustCompile("Key '.*' does not exist"),
			},
			{
				Config: testAccDataConsulKeysNonExistantKeyDefaultBehaviourConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.read", "read", ""),
				),
			},
		},
	})
}

func TestAccDataConsulKeysNonExistentKeyWithDefault(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeysNonExistantKeyWithDefaultConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.read", "read", "myvalue"),
				),
			},
		},
	})
}

func TestAccDataConsulKeysExistentKeyWithEmptyValueAndDefault(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeysExistantKeyWithDefaultAndEmptyValueConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.read", "read", "myvalue"),
				),
			},
		},
	})
}

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
				Config: testAccDataConsulKeysConfigNamespaceCE,

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

const testAccDataConsulKeysNonExistantKeyDefaultBehaviourConfig = `

data "consul_keys" "read" {
    datacenter = "dc1"
    key {
        path = "test/set"
        name = "read"
		
    }
}
`
const testAccDataConsulKeysNonExistantKeyConfig = `
provider "consul" {
    error_on_missing_key = true
}
data "consul_keys" "read" {
    datacenter = "dc1"
    key {
        path = "test/set"
        name = "read"
		
    }
}
`

const testAccDataConsulKeysNonExistantKeyWithDefaultConfig = `

data "consul_keys" "read" {
    # Create a dependency on the resource so we're sure to
    # have the value in place before we try to read it.
    datacenter = "dc1"
    key {
        path = "test/set"
        name = "read"
		default = "myvalue"
    }
}

`

const testAccDataConsulKeysExistantKeyWithDefaultAndEmptyValueConfig = `

resource "consul_keys" "write" {
    datacenter = "dc1"

    key {
        path = "test/set"
        value = ""
		delete = true
    }
}
data "consul_keys" "read" {
    # Create a dependency on the resource so we're sure to
    # have the value in place before we try to read it.
    datacenter = "dc1"
    key {
        path = "test/set"
        name = "read"
		default = "myvalue"
    }
}
`

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

const testAccDataConsulKeysConfigDatacenter = `
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
