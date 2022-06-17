package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSourceConsulPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulPeeringCreate,
		Read:   resourceConsulPeeringRead,
		Delete: resourceConsulPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// In
			"peer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peering_token": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
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

			// Out
			"deleted_at": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceConsulPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	name := d.Get("peer_name").(string)

	m := map[string]string{}
	for k, v := range d.Get("meta").(map[string]interface{}) {
		m[k] = v.(string)
	}

	req := api.PeeringEstablishRequest{
		PeerName:     name,
		PeeringToken: d.Get("peering_token").(string),
		Datacenter:   d.Get("datacenter").(string),
		Token:        d.Get("token").(string),
		Meta:         m,
	}

	_, _, err := client.Peerings().Establish(context.Background(), req, wOpts)
	if err != nil {
		return fmt.Errorf("failed to create peering: %w", err)
	}

	d.SetId(name)
	return resourceConsulPeeringRead(d, meta)
}

func resourceConsulPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Id()

	peer, _, err := client.Peerings().Read(context.Background(), name, qOpts)
	if err != nil {
		return fmt.Errorf("failed to list peerings: %w", err)
	}
	if peer == nil {
		d.SetId("")
		return nil
	}

	var deletedAt string
	if peer.DeletedAt != nil {
		deletedAt = peer.DeletedAt.String()
	}

	sw := newStateWriter(d)
	sw.set("peer_name", peer.Name)
	sw.set("deleted_at", deletedAt)
	sw.set("meta", peer.Meta)
	sw.set("state", peer.State)
	sw.set("peer_id", peer.PeerID)
	sw.set("peer_ca_pems", peer.PeerCAPems)
	sw.set("peer_server_name", peer.PeerServerName)
	sw.set("peer_server_addresses", peer.PeerServerAddresses)

	return sw.error()
}

func resourceConsulPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	name := d.Get("peer_name").(string)

	_, err := client.Peerings().Delete(context.Background(), name, wOpts)
	if err != nil {
		return fmt.Errorf("failed to delete peering %#v: %w", name, err)
	}

	return nil
}
