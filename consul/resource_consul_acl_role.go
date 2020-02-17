package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulACLRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLRoleCreate,
		Read:   resourceConsulACLRoleRead,
		Update: resourceConsulACLRoleUpdate,
		Delete: resourceConsulACLRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ACL role.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A free form human readable description of the role.",
			},

			"policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of policies that should be applied to the role.",
			},

			"service_identities": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the service.",
						},

						"datacenters": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The datacenters the effective policy is valid within. When no datacenters are provided the effective policy is valid in all datacenters including those which do not yet exist but may in the future.",
						},
					},
				},
				Description: "The list of service identities that should be applied to the role.",
			},
		},
	}
}

func resourceConsulACLRoleCreate(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	role := getRole(d)
	wOpts := &consulapi.WriteOptions{}

	name := role.Name
	role, _, err := ACL.RoleCreate(role, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to create role '%s': %s", name, err)
	}

	d.SetId(role.ID)
	return resourceConsulACLRoleRead(d, meta)
}

func resourceConsulACLRoleRead(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	qOpts := &consulapi.QueryOptions{}

	role, _, err := ACL.RoleRead(d.Id(), qOpts)
	if err != nil {
		return fmt.Errorf("Failed to read role '%s': %s", d.Id(), err)
	}
	if role == nil {
		d.SetId("")
		return nil
	}

	if err = d.Set("name", role.Name); err != nil {
		return fmt.Errorf("Failed to set 'name': %s", err)
	}

	if err = d.Set("description", role.Description); err != nil {
		return fmt.Errorf("Failed to set 'description': %s", err)
	}

	policies := make([]string, len(role.Policies))
	for i, policy := range role.Policies {
		policies[i] = policy.ID
	}
	if err = d.Set("policies", policies); err != nil {
		return fmt.Errorf("Failed to set 'policies': %s", err)
	}

	serviceIdentities := make([]map[string]interface{}, len(role.ServiceIdentities))
	for i, serviceIdentity := range role.ServiceIdentities {
		serviceIdentities[i] = map[string]interface{}{
			"service_name": serviceIdentity.ServiceName,
			"datacenters":  serviceIdentity.Datacenters,
		}
	}
	if err = d.Set("service_identities", serviceIdentities); err != nil {
		return fmt.Errorf("Failed to set 'service_identities': %s", err)
	}

	return nil
}

func resourceConsulACLRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	role := getRole(d)
	wOpts := &consulapi.WriteOptions{}

	role.ID = d.Id()

	role, _, err := ACL.RoleUpdate(role, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to update role '%s': %s", d.Id(), err)
	}

	d.SetId(role.ID)
	return resourceConsulACLRoleRead(d, meta)
}

func resourceConsulACLRoleDelete(d *schema.ResourceData, meta interface{}) error {
	ACL := getClient(meta).ACL()
	wOpts := &consulapi.WriteOptions{}

	if _, err := ACL.RoleDelete(d.Id(), wOpts); err != nil {
		return fmt.Errorf("Failed to delete role '%s': %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getRole(d *schema.ResourceData) *consulapi.ACLRole {
	roleName := d.Get("name").(string)
	role := &consulapi.ACLRole{
		Name:        roleName,
		Description: d.Get("description").(string),
	}
	policies := make([]*consulapi.ACLRolePolicyLink, 0)
	for _, raw := range d.Get("policies").(*schema.Set).List() {
		policies = append(policies, &consulapi.ACLRolePolicyLink{
			ID: raw.(string),
		})
	}
	role.Policies = policies

	serviceIdentities := make([]*consulapi.ACLServiceIdentity, 0)
	for _, raw := range d.Get("service_identities").(*schema.Set).List() {
		s := raw.(map[string]interface{})

		datacenters := make([]string, len(s["datacenters"].(*schema.Set).List()))
		for i, d := range s["datacenters"].(*schema.Set).List() {
			datacenters[i] = d.(string)
		}

		serviceIdentities = append(serviceIdentities, &consulapi.ACLServiceIdentity{
			ServiceName: s["service_name"].(string),
			Datacenters: datacenters,
		})
	}
	role.ServiceIdentities = serviceIdentities

	return role
}
