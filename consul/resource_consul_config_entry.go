package consul

import (
	"fmt"
	"strings"

	"github.com/hashicorp/consul/agent/structs"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConsulConfigurationEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulConfigurationEntryCreate,
		Update: resourceConsulConfigurationEntryCreate,
		Read:   resourceConsulConfigurationEntryRead,
		Delete: resourceConsulConfigurationEntryDelete,

		Schema: map[string]*schema.Schema{
			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceConsulConfigurationEntryCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceConsulConfigurationEntryUpdate(d, meta)
}

func resourceConsulConfigurationEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	kind := d.Get("kind").(string)
	name := d.Get("name").(string)
	config := d.Get("config").(map[string]interface{})

	config["kind"] = kind
	config["name"] = name

	configEntry, err := structs.DecodeConfigEntry(config)
	if err != nil {
		config = map[string]interface{}{
			"kind":   kind,
			"name":   name,
			"config": d.Get("config").(map[string]interface{}),
		}
		configEntry, err = structs.DecodeConfigEntry(config)
		if err != nil {
			return fmt.Errorf("Failed to decode config entry: %v", err)
		}
	}

	return fmt.Errorf("Succes: %#v", configEntry)

	// wOpts := &consulapi.WriteOptions{
	// 	Token: d.Get("token").(string),
	// }
	// if _, _, err := configEntries.Set(config, wOpts); err != nil {
	// 	return fmt.Errorf("Failed to set '%s' config entry: %#v", configName, err)
	// }

	// return resourceConsulConfigurationEntryRead(d, meta)
}

func resourceConsulConfigurationEntryRead(d *schema.ResourceData, meta interface{}) error {
	configEntries := getClient(meta).ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	qOpts := &consulapi.QueryOptions{
		Token: d.Get("token").(string),
	}
	_, _, err := configEntries.Get(configKind, configName, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "Unexpected response code: 404") {
			// The config entry has been removed
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Failed to fetch '%s' config entry: %#v", configName, err)
	}

	d.SetId(fmt.Sprintf("%s-%s", configKind, configName))

	return nil
}

func resourceConsulConfigurationEntryDelete(d *schema.ResourceData, meta interface{}) error {
	configEntries := getClient(meta).ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	wOpts := &consulapi.WriteOptions{
		Token: d.Get("token").(string),
	}
	if _, err := configEntries.Delete(configKind, configName, wOpts); err != nil {
		return fmt.Errorf("Failed to delete '%s' config entry: %#v", configName, err)
	}
	d.SetId("")
	return nil
}
