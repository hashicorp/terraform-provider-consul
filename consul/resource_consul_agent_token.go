package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulAgentToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulAgentTokenCreate,
		Update: nil,
		Read:   resourceConsulAgentTokenRead,
		Delete: resourceConsulAgentTokenDelete,

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"accessor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceConsulAgentTokenCreate(d *schema.ResourceData, meta interface{}) error {
	config := *meta.(*Config)
	config.Address = d.Get("address").(string)

	client, err := config.Client()
	if err != nil {
		return err
	}

	tokenType := d.Get("type").(string)
	accessorID := d.Get("accessor_id").(string)

	aclToken, _, err := client.ACL().TokenRead(accessorID, nil)
	if err != nil {
		return err
	}

	wOpts := &consulapi.WriteOptions{Token: config.Token}

	agent := client.Agent()
	switch tokenType {
	case "default":
		_, err = agent.UpdateDefaultACLToken(aclToken.SecretID, wOpts)
	case "agent":
		_, err = agent.UpdateAgentACLToken(aclToken.SecretID, wOpts)
	case "master":
		_, err = agent.UpdateAgentMasterACLToken(aclToken.SecretID, wOpts)
	case "replication":
		_, err = agent.UpdateReplicationACLToken(aclToken.SecretID, wOpts)
	default:
		return fmt.Errorf("Unknown token type '%s'", tokenType)
	}

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s-%s", config.Address, tokenType))
	return nil
}

func resourceConsulAgentTokenRead(d *schema.ResourceData, meta interface{}) error {
	// Not implemented. Consul doesn't provide an api to read agent tokens
	return nil
}

func resourceConsulAgentTokenDelete(d *schema.ResourceData, meta interface{}) error {
	// Delete is not implemented. We can call the api with an empty token, but if the
	// node doesn't exists anymore, the provider will always fail.

	d.SetId("")
	return nil
}
