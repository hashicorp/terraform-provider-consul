// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaQueryOpts() *schema.Schema {
	return &schema.Schema{
		Optional: true,
		Type:     schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allow_stale": {
					Optional: true,
					Default:  true,
					Type:     schema.TypeBool,
				},
				"datacenter": {
					// Optional because we'll pull the default from the local agent if it's
					// not specified, but we can query remote data centers as a result.
					Optional: true,
					Type:     schema.TypeString,
				},
				"partition": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"near": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"node_meta": {
					Optional: true,
					Type:     schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"require_consistent": {
					Optional: true,
					Default:  false,
					Type:     schema.TypeBool,
				},
				"token": {
					Optional:  true,
					Type:      schema.TypeString,
					Sensitive: true,
				},
				"wait_index": {
					Optional: true,
					Type:     schema.TypeInt,
					ValidateFunc: makeValidationFunc("wait_index", []interface{}{
						validateIntMin(0),
					}),
				},
				"wait_time": {
					Optional: true,
					Type:     schema.TypeString,
					ValidateFunc: makeValidationFunc("wait_time", []interface{}{
						validateDurationMin("0ns"),
					}),
				},
			},
		},
	}
}

func getQueryOpts(queryOpts *consulapi.QueryOptions, d *schema.ResourceData, meta interface{}) {
	if filter, ok := d.GetOk("filter"); ok {
		queryOpts.Filter = filter.(string)
	}

	if v, ok := d.GetOk("query_options"); ok {
		for _, config := range v.(*schema.Set).List() {
			queryOptions := config.(map[string]interface{})
			if v, ok := queryOptions["allow_stale"]; ok {
				queryOpts.AllowStale = v.(bool)
			}

			if v, ok := queryOptions["datacenter"]; ok {
				if v.(string) != "" {
					queryOpts.Datacenter = v.(string)
				}
			}

			if v, ok := queryOptions["partition"]; ok {
				if v.(string) != "" {
					queryOpts.Partition = v.(string)
				}
			}

			if v, ok := queryOptions["namespace"]; ok {
				queryOpts.Namespace = v.(string)
			}

			if v, ok := queryOptions["near"]; ok {
				queryOpts.Near = v.(string)
			}

			if v, ok := queryOptions["require_consistent"]; ok {
				queryOpts.RequireConsistent = v.(bool)
			}

			if v, ok := queryOptions["node_meta"]; ok {
				m := v.(map[string]interface{})
				nodeMetaMap := make(map[string]string, len("node_meta"))
				for s, t := range m {
					nodeMetaMap[s] = t.(string)
				}
				queryOpts.NodeMeta = nodeMetaMap
			}

			if v, ok := queryOptions["token"]; ok {
				queryOpts.Token = v.(string)
			}

			if v, ok := queryOptions["wait_index"]; ok {
				queryOpts.WaitIndex = uint64(v.(int))
			}

			if v, ok := queryOptions["wait_time"]; ok {
				d, _ := time.ParseDuration(v.(string))
				queryOpts.WaitTime = d
			}
		}
	}
}
