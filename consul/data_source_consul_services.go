package consul

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// Datasource predicates
	catalogServicesServiceName = "name"

	// Out parameters
	catalogServicesDatacenter  = "datacenter"
	catalogServicesNames       = "names"
	catalogServicesServices    = "services"
	catalogServicesServiceTags = "tags"
)

func dataSourceConsulServices() *schema.Resource {
	queryOpts := schemaQueryOpts()
	queryOpts.Elem.(*schema.Resource).Schema["namespace"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	return &schema.Resource{
		Read: dataSourceConsulServicesRead,
		Schema: map[string]*schema.Schema{
			// Data Source Predicate(s)
			catalogServicesDatacenter: {
				// Used in the query, must be stored and force a refresh if the value
				// changes.
				Computed: true,
				Type:     schema.TypeString,
			},
			catalogNodesQueryOpts: queryOpts,

			// Out parameters
			catalogServicesNames: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			catalogServicesServices: {
				Computed: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceConsulServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)

	// Parse out data source filters to populate Consul's query options
	queryOpts, err := getQueryOpts(d, client, meta)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog services: {{err}}", err)
	}

	services, meta, err := client.Catalog().Services(queryOpts)
	if err != nil {
		return err
	}

	catalogServices := make(map[string]interface{}, len(services))
	for name, tags := range services {
		tagList := make([]string, 0, len(tags))
		for _, tag := range tags {
			tagList = append(tagList, tag)
		}

		sort.Strings(tagList)
		catalogServices[name] = strings.Join(tagList, " ")
	}

	serviceNames := make([]interface{}, 0, len(services))
	for k := range catalogServices {
		serviceNames = append(serviceNames, k)
	}

	const idKeyFmt = "catalog-services-%s"
	d.SetId(fmt.Sprintf(idKeyFmt, queryOpts.Datacenter))

	d.Set(catalogServicesDatacenter, queryOpts.Datacenter)
	if err := d.Set(catalogServicesServices, catalogServices); err != nil {
		return errwrap.Wrapf("Unable to store services: {{err}}", err)
	}

	if err := d.Set(catalogServicesNames, serviceNames); err != nil {
		return errwrap.Wrapf("Unable to store service names: {{err}}", err)
	}

	return nil
}
