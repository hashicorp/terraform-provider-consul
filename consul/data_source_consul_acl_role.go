package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulACLRole() *schema.Resource {
	return &schema.Resource{
		Read: datasourceConsulACLRoleRead,

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

			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},

			"service_identities": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"datacenters": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func datasourceConsulACLRoleRead(d *schema.ResourceData, meta interface{}) error {
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

	role, _, err := client.ACL().RoleReadByName(name, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to get role: %v", err)
	}
	if role == nil {
		return fmt.Errorf("Could not find role '%s'", name)
	}

	d.SetId(role.ID)
	if err = d.Set("description", role.Description); err != nil {
		return fmt.Errorf("Failed to set 'description': %v", err)
	}

	policies := make([]map[string]interface{}, len(role.Policies))
	for i, p := range role.Policies {
		policies[i] = map[string]interface{}{
			"name": p.Name,
			"id":   p.ID,
		}
	}
	if err = d.Set("policies", policies); err != nil {
		return fmt.Errorf("Failed to set 'policies': %v", err)
	}

	identities := make([]map[string]interface{}, len(role.ServiceIdentities))
	for i, si := range role.ServiceIdentities {
		identities[i] = map[string]interface{}{
			"service_name": si.ServiceName,
			"datacenters":  si.Datacenters,
		}
	}
	if err = d.Set("service_identities", identities); err != nil {
		return fmt.Errorf("Failed to set 'service_identities': %v", err)
	}

	return nil
}
