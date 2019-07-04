package consul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
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
	client := getClient(meta)

	log.Printf("[DEBUG] Creating ACL token policy attachment")

	tokenID := d.Get("token_id").(string)

	aclToken, _, err := client.ACL().TokenRead(tokenID, nil)
	if err != nil {
		return fmt.Errorf("Token '%s' not found", tokenID)
	}

	newPolicyName := d.Get("policy").(string)
	for _, iPolicy := range aclToken.Policies {
		if iPolicy.Name == newPolicyName {
			return fmt.Errorf("Policy '%s' already attached to token", newPolicyName)
		}
	}

	aclToken.Policies = append(aclToken.Policies, &consulapi.ACLTokenPolicyLink{
		Name: newPolicyName,
	})

	_, _, err = client.ACL().TokenUpdate(aclToken, nil)
	if err != nil {
		return fmt.Errorf("Error updating ACL token '%q' to set new policy attachment: '%s'", tokenID, err)
	}

	id := buildTwoPartID(&tokenID, &newPolicyName)

	log.Printf("[DEBUG] Created ACL token policy attachment '%q'", id)

	d.SetId(id)

	return resourceConsulACLTokenPolicyAttachmentRead(d, meta)
}

func resourceConsulACLTokenPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL token policy attachment '%q'", id)

	tokenID, policyName, err := parseTwoPartID(id)
	if err != nil {
		return fmt.Errorf("Invalid ACL token policy attachment id '%q'", id)
	}

	aclToken, _, err := client.ACL().TokenRead(tokenID, nil)
	if err != nil {
		log.Printf("[WARN] ACL token not found, removing from state")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read ACL token %q", tokenID)

	policyFound := false
	for _, iPolicy := range aclToken.Policies {
		if iPolicy.Name == policyName {
			policyFound = true
			break
		}
	}
	if !policyFound {
		log.Printf("[WARN] ACL policy not found in token, removing from state")
		d.SetId("")
		return nil
	}

	if err = d.Set("token_id", tokenID); err != nil {
		return fmt.Errorf("Error while setting 'token_id': %s", err)
	}
	if err = d.Set("policy", policyName); err != nil {
		return fmt.Errorf("Error while setting 'policyName': %s", err)
	}

	return nil
}

func resourceConsulACLTokenPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL token policy attachment '%q'", id)

	tokenID, policyName, err := parseTwoPartID(id)
	if err != nil {
		return fmt.Errorf("Invalid ACL token policy attachment id '%q'", id)
	}

	aclToken, _, err := client.ACL().TokenRead(tokenID, nil)
	if err != nil {
		return fmt.Errorf("Token '%s' not found", tokenID)
	}

	for i, iPolicy := range aclToken.Policies {
		if iPolicy.Name == policyName {
			aclToken.Policies = append(aclToken.Policies[:i], aclToken.Policies[i+1:]...)
			break
		}
	}

	_, _, err = client.ACL().TokenUpdate(aclToken, nil)
	if err != nil {
		return fmt.Errorf("Error updating ACL token '%q' to set new policy attachment: '%s'", tokenID, err)
	}
	log.Printf("[DEBUG] Deleted ACL token attachment policy %q", id)

	return nil
}
