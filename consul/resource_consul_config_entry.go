package consul

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/consul/api"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulConfigEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulConfigEntryUpdate,
		Update: resourceConsulConfigEntryUpdate,
		Read:   resourceConsulConfigEntryRead,
		Delete: resourceConsulConfigEntryDelete,

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

			"config_json": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffConfigJSON,
			},
		},
	}
}

func resourceConsulConfigEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	configEntries := getClient(meta).ConfigEntries()

	kind := d.Get("kind").(string)
	name := d.Get("name").(string)

	configEntry, err := makeConfigEntry(kind, name, d.Get("config_json").(string))
	if err != nil {
		return fmt.Errorf("Failed to decode config entry: %v", err)
	}

	wOpts := &consulapi.WriteOptions{}
	if _, _, err := configEntries.Set(configEntry, wOpts); err != nil {
		return fmt.Errorf("Failed to set '%s' config entry: %#v", name, err)
	}

	d.SetId(fmt.Sprintf("%s-%s", kind, name))
	return resourceConsulConfigEntryRead(d, meta)
}

func resourceConsulConfigEntryRead(d *schema.ResourceData, meta interface{}) error {
	configEntries := getClient(meta).ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	qOpts := &consulapi.QueryOptions{}
	configEntry, _, err := configEntries.Get(configKind, configName, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "Unexpected response code: 404") {
			// The config entry has been removed
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Failed to fetch '%s' config entry: %#v", configName, err)
	}

	_, _, configJSON, err := parseConfigEntry(configEntry)
	if err != nil {
		return fmt.Errorf("Failed to parse ConfigEntry: %v", err)
	}

	if err = d.Set("config_json", string(configJSON)); err != nil {
		return fmt.Errorf("Failed to set 'config_json': %v", err)
	}

	return nil
}

func resourceConsulConfigEntryDelete(d *schema.ResourceData, meta interface{}) error {
	configEntries := getClient(meta).ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	wOpts := &consulapi.WriteOptions{}
	if _, err := configEntries.Delete(configKind, configName, wOpts); err != nil {
		return fmt.Errorf("Failed to delete '%s' config entry: %#v", configName, err)
	}
	d.SetId("")
	return nil
}

func makeConfigEntry(kind, name, config string) (api.ConfigEntry, error) {
	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(config), &configMap); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal configMap: %v", err)
	}

	configMap["kind"] = kind
	configMap["name"] = name

	configEntry, err := api.DecodeConfigEntry(configMap)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode config entry: %v", err)
	}

	return configEntry, nil
}

func configEntryToMap(configEntry api.ConfigEntry) (map[string]interface{}, error) {
	marshalled, err := json.Marshal(configEntry)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal %v: %v", configEntry, err)
	}

	var configMap map[string]interface{}
	if err = json.Unmarshal(marshalled, &configMap); err != nil {
		// This should never happen
		return nil, fmt.Errorf("Failed to unmarshal %v: %v", marshalled, err)
	}

	// Remove the fields unrelated to the configEntry
	delete(configMap, "CreateIndex")
	delete(configMap, "ModifyIndex")
	delete(configMap, "Kind")
	delete(configMap, "Name")

	return configMap, nil
}

func parseConfigEntry(configEntry api.ConfigEntry) (string, string, string, error) {
	// We need to transform the ConfigEntry to a representation that works with
	// the config_json attribute
	name := configEntry.GetName()
	kind := configEntry.GetKind()

	configMap, err := configEntryToMap(configEntry)
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to convert config entry to map")
	}

	configJSON, err := json.Marshal(configMap)
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to marshal %v: %v", configMap, err)
	}

	return kind, name, string(configJSON), nil
}
func diffConfigJSON(k, old, new string, d *schema.ResourceData) bool {
	kind := d.Get("kind").(string)
	name := d.Get("name").(string)

	oldEntry, err := makeConfigEntry(kind, name, old)
	if err != nil {
		return false
	}
	newEntry, err := makeConfigEntry(kind, name, new)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(oldEntry, newEntry)
}
