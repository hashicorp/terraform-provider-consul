package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	consulACLBindTypeService = "service"
	consulACLBindTypeRole    = "role"
)

func resourceConsulACLBindingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLBindingRuleCreate,
		Read:   resourceConsulACLBindingRuleRead,
		Update: resourceConsulACLBindingRuleUpdate,
		Delete: resourceConsulACLBindingRuleDelete,

		Schema: map[string]*schema.Schema{
			"auth_method": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the ACL auth method this rule apply.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A free form human readable description of the binding rule.",
			},

			"selector": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The expression used to math this rule against valid identities returned from an auth method validation.",
			},

			"bind_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the way the binding rule affects a token created at login.",
				ValidateFunc: validation.StringInSlice(
					[]string{
						consulACLBindTypeService,
						consulACLBindTypeRole,
					},
					false,
				),
			},

			"bind_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name to bind to a token at login-time.",
			},

			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceConsulACLBindingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	wOpts := &consulapi.WriteOptions{}

	rule := getBindingRule(d, meta)

	rule, _, err := ACL.BindingRuleCreate(rule, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to create binding rule: %#v", err)
	}

	d.SetId(rule.ID)

	return resourceConsulACLBindingRuleRead(d, meta)
}

func resourceConsulACLBindingRuleRead(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	namespace := getNamespace(d, meta)
	qOpts := &consulapi.QueryOptions{
		Namespace: namespace,
	}

	rule, _, err := ACL.BindingRuleRead(d.Id(), qOpts)
	if err != nil {
		return fmt.Errorf("Failed to read binding rule '%s': %#v", d.Id(), err)
	}
	if rule == nil {
		d.SetId("")
		return nil
	}

	if err = d.Set("description", rule.Description); err != nil {
		return fmt.Errorf("Failed to set 'description': %#v", err)
	}

	if err = d.Set("selector", rule.Selector); err != nil {
		return fmt.Errorf("Failed to set 'selector': %#v", err)
	}

	if err = d.Set("bind_type", rule.BindType); err != nil {
		return fmt.Errorf("Failed to set 'bind_type': %#v", err)
	}

	if err = d.Set("bind_name", rule.BindName); err != nil {
		return fmt.Errorf("Failed to set 'bind_name': %#v", err)
	}

	return nil
}

func resourceConsulACLBindingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	wOpts := &consulapi.WriteOptions{}

	rule := getBindingRule(d, meta)

	rule, _, err := ACL.BindingRuleUpdate(rule, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to update binding rule '%s': %#v", d.Id(), err)
	}

	return resourceConsulACLBindingRuleRead(d, meta)
}

func resourceConsulACLBindingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	namespace := getNamespace(d, meta)
	wOpts := &consulapi.WriteOptions{
		Namespace: namespace,
	}

	if _, err := ACL.BindingRuleDelete(d.Id(), wOpts); err != nil {
		return fmt.Errorf("Failed to delete binding rule '%s': %#v", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func getBindingRule(d *schema.ResourceData, meta interface{}) *consulapi.ACLBindingRule {
	namespace := getNamespace(d, meta)
	rule := &consulapi.ACLBindingRule{
		ID:          d.Id(),
		Description: d.Get("description").(string),
		AuthMethod:  d.Get("auth_method").(string),
		Selector:    d.Get("selector").(string),
		BindName:    d.Get("bind_name").(string),
		Namespace:   namespace,
	}

	bindType := d.Get("bind_type").(string)
	if bindType == consulACLBindTypeService {
		rule.BindType = consulapi.BindingRuleBindTypeService
	} else {
		// bindType == consulACLBindTypeRole
		rule.BindType = consulapi.BindingRuleBindTypeRole
	}

	return rule
}
