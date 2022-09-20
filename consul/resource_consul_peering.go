package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSourceConsulPeering() *schema.Resource {
	return &schema.Resource{
		Description: `
[Cluster Peering](https://www.consul.io/docs/connect/cluster-peering) can be used to create connections between two or more independent clusters so that services deployed to different partitions or datacenters can communicate.

The ` + "`cluster_peering`" + ` resource can be used to establish the peering after a peering token has been generated.

~> **Cluster peering is currently in technical preview:** Functionality associated with cluster peering is subject to change. You should never use the technical preview release in secure environments or production scenarios. Features in technical preview may have performance issues, scaling issues, and limited support.

The functionality described here is available only in Consul version 1.13.0 and later.
`,
		Create: resourceConsulPeeringCreate,
		Read:   resourceConsulPeeringRead,
		Delete: resourceConsulPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// In
			"peer_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name assigned to the peer cluster. The `peer_name` is used to reference the peer cluster in service discovery queries and configuration entries such as `service-intentions`. This field must be a valid DNS hostname label.",
			},
			"peering_token": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The peering token fetched from the peer cluster.",
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
			"partition": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"imported_service_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"exported_service_count": {
				Type:     schema.TypeInt,
				Computed: true,
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
		Meta:         m,
		Partition:    d.Get("partition").(string),
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
	sw.set("imported_service_count", peer.ImportedServiceCount)
	sw.set("exported_service_count", peer.ExportedServiceCount)

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
