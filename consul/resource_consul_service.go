package consul

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var headerResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeList,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	},
}

const (
	// ConsulSourceKey is the name of the meta attribute used by Consul to
	// record the origin of a service.
	consulSourceKey = "external-source"
	// ConsulSourceValue is its value.
	consulSourceValue = "terraform"
)

var NoServiceRegistered error = errors.New("No service was found in consul catalog")

func resourceConsulService() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulServiceCreate,
		Update: resourceConsulServiceUpdate,
		Read:   resourceConsulServiceRead,
		Delete: resourceConsulServiceDelete,

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"service_id": {
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

			"node": {
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

			"external": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "The external field has been deprecated and does nothing.",
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"meta": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"check": {
				Type: schema.TypeSet,
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					headers := []string{}
					for _, h := range m["header"].(*schema.Set).List() {
						name := h.(map[string]interface{})["name"].(string)
						value := ""
						for _, v := range h.(map[string]interface{})["value"].([]interface{}) {
							value += "-" + v.(string)
						}
						headers = append(headers, fmt.Sprintf("%s=%s", name, value))
					}

					attrs := []string{
						m["check_id"].(string),
						m["name"].(string),
						m["notes"].(string),
						m["tcp"].(string),
						m["http"].(string),
						strconv.FormatBool(m["tls_skip_verify"].(bool)),
						m["method"].(string),
						m["interval"].(string),
						m["timeout"].(string),
						m["deregister_critical_service_after"].(string),
					}
					attrs = append(attrs, headers...)

					return hashcode.String(hashcode.Strings(attrs))
				},
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"notes": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tcp": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"http": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"header": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     headerResource,
						},

						"tls_skip_verify": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"method": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "GET",
						},

						"interval": {
							Type:     schema.TypeString,
							Required: true,
						},

						"timeout": {
							Type:     schema.TypeString,
							Required: true,
						},

						"deregister_critical_service_after": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "30s",
						},
					},
				},
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"enable_tag_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceConsulServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	catalog := client.Catalog()

	name := d.Get("name").(string)
	node := d.Get("node").(string)

	registration, ident, err := getCatalogRegistration(d, meta)
	if err != nil {
		return err
	}

	if _, err := catalog.Register(registration, wOpts); err != nil {
		return fmt.Errorf("Failed to register service (dc: '%s'): %v", wOpts.Datacenter, err)
	}

	// Retrieve the service again to get the canonical service ID. We can't
	// get this back from the register call or through
	service, err := retrieveService(client, name, ident, node, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to retrieve service '%s' after registration. This may mean that the service should be manually deregistered. %v", ident, err)
	}

	d.SetId(service.ServiceID)

	return resourceConsulServiceRead(d, meta)
}

func resourceConsulServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	catalog := client.Catalog()

	registration, _, err := getCatalogRegistration(d, meta)
	if err != nil {
		return err
	}

	if _, err := catalog.Register(registration, wOpts); err != nil {
		return fmt.Errorf("Failed to update service (dc: '%s'): %v", wOpts.Datacenter, err)
	}

	return resourceConsulServiceRead(d, meta)
}

func resourceConsulServiceRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	id := d.Id()
	name := d.Get("name").(string)
	node := d.Get("node").(string)

	service, err := retrieveService(client, name, id, node, qOpts)
	if err != nil {
		if err == NoServiceRegistered {
			d.SetId("")
			return nil
		}
		return err
	}

	if err = d.Set("address", service.ServiceAddress); err != nil {
		return fmt.Errorf("Failed to store 'address': %s", err)
	}
	if err = d.Set("service_id", service.ServiceID); err != nil {
		return fmt.Errorf("Failed to store 'service_id': %s", err)
	}
	if err = d.Set("datacenter", service.Datacenter); err != nil {
		return fmt.Errorf("Failed to store 'datacenter': %s", err)
	}
	if err = d.Set("name", service.ServiceName); err != nil {
		return fmt.Errorf("Failed to store 'name': %s", err)
	}
	if err = d.Set("port", service.ServicePort); err != nil {
		return fmt.Errorf("Failed to store 'port': %s", err)
	}
	tags := make([]string, 0, len(service.ServiceTags))
	for _, tag := range service.ServiceTags {
		tags = append(tags, tag)
	}
	if err = d.Set("tags", tags); err != nil {
		return fmt.Errorf("Failed to store 'tags': %s", err)
	}
	if err = d.Set("node", service.Node); err != nil {
		return fmt.Errorf("Failed to store 'node': %s", err)
	}
	serviceMeta := service.ServiceMeta
	delete(serviceMeta, consulSourceKey)
	if err = d.Set("meta", serviceMeta); err != nil {
		return fmt.Errorf("Failed to store 'meta': %s", err)
	}

	checks := make([]map[string]interface{}, 0)
	for _, check := range service.Checks {
		m := make(map[string]interface{})
		m["check_id"] = check.CheckID
		m["name"] = check.Name
		m["notes"] = check.Notes
		m["status"] = check.Status
		m["tcp"] = check.Definition.TCP
		m["http"] = check.Definition.HTTP
		m["tls_skip_verify"] = check.Definition.TLSSkipVerify
		m["method"] = check.Definition.Method
		m["interval"] = check.Definition.Interval.String()
		m["timeout"] = check.Definition.Timeout.String()
		m["deregister_critical_service_after"] = check.Definition.DeregisterCriticalServiceAfter.String()
		headers := make([]interface{}, 0)
		for name, value := range check.Definition.Header {
			header := make(map[string]interface{})
			header["name"] = name

			valueInterface := make([]interface{}, 0)
			for _, v := range value {
				valueInterface = append(valueInterface, v)
			}

			header["value"] = valueInterface
			headers = append(headers, header)
		}

		// Setting a Set in a List does not work correctly
		// see https://github.com/hashicorp/terraform/issues/16331 for details
		m["header"] = schema.NewSet(
			schema.HashResource(headerResource),
			headers,
		)

		checks = append(checks, m)
	}
	if err := d.Set("check", checks); err != nil {
		return errwrap.Wrapf("Unable to store checks: {{err}}", err)
	}

	if err = d.Set("enable_tag_override", service.ServiceEnableTagOverride); err != nil {
		return fmt.Errorf("Failed to set 'enable_tag_override': %v", err)
	}

	return nil
}

func resourceConsulServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	catalog := client.Catalog()
	id := d.Id()
	node := d.Get("node").(string)

	// If we specified a custom service_id, we need
	// to utilize it for the delete
	if serviceID, ok := d.GetOk("service_id"); ok {
		id = serviceID.(string)
	}

	deregistration := consulapi.CatalogDeregistration{
		Datacenter: wOpts.Datacenter,
		Node:       node,
		ServiceID:  id,
	}

	if _, err := catalog.Deregister(&deregistration, wOpts); err != nil {
		return fmt.Errorf("Failed to deregister Consul service with id '%s' in %s: %v",
			id, wOpts.Datacenter, err)
	}

	// Clear the ID
	d.SetId("")
	return nil
}

func retrieveService(client *consulapi.Client, name, ident, node string, qOpts *consulapi.QueryOptions) (*consulapi.CatalogService, error) {
	services, _, err := client.Catalog().Service(name, "", qOpts)
	if err != nil {
		return nil, err
	}

	// Only one service with a given ID may be present per node
	for _, s := range services {
		if (s.ServiceID == ident) && (s.Node == node) {
			// Fetch health-checks for this service
			healthChecks, _, err := client.Health().Checks(name, qOpts)
			if err != nil {
				return nil, fmt.Errorf("Failed to fetch health-checks: %v", err)
			}
			// Filter the checks that correspond to this specific service instance
			s.Checks = make([]*consulapi.HealthCheck, 0)
			for _, h := range healthChecks {
				if h.Node == node && h.ServiceID == ident {
					s.Checks = append(s.Checks, h)
				}
			}
			return s, nil
		}
	}

	// No matching service has been found
	return nil, NoServiceRegistered
}

