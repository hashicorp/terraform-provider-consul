package consulyaml

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
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
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CONSUL_TOKEN",
					"CONSUL_HTTP_TOKEN",
				}, ""),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"consul_yaml": dataSourceConsulKeyPrefixFromFile(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"consul_yaml": resourceConsulKeyPrefixFromFile(),
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

func yamlhash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

// getDC is used to get the datacenter of the local agent
func getDC(d *schema.ResourceData, client *consulapi.Client) (string, error) {
	if v, ok := d.GetOk("datacenter"); ok {
		return v.(string), nil
	}
	info, err := client.Agent().Self()
	if err != nil {
		return "", fmt.Errorf("Failed to get datacenter from Consul agent: %v", err)
	}
	return info["Config"]["Datacenter"].(string), nil
}

func resourceConsulKeysDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	kv := client.KV()
	token := d.Get("token").(string)
	dc, err := getDC(d, client)
	if err != nil {
		return err
	}

	keyClient := newKeyClient(kv, dc, token)

	// Clean up any keys that we're explicitly managing
	keys := d.Get("key").(*schema.Set).List()
	for _, raw := range keys {
		_, path, sub, err := parseKey(raw)
		if err != nil {
			return err
		}

		// Skip if the key is non-managed
		shouldDelete, ok := sub["delete"].(bool)
		if !ok || !shouldDelete {
			continue
		}

		if err := keyClient.Delete(path); err != nil {
			return err
		}
	}

	// Clear the ID
	d.SetId("")
	return nil
}

// parseKey is used to parse a key into a name, path, config or error
func parseKey(raw interface{}) (string, string, map[string]interface{}, error) {
	sub, ok := raw.(map[string]interface{})
	if !ok {
		return "", "", nil, fmt.Errorf("Failed to unroll: %#v", raw)
	}

	key := sub["name"].(string)

	path, ok := sub["path"].(string)
	if !ok {
		return "", "", nil, fmt.Errorf("Failed to get path for key '%s'", key)
	}
	return key, path, sub, nil
}

// attributeValue determines the value for a key, potentially
// using a default value if provided.
func attributeValue(sub map[string]interface{}, readValue string) string {
	// Use the value if given
	if readValue != "" {
		return readValue
	}

	// Use a default if given
	if raw, ok := sub["default"]; ok {
		switch def := raw.(type) {
		case string:
			return def
		case bool:
			return strconv.FormatBool(def)
		}
	}

	// No value
	return ""
}
