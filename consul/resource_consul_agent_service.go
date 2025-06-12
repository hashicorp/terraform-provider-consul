// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConsulAgentService() *schema.Resource {
	return &schema.Resource{
		Create:             resourceConsulAgentServiceCreate,
		Update:             resourceConsulAgentServiceCreate,
		Read:               resourceConsulAgentServiceRead,
		Delete:             resourceConsulAgentServiceDelete,
		DeprecationMessage: "The consul_agent_service resource will be deprecated and removed in a future version. More information: https://github.com/hashicorp/terraform-provider-consul/issues/46",

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
		},
	}
}

func resourceConsulAgentServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, _ := getClient(d, meta)
	agent := client.Agent()

	name := d.Get("name").(string)
	registration := consulapi.AgentServiceRegistration{Name: name}

	if address, ok := d.GetOk("address"); ok {
		registration.Address = address.(string)
	}

	if port, ok := d.GetOk("port"); ok {
		registration.Port = port.(int)
	}

	if v, ok := d.GetOk("tags"); ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		registration.Tags = s
	}

	if err := agent.ServiceRegister(&registration); err != nil {
		return fmt.Errorf("failed to register service '%s' with Consul agent: %v", name, err)
	}

	// Update the resource
	if serviceMap, err := agent.Services(); err != nil {
		return fmt.Errorf("failed to read services from Consul agent: %v", err)
	} else if service, ok := serviceMap[name]; !ok {
		return fmt.Errorf("failed to read service '%s' from Consul agent: %v", name, err)
	} else {
		d.Set("address", service.Address)
		d.Set("id", service.ID)
		d.SetId(service.ID)
		d.Set("name", service.Service)
		d.Set("port", service.Port)
		d.Set("tags", service.Tags)
	}

	return nil
}

func resourceConsulAgentServiceRead(d *schema.ResourceData, meta interface{}) error {
	client, _, _ := getClient(d, meta)
	agent := client.Agent()

	name := d.Get("name").(string)

	if services, err := agent.Services(); err != nil {
		return fmt.Errorf("failed to get services from Consul agent: %v", err)
	} else if service, ok := services[name]; !ok {
		d.Set("id", "")
	} else {
		d.Set("address", service.Address)
		d.Set("id", service.ID)
		d.SetId(service.ID)
		d.Set("name", service.Service)
		d.Set("port", service.Port)
		d.Set("tags", service.Tags)
	}

	return nil
}

func resourceConsulAgentServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, _ := getClient(d, meta)
	catalog := client.Agent()

	id := d.Id()

	if err := catalog.ServiceDeregister(id); err != nil {
		return fmt.Errorf("failed to deregister service '%s' from Consul agent: %v", id, err)
	}

	// Clear the ID
	d.SetId("")
	return nil
}