func parseChecks(node string, serviceID string, d *schema.ResourceData) ([]*consulapi.HealthCheck, error) {
	// Get health checks definition
	checks := d.Get("check").(*schema.Set).List()
	s := []*consulapi.HealthCheck{}
	s = make([]*consulapi.HealthCheck, len(checks))
	for i, raw := range checks {
		check, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Failed to unroll: %#v", raw)
		}
		headers, err := parseHeaders(check)
		if err != nil {
			return nil, err
		}
		interval, err := time.ParseDuration(check["interval"].(string))
		if err != nil {
			return nil, fmt.Errorf("Failed to parse interval: %#v", interval)
		}
		timeout, err := time.ParseDuration(check["timeout"].(string))
		if err != nil {
			return nil, fmt.Errorf("Failed to parse timeout: %#v", timeout)
		}

		tcp := check["tcp"].(string)
		http := check["http"].(string)
		if tcp != "" && http != "" {
			return nil, fmt.Errorf("You cannot set both tcp and http in the same check")
		}
		var tlsSkipVerify bool
		if check["tls_skip_verify"] != nil {
			tlsSkipVerify = check["tls_skip_verify"].(bool)
		}
		var method string
		if check["method"] != nil {
			method = check["method"].(string)
		}
		healthCheck := consulapi.HealthCheckDefinition{
			HTTP:          http,
			Header:        headers,
			Method:        method,
			TLSSkipVerify: tlsSkipVerify,
			TCP:           tcp,
			Interval:      *consulapi.NewReadableDuration(interval),
			Timeout:       *consulapi.NewReadableDuration(timeout),
		}
		var deregisterCriticalServiceAfter string
		if check["deregister_critical_service_after"] == nil {
			deregisterCriticalServiceAfter = ""
		} else {
			deregisterCriticalServiceAfter = check["deregister_critical_service_after"].(string)
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
			ServiceID:  serviceID,
			CheckID:    check["check_id"].(string),
			Name:       check["name"].(string),
			Notes:      check["notes"].(string),
			Status:     check["status"].(string),
			Definition: healthCheck,
		}
	}

	return s, nil
}

func parseHeaders(check map[string]interface{}) (map[string][]string, error) {
	headers := make(map[string][]string, 0)
	header := check["header"].(*schema.Set).List()
	for _, h := range header {
		name := h.(map[string]interface{})["name"].(string)
		value := h.(map[string]interface{})["value"]
		for _, v := range value.([]interface{}) {
			headers[name] = append(headers[name], v.(string))
		}
	}
	return headers, nil
}

func getCatalogRegistration(d *schema.ResourceData, meta interface{}) (*consulapi.CatalogRegistration, string, error) {
	client, qOpts, _ := getClient(d, meta)

	name := d.Get("name").(string)
	node := d.Get("node").(string)
	ident := name

	// Check to see if the node exists. We do this because
	// the Consul API will upsert nodes that don't exist, but
	// Terraform won't be able to track that. Requiring
	// them to exist either ensures that it is knowlingly tracked
	// outside of TF state or that it is referencing a node
	// managed by the consul_node resource (or datasource)
	nodeCheck, _, err := client.Catalog().Node(node, qOpts)
	if err != nil {
		return nil, "", fmt.Errorf("Cannot retrieve node '%s': %v", node, err)
	}
	if nodeCheck == nil {
		return nil, "", fmt.Errorf("Node does not exist: '%s'", node)
	}

	registration := &consulapi.CatalogRegistration{
		Datacenter: qOpts.Datacenter,
		Node:       node,
		Service: &consulapi.AgentService{
			Service: name,
		},
		// Creating a service should not modify the node
		// See https://github.com/hashicorp/terraform-provider-consul/issues/101
		SkipNodeUpdate: true,
	}

	// If we have a service_id
	if serviceID, ok := d.GetOk("service_id"); ok {
		registration.Service.ID = serviceID.(string)
		ident = serviceID.(string)
	}

	// If the address is not specified, use the nodes
	if address, ok := d.GetOk("address"); ok {
		registration.Address = address.(string)
		registration.Service.Address = address.(string)
	} else {
		registration.Address = nodeCheck.Node.Address
		registration.Service.Address = nodeCheck.Node.Address
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

	checks, err := parseChecks(node, ident, d)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to fetch health-checks: %v", err)
	}
	registration.Checks = checks

	serviceMeta := map[string]string{
		consulSourceKey: consulSourceValue,
	}
	for k, v := range d.Get("meta").(map[string]interface{}) {
		serviceMeta[k] = v.(string)
	}
	registration.Service.Meta = serviceMeta

	registration.Service.EnableTagOverride = d.Get("enable_tag_override").(bool)

	return registration, ident, nil
}
