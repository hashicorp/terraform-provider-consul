package consul

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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
			"consul_agent_self":       dataSourceConsulAgentSelf(),
			"consul_agent_config":     dataSourceConsulAgentConfig(),
			"consul_autopilot_health": dataSourceConsulAutopilotHealth(),
			"consul_nodes":            dataSourceConsulNodes(),
			"consul_service":          dataSourceConsulService(),
			"consul_services":         dataSourceConsulServices(),
			"consul_keys":             dataSourceConsulKeys(),
			"consul_key_prefix":       dataSourceConsulKeyPrefix(),

			// Aliases to limit the impact of rename of catalog
			// datasources
			"consul_catalog_nodes":    dataSourceConsulNodes(),
			"consul_catalog_service":  dataSourceConsulService(),
			"consul_catalog_services": dataSourceConsulServices(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"consul_agent_service":  resourceConsulAgentService(),
			"consul_catalog_entry":  resourceConsulCatalogEntry(),
			"consul_keys":           resourceConsulKeys(),
			"consul_key_prefix":     resourceConsulKeyPrefix(),
			"consul_node":           resourceConsulNode(),
			"consul_prepared_query": resourceConsulPreparedQuery(),
			"consul_service":        resourceConsulService(),
			"consul_intention":      resourceConsulIntention(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var config Config
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}
	log.Printf("[INFO] Initializing Consul client")
	return config.Client()
}
