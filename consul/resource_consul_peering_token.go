package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSourceConsulPeeringToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulPeeringTokenCreate,
		Read: func(*schema.ResourceData, interface{}) error {
			return nil
		},
		Delete: func(*schema.ResourceData, interface{}) error {
			return nil
		},

		Schema: map[string]*schema.Schema{
			"peer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"partition": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"token": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"meta": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"peering_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceConsulPeeringTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	name := d.Get("peer_name").(string)

	m := map[string]string{}
	for k, v := range d.Get("meta").(map[string]interface{}) {
		m[k] = v.(string)
	}

	req := api.PeeringGenerateTokenRequest{
		PeerName:   name,
		Partition:  d.Get("partition").(string),
		Datacenter: d.Get("datacenter").(string),
		Token:      d.Get("token").(string),
		Meta:       m,
	}

	resp, _, err := client.Peerings().GenerateToken(context.Background(), req, wOpts)
	if err != nil {
		return fmt.Errorf("failed to create peering token: %w", err)
	}

	d.SetId(name)

	sw := newStateWriter(d)
	sw.set("peering_token", resp.PeeringToken)

	return sw.error()
}
