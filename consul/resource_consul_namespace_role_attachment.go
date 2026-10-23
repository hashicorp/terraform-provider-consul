// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConsulNamespaceRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulNamespaceRoleAttachmentCreate,
		Read:   resourceConsulNamespaceRoleAttachmentRead,
		Delete: resourceConsulNamespaceRoleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The namespace to attach the role to.",
			},
			"role": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The role name.",
			},
		},
	}
}

func resourceConsulNamespaceRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)

	name := d.Get("namespace").(string)
	role := d.Get("role").(string)

	namespace, err := findNamespace(client, qOpts, name)
	if err != nil {
		return err
	}

	for _, r := range namespace.ACLs.RoleDefaults {
		if r.Name == role {
			return fmt.Errorf("role %q already attached to the namespace", role)
		}
	}

	namespace.ACLs.RoleDefaults = append(namespace.ACLs.RoleDefaults, consulapi.ACLLink{
		Name: role,
	})

	_, _, err = client.Namespaces().Update(namespace, wOpts)
	if err != nil {
		return fmt.Errorf("failed to update namespace %q to attach role %q: %s", name, role, err)
	}

	d.SetId(fmt.Sprintf("%s:%s", name, role))

	return resourceConsulNamespaceRoleAttachmentRead(d, meta)
}

func resourceConsulNamespaceRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	name, role, err := parseTwoPartID(d.Id(), "namespace", "role")
	if err != nil {
		return err
	}

	namespace, _, err := client.Namespaces().Read(name, qOpts)
	if err != nil {
		return fmt.Errorf("failed to read namespace %q: %s", name, err)
	}
	if namespace == nil {
		d.SetId("")
		return nil
	}

	var found bool
	for _, l := range namespace.ACLs.RoleDefaults {
		if l.Name == role {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	if err = d.Set("namespace", name); err != nil {
		return fmt.Errorf("failed to set 'namespace': %s", err)
	}
	if err = d.Set("role", role); err != nil {
		return fmt.Errorf("failed to set 'role': %s", err)
	}

	return nil
}

func resourceConsulNamespaceRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	name, role, err := parseTwoPartID(d.Id(), "namespace", "role")
	if err != nil {
		return err
	}

	namespace, err := findNamespace(client, qOpts, name)
	if err != nil {
		return err
	}

	for i, p := range namespace.ACLs.RoleDefaults {
		if p.Name == role {
			namespace.ACLs.RoleDefaults = append(
				namespace.ACLs.RoleDefaults[:i],
				namespace.ACLs.RoleDefaults[i+1:]...,
			)
			break
		}
	}

	_, _, err = client.Namespaces().Update(namespace, wOpts)
	if err != nil {
		return fmt.Errorf("failed to remove role %q from namespace %q", role, name)
	}

	return nil
}

func findNamespace(client *consulapi.Client, qOpts *consulapi.QueryOptions, name string) (*consulapi.Namespace, error) {
	namespace, _, err := client.Namespaces().Read(name, qOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to read namespace %q: %s", name, err)
	}

	if namespace == nil {
		return nil, fmt.Errorf("namespace %q not found", name)
	}

	return namespace, nil
}
