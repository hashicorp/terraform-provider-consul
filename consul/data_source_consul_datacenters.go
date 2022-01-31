package consul

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulDatacenters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulDatacentersRead,
		Schema: map[string]*schema.Schema{
			// Out parameters
			"datacenters": {
				Computed: true,
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceConsulDatacentersRead(d *schema.ResourceData, meta interface{}) error {
	client, _, _ := getClient(d, meta)

	datacenters, err := client.Catalog().Datacenters()
	if err != nil {
		return err
	}

	d.SetId("-")
	sw := newStateWriter(d)
	sw.set("datacenters", datacenters)

	return sw.error()
}
