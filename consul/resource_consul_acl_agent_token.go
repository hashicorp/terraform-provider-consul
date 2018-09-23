package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/schema"

	"log"
)

func resourceConsulACLAgentToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLAgentTokenCreate,
		Read:   resourceConsulACLAgentTokenRead,
		Delete: resourceConsulACLAgentTokenDelete,

		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The ACL agent token to use.",
			},
		},
	}
}

func resourceConsulACLAgentTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	log.Printf("[DEBUG] Creating ACL agent token")

	var token string
	if v, ok := d.GetOk("token"); ok {
		token = v.(string)
	} else {
		var err error

		token, err = uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("error creating ACL agent token: %s", err)
		}
	}

	client.Agent().UpdateACLAgentToken(token, nil)

	_, err := client.Agent().UpdateACLAgentToken(token, nil)
	if err != nil {
		return fmt.Errorf("error creating ACL agent token: %s", err)
	}

	log.Printf("[DEBUG] Created ACL agent token %q", token)

	d.Set("token", token)
	d.SetId(token)

	return resourceConsulACLAgentTokenRead(d, meta)
}

func resourceConsulACLAgentTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()
	log.Printf("[DEBUG] Reading agent token %q", id)

	aclEntry, _, err := client.ACL().Info(id, nil)
	if err != nil {
		log.Printf("[WARN] ACL not found, removing from state")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read ACL %q", id)

	d.Set("name", aclEntry.Name)

	if aclEntry.Type == consulapi.ACLManagementType {
		d.Set("type", "management")
	} else {
		d.Set("type", "client")
	}

	d.Set("rules", aclEntry.Rules)

	return nil
}

func resourceConsulACLAgentTokenDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
