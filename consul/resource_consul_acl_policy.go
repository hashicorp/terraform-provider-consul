// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The partition the ACL policy is associated with.",
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

	aclPolicy, _, err := client.ACL().PolicyRead(id, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "ACL not found") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read policy '%s': %v", id, err)
	}

	sw := newStateWriter(d)
	sw.set("name", aclPolicy.Name)
	sw.set("description", aclPolicy.Description)
	sw.set("rules", aclPolicy.Rules)
	sw.set("datacenters", aclPolicy.Datacenters)
	sw.set("namespace", aclPolicy.Namespace)
	sw.set("partition", aclPolicy.Partition)

	return sw.error()
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
