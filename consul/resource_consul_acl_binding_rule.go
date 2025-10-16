// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulACLBindingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLBindingRuleCreate,
		Read:   resourceConsulACLBindingRuleRead,
		Update: resourceConsulACLBindingRuleUpdate,
		Delete: resourceConsulACLBindingRuleDelete,

		Description: "Starting with Consul 1.5.0, the consul_acl_binding_rule resource can be used to managed Consul ACL binding rules.",

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
				Description: "The expression used to match this rule against valid identities returned from an auth method validation.",
			},

			"bind_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the way the binding rule affects a token created at login.",
			},

			"bind_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name to bind to a token at login-time.",
			},

			"bind_vars": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The variables used when binding rule type is `templated-policy`. Can be lightly templated using HIL `${foo}` syntax from available field names.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of node, workload identity or service.",
						},
					},
				},
			},

			"namespace": {
				Type:        schema.TypeString,
				Description: "The namespace to create the binding rule within.",
				Optional:    true,
				ForceNew:    true,
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The partition the ACL binding rule is associated with.",
			},
		},
	}
}

func resourceConsulACLBindingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	ACL := client.ACL()

	rule := getBindingRule(d, meta)

	rule, _, err := ACL.BindingRuleCreate(rule, wOpts)
	if err != nil {
		return fmt.Errorf("failed to create binding rule: %v", err)
	}

	if err := waitForACLTokenReplication(client.ACL(), qOpts, rule.CreateIndex); err != nil {
		return err
	}

	d.SetId(rule.ID)

	return resourceConsulACLBindingRuleRead(d, meta)
}

func resourceConsulACLBindingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	ACL := client.ACL()

	rule, _, err := ACL.BindingRuleRead(d.Id(), qOpts)
	if err != nil {
		return fmt.Errorf("failed to read binding rule '%s': %v", d.Id(), err)
	}
	if rule == nil {
		d.SetId("")
		return nil
	}

	sw := newStateWriter(d)
	sw.set("description", rule.Description)
	sw.set("selector", rule.Selector)
	sw.set("bind_type", rule.BindType)
	sw.set("bind_name", rule.BindName)
	sw.set("namespace", rule.Namespace)
	sw.set("partition", rule.Partition)

	if rule.BindVars != nil {
		bindVars := []map[string]interface{}{
			{"name": rule.BindVars.Name},
		}
		sw.set("bind_vars", bindVars)
	}

	return sw.error()
}

func resourceConsulACLBindingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, wOpts := getClient(d, meta)
	ACL := client.ACL()

	rule := getBindingRule(d, meta)

	u, _, err := ACL.BindingRuleUpdate(rule, wOpts)
	if err != nil {
		return fmt.Errorf("failed to update binding rule '%s': %v", d.Id(), err)
	}

	if err := waitForACLTokenReplication(client.ACL(), qOpts, u.ModifyIndex); err != nil {
		return err
	}

	return resourceConsulACLBindingRuleRead(d, meta)
}

func resourceConsulACLBindingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()

	if _, err := ACL.BindingRuleDelete(d.Id(), wOpts); err != nil {
		return fmt.Errorf("failed to delete binding rule '%s': %v", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func getBindingRule(d *schema.ResourceData, meta interface{}) *consulapi.ACLBindingRule {
	_, _, wOpts := getClient(d, meta)
	bindingRule := &consulapi.ACLBindingRule{
		ID:          d.Id(),
		Description: d.Get("description").(string),
		AuthMethod:  d.Get("auth_method").(string),
		Selector:    d.Get("selector").(string),
		BindName:    d.Get("bind_name").(string),
		BindType:    consulapi.BindingRuleBindType(d.Get("bind_type").(string)),
		Namespace:   wOpts.Namespace,
	}

	if bindVars, ok := d.GetOk("bind_vars.0"); ok {
		tv := bindVars.(map[string]interface{})

		processedVars := &consulapi.ACLTemplatedPolicyVariables{}
		if tv["name"] != nil {
			processedVars.Name = tv["name"].(string)
		}

		bindingRule.BindVars = processedVars
	}

	return bindingRule
}
