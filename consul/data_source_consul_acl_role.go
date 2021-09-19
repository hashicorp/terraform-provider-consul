package consul

import (
	"fmt"

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
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Out parameters
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
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
				Computed: true,
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
			"node_identities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceConsulACLRoleRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Get("name").(string)

	role, _, err := client.ACL().RoleReadByName(name, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to get role: %v", err)
	}
	if role == nil {
		return fmt.Errorf("Could not find role '%s'", name)
	}

	policies := make([]map[string]interface{}, len(role.Policies))
	for i, p := range role.Policies {
		policies[i] = map[string]interface{}{
			"name": p.Name,
			"id":   p.ID,
		}
	}

	identities := make([]map[string]interface{}, len(role.ServiceIdentities))
	for i, si := range role.ServiceIdentities {
		identities[i] = map[string]interface{}{
			"service_name": si.ServiceName,
			"datacenters":  si.Datacenters,
		}
	}

	nodeIdentities := make([]interface{}, len(role.NodeIdentities))
	for i, ni := range role.NodeIdentities {
		nodeIdentities[i] = map[string]interface{}{
			"node_name":  ni.NodeName,
			"datacenter": ni.Datacenter,
		}
	}

	d.SetId(role.ID)

	sw := newStateWriter(d)
	sw.set("description", role.Description)
	sw.set("policies", policies)
	sw.set("service_identities", identities)
	sw.set("node_identities", nodeIdentities)

	return sw.error()
}
