package consul

import (
	"fmt"

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

			"local": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceConsulACLTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	accessorID := d.Get("accessor_id").(string)
	aclToken, _, err := client.ACL().TokenRead(accessorID, nil)
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

	d.SetId(accessorID)
	if err = d.Set("description", aclToken.Description); err != nil {
		return fmt.Errorf("Error while setting 'description': %s", err)
	}
	if err = d.Set("local", aclToken.Local); err != nil {
		return fmt.Errorf("Error while setting 'local': %s", err)
	}
	if err = d.Set("policies", policies); err != nil {
		return fmt.Errorf("Error while setting 'policies': %s", err)
	}

	return nil
}
