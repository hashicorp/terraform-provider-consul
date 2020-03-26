package consul

import (
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	queryOptAllowStale        = "allow_stale"
	queryOptDatacenter        = "datacenter"
	queryOptNear              = "near"
	queryOptNodeMeta          = "node_meta"
	queryOptRequireConsistent = "require_consistent"
	queryOptToken             = "token"
	queryOptWaitIndex         = "wait_index"
	queryOptWaitTime          = "wait_time"
)

func schemaQueryOpts() *schema.Schema {
	return &schema.Schema{
		Optional: true,
		Type:     schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				queryOptAllowStale: {
					Optional: true,
					Default:  true,
					Type:     schema.TypeBool,
				},
				queryOptDatacenter: {
					// Optional because we'll pull the default from the local agent if it's
					// not specified, but we can query remote data centers as a result.
					Optional: true,
					Type:     schema.TypeString,
				},
				queryOptNear: {
					Optional: true,
					Type:     schema.TypeString,
				},
				queryOptNodeMeta: {
					Optional: true,
					Type:     schema.TypeMap,
				},
				queryOptRequireConsistent: {
					Optional: true,
					Default:  false,
					Type:     schema.TypeBool,
				},
				queryOptToken: {
					Optional:  true,
					Type:      schema.TypeString,
					Sensitive: true,
				},
				queryOptWaitIndex: {
					Optional: true,
					Type:     schema.TypeInt,
					ValidateFunc: makeValidationFunc(queryOptWaitIndex, []interface{}{
						validateIntMin(0),
					}),
				},
				queryOptWaitTime: {
					Optional: true,
					Type:     schema.TypeString,
					ValidateFunc: makeValidationFunc(queryOptWaitTime, []interface{}{
						validateDurationMin("0ns"),
					}),
				},
			},
		},
	}
}

func getQueryOpts(d *schema.ResourceData, client *consulapi.Client, meta interface{}) (*consulapi.QueryOptions, error) {
	queryOpts := &consulapi.QueryOptions{}

	if v, ok := d.GetOk(catalogNodesQueryOpts); ok {
		for _, config := range v.(*schema.Set).List() {
			queryOptions := config.(map[string]interface{})
			if v, ok := queryOptions[queryOptAllowStale]; ok {
				queryOpts.AllowStale = v.(bool)
			}

			if v, ok := queryOptions[queryOptDatacenter]; ok {
				queryOpts.Datacenter = v.(string)
			}

			if v, ok := queryOptions["namespace"]; ok {
				queryOpts.Namespace = v.(string)
			}

			if queryOpts.Datacenter == "" {
				dc, err := getDC(d, client, meta)
				if err != nil {
					return nil, err
				}
				queryOpts.Datacenter = dc
			}

			if v, ok := queryOptions[queryOptNear]; ok {
				queryOpts.Near = v.(string)
			}

			if v, ok := queryOptions[queryOptRequireConsistent]; ok {
				queryOpts.RequireConsistent = v.(bool)
			}

			if v, ok := queryOptions[queryOptNodeMeta]; ok {
				m := v.(map[string]interface{})
				nodeMetaMap := make(map[string]string, len(queryOptNodeMeta))
				for s, t := range m {
					nodeMetaMap[s] = t.(string)
				}
				queryOpts.NodeMeta = nodeMetaMap
			}

			if v, ok := queryOptions[queryOptToken]; ok {
				queryOpts.Token = v.(string)
			}

			if v, ok := queryOptions[queryOptWaitIndex]; ok {
				queryOpts.WaitIndex = uint64(v.(int))
			}

			if v, ok := queryOptions[queryOptWaitTime]; ok {
				d, _ := time.ParseDuration(v.(string))
				queryOpts.WaitTime = d
			}
		}
	}

	return queryOpts, nil
}
