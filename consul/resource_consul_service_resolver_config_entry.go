// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

type serviceResolver struct{}

func (s *serviceResolver) GetKind() string {
	return consulapi.ServiceResolver
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
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specifies names for custom service subsets and the conditions under which service instances belong to each subset.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of subset",
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
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specifies controls for rerouting traffic to an alternate pool of service instances if the target service fails.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"subset_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of subset",
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
						Type:        schema.TypeList,
						Optional:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
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

	requestTimeout, err := time.ParseDuration(d.Get("connect_timeout").(string))
	if err != nil {
		return nil, err
	}
	configEntry.RequestTimeout = requestTimeout

	subsets := make(map[string]consulapi.ServiceResolverSubset)

	subsetsList := d.Get("subsets").([]interface{})
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

	failoverList := d.Get("failover").([]interface{})
	failover := make(map[string]consulapi.ServiceResolverFailover)
	for _, failoverElem := range failoverList {
		failoverMap := failoverElem.(map[string]interface{})
		var serviceResolverFailover consulapi.ServiceResolverFailover
		serviceResolverFailover.Service = failoverMap["service"].(string)
		serviceResolverFailover.ServiceSubset = failoverMap["service_subset"].(string)
		serviceResolverFailover.Namespace = failoverMap["namespace"].(string)
		serviceResolverFailover.SamenessGroup = failoverMap["sameness_group"].(string)
		serviceResolverFailover.Datacenters = failoverMap["datacenter"].([]string)
		serviceResolverFailoverTargets := make([]consulapi.ServiceResolverFailoverTarget, len(failoverMap["targets"].([]interface{})))
		for indx, targetElem := range failoverMap["targets"].([]map[string]interface{}) {
			var serviceResolverFailoverTarget consulapi.ServiceResolverFailoverTarget
			serviceResolverFailoverTarget.Service = targetElem["service"].(string)
			serviceResolverFailoverTarget.ServiceSubset = targetElem["service_subset"].(string)
			serviceResolverFailoverTarget.Namespace = targetElem["namespace"].(string)
			serviceResolverFailoverTarget.Partition = targetElem["partition"].(string)
			serviceResolverFailoverTarget.Datacenter = targetElem["datacenter"].(string)
			serviceResolverFailoverTarget.Peer = targetElem["peer"].(string)
			serviceResolverFailoverTargets[indx] = serviceResolverFailoverTarget
		}
		serviceResolverFailover.Targets = serviceResolverFailoverTargets
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
			var hashPolicyList []consulapi.HashPolicy
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

func (s *serviceResolver) Write(ce consulapi.ConfigEntry, sw *stateWriter) error {
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
	sw.set("connect_timeout", sr.ConnectTimeout)

	return nil
}
