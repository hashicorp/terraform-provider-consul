package consul

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulConfigEntry() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulConfigEntryRead,

		Schema: map[string]*schema.Schema{
			"kind": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The kind of config entry to read.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the config entry to read.",
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The partition the config entry is associated with.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The namespace the config entry is associated with.",
			},

			"config_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration of the config entry.",
			},
		},
	}
}

func dataSourceConsulConfigEntryRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	kind := d.Get("kind").(string)
	name := d.Get("name").(string)

	configEntry, _, err := client.ConfigEntries().Get(kind, name, qOpts)
	if err != nil {
		return fmt.Errorf("failed to read config entry %s/%s: %w", kind, name, err)
	}

	// Config Entries are too complex to write as maps for now so we save their JSON representation
	data, err := configEntryToMap(configEntry)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", kind, name))

	sw := newStateWriter(d)
	sw.setJson("config_json", data)

	return sw.error()
}
