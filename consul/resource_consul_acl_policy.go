package consul

import (
	"fmt"
	"log"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulACLPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLPolicyCreate,
		Read:   resourceConsulACLPolicyRead,
		Update: resourceConsulACLPolicyUpdate,
		Delete: resourceConsulACLPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ACL policy name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ACL policy description.",
			},
			"rules": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ACL policy rules.",
			},
			"datacenters": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The ACL policy datacenters.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceConsulACLPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	log.Printf("[DEBUG] Creating ACL policy")

	aclPolicy := consulapi.ACLPolicy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       d.Get("rules").(string),
		Namespace:   wOpts.Namespace,
	}

	if v, ok := d.GetOk("datacenters"); ok {
		vs := v.(*schema.Set).List()
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		aclPolicy.Datacenters = s
	}

	policy, _, err := client.ACL().PolicyCreate(&aclPolicy, wOpts)
	if err != nil {
		return fmt.Errorf("error creating ACL policy: %s", err)
	}

	log.Printf("[DEBUG] Created ACL policy %q", policy.ID)

	d.SetId(policy.ID)

	return resourceConsulACLPolicyRead(d, meta)
}

func resourceConsulACLPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL policy %q", id)

	aclPolicy, _, err := client.ACL().PolicyRead(id, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "ACL not found") {
			log.Printf("[INFO] ACL policy not found, removing from state")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Failed to read policy '%s': %v", id, err)
	}

	log.Printf("[DEBUG] Read ACL policy %q", id)

	if err = d.Set("name", aclPolicy.Name); err != nil {
		return fmt.Errorf("Error while setting 'name': %s", err)
	}
	if err = d.Set("description", aclPolicy.Description); err != nil {
		return fmt.Errorf("Error while setting 'description': %s", err)
	}
	if err = d.Set("rules", aclPolicy.Rules); err != nil {
		return fmt.Errorf("Error while setting 'rules': %s", err)
	}

	datacenters := make([]string, 0, len(aclPolicy.Datacenters))
	for _, datacenter := range aclPolicy.Datacenters {
		datacenters = append(datacenters, datacenter)
	}

	if err = d.Set("datacenters", datacenters); err != nil {
		return fmt.Errorf("Error while setting 'datacenters': %s", err)
	}

	return nil
}

func resourceConsulACLPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	id := d.Id()
	log.Printf("[DEBUG] Updating ACL policy %q", id)

	aclPolicy := consulapi.ACLPolicy{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       d.Get("rules").(string),
		Namespace:   wOpts.Namespace,
	}

	if v, ok := d.GetOk("datacenters"); ok {
		vs := v.(*schema.Set).List()
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		aclPolicy.Datacenters = s
	}

	_, _, err := client.ACL().PolicyUpdate(&aclPolicy, wOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL policy %q: %s", id, err)
	}
	log.Printf("[DEBUG] Updated ACL policy %q", id)

	return resourceConsulACLPolicyRead(d, meta)
}

func resourceConsulACLPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	id := d.Id()

	log.Printf("[DEBUG] Deleting ACL policy %q", id)
	_, err := client.ACL().PolicyDelete(id, wOpts)
	if err != nil {
		return fmt.Errorf("error deleting ACL policy %q: %s", id, err)
	}
	log.Printf("[DEBUG] Deleted ACL policy %q", id)

	return nil
}
