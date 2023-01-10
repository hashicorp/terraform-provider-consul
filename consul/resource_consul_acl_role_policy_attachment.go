package consul

import (
	"fmt"
	"strings"

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

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the role.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The policy name.",
			},
		},
	}
}

func resourceConsulACLRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)

	roleID := d.Get("role_id").(string)

	aclRole, _, err := client.ACL().RoleRead(roleID, qOpts)
	if err != nil {
		return fmt.Errorf("role '%s' not found", roleID)
	}

	policyID := d.Get("policy_id").(string)
	for _, policy := range aclRole.Policies {
		if policy.ID == policyID {
			return fmt.Errorf("policy '%s' already attached to role", policyID)
		}
	}

	aclRole.Policies = append(aclRole.Policies, &consulapi.ACLRolePolicyLink{
		ID: policyID,
	})

	_, _, err = client.ACL().RoleUpdate(aclRole, wOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL role '%q' to set new policy attachment: '%s'", roleID, err)
	}

	id := fmt.Sprintf("%s:%s", roleID, policyID)

	d.SetId(id)

	return resourceConsulACLRolePolicyAttachmentRead(d, meta)
}

func resourceConsulACLRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	id := d.Id()

	roleID, policyID, err := parseTwoPartID(id, "role_id", "policy_id")
	if err != nil {
		return fmt.Errorf("invalid ACL role policy attachment id '%q'", id)
	}

	aclRole, _, err := client.ACL().RoleRead(roleID, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "ACL not found") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read role '%s': %v", id, err)
	}

	policyFound := false
	for _, policy := range aclRole.Policies {
		if policy.ID == policyID {
			policyFound = true
			break
		}
	}
	if !policyFound {
		d.SetId("")
		return nil
	}

	if err = d.Set("role_id", roleID); err != nil {
		return fmt.Errorf("error while setting 'role': %s", err)
	}
	if err = d.Set("policy_id", policyID); err != nil {
		return fmt.Errorf("error while setting 'policy': %s", err)
	}

	return nil
}

func resourceConsulACLRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)

	id := d.Id()

	roleID, policyID, err := parseTwoPartID(id, "role_id", "policy_id")
	if err != nil {
		return fmt.Errorf("invalid ACL role policy attachment id '%q'", id)
	}

	aclRole, _, err := client.ACL().RoleRead(roleID, qOpts)
	if err != nil {
		return fmt.Errorf("role '%s' not found", roleID)
	}

	for i, policy := range aclRole.Policies {
		if policy.ID == policyID {
			aclRole.Policies = append(aclRole.Policies[:i], aclRole.Policies[i+1:]...)
			break
		}
	}

	_, _, err = client.ACL().RoleUpdate(aclRole, wOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL role '%q' to set new policy attachment: '%s'", roleID, err)
	}

	return nil
}
