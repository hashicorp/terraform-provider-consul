package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulPeerings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulPeeringsRead,
		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"partition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"peers": {
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
						"partition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"meta": {
							Type:     schema.TypeMap,
							Computed: true,

							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_ca_pems": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"peer_server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_server_addresses": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceConsulPeeringsRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	peerings, _, err := client.Peerings().List(context.Background(), qOpts)
	if err != nil {
		return fmt.Errorf("failed to list peerings: %w", err)
	}

	peers := make([]interface{}, len(peerings))
	for i, peer := range peerings {
		var deletedAt string
		if peer.DeletedAt != nil {
			deletedAt = peer.DeletedAt.String()
		}
		peers[i] = map[string]interface{}{
			"id":                    peer.ID,
			"name":                  peer.Name,
			"partition":             peer.Partition,
			"deleted_at":            deletedAt,
			"meta":                  peer.Meta,
			"state":                 peer.State,
			"peer_id":               peer.PeerID,
			"peer_ca_pems":          peer.PeerCAPems,
			"peer_server_name":      peer.PeerServerName,
			"peer_server_addresses": peer.PeerServerAddresses,
		}
	}

	d.SetId("peers")

	sw := newStateWriter(d)
	sw.set("peers", peers)

	return sw.error()
}
