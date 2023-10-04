// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
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
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"match": {
						Type:        schema.TypeSet,
						Description: "",
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"http": {
									Type:        schema.TypeSet,
									Description: "",
									Optional:    true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"path_exact": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"path_prefix": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"path_regex": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "",
											},
											"methods": {
												Type:     schema.TypeList,
												Elem:     &schema.Schema{Type: schema.TypeString},
												Optional: true,
											},
											"header": {
												Type:     schema.TypeList,
												Optional: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"name": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"present": {
															Type:     schema.TypeBool,
															Optional: true,
														},
														"exact": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"prefix": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"suffix": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"regex": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"invert": {
															Type:     schema.TypeBool,
															Optional: true,
														},
													},
												},
											},
											"query_param": {
												Type:     schema.TypeList,
												Optional: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"name": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"present": {
															Type:     schema.TypeBool,
															Optional: true,
														},
														"exact": {
															Type:     schema.TypeString,
															Optional: true,
														},
														"regex": {
															Type:     schema.TypeString,
															Optional: true,
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
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"service": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"service_subset": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"namespace": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"partition": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"prefix_rewrite": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"request_timeout": {
									Type:     schema.TypeString,
									Optional: true,
									Default:  "0s",
								},
								"idle_timeout": {
									Type:     schema.TypeString,
									Optional: true,
									Default:  "0s",
								},
								"num_retries": {
									Type:     schema.TypeInt,
									Optional: true,
									Default:  0,
								},
								"retry_on_connect_failure": {
									Type:     schema.TypeBool,
									Optional: true,
									Default:  false,
								},
								"retry_on": {
									Type:     schema.TypeList,
									Elem:     &schema.Schema{Type: schema.TypeString},
									Optional: true,
								},
								"retry_on_status_codes": {
									Type:     schema.TypeList,
									Elem:     &schema.Schema{Type: schema.TypeInt},
									Optional: true,
								},
								"request_headers": {
									Type:     schema.TypeSet,
									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"add": {
												Type:     schema.TypeMap,
												Optional: true,
												Elem:     &schema.Schema{Type: schema.TypeString},
											},
											"set": {
												Type:     schema.TypeMap,
												Optional: true,
												Elem:     &schema.Schema{Type: schema.TypeString},
											},
											"remove": {
												Type:     schema.TypeList,
												Optional: true,
												Elem:     &schema.Schema{Type: schema.TypeString},
											},
										},
									},
								},
								"response_headers": {
									Type:     schema.TypeSet,
									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"add": {
												Type:     schema.TypeMap,
												Optional: true,
												Elem:     &schema.Schema{Type: schema.TypeString},
											},
											"set": {
												Type:     schema.TypeMap,
												Optional: true,
												Elem:     &schema.Schema{Type: schema.TypeString},
											},
											"remove": {
												Type:     schema.TypeList,
												Optional: true,
												Elem:     &schema.Schema{Type: schema.TypeString},
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
			matchList := routListMap["match"].(*schema.Set).List()
			if len(matchList) > 0 {
				matchMap := matchList[0].(map[string]interface{})
				matchHTTPMap := matchMap["http"].(*schema.Set).List()
				if len(matchHTTPMap) > 0 {
					matchHTTP := matchHTTPMap[0].(map[string]interface{})
					var matchRoute *consulapi.ServiceRouteMatch
					matchRoute = new(consulapi.ServiceRouteMatch)
					var serviceRouteHTTPMatch *consulapi.ServiceRouteHTTPMatch
					serviceRouteHTTPMatch = new(consulapi.ServiceRouteHTTPMatch)
					serviceRouteHTTPMatch.PathExact = matchHTTP["path_exact"].(string)
					serviceRouteHTTPMatch.PathPrefix = matchHTTP["path_prefix"].(string)
					serviceRouteHTTPMatch.PathRegex = matchHTTP["path_regex"].(string)
					methods := make([]string, 0)
					for _, v := range matchHTTP["methods"].([]interface{}) {
						methods = append(methods, v.(string))
					}
					serviceRouteHTTPMatch.Methods = methods
					headers := make([]consulapi.ServiceRouteHTTPMatchHeader, 0)
					matchHeaders := matchHTTP["header"].([]interface{})
					for _, h := range matchHeaders {
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
						headers = append(headers, *matchHeader)
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
				}
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
					if destinationMap["request_timeout"] != "" {
						requestTimeout, err := time.ParseDuration(destinationMap["request_timeout"].(string))
						if err != nil {
							return nil, err
						}
						destination.RequestTimeout = requestTimeout
					}
					if destinationMap["idle_timeout"] != "" {
						idleTimeout, err := time.ParseDuration(destinationMap["idle_timeout"].(string))
						if err != nil {
							return nil, err
						}
						destination.IdleTimeout = idleTimeout
					}
					destination.NumRetries = uint32(destinationMap["num_retries"].(int))
					destination.RetryOnConnectFailure = destinationMap["retry_on_connect_failure"].(bool)
					retryOnList := make([]string, 0)
					for _, v := range destinationMap["retry_on"].([]interface{}) {
						retryOnList = append(retryOnList, v.(string))
					}
					destination.RetryOn = retryOnList
					retryOnCodes := make([]uint32, 0)
					for _, v := range destinationMap["retry_on_status_codes"].([]interface{}) {
						retryOnCodes = append(retryOnCodes, v.(uint32))
					}
					destination.RetryOnStatusCodes = retryOnCodes
					var requestMap *consulapi.HTTPHeaderModifiers
					requestMap = new(consulapi.HTTPHeaderModifiers)
					addMap := make(map[string]string)
					setMap := make(map[string]string)
					if destinationMap["request_headers"] != nil {
						reqHeadersList := destinationMap["request_headers"].(*schema.Set).List()
						if len(reqHeadersList) > 0 {
							destinationAddMap := reqHeadersList[0].(map[string]interface{})["add"]
							if destinationAddMap != nil {
								for k, v := range destinationAddMap.(map[string]string) {
									addMap[k] = v
								}
								requestMap.Add = addMap
							}
						}
					}
					if destinationMap["request_headers"] != nil {
						reqHeadersList := destinationMap["request_headers"].(*schema.Set).List()
						if len(reqHeadersList) > 0 {
							destinationSetMap := reqHeadersList[0].(map[string]interface{})["set"]
							if destinationSetMap != nil {
								for k, v := range destinationSetMap.(map[string]string) {
									setMap[k] = v
								}
								requestMap.Set = setMap
							}
						}
					}
					removeList := make([]string, 0)
					if destinationMap["request_headers"] != nil {
						reqHeadersList := destinationMap["request_headers"].(*schema.Set).List()
						if len(reqHeadersList) > 0 {
							destinationRemoveList := reqHeadersList[0].(map[string]interface{})["remove"]
							if destinationRemoveList != nil && len(destinationRemoveList.([]string)) > 0 {
								for _, v := range destinationRemoveList.([]string) {
									removeList = append(removeList, v)
								}
							}
						}
					}
					if len(removeList) > 0 {
						requestMap.Remove = removeList
					}
					destination.RequestHeaders = requestMap
					var responseMap *consulapi.HTTPHeaderModifiers
					responseMap = new(consulapi.HTTPHeaderModifiers)
					addMap = make(map[string]string)
					setMap = make(map[string]string)
					if destinationMap["response_headers"] != nil {
						resHeadersList := destinationMap["response_headers"].(*schema.Set).List()
						if len(resHeadersList) > 0 {
							destinationAddMap := resHeadersList[0].(map[string]interface{})["add"]
							if destinationAddMap != nil {
								for k, v := range destinationAddMap.(map[string]string) {
									addMap[k] = v
								}
								responseMap.Add = addMap
							}
						}
					}
					if destinationMap["response_headers"] != nil {
						resHeadersList := destinationMap["response_headers"].(*schema.Set).List()
						if len(resHeadersList) > 0 {
							destinationSetMap := resHeadersList[0].(map[string]interface{})["set"]
							if destinationSetMap != nil {
								for k, v := range destinationSetMap.(map[string]string) {
									setMap[k] = v
								}
								responseMap.Set = setMap
							}
						}
					}
					removeList = make([]string, 0)
					if destinationMap["response_headers"] != nil {
						resHeadersList := destinationMap["response_headers"].(*schema.Set).List()
						if len(resHeadersList) > 0 {
							destinationRemoveList := resHeadersList[0].(map[string]interface{})["remove"]
							if destinationRemoveList != nil && len(destinationRemoveList.([]string)) > 0 {
								for _, v := range destinationRemoveList.([]string) {
									removeList = append(removeList, v)
								}
							}
						}
					}
					if len(removeList) > 0 {
						responseMap.Remove = removeList
					}
					destination.ResponseHeaders = responseMap
					serviceRoutesList[indx].Destination = destination
				}
			}
		}
		configEntry.Routes = serviceRoutesList
	}
	return configEntry, nil
}

