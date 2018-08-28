package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceConsulACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLCreate,
		Read:   resourceConsulACLRead,
		Update: resourceConsulACLUpdate,
		Delete: resourceConsulACLDelete,

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The ACL ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: "The ACL name.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				Description:  "The ACL type.",
				ValidateFunc: validation.StringInSlice([]string{"client", "management"}, false),
				Default:      "client",
			},
			"rules": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: "The ACL rules.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ACL token.",
			},
		},
	}
}

func resourceConsulACLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	log.Printf("[DEBUG] Creating ACL")

	aclEntry := consulapi.ACLEntry{
		Name:  d.Get("name").(string),
		Type:  consulapi.ACLClientType,
		Rules: d.Get("rules").(string),
	}

	if d.Get("uuid") != "" {
		aclEntry.ID = d.Get("uuid").(string)
	}

	if d.Get("type") == "management" {
		aclEntry.Type = consulapi.ACLManagementType
	}

	token, _, err := client.ACL().Create(&aclEntry, nil)
	if err != nil {
		return fmt.Errorf("error creating ACL: %s", err)
	}

	log.Printf("[DEBUG] Created ACL %q", token)

	d.Set("token", token)

	d.SetId(token)

	return resourceConsulACLRead(d, meta)
}

func resourceConsulACLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL %q", id)

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

func resourceConsulACLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()
	log.Printf("[DEBUG] Updating ACL %q", id)

	aclEntry := consulapi.ACLEntry{
		ID:    id,
		Name:  d.Get("name").(string),
		Type:  consulapi.ACLClientType,
		Rules: d.Get("rules").(string),
	}

	if d.Get("type") == "management" {
		aclEntry.Type = consulapi.ACLManagementType
	}

	_, err := client.ACL().Update(&aclEntry, nil)
	if err != nil {
		return fmt.Errorf("error updating ACL %q: %s", id, err)
	}
	log.Printf("[DEBUG] Updated ACL %q", id)

	return resourceConsulACLRead(d, meta)
}

func resourceConsulACLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()

	log.Printf("[DEBUG] Deleting ACL %q", id)
	_, err := client.ACL().Destroy(id, nil)
	if err != nil {
		return fmt.Errorf("error deleting ACL %q: %s", id, err)
	}
	log.Printf("[DEBUG] Deleted ACL %q", id)

	return nil
}
