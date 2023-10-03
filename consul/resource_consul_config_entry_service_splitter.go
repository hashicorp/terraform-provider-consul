// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type serviceSplitter struct{}

func (s *serviceSplitter) GetKind() string {
	return consulapi.ServiceSplitter
}

func (s *serviceSplitter) GetDescription() string {
	return "The `consul_config_entry_service_splitter` resource configures a [service splitter](https://developer.hashicorp.com/consul/docs/connect/config-entries/service-splitter) that will redirect a percentage of incoming traffic requests for a service to one or more specific service instances."
}

func (s *serviceSplitter) GetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Specifies a name for the configuration entry.",
			Required:    true,
			ForceNew:    true,
		},
		"partition": {
			Type:        schema.TypeString,
			Description: "Specifies the admin partition to apply the configuration entry.",
			Optional:    true,
			ForceNew:    true,
		},
		"namespace": {
			Type:        schema.TypeString,
			Description: "Specifies the namespace to apply the configuration entry.",
			Optional:    true,
			ForceNew:    true,
		},
		"meta": {
			Type:        schema.TypeMap,
			Description: "Specifies key-value pairs to add to the KV store.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"splits": {
			Type:        schema.TypeList,
			Description: "Defines how much traffic to send to sets of service instances during a traffic split.",
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"weight": {
						Type:        schema.TypeFloat,
						Description: "Specifies the percentage of traffic sent to the set of service instances specified in the `service` field. Each weight must be a floating integer between `0` and `100`. The smallest representable value is `.01`. The sum of weights across all splits must add up to `100`.",
						Required:    true,
					},
					"service": {
						Type:        schema.TypeString,
						Description: "Specifies the name of the service to resolve.",
						Required:    true,
					},
					"service_subset": {
						Type:        schema.TypeString,
						Description: "Specifies a subset of the service to resolve. A service subset assigns a name to a specific subset of discoverable service instances within a datacenter, such as `version2` or `canary`. All services have an unnamed default subset that returns all healthy instances.",
						Optional:    true,
					},
					"namespace": {
						Type:        schema.TypeString,
						Description: "Specifies the namespace to use in the FQDN when resolving the service.",
						Optional:    true,
					},
					"partition": {
						Type:        schema.TypeString,
						Description: "Specifies the admin partition to use in the FQDN when resolving the service.",
						Optional:    true,
					},
					"request_headers": {
						Type:        schema.TypeList,
						MaxItems:    1,
						Description: "Specifies a set of HTTP-specific header modification rules applied to requests routed with the service split. You cannot configure request headers if the listener protocol is set to `tcp`.",
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"add": {
									Type:        schema.TypeMap,
									Description: "Map of one or more key-value pairs. Defines a set of key-value pairs to add to the header. Use header names as the keys. Header names are not case-sensitive. If header values with the same name already exist, the value is appended and Consul applies both headers.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"set": {
									Type:        schema.TypeMap,
									Description: "Map of one or more key-value pairs. Defines a set of key-value pairs to add to the request header or to replace existing header values with. Use header names as the keys. Header names are not case-sensitive. If header values with the same names already exist, Consul replaces the header values.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"remove": {
									Type:        schema.TypeList,
									Description: "Defines an list of headers to remove. Consul removes only headers containing exact matches. Header names are not case-sensitive.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
							},
						},
					},
					"response_headers": {
						Type:        schema.TypeList,
						MaxItems:    1,
						Description: "Specifies a set of HTTP-specific header modification rules applied to responses routed with the service split. You cannot configure request headers if the listener protocol is set to `tcp`.",
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"add": {
									Type:        schema.TypeMap,
									Description: "Map of one or more key-value pairs. Defines a set of key-value pairs to add to the header. Use header names as the keys. Header names are not case-sensitive. If header values with the same name already exist, the value is appended and Consul applies both headers.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"set": {
									Type:        schema.TypeMap,
									Description: "Map of one or more key-value pairs. Defines a set of key-value pairs to add to the request header or to replace existing header values with. Use header names as the keys. Header names are not case-sensitive. If header values with the same names already exist, Consul replaces the header values.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
								"remove": {
									Type:        schema.TypeList,
									Description: "Defines an list of headers to remove. Consul removes only headers containing exact matches. Header names are not case-sensitive.",
									Optional:    true,
									Elem:        &schema.Schema{Type: schema.TypeString},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s *serviceSplitter) Decode(d *schema.ResourceData) (consulapi.ConfigEntry, error) {
	configEntry := &consulapi.ServiceSplitterConfigEntry{
		Kind: consulapi.ServiceSplitter,
		Name: d.Get("name").(string),
		Meta: map[string]string{},
	}

	if d.Get("namespace") != nil {
		configEntry.Namespace = d.Get("namespace").(string)
	}

	if d.Get("partition") != nil {
		configEntry.Partition = d.Get("partition").(string)
	}

	if d.Get("meta") != nil {
		for k, v := range d.Get("meta").(map[string]interface{}) {
			configEntry.Meta[k] = v.(string)
		}
	}

	if d.Get("splits") != nil {
		for _, raw := range d.Get("splits").([]interface{}) {
			s := raw.(map[string]interface{})
			split := consulapi.ServiceSplit{
				Weight:        float32(s["weight"].(float64)),
				Service:       s["service"].(string),
				ServiceSubset: s["service_subset"].(string),
				Namespace:     s["namespace"].(string),
				Partition:     s["partition"].(string),
				RequestHeaders: &consulapi.HTTPHeaderModifiers{
					Add:    map[string]string{},
					Set:    map[string]string{},
					Remove: []string{},
				},
				ResponseHeaders: &consulapi.HTTPHeaderModifiers{
					Add:    map[string]string{},
					Set:    map[string]string{},
					Remove: []string{},
				},
			}

			addHeaders := func(modifier *consulapi.HTTPHeaderModifiers, path string) {
				elems := s[path].([]interface{})
				if len(elems) == 0 {
					return
				}

				headers := elems[0].(map[string]interface{})
				for k, v := range headers["add"].(map[string]interface{}) {
					modifier.Add[k] = v.(string)
				}
				for k, v := range headers["set"].(map[string]interface{}) {
					modifier.Set[k] = v.(string)
				}
				for _, v := range headers["remove"].([]interface{}) {
					modifier.Remove = append(modifier.Remove, v.(string))
				}
			}
			addHeaders(split.RequestHeaders, "request_headers")
			addHeaders(split.ResponseHeaders, "response_headers")

			configEntry.Splits = append(configEntry.Splits, split)
		}
	}

	return configEntry, nil
}

func (s *serviceSplitter) Write(ce consulapi.ConfigEntry, sw *stateWriter) error {
	sp, ok := ce.(*consulapi.ServiceSplitterConfigEntry)
	if !ok {
		return fmt.Errorf("expected '%s' but got '%s'", consulapi.ServiceSplitter, ce.GetKind())
	}

	sw.set("name", sp.Name)
	sw.set("partition", sp.Partition)
	sw.set("namespace", sp.Namespace)

	meta := map[string]interface{}{}
	for k, v := range sp.Meta {
		meta[k] = v
	}
	sw.set("meta", meta)

	splits := []interface{}{}
	for _, s := range sp.Splits {
		split := map[string]interface{}{
			"weight":         s.Weight,
			"service":        s.Service,
			"service_subset": s.ServiceSubset,
			"namespace":      s.Namespace,
			"partition":      s.Partition,
		}
		addHeaders := func(modifier *consulapi.HTTPHeaderModifiers, path string) {
			headers := map[string]interface{}{}

			add := map[string]interface{}{}
			for k, v := range modifier.Add {
				add[k] = v
			}
			headers["add"] = add

			set := map[string]interface{}{}
			for k, v := range modifier.Set {
				set[k] = v
			}
			headers["set"] = set

			remove := []interface{}{}
			for _, v := range modifier.Remove {
				remove = append(remove, v)
			}
			headers["remove"] = remove

			split[path] = []interface{}{headers}
		}
		addHeaders(s.RequestHeaders, "request_headers")
		addHeaders(s.ResponseHeaders, "response_headers")
		splits = append(splits, split)
	}
	sw.set("splits", splits)

	return nil
}
