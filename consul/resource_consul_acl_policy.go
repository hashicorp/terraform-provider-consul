package consul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	resourceACLPolicyName        = "name"
	resourceACLPolicyDescription = "description"
	resourceACLPolicyRules       = "rules"
	resourceACLPolicyDatacenters = "datacenters"
)

func resourceConsulACLPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLPolicyCreate,
		Read:   resourceConsulACLPolicyRead,
		Update: resourceConsulACLPolicyUpdate,
		Delete: resourceConsulACLPolicyDelete,

		Schema: map[string]*schema.Schema{
			requestOptions: schemaRequestOpts,
			resourceACLPolicyName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ACL policy name.",
			},
			resourceACLPolicyDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ACL policy description.",
			},
			resourceACLPolicyRules: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ACL policy rules.",
			},
			resourceACLPolicyDatacenters: {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The ACL policy datacenters.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceConsulACLPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	// Parse out data source filters to populate Consul's write options
	writeOpts, _, err := getRequestOpts(d, client)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog nodes: {{err}}", err)
	}

	log.Printf("[DEBUG] Creating ACL policy")

	aclPolicy := consulapi.ACLPolicy{
		Name:        d.Get(resourceACLPolicyName).(string),
		Description: d.Get(resourceACLPolicyDescription).(string),
		Rules:       d.Get(resourceACLPolicyRules).(string),
	}

	if v, ok := d.GetOk(resourceACLPolicyDatacenters); ok {
		vs := v.(*schema.Set).List()
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		aclPolicy.Datacenters = s
	}

	policy, _, err := client.ACL().PolicyCreate(&aclPolicy, writeOpts)
	if err != nil {
		return fmt.Errorf("error creating ACL policy: %s", err)
	}

	log.Printf("[DEBUG] Created ACL policy %q", policy.ID)

	d.SetId(policy.ID)

	return resourceConsulACLPolicyRead(d, meta)
}

func resourceConsulACLPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	// Parse out data source filters to populate Consul's write options
	_, queryOptions, err := getRequestOpts(d, client)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog nodes: {{err}}", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL policy %q", id)

	aclPolicy, _, err := client.ACL().PolicyRead(id, queryOptions)
	if err != nil {
		log.Printf("[WARN] ACL policy not found, removing from state")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read ACL policy %q", id)

	if err = d.Set(resourceACLPolicyName, aclPolicy.Name); err != nil {
		return fmt.Errorf("Error while setting 'name': %s", err)
	}
	if err = d.Set(resourceACLPolicyDescription, aclPolicy.Description); err != nil {
		return fmt.Errorf("Error while setting 'description': %s", err)
	}
	if err = d.Set(resourceACLPolicyRules, aclPolicy.Rules); err != nil {
		return fmt.Errorf("Error while setting 'rules': %s", err)
	}

	datacenters := make([]string, 0, len(aclPolicy.Datacenters))
	for _, datacenter := range aclPolicy.Datacenters {
		datacenters = append(datacenters, datacenter)
	}

	if err = d.Set(resourceACLPolicyDatacenters, datacenters); err != nil {
		return fmt.Errorf("Error while setting 'datacenters': %s", err)
	}

	return nil
}

func resourceConsulACLPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	// Parse out data source filters to populate Consul's write options
	writeOpts, _, err := getRequestOpts(d, client)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog nodes: {{err}}", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Updating ACL policy %q", id)

	aclPolicy := consulapi.ACLPolicy{
		ID:          id,
		Name:        d.Get(resourceACLPolicyName).(string),
		Description: d.Get(resourceACLPolicyDescription).(string),
		Rules:       d.Get(resourceACLPolicyRules).(string),
	}

	if v, ok := d.GetOk(resourceACLPolicyDatacenters); ok {
		vs := v.(*schema.Set).List()
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		aclPolicy.Datacenters = s
	}

	_, _, err = client.ACL().PolicyUpdate(&aclPolicy, writeOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL policy %q: %s", id, err)
	}
	log.Printf("[DEBUG] Updated ACL policy %q", id)

	return resourceConsulACLPolicyRead(d, meta)
}

func resourceConsulACLPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	// Parse out data source filters to populate Consul's write options
	writeOpts, _, err := getRequestOpts(d, client)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog nodes: {{err}}", err)
	}

	id := d.Id()

	log.Printf("[DEBUG] Deleting ACL policy %q", id)
	_, err = client.ACL().PolicyDelete(id, writeOpts)
	if err != nil {
		return fmt.Errorf("error deleting ACL policy %q: %s", id, err)
	}
	log.Printf("[DEBUG] Deleted ACL policy %q", id)

	return nil
}
