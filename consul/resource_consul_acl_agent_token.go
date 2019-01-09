package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"

	"log"
)

func resourceConsulACLAgentToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLAgentTokenCreate,
		Read:   resourceConsulACLAgentTokenRead,
		Update: resourceConsulACLAgentTokenUpdate,
		Delete: resourceConsulACLAgentTokenDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The token description.",
			},
			"policies": {
				Type:     schema.TypeSet,
				Required: false,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of policies.",
			},
			"local": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     false,
				Description: "Flag to set the token local to the current datacenter.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The token.",
				Sensitive:   true,
			},
		},
	}
}

func resourceConsulACLAgentTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	log.Printf("[DEBUG] Creating ACL agent token")

	aclToken := consulapi.ACLToken{
		Description: d.Get("description").(string),
		Local:       d.Get("local").(bool),
	}

	iPolicies := d.Get("policies").(*schema.Set).List()
	policyLinks := make([]*consulapi.ACLTokenPolicyLink, 0, len(iPolicies))
	for _, iPolicy := range iPolicies {
		policyLinks = append(policyLinks, &consulapi.ACLTokenPolicyLink{
			Name: iPolicy.(string),
		})
	}

	if len(policyLinks) > 0 {
		aclToken.Policies = policyLinks
	}

	token, _, err := client.ACL().TokenCreate(&aclToken, nil)
	if err != nil {
		return fmt.Errorf("error creating ACL agent token: %s", err)
	}

	log.Printf("[DEBUG] Created ACL agent token %q", token.AccessorID)

	d.Set("token", token.SecretID)

	d.SetId(token.AccessorID)

	return resourceConsulACLAgentTokenRead(d, meta)
}

func resourceConsulACLAgentTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL agent token %q", id)

	aclToken, _, err := client.ACL().TokenRead(id, nil)
	if err != nil {
		log.Printf("[WARN] ACL agent token not found, removing from state")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read ACL agent token %q", id)

	d.Set("description", aclToken.Description)

	policies := make([]string, 0, len(aclToken.Policies))
	for _, policyLink := range aclToken.Policies {
		policies = append(policies, policyLink.Name)
	}

	d.Set("policies", policies)
	d.Set("local", aclToken.Local)

	return nil
}

func resourceConsulACLAgentTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()
	log.Printf("[DEBUG] Updating ACL token %q", id)

	aclToken := consulapi.ACLToken{
		AccessorID:  id,
		SecretID:    d.Get("token").(string),
		Description: d.Get("description").(string),
		Local:       d.Get("local").(bool),
	}

	if v, ok := d.GetOk("policies"); ok {
		vs := v.([]interface{})
		s := make([]*consulapi.ACLTokenPolicyLink, len(vs))
		for i, raw := range vs {
			s[i].Name = raw.(string)
		}
		aclToken.Policies = s
	}

	_, _, err := client.ACL().TokenUpdate(&aclToken, nil)
	if err != nil {
		return fmt.Errorf("error updating ACL token %q: %s", id, err)
	}
	log.Printf("[DEBUG] Updated ACL token %q", id)

	return resourceConsulACLTokenRead(d, meta)
}

func resourceConsulACLAgentTokenDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
