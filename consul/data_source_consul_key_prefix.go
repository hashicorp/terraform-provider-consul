package consul

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceConsulKeyPrefix() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulKeyPrefixRead,

		Schema: map[string]*schema.Schema{
			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"token": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"path_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"subkey": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"path": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"default": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"var": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},

			"subkeys": &schema.Schema{
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
	client := meta.(*consulapi.Client)
	kv := client.KV()
	token := d.Get("token").(string)
	dc, err := getDC(d, client)
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
		value, err := keyClient.Get(fullPath)
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
		subKeys, err := keyClient.GetUnderPrefix(pathPrefix)
		if err != nil {
			return err
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
