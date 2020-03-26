package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulNetworkAreaMembers() *schema.Resource {
	return &schema.Resource{
		Read: datasourceConsulNetworkAreaMembersRead,

		Schema: map[string]*schema.Schema{
			// Input
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Output
			"members": {
				Type:     schema.TypeList,
				Computed: true,
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
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rtt": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceConsulNetworkAreaMembersRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	operator := client.Operator()

	uuid := d.Get("uuid").(string)
	d.SetId(fmt.Sprintf("consul-network-area-members-%s", uuid))

	token := d.Get("token").(string)
	dc, err := getDC(d, client, meta)
	if err != nil {
		return err
	}

	qOpts := &consulapi.QueryOptions{
		Token:      token,
		Datacenter: dc,
	}
	members, _, err := operator.AreaMembers(uuid, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to fetch the list of members: %v", err)
	}

	res := make([]map[string]interface{}, len(members))
	for i, m := range members {
		res[i] = map[string]interface{}{
			"id":         m.ID,
			"name":       m.Name,
			"address":    m.Addr.String(),
			"port":       m.Port,
			"datacenter": m.Datacenter,
			"role":       m.Role,
			"build":      m.Build,
			"protocol":   m.Protocol,
			"status":     m.Status,
			"rtt":        m.RTT,
		}
	}
	if err = d.Set("members", res); err != nil {
		return fmt.Errorf("Failed to set 'members': %v", err)
	}

	return nil
}
