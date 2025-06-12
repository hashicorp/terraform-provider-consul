// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	client, qOpts, wOpts := getClient(d, meta)

	tokenID := d.Get("token_id").(string)

	aclToken, _, err := client.ACL().TokenRead(tokenID, qOpts)
	if err != nil {
		return fmt.Errorf("token '%s' not found", tokenID)
	}

	roleName := d.Get("role").(string)
	for _, role := range aclToken.Roles {
		if role.Name == roleName {
			return fmt.Errorf("role '%s' already attached to token", roleName)
		}
	}

	aclToken.Roles = append(aclToken.Roles, &consulapi.ACLTokenRoleLink{
		Name: roleName,
	})

	_, _, err = client.ACL().TokenUpdate(aclToken, wOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL token '%q' to set new role attachment: '%s'", tokenID, err)
	}

	id := fmt.Sprintf("%s:%s", tokenID, roleName)

	d.SetId(id)

	return resourceConsulACLTokenRoleAttachmentRead(d, meta)
}

func resourceConsulACLTokenRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	id := d.Id()

	tokenID, roleName, err := parseTwoPartID(id, "token", "role")
	if err != nil {
		return fmt.Errorf("invalid ACL token role attachment id '%q'", id)
	}

	aclToken, _, err := client.ACL().TokenRead(tokenID, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "ACL not found") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read token '%s': %v", id, err)
	}

	roleFound := false
	for _, role := range aclToken.Roles {
		if role.Name == roleName {
			roleFound = true
			break
		}
	}
	if !roleFound {
		d.SetId("")
		return nil
	}

	if err = d.Set("token_id", tokenID); err != nil {
		return fmt.Errorf("error while setting 'token_id': %s", err)
	}
	if err = d.Set("role", roleName); err != nil {
		return fmt.Errorf("error while setting 'role': %s", err)
	}

	return nil
}

func resourceConsulACLTokenRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)

	id := d.Id()

	tokenID, roleName, err := parseTwoPartID(id, "token", "role")
	if err != nil {
		return fmt.Errorf("invalid ACL token role attachment id '%q'", id)
	}

	aclToken, _, err := client.ACL().TokenRead(tokenID, qOpts)
	if err != nil {
		return fmt.Errorf("token '%s' not found", tokenID)
	}

	for i, role := range aclToken.Roles {
		if role.Name == roleName {
			aclToken.Roles = append(aclToken.Roles[:i], aclToken.Roles[i+1:]...)
			break
		}
	}

	_, _, err = client.ACL().TokenUpdate(aclToken, wOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL token '%q' to set new role attachment: '%s'", tokenID, err)
	}

	return nil
}
