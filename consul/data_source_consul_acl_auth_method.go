package consul

import (
	"encoding/json"
	"fmt"

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
			},

			// Out parameters
			"type": {
				Type:     schema.TypeString,
				Computed: true,
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
				Computed: true,
			},

			"config": {
				Type:       schema.TypeMap,
				Computed:   true,
				Deprecated: "The config attribute is deprecated, please use config_json instead.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"config_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The raw configuration for this ACL auth method.",
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
	client, qOpts, _ := getClient(d, meta)
	name := d.Get("name").(string)

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

	configJson, err := json.Marshal(authMethod.Config)
	if err != nil {
		return fmt.Errorf("Failed to marshal 'config_json': %v", err)
	}
	if err = d.Set("config_json", string(configJson)); err != nil {
		return fmt.Errorf("Failed to set 'config_json': %v", err)
	}

	if err = d.Set("config", authMethod.Config); err != nil {
		// When a complex configuration is used we can fail to set config as it
		// will not support fields with maps or lists in them. In this case it
		// means that the user used the 'config_json' field, and since we
		// succeeded to set that and 'config' is deprecated, we can just use
		// an empty placeholder value and ignore the error.
		if c := d.Get("config_json").(string); c != "" {
			if err = d.Set("config", map[string]interface{}{}); err != nil {
				return fmt.Errorf("Failed to set 'config': %v", err)
			}
		} else {
			return fmt.Errorf("Failed to set 'config': %v", err)
		}
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
