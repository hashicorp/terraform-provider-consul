// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulCatalogEntry() *schema.Resource {
	return &schema.Resource{
		Create:             resourceConsulCatalogEntryCreate,
		Update:             resourceConsulCatalogEntryCreate,
		Read:               resourceConsulCatalogEntryRead,
		Delete:             resourceConsulCatalogEntryDelete,
		DeprecationMessage: "The consul_catalog_entry resource will be deprecated and removed in a future version. More information: https://github.com/hashicorp/terraform-provider-consul/issues/46",

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"node": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"service": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},

						"tags": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      resourceConsulCatalogEntryServiceTagsHash,
						},
					},
				},
				Set: resourceConsulCatalogEntryServicesHash,
			},

			"token": {
				Type:       schema.TypeString,
				Optional:   true,
				Sensitive:  true,
				Deprecated: tokenDeprecationMessage,
			},
		},
	}
}

func resourceConsulCatalogEntryServiceTagsHash(v interface{}) int {
	return hashcode.String(v.(string))
}

func resourceConsulCatalogEntryServicesHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["id"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["address"].(string)))
	buf.WriteString(fmt.Sprintf("%d-", m["port"].(int)))
	if v, ok := m["tags"]; ok {
		vs := v.(*schema.Set).List()
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		sort.Strings(s)

		for _, v := range s {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}
	return hashcode.String(buf.String())
}

func resourceConsulCatalogEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	catalog := client.Catalog()

	address := d.Get("address").(string)
	node := d.Get("node").(string)

	var serviceIDs []string
	if service, ok := d.GetOk("service"); ok {
		serviceList := service.(*schema.Set).List()
		serviceIDs = make([]string, len(serviceList))
		for i, rawService := range serviceList {
			serviceData := rawService.(map[string]interface{})

			if len(serviceData["id"].(string)) == 0 {
				serviceData["id"] = serviceData["name"].(string)
			}
			serviceID := serviceData["id"].(string)
			serviceIDs[i] = serviceID

			var tags []string
			if v := serviceData["tags"].(*schema.Set).List(); len(v) > 0 {
				tags = make([]string, len(v))
				for i, raw := range v {
					tags[i] = raw.(string)
				}
			}

			registration := &consulapi.CatalogRegistration{
				Address:    address,
				Datacenter: wOpts.Datacenter,
				Node:       node,
				Service: &consulapi.AgentService{
					Address: serviceData["address"].(string),
					ID:      serviceID,
					Service: serviceData["name"].(string),
					Port:    serviceData["port"].(int),
					Tags:    tags,
				},
			}

			if _, err := catalog.Register(registration, wOpts); err != nil {
				return fmt.Errorf("failed to register Consul catalog entry with node '%s' at address '%s' in %s: %v",
					node, address, wOpts.Datacenter, err)
			}
		}
	} else {
		registration := &consulapi.CatalogRegistration{
			Address:    address,
			Datacenter: wOpts.Datacenter,
			Node:       node,
		}

		if _, err := catalog.Register(registration, wOpts); err != nil {
			return fmt.Errorf("failed to register Consul catalog entry with node '%s' at address '%s' in %s: %v",
				node, address, wOpts.Datacenter, err)
		}
	}

	// Update the resource
	if _, _, err := catalog.Node(node, qOpts); err != nil {
		return fmt.Errorf("failed to read Consul catalog entry for node '%s' at address '%s' in %s: %v",
			node, address, qOpts.Datacenter, err)
	} else {
		d.Set("datacenter", qOpts.Datacenter)
	}

	sort.Strings(serviceIDs)
	serviceIDsJoined := strings.Join(serviceIDs, ",")

	d.SetId(fmt.Sprintf("%s-%s-[%s]", node, address, serviceIDsJoined))

	return nil
}

func resourceConsulCatalogEntryRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	catalog := client.Catalog()

	node := d.Get("node").(string)

	cNode, _, err := catalog.Node(node, qOpts)
	if err != nil {
		return fmt.Errorf("failed to get node '%s' from Consul catalog: %v", node, err)
	}
	if cNode == nil || cNode.Node == nil {
		d.SetId("")
	}

	return nil
}

func resourceConsulCatalogEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	catalog := client.Catalog()

	address := d.Get("address").(string)
	node := d.Get("node").(string)

	deregistration := consulapi.CatalogDeregistration{
		Address:    address,
		Datacenter: wOpts.Datacenter,
		Node:       node,
	}

	if _, err := catalog.Deregister(&deregistration, wOpts); err != nil {
		return fmt.Errorf("failed to deregister Consul catalog entry with node '%s' at address '%s' in %s: %v",
			node, address, wOpts.Datacenter, err)
	}

	// Clear the ID
	d.SetId("")
	return nil
}
