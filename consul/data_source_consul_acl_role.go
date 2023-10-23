// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulACLRole() *schema.Resource {
	return &schema.Resource{
		Read: datasourceConsulACLRoleRead,

		Description: "The `consul_acl_role` data source returns the information related to a [Consul ACL Role](https://www.consul.io/api/acl/roles.html).",

		Schema: map[string]*schema.Schema{
			// Filters
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the ACL Role.",
				Required:    true,
			},
			"namespace": {
				Type:        schema.TypeString,
				Description: "The namespace to lookup the role.",
				Optional:    true,
			},
			"partition": {
				Type:        schema.TypeString,
				Description: "The partition to lookup the role.",
				Optional:    true,
			},

			// Out parameters
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the ACL Role.",
				Computed:    true,
			},
			"policies": {
				Type:        schema.TypeList,
				Description: "The list of policies associated with the ACL Role.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "The name of the policy.",
						},
						"id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "The ID of the policy.",
						},
					},
				},
			},
			"service_identities": {
				Type:        schema.TypeList,
				Description: "The list of service identities associated with the ACL Role.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Description: "The name of the service.",
							Optional:    true,
						},

						"datacenters": {
							Type:        schema.TypeSet,
							Description: "Specifies the datacenters the effective policy is valid within.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"node_identities": {
				Type:        schema.TypeList,
				Description: "The list of node identities associated with the ACL Role.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:        schema.TypeString,
							Description: "The name of the node.",
							Computed:    true,
						},
						"datacenter": {
							Type:        schema.TypeString,
							Description: "Specifies the nodes datacenter. This will result in effective policy only being valid in that datacenter.",
							Computed:    true,
						},
					},
				},
			},
			"templated_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of templated policies that should be applied to the token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the templated policies.",
						},
						"template_variables": {
							Type:        schema.TypeList,
							Description: "The templated policy variables.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of node, workload identity or service.",
									},
								},
							},
						},
						"datacenters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the datacenters the effective policy is valid within.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func datasourceConsulACLRoleRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Get("name").(string)

	role, _, err := client.ACL().RoleReadByName(name, qOpts)
	if err != nil {
		return fmt.Errorf("failed to get role: %v", err)
	}
	if role == nil {
		return fmt.Errorf("could not find role '%s'", name)
	}

	policies := make([]map[string]interface{}, len(role.Policies))
	for i, p := range role.Policies {
		policies[i] = map[string]interface{}{
			"name": p.Name,
			"id":   p.ID,
		}
	}

	identities := make([]map[string]interface{}, len(role.ServiceIdentities))
	for i, si := range role.ServiceIdentities {
		identities[i] = map[string]interface{}{
			"service_name": si.ServiceName,
			"datacenters":  si.Datacenters,
		}
	}

	nodeIdentities := make([]interface{}, len(role.NodeIdentities))
	for i, ni := range role.NodeIdentities {
		nodeIdentities[i] = map[string]interface{}{
			"node_name":  ni.NodeName,
			"datacenter": ni.Datacenter,
		}
	}

	templatedPolicies := make([]map[string]interface{}, len(role.TemplatedPolicies))
	for i, tp := range role.TemplatedPolicies {
		templatedPolicies[i] = map[string]interface{}{
			"template_name":      tp.TemplateName,
			"datacenters":        tp.Datacenters,
			"template_variables": getTemplateVariables(tp),
		}
	}

	d.SetId(role.ID)

	sw := newStateWriter(d)
	sw.set("description", role.Description)
	sw.set("policies", policies)
	sw.set("service_identities", identities)
	sw.set("node_identities", nodeIdentities)
	sw.set("templated_policies", templatedPolicies)

	return sw.error()
}
