// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type serviceIntentions struct{}

func (s *serviceIntentions) GetKind() string {
	return consulapi.ServiceIntentions
}

func (s *serviceIntentions) GetDescription() string {
	return "The `consul_service_intentions_config_entry` resource configures [service intentions](https://developer.hashicorp.com/consul/docs/connect/config-entries/service-intentions) that are configurations for controlling access between services in the service mesh. A single service intentions configuration entry specifies one destination service and one or more L4 traffic sources, L7 traffic sources, or combination of traffic sources."
}

func (s *serviceIntentions) GetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Specifies a name of the destination service for all intentions defined in the configuration entry.",
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
		"jwt": {
			Type:        schema.TypeSet,
			Description: "Specifies a JSON Web Token provider configured in a JWT provider configuration entry, as well as additional configurations for verifying a service's JWT before authorizing communication between services",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"providers": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Specifies the names of one or more previously configured JWT provider configuration entries, which include the information necessary to validate a JSON web token.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Description: "Specifies the name of a JWT provider defined in the Name field of the jwt-provider configuration entry.",
									Optional:    true,
								},
								"verify_claims": {
									Type:        schema.TypeList,
									Description: "Specifies additional token information to verify beyond what is configured in the JWT provider configuration entry.",
									Optional:    true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"path": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "Specifies the path to the claim in the JSON web token.",
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
											"value": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Specifies the value to match on when verifying the the claim designated in path.",
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
		"sources": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of configurations that define intention sources and the authorization granted to the sources.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "Specifies the name of the source that the intention allows or denies traffic from.",
						Optional:    true,
					},
					"peer": {
						Type:        schema.TypeString,
						Description: "Specifies the name of a peered Consul cluster that the intention allows or denies traffic from",
						Optional:    true,
					},
					"namespace": {
						Type:        schema.TypeString,
						Description: "Specifies the traffic source namespace that the intention allows or denies traffic from.",
						Optional:    true,
					},
					"partition": {
						Type:        schema.TypeString,
						Description: "Specifies the name of an admin partition that the intention allows or denies traffic from.",
						Optional:    true,
					},
					"sameness_group": {
						Type:        schema.TypeString,
						Description: "Specifies the name of a sameness group that the intention allows or denies traffic from.",
						Optional:    true,
					},
					"action": {
						Type:        schema.TypeString,
						Description: "Specifies the action to take when the source sends traffic to the destination service.",
						Optional:    true,
					},
					"permissions": {
						Type:        schema.TypeList,
						Description: "Specifies a list of permissions for L7 traffic sources. The list contains one or more actions and a set of match criteria for each action.",
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"action": {
									Required:    true,
									Description: "Specifies the action to take when the source sends traffic to the destination service. The value is either allow or deny.",
									Type:        schema.TypeString,
								},
								"http": {
									Type:        schema.TypeSet,
									Required:    true,
									Description: "Specifies a set of HTTP-specific match criteria. ",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"path_exact": {
												Type:        schema.TypeString,
												Description: "Specifies an exact path to match on the HTTP request path.",
												Optional:    true,
											},
											"path_prefix": {
												Type:        schema.TypeString,
												Description: "Specifies a path prefix to match on the HTTP request path.",
												Optional:    true,
											},
											"path_regex": {
												Type:        schema.TypeString,
												Description: "Defines a regular expression to match on the HTTP request path.",
												Optional:    true,
											},
											"methods": {
												Type:        schema.TypeList,
												Description: "Specifies a list of HTTP methods.",
												Optional:    true,
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
											"headers": {
												Type:        schema.TypeList,
												Description: "Specifies a header name and matching criteria for HTTP request headers.",
												Optional:    true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"name": {
															Type:        schema.TypeString,
															Required:    true,
															Description: "Specifies the name of the header to match.",
														},
														"present": {
															Type:        schema.TypeBool,
															Default:     false,
															Optional:    true,
															Description: "Enables a match if the header configured in the Name field appears in the request. Consul matches on any value as long as the header key appears in the request.",
														},
														"exact": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies a value for the header key set in the Name field. If the request header value matches the Exact value, Consul applies the permission.",
														},
														"prefix": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies a prefix value for the header key set in the Name field.",
														},
														"suffix": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies a suffix value for the header key set in the Name field.",
														},
														"regex": {
															Type:        schema.TypeString,
															Optional:    true,
															Description: "Specifies a regular expression pattern as the value for the header key set in the Name field.",
														},
														"invert": {
															Type:        schema.TypeBool,
															Optional:    true,
															Description: "Inverts the matching logic configured in the Header.",
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
					"precedence": {
						Type:        schema.TypeInt,
						Description: "The Precedence field contains a read-only integer. Consul generates the value based on name configurations for the source and destination services.",
						Optional:    true,
					},
					"type": {
						Type:        schema.TypeString,
						Default:     "consul",
						Description: "Specifies the type of destination service that the configuration entry applies to.",
						Optional:    true,
					},
					"description": {
						Type:        schema.TypeString,
						Description: "Specifies a description of the intention.",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (s *serviceIntentions) Decode(d *schema.ResourceData) (consulapi.ConfigEntry, error) {
	configEntry := &consulapi.ServiceIntentionsConfigEntry{
		Kind: consulapi.ServiceIntentions,
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

	jwt := d.Get("jwt").(*schema.Set).List()

	if len(jwt) > 0 {
		jwtMap := jwt[0].(map[string]interface{})
		var jwtReq *consulapi.IntentionJWTRequirement
		jwtReq = new(consulapi.IntentionJWTRequirement)
		providers := make([]*consulapi.IntentionJWTProvider, 0)
		providerList := jwtMap["providers"].([]interface{})
		for _, pv := range providerList {
			pvm := pv.(map[string]interface{})
			var provider *consulapi.IntentionJWTProvider
			provider = new(consulapi.IntentionJWTProvider)
			if pvm["name"] != nil {
				provider.Name = pvm["name"].(string)
			}
			verifyClaims := make([]*consulapi.IntentionJWTClaimVerification, 0)
			if pvm["verify_claims"] != nil {
				verifyClaimsList := pvm["verify_claims"].([]interface{})
				for _, vcv := range verifyClaimsList {
					vcMap := vcv.(map[string]interface{})
					var verifyClaim *consulapi.IntentionJWTClaimVerification
					verifyClaim = new(consulapi.IntentionJWTClaimVerification)
					verifyClaimPath := make([]string, 0)
					for _, vcp := range vcMap["path"].([]interface{}) {
						verifyClaimPath = append(verifyClaimPath, vcp.(string))
					}
					verifyClaim.Path = verifyClaimPath
					if vcMap["value"] != nil {
						verifyClaim.Value = vcMap["value"].(string)
					}
					verifyClaims = append(verifyClaims, verifyClaim)
				}
				provider.VerifyClaims = verifyClaims
			}
			providers = append(providers, provider)
		}
		jwtReq.Providers = providers
		configEntry.JWT = jwtReq
	}

	sources := d.Get("sources")
	sourcesList := sources.([]interface{})

	if len(sourcesList) > 0 {
		sourcesIntentions := make([]*consulapi.SourceIntention, 0)
		for _, sr := range sourcesList {
			var sourceIntention *consulapi.SourceIntention
			sourceIntention = new(consulapi.SourceIntention)
			sourceMap := sr.(map[string]interface{})
			if sourceMap["name"] != nil {
				sourceIntention.Name = sourceMap["name"].(string)
			}
			if sourceMap["peer"] != nil {
				sourceIntention.Peer = sourceMap["peer"].(string)
			}
			if sourceMap["namespace"] != nil {
				sourceIntention.Namespace = sourceMap["namespace"].(string)
			}
			if sourceMap["partition"] != nil {
				sourceIntention.Partition = sourceMap["partition"].(string)
			}
			if sourceMap["sameness_group"] != nil {
				sourceIntention.SamenessGroup = sourceMap["sameness_group"].(string)
			}
			if sourceMap["action"] != nil {
				if sourceMap["action"].(string) == "allow" {
					sourceIntention.Action = consulapi.IntentionActionAllow
				} else if sourceMap["action"].(string) == "deny" {
					sourceIntention.Action = consulapi.IntentionActionDeny
				}
			}
			if sourceMap["permissions"] != nil {
				intentionPermissions := make([]*consulapi.IntentionPermission, 0)
				permissionList := sourceMap["permissions"].([]interface{})
				for _, permission := range permissionList {
					var intentionPermission *consulapi.IntentionPermission
					intentionPermission = new(consulapi.IntentionPermission)
					permissionMap := permission.(map[string]interface{})
					if permissionMap["action"] != nil && permissionMap["action"].(string) == "allow" {
						intentionPermission.Action = consulapi.IntentionActionAllow
					} else if permissionMap["action"] != nil && permissionMap["action"].(string) == "deny" {
						intentionPermission.Action = consulapi.IntentionActionDeny
					} else {
						return nil, fmt.Errorf("action is invalid. it should either be allow or deny")
					}
					if permissionMap["http"] != nil {
						var intentionPermissionHTTP *consulapi.IntentionHTTPPermission
						intentionPermissionHTTP = new(consulapi.IntentionHTTPPermission)
						httpMap := permissionMap["http"].(*schema.Set).List()
						if len(httpMap) > 0 {
							httpMapFirst := httpMap[0].(map[string]interface{})
							if httpMapFirst["path_exact"] != nil {
								intentionPermissionHTTP.PathExact = httpMapFirst["path_exact"].(string)
							}
							if httpMapFirst["path_prefix"] != nil {
								intentionPermissionHTTP.PathPrefix = httpMapFirst["path_prefix"].(string)
							}
							if httpMapFirst["path_regex"] != nil {
								intentionPermissionHTTP.PathPrefix = httpMapFirst["path_regex"].(string)
							}
							if httpMapFirst["methods"] != nil {
								httpMethods := make([]string, 0)
								for _, v := range httpMapFirst["methods"].([]interface{}) {
									httpMethods = append(httpMethods, v.(string))
								}
								intentionPermissionHTTP.Methods = httpMethods
							}
							intentionPermission.HTTP = intentionPermissionHTTP
							if httpMapFirst["headers"] != nil {
								httpHeaderPermissions := make([]consulapi.IntentionHTTPHeaderPermission, 0)
								for _, v := range httpMapFirst["headers"].([]interface{}) {
									var httpHeaderPermission consulapi.IntentionHTTPHeaderPermission
									headerPermissionMap := v.(map[string]interface{})
									if headerPermissionMap["name"] != nil {
										httpHeaderPermission.Name = headerPermissionMap["name"].(string)
									}
									if headerPermissionMap["present"] != nil {
										httpHeaderPermission.Present = headerPermissionMap["present"].(bool)
									}
									if headerPermissionMap["exact"] != nil {
										httpHeaderPermission.Exact = headerPermissionMap["exact"].(string)
									}
									if headerPermissionMap["prefix"] != nil {
										httpHeaderPermission.Prefix = headerPermissionMap["prefix"].(string)
									}
									if headerPermissionMap["suffix"] != nil {
										httpHeaderPermission.Suffix = headerPermissionMap["suffix"].(string)
									}
									if headerPermissionMap["regex"] != nil {
										httpHeaderPermission.Regex = headerPermissionMap["regex"].(string)
									}
									if headerPermissionMap["invert"] != nil {
										httpHeaderPermission.Invert = headerPermissionMap["invert"].(bool)
									}
									httpHeaderPermissions = append(httpHeaderPermissions, httpHeaderPermission)
								}
							}
						}
					}
					intentionPermissions = append(intentionPermissions, intentionPermission)
				}
				sourceIntention.Permissions = intentionPermissions
			}
			if sourceMap["precedence"] != nil {
				sourceIntention.Precedence = sourceMap["precedence"].(int)
			}
			if sourceMap["type"] != nil {
				typeSourceIntention := consulapi.IntentionSourceType(sourceMap["type"].(string))
				sourceIntention.Type = typeSourceIntention
			}
			if sourceMap["description"] != nil {
				sourceIntention.Description = sourceMap["description"].(string)
			}
			sourcesIntentions = append(sourcesIntentions, sourceIntention)
		}
		configEntry.Sources = sourcesIntentions
	}

	return configEntry, nil
}

func (s *serviceIntentions) Write(ce consulapi.ConfigEntry, d *schema.ResourceData, sw *stateWriter) error {
	si, ok := ce.(*consulapi.ServiceIntentionsConfigEntry)
	if !ok {
		return fmt.Errorf("expected '%s' but got '%s'", consulapi.ServiceRouter, ce.GetKind())
	}

	sw.set("name", si.Name)
	sw.set("partition", si.Partition)
	sw.set("namespace", si.Namespace)

	meta := map[string]interface{}{}
	for k, v := range si.Meta {
		meta[k] = v
	}
	sw.set("meta", meta)

	if si.JWT != nil {
		jwt := make([]map[string]interface{}, 1)
		jwt[0] = make(map[string]interface{})
		jwt[0]["providers"] = make([]map[string]interface{}, 0)
		jwtProviders := make([]map[string]interface{}, 0)
		for _, jwtProvider := range si.JWT.Providers {
			jwtProviderMap := make(map[string]interface{})
			jwtProviderMap["name"] = jwtProvider.Name
			verifyClaims := make([]map[string]interface{}, 0)
			for _, vc := range jwtProvider.VerifyClaims {
				vcMap := make(map[string]interface{})
				vcPaths := make([]string, 0)
				for _, p := range vc.Path {
					vcPaths = append(vcPaths, p)
				}
				vcMap["path"] = vcPaths
				vcMap["value"] = vc.Value
				verifyClaims = append(verifyClaims, vcMap)
			}
			jwtProviderMap["verify_claims"] = verifyClaims
			jwtProviders = append(jwtProviders, jwtProviderMap)
		}
		jwt[0]["providers"] = jwtProviders

		sw.set("jwt", jwt)
	}

	sources := make([]map[string]interface{}, 0)
	for _, source := range si.Sources {
		sourceMap := make(map[string]interface{})
		sourceMap["name"] = source.Name
		sourceMap["peer"] = source.Peer
		sourceMap["namespace"] = source.Namespace
		sourceMap["partition"] = source.Partition
		sourceMap["sameness_group"] = source.SamenessGroup
		sourceMap["action"] = source.Action
		sourceMap["precedence"] = source.Precedence
		sourceMap["type"] = source.Type
		sourceMap["description"] = source.Description
		permissions := make([]map[string]interface{}, 0)
		for _, permission := range source.Permissions {
			permissionMap := make(map[string]interface{})
			permissionMap["action"] = permission.Action
			permissionHttp := make([]map[string]interface{}, 1)
			permissionHttp[0] = make(map[string]interface{})
			permissionHttp[0]["path_exact"] = permission.HTTP.PathExact
			permissionHttp[0]["path_prefix"] = permission.HTTP.PathPrefix
			permissionHttp[0]["path_regex"] = permission.HTTP.PathRegex
			permissionHttp[0]["methods"] = permission.HTTP.Methods
			headers := make([]map[string]interface{}, 0)
			for _, header := range permission.HTTP.Header {
				headerMap := make(map[string]interface{})
				headerMap["name"] = header.Name
				headerMap["present"] = header.Present
				headerMap["exact"] = header.Exact
				headerMap["prefix"] = header.Prefix
				headerMap["suffix"] = header.Suffix
				headerMap["regex"] = header.Regex
				headerMap["invert"] = header.Invert
				headers = append(headers, headerMap)
			}
			permissionHttp[0]["headers"] = headers
			permissionMap["http"] = permissionHttp
			permissions = append(permissions, permissionMap)
		}
		sourceMap["permissions"] = permissions
		sources = append(sources, sourceMap)
	}

	sw.set("sources", sources)

	return nil
}
