// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var headerResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the header.",
		},
		"value": {
			Type:        schema.TypeList,
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The header's list of values.",
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

var ErrNoServiceRegistered error = errors.New("no service was found in consul catalog")

func resourceConsulService() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulServiceCreate,
		Update: resourceConsulServiceUpdate,
		Read:   resourceConsulServiceRead,
		Delete: resourceConsulServiceDelete,

		Description: `
A high-level resource for creating a Service in Consul in the Consul catalog. This
is appropriate for registering [external services](https://www.consul.io/docs/guides/external.html) and
can be used to create services addressable by Consul that cannot be registered
with a [local agent](https://www.consul.io/docs/agent/basics.html).

-> **NOTE:** If a Consul agent is running on the node where this service is
registered, it is not recommended to use this resource as the service will be
removed during the next [anti-entropy synchronization](https://www.consul.io/docs/architecture/anti-entropy).
`,

		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The address of the service. Defaults to the address of the node.",
			},

			"service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "If the service ID is not provided, it will be defaulted to the value of the `name` attribute.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the service.",
			},

			"node": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the node the to register the service on.",
			},

			"datacenter": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The datacenter to use. This overrides the agent's default datacenter and the datacenter in the provider setup.",
			},

			"external": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "The external field has been deprecated and does nothing.",
			},

			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port of the service.",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of values that are opaque to Consul, but can be used to distinguish between services or nodes.",
			},

			"meta": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A map of arbitrary KV metadata linked to the service instance.",
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
							Type:        schema.TypeString,
							Required:    true,
							Description: "An ID, *unique per agent*.",
						},

						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the health-check.",
						},

						"notes": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An opaque field meant to hold human readable text.",
						},

						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The initial health-check status.",
						},
						"tcp": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The TCP address and port to connect to for a TCP check.",
						},

						"http": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP endpoint to call for an HTTP check.",
						},

						"header": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        headerResource,
							Description: "The headers to send for an HTTP check. The attributes of each header is given below.",
						},

						"tls_skip_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to deactivate certificate verification for HTTP health-checks. Defaults to `false`.",
						},

						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "GET",
							Description: "The method to use for HTTP health-checks. Defaults to `GET`.",
						},

						"interval": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The interval to wait between each health-check invocation.",
						},

						"timeout": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies a timeout for outgoing connections in the case of a HTTP or TCP check.",
						},

						"deregister_critical_service_after": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "30s",
							Description: "The time after which the service is automatically deregistered when in the `critical` state. Defaults to `30s`. Setting to `0` will disable.",
						},
					},
				},
			},

			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The namespace to create the service within.",
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The partition the service is associated with.",
			},

			"enable_tag_override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies to disable the anti-entropy feature for this service's tags. Defaults to `false`.",
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
		return fmt.Errorf("failed to register service (dc: '%s'): %v", wOpts.Datacenter, err)
	}

	// Retrieve the service again to get the canonical service ID. We can't
	// get this back from the register call or through
	service, err := retrieveService(client, name, ident, node, qOpts)
	if err != nil {
		return fmt.Errorf("failed to retrieve service '%s' after registration. This may mean that the service should be manually deregistered. %v", ident, err)
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
		return fmt.Errorf("failed to update service (dc: '%s'): %v", wOpts.Datacenter, err)
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
		if err == ErrNoServiceRegistered {
			d.SetId("")
			return nil
		}
		return err
	}

	sw := newStateWriter(d)

	sw.set("address", service.ServiceAddress)
	sw.set("service_id", service.ServiceID)
	sw.set("datacenter", service.Datacenter)
	sw.set("name", service.ServiceName)
	sw.set("port", service.ServicePort)
	sw.set("tags", service.ServiceTags)
	sw.set("node", service.Node)

	serviceMeta := service.ServiceMeta
	delete(serviceMeta, consulSourceKey)
	sw.set("meta", serviceMeta)

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
	sw.set("check", checks)
	sw.set("enable_tag_override", service.ServiceEnableTagOverride)
	sw.set("namespace", service.Namespace)
	sw.set("partition", service.Partition)

	return sw.error()
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
		return fmt.Errorf("failed to deregister Consul service with id '%s' in %s: %v",
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
				return nil, fmt.Errorf("failed to fetch health-checks: %v", err)
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
	return nil, ErrNoServiceRegistered
}

func parseChecks(node string, serviceID string, d *schema.ResourceData) ([]*consulapi.HealthCheck, error) {
	// Get health checks definition
	checks := d.Get("check").(*schema.Set).List()
	s := make([]*consulapi.HealthCheck, len(checks))
	for i, raw := range checks {
		check, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to unroll: %#v", raw)
		}
		headers, err := parseHeaders(check)
		if err != nil {
			return nil, err
		}
		interval, err := time.ParseDuration(check["interval"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse interval: %#v", interval)
		}
		timeout, err := time.ParseDuration(check["timeout"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse timeout: %#v", timeout)
		}

		tcp := check["tcp"].(string)
		http := check["http"].(string)
		if tcp != "" && http != "" {
			return nil, fmt.Errorf("you cannot set both tcp and http in the same check")
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
				return nil, fmt.Errorf("failed to parse deregister_critical_service_after: %#v", deregisterCriticalServiceAfter)
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
	headers := make(map[string][]string)
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
		return nil, "", fmt.Errorf("cannot retrieve node '%s': %v", node, err)
	}
	if nodeCheck == nil {
		return nil, "", fmt.Errorf("node does not exist: '%s'", node)
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
		return nil, "", fmt.Errorf("failed to fetch health-checks: %v", err)
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
