package consul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	resourceACLTokenDescription = "description"
	resourceACLTokenPolicies    = "policies"
	resourceACLTokenLocal       = "local"

	resourceACLTokenSecret = "secret"
)

func resourceConsulACLToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLTokenCreate,
		Read:   resourceConsulACLTokenRead,
		Update: resourceConsulACLTokenUpdate,
		Delete: resourceConsulACLTokenDelete,

		Schema: map[string]*schema.Schema{
			requestOptions: schemaRequestOpts,
			resourceACLTokenDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The token description.",
			},
			resourceACLTokenPolicies: {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of policies.",
			},
			resourceACLTokenLocal: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag to set the token local to the current datacenter.",
			},
			resourceACLTokenSecret: {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceConsulACLTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	// Parse out data source filters to populate Consul's write options
	writeOpts, _, err := getRequestOpts(d, client)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog nodes: {{err}}", err)
	}

	log.Printf("[DEBUG] Creating ACL token")

	aclToken := consulapi.ACLToken{
		Description: d.Get(resourceACLTokenDescription).(string),
		Local:       d.Get(resourceACLTokenLocal).(bool),
	}

	iPolicies := d.Get(resourceACLTokenPolicies).(*schema.Set).List()
	policyLinks := make([]*consulapi.ACLTokenPolicyLink, 0, len(iPolicies))
	for _, iPolicy := range iPolicies {
		policyLinks = append(policyLinks, &consulapi.ACLTokenPolicyLink{
			Name: iPolicy.(string),
		})
	}

	if len(policyLinks) > 0 {
		aclToken.Policies = policyLinks
	}

	token, _, err := client.ACL().TokenCreate(&aclToken, writeOpts)
	if err != nil {
		return fmt.Errorf("error creating ACL token: %s", err)
	}

	if err = d.Set(resourceACLTokenSecret, token.SecretID); err != nil {
		return fmt.Errorf("Error while setting 'secret': %s", err)
	}

	log.Printf("[DEBUG] Created ACL token %q", token.AccessorID)

	d.SetId(token.AccessorID)

	return resourceConsulACLTokenRead(d, meta)
}

func resourceConsulACLTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	// Parse out data source filters to populate Consul's write options
	_, queryOptions, err := getRequestOpts(d, client)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog nodes: {{err}}", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL token %q", id)

	aclToken, _, err := client.ACL().TokenRead(id, queryOptions)
	if err != nil {
		log.Printf("[WARN] ACL token not found, removing from state")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read ACL token %q", id)

	if err = d.Set(resourceACLTokenDescription, aclToken.Description); err != nil {
		return fmt.Errorf("Error while setting 'description': %s", err)
	}

	policies := make([]string, 0, len(aclToken.Policies))
	for _, policyLink := range aclToken.Policies {
		policies = append(policies, policyLink.Name)
	}

	if err = d.Set(resourceACLTokenPolicies, policies); err != nil {
		return fmt.Errorf("Error while setting 'policies': %s", err)
	}
	if err = d.Set(resourceACLTokenLocal, aclToken.Local); err != nil {
		return fmt.Errorf("Error while setting 'local': %s", err)
	}
	if err = d.Set(resourceACLTokenSecret, aclToken.SecretID); err != nil {
		return fmt.Errorf("Error while setting 'secret': %s", err)
	}

	return nil
}

func resourceConsulACLTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	// Parse out data source filters to populate Consul's write options
	writeOpts, _, err := getRequestOpts(d, client)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog nodes: {{err}}", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Updating ACL token %q", id)

	aclToken := consulapi.ACLToken{
		AccessorID:  id,
		Description: d.Get(resourceACLTokenDescription).(string),
		Local:       d.Get(resourceACLTokenLocal).(bool),
	}

	if v, ok := d.GetOk(resourceACLTokenPolicies); ok {
		vs := v.([]interface{})
		s := make([]*consulapi.ACLTokenPolicyLink, len(vs))
		for i, raw := range vs {
			s[i].Name = raw.(string)
		}
		aclToken.Policies = s
	}

	token, _, err := client.ACL().TokenUpdate(&aclToken, writeOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL token %q: %s", id, err)
	}

	if err = d.Set(resourceACLTokenSecret, token.SecretID); err != nil {
		return fmt.Errorf("Error while setting 'secret': %s", err)
	}

	log.Printf("[DEBUG] Updated ACL token %q", id)

	return resourceConsulACLTokenRead(d, meta)
}

func resourceConsulACLTokenDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	// Parse out data source filters to populate Consul's write options
	writeOpts, _, err := getRequestOpts(d, client)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog nodes: {{err}}", err)
	}

	id := d.Id()

	log.Printf("[DEBUG] Deleting ACL token %q", id)
	_, err = client.ACL().TokenDelete(id, writeOpts)
	if err != nil {
		return fmt.Errorf("error deleting ACL token %q: %s", id, err)
	}
	log.Printf("[DEBUG] Deleted ACL token %q", id)

	return nil
}
