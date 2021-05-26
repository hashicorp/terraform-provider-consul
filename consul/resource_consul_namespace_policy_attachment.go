package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulNamespacePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulNamespacePolicyAttachmentCreate,
		Read:   resourceConsulNamespacePolicyAttachmentRead,
		Delete: resourceConsulNamespacePolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The namespace to attach the policy to.",
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

func resourceConsulNamespacePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)

	name := d.Get("namespace").(string)
	namespace, _, err := client.Namespaces().Read(name, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to read namespace %q: %s", name, err)
	}

	if namespace == nil {
		return fmt.Errorf("Namespace %q not found", name)
	}

	policy := d.Get("policy").(string)
	for _, p := range namespace.ACLs.PolicyDefaults {
		if p.Name == policy {
			return fmt.Errorf("Policy %q already attached to the namespace", policy)
		}
	}

	namespace.ACLs.PolicyDefaults = append(namespace.ACLs.PolicyDefaults, consulapi.ACLLink{
		Name: policy,
	})

	_, _, err = client.Namespaces().Update(namespace, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to update namespace %q to attach policy %q", name, policy)
	}

	d.SetId(fmt.Sprintf("%s:%s", name, policy))

	return resourceConsulNamespacePolicyAttachmentRead(d, meta)
}

func resourceConsulNamespacePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	name, policy, err := parseTwoPartID(d.Id(), "namespace", "policy")
	if err != nil {
		return err
	}

	namespace, _, err := client.Namespaces().Read(name, qOpts)
	if namespace == nil {
		d.SetId("")
		return nil
	}

	var found bool
	for _, l := range namespace.ACLs.PolicyDefaults {
		if l.Name == policy {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	if err = d.Set("namespace", name); err != nil {
		return fmt.Errorf("Failed to set 'namespace': %s", err)
	}
	if err = d.Set("policy", policy); err != nil {
		return fmt.Errorf("Failed to set 'policy': %s", err)
	}

	return nil
}

func resourceConsulNamespacePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	name, policy, err := parseTwoPartID(d.Id(), "namespace", "policy")
	if err != nil {
		return err
	}

	namespace, _, err := client.Namespaces().Read(name, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to get namespace %q: %s", name, err)
	}
	if namespace == nil {
		return fmt.Errorf("Namespace %q not found", name)
	}

	for i, p := range namespace.ACLs.PolicyDefaults {
		if p.Name == policy {
			namespace.ACLs.PolicyDefaults = append(
				namespace.ACLs.PolicyDefaults[:i],
				namespace.ACLs.PolicyDefaults[i+1:]...,
			)
			break
		}
	}

	_, _, err = client.Namespaces().Update(namespace, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to remove policy %q from namespace %q", policy, name)
	}

	return nil
}
