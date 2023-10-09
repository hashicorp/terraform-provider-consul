// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type serviceDefaults struct{}

func (s *serviceDefaults) GetKind() string {
	return consulapi.ServiceDefaults
}

func (s *serviceDefaults) GetDescription() string {
	return "The `consul_config_entry_service_defaults` resource configures a [service defaults](https://developer.hashicorp.com/consul/docs/connect/config-entries/service-defaults) that contains common configuration settings for service mesh services, such as upstreams and gateways."
}

func (s *serviceDefaults) GetSchema() map[string]*schema.Schema {
	upstreamConfigSchemaOverrides := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the service you are setting the defaults for.",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the name of the Consul admin partition that the configuration entry applies to.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the namespace containing the upstream service that the configuration applies to.",
			},
			"peer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the peer name of the upstream service that the configuration applies to.",
			},
			"envoy_listener_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the default protocol for the service.",
			},
			"connect_timeout_ms": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limits": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Map that specifies a set of limits to apply to when connecting upstream services.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_connections": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum number of connections a service instance can establish against the upstream.",
						},
						"max_pending_requests": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum number of requests that are queued while waiting for a connection to establish.",
						},
						"max_concurrent_requests": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum number of concurrent requests.",
						},
					},
				},
			},
			"passive_health_check": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Map that specifies a set of rules that enable Consul to remove hosts from the upstream cluster that are unreachable or that return errors.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the time between checks.",
						},
						"max_failures": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the number of consecutive failures allowed per check interval. If exceeded, Consul removes the host from the load balancer.",
						},
						"enforcing_consecutive_5xx": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies a percentage that indicates how many times out of 100 that Consul ejects the host when it detects an outlier status.",
						},
						"max_ejection_percent": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum percentage of an upstream cluster that Consul ejects when the proxy reports an outlier.",
						},
						"base_ejection_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the minimum amount of time that an ejected host must remain outside the cluster before rejoining.",
						},
					},
				},
			},
			"mesh_gateway": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Specifies the default mesh gateway mode field for all upstreams.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"balance_outbound_connections": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the strategy for allocating outbound connections from upstreams across Envoy proxy threads.",
			},
		},
	}
	upstreamConfigSchemaDefaults := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the default protocol for the service.",
			},
			"connect_timeout_ms": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limits": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Map that specifies a set of limits to apply to when connecting upstream services.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_connections": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum number of connections a service instance can establish against the upstream.",
						},
						"max_pending_requests": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum number of requests that are queued while waiting for a connection to establish.",
						},
						"max_concurrent_requests": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum number of concurrent requests.",
						},
					},
				},
			},
			"passive_health_check": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Map that specifies a set of rules that enable Consul to remove hosts from the upstream cluster that are unreachable or that return errors.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the time between checks.",
						},
						"max_failures": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the number of consecutive failures allowed per check interval. If exceeded, Consul removes the host from the load balancer.",
						},
						"enforcing_consecutive_5xx": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies a percentage that indicates how many times out of 100 that Consul ejects the host when it detects an outlier status.",
						},
						"max_ejection_percent": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum percentage of an upstream cluster that Consul ejects when the proxy reports an outlier.",
						},
						"base_ejection_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the minimum amount of time that an ejected host must remain outside the cluster before rejoining.",
						},
					},
				},
			},
			"mesh_gateway": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Specifies the default mesh gateway mode field for all upstreams.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"balance_outbound_connections": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sets the strategy for allocating outbound connections from upstreams across Envoy proxy threads.",
			},
		},
	}
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Specifies the name of the service you are setting the defaults for.",
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Specifies the Consul namespace that the configuration entry applies to.",
		},
		"partition": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Specifies the name of the name of the Consul admin partition that the configuration entry applies to. Refer to Admin Partitions for additional information.",
		},
		"meta": {
			Type:        schema.TypeMap,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Specifies a set of custom key-value pairs to add to the Consul KV store.",
		},
		"protocol": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies the default protocol for the service.",
		},
		"balance_inbound_connections": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the strategy for allocating inbound connections to the service across Envoy proxy threads.",
		},
		"mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies a mode for how the service directs inbound and outbound traffic.",
		},
		"upstream_config": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Controls default upstream connection settings and custom overrides for individual upstream services.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"overrides": {
						Type:        schema.TypeList,
						Optional:    true,
						Elem:        upstreamConfigSchemaOverrides,
						Description: "Specifies options that override the default upstream configurations for individual upstreams.",
					},
					"defaults": {
						Type:        schema.TypeSet,
						Optional:    true,
						Elem:        upstreamConfigSchemaDefaults,
						Description: "Specifies configurations that set default upstream settings. For information about overriding the default configurations for in for individual upstreams, refer to UpstreamConfig.Overrides.",
					},
				},
			},
		},
		"transparent_proxy": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Controls configurations specific to proxies in transparent mode. Refer to Transparent Proxy Mode for additional information.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"outbound_listener_port": {
						Required: true,
						Type:     schema.TypeInt,
					},
					"dialed_directly": {
						Required: true,
						Type:     schema.TypeBool,
					},
				},
			},
		},
		"mutual_tls_mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Controls whether mutual TLS is required for incoming connections to this service. This setting is only supported for services with transparent proxy enabled.",
		},
		"envoy_extensions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of extensions to modify Envoy proxy configuration.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"required": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"arguments": {
						Type:     schema.TypeMap,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"consul_version": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"envoy_version": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"destination": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Configures the destination for service traffic through terminating gateways.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"port": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"addresses": {
						Type:     schema.TypeList,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		"local_connect_timeout_ms": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Specifies the number of milliseconds allowed for establishing connections to the local application instance before timing out.",
		},
		"max_inbound_connections": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Specifies the maximum number of concurrent inbound connections to each service instance.",
		},
		"local_request_timeout_ms": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Specifies the timeout for HTTP requests to the local application instance.",
		},
		"mesh_gateway": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Specifies the default mesh gateway mode field for the service.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"mode": {
						Required: true,
						Type:     schema.TypeString,
					},
				},
			},
		},
		"external_sni": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the TLS server name indication (SNI) when federating with an external system.",
		},
		"expose": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "Specifies default configurations for exposing HTTP paths through Envoy.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"checks": {
						Type:     schema.TypeBool,
						Optional: true,
						ForceNew: true,
					},
					"paths": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"path": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"local_path_port": {
									Type:     schema.TypeInt,
									Optional: true,
								},
								"listener_port": {
									Type:     schema.TypeInt,
									Optional: true,
								},
								"protocol": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s *serviceDefaults) Decode(d *schema.ResourceData) (consulapi.ConfigEntry, error) {
	configEntry := &consulapi.ServiceConfigEntry{
		Kind:      consulapi.ServiceDefaults,
		Name:      d.Get("name").(string),
		Partition: d.Get("partition").(string),
	}
	if d.Get("namespace") != nil {
		configEntry.Namespace = d.Get("namespace").(string)
	}
	if d.Get("protocol") != nil {
		configEntry.Protocol = d.Get("protocol").(string)
	}
	if d.Get("balance_inbound_connections") != nil {
		configEntry.BalanceInboundConnections = d.Get("balance_inbound_connections").(string)
	}
	if d.Get("local_connect_timeout_ms") != nil {
		configEntry.LocalConnectTimeoutMs = d.Get("local_connect_timeout_ms").(int)
	}
	if d.Get("max_inbound_connections") != nil {
		configEntry.MaxInboundConnections = d.Get("max_inbound_connections").(int)
	}
	if d.Get("local_request_timeout_ms") != nil {
		configEntry.LocalRequestTimeoutMs = d.Get("local_request_timeout_ms").(int)
	}
	if d.Get("external_sni") != nil {
		configEntry.ExternalSNI = d.Get("external_sni").(string)
	}
	if d.Get("mode") != nil {
		configEntry.Mode = consulapi.ProxyMode(d.Get("mode").(string))
	}
	if d.Get("mutual_tls_mode") != nil {
		configEntry.MutualTLSMode = consulapi.MutualTLSMode(d.Get("mutual_tls_mode").(string))
	}
	configEntry.Meta = map[string]string{}

	if d.Get("meta") != nil {
		for k, v := range d.Get("meta").(map[string]interface{}) {
			configEntry.Meta[k] = v.(string)
		}
	}

	getLimits := func(limitsMap map[string]interface{}) *consulapi.UpstreamLimits {
		intPtr := func(i int) *int {
			if i == 0 {
				return nil
			}
			return &i
		}
		return &consulapi.UpstreamLimits{
			MaxPendingRequests:    intPtr(limitsMap["max_pending_requests"].(int)),
			MaxConnections:        intPtr(limitsMap["max_connections"].(int)),
			MaxConcurrentRequests: intPtr(limitsMap["max_concurrent_requests"].(int)),
		}
	}

	getPassiveHealthCheck := func(passiveHealthCheckSet interface{}) (*consulapi.PassiveHealthCheck, error) {
		passiveHealthCheck := passiveHealthCheckSet.(*schema.Set).List()
		if len(passiveHealthCheck) > 0 {
			passiveHealthCheckMap := passiveHealthCheck[0].(map[string]interface{})
			uint32Ptr := func(i int) *uint32 {
				if i == 0 {
					return nil
				}
				ui := uint32(i)
				return &ui
			}
			passiveHealthCheck := &consulapi.PassiveHealthCheck{
				MaxFailures:             uint32(passiveHealthCheckMap["max_failures"].(int)),
				EnforcingConsecutive5xx: uint32Ptr(passiveHealthCheckMap["enforcing_consecutive_5xx"].(int)),
				MaxEjectionPercent:      uint32Ptr(passiveHealthCheckMap["max_ejection_percent"].(int)),
			}
			duration, err := time.ParseDuration(passiveHealthCheckMap["interval"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to parse interval: %w", err)
			}
			passiveHealthCheck.Interval = duration
			baseEjectionTime, err := time.ParseDuration(passiveHealthCheckMap["base_ejection_time"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to parse base_ejection_time: %w", err)
			}
			passiveHealthCheck.BaseEjectionTime = &baseEjectionTime
			return passiveHealthCheck, nil
		}
		return nil, nil
	}

	getMeshGateway := func(meshGateway interface{}) *consulapi.MeshGatewayConfig {
		meshGatewayList := meshGateway.(*schema.Set).List()
		if len(meshGatewayList) > 0 {
			meshGatewayData := meshGatewayList[0].(map[string]interface{})
			return &consulapi.MeshGatewayConfig{
				Mode: consulapi.MeshGatewayMode(meshGatewayData["mode"].(string)),
			}
		}
		return nil
	}

	getUpstreamConfigOverrides := func(upstreamConfigMap map[string]interface{}) (*consulapi.UpstreamConfig, error) {
		upstreamConfig := &consulapi.UpstreamConfig{}
		if upstreamConfigMap["name"] != nil {
			upstreamConfig.Name = upstreamConfigMap["name"].(string)
		}
		if upstreamConfigMap["partition"] != nil {
			upstreamConfig.Partition = upstreamConfigMap["partition"].(string)
		}
		if upstreamConfigMap["namespace"] != nil {
			upstreamConfig.Namespace = upstreamConfigMap["namespace"].(string)
		}
		if upstreamConfigMap["peer"] != nil {
			upstreamConfig.Peer = upstreamConfigMap["peer"].(string)
		}
		if upstreamConfigMap["protocol"] != nil {
			upstreamConfig.Protocol = upstreamConfigMap["protocol"].(string)
		}
		if upstreamConfigMap["connect_timeout_ms"] != nil {
			upstreamConfig.ConnectTimeoutMs = upstreamConfigMap["connect_timeout_ms"].(int)
		}
		if upstreamConfigMap["limits"] != nil && len(upstreamConfigMap["limits"].(*schema.Set).List()) > 0 {
			upstreamConfig.Limits = getLimits(upstreamConfigMap["limits"].(*schema.Set).List()[0].(map[string]interface{}))
		}
		if upstreamConfigMap["passive_health_check"] != nil {
			passiveHealthCheck, err := getPassiveHealthCheck(upstreamConfigMap["passive_health_check"])
			if err != nil {
				return nil, err
			}
			upstreamConfig.PassiveHealthCheck = passiveHealthCheck
		}
		if upstreamConfigMap["mesh_gateway"] != nil {
			upstreamConfig.MeshGateway = *getMeshGateway(upstreamConfigMap["mesh_gateway"])
		}
		if upstreamConfigMap["balance_outbound_connections"] != nil {
			upstreamConfig.BalanceOutboundConnections = upstreamConfigMap["balance_outbound_connections"].(string)
		}
		return upstreamConfig, nil
	}

	getUpstreamConfigDefaults := func(upstreamConfigMap map[string]interface{}) (*consulapi.UpstreamConfig, error) {
		upstreamConfig := &consulapi.UpstreamConfig{}
		if upstreamConfigMap["protocol"] != nil {
			upstreamConfig.Protocol = upstreamConfigMap["protocol"].(string)
		}
		if upstreamConfigMap["connect_timeout_ms"] != nil {
			upstreamConfig.ConnectTimeoutMs = upstreamConfigMap["connect_timeout_ms"].(int)
		}
		if upstreamConfigMap["limits"] != nil && len(upstreamConfigMap["limits"].(*schema.Set).List()) > 0 {
			upstreamConfig.Limits = getLimits(upstreamConfigMap["limits"].(*schema.Set).List()[0].(map[string]interface{}))
		}
		if upstreamConfigMap["passive_health_check"] != nil {
			passiveHealthCheck, err := getPassiveHealthCheck(upstreamConfigMap["passive_health_check"])
			if err != nil {
				return nil, err
			}
			upstreamConfig.PassiveHealthCheck = passiveHealthCheck
		}
		if upstreamConfigMap["mesh_gateway"] != nil {
			upstreamConfig.MeshGateway = *getMeshGateway(upstreamConfigMap["mesh_gateway"])
		}
		if upstreamConfigMap["balance_outbound_connections"] != nil {
			upstreamConfig.BalanceOutboundConnections = upstreamConfigMap["balance_outbound_connections"].(string)
		}
		return upstreamConfig, nil
	}

	getTransparentProxy := func(transparentProxy map[string]interface{}) *consulapi.TransparentProxyConfig {
		transparentProxyConfig := &consulapi.TransparentProxyConfig{}
		if transparentProxy["outbound_listener_port"] != nil {
			transparentProxyConfig.OutboundListenerPort = transparentProxy["outbound_listener_port"].(int)
		}
		if transparentProxy["dialed_directly"] != nil {
			transparentProxyConfig.DialedDirectly = transparentProxy["dialed_directly"].(bool)
		}
		return transparentProxyConfig
	}

	getEnvoyExtension := func(envoyExtensionMap map[string]interface{}) consulapi.EnvoyExtension {
		envoyExtension := consulapi.EnvoyExtension{}
		if envoyExtensionMap["name"] != nil {
			envoyExtension.Name = envoyExtensionMap["name"].(string)
		}
		if envoyExtensionMap["required"] != nil {
			envoyExtension.Required = envoyExtensionMap["required"].(bool)
		}
		if envoyExtensionMap["arguments"] != nil {
			envoyExtension.Arguments = envoyExtensionMap["arguments"].(map[string]interface{})
		}
		if envoyExtensionMap["consul_version"] != nil {
			envoyExtension.ConsulVersion = envoyExtensionMap["consul_version"].(string)
		}
		if envoyExtensionMap["envoy_version"] != nil {
			envoyExtension.EnvoyVersion = envoyExtensionMap["envoy_version"].(string)
		}
		return envoyExtension
	}

	getDestination := func(destinationMap map[string]interface{}) *consulapi.DestinationConfig {
		if destinationMap != nil {
			destination := &consulapi.DestinationConfig{}
			if destinationMap["port"] != nil {
				destination.Port = destinationMap["port"].(int)
			}
			if destinationMap["addresses"] != nil {
				for _, addr := range destinationMap["addresses"].([]interface{}) {
					destination.Addresses = append(destination.Addresses, addr.(string))
				}
			}
			return destination
		}
		return nil
	}

	getExposePath := func(exposePathMap map[string]interface{}) *consulapi.ExposePath {
		exposePath := &consulapi.ExposePath{}
		if exposePathMap["path"] != nil {
			exposePath.Path = exposePathMap["path"].(string)
		}
		if exposePathMap["local_path_port"] != nil {
			exposePath.LocalPathPort = exposePathMap["local_path_port"].(int)
		}
		if exposePathMap["listener_port"] != nil {
			exposePath.ListenerPort = exposePathMap["listener_port"].(int)
		}
		if exposePathMap["protocol"] != nil {
			exposePath.Protocol = exposePathMap["protocol"].(string)
		}
		return exposePath
	}

	getExpose := func(exposeMap map[string]interface{}) consulapi.ExposeConfig {
		exposeConfig := consulapi.ExposeConfig{}
		if exposeMap["checks"] != nil {
			exposeConfig.Checks = exposeMap["checks"].(bool)
		}
		if exposeMap["paths"] != nil {
			for _, elem := range exposeMap["paths"].([]interface{}) {
				exposeConfig.Paths = append(exposeConfig.Paths, *getExposePath(elem.(map[string]interface{})))
			}
		}
		return exposeConfig
	}

	upstreamConfigList := d.Get("upstream_config").(*schema.Set).List()
	if len(upstreamConfigList) > 0 {
		configEntry.UpstreamConfig = &consulapi.UpstreamConfiguration{}
		upstreamConfigMap := upstreamConfigList[0].(map[string]interface{})
		if upstreamConfigMap["defaults"] != nil {
			defaultsUpstreamConfigMapList := upstreamConfigMap["defaults"].(*schema.Set).List()
			if len(defaultsUpstreamConfigMapList) > 0 {
				defaultsUpstreamConfig, err := getUpstreamConfigDefaults(defaultsUpstreamConfigMapList[0].(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				configEntry.UpstreamConfig.Defaults = defaultsUpstreamConfig
			}
		}

		if upstreamConfigMap["overrides"] != nil {
			overrideUpstreamConfigList := upstreamConfigMap["overrides"].([]interface{})
			var overrideUpstreamConfig []*consulapi.UpstreamConfig
			for _, elem := range overrideUpstreamConfigList {
				overrideUpstreamConfigElem, err := getUpstreamConfigOverrides(elem.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				overrideUpstreamConfig = append(overrideUpstreamConfig, overrideUpstreamConfigElem)
			}
			configEntry.UpstreamConfig.Overrides = overrideUpstreamConfig
		}
	}

	if d.Get("transparent_proxy") != nil {
		transparentProxyList := d.Get("transparent_proxy").(*schema.Set).List()
		if len(transparentProxyList) > 0 {
			transparentProxy := transparentProxyList[0].(map[string]interface{})
			configEntry.TransparentProxy = getTransparentProxy(transparentProxy)
		}
	}

	if d.Get("envoy_extensions") != nil {
		for _, elem := range d.Get("envoy_extensions").([]interface{}) {
			envoyExtensionMap := elem.(map[string]interface{})
			configEntry.EnvoyExtensions = append(configEntry.EnvoyExtensions, getEnvoyExtension(envoyExtensionMap))
		}
	}

	if d.Get("destination") != nil {
		destinationList := d.Get("destination").(*schema.Set).List()
		if len(destinationList) > 0 {
			configEntry.Destination = getDestination(destinationList[0].(map[string]interface{}))
		}
	}

	if d.Get("mesh_gateway") != nil {
		if v := getMeshGateway(d.Get("mesh_gateway")); v != nil {
			configEntry.MeshGateway = *v
		}
	}

	if d.Get("expose") != nil {
		exposeList := d.Get("expose").(*schema.Set).List()
		if len(exposeList) > 0 {
			configEntry.Expose = getExpose(exposeList[0].(map[string]interface{}))
		}
	}

	return configEntry, nil
}

func (s *serviceDefaults) Write(ce consulapi.ConfigEntry, d *schema.ResourceData, sw *stateWriter) error {
	sd, ok := ce.(*consulapi.ServiceConfigEntry)
	if !ok {
		return fmt.Errorf("expected '%s' but got '%s'", consulapi.ServiceDefaults, ce.GetKind())
	}

	sw.set("name", sd.Name)
	sw.set("partition", sd.Partition)
	sw.set("namespace", sd.Namespace)

	meta := make(map[string]interface{})
	for k, v := range sd.Meta {
		meta[k] = v
	}
	sw.set("meta", meta)

	sw.set("protocol", sd.Protocol)
	sw.set("balance_inbound_connections", sd.BalanceInboundConnections)
	sw.set("mode", sd.Mode)

	getUpstreamConfigOverrides := func(elem *consulapi.UpstreamConfig) map[string]interface{} {
		upstreamConfig := make(map[string]interface{})
		upstreamConfig["name"] = elem.Name
		upstreamConfig["partition"] = elem.Partition
		upstreamConfig["namespace"] = elem.Namespace
		upstreamConfig["peer"] = elem.Peer
		upstreamConfig["protocol"] = elem.Protocol
		upstreamConfig["connect_timeout_ms"] = elem.ConnectTimeoutMs
		limits := make([]map[string]interface{}, 1)
		limits[0] = make(map[string]interface{})
		limits[0]["max_connections"] = elem.Limits.MaxConnections
		limits[0]["max_pending_requests"] = elem.Limits.MaxPendingRequests
		limits[0]["max_concurrent_requests"] = elem.Limits.MaxConcurrentRequests
		upstreamConfig["limits"] = limits
		passiveHealthCheck := make([]map[string]interface{}, 1)
		passiveHealthCheck[0] = make(map[string]interface{})
		passiveHealthCheck[0]["interval"] = elem.PassiveHealthCheck.Interval.String()
		passiveHealthCheck[0]["max_failures"] = elem.PassiveHealthCheck.MaxFailures
		passiveHealthCheck[0]["enforcing_consecutive_5xx"] = elem.PassiveHealthCheck.EnforcingConsecutive5xx
		passiveHealthCheck[0]["max_ejection_percent"] = elem.PassiveHealthCheck.MaxEjectionPercent
		passiveHealthCheck[0]["base_ejection_time"] = elem.PassiveHealthCheck.BaseEjectionTime.String()
		upstreamConfig["passive_health_check"] = passiveHealthCheck

		meshGateway := make([]map[string]interface{}, 1)
		meshGateway[0] = make(map[string]interface{})
		meshGateway[0]["mode"] = elem.MeshGateway.Mode
		upstreamConfig["mesh_gateway"] = meshGateway

		upstreamConfig["balance_outbound_connections"] = elem.BalanceOutboundConnections
		return upstreamConfig
	}

	getUpstreamConfigDefaults := func(elem *consulapi.UpstreamConfig) map[string]interface{} {
		upstreamConfig := make(map[string]interface{})
		upstreamConfig["protocol"] = elem.Protocol
		upstreamConfig["connect_timeout_ms"] = elem.ConnectTimeoutMs
		limits := make([]map[string]interface{}, 1)
		limits[0] = make(map[string]interface{})
		limits[0]["max_connections"] = elem.Limits.MaxConnections
		limits[0]["max_pending_requests"] = elem.Limits.MaxPendingRequests
		limits[0]["max_concurrent_requests"] = elem.Limits.MaxConcurrentRequests
		upstreamConfig["limits"] = limits
		passiveHealthCheck := make([]map[string]interface{}, 1)
		passiveHealthCheck[0] = make(map[string]interface{})
		passiveHealthCheck[0]["interval"] = elem.PassiveHealthCheck.Interval.String()
		passiveHealthCheck[0]["max_failures"] = elem.PassiveHealthCheck.MaxFailures
		passiveHealthCheck[0]["enforcing_consecutive_5xx"] = elem.PassiveHealthCheck.EnforcingConsecutive5xx
		passiveHealthCheck[0]["max_ejection_percent"] = elem.PassiveHealthCheck.MaxEjectionPercent
		passiveHealthCheck[0]["base_ejection_time"] = elem.PassiveHealthCheck.BaseEjectionTime.String()
		upstreamConfig["passive_health_check"] = passiveHealthCheck

		meshGateway := make([]map[string]interface{}, 1)
		meshGateway[0] = make(map[string]interface{})
		meshGateway[0]["mode"] = elem.MeshGateway.Mode
		upstreamConfig["mesh_gateway"] = meshGateway

		upstreamConfig["balance_outbound_connections"] = elem.BalanceOutboundConnections
		return upstreamConfig
	}

	if sd.UpstreamConfig != nil {
		var overrides []interface{}
		for _, elem := range sd.UpstreamConfig.Overrides {
			overrides = append(overrides, getUpstreamConfigOverrides(elem))
		}

		upstreamConfig := make(map[string]interface{})
		upstreamConfig["overrides"] = overrides
		defaultsSlice := make([]map[string]interface{}, 1)
		defaultsSlice[0] = getUpstreamConfigDefaults(sd.UpstreamConfig.Defaults)
		upstreamConfig["defaults"] = defaultsSlice
		upstreamConfigSlice := make([]map[string]interface{}, 1)
		upstreamConfigSlice[0] = upstreamConfig
		sw.set("upstream_config", upstreamConfigSlice)
	}

	transparentProxy := make([]map[string]interface{}, 1)
	transparentProxy[0] = make(map[string]interface{})
	transparentProxy[0]["outbound_listener_port"] = sd.TransparentProxy.OutboundListenerPort
	transparentProxy[0]["dialed_directly"] = sd.TransparentProxy.DialedDirectly
	sw.set("transparent_proxy", transparentProxy)

	sw.set("mutual_tls_mode", sd.MutualTLSMode)

	getEnvoyExtension := func(elem consulapi.EnvoyExtension) map[string]interface{} {
		envoyExtension := make(map[string]interface{})
		envoyExtension["name"] = elem.Name
		envoyExtension["required"] = elem.Required
		arguments := make(map[string]interface{})
		for k, v := range elem.Arguments {
			arguments[k] = v
		}
		envoyExtension["arguments"] = arguments
		envoyExtension["consul_version"] = elem.ConsulVersion
		envoyExtension["envoy_version"] = elem.EnvoyVersion
		return envoyExtension
	}

	var envoyExtensions []map[string]interface{}
	for _, elem := range sd.EnvoyExtensions {
		envoyExtensions = append(envoyExtensions, getEnvoyExtension(elem))
	}
	sw.set("envoy_extensions", envoyExtensions)

	destination := make([]map[string]interface{}, 1)
	if sd.Destination != nil {
		destination[0] = make(map[string]interface{})
		destination[0]["port"] = sd.Destination.Port
		destination[0]["addresses"] = sd.Destination.Addresses
		sw.set("destination", destination)
	}

	sw.set("local_connect_timeout_ms", sd.LocalConnectTimeoutMs)
	sw.set("max_inbound_connections", sd.MaxInboundConnections)
	sw.set("local_request_timeout_ms", sd.LocalRequestTimeoutMs)

	meshGateway := make([]map[string]interface{}, 1)
	meshGateway[0] = make(map[string]interface{})
	meshGateway[0]["mode"] = sd.MeshGateway.Mode
	sw.set("mesh_gateway", meshGateway)

	sw.set("external_sni", sd.ExternalSNI)

	expose := make([]map[string]interface{}, 1)
	expose[0] = make(map[string]interface{})
	expose[0]["checks"] = sd.Expose.Checks
	var paths []map[string]interface{}
	for _, elem := range sd.Expose.Paths {
		path := make(map[string]interface{})
		path["path"] = elem.Path
		path["local_path_port"] = elem.LocalPathPort
		path["listener_port"] = elem.ListenerPort
		path["protocol"] = elem.Protocol
		paths = append(paths, path)
	}
	expose[0]["paths"] = paths
	sw.set("expose", expose)

	return nil
}
