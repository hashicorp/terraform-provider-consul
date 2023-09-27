// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"encoding/json"
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const KindServiceDefaults = "service-defaults"

var upstreamConfigSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"partition": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"namespace": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"peer": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"envoy_listener_json": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"envoy_cluster_json": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"protocol": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"connect_timeout_ms": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"limits": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"max_connections": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"max_pending_requests": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"max_concurrent_requests": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		"passive_health_check": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"interval": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"max_failures": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"enforcing_consecutive_5xx": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"max_ejection_percent": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"base_ejection_time": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		"mesh_gateway": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"mode": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"balance_outbound_connections": {
			Type:     schema.TypeString,
			Optional: true,
		},
	},
}

var serviceDefaultsConfigEntrySchema = map[string]*schema.Schema{
	"kind": {
		Type:     schema.TypeString,
		Required: false,
		ForceNew: true,
		Computed: true,
	},

	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},

	"namespace": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},

	"partition": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "The partition the config entry is associated with.",
	},

	"meta": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},

	"protocol": {
		Type:     schema.TypeString,
		Required: true,
	},

	"balance_inbound_connections": {
		Type:     schema.TypeString,
		Optional: true,
	},

	"mode": {
		Type:     schema.TypeString,
		Optional: true,
	},

	"upstream_config": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"overrides": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     upstreamConfigSchema,
				},
				"defaults": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem:     upstreamConfigSchema,
				},
			},
		},
	},

	"transparent_proxy": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"outbound_listener_port": {
					Required: true,
					Type:     schema.TypeInt,
				},
				"dialed_directly": {
					Required: true,
					Type:     schema.TypeBool,
				},
			},
		},
	},

	"mutual_tls_mode": {
		Type:     schema.TypeString,
		Optional: true,
	},

	"envoy_extensions": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"required": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"arguments": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"consul_version": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"envoy_version": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	},

	"destination": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"port": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"addresses": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	},

	"local_connect_timeout_ms": {
		Type:     schema.TypeInt,
		Optional: true,
	},

	"max_inbound_connections": {
		Type:     schema.TypeInt,
		Optional: true,
	},

	"local_request_timeout_ms": {
		Type:     schema.TypeInt,
		Optional: true,
	},

	"mesh_gateway": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"mode": {
					Required: true,
					Type:     schema.TypeString,
				},
			},
		},
	},

	"external_sni": {
		Type:     schema.TypeString,
		Optional: true,
	},

	"expose": {
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"checks": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
				},
				"paths": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"path": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"local_path_port": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"listener_port": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"protocol": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
		},
	},
}

func resourceServiceDefaultsConfigEntry() *schema.Resource {

	return &schema.Resource{
		Create: resourceConsulServiceDefaultsConfigEntryUpdate,
		Update: resourceConsulServiceDefaultsConfigEntryUpdate,
		Read:   resourceConsulServiceDefaultsConfigEntryRead,
		Delete: resourceConsulServiceDefaultsConfigEntryDelete,
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
		Schema: serviceDefaultsConfigEntrySchema,
	}
}

func resourceConsulServiceDefaultsConfigEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	configEntries := client.ConfigEntries()

	name := d.Get("name").(string)

	configMap := make(map[string]interface{})

	configMap["kind"] = KindServiceDefaults

	configMap["name"] = name

	kind := configMap["kind"].(string)
	err := d.Set("kind", kind)

	if err != nil {
		return err
	}

	var attributes []string

	for key, _ := range serviceDefaultsConfigEntrySchema {
		attributes = append(attributes, key)
	}

	for _, attribute := range attributes {
		configMap[attribute] = d.Get(attribute)
	}

	formattedMap, err := formatKeys(configMap, formatKey)

	if err != nil {
		return err
	}

	configEntry, err := makeServiceDefaultsConfigEntry(name, formattedMap.(map[string]interface{}), wOpts.Namespace, wOpts.Partition)
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
	return resourceConsulServiceDefaultsConfigEntryRead(d, meta)
}

func resourceConsulServiceDefaultsConfigEntryRead(d *schema.ResourceData, meta interface{}) error {
	var attributes []string

	for key, _ := range serviceDefaultsConfigEntrySchema {
		attributes = append(attributes, key)
	}

	client, qOpts, _ := getClient(d, meta)
	configEntries := client.ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	configEntry, _, err := configEntries.Get(configKind, configName, qOpts)

	_, _, configJSON, err := parseConfigEntry(configEntry)
	if err != nil {
		return fmt.Errorf("failed to parse ConfigEntry: %v", err)
	}

	configJsonMap := make(map[string]interface{})

	err = json.Unmarshal([]byte(configJSON), &configJsonMap)
	if err != nil {
		return err
	}

	for _, attribute := range attributes {
		value := configJsonMap[formatKey(attribute)]
		if err = d.Set(attribute, value); err != nil {
			return fmt.Errorf("failed to set attribute': %v with value %v", attribute, configJsonMap[formatKey(attribute)])
		}
	}

	return err
}

func resourceConsulServiceDefaultsConfigEntryDelete(d *schema.ResourceData, meta interface{}) error {
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

func makeServiceDefaultsConfigEntry(name string, configMap map[string]interface{}, namespace, partition string) (consulapi.ConfigEntry, error) {
	configMap["kind"] = KindServiceDefaults
	configMap["name"] = name
	configMap["Namespace"] = namespace
	configMap["Partition"] = partition

	configEntry, err := consulapi.DecodeConfigEntry(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config entry: %v", err)
	}

	return configEntry, nil
}
