package consul

import (
	"log"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/mapstructure"
)

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
				}, "http"),
			},

			"http_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONSUL_HTTP_AUTH", ""),
			},

			"ca_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONSUL_CA_FILE", ""),
			},

			"cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONSUL_CERT_FILE", ""),
			},

			"key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONSUL_KEY_FILE", ""),
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
		},

		DataSourcesMap: map[string]*schema.Resource{
			"consul_agent_self":          dataSourceConsulAgentSelf(),
			"consul_agent_config":        dataSourceConsulAgentConfig(),
			"consul_autopilot_health":    dataSourceConsulAutopilotHealth(),
			"consul_nodes":               dataSourceConsulNodes(),
			"consul_service":             dataSourceConsulService(),
			"consul_service_health":      dataSourceConsulServiceHealth(),
			"consul_services":            dataSourceConsulServices(),
			"consul_keys":                dataSourceConsulKeys(),
			"consul_key_prefix":          dataSourceConsulKeyPrefix(),
			"consul_acl_auth_method":     dataSourceConsulACLAuthMethod(),
			"consul_acl_policy":          dataSourceConsulACLPolicy(),
			"consul_acl_role":            dataSourceConsulACLRole(),
			"consul_acl_token":           dataSourceConsulACLToken(),
			"consul_acl_token_secret_id": dataSourceConsulACLTokenSecretID(),

			// Aliases to limit the impact of rename of catalog
			// datasources
			"consul_catalog_nodes":    dataSourceConsulNodes(),
			"consul_catalog_service":  dataSourceConsulService(),
			"consul_catalog_services": dataSourceConsulServices(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"consul_acl_auth_method":             resourceConsulACLAuthMethod(),
			"consul_acl_binding_rule":            resourceConsulACLBindingRule(),
			"consul_acl_policy":                  resourceConsulACLPolicy(),
			"consul_acl_role":                    resourceConsulACLRole(),
			"consul_acl_token":                   resourceConsulACLToken(),
			"consul_acl_token_policy_attachment": resourceConsulACLTokenPolicyAttachment(),
			"consul_agent_service":               resourceConsulAgentService(),
			"consul_catalog_entry":               resourceConsulCatalogEntry(),
			"consul_keys":                        resourceConsulKeys(),
			"consul_key_prefix":                  resourceConsulKeyPrefix(),
			"consul_node":                        resourceConsulNode(),
			"consul_prepared_query":              resourceConsulPreparedQuery(),
			"consul_autopilot_config":            resourceConsulAutopilotConfig(),
			"consul_service":                     resourceConsulService(),
			"consul_intention":                   resourceConsulIntention(),
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
	if _, err := config.Client(); err != nil {
		// The provider must error if the configuration is incorrect. We must
		// check this here.
		return nil, err
	}
	return config, nil
}

func getClient(meta interface{}) *consulapi.Client {
	// We can ignore err since we checked the configuration in providerConfigure()
	client, _ := meta.(*Config).Client()
	return client
}
