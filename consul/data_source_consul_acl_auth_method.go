package consul

import (
	"encoding/json"
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

			// Out parameters
			"type": {
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

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	return nil
}
