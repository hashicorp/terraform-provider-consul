// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccConsulCatalogEntry_basic(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    providers,
		CheckDestroy: testAccCheckConsulCatalogEntryDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccConsulCatalogEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulCatalogEntryExists(client),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "address", "127.0.0.1"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "node", "bastion"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.#", "1"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.address", "www.google.com"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.id", "google1"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.name", "google"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.port", "80"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.tags.#", "2"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.tags.2154398732", "tag0"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.tags.4151227546", "tag1"),
				),
			},
		},
	})
}

func TestAccConsulCatalogEntry_extremove(t *testing.T) {
	providers, client := startTestServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() {},
		Providers:    providers,
		CheckDestroy: testAccCheckConsulCatalogEntryDestroy(client),
		Steps: []resource.TestStep{
			{
				Config:             testAccConsulCatalogEntryConfig,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulCatalogEntryExists(client),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "address", "127.0.0.1"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "node", "bastion"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.#", "1"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.address", "www.google.com"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.id", "google1"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.name", "google"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.port", "80"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.tags.#", "2"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.tags.2154398732", "tag0"),
					testAccCheckConsulCatalogEntryValue("consul_catalog_entry.app", "service.3112399829.tags.4151227546", "tag1"),
					testAccCheckConsulCatalogEntryDeregister(client, "bastion"),
				),
			},
		},
	})
}

func testAccCheckConsulCatalogEntryDestroy(client *consulapi.Client) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		catalog := client.Catalog()
		qOpts := consulapi.QueryOptions{}
		services, _, err := catalog.Services(&qOpts)
		if err != nil {
			return fmt.Errorf("Could not retrieve services: %#v", err)
		}
		_, ok := services["google"]
		if ok {
			return fmt.Errorf("Service still exists: %#v", "google")
		}
		return nil
	}
}

func testAccCheckConsulCatalogEntryDeregister(client *consulapi.Client, node string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		catalog := client.Catalog()
		wOpts := consulapi.WriteOptions{}

		deregistration := consulapi.CatalogDeregistration{
			Node: node,
		}
		_, err := catalog.Deregister(&deregistration, &wOpts)
		if err != nil {
			return err
		}

		qOpts := consulapi.QueryOptions{}
		services, _, err := catalog.Services(&qOpts)
		if err != nil {
			return fmt.Errorf("Could not retrieve services: %#v", err)
		}
		_, ok := services["google"]
		if ok {
			return fmt.Errorf("Service still exists: %#v", "google")
		}
		return nil
	}
}

func testAccCheckConsulCatalogEntryExists(client *consulapi.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		catalog := client.Catalog()
		qOpts := consulapi.QueryOptions{}
		services, _, err := catalog.Services(&qOpts)
		if err != nil {
			return err
		}
		_, ok := services["google"]
		if !ok {
			return fmt.Errorf("Service does not exist: %#v", "google")
		}
		return nil
	}
}

func testAccCheckConsulCatalogEntryValue(n, attr, val string) resource.TestCheckFunc {
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

const testAccConsulCatalogEntryConfig = `
resource "consul_catalog_entry" "app" {
	address = "127.0.0.1"
	node = "bastion"
	service {
		address = "www.google.com"
		id = "google1"
		name = "google"
		port = 80
		tags = ["tag0", "tag1"]
	}
}
`
