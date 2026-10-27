// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConsulACLToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulACLTokenRead,

		Description: "The `consul_acl_token` data source returns the information related to the `consul_acl_token` resource with the exception of its secret ID.\n\nIf you want to get the secret ID associated with a token, use the [`consul_acl_token_secret_id` data source](/docs/providers/consul/d/acl_token_secret_id.html).",

		Schema: map[string]*schema.Schema{

			// Filters
			"accessor_id": {
				Required:    true,
				Description: "The accessor ID of the ACL token.",
				Type:        schema.TypeString,
			},
			"namespace": {
				Type:        schema.TypeString,
				Description: "The namespace to lookup the ACL token.",
				Optional:    true,
			},
			"partition": {
				Type:        schema.TypeString,
				Description: "The partition to lookup the ACL token.",
				Optional:    true,
			},

			// Out parameters
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the ACL token.",
				Computed:    true,
			},
			"policies": {
				Type:        schema.TypeList,
				Description: "A list of policies associated with the ACL token.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"roles": {
				Type:        schema.TypeList,
				Description: "List of roles linked to the token",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"service_identities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of service identities attached to the token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the service.",
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
			"node_identities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of node identities attached to the token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The list of node identities that should be applied to the token.",
						},
						"datacenter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the node's datacenter.",
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
			"local": {
				Type:        schema.TypeBool,
				Description: "Whether the ACL token is local to the datacenter it was created within.",
				Computed:    true,
			},
			"expiration_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set this represents the point after which a token should be considered revoked and is eligible for destruction.",
			},
		},
	}
}

func dataSourceConsulACLTokenRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	accessorID := d.Get("accessor_id").(string)

	aclToken, _, err := client.ACL().TokenRead(accessorID, qOpts)
	if err != nil {
		return err
	}

	policies := make([]map[string]interface{}, len(aclToken.Policies))
	for i, policyLink := range aclToken.Policies {
		policies[i] = map[string]interface{}{
			"name": policyLink.Name,
			"id":   policyLink.ID,
		}
	}

	roles := make([]interface{}, len(aclToken.Roles))
	for i, r := range aclToken.Roles {
		roles[i] = map[string]interface{}{
			"id":   r.ID,
			"name": r.Name,
		}
	}

	serviceIdentities := make([]map[string]interface{}, len(aclToken.ServiceIdentities))
	for i, si := range aclToken.ServiceIdentities {
		serviceIdentities[i] = map[string]interface{}{
			"service_name": si.ServiceName,
			"datacenters":  si.Datacenters,
		}
	}

	nodeIdentities := make([]map[string]interface{}, len(aclToken.NodeIdentities))
	for i, ni := range aclToken.NodeIdentities {
		nodeIdentities[i] = map[string]interface{}{
			"node_name":  ni.NodeName,
			"datacenter": ni.Datacenter,
		}
	}

	templatedPolicies := make([]map[string]interface{}, len(aclToken.TemplatedPolicies))
	for i, tp := range aclToken.TemplatedPolicies {
		templatedPolicies[i] = map[string]interface{}{
			"template_name":      tp.TemplateName,
			"datacenters":        tp.Datacenters,
			"template_variables": getTemplateVariables(tp),
		}
	}

	var expirationTime string
	if aclToken.ExpirationTime != nil {
		expirationTime = aclToken.ExpirationTime.Format(time.RFC3339)
	}

	d.SetId(accessorID)

	sw := newStateWriter(d)
	sw.set("description", aclToken.Description)
	sw.set("local", aclToken.Local)
	sw.set("policies", policies)
	sw.set("roles", roles)
	sw.set("service_identities", serviceIdentities)
	sw.set("node_identities", nodeIdentities)
	sw.set("templated_policies", templatedPolicies)
	sw.set("expiration_time", expirationTime)

	return sw.error()
}
