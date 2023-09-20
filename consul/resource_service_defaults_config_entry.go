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

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transparent_proxy": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"mutual_tls_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mesh_gateway": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"expose": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"external_sni": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"upstream_config": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"max_inbound_connections": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"local_connect_timeout_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"local_request_timeout_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"balance_inbound_connections": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rate_limits": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"envoy_extensions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
			},
			"meta": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
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

func resourceConsulServiceDefaultsConfigEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	configEntries := client.ConfigEntries()

	kind := d.Get("kind").(string)
	name := d.Get("name").(string)

	fixQOptsForServiceDefaultsConfigEntry(name, kind, qOpts)

	attributes := []string{"partition", "namespace", "protocol", "mode", "transparent_proxy", "mutual_tls_mode", "mesh_gateway",
		"expose", "external_sni", "upstream_config", "destination", "max_inbound_connections", "local_connect_timeout_ms",
		"local_request_timeout_ms", "balance_inbound_connections", "rate_limits", "envoy_extensions", "meta"}

	var configJson map[string]interface{}

	for _, attribute := range attributes {
		value := d.Get(attribute)
		if value != nil {
			configJson[attribute] = value
		}
	}

	err := d.Set("config_json", configJson)

	if err != nil {
		return fmt.Errorf("failed to create config json to make config entry")
	}

	configEntry, err := makeServiceDefaultsConfigEntry(kind, name, d.Get("config_json").(string), wOpts.Namespace, wOpts.Partition)
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
	client, qOpts, _ := getClient(d, meta)
	configEntries := client.ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	fixQOptsForConfigEntry(configName, configKind, qOpts)

	_, _, err := configEntries.Get(configKind, configName, qOpts)
	if err != nil {
		return err
	}

	return nil
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

func makeServiceDefaultsConfigEntry(kind, name, config, namespace, partition string) (consulapi.ConfigEntry, error) {
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
