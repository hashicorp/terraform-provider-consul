package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConsulService() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulServiceCreate,
		Update: resourceConsulServiceUpdate,
		Read:   resourceConsulServiceRead,
		Delete: resourceConsulServiceDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"service_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"node": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceConsulServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	catalog := client.Catalog()

	name := d.Get("name").(string)
	node := d.Get("node").(string)

	dc := ""
	if _, ok := d.GetOk("datacenter"); ok {
		dc = d.Get("datacenter").(string)
	}

	// Check to see if the node exists. We do this because
	// the Consul API will upsert nodes that don't exist, but
	// Terraform won't be able to track that. Requiring
	// them to exist either ensures that it is knowlingly tracked
	// outside of TF state or that it is referencing a node
	// managed by the consul_node resource (or datasource)
	nodeCheck, _, err := client.Catalog().Node(node, &consulapi.QueryOptions{Datacenter: dc})
	if err != nil {
		return fmt.Errorf("Cannot retrieve node: %v", err)
	}
	if nodeCheck == nil {
		return fmt.Errorf("Node does not exist: '%s'", node)
	}

	// Setup the operations using the datacenter
	wOpts := consulapi.WriteOptions{Datacenter: dc}

	registration := &consulapi.CatalogRegistration{
		Datacenter: dc,
		Node:       node,
		Service: &consulapi.AgentService{
			Service: name,
		},
	}

	// If the address is not specified, use the nodes
	if address, ok := d.GetOk("address"); ok {
		registration.Address = address.(string)
		registration.Service.Address = address.(string)
	} else {
		registration.Address = nodeCheck.Node.Address
		registration.Service.Address = nodeCheck.Node.Address
	}

	if serviceID, ok := d.GetOk("service_id"); ok {
		registration.Service.ID = serviceID.(string)
	}

	if port, ok := d.GetOk("port"); ok {
		registration.Service.Port = port.(int)
	}

	if v, ok := d.GetOk("tags"); ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		registration.Service.Tags = s
	}

	if _, err := catalog.Register(registration, &wOpts); err != nil {
		return fmt.Errorf("Failed to register service (dc: '%s'): %v", dc, err)
	}

	// Set the ID, attributes will be populated by read. We assume the
	// service name is the ID as the Consul API returned successfully.
	d.SetId(name)

	return resourceConsulServiceRead(d, meta)
}

func resourceConsulServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	catalog := client.Catalog()

	name := d.Get("name").(string)
	node := d.Get("node").(string)

	dc := ""
	if _, ok := d.GetOk("datacenter"); ok {
		dc = d.Get("datacenter").(string)
	}

	// Setup the operations using the datacenter
	wOpts := consulapi.WriteOptions{Datacenter: dc}

	registration := &consulapi.CatalogRegistration{
		Datacenter: dc,
		Node:       node,
		Service: &consulapi.AgentService{
			Service: name,
		},
	}

	// If we have a service_id
	if serviceID, ok := d.GetOk("service_id"); ok {
		registration.Service.ID = serviceID.(string)
	}

	if address, ok := d.GetOk("address"); ok {
		registration.Address = address.(string)
		registration.Service.Address = address.(string)
	} else {
		// If we don't have an address, skip updating the node
		registration.SkipNodeUpdate = true
	}

	if port, ok := d.GetOk("port"); ok {
		registration.Service.Port = port.(int)
	}

	if v, ok := d.GetOk("tags"); ok {
		vs := v.([]interface{})
		s := make([]string, len(vs))
		for i, raw := range vs {
			s[i] = raw.(string)
		}
		registration.Service.Tags = s
	}

	if _, err := catalog.Register(registration, &wOpts); err != nil {
		return fmt.Errorf("Failed to update service (dc: '%s'): %v", dc, err)
	}

	return resourceConsulServiceRead(d, meta)
}

func resourceConsulServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	dc := ""
	if _, ok := d.GetOk("datacenter"); ok {
		dc = d.Get("datacenter").(string)
	}

	name := d.Get("name").(string)

	qOpts := consulapi.QueryOptions{Datacenter: dc}
	services, meta, err := client.Catalog().Service(name, "", &qOpts)
	if err != nil {
		return fmt.Errorf("Failed to retrieve service (dc: '%s'): %v", dc, err)
	}

	if len(services) > 1 {
		return fmt.Errorf("Multiple services returned after creation (dc: '%s'): %v", dc, services)
	}

	service := services[0]

	d.Set("address", service.ServiceAddress)
	d.Set("service_id", service.ID)
	d.Set("name", service.ServiceName)
	d.Set("port", service.ServicePort)
	tags := make([]string, 0, len(service.ServiceTags))
	for _, tag := range service.ServiceTags {
		tags = append(tags, tag)
	}
	d.Set("tags", tags)
	d.Set("node", service.Node)

	return nil
}

func resourceConsulServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	catalog := client.Catalog()
	id := d.Id()
	node := d.Get("node").(string)

	dc := ""
	if _, ok := d.GetOk("datacenter"); ok {
		dc = d.Get("datacenter").(string)
	}

	var token string
	if v, ok := d.GetOk("token"); ok {
		token = v.(string)
	}

	// If we specified a custom service_id, we need
	// to utilize it for the delete
	if serviceID, ok := d.GetOk("service_id"); ok {
		id = serviceID.(string)
	}

	// Setup the operations using the datacenter
	wOpts := consulapi.WriteOptions{Datacenter: dc, Token: token}

	// TODO(pearkes): Should we add an option to automatically
	// deregister nodes associated with the service?
	deregistration := consulapi.CatalogDeregistration{
		Datacenter: dc,
		Node:       node,
		ServiceID:  id,
	}

	if _, err := catalog.Deregister(&deregistration, &wOpts); err != nil {
		return fmt.Errorf("Failed to deregister Consul service with id '%s' in %s: %v",
			id, dc, err)
	}

	// Clear the ID
	d.SetId("")
	return nil
}
