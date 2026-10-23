// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConsulKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulKeysRead,

		Description: "The `consul_keys` datasource reads values from the Consul key/value store. This is a powerful way to dynamically set values in templates.",

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Description: "The datacenter to use. This overrides the agent's default datacenter and the datacenter in the provider setup.",
				Optional:    true,
				Computed:    true,
			},

			"error_on_missing_keys": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to return an error when a key is absent from the KV store and no default is configured. This defaults to `false`.",
			},

			"token": {
				Type:        schema.TypeString,
				Description: "The ACL token to use. This overrides the token that the agent provides by default.",
				Deprecated:  tokenDeprecationMessage,
				Optional:    true,
				Sensitive:   true,
			},

			"key": {
				Type:        schema.TypeSet,
				Description: "Specifies a key in Consul to be read. Supported values documented below. Multiple blocks supported.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "This is the name of the key. This value of the key is exposed as `var.<name>`. This is not the path of the key in Consul.",
							Required:    true,
						},

						"path": {
							Type:        schema.TypeString,
							Description: "This is the path in Consul that should be read or written to.",
							Required:    true,
						},

						"default": {
							Type:        schema.TypeString,
							Description: "This is the default value to set for `var.<name>` if the key does not exist in Consul. Defaults to an empty string.",
							Optional:    true,
						},
					},
				},
			},

			"var": {
				Type:        schema.TypeMap,
				Description: "For each name given, the corresponding attribute has the value of the key.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"namespace": {
				Type:        schema.TypeString,
				Description: "The namespace to lookup the keys.",
				Optional:    true,
			},

			"partition": {
				Type:        schema.TypeString,
				Description: "The partition to lookup the keys.",
				Optional:    true,
			},
		},
	}
}

func dataSourceConsulKeysRead(d *schema.ResourceData, meta interface{}) error {
	keyClient := newKeyClient(d, meta)

	vars := make(map[string]string)

	keys := d.Get("key").(*schema.Set).List()
	for _, raw := range keys {
		key, path, sub, err := parseKey(raw)
		if err != nil {
			return err
		}

		exist, value, _, err := keyClient.Get(path)
		if err != nil {
			return err
		}

		// This returns the value if it exists or the default value if one is set.
		// If the key does not exist and there is no default, value will be the
		// empty string.
		value = attributeValue(sub, value)

		if !exist && value == "" && d.Get("error_on_missing_keys").(bool) {
			// We return an error when the key does not exist, there is no default
			// and error_on_missing_keys has been set in the config.
			return fmt.Errorf("Key %q does not exist", path)
		}

		vars[key] = value
	}

	if err := d.Set("var", vars); err != nil {
		return err
	}

	// Store the datacenter on this resource, which can be helpful for reference
	// in case it was read from the provider
	d.Set("datacenter", keyClient.qOpts.Datacenter)

	d.SetId("-")

	return nil
}