func (s *serviceRouter) Write(ce consulapi.ConfigEntry, sw *stateWriter) error {
	sr, ok := ce.(*consulapi.ServiceRouterConfigEntry)
	if !ok {
		return fmt.Errorf("expected '%s' but got '%s'", consulapi.ServiceDefaults, ce.GetKind())
	}

	sw.set("name", sr.Name)
	sw.set("partition", sr.Partition)
	sw.set("namespace", sr.Namespace)

	meta := make(map[string]interface{})
	for k, v := range sr.Meta {
		meta[k] = v
	}
	sw.set("meta", meta)

	routes := make([]map[string]interface{}, 0)
	if len(sr.Routes) > 0 {
		route := make(map[string]interface{})
		for _, routesValue := range sr.Routes {
			match := make([]map[string]interface{}, 1)
			match[0] = make(map[string]interface{})
			matchHTTP := make([]map[string]interface{}, 1)
			matchHTTP[0] = make(map[string]interface{})
			matchHTTP[0]["path_exact"] = routesValue.Match.HTTP.PathExact
			matchHTTP[0]["path_prefix"] = routesValue.Match.HTTP.PathPrefix
			matchHTTP[0]["path_regex"] = routesValue.Match.HTTP.PathRegex
			matchHTTP[0]["methods"] = routesValue.Match.HTTP.Methods
			headerList := make([]map[string]interface{}, 0)
			for _, headerValue := range routesValue.Match.HTTP.Header {
				headerMap := make(map[string]interface{})
				headerMap["name"] = headerValue.Name
				headerMap["present"] = headerValue.Present
				headerMap["exact"] = headerValue.Exact
				headerMap["prefix"] = headerValue.Prefix
				headerMap["suffix"] = headerValue.Suffix
				headerMap["regex"] = headerValue.Regex
				headerMap["invert"] = headerValue.Invert
				headerList = append(headerList, headerMap)
			}
			queryParamList := make([]map[string]interface{}, 0)
			for _, queryParamValue := range routesValue.Match.HTTP.QueryParam {
				queryParamMap := make(map[string]interface{})
				queryParamMap["name"] = queryParamValue.Name
				queryParamMap["present"] = queryParamValue.Present
				queryParamMap["exact"] = queryParamValue.Exact
				queryParamMap["regex"] = queryParamValue.Regex
				queryParamList = append(queryParamList, queryParamMap)
			}
			matchHTTP[0]["header"] = headerList
			matchHTTP[0]["query_param"] = queryParamList
			match[0]["http"] = matchHTTP
			destination := make([]map[string]interface{}, 1)
			destination[0] = make(map[string]interface{})
			if routesValue.Destination != nil {
				destination[0]["service"] = routesValue.Destination.Service
				destination[0]["service_subset"] = routesValue.Destination.ServiceSubset
				destination[0]["namespace"] = routesValue.Destination.Namespace
				destination[0]["partition"] = routesValue.Destination.Partition
				destination[0]["prefix_rewrite"] = routesValue.Destination.PrefixRewrite
				destination[0]["request_timeout"] = routesValue.Destination.RequestTimeout.String()
				destination[0]["idle_timeout"] = routesValue.Destination.IdleTimeout.String()
				destination[0]["num_retries"] = routesValue.Destination.NumRetries
				destination[0]["retry_on_connect_failure"] = routesValue.Destination.RetryOnConnectFailure
				destination[0]["retry_on"] = routesValue.Destination.RetryOn
				destination[0]["retry_on_status_codes"] = routesValue.Destination.RetryOnStatusCodes
				requestHeaders := make([]map[string]interface{}, 1)
				requestHeaders[0] = make(map[string]interface{})
				addMap := make(map[string]interface{})
				if routesValue.Destination.RequestHeaders != nil && routesValue.Destination.RequestHeaders.Add != nil {
					for k, v := range routesValue.Destination.RequestHeaders.Add {
						addMap[k] = v
					}
				}
				requestHeaders[0]["add"] = addMap
				setMap := make(map[string]interface{})
				if routesValue.Destination.RequestHeaders != nil && routesValue.Destination.RequestHeaders.Set != nil {
					for k, v := range routesValue.Destination.RequestHeaders.Set {
						setMap[k] = v
					}
				}
				requestHeaders[0]["set"] = setMap
				removeList := make([]string, 0)
				if routesValue.Destination.RequestHeaders != nil && routesValue.Destination.RequestHeaders.Remove != nil {
					for _, v := range routesValue.Destination.RequestHeaders.Remove {
						removeList = append(removeList, v)
					}
				}
				if len(removeList) > 0 {
					requestHeaders[0]["remove"] = removeList
				}
				destination[0]["request_headers"] = requestHeaders
				responseHeaders := make([]map[string]interface{}, 1)
				responseHeaders[0] = make(map[string]interface{})
				responseHeaders[0]["add"] = make(map[string]interface{})
				addMap = make(map[string]interface{})
				if routesValue.Destination.ResponseHeaders != nil && routesValue.Destination.ResponseHeaders.Add != nil {
					for k, v := range routesValue.Destination.ResponseHeaders.Add {
						addMap[k] = v
					}
				}
				responseHeaders[0]["add"] = addMap
				setMap = make(map[string]interface{})
				if routesValue.Destination.ResponseHeaders != nil && routesValue.Destination.ResponseHeaders.Set != nil {
					for k, v := range routesValue.Destination.ResponseHeaders.Set {
						setMap[k] = v
					}
				}
				responseHeaders[0]["set"] = setMap
				removeList = make([]string, 0)
				if routesValue.Destination.ResponseHeaders != nil && routesValue.Destination.ResponseHeaders.Remove != nil {
					for _, v := range routesValue.Destination.ResponseHeaders.Remove {
						removeList = append(removeList, v)
					}
				}
				if len(removeList) > 0 {
					responseHeaders[0]["remove"] = removeList
				}
				destination[0]["response_headers"] = responseHeaders
			}
			route["match"] = match
			route["destination"] = destination
		}
		routes = append(routes, route)
	}

	sw.set("routes", routes)

	return nil
}
