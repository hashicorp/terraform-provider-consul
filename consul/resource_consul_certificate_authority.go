// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"encoding/json"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulCertificateAuthority() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulCertificateAuthorityCreate,
		Read:   resourceConsulCertificateAuthorityRead,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Description: "The `consul_certificate_authority` resource can be used to manage the configuration of the Certificate Authority used by [Consul Connect](https://www.consul.io/docs/connect/ca).\n\n-> **Note:** The keys in the `config` argument must be using Camel case.",

		Schema: map[string]*schema.Schema{
			"connect_provider": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the CA provider type to use.",
			},

			"config": {
				Type:          schema.TypeMap,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "The raw configuration to use for the chosen provider. For more information on configuring the Connect CA providers, see [Provider Config](https://developer.hashicorp.com/consul/docs/connect/ca).",
				Deprecated:    "The config attribute is deprecated, please use config_json instead.",
				ConflictsWith: []string{"config_json"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "" || new == "0"
				},
			},

			"config_json": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "The raw configuration to use for the chosen provider. For more information on configuring the Connect CA providers, see [Provider Config](https://developer.hashicorp.com/consul/docs/connect/ca).",
				ConflictsWith: []string{"config"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "" || new == "0"
				},
			},
		},
	}
}

func resourceConsulCertificateAuthorityCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	var config map[string]interface{}
	if c := d.Get("config_json").(string); c != "" {
		err := json.Unmarshal([]byte(c), &config)
		if err != nil {
			return fmt.Errorf("failed to read 'config_json': %v", err)
		}
	} else {
		config = d.Get("config").(map[string]interface{})
	}

	if len(config) == 0 {
		return fmt.Errorf("one of 'config' or 'config_json' must be set")
	}

	caConfig := &consulapi.CAConfig{
		Provider: d.Get("connect_provider").(string),
		Config:   config,
	}

	if _, err := client.Connect().CASetConfig(caConfig, wOpts); err != nil {
		return fmt.Errorf("failed to set CA configuration: %v", err)
	}

	d.SetId("consul-ca")

	return resourceConsulCertificateAuthorityRead(d, meta)
}

func resourceConsulCertificateAuthorityRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	conf, _, err := client.Connect().CAGetConfig(qOpts)
	if err != nil {
		return fmt.Errorf("failed to get CA configuration: %v", err)
	}

	sw := newStateWriter(d)

	sw.set("connect_provider", conf.Provider)
	sw.setJson("config_json", conf.Config)

	if err = d.Set("config", conf.Config); err != nil {
		// When a complex configuration is used we can fail to set config as it
		// will not support fields with maps or lists in them. In this case it
		// means that the user used the 'config_json' field, and since we
		// succeeded to set that and 'config' is deprecated, we can just use
		// an empty placeholder value and ignore the error.
		if c := d.Get("config_json").(string); c != "" {
			sw.set("config", map[string]interface{}{})
		} else {
			return fmt.Errorf("failed to set 'config': %v", err)
		}
	}

	return sw.error()
}
