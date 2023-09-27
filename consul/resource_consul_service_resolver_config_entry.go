// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const KindServiceResolver = "service-resolver"

var serviceResolverConfigEntrySchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
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
		ForceNew: true,
	},
	"namespace": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"meta": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"connect_timeout": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"request_timeout": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"subsets": {
		Type:     schema.TypeMap,
		Required: true,
		Elem: &schema.Schema{
			Type:     schema.TypeMap,
			Required: true,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"filter": {
					Type:     schema.TypeString,
					Required: true,
				},
				"only_passing": {
					Type:     schema.TypeBool,
					Required: true,
				},
			}},
		},
	},
	"default_subset": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"redirect": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service": {
					Type:     schema.TypeString,
					Optional: true,
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
				"sameness_group": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"datacenter": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"peer": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	},
	"failover": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"service_subset": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"namespace": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"sameness_group": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"datacenters": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     schema.TypeString,
				},
				"targets": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"service": {
								Type:     schema.TypeString,
								Optional: true,
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
							"datacenter": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"peer": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
		},
	},
	"load_balancer": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policy": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"least_request_config": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"choice_count": {
								Type:     schema.TypeInt,
								Optional: true,
							},
						},
					},
				},
				"ring_hash_config": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"minimum_ring_size": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"maximum_ring_size": {
								Type:     schema.TypeInt,
								Optional: true,
							},
						},
					},
				},
				"hash_policies": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"field": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"field_value": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"cookie_config": {
								Type:     schema.TypeString,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"session": {
											Type:     schema.TypeBool,
											Optional: true,
										},
										"ttl": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"path": {
											Type:     schema.TypeString,
											Optional: true,
										},
									},
								},
							},
							"source_ip": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"terminal": {
								Type:     schema.TypeBool,
								Optional: true,
							},
						},
					},
				},
			},
		},
	},
}

func resourceServiceResolverConfigEntry() *schema.Resource {

	return &schema.Resource{
		Create: resourceConsulServiceResolverConfigEntryUpdate,
		Update: resourceConsulServiceResolverConfigEntryUpdate,
		Read:   resourceConsulServiceResolverConfigEntryRead,
		Delete: resourceConsulServiceResolverConfigEntryDelete,
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
		Schema: serviceResolverConfigEntrySchema,
	}
}

func resourceConsulServiceResolverConfigEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	configEntries := client.ConfigEntries()

	name := d.Get("name").(string)

	configMap := make(map[string]interface{})
	configMap["kind"] = KindServiceResolver

	configMap["name"] = name

	kind := configMap["kind"].(string)
	err := d.Set("kind", kind)

	if err != nil {
		return err
	}

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

	configEntry, err := makeServiceResolverConfigEntry(name, formattedMap.(map[string]interface{}), wOpts.Namespace, wOpts.Partition)
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

func resourceConsulServiceResolverConfigEntryRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	configEntries := client.ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	fixQOptsForConfigEntry(configName, configKind, qOpts)

	_, _, err := configEntries.Get(configKind, configName, qOpts)
	return err
}

func resourceConsulServiceResolverConfigEntryDelete(d *schema.ResourceData, meta interface{}) error {
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

func makeServiceResolverConfigEntry(name string, configMap map[string]interface{}, namespace, partition string) (consulapi.ConfigEntry, error) {
	configMap["kind"] = KindServiceResolver
	configMap["name"] = name
	configMap["Namespace"] = namespace
	configMap["Partition"] = partition

	configEntry, err := consulapi.DecodeConfigEntry(configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config entry: %v", err)
	}

	return configEntry, nil
}
