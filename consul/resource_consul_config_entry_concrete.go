// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ConfigEntryImplementation is the common implementation for all specific
// config entries.
type ConfigEntryImplementation interface {
	GetKind() string
	GetDescription() string
	GetSchema() map[string]*schema.Schema
	Decode(d *schema.ResourceData) (consulapi.ConfigEntry, error)
	Write(ce consulapi.ConfigEntry, sw *stateWriter) error
}

func resourceFromConfigEntryImplementation(c ConfigEntryImplementation) *schema.Resource {
	return &schema.Resource{
		Description: c.GetDescription(),
		Schema:      c.GetSchema(),
		Create:      configEntryImplementationWrite(c),
		Update:      configEntryImplementationWrite(c),
		Read:        configEntryImplementationRead(c),
		Delete:      configEntryImplementationDelete(c),
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				var name, partition, namespace string
				switch len(parts) {
				case 1:
					name = parts[0]
				case 3:
					partition = parts[0]
					namespace = parts[1]
					name = parts[2]
				default:
					return nil, fmt.Errorf(`expected path of the form "<name>" or "<partition>/<namespace>/<name>"`)
				}

				d.SetId(name)
				sw := newStateWriter(d)
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
	}
}

func configEntryImplementationWrite(impl ConfigEntryImplementation) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		client, qOpts, wOpts := getClient(d, meta)

		configEntry, err := impl.Decode(d)
		if err != nil {
			return err
		}

		if _, _, err := client.ConfigEntries().Set(configEntry, wOpts); err != nil {
			return fmt.Errorf("failed to set '%s' config entry: %v", configEntry.GetName(), err)
		}
		_, _, err = client.ConfigEntries().Get(configEntry.GetKind(), configEntry.GetName(), qOpts)
		if err != nil {
			if strings.Contains(err.Error(), "Unexpected response code: 404") {
				return fmt.Errorf("failed to read config entry after setting it")
			}
			return fmt.Errorf("failed to read config entry: %v", err)
		}

		d.SetId(configEntry.GetName())
		return configEntryImplementationRead(impl)(d, meta)
	}
}

func configEntryImplementationRead(impl ConfigEntryImplementation) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		client, qOpts, _ := getClient(d, meta)
		name := d.Get("name").(string)

		fixQOptsForConfigEntry(name, impl.GetKind(), qOpts)

		ce, _, err := client.ConfigEntries().Get(impl.GetKind(), name, qOpts)
		if err != nil {
			if strings.Contains(err.Error(), "Unexpected response code: 404") {
				// The config entry has been removed
				d.SetId("")
				return nil
			}
			return fmt.Errorf("failed to fetch '%s' config entry: %v", name, err)
		}
		if ce == nil {
			d.SetId("")
			return nil
		}

		sw := newStateWriter(d)
		if err := impl.Write(ce, sw); err != nil {
			return err
		}
		return sw.error()
	}
}

func configEntryImplementationDelete(impl ConfigEntryImplementation) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		client, _, wOpts := getClient(d, meta)
		name := d.Get("name").(string)

		if _, err := client.ConfigEntries().Delete(impl.GetKind(), name, wOpts); err != nil {
			return fmt.Errorf("failed to delete '%s' config entry: %v", name, err)
		}
		d.SetId("")
		return nil
	}
}
