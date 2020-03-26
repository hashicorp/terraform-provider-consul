package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulACLAuthMethod() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulACLAuthMethodRead,

		Schema: map[string]*schema.Schema{
			// Filters
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Out parameters
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"config": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceConsulACLAuthMethodRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	name := d.Get("name").(string)
	dc, err := getDC(d, client, meta)
	if err != nil {
		return fmt.Errorf("Failed to get DC: %v", err)
	}
	qOpts := &consulapi.QueryOptions{
		Datacenter: dc,
		Namespace:  getNamespace(d, meta),
	}

	authMethod, _, err := client.ACL().AuthMethodRead(name, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to get auth-method: %v", err)
	}
	if authMethod == nil {
		return fmt.Errorf("Could not find auth-method '%s'", name)
	}

	d.SetId(authMethod.Name)
	if err = d.Set("type", authMethod.Type); err != nil {
		return fmt.Errorf("Failed to set 'type': %v", err)
	}
	if err = d.Set("description", authMethod.Description); err != nil {
		return fmt.Errorf("Failed to set 'description': %v", err)
	}
	if err = d.Set("config", authMethod.Config); err != nil {
		return fmt.Errorf("Failed to set 'config': %v", err)
	}
	return nil
}
