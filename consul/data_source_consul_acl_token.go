package consul

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulACLToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulACLTokenRead,

		Schema: map[string]*schema.Schema{

			// Filters
			"accessor_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Out parameters
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Required: true,
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
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "List of roles.",
			},
			"service_identities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of service identities that should be applied to the token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the service.",
						},
						"datacenters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the datacenters the effective policy is valid within.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"node_identities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of node identities that should be applied to the token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The list of node identities that should be applied to the token.",
						},
						"datacenter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the node's datacenter.",
						},
					},
				},
			},
			"local": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"expiration_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set this represents the point after which a token should be considered revoked and is eligible for destruction.",
			},
		},
	}
}

func dataSourceConsulACLTokenRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	accessorID := d.Get("accessor_id").(string)

	aclToken, _, err := client.ACL().TokenRead(accessorID, qOpts)
	if err != nil {
		return err
	}

	policies := make([]map[string]interface{}, len(aclToken.Policies))
	for i, policyLink := range aclToken.Policies {
		policies[i] = map[string]interface{}{
			"name": policyLink.Name,
			"id":   policyLink.ID,
		}
	}

	roles := make([]interface{}, len(aclToken.Roles))
	for i, r := range aclToken.Roles {
		roles[i] = map[string]interface{}{
			"id":   r.ID,
			"name": r.Name,
		}
	}

	serviceIdentities := make([]map[string]interface{}, len(aclToken.ServiceIdentities))
	for i, si := range aclToken.ServiceIdentities {
		serviceIdentities[i] = map[string]interface{}{
			"service_name": si.ServiceName,
			"datacenters":  si.Datacenters,
		}
	}

	nodeIdentities := make([]map[string]interface{}, len(aclToken.NodeIdentities))
	for i, ni := range aclToken.NodeIdentities {
		nodeIdentities[i] = map[string]interface{}{
			"node_name":  ni.NodeName,
			"datacenter": ni.Datacenter,
		}
	}

	var expirationTime string
	if aclToken.ExpirationTime != nil {
		expirationTime = aclToken.ExpirationTime.Format(time.RFC3339)
	}

	d.SetId(accessorID)

	sw := newStateWriter(d)
	sw.set("description", aclToken.Description)
	sw.set("local", aclToken.Local)
	sw.set("policies", policies)
	sw.set("roles", roles)
	sw.set("service_identities", serviceIdentities)
	sw.set("node_identities", nodeIdentities)
	sw.set("expiration_time", expirationTime)

	return sw.error()
}
