// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulACLTokenPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLTokenPolicyAttachmentCreate,
		Read:   resourceConsulACLTokenPolicyAttachmentRead,
		Delete: resourceConsulACLTokenPolicyAttachmentDelete,
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
			"policy": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The policy name.",
			},
		},
	}
}

func resourceConsulACLTokenPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)

	tokenID := d.Get("token_id").(string)

	aclToken, _, err := client.ACL().TokenRead(tokenID, qOpts)
	if err != nil {
		return fmt.Errorf("token '%s' not found", tokenID)
	}

	newPolicyName := d.Get("policy").(string)
	for _, iPolicy := range aclToken.Policies {
		if iPolicy.Name == newPolicyName {
			return fmt.Errorf("policy '%s' already attached to token", newPolicyName)
		}
	}

	aclToken.Policies = append(aclToken.Policies, &consulapi.ACLTokenPolicyLink{
		Name: newPolicyName,
	})

	u, _, err := client.ACL().TokenUpdate(aclToken, wOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL token '%q' to set new policy attachment: '%s'", tokenID, err)
	}

	if !aclToken.Local {
		if err := waitForACLTokenReplication(client.ACL(), qOpts, u.ModifyIndex); err != nil {
			return err
		}
	}

	id := fmt.Sprintf("%s:%s", tokenID, newPolicyName)

	d.SetId(id)

	return resourceConsulACLTokenPolicyAttachmentRead(d, meta)
}

func resourceConsulACLTokenPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	id := d.Id()

	tokenID, policyName, err := parseTwoPartID(id, "token", "policy")
	if err != nil {
		return fmt.Errorf("invalid ACL token policy attachment id '%q'", id)
	}

	aclToken, _, err := client.ACL().TokenRead(tokenID, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "ACL not found") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read token '%s': %v", id, err)
	}

	policyFound := false
	for _, iPolicy := range aclToken.Policies {
		if iPolicy.Name == policyName {
			policyFound = true
			break
		}
	}
	if !policyFound {
		d.SetId("")
		return nil
	}

	if err = d.Set("token_id", tokenID); err != nil {
		return fmt.Errorf("error while setting 'token_id': %s", err)
	}
	if err = d.Set("policy", policyName); err != nil {
		return fmt.Errorf("error while setting 'policyName': %s", err)
	}

	return nil
}

func resourceConsulACLTokenPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)

	id := d.Id()

	tokenID, policyName, err := parseTwoPartID(id, "token", "policy")
	if err != nil {
		return fmt.Errorf("invalid ACL token policy attachment id '%q'", id)
	}

	aclToken, _, err := client.ACL().TokenRead(tokenID, qOpts)
	if err != nil {
		return fmt.Errorf("token '%s' not found", tokenID)
	}

	for i, iPolicy := range aclToken.Policies {
		if iPolicy.Name == policyName {
			aclToken.Policies = append(aclToken.Policies[:i], aclToken.Policies[i+1:]...)
			break
		}
	}

	_, _, err = client.ACL().TokenUpdate(aclToken, wOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL token '%q' to set new policy attachment: '%s'", tokenID, err)
	}

	return nil
}

// return the pieces of id `a:b` as a, b
func parseTwoPartID(id, resource, name string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected ID format (%q). Expected %s_id:%s_name", id, resource, name)
	}

	return parts[0], parts[1], nil
}
