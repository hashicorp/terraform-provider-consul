package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulNetworkSegments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulNetworkSegmentsRead,

		Schema: map[string]*schema.Schema{
			// Inputs
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

			// Outputs
			"segments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceConsulNetworkSegmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	operator := client.Operator()

	token := d.Get("token").(string)
	dc, err := getDC(d, client, meta)
	if err != nil {
		return err
	}

	qOpts := &consulapi.QueryOptions{
		Token:      token,
		Datacenter: dc,
	}
	segments, _, err := operator.SegmentList(qOpts)
	if err != nil {
		return fmt.Errorf("Failed to get segment list: %v", err)
	}

	d.SetId(fmt.Sprintf("consul-segments-%s", dc))
	if err := d.Set("segments", segments); err != nil {
		return fmt.Errorf("Failed to set 'segments': %v", err)
	}

	return nil
}
