package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceConsulACLPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLPolicyCreate,
		Read:   resourceConsulACLPolicyRead,
		Update: resourceConsulACLPolicyUpdate,
		Delete: resourceConsulACLPolicyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: "The ACL policy name.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The ACL policy description.",
			},
			"rules": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: "The ACL policy rules.",
			},
			"datacenters": {
				Type:        schema.TypeSet,
				Required:    false,
				Optional:    true,
				Description: "The ACL policy datacenters.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceConsulACLPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	log.Printf("[DEBUG] Creating ACL policy")

	aclPolicy := consulapi.ACLPolicy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       d.Get("rules").(string),
	}

	if v, ok := d.GetOk("datacenters"); ok {
		vs := v.(*schema.Set).List()
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		aclPolicy.Datacenters = s
	}

	policy, _, err := client.ACL().PolicyCreate(&aclPolicy, nil)
	if err != nil {
		return fmt.Errorf("error creating ACL policy: %s", err)
	}

	log.Printf("[DEBUG] Created ACL policy %q", policy.ID)

	d.SetId(policy.ID)

	return resourceConsulACLPolicyRead(d, meta)
}

func resourceConsulACLPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL policy %q", id)

	aclPolicy, _, err := client.ACL().PolicyRead(id, nil)
	if err != nil {
		log.Printf("[WARN] ACL policy not found, removing from state")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read ACL policy %q", id)

	d.Set("name", aclPolicy.Name)
	d.Set("description", aclPolicy.Description)
	d.Set("rules", aclPolicy.Rules)

	datacenters := make([]string, 0, len(aclPolicy.Datacenters))
	for _, datacenter := range aclPolicy.Datacenters {
		datacenters = append(datacenters, datacenter)
	}

	d.Set("datacenters", datacenters)

	return nil
}

func resourceConsulACLPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()
	log.Printf("[DEBUG] Updating ACL policy %q", id)

	aclPolicy := consulapi.ACLPolicy{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       d.Get("rules").(string),
	}

	if v, ok := d.GetOk("datacenters"); ok {
		vs := v.(*schema.Set).List()
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		aclPolicy.Datacenters = s
	}

	_, _, err := client.ACL().PolicyUpdate(&aclPolicy, nil)
	if err != nil {
		return fmt.Errorf("error updating ACL policy %q: %s", id, err)
	}
	log.Printf("[DEBUG] Updated ACL policy %q", id)

	return resourceConsulACLPolicyRead(d, meta)
}

func resourceConsulACLPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	id := d.Id()

	log.Printf("[DEBUG] Deleting ACL policy %q", id)
	_, err := client.ACL().PolicyDelete(id, nil)
	if err != nil {
		return fmt.Errorf("error deleting ACL policy %q: %s", id, err)
	}
	log.Printf("[DEBUG] Deleted ACL policy %q", id)

	return nil
}
