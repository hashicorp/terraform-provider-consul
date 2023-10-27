// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// fadia you have added this test
// here we are testing for no exitant key and no default value so we are expecting an error.
func TestAccDataConsulKeysNonExistentKeysDefaultBehaviour(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsulKeysNonExistantKeyDefaultBehaviourConfig,
				//ExpectError: regexp.MustCompile("Key '.*' does not exist"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulKeysValue("data.consul_keys.read", "read", ""),
				),
			},
		},
	})
}
func TestAccDataConsulKeysNonExistentKeys(t *testing.T) {
	providers, _ := startTestServer(t)

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsulKeysNonExistantKeyConfig,
				ExpectError: regexp.MustCompile("Key '.*' does not exist"),
			},
		},
	})
}

// here they key doesn't exist but we have a default value so we are checking if we get the default value correctly.
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

//fadia end of what you have added

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
					//testAccCheckConsulKeysValue("data.consul_keys.dc1", "read", ""),
					// I removed the previous line since now we have a correct behaviour of launching an error when they key doesn't exist
					testAccCheckConsulKeysValue("data.consul_keys.dc2", "read", "dc2"),
				),
				ExpectError: regexp.MustCompile("Key '.*' does not exist"), // added here becuase test/set doesn't exist in dc1.
			},
		},
	})
}

// fadia you have added the following
// A non existent key with with no default value, error expected in this config.
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
    new_behaviour = true
}
data "consul_keys" "read" {
    datacenter = "dc1"
    key {
        path = "test/set"
        name = "read"
		
    }
}
`

// A non existent key with a default value, no error expected here.
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

// exitant key with empty value and default value
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

// end of what you have added
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
