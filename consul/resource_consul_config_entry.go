package consul

import (
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	consulConfigEntryServiceDefaults = "service-defaults"
	consulConfigEntryProxyDefaults   = "proxy-defaults"
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
				ValidateFunc: validation.StringInSlice(
					[]string{
						consulConfigEntryServiceDefaults,
						consulConfigEntryProxyDefaults,
					},
					false,
				),
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"config": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceConsulConfigurationEntryCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceConsulConfigurationEntryUpdate(d, meta)
}

func resourceConsulConfigurationEntryUpdate(d *schema.ResourceData, meta interface{}) error {
	var config consulapi.ConfigEntry
	configEntries := getClient(meta).ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	switch configKind {
	case consulConfigEntryServiceDefaults:
		config = &consulapi.ServiceConfigEntry{
			Kind:     consulConfigEntryServiceDefaults,
			Name:     configName,
			Protocol: d.Get("protocol").(string),
		}
	case consulConfigEntryProxyDefaults:
		config = &consulapi.ProxyConfigEntry{
			Kind:   consulConfigEntryProxyDefaults,
			Name:   configName,
			Config: d.Get("config").(map[string]interface{}),
		}
	default:
		return fmt.Errorf("Config kind '%s' is not supported", configKind)
	}

	wOpts := &consulapi.WriteOptions{
		Token: d.Get("token").(string),
	}
	if _, _, err := configEntries.Set(config, wOpts); err != nil {
		return fmt.Errorf("Failed to set '%s' config entry: %#v", configName, err)
	}

	return resourceConsulConfigurationEntryRead(d, meta)
}

func resourceConsulConfigurationEntryRead(d *schema.ResourceData, meta interface{}) error {
	configEntries := getClient(meta).ConfigEntries()
	configKind := d.Get("kind").(string)
	configName := d.Get("name").(string)

	qOpts := &consulapi.QueryOptions{
		Token: d.Get("token").(string),
	}
	configEntry, _, err := configEntries.Get(configKind, configName, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "Unexpected response code: 404") {
			// The config entry has been removed
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Failed to fetch '%s' config entry: %#v", configName, err)
	}

	switch configKind {
	case consulConfigEntryProxyDefaults:
		if err = d.Set("config", configEntry.(*consulapi.ProxyConfigEntry).Config); err != nil {
			return fmt.Errorf("Failed to set 'config': %#v", err)
		}
	case consulConfigEntryServiceDefaults:
		if err = d.Set("protocol", configEntry.(*consulapi.ServiceConfigEntry).Protocol); err != nil {
			return fmt.Errorf("Failed to set 'protocol': %#v", err)
		}
	default:
		return fmt.Errorf("Config kind '%s' is not supported", configKind)
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
