package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	autopilotHealthDatacenter        = "datacenter"
	autopilotHealthHealthy           = "healthy"
	autopilotHealthFailureTolerance  = "failure_tolerance"
	autopilotHealthServers           = "servers"
	autopilotHealthServerID          = "id"
	autopilotHealthServerName        = "name"
	autopilotHealthServerAddress     = "address"
	autopilotHealthServerSerfStatus  = "serf_status"
	autopilotHealthServerVersion     = "version"
	autopilotHealthServerLeader      = "leader"
	autopilotHealthServerLastContact = "last_contact"
	autopilotHealthServerLastTerm    = "last_term"
	autopilotHealthServerLastIndex   = "last_index"
	autopilotHealthServerHealthy     = "healthy"
	autopilotHealthServerVoter       = "voter"
	autopilotHealthServerStableSince = "stable_since"
)

func dataSourceConsulAutopilotHealth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulAutopilotHealthRead,
		Schema: map[string]*schema.Schema{
			// Filters
			autopilotHealthDatacenter: &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
			},

			// Out parameters
			autopilotHealthHealthy: &schema.Schema{
				Computed: true,
				Type:     schema.TypeBool,
			},
			autopilotHealthFailureTolerance: &schema.Schema{
				Computed: true,
				Type:     schema.TypeInt,
			},
			autopilotHealthServers: &schema.Schema{
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						autopilotHealthServerID: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						autopilotHealthServerName: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						autopilotHealthServerAddress: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						autopilotHealthServerSerfStatus: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						autopilotHealthServerVersion: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						autopilotHealthServerLeader: &schema.Schema{
							Computed: true,
							Type:     schema.TypeBool,
						},
						autopilotHealthServerLastContact: &schema.Schema{
							Computed: true,
							Type:     schema.TypeString,
						},
						autopilotHealthServerLastTerm: &schema.Schema{
							Computed: true,
							Type:     schema.TypeInt,
						},
						autopilotHealthServerLastIndex: &schema.Schema{
							Computed: true,
							Type:     schema.TypeInt,
						},
						autopilotHealthServerHealthy: &schema.Schema{
							Computed: true,
							Type:     schema.TypeBool,
						},
						autopilotHealthServerVoter: &schema.Schema{
							Computed: true,
							Type:     schema.TypeBool,
						},
						autopilotHealthServerStableSince: &schema.Schema{
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
	if datacenter, ok := d.GetOk(autopilotHealthDatacenter); ok {
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

		h[autopilotHealthServerID] = server.ID
		h[autopilotHealthServerName] = server.Name
		h[autopilotHealthServerAddress] = server.Address
		h[autopilotHealthServerSerfStatus] = server.SerfStatus
		h[autopilotHealthServerVersion] = server.Version
		h[autopilotHealthServerLeader] = server.Leader
		h[autopilotHealthServerLastContact] = server.LastContact.String()
		h[autopilotHealthServerLastTerm] = server.LastTerm
		h[autopilotHealthServerLastIndex] = server.LastIndex
		h[autopilotHealthServerHealthy] = server.Healthy
		h[autopilotHealthServerVoter] = server.Voter
		h[autopilotHealthServerStableSince] = server.StableSince.String()

		serversHealth = append(serversHealth, h)
	}

	if err := d.Set("servers", serversHealth); err != nil {
		return errwrap.Wrapf("Unable to store servers health: {{err}}", err)
	}
	return nil
}
