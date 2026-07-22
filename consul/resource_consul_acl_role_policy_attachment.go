// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulACLRolePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLRolePolicyAttachmentCreate,
		Read:   resourceConsulACLRolePolicyAttachmentRead,
		Delete: resourceConsulACLRolePolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Description: "The `consul_acl_role_policy_attachment` resource links a Consul ACL role and an ACL policy. The link is implemented through an update to the Consul ACL role.\n\n~> **NOTE:** This resource is only useful to attach policies to an ACL role that has been created outside the current Terraform configuration. If the ACL role you need to attach a policy to has been created in the current Terraform configuration and will only be used in it, you should use the `policies` attribute of [`consul_acl_role`](/docs/providers/consul/r/acl_role.html).",

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The id of the role.",
			},
			"policy": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The policy name.",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The partition the ACL policy and role are associated with.",
			},
		},
	}
}

func resourceConsulACLRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)

	roleID := d.Get("role_id").(string)

	role, _, err := client.ACL().RoleRead(roleID, qOpts)
	if err != nil {
		return fmt.Errorf("failed to find role %q: %w", roleID, err)
	}
	if role == nil {
		return fmt.Errorf("role %q not found", roleID)
	}

	newPolicyName := d.Get("policy").(string)
	for _, iPolicy := range role.Policies {
		if iPolicy.Name == newPolicyName {
			return fmt.Errorf("policy '%s' already attached to role", newPolicyName)
		}
	}

	role.Policies = append(role.Policies, &consulapi.ACLRolePolicyLink{
		Name: newPolicyName,
	})

	_, _, err = client.ACL().RoleUpdate(role, wOpts)
	if err != nil {
		return fmt.Errorf("error updating role '%q' to set new policy attachment: '%s'", roleID, err)
	}

	id := fmt.Sprintf("%s:%s", roleID, newPolicyName)

	d.SetId(id)

	return resourceConsulACLRolePolicyAttachmentRead(d, meta)
}

func resourceConsulACLRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	id := d.Id()

	roleID, policyName, err := parseTwoPartID(id, "role", "policy")
	if err != nil {
		return fmt.Errorf("invalid role policy attachment id '%q'", id)
	}

	role, _, err := client.ACL().RoleRead(roleID, qOpts)
	if err != nil {
		return fmt.Errorf("failed to read token '%s': %v", id, err)
	}
	if role == nil {
		d.SetId("")
		return nil
	}

	policyFound := false
	for _, iPolicy := range role.Policies {
		if iPolicy.Name == policyName {
			policyFound = true
			break
		}
	}
	if !policyFound {
		d.SetId("")
		return nil
	}

	sw := newStateWriter(d)
	sw.set("role_id", roleID)
	sw.set("policy", policyName)

	return sw.error()
}

func resourceConsulACLRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)

	id := d.Id()

	roleID, policyName, err := parseTwoPartID(id, "role", "policy")
	if err != nil {
		return fmt.Errorf("invalid role policy attachment id '%q'", id)
	}

	role, _, err := client.ACL().RoleRead(roleID, qOpts)
	if err != nil {
		return fmt.Errorf("role '%s' not found", roleID)
	}
	if role == nil {
		// If the role does not exist there is no policy attachment to remove
		return nil
	}

	for i, iPolicy := range role.Policies {
		if iPolicy.Name == policyName {
			role.Policies = append(role.Policies[:i], role.Policies[i+1:]...)
			break
		}
	}

	_, _, err = client.ACL().RoleUpdate(role, wOpts)
	if err != nil {
		return fmt.Errorf("error updating role '%q' to remove policy attachment: '%s'", roleID, err)
	}

	return nil
}
