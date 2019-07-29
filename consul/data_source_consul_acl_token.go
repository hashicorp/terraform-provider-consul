package consul

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

	policies := make([]string, 0, len(aclToken.Policies))
	for _, policyLink := range aclToken.Policies {
		policies = append(policies, policyLink.Name)
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
