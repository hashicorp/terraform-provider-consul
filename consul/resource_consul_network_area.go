package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulNetworkArea() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulNetworkAreaCreate,
		Read:   resourceConsulNetworkAreaRead,
		Update: resourceConsulNetworkAreaUpdate,
		Delete: resourceConsulNetworkAreaDelete,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"peer_datacenter": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"retry_join": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"use_tls": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceConsulNetworkAreaCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	operator := client.Operator()

	area := &consulapi.Area{
		PeerDatacenter: d.Get("peer_datacenter").(string),
		UseTLS:         d.Get("use_tls").(bool),
	}

	if v, ok := d.GetOk("retry_join"); ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		area.RetryJoin = s
	}

	id, _, err := operator.AreaCreate(area, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to create network area: %v", err)
	}

	d.SetId(id)
	return resourceConsulNetworkAreaRead(d, meta)
}

func resourceConsulNetworkAreaRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	operator := client.Operator()

	id := d.Id()

	area, _, err := operator.AreaGet(id, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to get %s area: %v", id, err)
	}

	if len(area) == 0 {
		d.SetId("")
		return nil
	}

	if len(area) != 1 {
		return fmt.Errorf("There should be only one area")
	}

	peerDatacenter := area[0].PeerDatacenter
	retryJoin := area[0].RetryJoin
	useTLS := area[0].UseTLS

	if err = d.Set("peer_datacenter", peerDatacenter); err != nil {
		return fmt.Errorf("Failed to set 'peer_datacenter': %v", err)
	}
	if err = d.Set("retry_join", retryJoin); err != nil {
		return fmt.Errorf("Failed to set 'retry_join': %v", err)
	}
	if err = d.Set("use_tls", useTLS); err != nil {
		return fmt.Errorf("Failed to set 'use_tls': %v", err)
	}

	return nil
}

func resourceConsulNetworkAreaUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	operator := client.Operator()

	id := d.Id()

	// Consul Enterprise version <= 1.4.? needs to have PeerDatacenter and
	// RetryJoin set during updates.
	// We still mark `retry_join` as ForceNew as this may change in the future.
	// Issue: https://github.com/hashicorp/consul/issues/5727
	area := &consulapi.Area{
		PeerDatacenter: d.Get("peer_datacenter").(string),
		UseTLS:         d.Get("use_tls").(bool),
	}

	if v, ok := d.GetOk("retry_join"); ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		area.RetryJoin = s
	}

	_id, _, err := operator.AreaUpdate(id, area, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to update '%s' network area: %v", id, err)
	}

	if id != _id {
		return fmt.Errorf("This should not happen")
	}

	d.SetId(id)
	return nil
}

func resourceConsulNetworkAreaDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	operator := client.Operator()

	id := d.Id()

	_, err := operator.AreaDelete(id, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to delete '%s' network area: %v", err, id)
	}

	d.SetId("")
	return nil
}
