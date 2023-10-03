// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
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
			Description: "",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"match": {
						Type:        schema.TypeSet,
						Description: "",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"http": {
									Type:        schema.TypeSet,
									Description: "",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"path_exact": {
												Type:        schema.TypeString,
												Description: "",
											},
											"path_prefix": {
												Type:        schema.TypeString,
												Description: "",
											},
											"path_regex": {
												Type:        schema.TypeString,
												Description: "",
											},
											"methods": {
												Type: schema.TypeList,
												Elem: schema.TypeString,
											},
											"header": {
												Type: schema.TypeList,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"name": {
															Type: schema.TypeString,
														},
														"present": {
															Type: schema.TypeBool,
														},
														"exact": {
															Type: schema.TypeString,
														},
														"prefix": {
															Type: schema.TypeString,
														},
														"suffix": {
															Type: schema.TypeString,
														},
														"regex": {
															Type: schema.TypeString,
														},
														"invert": {
															Type: schema.TypeBool,
														},
													},
												},
											},
											"query_param": {
												Type: schema.TypeList,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"name": {
															Type: schema.TypeString,
														},
														"present": {
															Type: schema.TypeBool,
														},
														"exact": {
															Type: schema.TypeString,
														},
														"regex": {
															Type: schema.TypeString,
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
						Type: schema.TypeSet,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"service": {
									Type: schema.TypeString,
								},
								"service_subset": {
									Type: schema.TypeString,
								},
								"namespace": {
									Type: schema.TypeString,
								},
								"partition": {
									Type: schema.TypeString,
								},
								"prefix_rewrite": {
									Type: schema.TypeString,
								},
								"request_timeout": {
									Type: schema.TypeString,
								},
								"idle_timeout": {
									Type: schema.TypeString,
								},
								"num_retries": {
									Type: schema.TypeInt,
								},
								"retry_on_connect_failure": {
									Type: schema.TypeBool,
								},
								"retry_on": {
									Type: schema.TypeList,
									Elem: schema.TypeString,
								},
								"retry_on_status_code": {
									Type: schema.TypeList,
									Elem: schema.TypeInt,
								},
								"request_headers": {
									Type: schema.TypeSet,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"add": {
												Type: schema.TypeMap,
												Elem: schema.TypeString,
											},
											"set": {
												Type: schema.TypeMap,
												Elem: schema.TypeString,
											},
											"remote": {
												Type: schema.TypeMap,
												Elem: schema.TypeString,
											},
										},
									},
								},
								"response_headers": {
									Type: schema.TypeMap,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"add": {
												Type: schema.TypeMap,
												Elem: schema.TypeString,
											},
											"set": {
												Type: schema.TypeMap,
												Elem: schema.TypeString,
											},
											"remote": {
												Type: schema.TypeMap,
												Elem: schema.TypeString,
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

	routes := d.Get("routes")
	if routes != nil {
		routeList := routes.([]interface{})
		serviceRoutesList := make([]consulapi.ServiceRoute, len(routeList))
		for indx, r := range routeList {
			routListMap := r.(map[string]interface{})
			matchMap := routListMap["match"].(map[string]interface{})
			matchHTTP := matchMap["HTTP"].(map[string]interface{})
			var matchRoute *consulapi.ServiceRouteMatch
			matchRoute = new(consulapi.ServiceRouteMatch)
			var serviceRouteHTTPMatch *consulapi.ServiceRouteHTTPMatch
			serviceRouteHTTPMatch = new(consulapi.ServiceRouteHTTPMatch)
			serviceRouteHTTPMatch.PathExact = matchHTTP["path_exact"].(string)
			serviceRouteHTTPMatch.PathPrefix = matchHTTP["path_prefix"].(string)
			serviceRouteHTTPMatch.PathRegex = matchHTTP["path_regex"].(string)
			serviceRouteHTTPMatch.Methods = matchHTTP["methods"].([]string)
			var headers []consulapi.ServiceRouteHTTPMatchHeader
			matchHeaders := matchHTTP["header"].([]interface{})
			for index, h := range matchHeaders {
				header := h.(map[string]interface{})
				matchHeader := &consulapi.ServiceRouteHTTPMatchHeader{
					Name:    header["name"].(string),
					Present: header["present"].(bool),
					Exact:   header["exact"].(string),
					Prefix:  header["prefix"].(string),
					Suffix:  header["suffix"].(string),
					Regex:   header["regex"].(string),
					Invert:  header["present"].(bool),
				}
				headers[index] = *matchHeader
			}
			serviceRouteHTTPMatch.Header = headers
			queryParam := matchHTTP["query_param"].([]interface{})
			queryParamList := make([]consulapi.ServiceRouteHTTPMatchQueryParam, len(queryParam))
			for index, q := range queryParam {
				queryParamMap := q.(map[string]interface{})
				queryParamList[index].Name = queryParamMap["name"].(string)
				queryParamList[index].Regex = queryParamMap["regex"].(string)
				queryParamList[index].Present = queryParamMap["present"].(bool)
				queryParamList[index].Exact = queryParamMap["exact"].(string)
			}
			serviceRouteHTTPMatch.QueryParam = queryParamList
			matchRoute.HTTP = serviceRouteHTTPMatch
			serviceRoutesList[indx].Match = matchRoute
			var destination *consulapi.ServiceRouteDestination
			destination = new(consulapi.ServiceRouteDestination)
			destinationList := (routListMap["destination"].(*schema.Set)).List()
			if len(destinationList) > 0 {
				destinationMap := destinationList[0].(map[string]interface{})
				destination.Service = destinationMap["service"].(string)
				destination.ServiceSubset = destinationMap["service_subset"].(string)
				destination.Namespace = destinationMap["namespace"].(string)
				destination.Partition = destinationMap["partition"].(string)
				destination.PrefixRewrite = destinationMap["prefix_rewrite"].(string)
				requestTimeout, err := time.ParseDuration(destinationMap["request_timeout"].(string))
				if err != nil {
					return nil, err
				}
				destination.RequestTimeout = requestTimeout
				idleTimeout, err := time.ParseDuration(destinationMap["idle_timeout"].(string))
				if err != nil {
					return nil, err
				}
				destination.IdleTimeout = idleTimeout
			}
			serviceRoutesList[indx].Destination = destination
		}
		configEntry.Routes = serviceRoutesList
	}
	return configEntry, nil
}

func (s *serviceRouter) Write(ce consulapi.ConfigEntry, sw *stateWriter) error {
	return nil
}
