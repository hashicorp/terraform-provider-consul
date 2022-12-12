package consul

import (
	"fmt"

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
				Type:       schema.TypeString,
				Optional:   true,
				Sensitive:  true,
				Deprecated: tokenDeprecationMessage,
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
	client, qOpts, _ := getClient(d, meta)
	operator := client.Operator()

	segments, _, err := operator.SegmentList(qOpts)
	if err != nil {
		return fmt.Errorf("Failed to get segment list: %v", err)
	}

	d.SetId(fmt.Sprintf("consul-segments-%s", qOpts.Datacenter))
	if err := d.Set("segments", segments); err != nil {
		return fmt.Errorf("Failed to set 'segments': %v", err)
	}

	return nil
}
