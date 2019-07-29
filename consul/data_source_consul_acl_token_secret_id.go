package consul

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceConsulACLTokenSecretID() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulACLTokenSecretIDRead,

		Schema: map[string]*schema.Schema{

			// Filters
			"accessor_id": {
				Required: true,
				Type:     schema.TypeString,
			},

			// Out parameters
			"secret_id": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceConsulACLTokenSecretIDRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	accessorID := d.Get("accessor_id").(string)
	aclToken, _, err := client.ACL().TokenRead(accessorID, nil)
	if err != nil {
		return err
	}

	d.SetId(accessorID)
	if err = d.Set("secret_id", aclToken.SecretID); err != nil {
		return fmt.Errorf("Error while setting '%s': %s", "secret_id", err)
	}
	return nil
}
