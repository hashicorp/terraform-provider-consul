package consul

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/mapstructure"
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
				Type:     schema.TypeString,
				Optional: true,
			},

			"address": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CONSUL_ADDRESS",
					"CONSUL_HTTP_ADDR",
				}, "localhost:8500"),
			},

			"scheme": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CONSUL_SCHEME",
					"CONSUL_HTTP_SCHEME",
				}, ""),
			},

			"http_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONSUL_HTTP_AUTH", ""),
			},

			"ca_file": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("CONSUL_CA_FILE", nil),
				ConflictsWith: []string{"ca_pem"},
			},

			"ca_pem": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ca_file"},
			},

			"cert_file": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("CONSUL_CERT_FILE", nil),
				ConflictsWith: []string{"cert_pem"},
			},

			"cert_pem": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cert_file"},
			},

			"key_file": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("CONSUL_KEY_FILE", nil),
				ConflictsWith: []string{"key_pem"},
			},

			"key_pem": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key_file"},
			},

			"ca_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONSUL_CAPATH", ""),
			},

			"insecure_https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CONSUL_TOKEN",
					"CONSUL_HTTP_TOKEN",
				}, ""),
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"header": {
				Type:        schema.TypeList,
				Optional:    true,
				Sensitive:   true,
				Description: "Additional headers to send with each Consul request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The header name",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The header value",
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
			"consul_peering":              dataSourceConsulPeering(),
			"consul_peerings":             dataSourceConsulPeerings(),

			// Aliases to limit the impact of rename of catalog
			// datasources
			"consul_catalog_nodes":    deprecated("consul_catalog_nodes", dataSourceConsulNodes()),
			"consul_catalog_service":  deprecated("consul_catalog_service", dataSourceConsulService()),
			"consul_catalog_services": deprecated("consul_catalog_services", dataSourceConsulServices()),
		},

		ResourcesMap: map[string]*schema.Resource{
			"consul_acl_auth_method":             resourceConsulACLAuthMethod(),
			"consul_acl_binding_rule":            resourceConsulACLBindingRule(),
			"consul_acl_policy":                  resourceConsulACLPolicy(),
			"consul_acl_role":                    resourceConsulACLRole(),
			"consul_acl_token":                   resourceConsulACLToken(),
			"consul_acl_token_policy_attachment": resourceConsulACLTokenPolicyAttachment(),
			"consul_acl_token_role_attachment":   resourceConsulACLTokenRoleAttachment(),
			"consul_admin_partition":             resourceConsulAdminPartition(),
			"consul_agent_service":               resourceConsulAgentService(),
			"consul_catalog_entry":               resourceConsulCatalogEntry(),
			"consul_certificate_authority":       resourceConsulCertificateAuthority(),
			"consul_config_entry":                resourceConsulConfigEntry(),
			"consul_keys":                        resourceConsulKeys(),
			"consul_key_prefix":                  resourceConsulKeyPrefix(),
			"consul_license":                     resourceConsulLicense(),
			"consul_namespace":                   resourceConsulNamespace(),
			"consul_namespace_policy_attachment": resourceConsulNamespacePolicyAttachment(),
			"consul_namespace_role_attachment":   resourceConsulNamespaceRoleAttachment(),
			"consul_node":                        resourceConsulNode(),
			"consul_prepared_query":              resourceConsulPreparedQuery(),
			"consul_autopilot_config":            resourceConsulAutopilotConfig(),
			"consul_service":                     resourceConsulService(),
			"consul_intention":                   resourceConsulIntention(),
			"consul_network_area":                resourceConsulNetworkArea(),
			"consul_peering_token":               resourceSourceConsulPeeringToken(),
			"consul_peering":                     resourceSourceConsulPeering(),
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
	return config, nil
}

func getClient(d *schema.ResourceData, meta interface{}) (*consulapi.Client, *consulapi.QueryOptions, *consulapi.WriteOptions) {
	client := meta.(*Config).client
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
		if meta.(*Config).Datacenter != "" {
			dc = meta.(*Config).Datacenter
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
	return client, qOpts, wOpts
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
			fmt.Sprintf("failed to marshal '%s': %v", key, err),
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
