package consul

import (
	"fmt"
	"time"

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

			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"external": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

			"check": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"notes": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"status": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"definition": &schema.Schema{
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tcp": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},

									"http": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},

									"tls_skip_verify": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"method": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Default:  "GET",
									},

									"header": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},

												"value": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},

									"interval": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},

									"timeout": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},

									"deregister_critical_service_after": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
								},
							},
						},
					},
				},
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
		return fmt.Errorf("Cannot retrieve node '%s': %v", node, err)
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

	// By default, the ID will match the name of the service
	// which we use later to query the catalog entry
	ident := name

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
		// If we are specifying an ID, we need to
		// query it as such
		ident = serviceID.(string)
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

	var nodeMeta map[string]string
	nodeMeta = make(map[string]string)
	if d.Get("external").(bool) {
		nodeMeta["external-node"] = "true"
		nodeMeta["external-probe"] = "true"
	}
	registration.NodeMeta = nodeMeta

	checks, err := parseChecks(node, name, d)
	if err != nil {
		return fmt.Errorf("Failed to fecth health-checks: %v", err)
	}
	registration.Checks = checks

	if _, err := catalog.Register(registration, &wOpts); err != nil {
		return fmt.Errorf("Failed to register service (dc: '%s'): %v", dc, err)
	}

	// Retrieve the service again to get the canonical service ID. We can't
	// get this back from the register call or through
	service, err := retrieveService(client, name, ident, node, dc)
	if err != nil {
		return fmt.Errorf("Failed to retrieve service '%s' after registration. This may mean that the service should be manually deregistered. %v", ident, err)
	}

	d.SetId(service.ServiceID)

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

	var nodeMeta map[string]string
	nodeMeta = make(map[string]string)
	if d.Get("external").(bool) {
		nodeMeta["external-node"] = "true"
		nodeMeta["external-probe"] = "true"
	}
	registration.NodeMeta = nodeMeta

	checks, err := parseChecks(node, name, d)
	if err != nil {
		return fmt.Errorf("Failed to fetch health-checks: %v", err)
	}
	registration.Checks = checks

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

	id := d.Id()
	name := d.Get("name").(string)
	node := d.Get("node").(string)

	service, err := retrieveService(client, name, id, node, dc)
	if err != nil {
		return err
	}

	d.Set("address", service.ServiceAddress)
	d.Set("service_id", service.ServiceID)
	d.Set("datacenter", service.Datacenter)
	d.Set("name", service.ServiceName)
	d.Set("port", service.ServicePort)
	tags := make([]string, 0, len(service.ServiceTags))
	for _, tag := range service.ServiceTags {
		tags = append(tags, tag)
	}
	d.Set("tags", tags)
	d.Set("node", service.Node)
	if externalNode, present := service.NodeMeta["external-node"]; present && externalNode == "true" {
		d.Set("external", true)
	} else {
		d.Set("external", false)
	}

	checks := make([]map[string]interface{}, 0)
	for _, check := range service.Checks {
		m := make(map[string]interface{})
		m["check_id"] = check.CheckID
		m["name"] = check.Name
		m["notes"] = check.Notes
		m["status"] = check.Status
		definition := make(map[string]interface{})
		definition["tcp"] = check.Definition.TCP
		definition["http"] = check.Definition.HTTP
		definition["tls_skip_verify"] = check.Definition.TLSSkipVerify
		definition["method"] = check.Definition.Method
		// definition["header"] = check.Definition.Header
		definition["interval"] = check.Definition.Interval
		definition["timeout"] = check.Definition.Timeout
		definition["deregister_critical_service_after"] = check.Definition.DeregisterCriticalServiceAfter
		m["definition"] = definition
		checks = append(checks, m)
	}
	d.Set("checks", checks)

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

func retrieveService(client *consulapi.Client, name string, ident string, node string, dc string) (*consulapi.CatalogService, error) {
	qOpts := consulapi.QueryOptions{Datacenter: dc}
	services, _, err := client.Catalog().Service(name, "", &qOpts)
	if err != nil {
		return nil, err
	}

	// Only one service with a given ID may be present per node
	for _, s := range services {
		if (s.ServiceID == ident) && (s.Node == node) {
			healthChecks, _, err := client.Health().Checks(name, &qOpts)
			if err != nil {
				return nil, fmt.Errorf("Failed to fetch health-checks: %v", err)
			}
			s.Checks = healthChecks
			return s, nil
		}
	}

	return nil, fmt.Errorf("Failed to retrieve service: '%s', services: %v", name, len(services))
}

func parseChecks(node string, name string, d *schema.ResourceData) ([]*consulapi.HealthCheck, error) {
	// Get health checks definition
	checks := d.Get("check").([]interface{})
	s := []*consulapi.HealthCheck{}
	s = make([]*consulapi.HealthCheck, len(checks))
	for i, raw := range checks {
		sub, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Failed to unroll: %#v", raw)
		}
		definition, ok := sub["definition"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Failed to unroll: %#v", raw)
		}
		interval, err := time.ParseDuration(definition["interval"].(string))
		if err != nil {
			return nil, fmt.Errorf("Failed to parse interval: %#v", interval)
		}
		timeout, err := time.ParseDuration(definition["timeout"].(string))
		if err != nil {
			return nil, fmt.Errorf("Failed to parse timeout: %#v", timeout)
		}

		if definition["tcp"] != nil && definition["http"] != nil {
			return nil, fmt.Errorf("You cannot set both tcp and http in the same check")
		}

		var tcp string
		if definition["tcp"] != nil {
			tcp = definition["tcp"].(string)
		}

		var http string
		if definition["http"] != nil {
			http = definition["http"].(string)
		}
		// var header map[string][]string
		// header = make(map[string][]string, 0)
		var tlsSkipVerify bool
		if definition["tls_skip_verify"] != nil {
			tlsSkipVerify = definition["tls_skip_verify"].(string) == "true"
		}
		var method string
		if definition["method"] != nil {
			method = definition["method"].(string)
		}
		healthCheck := &consulapi.HealthCheckDefinition{
			HTTP: http,
			// Header:        header,
			Method:        method,
			TLSSkipVerify: tlsSkipVerify,
			TCP:           tcp,
			Interval:      *consulapi.NewReadableDuration(interval),
			Timeout:       *consulapi.NewReadableDuration(timeout),
		}
		var deregisterCriticalServiceAfter string
		if definition["deregister_critical_service_after"] == nil {
			deregisterCriticalServiceAfter = ""
		} else {
			deregisterCriticalServiceAfter = definition["deregister_critical_service_after"].(string)
		}
		if deregisterCriticalServiceAfter != "" {
			deregisterCriticalServiceAfter, err := time.ParseDuration(deregisterCriticalServiceAfter)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse deregister_critical_service_after: %#v", deregisterCriticalServiceAfter)
			}
			healthCheck.DeregisterCriticalServiceAfter = *consulapi.NewReadableDuration(deregisterCriticalServiceAfter)
		}

		s[i] = &consulapi.HealthCheck{
			Node:       node,
			ServiceID:  name,
			CheckID:    sub["check_id"].(string),
			Name:       sub["name"].(string),
			Notes:      sub["notes"].(string),
			Status:     sub["status"].(string),
			Definition: *healthCheck,
		}
	}

	return s, nil
}
