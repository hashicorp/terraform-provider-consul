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

			// Out parameters
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datacenters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
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
		return fmt.Errorf("Could not list policies: %v", err)
	}
	for _, pe := range policyEntries {
		if pe.Name == name {
			policyEntry = pe
			break
		}
	}
	if policyEntry == nil {
		return fmt.Errorf("Could not find policy '%s'", name)
	}

	policy, _, err := client.ACL().PolicyRead(policyEntry.ID, qOpts)
	if err != nil {
		return fmt.Errorf("Could not read policy '%s': %v", name, err)
	}

	d.SetId(policy.ID)
	if err = d.Set("description", policy.Description); err != nil {
		return fmt.Errorf("Could not set 'description': %v", err)
	}
	if err = d.Set("rules", policy.Rules); err != nil {
		return fmt.Errorf("Could not set 'rules': %v", err)
	}
	if err = d.Set("datacenters", policy.Datacenters); err != nil {
		return fmt.Errorf("Could not set 'datacenters': %v", err)
	}

	return nil
}
