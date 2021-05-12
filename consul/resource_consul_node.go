package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulNodeCreate,
		Update: resourceConsulNodeCreate,
		Read:   resourceConsulNodeRead,
		Delete: resourceConsulNodeDelete,

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"meta": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: false,
			},

			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceConsulNodeCreate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	catalog := client.Catalog()

	address := d.Get("address").(string)
	name := d.Get("name").(string)

	registration := &consulapi.CatalogRegistration{
		Address:    address,
		Datacenter: wOpts.Datacenter,
		Node:       name,
	}

	if v, ok := d.GetOk("meta"); ok {
		nodeMeta := make(map[string]string)
		for k, j := range v.(map[string]interface{}) {
			nodeMeta[k] = j.(string)
		}
		registration.NodeMeta = nodeMeta
	}

	if _, err := catalog.Register(registration, wOpts); err != nil {
		return fmt.Errorf("Failed to register Consul catalog node with name '%s' at address '%s' in %s: %v",
			name, address, wOpts.Datacenter, err)
	}

	// Update the resource
	if _, _, err := catalog.Node(name, qOpts); err != nil {
		return fmt.Errorf("Failed to read Consul catalog node with name '%s' at address '%s' in %s: %v",
			name, address, qOpts.Datacenter, err)
	} else {
		d.Set("datacenter", qOpts.Datacenter)
	}

	d.SetId(fmt.Sprintf("%s-%s", name, address))

	return nil
}

func resourceConsulNodeRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	catalog := client.Catalog()

	name := d.Get("name").(string)

	n, _, err := catalog.Node(name, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to get name '%s' from Consul catalog: %v", name, err)
	}

	if n == nil {
		d.SetId("")
		return nil
	}

	if err = d.Set("address", n.Node.Address); err != nil {
		return fmt.Errorf("Failed to set 'address': %v", err)
	}
	if err = d.Set("meta", n.Node.Meta); err != nil {
		return fmt.Errorf("Failed to set 'meta': %v", err)
	}
	return nil
}

func resourceConsulNodeDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	catalog := client.Catalog()

	address := d.Get("address").(string)
	name := d.Get("name").(string)

	deregistration := consulapi.CatalogDeregistration{
		Address:    address,
		Datacenter: wOpts.Datacenter,
		Node:       name,
	}

	if _, err := catalog.Deregister(&deregistration, wOpts); err != nil {
		return fmt.Errorf("Failed to deregister Consul catalog node with name '%s' at address '%s' in %s: %v",
			name, address, wOpts.Datacenter, err)
	}

	// Clear the ID
	d.SetId("")
	return nil
}
