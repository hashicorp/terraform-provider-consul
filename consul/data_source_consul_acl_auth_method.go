// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consul

import (
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

			"partition": {
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
		return fmt.Errorf("failed to get auth-method: %v", err)
	}
	if authMethod == nil {
		return fmt.Errorf("could not find auth-method '%s'", name)
	}

	d.SetId(authMethod.Name)

	sw := newStateWriter(d)
	sw.set("type", authMethod.Type)
	sw.set("description", authMethod.Description)
	sw.setJson("config_json", authMethod.Config)

	if err = d.Set("config", authMethod.Config); err != nil {
		// When a complex configuration is used we can fail to set config as it
		// will not support fields with maps or lists in them. In this case it
		// means that the user used the 'config_json' field, and since we
		// succeeded to set that and 'config' is deprecated, we can just use
		// an empty placeholder value and ignore the error.
		if c := d.Get("config_json").(string); c != "" {
			sw.set("config", map[string]interface{}{})
		} else {
			return fmt.Errorf("failed to set 'config': %v", err)
		}
	}
	sw.set("display_name", authMethod.DisplayName)
	sw.set("max_token_ttl", authMethod.MaxTokenTTL.String())
	sw.set("token_locality", authMethod.TokenLocality)

	rules := make([]interface{}, 0)
	for _, rule := range authMethod.NamespaceRules {
		rules = append(rules, map[string]interface{}{
			"selector":       rule.Selector,
			"bind_namespace": rule.BindNamespace,
		})
	}
	sw.set("namespace_rule", rules)

	return sw.error()
}
