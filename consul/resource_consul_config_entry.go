package consul

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
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
		return err
	}

	wOpts := &consulapi.WriteOptions{}
	if _, _, err := configEntries.Set(configEntry, wOpts); err != nil {
		return fmt.Errorf("Failed to set '%s' config entry: %#v", name, err)
	}
	qOpts := &consulapi.QueryOptions{}
	_, _, err = configEntries.Get(configEntry.GetKind(), configEntry.GetName(), qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "Unexpected response code: 404") {
			return fmt.Errorf(`failed to read config entry after setting it.
This may happen when some attributes have an unexpected value.
Read the documentation at https://www.consul.io/docs/agent/config-entries/%s.html
to see what values are expected.`, configEntry.GetKind())
		}
		return fmt.Errorf("failed to read config entry: %v", err)
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

func makeConfigEntry(kind, name, config string) (consulapi.ConfigEntry, error) {
	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(config), &configMap); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal configMap: %v", err)
	}

	configMap["kind"] = kind
	configMap["name"] = name

	configEntry, err := decodeConfigEntry(configMap)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode config entry: %v", err)
	}

	return configEntry, nil
}

func configEntryToMap(configEntry consulapi.ConfigEntry) (map[string]interface{}, error) {
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

func parseConfigEntry(configEntry consulapi.ConfigEntry) (string, string, string, error) {
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

// This could be removed if 'ErrorUnused: true' is upstreamed in
// https://github.com/hashicorp/consul/blob/fdd10dd8b872466f8a614c7ed76b3becf2c5fc4c/api/config_entry.go#L178
func decodeConfigEntry(raw map[string]interface{}) (consulapi.ConfigEntry, error) {
	var entry consulapi.ConfigEntry

	kindVal, ok := raw["Kind"]
	if !ok {
		kindVal, ok = raw["kind"]
	}
	if !ok {
		return nil, fmt.Errorf("Payload does not contain a kind/Kind key at the top level")
	}

	if kindStr, ok := kindVal.(string); ok {
		newEntry, err := consulapi.MakeConfigEntry(kindStr, "")
		if err != nil {
			return nil, err
		}
		entry = newEntry
	} else {
		return nil, fmt.Errorf("Kind value in payload is not a string")
	}

	decodeConf := &mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		Result:           &entry,
		WeaklyTypedInput: true,
		ErrorUnused:      true,
	}

	decoder, err := mapstructure.NewDecoder(decodeConf)
	if err != nil {
		return nil, err
	}

	return entry, decoder.Decode(raw)
}
