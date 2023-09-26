// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulKeys() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulKeysCreateUpdate,
		Update: resourceConsulKeysCreateUpdate,
		Read:   resourceConsulKeysRead,
		Delete: resourceConsulKeysDelete,

		SchemaVersion: 1,
		MigrateState:  resourceConsulKeysMigrateState,

		CustomizeDiff: func(d *schema.ResourceDiff, _ interface{}) error {
			if d.HasChange("key") {
				d.SetNewComputed("var")
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
							Type:       schema.TypeString,
							Optional:   true,
							Default:    "",
							Deprecated: "Using consul_keys resource to *read* is deprecated; please use consul_keys data source instead",
						},

						"path": {
							Type:     schema.TypeString,
							Required: true,
						},

						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"flags": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"cas": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						"default": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},

						"delete": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
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
				ForceNew: true,
			},

			"partition": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceConsulKeysCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	keyClient := newKeyClient(d, meta)
	if d.HasChange("key") {
		o, n := d.GetChange("key")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		// We'll keep track of what keys we add so that if a key is
		// in both the "remove" and "add" sets -- which will happen if
		// its value is changed in-place -- we will avoid writing the
		// value and then immediately removing it.
		addedPaths := make(map[string]bool)

		// We add before we remove because then it's possible to change
		// a key name (which will result in both an add and a remove)
		// without very temporarily having *neither* value in the store.
		// Instead, both will briefly be present, which should be less
		// disruptive in most cases.
		for _, raw := range add {
			_, path, sub, err := parseKey(raw)
			if err != nil {
				return err
			}

			// The name attribute is set when using consul_keys to read values
			// from the KV store. We must not overwrite the value when are
			// reading.
			name := sub["name"].(string)
			value := sub["value"].(string)
			if name != "" && value == "" {
				continue
			}

			flags := sub["flags"].(int)
			cas := sub["cas"].(int)
			if cas >= 0 {
				_, err := keyClient.Cas(path, value, flags, cas)
				if err != nil {
					return err
				}
				addedPaths[path] = true
				continue
			}
			if err := keyClient.Put(path, value, flags); err != nil {
				return err
			}
			addedPaths[path] = true
		}

		for _, raw := range remove {
			_, path, sub, err := parseKey(raw)
			if err != nil {
				return err
			}

			// Don't delete something we've just added.
			// (See explanation at the declaration of this variable above.)
			if addedPaths[path] {
				continue
			}

			shouldDelete, ok := sub["delete"].(bool)
			if !ok || !shouldDelete {
				continue
			}

			if err := keyClient.Delete(path); err != nil {
				return err
			}
		}
	}

	// Store the datacenter on this resource, which can be helpful for reference
	// in case it was read from the provider
	d.Set("datacenter", keyClient.qOpts.Datacenter)

	// The ID doesn't matter, since we use provider config, datacenter,
	// and key paths to address consul properly. So we just need to fill it in
	// with some value to indicate the resource has been created.
	d.SetId("consul")

	return resourceConsulKeysRead(d, meta)
}

func resourceConsulKeysRead(d *schema.ResourceData, meta interface{}) error {
	keyClient := newKeyClient(d, meta)

	vars := make(map[string]string)

	keys := d.Get("key").(*schema.Set).List()
	for _, raw := range keys {
		name, path, sub, err := parseKey(raw)
		if err != nil {
			return err
		}

		value, flags, err := keyClient.Get(path)
		if err != nil {
			return err
		}
		sub["flags"] = flags

		value = attributeValue(sub, value)
		if name != "" {
			// If 'name' is set then we'll update vars, for backward-compatibilty
			// with the pre-0.7 capability to read from Consul with this
			// resource.
			vars[name] = value
		} else {
			// If the 'name' attribute is set for this key then it was created
			// as a "write" block. We need to update the given value within the
			// block itself so that Terraform can detect when the
			// Consul-stored value has drifted from what was most recently
			// written by Terraform.
			// We don't do this for "read" blocks; that causes confusing diffs
			// because "value" should not be set for read-only key blocks.
			sub["value"] = value
		}
	}

	if err := d.Set("var", vars); err != nil {
		return err
	}
	if err := d.Set("key", keys); err != nil {
		return err
	}

	// Store the datacenter on this resource, which can be helpful for reference
	// in case it was read from the provider
	d.Set("datacenter", keyClient.qOpts.Datacenter)

	return nil
}

func resourceConsulKeysDelete(d *schema.ResourceData, meta interface{}) error {
	keyClient := newKeyClient(d, meta)

	// Clean up any keys that we're explicitly managing
	keys := d.Get("key").(*schema.Set).List()
	for _, raw := range keys {
		_, path, sub, err := parseKey(raw)
		if err != nil {
			return err
		}

		// Skip if the key is non-managed
		shouldDelete, ok := sub["delete"].(bool)
		if !ok || !shouldDelete {
			continue
		}
		if err := keyClient.Delete(path); err != nil {
			return err
		}
	}

	// Clear the ID
	d.SetId("")
	return nil
}

// parseKey is used to parse a key into a name, path, config or error
func parseKey(raw interface{}) (string, string, map[string]interface{}, error) {
	sub, ok := raw.(map[string]interface{})
	if !ok {
		return "", "", nil, fmt.Errorf("failed to unroll: %#v", raw)
	}

	key := sub["name"].(string)

	path, ok := sub["path"].(string)
	if !ok {
		return "", "", nil, fmt.Errorf("failed to get path for key '%s'", key)
	}
	return key, path, sub, nil
}

// attributeValue determines the value for a key, potentially
// using a default value if provided.
func attributeValue(sub map[string]interface{}, readValue string) string {
	// Use the value if given
	if readValue != "" {
		return readValue
	}

	// Use a default if given
	if raw, ok := sub["default"]; ok {
		switch def := raw.(type) {
		case string:
			return def
		case bool:
			return strconv.FormatBool(def)
		}
	}

	// No value
	return ""
}
