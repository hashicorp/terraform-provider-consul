package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
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

	var policyEntry *consulapi.ACLPolicyListEntry
	policyEntries, _, err := client.ACL().PolicyList(qOpts)
	if err != nil {
		return fmt.Errorf("could not list policies: %v", err)
	}
	for _, pe := range policyEntries {
		if pe.Name == name {
			policyEntry = pe
			break
		}
	}
	if policyEntry == nil {
		return fmt.Errorf("could not find policy '%s'", name)
	}

	policy, _, err := client.ACL().PolicyRead(policyEntry.ID, qOpts)
	if err != nil {
		return fmt.Errorf("could not read policy '%s': %v", name, err)
	}

	d.SetId(policy.ID)

	sw := newStateWriter(d)
	sw.set("description", policy.Description)
	sw.set("rules", policy.Rules)
	sw.set("datacenters", policy.Datacenters)

	return sw.error()
}
