package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulNamespaceCreate,
		Read:   resourceConsulNamespaceRead,
		Update: resourceConsulNamespaceUpdate,
		Delete: resourceConsulNamespaceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_defaults": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"role_defaults": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"meta": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceConsulNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	namespace := getNamespaceFromResourceData(d)
	namespace, _, err := client.Namespaces().Create(namespace, wOpts)
	if err != nil {
		return fmt.Errorf("failed to create namespace: %v", err)
	}
	d.SetId(namespace.Name)
	return resourceConsulNamespaceRead(d, meta)
}

func resourceConsulNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Id()

	namespace, _, err := client.Namespaces().Read(name, qOpts)
	if namespace == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to read namespace '%s': %v", name, err)
	}

	if err = d.Set("description", namespace.Description); err != nil {
		return fmt.Errorf("failed to set 'description': %v", err)
	}
	if err = d.Set("meta", namespace.Meta); err != nil {
		return fmt.Errorf("failed to set 'meta': %v", err)
	}
	roleDefaults := make([]interface{}, 0)
	for _, r := range namespace.ACLs.RoleDefaults {
		roleDefaults = append(roleDefaults, r.Name)
	}
	if err = d.Set("role_defaults", roleDefaults); err != nil {
		return fmt.Errorf("failed to set 'role_defaults': %v", err)
	}
	policyDefaults := make([]interface{}, 0)
	for _, p := range namespace.ACLs.PolicyDefaults {
		policyDefaults = append(policyDefaults, p.Name)
	}
	if err = d.Set("policy_defaults", policyDefaults); err != nil {
		return fmt.Errorf("failed to set 'policy_defaults': %v", err)
	}
	return nil
}

func resourceConsulNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	namespace := getNamespaceFromResourceData(d)
	namespace, _, err := client.Namespaces().Update(namespace, wOpts)
	if err != nil {
		return fmt.Errorf("failed to update namespace '%s': %v", namespace.Name, err)
	}

	return resourceConsulNamespaceRead(d, meta)
}

func resourceConsulNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	_, err := client.Namespaces().Delete(d.Id(), wOpts)
	if err != nil {
		return fmt.Errorf("failed to delete namespace '%s': %v", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getNamespaceFromResourceData(d *schema.ResourceData) *consulapi.Namespace {
	m := make(map[string]string)
	for name, value := range d.Get("meta").(map[string]interface{}) {
		m[name] = value.(string)
	}

	policyDefaults := make([]consulapi.ACLLink, 0)
	for _, p := range d.Get("policy_defaults").([]interface{}) {
		policyDefaults = append(policyDefaults, consulapi.ACLLink{
			Name: p.(string),
		})
	}

	roleDefaults := make([]consulapi.ACLLink, 0)
	for _, r := range d.Get("role_defaults").([]interface{}) {
		roleDefaults = append(roleDefaults, consulapi.ACLLink{
			Name: r.(string),
		})
	}

	return &consulapi.Namespace{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Meta:        m,
		ACLs: &consulapi.NamespaceACLConfig{
			PolicyDefaults: policyDefaults,
			RoleDefaults:   roleDefaults,
		},
	}
}
