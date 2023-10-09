// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type serviceResolver struct{}

func (s *serviceResolver) GetKind() string {
	return consulapi.ServiceResolver
}

func (s *serviceResolver) GetDescription() string {
	return "The `consul_config_entry_service_resolver` resource configures a [service resolver](https://developer.hashicorp.com/consul/docs/connect/config-entries/service-resolver) that creates named subsets of service instances and define their behavior when satisfying upstream requests."
}

func (s *serviceResolver) GetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Specifies a name for the configuration entry.",
		},
		"partition": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Specifies the admin partition that the service resolver applies to.",
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Specifies the namespace that the service resolver applies to.",
		},
		"meta": {
			Type:        schema.TypeMap,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Specifies key-value pairs to add to the KV store.",
		},
		"connect_timeout": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the timeout duration for establishing new network connections to this service.",
		},
		"request_timeout": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the timeout duration for receiving an HTTP response from this service.",
		},
		"subsets": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Specifies names for custom service subsets and the conditions under which service instances belong to each subset.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of subset.",
					},
					"filter": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Specifies an expression that filters the DNS elements of service instances that belong to the subset. If empty, all healthy instances of a service are returned.",
					},
					"only_passing": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Determines if instances that return a warning from a health check are allowed to resolve a request. When set to false, instances with passing and warning states are considered healthy. When set to true, only instances with a passing health check state are considered healthy.",
					},
				},
			},
		},
		"default_subset": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies a defined subset of service instances to use when no explicit subset is requested. If this parameter is not specified, Consul uses the unnamed default subset.",
		},
		"redirect": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Specifies redirect instructions for local service traffic so that services deployed to a different network location resolve the upstream request instead.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the name of a service at the redirect’s destination that resolves local upstream requests.",
					},
					"service_subset": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the name of a subset of services at the redirect’s destination that resolves local upstream requests. If empty, the default subset is used. If specified, you must also specify at least one of the following in the same Redirect map: Service, Namespace, andDatacenter.",
					},
					"namespace": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the namespace at the redirect’s destination that resolves local upstream requests.",
					},
					"partition": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the admin partition at the redirect’s destination that resolves local upstream requests.",
					},
					"sameness_group": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the sameness group at the redirect’s destination that resolves local upstream requests.",
					},
					"datacenter": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the datacenter at the redirect’s destination that resolves local upstream requests.",
					},
					"peer": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the cluster with an active cluster peering connection at the redirect’s destination that resolves local upstream requests.",
					},
				},
			},
		},
		"failover": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Specifies controls for rerouting traffic to an alternate pool of service instances if the target service fails.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"subset_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of subset.",
					},
					"service": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the name of the service to resolve at the failover location during a failover scenario.",
					},
					"service_subset": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the name of a subset of service instances to resolve at the failover location during a failover scenario.",
					},
					"namespace": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the namespace at the failover location where the failover services are deployed.",
					},
					"sameness_group": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the sameness group at the failover location where the failover services are deployed.",
					},
					"datacenters": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Description: "Specifies an ordered list of datacenters at the failover location to attempt connections to during a failover scenario. When Consul cannot establish a connection with the first datacenter in the list, it proceeds sequentially until establishing a connection with another datacenter.",
					},
					"targets": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Specifies a fixed list of failover targets to try during failover. This list can express complicated failover scenarios.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"service": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the service name to use for the failover target. If empty, the current service name is used.",
								},
								"service_subset": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the named subset to use for the failover target. If empty, the default subset for the requested service name is used.",
								},
								"namespace": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the namespace to use for the failover target. If empty, the default namespace is used.",
								},
								"partition": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the admin partition within the same datacenter to use for the failover target. If empty, the default partition is used.",
								},
								"datacenter": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the WAN federated datacenter to use for the failover target. If empty, the current datacenter is used.",
								},
								"peer": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the destination cluster peer to resolve the target service name from.",
								},
							},
						},
					},
				},
			},
			Set: resourceConsulConfigEntryServiceResolverFailoverSetHash,
		},
		"load_balancer": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Specifies the load balancing policy and configuration for services issuing requests to this upstream.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"policy": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specifies the type of load balancing policy for selecting a host. ",
					},
					"least_request_config": {
						Type:        schema.TypeSet,
						Optional:    true,
						Description: "Specifies configuration for the least_request policy type.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"choice_count": {
									Type:     schema.TypeInt,
									Optional: true,
								},
							},
						},
					},
					"ring_hash_config": {
						Type:        schema.TypeSet,
						Optional:    true,
						Description: "Specifies configuration for the ring_hash policy type.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"minimum_ring_size": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Determines the minimum number of entries in the hash ring.",
								},
								"maximum_ring_size": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Determines the maximum number of entries in the hash ring.",
								},
							},
						},
					},
					"hash_policies": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Specifies a list of hash policies to use for hashing load balancing algorithms. Consul evaluates hash policies individually and combines them so that identical lists result in the same hash.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"field": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the attribute type to hash on. You cannot specify the Field parameter if SourceIP is also configured.",
								},
								"field_value": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Specifies the value to hash, such as a header name, cookie name, or a URL query parameter name.",
								},
								"cookie_config": {
									Type:        schema.TypeSet,
									Optional:    true,
									Description: "Specifies additional configuration options for the cookie hash policy type.",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"session": {
												Type:        schema.TypeBool,
												Optional:    true,
												Description: "Directs Consul to generate a session cookie with no expiration.",
											},
											"ttl": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Specifies the TTL for generated cookies. Cannot be specified for session cookies.",
											},
											"path": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Specifies the path to set for the cookie.",
											},
										},
									},
								},
								"source_ip": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Determines if the hash type should be source IP address.",
								},
								"terminal": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Determines if Consul should stop computing the hash when multiple hash policies are present.",
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceConsulConfigEntryServiceResolverFailoverSetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["subset_name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["subset_name"].(string)))
	}
	if m["service"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["service"].(string)))
	}
	if m["service_subset"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["service_subset"].(string)))
	}
	if m["namespace"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["namespace"].(string)))
	}
	if m["sameness_group"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["sameness_group"].(string)))
	}
	if m["datacenters"] != nil {
		datacenters := make([]string, 0)
		if strings.HasPrefix(reflect.ValueOf(m["datacenters"]).String(), "<[]interface") {
			for _, v := range m["datacenters"].([]interface{}) {
				datacenters = append(datacenters, v.(string))
			}
		} else {
			for _, v := range m["datacenters"].([]string) {
				datacenters = append(datacenters, v)
			}
		}
		sort.Strings(datacenters)
		for _, v := range datacenters {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}
	if m["targets"] != nil {
		for _, target := range m["targets"].([]interface{}) {
			var keys []string
			targetMap := target.(map[string]interface{})
			for k, _ := range targetMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				buf.WriteString(fmt.Sprintf("%s:%s-", k, targetMap[k].(string)))
			}
		}
	}
	return hashcode.String(buf.String())
}

