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

var serviceSplitterConfigEntrySchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"kind": {
		Type:     schema.TypeString,
		Required: false,
		ForceNew: true,
		Computed: true,
	},
	"partition": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"namespace": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"meta": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"splits": {
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"weight": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"service": {
					Type:     schema.TypeString,
					Required: true,
				},
				"service_subset": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"namespace": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"partition": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"request_headers": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"add": {
								Type:     schema.TypeMap,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"set": {
								Type:     schema.TypeMap,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"remove": {
								Type:     schema.TypeMap,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
				"response_headers": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"add": {
								Type:     schema.TypeMap,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"set": {
								Type:     schema.TypeMap,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"remove": {
								Type:     schema.TypeMap,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
		},
	},
}

func resourceServiceSplitterConfigEntry() *schema.Resource {

	return &schema.Resource{
		Create: resourceConsulServiceSplitterConfigEntryUpdate,
		Update: resourceConsulServiceSplitterConfigEntryUpdate,
		Read:   resourceConsulServiceSplitterConfigEntryRead,
		Delete: resourceConsulServiceSplitterConfigEntryDelete,
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
		Schema: serviceSplitterConfigEntrySchema,
	}
}

func fixQOptsForServiceDefaultsConfigEntry(name, kind string, qOpts *consulapi.QueryOptions) {
	// exported-services config entries are weird in that their name correspond
	// to the partition they are created in, see
	// https://www.consul.io/docs/connect/config-entries/exported-services#configuration-parameters
	if kind == "exported-services" && name != "default" {
		qOpts.Partition = name
	}
}

func formatKey(key string) string {
	tokens := strings.Split(key, "_")
	keyToReturn := ""
	for _, token := range tokens {
		keyToReturn += strings.ToTitle(token)
	}
	return keyToReturn
}

func isSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array
}

func isMap(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}

func isSetSchema(v interface{}) bool {
	return reflect.TypeOf(v).String() == "*schema.Set"
}

func isStruct(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

func formatKeys(config interface{}, formatFunc func(string) string) (interface{}, error) {
	if isMap(config) {
		formattedMap := make(map[string]interface{})
		for key, value := range config.(map[string]interface{}) {
			formattedKey := formatFunc(key)
			formattedValue, err := formatKeys(value, formatKey)
			if err != nil {
				return nil, err
			}
			if formattedValue != nil {
				formattedMap[formattedKey] = formattedValue
			}
		}
		return formattedMap, nil
	} else if isSlice(config) {
		var newSlice []interface{}
		listValue := config.([]interface{})
		for _, elem := range listValue {
			newElem, err := formatKeys(elem, formatKey)
			if err != nil {
				return nil, err
			}
			newSlice = append(newSlice, newElem)
		}
		return newSlice, nil
	} else if isStruct(config) {
		var modifiedStruct map[string]interface{}
		jsonValue, err := json.Marshal(config)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonValue, &modifiedStruct)
		if err != nil {
			return nil, err
		}
		formattedStructKeys, err := formatKeys(modifiedStruct, formatKey)
		if err != nil {
			return nil, err
		}
		return formattedStructKeys, nil
	} else if isSetSchema(config) {
		valueList := config.(*schema.Set).List()
		if len(valueList) > 0 {
			formattedSetValue, err := formatKeys(valueList[0], formatKey)
			if err != nil {
				return nil, err
			}
			return formattedSetValue, nil
		}
		return nil, nil
	}
	return config, nil
}

func resourceConsulServiceSplitterConfigEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	configEntries := client.ConfigEntries()

	name := d.Get("name").(string)

	configMap := make(map[string]interface{})
	configMap["kind"] = "service-defaults"

	configMap["name"] = name

	kind := configMap["kind"].(string)
	err := d.Set("kind", kind)

	if err != nil {
		return err
	}

	fixQOptsForServiceDefaultsConfigEntry(name, kind, qOpts)

	var attributes []string

	for key, _ := range serviceSplitterConfigEntrySchema {
		attributes = append(attributes, key)
	}

	for _, attribute := range attributes {
		configMap[attribute] = d.Get(attribute)
	}

	formattedMap, err := formatKeys(configMap, formatKey)
	if err != nil {
		return err
	}

	configEntry, err := makeServiceDefaultsConfigEntry(kind, name, formattedMap.(map[string]interface{}), wOpts.Namespace, wOpts.Partition)
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
	return resourceConsulServiceSplitterConfigEntryRead(d, meta)
}

func resourceConsulServiceSplitterConfigEntryRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	configEntries := client.ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	fixQOptsForConfigEntry(configName, configKind, qOpts)

	_, _, err := configEntries.Get(configKind, configName, qOpts)
	return err
}

func resourceConsulServiceSplitterConfigEntryDelete(d *schema.ResourceData, meta interface{}) error {
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

func makeServiceDefaultsConfigEntry(kind, name string, configMap map[string]interface{}, namespace, partition string) (consulapi.ConfigEntry, error) {
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
