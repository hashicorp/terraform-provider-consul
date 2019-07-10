package consul

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	// Datasource predicates
	aclTokenAccessorID = "accessor_id"

	// Output
	aclTokenSecretID    = "secret_id"
	aclTokenDescription = "description"
	aclTokenPolicies    = "policies"
	aclTokenLocal       = "local"
)

func dataSourceConsulACLToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulACLTokenRead,

		Schema: map[string]*schema.Schema{

			// Filters
			aclTokenAccessorID: {
				Required: true,
				Type:     schema.TypeString,
			},

			// Out parameters
			aclTokenSecretID: {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			aclTokenDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			aclTokenPolicies: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			aclTokenLocal: {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceConsulACLTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	accessorID := d.Get(aclTokenAccessorID).(string)
	aclToken, _, err := client.ACL().TokenRead(accessorID, nil)
	if err != nil {
		return err
	}

	policies := make([]string, 0, len(aclToken.Policies))
	for _, policyLink := range aclToken.Policies {
		policies = append(policies, policyLink.Name)
	}

	d.SetId(accessorID)
	if err = d.Set(aclTokenSecretID, aclToken.SecretID); err != nil {
		return fmt.Errorf("Error while setting '%s': %s", aclTokenSecretID, err)
	}
	if err = d.Set(aclTokenDescription, aclToken.Description); err != nil {
		return fmt.Errorf("Error while setting %s: %s", aclTokenDescription, accessorID)
	}
	if err = d.Set(aclTokenLocal, aclToken.Local); err != nil {
		return fmt.Errorf("Error while setting %s: %s", aclTokenLocal, accessorID)
	}
	if err = d.Set(aclTokenPolicies, policies); err != nil {
		return fmt.Errorf("Error while setting %s: %s", aclTokenPolicies, accessorID)
	}

	return nil
}