func (s *serviceResolver) Decode(d *schema.ResourceData) (consulapi.ConfigEntry, error) {
	configEntry := &consulapi.ServiceResolverConfigEntry{
		Kind:      consulapi.ServiceResolver,
		Name:      d.Get("name").(string),
		Partition: d.Get("partition").(string),
		Namespace: d.Get("namespace").(string),
		Meta:      map[string]string{},
	}

	for k, v := range d.Get("meta").(map[string]interface{}) {
		configEntry.Meta[k] = v.(string)
	}

	connectTimeout, err := time.ParseDuration(d.Get("connect_timeout").(string))
	if err != nil {
		return nil, err
	}
	configEntry.ConnectTimeout = connectTimeout

	requestTimeout, err := time.ParseDuration(d.Get("request_timeout").(string))
	if err != nil {
		return nil, err
	}
	configEntry.RequestTimeout = requestTimeout

	subsets := make(map[string]consulapi.ServiceResolverSubset)

	subsetsList := d.Get("subsets").(*schema.Set).List()
	for _, subset := range subsetsList {
		subsetMap := subset.(map[string]interface{})
		var serviceResolverSubset consulapi.ServiceResolverSubset
		serviceResolverSubset.Filter = subsetMap["filter"].(string)
		serviceResolverSubset.OnlyPassing = subsetMap["only_passing"].(bool)
		subsets[subsetMap["name"].(string)] = serviceResolverSubset
	}
	configEntry.Subsets = subsets

	configEntry.DefaultSubset = d.Get("default_subset").(string)

	if v := (d.Get("redirect").(*schema.Set)).List(); len(v) == 1 {
		redirectMap := v[0].(map[string]interface{})
		var serviceResolverRedirect *consulapi.ServiceResolverRedirect
		serviceResolverRedirect = new(consulapi.ServiceResolverRedirect)
		serviceResolverRedirect.Service = redirectMap["service"].(string)
		serviceResolverRedirect.ServiceSubset = redirectMap["service_subset"].(string)
		serviceResolverRedirect.Namespace = redirectMap["namespace"].(string)
		serviceResolverRedirect.Partition = redirectMap["partition"].(string)
		serviceResolverRedirect.SamenessGroup = redirectMap["sameness_group"].(string)
		serviceResolverRedirect.Datacenter = redirectMap["datacenter"].(string)
		serviceResolverRedirect.Peer = redirectMap["peer"].(string)
		configEntry.Redirect = serviceResolverRedirect
	}

	failoverList := d.Get("failover").(*schema.Set).List()
	failover := make(map[string]consulapi.ServiceResolverFailover)
	for _, failoverElem := range failoverList {
		failoverMap := failoverElem.(map[string]interface{})
		var serviceResolverFailover consulapi.ServiceResolverFailover
		if value, ok := failoverMap["service"]; ok {
			serviceResolverFailover.Service = value.(string)
		}
		if value, ok := failoverMap["service_subset"]; ok {
			serviceResolverFailover.ServiceSubset = value.(string)
		}
		if value, ok := failoverMap["namespace"]; ok {
			serviceResolverFailover.Namespace = value.(string)
		}
		if value, ok := failoverMap["sameness_group"]; ok {
			serviceResolverFailover.SamenessGroup = value.(string)
		}
		if value, ok := failoverMap["datacenters"]; ok {
			datacenters := make([]string, 0)
			for _, v := range value.([]interface{}) {
				datacenters = append(datacenters, v.(string))
			}
			serviceResolverFailover.Datacenters = datacenters
		}
		if (failoverMap["targets"] != nil) && len(failoverMap["targets"].([]interface{})) > 0 {
			serviceResolverFailoverTargets := make([]consulapi.ServiceResolverFailoverTarget, len(failoverMap["targets"].([]interface{})))
			for indx, target := range failoverMap["targets"].([]interface{}) {
				targetElem := target.(map[string]interface{})
				var serviceResolverFailoverTarget consulapi.ServiceResolverFailoverTarget
				if value, ok := targetElem["service"]; ok {
					serviceResolverFailoverTarget.Service = value.(string)
				}
				if value, ok := targetElem["service_subset"]; ok {
					serviceResolverFailoverTarget.ServiceSubset = value.(string)
				}
				if value, ok := targetElem["namespace"]; ok {
					serviceResolverFailoverTarget.Namespace = value.(string)
				}
				if value, ok := targetElem["partition"]; ok {
					serviceResolverFailoverTarget.Partition = value.(string)
				}
				if value, ok := targetElem["datacenter"]; ok {
					serviceResolverFailoverTarget.Datacenter = value.(string)
				}
				if value, ok := targetElem["peer"]; ok {
					serviceResolverFailoverTarget.Peer = value.(string)
				}
				serviceResolverFailoverTargets[indx] = serviceResolverFailoverTarget
			}
			serviceResolverFailover.Targets = serviceResolverFailoverTargets
		}
		failover[failoverMap["subset_name"].(string)] = serviceResolverFailover
	}
	configEntry.Failover = failover

	if lb := (d.Get("load_balancer").(*schema.Set)).List(); len(lb) == 1 {
		loadBalancer := lb[0].(map[string]interface{})
		var ceLoadBalancer *consulapi.LoadBalancer
		ceLoadBalancer = new(consulapi.LoadBalancer)
		ceLoadBalancer.Policy = loadBalancer["policy"].(string)
		if lrc := (loadBalancer["least_request_config"].(*schema.Set)).List(); len(lrc) == 1 {
			var lreqConfig *consulapi.LeastRequestConfig
			lreqConfig = new(consulapi.LeastRequestConfig)
			lreqConfig.ChoiceCount = uint32(((lrc[0].(map[string]interface{}))["choice_count"]).(int))
			ceLoadBalancer.LeastRequestConfig = lreqConfig
		}
		if rhc := (loadBalancer["ring_hash_config"].(*schema.Set)).List(); len(rhc) == 1 {
			var rhConfig *consulapi.RingHashConfig
			rhConfig = new(consulapi.RingHashConfig)
			rhConfig.MaximumRingSize = uint64(rhc[0].(map[string]interface{})["maximum_ring_size"].(int))
			rhConfig.MinimumRingSize = uint64(rhc[0].(map[string]interface{})["minimum_ring_size"].(int))
			ceLoadBalancer.RingHashConfig = rhConfig
		}
		if hp := loadBalancer["hash_policies"].([]interface{}); len(hp) > 0 {
			hashPolicyList := make([]consulapi.HashPolicy, len(hp))
			for indx, hashPolicy := range hp {
				hashPolicyMap := hashPolicy.(map[string]interface{})
				hashPolicyList[indx].Field = hashPolicyMap["field"].(string)
				hashPolicyList[indx].FieldValue = hashPolicyMap["field_value"].(string)
				var cookieConfig *consulapi.CookieConfig
				cookieConfig = new(consulapi.CookieConfig)
				if cc := hashPolicyMap["cookie_config"].(*schema.Set).List(); len(cc) == 1 {
					cookieConfigMap := cc[0].(map[string]interface{})
					cookieConfig.Path = cookieConfigMap["path"].(string)
					cookieConfig.Session = cookieConfigMap["session"].(bool)
					ttl, err := time.ParseDuration(cookieConfigMap["ttl"].(string))
					if err != nil {
						return nil, err
					}
					cookieConfig.TTL = ttl
					hashPolicyList[indx].CookieConfig = cookieConfig
				}
				hashPolicyList[indx].SourceIP = hashPolicyMap["source_ip"].(bool)
				hashPolicyList[indx].Terminal = hashPolicyMap["terminal"].(bool)
			}
			ceLoadBalancer.HashPolicies = hashPolicyList
		}
		configEntry.LoadBalancer = ceLoadBalancer
	}

	return configEntry, nil
}

