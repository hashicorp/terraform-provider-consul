// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/mapstructure"
)

var (
	tokenDeprecationMessage = `The token argument has been deprecated and will be removed in a future release.
Please use the token argument in the provider configuration`
)

func deprecated(name string, resource *schema.Resource) *schema.Resource {
	resource.DeprecationMessage = fmt.Sprintf("%s is deprecated and will be removed in a future version.", name)
	return resource
}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The datacenter to use. Defaults to that of the agent.",
			},

			"address": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CONSUL_ADDRESS",
					"CONSUL_HTTP_ADDR",
				}, "localhost:8500"),
				Description: `The HTTP(S) API address of the agent to use. Defaults to "127.0.0.1:8500".`,
			},

			"scheme": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CONSUL_SCHEME",
					"CONSUL_HTTP_SCHEME",
				}, ""),
				Description: `The URL scheme of the agent to use ("http" or "https"). Defaults to "http".`,
			},

			"http_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONSUL_HTTP_AUTH", nil),
				Description: "HTTP Basic Authentication credentials to be used when communicating with Consul, in the format of either `user` or `user:pass`. This may also be specified using the `CONSUL_HTTP_AUTH` environment variable.",
			},

			"ca_file": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("CONSUL_CA_FILE", nil),
				ConflictsWith: []string{"ca_pem"},
				Description:   "A path to a PEM-encoded certificate authority used to verify the remote agent's certificate.",
			},

			"ca_pem": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ca_file"},
				Description:   "PEM-encoded certificate authority used to verify the remote agent's certificate.",
			},

			"cert_file": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("CONSUL_CERT_FILE", nil),
				ConflictsWith: []string{"cert_pem"},
				Description:   "A path to a PEM-encoded certificate provided to the remote agent; requires use of `key_file` or `key_pem`.",
			},

			"cert_pem": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cert_file"},
				Description:   "PEM-encoded certificate provided to the remote agent; requires use of `key_file` or `key_pem`.",
			},

			"key_file": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("CONSUL_KEY_FILE", nil),
				ConflictsWith: []string{"key_pem"},
				Description:   "A path to a PEM-encoded private key, required if `cert_file` or `cert_pem` is specified.",
			},

			"key_pem": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key_file"},
				Description:   "PEM-encoded private key, required if `cert_file` or `cert_pem` is specified.",
			},

			"ca_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONSUL_CAPATH", ""),
				Description: "A path to a directory of PEM-encoded certificate authority files to use to check the authenticity of client and server connections. Can also be specified with the `CONSUL_CAPATH` environment variable.",
			},

			"insecure_https": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Boolean value to disable SSL certificate verification; setting this value to true is not recommended for production use. Only use this with scheme set to "https".`,
			},

			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CONSUL_TOKEN",
					"CONSUL_HTTP_TOKEN",
				}, nil),
				Description: "The ACL token to use by default when making requests to the agent. Can also be specified with `CONSUL_HTTP_TOKEN` or `CONSUL_TOKEN` as an environment variable.",
			},

			"auth_jwt": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Authenticates to Consul using a JWT authentication method.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth_method": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the auth method to use for login.",
						},
						"bearer_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The bearer token to present to the auth method during login for authentication purposes. For the Kubernetes auth method this is a [Service Account Token (JWT)](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#service-account-tokens).",
						},
						"use_terraform_cloud_workload_identity": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to use a [Terraform Workload Identity token](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials/workload-identity-tokens). The token will be read from the `TFC_WORKLOAD_IDENTITY_TOKEN` environment variable.",
						},
						"meta": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Specifies arbitrary KV metadata linked to the token. Can be useful to track origins.",
						},
					},
				},
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"header": {
				Type:        schema.TypeList,
				Optional:    true,
				Sensitive:   true,
				Description: "A configuration block, described below, that provides additional headers to be sent along with all requests to the Consul server. This block can be specified multiple times.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the header.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the header.",
						},
					},
				},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"consul_agent_self":           dataSourceConsulAgentSelf(),
			"consul_agent_config":         dataSourceConsulAgentConfig(),
			"consul_autopilot_health":     dataSourceConsulAutopilotHealth(),
			"consul_nodes":                dataSourceConsulNodes(),
			"consul_service":              dataSourceConsulService(),
			"consul_service_health":       dataSourceConsulServiceHealth(),
			"consul_services":             dataSourceConsulServices(),
			"consul_keys":                 dataSourceConsulKeys(),
			"consul_key_prefix":           dataSourceConsulKeyPrefix(),
			"consul_acl_auth_method":      dataSourceConsulACLAuthMethod(),
			"consul_acl_policy":           dataSourceConsulACLPolicy(),
			"consul_acl_role":             dataSourceConsulACLRole(),
			"consul_acl_token":            dataSourceConsulACLToken(),
			"consul_acl_token_secret_id":  dataSourceConsulACLTokenSecretID(),
			"consul_network_segments":     dataSourceConsulNetworkSegments(),
			"consul_network_area_members": dataSourceConsulNetworkAreaMembers(),
			"consul_datacenters":          dataSourceConsulDatacenters(),
			"consul_config_entry":         dataSourceConsulConfigEntry(),
			"consul_peering":              dataSourceConsulPeering(),
			"consul_peerings":             dataSourceConsulPeerings(),

			// Aliases to limit the impact of rename of catalog
			// datasources
			"consul_catalog_nodes":    deprecated("consul_catalog_nodes", dataSourceConsulNodes()),
			"consul_catalog_service":  deprecated("consul_catalog_service", dataSourceConsulService()),
			"consul_catalog_services": deprecated("consul_catalog_services", dataSourceConsulServices()),
		},

		ResourcesMap: map[string]*schema.Resource{
			"consul_acl_auth_method":                 resourceConsulACLAuthMethod(),
			"consul_acl_binding_rule":                resourceConsulACLBindingRule(),
			"consul_acl_policy":                      resourceConsulACLPolicy(),
			"consul_acl_role_policy_attachment":      resourceConsulACLRolePolicyAttachment(),
			"consul_acl_role":                        resourceConsulACLRole(),
			"consul_acl_token_policy_attachment":     resourceConsulACLTokenPolicyAttachment(),
			"consul_acl_token_role_attachment":       resourceConsulACLTokenRoleAttachment(),
			"consul_acl_token":                       resourceConsulACLToken(),
			"consul_admin_partition":                 resourceConsulAdminPartition(),
			"consul_agent_service":                   resourceConsulAgentService(),
			"consul_autopilot_config":                resourceConsulAutopilotConfig(),
			"consul_catalog_entry":                   resourceConsulCatalogEntry(),
			"consul_certificate_authority":           resourceConsulCertificateAuthority(),
			"consul_config_entry_service_defaults":   resourceFromConfigEntryImplementation(&serviceDefaults{}),
			"consul_config_entry_service_intentions": resourceFromConfigEntryImplementation(&serviceIntentions{}),
			"consul_config_entry_service_splitter":   resourceFromConfigEntryImplementation(&serviceSplitter{}),
			"consul_config_entry":                    resourceConsulConfigEntry(),
			"consul_intention":                       resourceConsulIntention(),
			"consul_key_prefix":                      resourceConsulKeyPrefix(),
			"consul_keys":                            resourceConsulKeys(),
			"consul_license":                         resourceConsulLicense(),
			"consul_namespace_policy_attachment":     resourceConsulNamespacePolicyAttachment(),
			"consul_namespace_role_attachment":       resourceConsulNamespaceRoleAttachment(),
			"consul_namespace":                       resourceConsulNamespace(),
			"consul_network_area":                    resourceConsulNetworkArea(),
			"consul_node":                            resourceConsulNode(),
			"consul_peering_token":                   resourceSourceConsulPeeringToken(),
			"consul_peering":                         resourceSourceConsulPeering(),
			"consul_prepared_query":                  resourceConsulPreparedQuery(),
			"consul_service":                         resourceConsulService(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var config *Config
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}
	log.Printf("[INFO] Initializing Consul client")
	client, err := config.Client()
	if err != nil {
		return nil, err
	}
	config.client = client

	// Set headers if provided
	headers := d.Get("header").([]interface{})
	parsedHeaders := client.Headers().Clone()

	if parsedHeaders == nil {
		parsedHeaders = make(http.Header)
	}

	for _, h := range headers {
		header := h.(map[string]interface{})
		parsedHeaders.Add(header["name"].(string), header["value"].(string))
	}
	client.SetHeaders(parsedHeaders)

	authJWT := d.Get("auth_jwt").([]interface{})
	if len(authJWT) > 0 {
		authConfig := authJWT[0].(map[string]interface{})
		authMethod := authConfig["auth_method"].(string)
		tfeWorkloadIdentity := authConfig["use_terraform_cloud_workload_identity"].(bool)
		bearerToken := authConfig["bearer_token"].(string)

		if tfeWorkloadIdentity {
			bearerToken = os.Getenv("TFC_WORKLOAD_IDENTITY_TOKEN")
			if bearerToken == "" {
				return nil, fmt.Errorf("auth_jwt.use_terraform_cloud_workload_identity has been set but no token found in TFC_WORKLOAD_IDENTITY_TOKEN environment variable")
			}

		} else if bearerToken == "" {
			return nil, fmt.Errorf("either auth_jwt.bearer_token or auth_jwt.use_terraform_cloud_workload_identity should be set")
		}

		meta := map[string]string{}
		for k, v := range authConfig["meta"].(map[string]interface{}) {
			meta[k] = v.(string)
		}
		_, wOpts := getOptions(d, config)
		token, _, err := client.ACL().Login(&consulapi.ACLLoginParams{
			AuthMethod:  authMethod,
			BearerToken: bearerToken,
			Meta:        meta,
		}, wOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to login using JWT auth method %q: %v", authMethod, err)
		}
		config.Token = token.SecretID
	}

	return config, nil
}

func getClient(d *schema.ResourceData, meta interface{}) (*consulapi.Client, *consulapi.QueryOptions, *consulapi.WriteOptions) {
	config := meta.(*Config)
	client := config.client
	qOpts, wOpts := getOptions(d, config)
	return client, qOpts, wOpts
}

func getOptions(d *schema.ResourceData, meta interface{}) (*consulapi.QueryOptions, *consulapi.WriteOptions) {
	config := meta.(*Config)
	client := config.client
	var dc, token, namespace, partition string
	if v, ok := d.GetOk("datacenter"); ok {
		dc = v.(string)
	}
	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
	}
	if v, ok := d.GetOk("token"); ok {
		token = v.(string)
	}
	if v, ok := d.GetOk("partition"); ok {
		partition = v.(string)
	}

	if dc == "" {
		if config.Datacenter != "" {
			dc = config.Datacenter
		} else {
			info, _ := client.Agent().Self()
			if info != nil {
				dc = info["Config"]["Datacenter"].(string)
			}
		}
	}

	qOpts := &consulapi.QueryOptions{
		Datacenter: dc,
		Namespace:  namespace,
		Partition:  partition,
		Token:      token,
	}
	wOpts := &consulapi.WriteOptions{
		Datacenter: dc,
		Namespace:  namespace,
		Partition:  partition,
		Token:      token,
	}

	return qOpts, wOpts
}

type stateWriter struct {
	d      *schema.ResourceData
	errors []string
}

func newStateWriter(d *schema.ResourceData) *stateWriter {
	return &stateWriter{d: d}
}

func (sw *stateWriter) set(key string, value interface{}) {
	if key == "namespace" || key == "partition" {
		// Consul Enterprise will change "" to "default" but Community Edition only
		// understands the first one.
		if sw.d.Get(key).(string) == "" && value.(string) == "default" {
			value = ""
		}
	}

	err := sw.d.Set(key, value)
	if err != nil {
		sw.errors = append(
			sw.errors,
			fmt.Sprintf(" - failed to set '%s': %v", key, err),
		)
	}
}

func (sw *stateWriter) setJson(key string, value interface{}) {
	marshaled, err := json.Marshal(value)
	if err != nil {
		sw.errors = append(
			sw.errors,
			fmt.Sprintf(" - failed to marshal '%s': %v", key, err),
		)
		return
	}

	sw.set(key, string(marshaled))
}

func (sw *stateWriter) error() error {
	if sw.errors == nil {
		return nil
	}
	errors := strings.Join(sw.errors, "\n")
	return fmt.Errorf("failed to write the state:\n%s", errors)
}
