package consul

import (
	"fmt"

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

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A free form human readable description of the auth method.",
			},

			"config": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "The raw configuration for this ACL auth method.",
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
	namespace := getNamespace(d, meta)

	authMethod := &consulapi.ACLAuthMethod{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Config:      d.Get("config").(map[string]interface{}),
		Namespace:   namespace,
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

	return nil
}

func resourceConsulACLAuthMethodUpdate(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	wOpts := &consulapi.WriteOptions{}
	namespace := getNamespace(d, meta)

	authMethod := &consulapi.ACLAuthMethod{
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Config:      d.Get("config").(map[string]interface{}),
		Namespace:   namespace,
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