func (s *serviceResolver) Write(ce consulapi.ConfigEntry, d *schema.ResourceData, sw *stateWriter) error {
	sr, ok := ce.(*consulapi.ServiceResolverConfigEntry)
	if !ok {
		return fmt.Errorf("expected '%s' but got '%s'", consulapi.ServiceResolver, ce.GetKind())
	}

	sw.set("name", sr.Name)
	sw.set("partition", sr.Partition)
	sw.set("namespace", sr.Namespace)

	meta := map[string]interface{}{}
	for k, v := range sr.Meta {
		meta[k] = v
	}
	sw.set("meta", meta)
	sw.set("connect_timeout", sr.ConnectTimeout.String())
	sw.set("request_timeout", sr.RequestTimeout.String())

	subsets := make([]map[string]interface{}, len(sr.Subsets))
	indx := 0
	for name, _ := range sr.Subsets {
		subsetMap := make(map[string]interface{})
		subsetMap["name"] = name
		subsetSt := sr.Subsets[name]
		subsetMap["filter"] = subsetSt.Filter
		subsetMap["only_passing"] = subsetSt.OnlyPassing
		subsets[indx] = subsetMap
		indx++
	}
	sw.set("subsets", subsets)

	sw.set("default_subset", sr.DefaultSubset)

	redirect := make([]map[string]interface{}, 1)
	if sr.Redirect != nil {
		redirect[0] = make(map[string]interface{})
		redirect[0]["service"] = sr.Redirect.Service
		redirect[0]["service_subset"] = sr.Redirect.ServiceSubset
		redirect[0]["namespace"] = sr.Redirect.Namespace
		redirect[0]["partition"] = sr.Redirect.Partition
		redirect[0]["sameness_group"] = sr.Redirect.SamenessGroup
		redirect[0]["datacenter"] = sr.Redirect.Datacenter
		redirect[0]["peer"] = sr.Redirect.Peer
		sw.set("redirect", redirect)
	}

	var failover *schema.Set
	failover = new(schema.Set)
	failover.F = resourceConsulConfigEntryServiceResolverFailoverSetHash
	for name, failoverValue := range sr.Failover {
		failoverMap := make(map[string]interface{})
		failoverMap["subset_name"] = name
		failoverMap["service"] = failoverValue.Service
		failoverMap["service_subset"] = failoverValue.ServiceSubset
		failoverMap["namespace"] = failoverValue.Namespace
		failoverMap["sameness_group"] = failoverValue.SamenessGroup
		if len(failoverValue.Datacenters) > 0 {
			failoverDatacenters := make([]string, len(failoverValue.Datacenters))
			for index, fd := range failoverValue.Datacenters {
				failoverDatacenters[index] = fd
			}
			failoverMap["datacenters"] = failoverDatacenters
		}
		failoverTargets := make([]interface{}, len(sr.Failover[name].Targets))
		for index, ft := range sr.Failover[name].Targets {
			failoverTargetMap := make(map[string]interface{})
			failoverTargetMap["service"] = ft.Service
			failoverTargetMap["service_subset"] = ft.ServiceSubset
			failoverTargetMap["namespace"] = ft.Namespace
			failoverTargetMap["partition"] = ft.Partition
			failoverTargetMap["datacenter"] = ft.Datacenter
			failoverTargetMap["peer"] = ft.Peer
			failoverTargets[index] = failoverTargetMap
		}
		failoverMap["targets"] = failoverTargets
		failover.Add(failoverMap)
	}
	sw.set("failover", failover)

	if sr.LoadBalancer != nil {
		loadBalancer := make([]map[string]interface{}, 1)
		loadBalancer[0] = make(map[string]interface{})
		loadBalancer[0]["policy"] = sr.LoadBalancer.Policy
		if sr.LoadBalancer.LeastRequestConfig != nil {
			leastRequestConfig := make([]map[string]interface{}, 1)
			leastRequestConfig[0] = make(map[string]interface{})
			leastRequestConfig[0]["choice_count"] = sr.LoadBalancer.LeastRequestConfig.ChoiceCount
			loadBalancer[0]["least_request_config"] = leastRequestConfig
		}
		if sr.LoadBalancer.RingHashConfig != nil {
			ringHashConfig := make([]map[string]interface{}, 1)
			ringHashConfig[0] = make(map[string]interface{})
			ringHashConfig[0]["minimum_ring_size"] = sr.LoadBalancer.RingHashConfig.MinimumRingSize
			ringHashConfig[0]["maximum_ring_size"] = sr.LoadBalancer.RingHashConfig.MaximumRingSize
			loadBalancer[0]["ring_hash_config"] = ringHashConfig
		}

		if sr.LoadBalancer.HashPolicies != nil {
			hashPolicyList := make([]map[string]interface{}, len(sr.LoadBalancer.HashPolicies))
			for index, hashPolicy := range sr.LoadBalancer.HashPolicies {
				hashPolicyMap := make(map[string]interface{})
				hashPolicyMap["field"] = hashPolicy.Field
				hashPolicyMap["field_value"] = hashPolicy.FieldValue
				if hashPolicy.CookieConfig != nil {
					cookieConfigSet := make([]map[string]interface{}, 1)
					cookieConfigSet[0] = make(map[string]interface{})
					cookieConfigSet[0]["session"] = hashPolicy.CookieConfig.Session
					cookieConfigSet[0]["ttl"] = hashPolicy.CookieConfig.TTL
					cookieConfigSet[0]["path"] = hashPolicy.CookieConfig.TTL
					hashPolicyMap["cookie_config"] = cookieConfigSet
				}
				hashPolicyMap["source_ip"] = hashPolicy.SourceIP
				hashPolicyMap["terminal"] = hashPolicy.Terminal
				hashPolicyList[index] = hashPolicyMap
			}
			loadBalancer[0]["hash_policies"] = hashPolicyList
		}
		sw.set("load_balancer", loadBalancer)
	}

	return nil
}
