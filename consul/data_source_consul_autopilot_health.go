package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceConsulAutopilotHealth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulAutopilotHealthRead,
		Schema: map[string]*schema.Schema{
			// Filters
			"datacenter": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
			},

			// Out parameters
			"healthy": &schema.Schema{
				Computed: true,
				Type:     schema.TypeBool,
			},
			"failure_tolerance": &schema.Schema{
				Computed: true,
				Type:     schema.TypeInt,
			},
			"servers": &schema.Schema{
				Computed: true,
				Type:     schema.TypeList,
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
						"serf_status": &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						"version": &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						"leader": &schema.Schema{
							Computed: true,
							Type:     schema.TypeBool,
						},
						"last_contact": &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						"last_term": &schema.Schema{
							Computed: true,
							Type:     schema.TypeInt,
						},
						"last_index": &schema.Schema{
							Computed: true,
							Type:     schema.TypeInt,
						},
						"healthy": &schema.Schema{
							Computed: true,
							Type:     schema.TypeBool,
						},
						"voter": &schema.Schema{
							Computed: true,
							Type:     schema.TypeBool,
						},
						"stable_since": &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceConsulAutopilotHealthRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	operator := client.Operator()

	queryOpts, err := getQueryOpts(d, client)
	if datacenter, ok := d.GetOk("datacenter"); ok {
		queryOpts.Datacenter = datacenter.(string)
	}

	health, err := operator.AutopilotServerHealth(queryOpts)
	if err != nil {
		return err
	}
	const idKeyFmt = "autopilot-health-%s"
	d.SetId(fmt.Sprintf(idKeyFmt, queryOpts.Datacenter))

	d.Set("healthy", health.Healthy)
	d.Set("failure_tolerance", health.FailureTolerance)

	serversHealth := make([]interface{}, 0, len(health.Servers))
	for _, server := range health.Servers {
		h := make(map[string]interface{}, 12)

		h["id"] = server.ID
		h["name"] = server.Name
		h["address"] = server.Address
		h["serf_status"] = server.SerfStatus
		h["version"] = server.Version
		h["leader"] = server.Leader
		h["last_contact"] = server.LastContact.String()
		h["last_term"] = server.LastTerm
		h["last_index"] = server.LastIndex
		h["healthy"] = server.Healthy
		h["voter"] = server.Voter
		h["stable_since"] = server.StableSince.String()

		serversHealth = append(serversHealth, h)
	}

	if err := d.Set("servers", serversHealth); err != nil {
		return errwrap.Wrapf("Unable to store servers health: {{err}}", err)
	}
	return nil
}
