// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulConfigEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulConfigEntryUpdate,
		Update: resourceConsulConfigEntryUpdate,
		Read:   resourceConsulConfigEntryRead,
		Delete: resourceConsulConfigEntryDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				var kind, name, partition, namespace string
				switch len(parts) {
				case 2:
					kind = parts[0]
					name = parts[1]
				case 4:
					partition = parts[0]
					namespace = parts[1]
					kind = parts[2]
					name = parts[3]
				default:
					return nil, fmt.Errorf(`expected path of the form "<kind>/<name>" or "<partition>/<namespace>/<kind>/<name>"`)
				}

				d.SetId(fmt.Sprintf("%s-%s", kind, name))
				sw := newStateWriter(d)
				sw.set("kind", kind)
				sw.set("name", name)
				sw.set("partition", partition)
				sw.set("namespace", namespace)

				err := sw.error()
				if err != nil {
					return nil, err
				}

				return []*schema.ResourceData{d}, nil
			},
		},

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

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The partition the config entry is associated with.",
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
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

func fixQOptsForConfigEntry(name, kind string, qOpts *consulapi.QueryOptions) {
	// exported-services config entries are weird in that their name correspond
	// to the partition they are created in, see
	// https://www.consul.io/docs/connect/config-entries/exported-services#configuration-parameters
	if kind == "exported-services" && name != "default" {
		qOpts.Partition = name
	}
}

func resourceConsulConfigEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	configEntries := client.ConfigEntries()

	kind := d.Get("kind").(string)
	name := d.Get("name").(string)

	fixQOptsForConfigEntry(name, kind, qOpts)

	configEntry, err := makeConfigEntry(kind, name, d.Get("config_json").(string), wOpts.Namespace, wOpts.Partition)
	if err != nil {
		return err
	}

	if _, _, err := configEntries.Set(configEntry, wOpts); err != nil {
		return fmt.Errorf("failed to set '%s' config entry: %v", name, err)
	}
	_, _, err = configEntries.Get(configEntry.GetKind(), configEntry.GetName(), qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "Unexpected response code: 404") {
			return fmt.Errorf(`failed to read config entry after setting it.
This may happen when some attributes have an unexpected value.
Read the documentation at https://www.consul.io/docs/agent/config-entries/%s.html
to see what values are expected`, configEntry.GetKind())
		}
		return fmt.Errorf("failed to read config entry: %v", err)
	}

	d.SetId(fmt.Sprintf("%s-%s", kind, name))
	return resourceConsulConfigEntryRead(d, meta)
}

func resourceConsulConfigEntryRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	configEntries := client.ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	fixQOptsForConfigEntry(configName, configKind, qOpts)

	configEntry, _, err := configEntries.Get(configKind, configName, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "Unexpected response code: 404") {
			// The config entry has been removed
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to fetch '%s' config entry: %v", configName, err)
	}

	_, _, configJSON, err := parseConfigEntry(configEntry)
	if err != nil {
		return fmt.Errorf("failed to parse ConfigEntry: %v", err)
	}

	if err = d.Set("config_json", string(configJSON)); err != nil {
		return fmt.Errorf("failed to set 'config_json': %v", err)
	}

	return nil
}

func resourceConsulConfigEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	configEntries := client.ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	if _, err := configEntries.Delete(configKind, configName, wOpts); err != nil {
		return fmt.Errorf("failed to delete '%s' config entry: %v", configName, err)
	}
	d.SetId("")
	return nil
}

func makeConfigEntry(kind, name, config, namespace, partition string) (consulapi.ConfigEntry, error) {
	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(config), &configMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configMap: %v", err)
	}

	configMap["kind"] = kind
	configMap["name"] = name
	configMap["Namespace"] = namespace
	configMap["Partition"] = partition

	configEntry, err := consulapi.DecodeConfigEntry(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config entry: %v", err)
	}

	return configEntry, nil
}

func configEntryToMap(configEntry consulapi.ConfigEntry) (map[string]interface{}, error) {
	marshalled, err := json.Marshal(configEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %v: %v", configEntry, err)
	}

	var configMap map[string]interface{}
	if err = json.Unmarshal(marshalled, &configMap); err != nil {
		// This should never happen
		return nil, fmt.Errorf("failed to unmarshal %v: %v", marshalled, err)
	}

	// Remove the fields unrelated to the configEntry
	delete(configMap, "CreateIndex")
	delete(configMap, "ModifyIndex")
	delete(configMap, "Kind")
	delete(configMap, "Name")
	delete(configMap, "Namespace")
	delete(configMap, "Partition")

	return configMap, nil
}

func parseConfigEntry(configEntry consulapi.ConfigEntry) (string, string, string, error) {
	// We need to transform the ConfigEntry to a representation that works with
	// the config_json attribute
	name := configEntry.GetName()
	kind := configEntry.GetKind()

	configMap, err := configEntryToMap(configEntry)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to convert config entry to map")
	}

	configJSON, err := json.Marshal(configMap)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to marshal %v: %v", configMap, err)
	}

	return kind, name, string(configJSON), nil
}
func diffConfigJSON(k, old, new string, d *schema.ResourceData) bool {
	kind := d.Get("kind").(string)
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	partition := d.Get("partition").(string)

	oldEntry, err := makeConfigEntry(kind, name, old, namespace, partition)
	if err != nil {
		return false
	}
	newEntry, err := makeConfigEntry(kind, name, new, namespace, partition)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(oldEntry, newEntry)
}
