package consul

import (
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

		Schema: map[string]*schema.Schema{
			"connect_provider": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"config": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceConsulCertificateAuthorityCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	caConfig := &consulapi.CAConfig{
		Provider: d.Get("connect_provider").(string),
		Config:   d.Get("config").(map[string]interface{}),
	}

	if _, err := client.Connect().CASetConfig(caConfig, wOpts); err != nil {
		return fmt.Errorf("Failed to set CA configuration: %v", err)
	}

	d.SetId("consul-ca")

	return resourceConsulCertificateAuthorityRead(d, meta)
}

func resourceConsulCertificateAuthorityRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	conf, _, err := client.Connect().CAGetConfig(qOpts)
	if err != nil {
		return fmt.Errorf("Failed to get CA configuration: %v", err)
	}

	if err = d.Set("connect_provider", conf.Provider); err != nil {
		return fmt.Errorf("Failed to set 'connect_provider': %v", err)
	}

	if err = d.Set("config", conf.Config); err != nil {
		return fmt.Errorf("Failed to set 'config': %v", err)
	}

	return nil
}
