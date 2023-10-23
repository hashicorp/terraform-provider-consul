// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceConsulACLRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLRoleCreate,
		Read:   resourceConsulACLRoleRead,
		Update: resourceConsulACLRoleUpdate,
		Delete: resourceConsulACLRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Description: "Starting with Consul 1.5.0, the `consul_acl_role` can be used to managed Consul ACL roles.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ACL role.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A free form human readable description of the role.",
			},
			"policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					ValidateFunc: validation.IsUUID,
					Type:         schema.TypeString,
				},
				Description: "The list of policies that should be applied to the role.",
			},
			"service_identities": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the service.",
						},

						"datacenters": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The datacenters the effective policy is valid within. When no datacenters are provided the effective policy is valid in all datacenters including those which do not yet exist but may in the future.",
						},
					},
				},
				Description: "The list of service identities that should be applied to the role.",
			},
			"node_identities": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of node identities that should be applied to the role.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the node.",
						},
						"datacenter": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the node's datacenter.",
						},
					},
				},
			},
			"templated_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of templated policies that should be applied to the token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the templated policies.",
						},
						"template_variables": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "The templated policy variables.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of node, workload identity or service.",
									},
								},
							},
						},
						"datacenters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the datacenters the effective policy is valid within.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"namespace": {
				Type:        schema.TypeString,
				Description: "The namespace to create the role within.",
				Optional:    true,
				ForceNew:    true,
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The partition the ACL role is associated with.",
			},
		},
	}
}

func resourceConsulACLRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()
	role := getRole(d, meta)

	name := role.Name
	role, _, err := ACL.RoleCreate(role, wOpts)
	if err != nil {
		return fmt.Errorf("failed to create role '%s': %s", name, err)
	}

	d.SetId(role.ID)
	return resourceConsulACLRoleRead(d, meta)
}

func resourceConsulACLRoleRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	ACL := client.ACL()

	role, _, err := ACL.RoleRead(d.Id(), qOpts)
	if err != nil {
		return fmt.Errorf("failed to read role '%s': %s", d.Id(), err)
	}
	if role == nil {
		d.SetId("")
		return nil
	}

	policies := make([]string, len(role.Policies))
	for i, policy := range role.Policies {
		policies[i] = policy.ID
	}

	serviceIdentities := make([]map[string]interface{}, len(role.ServiceIdentities))
	for i, serviceIdentity := range role.ServiceIdentities {
		serviceIdentities[i] = map[string]interface{}{
			"service_name": serviceIdentity.ServiceName,
			"datacenters":  serviceIdentity.Datacenters,
		}
	}

	nodeIdentities := make([]interface{}, len(role.NodeIdentities))
	for i, ni := range role.NodeIdentities {
		nodeIdentities[i] = map[string]interface{}{
			"node_name":  ni.NodeName,
			"datacenter": ni.Datacenter,
		}
	}

	templatedPolicies := make([]interface{}, len(role.TemplatedPolicies))
	for i, tp := range role.TemplatedPolicies {
		templatedPolicies[i] = map[string]interface{}{
			"template_name":      tp.TemplateName,
			"datacenters":        tp.Datacenters,
			"template_variables": getTemplateVariables(tp),
		}
	}

	sw := newStateWriter(d)

	sw.set("name", role.Name)
	sw.set("description", role.Description)
	sw.set("policies", policies)
	sw.set("service_identities", serviceIdentities)
	sw.set("node_identities", nodeIdentities)
	sw.set("templated_policies", templatedPolicies)
	sw.set("namespace", role.Namespace)
	sw.set("partition", role.Partition)

	return sw.error()
}

func resourceConsulACLRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()
	role := getRole(d, meta)

	role.ID = d.Id()

	role, _, err := ACL.RoleUpdate(role, wOpts)
	if err != nil {
		return fmt.Errorf("failed to update role '%s': %s", d.Id(), err)
	}

	d.SetId(role.ID)
	return resourceConsulACLRoleRead(d, meta)
}

func resourceConsulACLRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()

	if _, err := ACL.RoleDelete(d.Id(), wOpts); err != nil {
		return fmt.Errorf("failed to delete role '%s': %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getRole(d *schema.ResourceData, meta interface{}) *consulapi.ACLRole {
	_, qOpts, _ := getClient(d, meta)
	roleName := d.Get("name").(string)
	role := &consulapi.ACLRole{
		Name:        roleName,
		Description: d.Get("description").(string),
		Namespace:   qOpts.Namespace,
	}
	policies := make([]*consulapi.ACLRolePolicyLink, 0)
	for _, raw := range d.Get("policies").(*schema.Set).List() {
		policies = append(policies, &consulapi.ACLRolePolicyLink{
			ID: raw.(string),
		})
	}
	role.Policies = policies

	for _, raw := range d.Get("service_identities").(*schema.Set).List() {
		s := raw.(map[string]interface{})

		datacenters := make([]string, len(s["datacenters"].(*schema.Set).List()))
		for i, d := range s["datacenters"].(*schema.Set).List() {
			datacenters[i] = d.(string)
		}

		role.ServiceIdentities = append(role.ServiceIdentities, &consulapi.ACLServiceIdentity{
			ServiceName: s["service_name"].(string),
			Datacenters: datacenters,
		})
	}

	for _, ni := range d.Get("node_identities").([]interface{}) {
		n := ni.(map[string]interface{})
		role.NodeIdentities = append(role.NodeIdentities, &consulapi.ACLNodeIdentity{
			NodeName:   n["node_name"].(string),
			Datacenter: n["datacenter"].(string),
		})
	}

	for key, tp := range d.Get("templated_policies").([]interface{}) {
		t := tp.(map[string]interface{})

		datacenters := []string{}
		for _, d := range t["datacenters"].([]interface{}) {
			datacenters = append(datacenters, d.(string))
		}

		templatedPolicy := &consulapi.ACLTemplatedPolicy{
			Datacenters:  datacenters,
			TemplateName: t["template_name"].(string),
		}

		if templatedVariables, ok := d.GetOk(fmt.Sprint("templated_policies.", key, ".template_variables.0")); ok {
			tv := templatedVariables.(map[string]interface{})
			templatedPolicy.TemplateVariables = &consulapi.ACLTemplatedPolicyVariables{}

			if tv["name"] != nil {
				templatedPolicy.TemplateVariables.Name = tv["name"].(string)
			}
		}
		role.TemplatedPolicies = append(role.TemplatedPolicies, templatedPolicy)
	}

	return role
}
