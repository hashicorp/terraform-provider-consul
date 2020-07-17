package consul

import (
	"encoding/json"
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulACLAuthMethod() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLAuthMethodCreate,
		Read:   resourceConsulACLAuthMethodRead,
		Update: resourceConsulACLAuthMethodUpdate,
		Delete: resourceConsulACLAuthMethodDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the ACL auth method.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the ACL auth method.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional name to use instead of the name attribute when displaying information about this auth method.",
			},

			"max_token_ttl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0s",
				Description: "The maximum life of any token created by this auth method.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					o, err := time.ParseDuration(old)
					if err != nil {
						return false
					}
					n, err := time.ParseDuration(new)
					if err != nil {
						return false
					}
					return o.Seconds() == n.Seconds()
				},
			},

			"token_locality": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The kind of token that this auth method produces. This can be either 'local' or 'global'.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A free form human readable description of the auth method.",
			},

			"config": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The raw configuration for this ACL auth method.",
				Deprecated:  "The config attribute is deprecated, please use config_json instead.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"config_json"},
			},

			"config_json": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The raw configuration for this ACL auth method.",
				ConflictsWith: []string{"config"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					config := d.Get("config").(map[string]interface{})
					return len(config) != 0
				},
			},

			"namespace_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"selector": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"bind_namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceConsulACLAuthMethodCreate(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	wOpts := &consulapi.WriteOptions{}

	authMethod, err := getAuthMethod(d, meta)
	if err != nil {
		return err
	}

	if _, _, err := ACL.AuthMethodCreate(authMethod, wOpts); err != nil {
		return fmt.Errorf("Failed to create auth method '%s': %v", authMethod.Name, err)
	}

	return resourceConsulACLAuthMethodRead(d, meta)
}

func resourceConsulACLAuthMethodRead(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	namespace := getNamespace(d, meta)
	qOpts := &consulapi.QueryOptions{
		Namespace: namespace,
	}

	name := d.Get("name").(string)
	authMethod, _, err := ACL.AuthMethodRead(name, qOpts)
	if err != nil {
		return fmt.Errorf("Failed to read auth method '%s': %v", name, err)
	}
	if authMethod == nil {
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("auth-method-%s", authMethod.Name))

	if err = d.Set("type", authMethod.Type); err != nil {
		return fmt.Errorf("Failed to set 'type': %v", err)
	}

	if err = d.Set("description", authMethod.Description); err != nil {
		return fmt.Errorf("Failed to set 'description': %v", err)
	}

	configJson, err := json.Marshal(authMethod.Config)
	if err != nil {
		return fmt.Errorf("Failed to marshal 'config_json': %v", err)
	}
	if err = d.Set("config_json", string(configJson)); err != nil {
		return fmt.Errorf("Failed to set 'config_json': %v", err)
	}

	if err = d.Set("config", authMethod.Config); err != nil {
		// When a complex configuration is used we can fail to set config as it
		// will not support fields with maps or lists in them. In this case it
		// means that the user used the 'config_json' field, and since we
		// succeeded to set that and 'config' is deprecated, we can just use
		// an empty placeholder value and ignore the error.
		if c := d.Get("config_json").(string); c != "" {
			if err = d.Set("config", map[string]interface{}{}); err != nil {
				return fmt.Errorf("Failed to set 'config': %v", err)
			}
		} else {
			return fmt.Errorf("Failed to set 'config': %v", err)
		}
	}

	if err = d.Set("display_name", authMethod.DisplayName); err != nil {
		return fmt.Errorf("Failed to set 'display_name': %#v", err)
	}

	if err = d.Set("max_token_ttl", authMethod.MaxTokenTTL.String()); err != nil {
		return fmt.Errorf("Failed to set 'max_token_ttl': %#v", err)
	}

	if err = d.Set("token_locality", authMethod.TokenLocality); err != nil {
		return fmt.Errorf("Failed to set 'token_locality': %#v", err)
	}

	rules := make([]interface{}, 0)
	for _, rule := range authMethod.NamespaceRules {
		rules = append(rules, map[string]interface{}{
			"selector":       rule.Selector,
			"bind_namespace": rule.BindNamespace,
		})
	}
	if err = d.Set("namespace_rule", rules); err != nil {
		return fmt.Errorf("Failed to set 'namespace_rule': %v", err)
	}

	return nil
}

func resourceConsulACLAuthMethodUpdate(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	wOpts := &consulapi.WriteOptions{}

	authMethod, err := getAuthMethod(d, meta)
	if err != nil {
		return err
	}

	if _, _, err := ACL.AuthMethodUpdate(authMethod, wOpts); err != nil {
		return fmt.Errorf("Failed to update the auth method '%s': %v", authMethod.Name, err)
	}

	return resourceConsulACLAuthMethodRead(d, meta)
}

func resourceConsulACLAuthMethodDelete(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	namespace := getNamespace(d, meta)
	wOpts := &consulapi.WriteOptions{
		Namespace: namespace,
	}

	authMethodName := d.Get("name").(string)
	if _, err := ACL.AuthMethodDelete(authMethodName, wOpts); err != nil {
		return fmt.Errorf("Failed to delete auth method '%s': %v", authMethodName, err)
	}

	d.SetId("")
	return nil
}

func getAuthMethod(d *schema.ResourceData, meta interface{}) (*consulapi.ACLAuthMethod, error) {
	namespace := getNamespace(d, meta)

	var config map[string]interface{}
	if c := d.Get("config_json").(string); c != "" {
		err := json.Unmarshal([]byte(c), &config)
		if err != nil {
			return nil, fmt.Errorf("Failed to read 'config_json': %v", err)
		}
	} else {
		config = d.Get("config").(map[string]interface{})
	}

	authMethod := &consulapi.ACLAuthMethod{
		Name:          d.Get("name").(string),
		DisplayName:   d.Get("display_name").(string),
		TokenLocality: d.Get("token_locality").(string),
		Type:          d.Get("type").(string),
		Description:   d.Get("description").(string),
		Config:        config,
		Namespace:     namespace,
	}

	if mtt, ok := d.GetOk("max_token_ttl"); ok {
		maxTokenTTL, err := time.ParseDuration(mtt.(string))
		if err != nil {
			return nil, err
		}
		authMethod.MaxTokenTTL = maxTokenTTL
	}

	authMethod.NamespaceRules = make([]*consulapi.ACLAuthMethodNamespaceRule, 0)
	for _, r := range d.Get("namespace_rule").([]interface{}) {
		rule := r.(map[string]interface{})
		namespaceRule := &consulapi.ACLAuthMethodNamespaceRule{
			Selector:      rule["selector"].(string),
			BindNamespace: rule["bind_namespace"].(string),
		}
		authMethod.NamespaceRules = append(authMethod.NamespaceRules, namespaceRule)
	}

	return authMethod, nil
}
