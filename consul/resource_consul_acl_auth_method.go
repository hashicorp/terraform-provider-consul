// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"encoding/json"
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConsulACLAuthMethod() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLAuthMethodCreate,
		Read:   resourceConsulACLAuthMethodRead,
		Update: resourceConsulACLAuthMethodUpdate,
		Delete: resourceConsulACLAuthMethodDelete,

		Description: "Starting with Consul 1.5.0, the `consul_acl_auth_method` resource can be used to managed [Consul ACL auth methods](https://www.consul.io/docs/acl/auth-methods).",

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
				Description: "The maximum life of any token created by this auth method. **This attribute is required and must be set to a nonzero for the OIDC auth method.**",
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
				Deprecated:  "The config attribute is deprecated, please use `config_json` instead.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"config_json"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "" || new == "0"
				},
			},

			"config_json": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The raw configuration for this ACL auth method.",
				ConflictsWith: []string{"config"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "" || new == "0"
				},
			},

			"namespace_rule": {
				Type:        schema.TypeList,
				Description: "A set of rules that control which namespace tokens created via this auth method will be created within.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"selector": {
							Type:        schema.TypeString,
							Description: "Specifies the expression used to match this namespace rule against valid identities returned from an auth method validation.",
							Optional:    true,
						},
						"bind_namespace": {
							Type:        schema.TypeString,
							Description: "If the namespace rule's `selector` matches then this is used to control the namespace where the token is created.",
							Required:    true,
						},
					},
				},
			},

			"namespace": {
				Type:        schema.TypeString,
				Description: "The namespace in which to create the auth method.",
				Optional:    true,
				ForceNew:    true,
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The partition the ACL auth method is associated with.",
			},
		},
	}
}

func resourceConsulACLAuthMethodCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()

	authMethod, err := getAuthMethod(d, meta)
	if err != nil {
		return err
	}

	if _, _, err := ACL.AuthMethodCreate(authMethod, wOpts); err != nil {
		return fmt.Errorf("failed to create auth method '%s': %v", authMethod.Name, err)
	}

	return resourceConsulACLAuthMethodRead(d, meta)
}

func resourceConsulACLAuthMethodRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	ACL := client.ACL()

	name := d.Get("name").(string)
	authMethod, _, err := ACL.AuthMethodRead(name, qOpts)
	if err != nil {
		return fmt.Errorf("failed to read auth method '%s': %v", name, err)
	}
	if authMethod == nil {
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("auth-method-%s", authMethod.Name))

	sw := newStateWriter(d)
	sw.set("type", authMethod.Type)
	sw.set("description", authMethod.Description)
	sw.setJson("config_json", authMethod.Config)

	if err = d.Set("config", authMethod.Config); err != nil {
		// When a complex configuration is used we can fail to set config as it
		// will not support fields with maps or lists in them. In this case it
		// means that the user used the 'config_json' field, and since we
		// succeeded to set that and 'config' is deprecated, we can just use
		// an empty placeholder value and ignore the error.
		if c := d.Get("config_json").(string); c != "" {
			sw.set("config", map[string]interface{}{})
		} else {
			return fmt.Errorf("failed to set 'config': %v", err)
		}
	}

	sw.set("display_name", authMethod.DisplayName)
	sw.set("max_token_ttl", authMethod.MaxTokenTTL.String())
	sw.set("token_locality", authMethod.TokenLocality)

	rules := make([]interface{}, 0)
	for _, rule := range authMethod.NamespaceRules {
		rules = append(rules, map[string]interface{}{
			"selector":       rule.Selector,
			"bind_namespace": rule.BindNamespace,
		})
	}
	sw.set("namespace_rule", rules)
	sw.set("namespace", authMethod.Namespace)
	sw.set("partition", authMethod.Partition)

	return sw.error()
}

func resourceConsulACLAuthMethodUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()

	authMethod, err := getAuthMethod(d, meta)
	if err != nil {
		return err
	}

	if _, _, err := ACL.AuthMethodUpdate(authMethod, wOpts); err != nil {
		return fmt.Errorf("failed to update the auth method '%s': %v", authMethod.Name, err)
	}

	return resourceConsulACLAuthMethodRead(d, meta)
}

func resourceConsulACLAuthMethodDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()

	authMethodName := d.Get("name").(string)
	if _, err := ACL.AuthMethodDelete(authMethodName, wOpts); err != nil {
		return fmt.Errorf("failed to delete auth method '%s': %v", authMethodName, err)
	}

	d.SetId("")
	return nil
}

func getAuthMethod(d *schema.ResourceData, meta interface{}) (*consulapi.ACLAuthMethod, error) {
	_, qOpts, _ := getClient(d, meta)

	var config map[string]interface{}
	if c := d.Get("config_json").(string); c != "" {
		err := json.Unmarshal([]byte(c), &config)
		if err != nil {
			return nil, fmt.Errorf("failed to read 'config_json': %v", err)
		}
	} else {
		config = d.Get("config").(map[string]interface{})
	}

	if len(config) == 0 {
		return nil, fmt.Errorf("one of 'config' or 'config_json' must be set")
	}

	authMethod := &consulapi.ACLAuthMethod{
		Name:          d.Get("name").(string),
		DisplayName:   d.Get("display_name").(string),
		TokenLocality: d.Get("token_locality").(string),
		Type:          d.Get("type").(string),
		Description:   d.Get("description").(string),
		Config:        config,
		Namespace:     qOpts.Namespace,
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
