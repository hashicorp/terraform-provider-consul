// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulACLPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulACLPolicyRead,

		Schema: map[string]*schema.Schema{
			// Filters
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"partition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Out parameters
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datacenters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceConsulACLPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Get("name").(string)

	policy, _, err := client.ACL().PolicyReadByName(name, qOpts)
	if err != nil {
		return fmt.Errorf("could not read policy '%s': %v", name, err)
	}
	if policy == nil {
		return fmt.Errorf("could not find policy %q", name)
	}

	d.SetId(policy.ID)

	sw := newStateWriter(d)
	sw.set("description", policy.Description)
	sw.set("rules", policy.Rules)
	sw.set("datacenters", policy.Datacenters)

	return sw.error()
}
