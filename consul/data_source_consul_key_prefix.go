package consul

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceConsulKeyPrefix() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulKeyPrefixRead,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"path_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},

			"subkey": {
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
			},

			"subkeys": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceConsulKeyPrefixRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)

	kv := client.KV()
	token := d.Get("token").(string)
	dc, err := getDC(d, client, meta)
	if err != nil {
		return err
	}

	keyClient := newKeyClient(kv, dc, token)

	pathPrefix := d.Get("path_prefix").(string)

	vars := make(map[string]string)

	keys := d.Get("subkey").(*schema.Set).List()
	for _, raw := range keys {
		key, path, sub, err := parseKey(raw)
		if err != nil {
			return err
		}

		fullPath := pathPrefix + path
		value, _, err := keyClient.Get(fullPath)
		if err != nil {
			return err
		}

		value = attributeValue(sub, value)
		vars[key] = value
	}

	if err := d.Set("var", vars); err != nil {
		return err
	}

	if len(keys) <= 0 {
		pairs, err := keyClient.GetUnderPrefix(pathPrefix)
		if err != nil {
			return err
		}
		subKeys := map[string]string{}
		for _, pair := range pairs {
			subKey := pair.Key[len(pathPrefix):]
			subKeys[subKey] = string(pair.Value)
		}
		d.Set("subkeys", subKeys)
	}

	// Store the datacenter on this resource, which can be helpful for reference
	// in case it was read from the provider
	d.Set("datacenter", dc)
	d.Set("path_prefix", pathPrefix)

	d.SetId("-")

	return nil
}
