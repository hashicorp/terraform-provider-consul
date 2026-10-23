// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type serviceRouter struct{}

func (s *serviceRouter) GetKind() string {
	return consulapi.ServiceRouter
}

func (s *serviceRouter) GetDescription() string {
	return "The `consul_config_entry_service_router` resource configures a [service router](https://developer.hashicorp.com/consul/docs/connect/config-entries/service-router) to redirect a traffic request for a service to one or more specific service instances."
}

func (s *serviceRouter) GetSchema() map[string]*schema.Schema {
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
		"routes": {
			Type:        schema.TypeList,
			Description: "Defines the possible routes for L7 requests.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"match": {
						Type:        schema.TypeList,
						MaxItems:    1,
						Description: "Describes a set of criteria that Consul compares incoming L7 traffic with.",
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"http": {
									Type:        schema.TypeList,
									MaxItems:    1,
									Description: "Specifies a set of HTTP criteria used to evaluate incoming L7 traffic for matches.",
									Optional:    true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"path_exact": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Specifies the exact path to match on the HTTP request path.",
											},
											"path_prefix": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Specifies the path prefix to match on the HTTP request path.",
											},
											"path_regex": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Specifies a regular expression to match on the HTTP request path.",
											},
											"methods": {
												Type:        schema.TypeList,
												Description: "Specifies HTTP methods that the match applies to.",
												Elem:        &schema.Schema{Type: schema.TypeString},
												Optional:    true,
											},
											"header": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "Specifies information in the HTTP request header to match with.",
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"name": {
															Type:        schema.TypeString,
															Description: "Specifies the name of the HTTP header to match.",
															Optional:    true,
														},
														"present": {
															Type:        schema.TypeBool,
															Optional:    true,
															Description: "Specifies that a request matches when the value in the `name` argument is present anywhere in the HTTP header.",
														},
														"exact": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies that a request matches when the header with the given name is this exact value.",
														},
														"prefix": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies that a request matches when the header with the given name has this prefix.",
														},
														"suffix": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies that a request matches when the header with the given name has this suffix.",
														},
														"regex": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies that a request matches when the header with the given name matches this regular expression.",
														},
														"invert": {
															Type:        schema.TypeBool,
															Optional:    true,
															Description: "Specifies that the logic for the HTTP header match should be inverted.",
														},
													},
												},
											},
											"query_param": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "Specifies information to match to on HTTP query parameters.",
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"name": {
															Type:        schema.TypeString,
															Description: "Specifies the name of the HTTP query parameter to match.",
															Optional:    true,
														},
														"present": {
															Type:        schema.TypeBool,
															Optional:    true,
															Description: "Specifies that a request matches when the value in the `name` argument is present anywhere in the HTTP query parameter.",
														},
														"exact": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies that a request matches when the query parameter with the given name is this exact value.",
														},
														"regex": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies that a request matches when the query parameter with the given name matches this regular expression.",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"destination": {
						Type:        schema.TypeList,
						MaxItems:    1,
						Optional:    true,
						Description: "Specifies the target service to route matching requests to, as well as behavior for the request to follow when routed.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"service": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the name of the service to resolve.",
								},
								"service_subset": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies a named subset of the given service to resolve instead of the one defined as that service's `default_subset` in the service resolver configuration entry.",
								},
								"namespace": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the Consul namespace to resolve the service from instead of the current namespace.",
								},
								"partition": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the Consul admin partition to resolve the service from instead of the current partition.",
								},
								"prefix_rewrite": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies rewrites to the HTTP request path before proxying it to its final destination.",
								},
								"request_timeout": {
									Type:        schema.TypeString,
									Optional:    true,
									Default:     "0s",
									Description: "Specifies the total amount of time permitted for the entire downstream request to be processed, including retry attempts.",
								},
								"idle_timeout": {
									Type:        schema.TypeString,
									Optional:    true,
									Default:     "0s",
									Description: "Specifies the total amount of time permitted for the request stream to be idle.",
								},
								"num_retries": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Specifies the number of times to retry the request when a retry condition occurs.",
								},
								"retry_on_connect_failure": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Specifies that connection failure errors that trigger a retry request.",
								},
								"retry_on": {
									Type:        schema.TypeList,
									Elem:        &schema.Schema{Type: schema.TypeString},
									Optional:    true,
									Description: "Specifies a list of conditions for Consul to retry requests based on the response from an upstream service.",
								},
								"retry_on_status_codes": {
									Type:        schema.TypeList,
									Elem:        &schema.Schema{Type: schema.TypeInt},
									Optional:    true,
									Description: "Specifies a list of integers for HTTP response status codes that trigger a retry request.",
								},
								"request_headers": {
									Type:        schema.TypeList,
									MaxItems:    1,
									Optional:    true,
									Description: "Specifies a set of HTTP-specific header modification rules applied to requests routed with the service router.",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"add": {
												Type:        schema.TypeMap,
												Description: "Defines a set of key-value pairs to add to the header. Use header names as the keys.",
												Optional:    true,
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
											"set": {
												Type:        schema.TypeMap,
												Optional:    true,
												Description: "Defines a set of key-value pairs to add to the request header or to replace existing header values with.",
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
											"remove": {
												Type:        schema.TypeList,
												Description: "Defines a list of headers to remove.",
												Optional:    true,
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
										},
									},
								},
								"response_headers": {
									Type:        schema.TypeList,
									MaxItems:    1,
									Optional:    true,
									Description: "Specifies a set of HTTP-specific header modification rules applied to responses routed with the service router.",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"add": {
												Type:        schema.TypeMap,
												Optional:    true,
												Description: "Defines a set of key-value pairs to add to the header. Use header names as the keys",
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
											"set": {
												Type:        schema.TypeMap,
												Optional:    true,
												Description: "Defines a set of key-value pairs to add to the response header or to replace existing header values with",
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
											"remove": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "Defines a list of headers to remove.",
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
										},
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

func (s *serviceRouter) Decode(d *schema.ResourceData) (consulapi.ConfigEntry, error) {
	configEntry := &consulapi.ServiceRouterConfigEntry{
		Kind:      consulapi.ServiceRouter,
		Name:      d.Get("name").(string),
		Partition: d.Get("partition").(string),
		Namespace: d.Get("namespace").(string),
		Meta:      map[string]string{},
	}

	for k, v := range d.Get("meta").(map[string]interface{}) {
		configEntry.Meta[k] = v.(string)
	}

	for i, r := range d.Get("routes").([]interface{}) {
		route := r.(map[string]interface{})

		sr := consulapi.ServiceRoute{
			Destination: &consulapi.ServiceRouteDestination{},
		}

		for _, m := range route["match"].([]interface{}) {
			match := m.(map[string]interface{})

			for _, h := range match["http"].([]interface{}) {
				http := h.(map[string]interface{})

				sr.Match = &consulapi.ServiceRouteMatch{
					HTTP: &consulapi.ServiceRouteHTTPMatch{
						PathExact:  http["path_exact"].(string),
						PathPrefix: http["path_prefix"].(string),
						PathRegex:  http["path_regex"].(string),
					},
				}

				for _, method := range http["methods"].([]interface{}) {
					sr.Match.HTTP.Methods = append(sr.Match.HTTP.Methods, method.(string))
				}

				for _, h := range http["header"].([]interface{}) {
					header := h.(map[string]interface{})

					sr.Match.HTTP.Header = append(sr.Match.HTTP.Header, consulapi.ServiceRouteHTTPMatchHeader{
						Name:    header["name"].(string),
						Exact:   header["exact"].(string),
						Prefix:  header["prefix"].(string),
						Suffix:  header["suffix"].(string),
						Regex:   header["regex"].(string),
						Present: header["present"].(bool),
						Invert:  header["invert"].(bool),
					})
				}

				for _, q := range http["query_param"].([]interface{}) {
					queryParam := q.(map[string]interface{})

					sr.Match.HTTP.QueryParam = append(sr.Match.HTTP.QueryParam, consulapi.ServiceRouteHTTPMatchQueryParam{
						Name:    queryParam["name"].(string),
						Exact:   queryParam["exact"].(string),
						Regex:   queryParam["regex"].(string),
						Present: queryParam["present"].(bool),
					})
				}
			}
		}

		for _, d := range route["destination"].([]interface{}) {
			destination := d.(map[string]interface{})

			sr.Destination = &consulapi.ServiceRouteDestination{
				Service:               destination["service"].(string),
				ServiceSubset:         destination["service_subset"].(string),
				Namespace:             destination["namespace"].(string),
				Partition:             destination["partition"].(string),
				PrefixRewrite:         destination["prefix_rewrite"].(string),
				NumRetries:            uint32(destination["num_retries"].(int)),
				RetryOnConnectFailure: destination["retry_on_connect_failure"].(bool),
			}

			parseDuration := func(name string) (time.Duration, error) {
				dur, err := time.ParseDuration(destination[name].(string))
				if err != nil {
					return 0, fmt.Errorf("failed to parse routes[%d].destination.%s: %w", i, name, err)
				}
				return dur, nil
			}

			dur, err := parseDuration("request_timeout")
			if err != nil {
				return nil, err
			}
			sr.Destination.RequestTimeout = dur

			dur, err = parseDuration("idle_timeout")
			if err != nil {
				return nil, err
			}
			sr.Destination.IdleTimeout = dur

			for _, r := range destination["retry_on_status_codes"].([]interface{}) {
				sr.Destination.RetryOnStatusCodes = append(sr.Destination.RetryOnStatusCodes, uint32(r.(int)))
			}

			for _, r := range destination["retry_on"].([]interface{}) {
				sr.Destination.RetryOn = append(sr.Destination.RetryOn, r.(string))
			}

			parseHTTPHeaderModifiers := func(name string) *consulapi.HTTPHeaderModifiers {
				if len(destination[name].([]interface{})) == 0 {
					return nil
				}

				headers := destination[name].([]interface{})[0].(map[string]interface{})
				result := &consulapi.HTTPHeaderModifiers{
					Add: map[string]string{},
					Set: map[string]string{},
				}

				for k, v := range headers["add"].(map[string]interface{}) {
					result.Add[k] = v.(string)
				}

				for k, v := range headers["set"].(map[string]interface{}) {
					result.Add[k] = v.(string)
				}

				for _, v := range headers["remove"].([]interface{}) {
					result.Remove = append(result.Remove, v.(string))
				}

				return result
			}

			sr.Destination.RequestHeaders = parseHTTPHeaderModifiers("request_headers")
			sr.Destination.ResponseHeaders = parseHTTPHeaderModifiers("response_headers")
		}

		configEntry.Routes = append(configEntry.Routes, sr)
	}

	return configEntry, nil
}

func (s *serviceRouter) Write(ce consulapi.ConfigEntry, d *schema.ResourceData, sw *stateWriter) error {
	sr, ok := ce.(*consulapi.ServiceRouterConfigEntry)
	if !ok {
		return fmt.Errorf("expected '%s' but got '%s'", consulapi.ServiceDefaults, ce.GetKind())
	}

	sw.set("name", sr.Name)
	sw.set("partition", sr.Partition)
	sw.set("namespace", sr.Namespace)

	meta := map[string]interface{}{}
	for k, v := range sr.Meta {
		meta[k] = v
	}
	sw.set("meta", meta)

	routes := make([]map[string]interface{}, 0)
	for _, route := range sr.Routes {
		result := map[string]interface{}{}

		var http map[string]interface{}

		shouldSet := func() bool {
			isEmpty := route.Match == nil || route.Match.HTTP == nil || (route.Match.HTTP.PathExact == "" &&
				route.Match.HTTP.PathPrefix == "" &&
				route.Match.HTTP.PathRegex == "" &&
				len(route.Match.HTTP.Header) == 0 &&
				len(route.Match.HTTP.QueryParam) == 0 &&
				len(route.Match.HTTP.Methods) == 0)

			if !isEmpty {
				return true
			}

			routes := d.Get("routes").([]interface{})
			if len(routes) == 0 {
				return false
			}
			match := routes[0].(map[string]interface{})["match"].([]interface{})
			if len(match) == 0 {
				return false
			}
			http := match[0].(map[string]interface{})["http"].([]interface{})
			if len(http) == 0 {
				return false
			}
			if len(http[0].(map[string]interface{})["header"].([]interface{})) != 0 {
				return true
			}
			if len(http[0].(map[string]interface{})["query_param"].([]interface{})) != 0 {
				return true
			}

			return false
		}

		if shouldSet() {
			http = map[string]interface{}{
				"path_exact":  route.Match.HTTP.PathExact,
				"path_prefix": route.Match.HTTP.PathPrefix,
				"path_regex":  route.Match.HTTP.PathRegex,
				"methods":     route.Match.HTTP.Methods,
			}

			header := []interface{}{}
			for _, h := range route.Match.HTTP.Header {
				header = append(header, map[string]interface{}{
					"name":    h.Name,
					"present": h.Present,
					"exact":   h.Exact,
					"prefix":  h.Prefix,
					"suffix":  h.Suffix,
					"regex":   h.Regex,
					"invert":  h.Invert,
				})
			}
			http["header"] = header

			queryParam := []interface{}{}
			for _, q := range route.Match.HTTP.QueryParam {
				queryParam = append(queryParam, map[string]interface{}{
					"name":    q.Name,
					"present": q.Present,
					"exact":   q.Exact,
					"regex":   q.Regex,
				})
			}
			http["query_param"] = queryParam

			result["match"] = []interface{}{
				map[string]interface{}{
					"http": []interface{}{http},
				},
			}
		}

		shouldSet = func() bool {
			isEmpty := route.Destination == nil || (route.Destination.Service == "" &&
				route.Destination.ServiceSubset == "" &&
				route.Destination.Namespace == "" &&
				route.Destination.Partition == "" &&
				route.Destination.PrefixRewrite == "" &&
				route.Destination.RequestTimeout == 0 &&
				route.Destination.IdleTimeout == 0 &&
				route.Destination.NumRetries == 0 &&
				!route.Destination.RetryOnConnectFailure &&
				len(route.Destination.RetryOnStatusCodes) == 0 &&
				len(route.Destination.RetryOn) == 0 && (route.Destination.RequestHeaders == nil || len(route.Destination.RequestHeaders.Add)+len(route.Destination.RequestHeaders.Set)+len(route.Destination.RequestHeaders.Remove) == 0) &&
				route.Destination.ResponseHeaders == nil)

			if !isEmpty {
				return true
			}

			routes := d.Get("routes").([]interface{})
			if len(routes) == 0 {
				return false
			}
			destination := routes[0].(map[string]interface{})["destination"].([]interface{})
			return len(destination) != 0
		}

		if shouldSet() {
			destination := map[string]interface{}{
				"service":                  route.Destination.Service,
				"service_subset":           route.Destination.ServiceSubset,
				"namespace":                route.Destination.Namespace,
				"partition":                route.Destination.Partition,
				"prefix_rewrite":           route.Destination.PrefixRewrite,
				"request_timeout":          route.Destination.RequestTimeout.String(),
				"idle_timeout":             route.Destination.IdleTimeout.String(),
				"num_retries":              route.Destination.NumRetries,
				"retry_on_connect_failure": route.Destination.RetryOnConnectFailure,
				"retry_on":                 route.Destination.RetryOn,
				"retry_on_status_codes":    route.Destination.RetryOnStatusCodes,
			}

			convertHeaders := func(headers *consulapi.HTTPHeaderModifiers) []interface{} {
				if headers == nil {
					return []interface{}{}
				}

				add := map[string]interface{}{}
				for k, v := range headers.Add {
					add[k] = v
				}

				set := map[string]interface{}{}
				for k, v := range headers.Set {
					set[k] = v
				}

				remove := []interface{}{}
				for _, v := range headers.Remove {
					remove = append(remove, v)
				}

				return []interface{}{
					map[string]interface{}{
						"add":    add,
						"set":    set,
						"remove": remove,
					},
				}
			}
			destination["request_headers"] = convertHeaders(route.Destination.RequestHeaders)
			result["destination"] = []interface{}{destination}
		}

		routes = append(routes, result)
	}
	sw.set("routes", routes)

	return sw.error()
}
