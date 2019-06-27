package consul

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccConsulConfigurationEntry_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConsulConfigurationEntry_ServiceDefaults,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_configuration_entry.foo", "name", "foo"),
					resource.TestCheckResourceAttr("consul_configuration_entry.foo", "kind", "service-defaults"),
					resource.TestCheckResourceAttr("consul_configuration_entry.foo", "protocol", "https"),
				),
			},
			{
				Config: testAccConsulConfigurationEntry_ProxyDefaults,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("consul_configuration_entry.foo", "name", "global"),
					resource.TestCheckResourceAttr("consul_configuration_entry.foo", "kind", "proxy-defaults"),
					resource.TestCheckResourceAttr("consul_configuration_entry.foo", "config.foo", "bar"),
				),
			},
		},
	})
}

func TestAccConsulConfigurationEntry_Errors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsulConfigurationEntry_ProxyDefaultsWrongName,
				ExpectError: regexp.MustCompile("'name' must be 'global' when 'kind' is 'proxy-defaults'"),
			},
			{
				Config:      testAccConsulConfigurationEntry_ProxyDefaultsProtocolSet,
				ExpectError: regexp.MustCompile("'protocol' must not be set when 'kind' is 'proxy-defaults'"),
			},
			{
				Config:      testAccConsulConfigurationEntry_ServiceDefaultsConfigSet,
				ExpectError: regexp.MustCompile("'config' must not be set when 'kind' is 'service-defaults'"),
			},
		},
	})
}

const testAccConsulConfigurationEntry_ServiceDefaults = `
resource "consul_configuration_entry" "foo" {
	name = "foo"
	kind = "service-defaults"

	protocol = "https"
}
`

const testAccConsulConfigurationEntry_ProxyDefaults = `
resource "consul_configuration_entry" "foo" {
	name = "global"
	kind = "proxy-defaults"

	config = {
		foo = "bar"
	}
}
`

const testAccConsulConfigurationEntry_ProxyDefaultsWrongName = `
resource "consul_configuration_entry" "foo" {
	name = "foo"
	kind = "proxy-defaults"

	config = {
		foo = "bar"
	}
}
`

const testAccConsulConfigurationEntry_ProxyDefaultsProtocolSet = `
resource "consul_configuration_entry" "foo" {
	name = "global"
	kind = "proxy-defaults"

	protocol = "https"
	config = {
		foo = "bar"
	}
}
`

const testAccConsulConfigurationEntry_ServiceDefaultsConfigSet = `
resource "consul_configuration_entry" "foo" {
	name = "foo"
	kind = "service-defaults"

	protocol = "https"
	config = {
		foo = "bar"
	}
}
`
