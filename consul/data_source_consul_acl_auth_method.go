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

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Out parameters
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"max_token_ttl": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"token_locality": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"config": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"namespace_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"selector": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bind_namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
	if err = d.Set("display_name", authMethod.DisplayName); err != nil {
		return fmt.Errorf("Failed to set 'display_name': %v", err)
	}
	if err = d.Set("max_token_ttl", authMethod.MaxTokenTTL.String()); err != nil {
		return fmt.Errorf("Failed to set 'max_token_ttl': %v", err)
	}
	if err = d.Set("token_locality", authMethod.TokenLocality); err != nil {
		return fmt.Errorf("Failed to set 'token_locality': %v", err)
	}

	rules := make([]interface{}, 0)
	for _, rule := range authMethod.NamespaceRules {
		rules = append(rules, map[string]interface{}{
			"selector":       rule.Selector,
			"bind_namespace": rule.BindNamespace,
		})
	}
	if err = d.Set("namespace_rule", rules); err != nil {
		return fmt.Errorf("Failed to set 'namespace_rule': %v", err)
	}

	return nil
}
