// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConsulAutopilotHealth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulAutopilotHealthRead,
		Schema: map[string]*schema.Schema{
			// Filters
			"datacenter": {
				Optional: true,
				Type:     schema.TypeString,
			},

			// Out parameters
			"healthy": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"failure_tolerance": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"servers": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"address": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"serf_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"version": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"leader": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"last_contact": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"last_term": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"last_index": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"healthy": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"voter": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"stable_since": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceConsulAutopilotHealthRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	operator := client.Operator()
	getQueryOpts(qOpts, d, meta)

	health, err := operator.AutopilotServerHealth(qOpts)
	if err != nil {
		return err
	}
	const idKeyFmt = "autopilot-health-%s"
	d.SetId(fmt.Sprintf(idKeyFmt, qOpts.Datacenter))

	d.Set("healthy", health.Healthy)
	d.Set("failure_tolerance", health.FailureTolerance)

	serversHealth := make([]interface{}, 0, len(health.Servers))
	for _, server := range health.Servers {
		h := make(map[string]interface{}, 12)

		h["id"] = server.ID
		h["name"] = server.Name
		h["address"] = server.Address
		h["serf_status"] = server.SerfStatus
		h["version"] = server.Version
		h["leader"] = server.Leader
		h["last_contact"] = server.LastContact.String()
		h["last_term"] = server.LastTerm
		h["last_index"] = server.LastIndex
		h["healthy"] = server.Healthy
		h["voter"] = server.Voter
		h["stable_since"] = server.StableSince.String()

		serversHealth = append(serversHealth, h)
	}

	if err := d.Set("servers", serversHealth); err != nil {
		return errwrap.Wrapf("Unable to store servers health: {{err}}", err)
	}
	return nil
}
