package consul

import (
	"fmt"
	"log"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulACLToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLTokenCreate,
		Read:   resourceConsulACLTokenRead,
		Update: resourceConsulACLTokenUpdate,
		Delete: resourceConsulACLTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"accessor_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Optional:    true,
				Description: "The token id.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The token description.",
			},
			"policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of policies.",
			},
			"roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of roles",
			},
			"local": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     false,
				Description: "Flag to set the token local to the current datacenter.",
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceConsulACLTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	log.Printf("[DEBUG] Creating ACL token")

	aclToken := getToken(d)

	token, _, err := client.ACL().TokenCreate(aclToken, wOpts)
	if err != nil {
		return fmt.Errorf("error creating ACL token: %s", err)
	}

	log.Printf("[DEBUG] Created ACL token %q", token.AccessorID)

	d.SetId(token.AccessorID)

	return resourceConsulACLTokenRead(d, meta)
}

func resourceConsulACLTokenRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL token %q", id)

	aclToken, _, err := client.ACL().TokenRead(id, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "ACL not found") {
			log.Printf("[WARN] ACL token not found, removing from state")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Failed to read token '%s': %v", id, err)
	}

	log.Printf("[DEBUG] Read ACL token %q", id)

	if err = d.Set("accessor_id", aclToken.AccessorID); err != nil {
		return fmt.Errorf("Error while setting 'accessor_id': %s", err)
	}

	if err = d.Set("description", aclToken.Description); err != nil {
		return fmt.Errorf("Error while setting 'description': %s", err)
	}

	policies := make([]string, 0, len(aclToken.Policies))
	for _, policyLink := range aclToken.Policies {
		policies = append(policies, policyLink.Name)
	}

	if err = d.Set("policies", policies); err != nil {
		return fmt.Errorf("Error while setting 'policies': %s", err)
	}

	roles := make([]string, 0, len(aclToken.Roles))
	for _, roleLink := range aclToken.Roles {
		roles = append(roles, roleLink.Name)
	}

	if err = d.Set("roles", roles); err != nil {
		return fmt.Errorf("Error while setting 'roles': %s", err)
	}

	if err = d.Set("local", aclToken.Local); err != nil {
		return fmt.Errorf("Error while setting 'local': %s", err)
	}

	return nil
}

func resourceConsulACLTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	id := d.Id()
	log.Printf("[DEBUG] Updating ACL token %q", id)

	aclToken := getToken(d)
	aclToken.AccessorID = id

	_, _, err := client.ACL().TokenUpdate(aclToken, wOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL token %q: %s", id, err)
	}
	log.Printf("[DEBUG] Updated ACL token %q", id)

	return resourceConsulACLTokenRead(d, meta)
}

func resourceConsulACLTokenDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	id := d.Id()

	log.Printf("[DEBUG] Deleting ACL token %q", id)
	_, err := client.ACL().TokenDelete(id, wOpts)
	if err != nil {
		return fmt.Errorf("error deleting ACL token %q: %s", id, err)
	}
	log.Printf("[DEBUG] Deleted ACL token %q", id)

	return nil
}

func getToken(d *schema.ResourceData) *consulapi.ACLToken {
	aclToken := &consulapi.ACLToken{
		AccessorID:  d.Get("accessor_id").(string),
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
	aclToken.Policies = policyLinks

	iRoles := d.Get("roles").(*schema.Set).List()
	roleLinks := make([]*consulapi.ACLTokenRoleLink, 0, len(iRoles))
	for _, iRole := range iRoles {
		roleLinks = append(roleLinks, &consulapi.ACLTokenRoleLink{
			Name: iRole.(string),
		})
	}
	aclToken.Roles = roleLinks

	return aclToken
}
