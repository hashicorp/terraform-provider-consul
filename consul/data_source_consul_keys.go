// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulKeysRead,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"token": {
				Type:       schema.TypeString,
				Optional:   true,
				Sensitive:  true,
				Deprecated: tokenDeprecationMessage,
			},

			"key": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"path": {
							Type:     schema.TypeString,
							Required: true,
						},

						"default": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"var": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"partition": {
				Type:     schema.TypeString,
				Optional: true,
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

		if !exist && value == "" && meta.(*Config).ErrorOnMissingKey {
			// We return an error when the key does not exist, there is no default
			// and error_on_missing_key has been set in the config.
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
