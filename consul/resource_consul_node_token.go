package consul

import (
	"errors"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*
Example:

provider "consul" {
  address        = "localhost:8500"
  token          = "<Token>"
}

resource "consul_acl_policy" "agent" {
  name  = "agent"
  rules = <<-RULE
    node_prefix "" {
      policy = "write"
    }
    RULE
}

resource "consul_acl_token" "agent_token" {
  description = "my test token"
  policies = ["${consul_acl_policy.agent.name}"]
}

resource "consul_node_token" "agent_token" {
  type        = "agent"
  accessor_id = consul_acl_token.agent_token.accessor_id
}

*/

func resourceConsulNodeToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulNodeTokenCreate,
		Update: resourceConsulNodeTokenCreate,
		Read:   resourceConsulNodeTokenRead,
		Delete: resourceConsulNodeTokenDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"accessor_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConsulNodeTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	agent := client.Agent()

	tokenType := d.Get("type").(string)
	accessorID := d.Get("accessor_id").(string)

	aclToken, _, err := client.ACL().TokenRead(accessorID, nil)
	if err != nil {
		return err
	}

	wOpts := &consulapi.WriteOptions{Token: "9a1c245e-53d1-457e-877e-602c451599a1"} // #todo provider token

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

	d.SetId(accessorID)

	return nil
}

func resourceConsulNodeTokenRead(d *schema.ResourceData, meta interface{}) error {
	// #todo Consul doenst provides an api entpoint to read the agent/master/... token
	return nil
}

func resourceConsulNodeTokenDelete(d *schema.ResourceData, meta interface{}) error {
	// #todo impl. Same as "*Create" but with an empty token
	return errors.New("delete not implemented")
}
