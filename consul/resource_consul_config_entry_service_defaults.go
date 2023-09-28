// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

type serviceDefaults struct{}

func (s *serviceDefaults) GetKind() string {
	return consulapi.ServiceSplitter
}

func (s *serviceDefaults) GetDescription() string {
	return "The `consul_config_entry_service_defaults` resource configures a [service defaults](https://developer.hashicorp.com/consul/docs/connect/config-entries/service-defaults) that contains common configuration settings for service mesh services, such as upstreams and gateways."
}

func (s *serviceDefaults) GetSchema() map[string]*schema.Schema {
	upstreamConfigSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"partition": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"envoy_listener_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"envoy_cluster_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connect_timeout_ms": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limits": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_connections": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_pending_requests": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_concurrent_requests": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"passive_health_check": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_failures": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"enforcing_consecutive_5xx": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_ejection_percent": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"base_ejection_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"mesh_gateway": {
				Type:     schema.TypeSet,
				Optional: true,
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
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	return map[string]*schema.Schema{
		"kind": {
			Type:     schema.TypeString,
			Required: false,
			ForceNew: true,
			Computed: true,
		},

		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"namespace": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"partition": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The partition the config entry is associated with.",
		},

		"meta": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},

		"protocol": {
			Type:     schema.TypeString,
			Required: true,
		},

		"balance_inbound_connections": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"mode": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"upstream_config": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"overrides": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     upstreamConfigSchema,
					},
					"defaults": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     upstreamConfigSchema,
					},
				},
			},
		},

		"transparent_proxy": {
			Type:     schema.TypeSet,
			Optional: true,
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
			Type:     schema.TypeString,
			Optional: true,
		},

		"envoy_extensions": {
			Type:     schema.TypeList,
			Optional: true,
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
			Type:     schema.TypeSet,
			Optional: true,
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
			Type:     schema.TypeInt,
			Optional: true,
		},

		"max_inbound_connections": {
			Type:     schema.TypeInt,
			Optional: true,
		},

		"local_request_timeout_ms": {
			Type:     schema.TypeInt,
			Optional: true,
		},

		"mesh_gateway": {
			Type:     schema.TypeSet,
			Optional: true,
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
			Type:     schema.TypeString,
			Optional: true,
		},

		"expose": {
			Type:     schema.TypeSet,
			Required: true,
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
		Namespace: d.Get("namespace").(string),
		Meta:      map[string]string{},
	}

	for k, v := range d.Get("meta").(map[string]interface{}) {
		configEntry.Meta[k] = v.(string)
	}

	configEntry.Protocol = d.Get("protocol").(string)
	configEntry.BalanceInboundConnections = d.Get("balance_inbound_connections").(string)

	proxyMode := consulapi.ProxyMode(d.Get("mode").(string))
	configEntry.Mode = proxyMode

	getLimits := func(limitsMap map[string]interface{}) *consulapi.UpstreamLimits {
		upstreamLimit := &consulapi.UpstreamLimits{}
		upstreamLimit.MaxPendingRequests = limitsMap["max_pending_requests"].(*int)
		upstreamLimit.MaxConnections = limitsMap["max_connections"].(*int)
		upstreamLimit.MaxConcurrentRequests = limitsMap["max_concurrent_requests"].(*int)
		return upstreamLimit
	}

	getPassiveHealthCheck := func(passiveHealthCheckSet interface{}) (*consulapi.PassiveHealthCheck, error) {
		passiveHealthCheck := passiveHealthCheckSet.(*schema.Set).List()
		if len(passiveHealthCheck) > 0 {
			passiveHealthCheckMap := passiveHealthCheck[0].(map[string]interface{})
			passiveHealthCheck := &consulapi.PassiveHealthCheck{}
			duration, err := time.ParseDuration(passiveHealthCheckMap["interval"].(string))
			if err != nil {
				return nil, err
			}
			passiveHealthCheck.Interval = duration
			passiveHealthCheck.MaxFailures = passiveHealthCheckMap["max_failures"].(uint32)
			passiveHealthCheck.EnforcingConsecutive5xx = passiveHealthCheckMap["enforcing_consecutive_5xx"].(*uint32)
			passiveHealthCheck.MaxEjectionPercent = passiveHealthCheckMap["max_ejection_percent"].(*uint32)
			baseEjectionTime, err := time.ParseDuration(passiveHealthCheckMap["base_ejection_time"].(string))
			if err != nil {
				return nil, err
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
			meshGatewayMode := consulapi.MeshGatewayMode(meshGatewayData["mode"].(string))
			meshGateway := &consulapi.MeshGatewayConfig{}
			meshGateway.Mode = meshGatewayMode
			return meshGateway
		}
		return nil
	}

	getUpstreamConfig := func(upstreamConfigMap map[string]interface{}) (*consulapi.UpstreamConfig, error) {
		upstreamConfig := &consulapi.UpstreamConfig{}
		upstreamConfig.Name = upstreamConfigMap["name"].(string)
		upstreamConfig.Partition = upstreamConfigMap["partition"].(string)
		upstreamConfig.Namespace = upstreamConfigMap["namespace"].(string)
		upstreamConfig.Peer = upstreamConfigMap["peer"].(string)
		upstreamConfig.EnvoyListenerJSON = upstreamConfigMap["envoy_listener_json"].(string)
		upstreamConfig.EnvoyClusterJSON = upstreamConfigMap["envoy_cluster_json"].(string)
		upstreamConfig.Protocol = upstreamConfigMap["protocol"].(string)
		upstreamConfig.ConnectTimeoutMs = upstreamConfigMap["connect_timeout_ms"].(int)
		upstreamConfig.Limits = getLimits(upstreamConfigMap["limits"].(map[string]interface{}))
		passiveHealthCheck, err := getPassiveHealthCheck(upstreamConfigMap["passive_health_check"])
		if err != nil {
			return nil, err
		}
		upstreamConfig.PassiveHealthCheck = passiveHealthCheck
		upstreamConfig.MeshGateway = *getMeshGateway(upstreamConfigMap["mesh_gateway"])
		upstreamConfig.BalanceOutboundConnections = upstreamConfigMap["balance_outbound_connections"].(string)
		return upstreamConfig, nil
	}

	getTransparentProxy := func(transparentProxy map[string]interface{}) *consulapi.TransparentProxyConfig {
		tProxy := &consulapi.TransparentProxyConfig{}
		tProxy.OutboundListenerPort = transparentProxy["outbound_listener_port"].(int)
		tProxy.DialedDirectly = transparentProxy["dialed_directly"].(bool)
		return tProxy
	}

	getEnvoyExtension := func(envoyExtensionMap map[string]interface{}) *consulapi.EnvoyExtension {
		envoyExtension := &consulapi.EnvoyExtension{}
		envoyExtension.Name = envoyExtensionMap["name"].(string)
		envoyExtension.Required = envoyExtensionMap["required"].(bool)
		envoyExtension.Arguments = envoyExtensionMap["arguments"].(map[string]interface{})
		envoyExtension.ConsulVersion = envoyExtensionMap["consul_version"].(string)
		envoyExtension.EnvoyVersion = envoyExtensionMap["envoy_version"].(string)
		return envoyExtension
	}

	getDestination := func(destinationMap map[string]interface{}) *consulapi.DestinationConfig {
		var destination consulapi.DestinationConfig
		destination.Port = destinationMap["port"].(int)
		destination.Addresses = destinationMap["addresses"].([]string)
		return &destination
	}

	getExposePath := func(exposePathMap map[string]interface{}) *consulapi.ExposePath {
		exposePath := &consulapi.ExposePath{}
		exposePath.Path = exposePathMap["path"].(string)
		exposePath.LocalPathPort = exposePathMap["local_path_port"].(int)
		exposePath.ListenerPort = exposePathMap["listener_port"].(int)
		exposePath.Protocol = exposePathMap["protocol"].(string)
		return exposePath
	}

	getExpose := func(exposeMap map[string]interface{}) consulapi.ExposeConfig {
		var exposeConfig consulapi.ExposeConfig
		exposeConfig.Checks = exposeMap["checks"].(bool)
		var paths []consulapi.ExposePath
		for _, elem := range exposeMap["paths"].([]interface{}) {
			paths = append(paths, *getExposePath(elem.(map[string]interface{})))
		}
		exposeConfig.Paths = paths
		return exposeConfig
	}

	upstreamConfigMap := d.Get("upstream_config").(map[string]interface{})
	defaultsUpstreamConfigMapList := upstreamConfigMap["defaults"].(*schema.Set).List()
	if len(defaultsUpstreamConfigMapList) > 0 {
		defaultsUpstreamConfig, err := getUpstreamConfig(defaultsUpstreamConfigMapList[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		configEntry.UpstreamConfig.Defaults = defaultsUpstreamConfig
	}

	overrideUpstreamConfigList := upstreamConfigMap["overrides"].([]interface{})
	var overrideUpstreamConfig []*consulapi.UpstreamConfig
	for _, elem := range overrideUpstreamConfigList {
		overrideUpstreamConfigList := elem.(*schema.Set).List()
		if len(overrideUpstreamConfigList) > 0 {
			overrideUpstreamConfigElem, err := getUpstreamConfig(overrideUpstreamConfigList[0].(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			overrideUpstreamConfig = append(overrideUpstreamConfig, overrideUpstreamConfigElem)
		}
	}
	configEntry.UpstreamConfig.Overrides = overrideUpstreamConfig

	transparentProxyList := d.Get("transparent_proxy").(*schema.Set).List()
	if len(transparentProxyList) > 0 {
		transparentProxy := transparentProxyList[0].(map[string]interface{})
		configEntry.TransparentProxy = getTransparentProxy(transparentProxy)
	}

	mutualTLSmode := consulapi.MutualTLSMode(d.Get("mutual_tls_mode").(string))
	configEntry.MutualTLSMode = mutualTLSmode

	var envoyExtensions []consulapi.EnvoyExtension

	for _, elem := range d.Get("envoy_extensions").([]interface{}) {
		envoyExtensionMap := elem.(map[string]interface{})
		envoyExtensions = append(envoyExtensions, *getEnvoyExtension(envoyExtensionMap))
	}
	configEntry.EnvoyExtensions = envoyExtensions

	destinationList := d.Get("destination").(*schema.Set).List()
	if len(destinationList) > 0 {
		configEntry.Destination = getDestination(destinationList[0].(map[string]interface{}))
	}

	configEntry.LocalConnectTimeoutMs = d.Get("local_connect_timeout_ms").(int)
	configEntry.MaxInboundConnections = d.Get("max_inbound_connections").(int)
	configEntry.LocalRequestTimeoutMs = d.Get("local_request_timeout_ms").(int)

	configEntry.MeshGateway = *getMeshGateway(d.Get("mesh_gateway"))
	configEntry.ExternalSNI = d.Get("external_sni").(string)

	exposeList := d.Get("expose").(*schema.Set).List()
	if len(exposeList) > 0 {
		configEntry.Expose = getExpose(exposeList[0].(map[string]interface{}))
	}

	return configEntry, nil
}

func (s *serviceDefaults) Write(ce consulapi.ConfigEntry, sw *stateWriter) error {
	sp, ok := ce.(*consulapi.ServiceConfigEntry)
	if !ok {
		return fmt.Errorf("expected '%s' but got '%s'", consulapi.ServiceDefaults, ce.GetKind())
	}

	sw.set("name", sp.Name)
	sw.set("partition", sp.Partition)
	sw.set("namespace", sp.Partition)

	var meta map[string]interface{}
	for k, v := range sp.Meta {
		meta[k] = v
	}
	sw.set("meta", meta)

	sw.set("protocol", sp.Protocol)
	sw.set("balance_inbound_connections", sp.BalanceInboundConnections)
	sw.set("mode", sp.Mode)

	getUpstreamConfig := func(elem *consulapi.UpstreamConfig) map[string]interface{} {
		var upstreamConfig map[string]interface{}
		upstreamConfig["name"] = elem.Name
		upstreamConfig["partition"] = elem.Partition
		upstreamConfig["namespace"] = elem.Namespace
		upstreamConfig["peer"] = elem.Peer
		upstreamConfig["envoy_listener_json"] = elem.EnvoyListenerJSON
		upstreamConfig["envoy_cluster_json"] = elem.EnvoyClusterJSON
		upstreamConfig["protocol"] = elem.Protocol
		upstreamConfig["connect_timeout_ms"] = elem.ConnectTimeoutMs
		var limits map[string]interface{}
		limits["max_connections"] = elem.Limits.MaxConnections
		limits["max_pending_requests"] = elem.Limits.MaxPendingRequests
		limits["max_concurrent_requests"] = elem.Limits.MaxConcurrentRequests
		upstreamConfig["limits"] = limits
		var passiveHealthCheck map[string]interface{}
		passiveHealthCheck["interval"] = elem.PassiveHealthCheck.Interval
		passiveHealthCheck["max_failures"] = elem.PassiveHealthCheck.MaxFailures
		passiveHealthCheck["enforcing_consecutive_5xx"] = elem.PassiveHealthCheck.EnforcingConsecutive5xx
		passiveHealthCheck["max_ejection_percent"] = elem.PassiveHealthCheck.MaxEjectionPercent
		passiveHealthCheck["base_ejection_time"] = elem.PassiveHealthCheck.BaseEjectionTime
		upstreamConfig["passive_health_check"] = passiveHealthCheck

		var meshGateway map[string]interface{}
		meshGateway["mode"] = elem.MeshGateway.Mode
		upstreamConfig["mesh_gateway"] = meshGateway

		upstreamConfig["balance_outbound_connections"] = elem.BalanceOutboundConnections
		return upstreamConfig
	}

	var overrides []interface{}
	for _, elem := range sp.UpstreamConfig.Overrides {
		overrides = append(overrides, getUpstreamConfig(elem))
	}

	var upstreamConfig map[string]interface{}
	upstreamConfig["overrides"] = overrides
	upstreamConfig["defaults"] = getUpstreamConfig(sp.UpstreamConfig.Defaults)
	sw.set("upstream_config", upstreamConfig)

	var transparentProxy map[string]interface{}
	transparentProxy["outbound_listener_port"] = sp.TransparentProxy.OutboundListenerPort
	transparentProxy["dialed_directly"] = sp.TransparentProxy.DialedDirectly
	sw.set("transparent_proxy", transparentProxy)

	sw.set("mutual_tls_mode", sp.MutualTLSMode)

	getEnvoyExtension := func(elem consulapi.EnvoyExtension) map[string]interface{} {
		var envoyExtension map[string]interface{}
		envoyExtension["name"] = elem.Name
		envoyExtension["required"] = elem.Required
		var arguments map[string]interface{}
		for k, v := range elem.Arguments {
			arguments[k] = v
		}
		envoyExtension["arguments"] = arguments
		envoyExtension["consul_version"] = elem.ConsulVersion
		envoyExtension["envoy_version"] = elem.EnvoyVersion
		return envoyExtension
	}

	var envoyExtensions []map[string]interface{}
	for _, elem := range sp.EnvoyExtensions {
		envoyExtensions = append(envoyExtensions, getEnvoyExtension(elem))
	}
	sw.set("envoy_extensions", envoyExtensions)

	destination := make(map[string]interface{})
	destination["port"] = sp.Destination.Port
	destination["addresses"] = sp.Destination.Addresses
	sw.set("destination", destination)

	sw.set("local_connect_timeout_ms", sp.LocalConnectTimeoutMs)
	sw.set("max_inbound_connections", sp.MaxInboundConnections)
	sw.set("local_request_timeout_ms", sp.LocalRequestTimeoutMs)

	meshGateway := make(map[string]interface{})
	meshGateway["mode"] = sp.MeshGateway.Mode
	sw.set("mesh_gateway", meshGateway)

	sw.set("external_sni", sp.ExternalSNI)

	expose := make(map[string]interface{})
	expose["checks"] = sp.Expose.Checks
	var paths []map[string]interface{}
	for _, elem := range sp.Expose.Paths {
		path := make(map[string]interface{})
		path["path"] = elem.Path
		path["local_path_port"] = elem.LocalPathPort
		path["listener_port"] = elem.ListenerPort
		path["protocol"] = elem.Protocol
		paths = append(paths, path)
	}
	expose["paths"] = paths
	sw.set("expose", expose)

	return nil
}
