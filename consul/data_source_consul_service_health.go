package consul

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulServiceHealth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulServiceHealthRead,
		Schema: map[string]*schema.Schema{
			// Filter parameters
			"datacenter": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				// ForceNew: true,
			},
			"name": &schema.Schema{
				Required: true,
				Type:     schema.TypeString,
				// ForceNew: true,
			},
			"near": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				// ForceNew: true,
			},
			"tag": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				// ForceNew: true,
			},
			"node_meta": &schema.Schema{
				Optional: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// ForceNew: true,
			},
			"passing": &schema.Schema{
				Optional: true,
				Type:     schema.TypeBool,
				// ForceNew: true,
				Default: true,
			},
			"wait_for": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
				// ForceNew: true,
			},

			// Out parameters
			"results": &schema.Schema{
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": &schema.Schema{
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							// TODO(remi): set AsSingle when releasing
							// a version for Terraform 0.12
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"name": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"address": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"datacenter": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"tagged_addresses": &schema.Schema{
										Computed: true,
										Type:     schema.TypeMap,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"meta": &schema.Schema{
										Computed: true,
										Type:     schema.TypeMap,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"service": &schema.Schema{
							Computed: true,
							Type:     schema.TypeList,
							MaxItems: 1,
							// TODO(remi): set AsSingle when releasing
							// a version for Terraform 0.12
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"name": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"tags": &schema.Schema{
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"address": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"meta": &schema.Schema{
										Computed: true,
										Type:     schema.TypeMap,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"port": &schema.Schema{
										Computed: true,
										Type:     schema.TypeInt,
									},
								},
							},
						},
						"checks": &schema.Schema{
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"id": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"name": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"status": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"notes": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"output": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"service_id": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"service_name": &schema.Schema{
										Computed: true,
										Type:     schema.TypeString,
									},
									"service_tags": &schema.Schema{
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
	client := getClient(meta)
	health := client.Health()

	serviceName := d.Get("name").(string)
	serviceTag := d.Get("tag").(string)
	passingOnly := d.Get("passing").(bool)
	near := d.Get("near").(string)
	nodeMeta := d.Get("node_meta").(map[string]interface{})

	dc := d.Get("datacenter").(string)
	if dc == "" {
		var err error
		dc, err = getDC(d, client, meta)
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
	var err error
	var serviceEntries []*consulapi.ServiceEntry
	if d.Get("wait_for").(string) == "" || !passingOnly {
		log.Printf("[INFO] Fetching health information for service '%s'", serviceName)
		serviceEntries, _, err = health.Service(serviceName, serviceTag, passingOnly, qOps)
		if err != nil {
			return fmt.Errorf("Failed to retrieve service health: %v", err)
		}
	} else {
		waitFor, err := time.ParseDuration(d.Get("wait_for").(string))
		if err != nil {
			return fmt.Errorf("Could not parse 'wait_for': %s", err)
		}
		log.Printf("[INFO] Fetching health information for service '%s' for %s", serviceName, waitFor)
		err = resource.Retry(waitFor, func() *resource.RetryError {

			serviceEntries, _, err = health.Service(serviceName, serviceTag, passingOnly, qOps)
			if err != nil {
				return resource.RetryableError(fmt.Errorf("Failed to retrieve service health: %v", err))
			}
			if len(serviceEntries) == 0 {
				return resource.RetryableError(fmt.Errorf("No healthy service found"))
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Failed to wait for '%s' to be healthy: %s", serviceName, err)
		}
	}

	results := make([]interface{}, 0, len(serviceEntries))
	for _, serviceEntry := range serviceEntries {
		m := make(map[string]interface{})

		node := make(map[string]interface{})
		node["id"] = serviceEntry.Node.ID
		node["name"] = serviceEntry.Node.Node
		node["address"] = serviceEntry.Node.Address
		node["datacenter"] = serviceEntry.Node.Datacenter
		node["tagged_addresses"] = serviceEntry.Node.TaggedAddresses
		node["meta"] = serviceEntry.Node.Meta

		m["node"] = []map[string]interface{}{
			node,
		}

		service := make(map[string]interface{})
		service["id"] = serviceEntry.Service.ID
		service["name"] = serviceEntry.Service.Service
		service["address"] = serviceEntry.Service.Address
		service["port"] = serviceEntry.Service.Port
		service["tags"] = serviceEntry.Service.Tags
		service["meta"] = serviceEntry.Service.Meta

		m["service"] = []map[string]interface{}{
			service,
		}

		checks := make([]interface{}, 0, len(serviceEntry.Checks))
		for _, healthCheck := range serviceEntry.Checks {
			check := make(map[string]interface{}, 8)

			check["node"] = healthCheck.Node
			check["id"] = healthCheck.CheckID
			check["name"] = healthCheck.Name
			check["status"] = healthCheck.Status
			check["notes"] = healthCheck.Notes
			check["output"] = healthCheck.Output
			check["service_id"] = healthCheck.ServiceID
			check["service_name"] = healthCheck.ServiceName
			check["service_tags"] = healthCheck.ServiceTags

			checks = append(checks, check)
		}
		m["checks"] = checks
		results = append(results, m)
	}

	const idKeyFmt = "service-health-%s-%q-%q"
	d.SetId(fmt.Sprintf(idKeyFmt, dc, serviceName, serviceTag))
	if err = d.Set("datacenter", dc); err != nil {
		return fmt.Errorf("Failed to set 'datacenter': %s", err)
	}
	if err = d.Set("near", near); err != nil {
		return fmt.Errorf("Failed to set 'near': %s", err)
	}
	if err = d.Set("tag", serviceTag); err != nil {
		return fmt.Errorf("Failed to set 'tag': %s", err)
	}
	if err = d.Set("node_meta", nodeMeta); err != nil {
		return fmt.Errorf("Failed to set 'node_meta': %s", err)
	}
	if err = d.Set("passing", passingOnly); err != nil {
		return fmt.Errorf("Failed to set 'passing': %s", err)
	}
	if err = d.Set("results", results); err != nil {
		return fmt.Errorf("Failed to set 'results': %s", err)
	}

	return nil
}
