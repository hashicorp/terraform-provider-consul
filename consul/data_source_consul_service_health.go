package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	serviceHealthService    = "service"
	serviceHealthDatacenter = "datacenter"
	serviceHealthNear       = "near"
	serviceHealthTag        = "tag"
	serviceHealthNodeMeta   = "node_meta"
	serviceHealthPassing    = "passing"

	serviceHealthNodes               = "nodes"
	serviceHealthNodeID              = "node_id"
	serviceHealthNodeName            = "node_name"
	serviceHealthNodeAddress         = "node_address"
	serviceHealthNodeDatacenter      = "node_datacenter"
	serviceHealthNodeTaggedAddresses = "node_tagged_addresses"

	serviceHealthServiceID      = "service_id"
	serviceHealthServiceName    = "service_name"
	serviceHealthServiceTags    = "service_tags"
	serviceHealthServiceAddress = "service_address"
	serviceHealthServiceMeta    = "service_meta"
	serviceHealthServicePort    = "service_port"

	serviceHealthChecks           = "checks"
	serviceHealthCheckNode        = "node"
	serviceHealthCheckID          = "check_id"
	serviceHealthCheckName        = "name"
	serviceHealthCheckStatus      = "status"
	serviceHealthCheckNotes       = "notes"
	serviceHealthCheckOutput      = "output"
	serviceHealthCheckServiceID   = "service_id"
	serviceHealthCheckServiceName = "service_name"
	serviceHealthCheckServiceTags = "service_tags"
)

func dataSourceConsulServiceHealth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulServiceHealthRead,
		Schema: map[string]*schema.Schema{
			// Filter parameters
			serviceHealthService: &schema.Schema{
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			serviceHealthDatacenter: &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			serviceHealthNear: &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			serviceHealthTag: &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			serviceHealthNodeMeta: &schema.Schema{
				Optional: true,
				Type:     schema.TypeMap,
				ForceNew: true,
			},
			serviceHealthPassing: &schema.Schema{
				Optional: true,
				Type:     schema.TypeBool,
				ForceNew: true,
			},

			// Out parameters
			serviceHealthNodes: &schema.Schema{
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						serviceHealthNodeID: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						serviceHealthNodeName: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						serviceHealthNodeAddress: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						serviceHealthNodeDatacenter: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						serviceHealthNodeTaggedAddresses: &schema.Schema{
							Computed: true,
							Type:     schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						serviceHealthNodeMeta: &schema.Schema{
							Computed: true,
							Type:     schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						serviceHealthServiceID: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						serviceHealthServiceName: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						serviceHealthServiceTags: &schema.Schema{
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						serviceHealthServiceAddress: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						serviceHealthServiceMeta: &schema.Schema{
							Computed: true,
							Type:     schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						serviceHealthServicePort: &schema.Schema{
							Computed: true,
							Type:     schema.TypeInt,
						},
						serviceHealthChecks: &schema.Schema{
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									serviceHealthCheckNode: &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									serviceHealthCheckID: &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									serviceHealthCheckName: &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									serviceHealthCheckStatus: &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									serviceHealthCheckNotes: &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									serviceHealthCheckOutput: &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									serviceHealthCheckServiceID: &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									serviceHealthCheckServiceName: &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									serviceHealthCheckServiceTags: &schema.Schema{
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceConsulServiceHealthRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	health := client.Health()

	serviceName := d.Get(serviceHealthService).(string)
	serviceTag := d.Get(serviceHealthTag).(string)
	passingOnly := d.Get(serviceHealthPassing).(bool)
	near := d.Get(serviceHealthNear).(string)
	nodeMeta := d.Get(serviceHealthNodeMeta).(map[string]interface{})

	dc := d.Get(serviceHealthDatacenter).(string)
	if dc == "" {
		var err error
		dc, err = getDC(d, client)
		if err != nil {
			return err
		}
	}

	queryNodeMeta := map[string]string{}
	for key, value := range nodeMeta {
		queryNodeMeta[key] = value.(string)
	}

	qOps := &consulapi.QueryOptions{
		Near:       near,
		NodeMeta:   queryNodeMeta,
		Datacenter: dc,
	}
	serviceEntries, _, err := health.Service(serviceName, serviceTag, passingOnly, qOps)
	if err != nil {
		return fmt.Errorf("Failed to retrieve service health: %v", err)
	}

	l := make([]interface{}, 0, len(serviceEntries))
	for _, serviceEntry := range serviceEntries {
		m := make(map[string]interface{}, 6+4)

		m[serviceHealthNodeID] = serviceEntry.Node.ID
		m[serviceHealthNodeName] = serviceEntry.Node.Node
		m[serviceHealthNodeAddress] = serviceEntry.Node.Address
		m[serviceHealthNodeDatacenter] = serviceEntry.Node.Datacenter
		m[serviceHealthNodeTaggedAddresses] = serviceEntry.Node.TaggedAddresses
		m[serviceHealthNodeMeta] = serviceEntry.Node.Meta

		m[serviceHealthServiceID] = serviceEntry.Service.ID
		m[serviceHealthServiceName] = serviceEntry.Service.Service
		m[serviceHealthServiceAddress] = serviceEntry.Service.Address
		m[serviceHealthServicePort] = serviceEntry.Service.Port
		m[serviceHealthServiceTags] = serviceEntry.Service.Tags
		m[serviceHealthServiceMeta] = serviceEntry.Service.Meta

		c := make([]interface{}, 0, len(serviceEntry.Checks))
		for _, healthCheck := range serviceEntry.Checks {
			check := make(map[string]interface{}, 8)

			check[serviceHealthCheckNode] = healthCheck.Node
			check[serviceHealthCheckID] = healthCheck.CheckID
			check[serviceHealthCheckName] = healthCheck.Name
			check[serviceHealthCheckStatus] = healthCheck.Status
			check[serviceHealthCheckNotes] = healthCheck.Notes
			check[serviceHealthCheckOutput] = healthCheck.Output
			check[serviceHealthCheckServiceID] = healthCheck.ServiceID
			check[serviceHealthCheckServiceName] = healthCheck.ServiceName
			check[serviceHealthCheckServiceTags] = healthCheck.ServiceTags

			c = append(c, check)
		}
		l = append(l, m)
	}

	const idKeyFmt = "service-health-%s-%q-%q"
	d.SetId(fmt.Sprintf(idKeyFmt, dc, serviceName, serviceTag))
	d.Set(serviceHealthDatacenter, dc)
	d.Set(serviceHealthNear, near)
	d.Set(serviceHealthTag, serviceTag)
	d.Set(serviceHealthNodeMeta, nodeMeta)
	d.Set(serviceHealthPassing, passingOnly)
	d.Set(serviceHealthNodes, l)

	return nil
}
