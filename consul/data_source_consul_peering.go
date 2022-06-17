package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulPeering() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulPeeringRead,
		Schema: map[string]*schema.Schema{
			// In
			"peer_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"partition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Out
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
	}
}

func dataSourceConsulPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Get("peer_name").(string)

	peer, _, err := client.Peerings().Read(context.Background(), name, qOpts)
	if err != nil {
		return fmt.Errorf("failed to get peer information: %w", err)
	}
	if peer == nil {
		return fmt.Errorf("no peer name %#v found", name)
	}
	d.SetId(peer.ID)

	var deletedAt string
	if peer.DeletedAt != nil {
		deletedAt = peer.DeletedAt.String()
	}

	sw := newStateWriter(d)
	sw.set("deleted_at", deletedAt)
	sw.set("meta", peer.Meta)
	sw.set("state", peer.State)
	sw.set("peer_id", peer.PeerID)
	sw.set("peer_ca_pems", peer.PeerCAPems)
	sw.set("peer_server_name", peer.PeerServerName)
	sw.set("peer_server_addresses", peer.PeerServerAddresses)

	return sw.error()
}
