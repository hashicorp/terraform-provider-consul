package consul

import (
	"fmt"
	"log"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulACLTokenRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLTokenRoleAttachmentCreate,
		Read:   resourceConsulACLTokenRoleAttachmentRead,
		Delete: resourceConsulACLTokenRoleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"token_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The token accessor id.",
			},
			"role": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The role name.",
			},
		},
	}
}

func resourceConsulACLTokenRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)

	log.Printf("[DEBUG] Creating ACL token role attachment")

	tokenID := d.Get("token_id").(string)

	aclToken, _, err := client.ACL().TokenRead(tokenID, nil)
	if err != nil {
		return fmt.Errorf("Token '%s' not found", tokenID)
	}

	roleName := d.Get("role").(string)
	for _, role := range aclToken.Roles {
		if role.Name == roleName {
			return fmt.Errorf("Role '%s' already attached to token", roleName)
		}
	}

	aclToken.Roles = append(aclToken.Roles, &consulapi.ACLTokenRoleLink{
		Name: roleName,
	})

	_, _, err = client.ACL().TokenUpdate(aclToken, nil)
	if err != nil {
		return fmt.Errorf("Error updating ACL token '%q' to set new role attachment: '%s'", tokenID, err)
	}

	id := fmt.Sprintf("%s:%s", tokenID, roleName)

	log.Printf("[DEBUG] Created ACL token role attachment '%q'", id)

	d.SetId(id)

	return resourceConsulACLTokenRoleAttachmentRead(d, meta)
}

func resourceConsulACLTokenRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL token role attachment '%q'", id)

	tokenID, roleName, err := parseTwoPartID(id)
	if err != nil {
		return fmt.Errorf("Invalid ACL token role attachment id '%q'", id)
	}

	aclToken, _, err := client.ACL().TokenRead(tokenID, nil)
	if err != nil {
		if strings.Contains(err.Error(), "ACL not found") {
			log.Printf("[WARN] ACL token not found, removing from state")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Failed to read token '%s': %v", id, err)
	}

	log.Printf("[DEBUG] Read ACL token %q", tokenID)

	roleFound := false
	for _, role := range aclToken.Roles {
		if role.Name == roleName {
			roleFound = true
			break
		}
	}
	if !roleFound {
		log.Printf("[WARN] ACL role not found in token, removing from state")
		d.SetId("")
		return nil
	}

	if err = d.Set("token_id", tokenID); err != nil {
		return fmt.Errorf("Error while setting 'token_id': %s", err)
	}
	if err = d.Set("role", roleName); err != nil {
		return fmt.Errorf("Error while setting 'role': %s", err)
	}

	return nil
}

func resourceConsulACLTokenRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL token role attachment '%q'", id)

	tokenID, roleName, err := parseTwoPartID(id)
	if err != nil {
		return fmt.Errorf("Invalid ACL token role attachment id '%q'", id)
	}

	aclToken, _, err := client.ACL().TokenRead(tokenID, nil)
	if err != nil {
		return fmt.Errorf("Token '%s' not found", tokenID)
	}

	for i, role := range aclToken.Roles {
		if role.Name == roleName {
			aclToken.Roles = append(aclToken.Roles[:i], aclToken.Roles[i+1:]...)
			break
		}
	}

	_, _, err = client.ACL().TokenUpdate(aclToken, nil)
	if err != nil {
		return fmt.Errorf("Error updating ACL token '%q' to set new role attachment: '%s'", tokenID, err)
	}
	log.Printf("[DEBUG] Deleted ACL token attachment role %q", id)

	return nil
}
