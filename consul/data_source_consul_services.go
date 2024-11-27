// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"datacenter": {
				// Used in the query, must be stored and force a refresh if the value
				// changes.
				Computed: true,
				Type:     schema.TypeString,
			},
			"filter": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"query_options": queryOpts,

			// Out parameters
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"services": {
				Computed: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
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
	client, qOpts, _ := getClient(d, meta)

	// Parse out data source filters to populate Consul's query options
	getQueryOpts(qOpts, d, meta)

	services, meta, err := client.Catalog().Services(qOpts)
	if err != nil {
		return err
	}

	catalogServices := make(map[string]interface{}, len(services))
	for name, tags := range services {
		tagList := make([]string, 0, len(tags))
		tagList = append(tagList, tags...)
		sort.Strings(tagList)
		catalogServices[name] = strings.Join(tagList, " ")
	}

	serviceNames := make([]interface{}, 0, len(services))
	for k := range catalogServices {
		serviceNames = append(serviceNames, k)
	}

	const idKeyFmt = "catalog-services-%s"
	d.SetId(fmt.Sprintf(idKeyFmt, qOpts.Datacenter))

	d.Set("datacenter", qOpts.Datacenter)
	if err := d.Set("services", catalogServices); err != nil {
		return errwrap.Wrapf("Unable to store services: {{err}}", err)
	}

	catalogTags := map[string]map[string]struct{}{}
	for serviceName, tags := range services {
		for _, tag := range tags {
			if _, found := catalogTags[tag]; !found {
				catalogTags[tag] = map[string]struct{}{}
			}
			catalogTags[tag][serviceName] = struct{}{}
		}
	}
	ct := map[string]string{}
	for tag, services := range catalogTags {
		serviceList := []string{}
		for s := range services {
			serviceList = append(serviceList, s)
		}
		sort.Strings(serviceList)
		ct[tag] = strings.Join(serviceList, " ")
	}
	if err := d.Set("tags", ct); err != nil {
		return errwrap.Wrapf("Unable to store tags: {{err}}", err)
	}

	if err := d.Set("names", serviceNames); err != nil {
		return errwrap.Wrapf("Unable to store service names: {{err}}", err)
	}

	return nil
}
