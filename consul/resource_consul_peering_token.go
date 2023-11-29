// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSourceConsulPeeringToken() *schema.Resource {
	return &schema.Resource{
		Description: `
[Cluster Peering](https://www.consul.io/docs/connect/cluster-peering) can be used to create connections between two or more independent clusters so that services deployed to different partitions or datacenters can communicate.

The ` + "`cluster_peering_token`" + ` resource can be used to generate a peering token that can later be used to establish a peering connection.

The functionality described here is available only in Consul version 1.13.0 and later.
`,

		Create: resourceConsulPeeringTokenCreate,
		Read:   resourceConsulPeeringTokenRead,
		Delete: resourceConsulPeeringTokenDelete,

		Schema: map[string]*schema.Schema{
			"peer_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name assigned to the peer cluster. The `peer_name` is used to reference the peer cluster in service discovery queries and configuration entries such as `service-intentions`. This field must be a valid DNS hostname label.",
			},
			"partition": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"meta": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies KV metadata to associate with the peering. This parameter is not required and does not directly impact the cluster peering process.",

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"peering_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The generated peering token",
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
		PeerName:  name,
		Partition: d.Get("partition").(string),
		Meta:      m,
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

func resourceConsulPeeringTokenRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Id()
	client, qOpts, _ := getClient(d, meta)

	peer, _, err := client.Peerings().Read(context.Background(), name, qOpts)
	if err != nil {
		return fmt.Errorf("failed to find peer %q: %s", name, err)
	}

	if peer == nil {
		d.SetId("")
	}

	return nil
}

func resourceConsulPeeringTokenDelete(d *schema.ResourceData, meta interface{}) error {
	name := d.Id()
	client, _, wOpts := getClient(d, meta)

	_, err := client.Peerings().Delete(context.Background(), name, wOpts)
	if err != nil {
		return fmt.Errorf("failed to delete peer %q: %s", name, err)
	}
	return nil
}
