package consul

import (
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
				Required:    true,
				Description: "The raw configuration for this ACL auth method.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
		return fmt.Errorf("Failed to create auth method '%s': %#v", authMethod.Name, err)
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
		return fmt.Errorf("Failed to read auth method '%s': %#v", name, err)
	}
	if authMethod == nil {
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("auth-method-%s", authMethod.Name))

	if err = d.Set("type", authMethod.Type); err != nil {
		return fmt.Errorf("Failed to set 'type': %#v", err)
	}

	if err = d.Set("description", authMethod.Description); err != nil {
		return fmt.Errorf("Failed to set 'description': %#v", err)
	}

	if err = d.Set("config", authMethod.Config); err != nil {
		return fmt.Errorf("Failed to set 'config': %#v", err)
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
		return fmt.Errorf("Failed to update the auth method '%s': %#v", authMethod.Name, err)
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
		return fmt.Errorf("Failed to delete auth method '%s': %#v", authMethodName, err)
	}

	d.SetId("")
	return nil
}

func getAuthMethod(d *schema.ResourceData, meta interface{}) (*consulapi.ACLAuthMethod, error) {
	namespace := getNamespace(d, meta)

	authMethod := &consulapi.ACLAuthMethod{
		Name:          d.Get("name").(string),
		DisplayName:   d.Get("display_name").(string),
		TokenLocality: d.Get("token_locality").(string),
		Type:          d.Get("type").(string),
		Description:   d.Get("description").(string),
		Config:        d.Get("config").(map[string]interface{}),
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
