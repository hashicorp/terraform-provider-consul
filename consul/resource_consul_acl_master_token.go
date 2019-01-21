package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceConsulACLMasterToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLMasterTokenCreate,
		Read:   resourceConsulACLMasterTokenRead,
		Delete: resourceConsulACLMasterTokenDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The token description.",
			},
			"policies": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of policies.",
			},
			"local": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to set the token local to the current datacenter.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ACL master token.",
				Sensitive:   true,
			},
		},
	}
}

func resourceConsulACLMasterTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	log.Printf("[DEBUG] Creating ACL master token")

	aclToken, _, err := client.ACL().Bootstrap()
	if err != nil {
		return fmt.Errorf("error creating ACL master token: %s", err)
	}

	log.Printf("[DEBUG] Created ACL master token %q", aclToken.AccessorID)

	d.Set("description", aclToken.Description)

	policies := make([]string, 0, len(aclToken.Policies))
	for _, policyLink := range aclToken.Policies {
		policies = append(policies, policyLink.Name)
	}

	d.Set("policies", policies)
	d.Set("local", aclToken.Local)
	d.Set("token", aclToken.SecretID)

	d.SetId(aclToken.AccessorID)

	return resourceConsulACLMasterTokenRead(d, meta)
}

func resourceConsulACLMasterTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL master token %q", id)

	q := consulapi.QueryOptions{
		Token: d.Get("token").(string),
	}
	_, _, err := client.ACL().TokenRead(id, &q)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] ACL token not found, removing from state")

			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[DEBUG] Read ACL token %q", id)

	return nil
}

func resourceConsulACLMasterTokenDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
