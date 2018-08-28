package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceConsulACLMasterToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLMasterTokenCreate,
		Read:   resourceConsulACLMasterTokenRead,
		Delete: resourceConsulACLMasterTokenDelete,

		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ACL master token.",
			},
		},
	}
}

func resourceConsulACLMasterTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	log.Printf("[DEBUG] Creating ACL master token")

	token, _, err := client.ACL().Bootstrap()
	if err != nil {
		return fmt.Errorf("error creating ACL master token: %s", err)
	}

	log.Printf("[DEBUG] Created ACL master token %q", token)

	d.Set("token", token)

	d.SetId(token)

	return resourceConsulACLMasterTokenRead(d, meta)
}

func resourceConsulACLMasterTokenRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceConsulACLMasterTokenDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
