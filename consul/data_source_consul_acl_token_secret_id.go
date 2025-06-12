// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"partition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Out parameters
			"secret_id": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"pgp_key": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"encrypted_secret_id": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceConsulACLTokenSecretIDRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	accessorID := d.Get("accessor_id").(string)

	aclToken, _, err := client.ACL().TokenRead(accessorID, qOpts)
	if err != nil {
		return err
	}

	d.SetId(accessorID)

	if err = d.Set("pgp_key", ""); err != nil {
		return fmt.Errorf("Error while setting 'pgp_key': %s", err)
	}

	if err = d.Set("secret_id", aclToken.SecretID); err != nil {
		return fmt.Errorf("Error while setting 'secret_id': %s", err)
	}

	if err = d.Set("encrypted_secret_id", ""); err != nil {
		return fmt.Errorf("Error while setting 'encrypted_secret_id': %s", err)
	}

	return nil
}
