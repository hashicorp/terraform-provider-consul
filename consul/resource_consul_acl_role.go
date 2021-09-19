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
			"node_identities": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of node identities that should be applied to the role.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the node.",
						},
						"datacenter": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the node's datacenter.",
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

func resourceConsulACLRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()
	role := getRole(d, meta)

	name := role.Name
	role, _, err := ACL.RoleCreate(role, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to create role '%s': %s", name, err)
	}

	d.SetId(role.ID)
	return resourceConsulACLRoleRead(d, meta)
}

func resourceConsulACLRoleRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	ACL := client.ACL()

	role, _, err := ACL.RoleRead(d.Id(), qOpts)
	if err != nil {
		return fmt.Errorf("Failed to read role '%s': %s", d.Id(), err)
	}
	if role == nil {
		d.SetId("")
		return nil
	}

	policies := make([]string, len(role.Policies))
	for i, policy := range role.Policies {
		policies[i] = policy.ID
	}

	serviceIdentities := make([]map[string]interface{}, len(role.ServiceIdentities))
	for i, serviceIdentity := range role.ServiceIdentities {
		serviceIdentities[i] = map[string]interface{}{
			"service_name": serviceIdentity.ServiceName,
			"datacenters":  serviceIdentity.Datacenters,
		}
	}

	nodeIdentities := make([]interface{}, len(role.NodeIdentities))
	for i, ni := range role.NodeIdentities {
		nodeIdentities[i] = map[string]interface{}{
			"node_name":  ni.NodeName,
			"datacenter": ni.Datacenter,
		}
	}

	sw := newStateWriter(d)

	// Consul Enterprise will change "" to "default" but Community Edition only
	// understands the first one.
	if d.Get("namespace").(string) != "" || role.Namespace != "default" {
		if err = d.Set("namespace", role.Namespace); err != nil {
			return fmt.Errorf("failed to set 'namespace': %v", err)
		}
	}

	sw.set("name", role.Name)
	sw.set("description", role.Description)
	sw.set("policies", policies)
	sw.set("service_identities", serviceIdentities)
	sw.set("node_identities", nodeIdentities)

	return sw.error()
}

func resourceConsulACLRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()
	role := getRole(d, meta)

	role.ID = d.Id()

	role, _, err := ACL.RoleUpdate(role, wOpts)
	if err != nil {
		return fmt.Errorf("Failed to update role '%s': %s", d.Id(), err)
	}

	d.SetId(role.ID)
	return resourceConsulACLRoleRead(d, meta)
}

func resourceConsulACLRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	ACL := client.ACL()

	if _, err := ACL.RoleDelete(d.Id(), wOpts); err != nil {
		return fmt.Errorf("Failed to delete role '%s': %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getRole(d *schema.ResourceData, meta interface{}) *consulapi.ACLRole {
	_, qOpts, _ := getClient(d, meta)
	roleName := d.Get("name").(string)
	role := &consulapi.ACLRole{
		Name:        roleName,
		Description: d.Get("description").(string),
		Namespace:   qOpts.Namespace,
	}
	policies := make([]*consulapi.ACLRolePolicyLink, 0)
	for _, raw := range d.Get("policies").(*schema.Set).List() {
		policies = append(policies, &consulapi.ACLRolePolicyLink{
			ID: raw.(string),
		})
	}
	role.Policies = policies

	for _, raw := range d.Get("service_identities").(*schema.Set).List() {
		s := raw.(map[string]interface{})

		datacenters := make([]string, len(s["datacenters"].(*schema.Set).List()))
		for i, d := range s["datacenters"].(*schema.Set).List() {
			datacenters[i] = d.(string)
		}

		role.ServiceIdentities = append(role.ServiceIdentities, &consulapi.ACLServiceIdentity{
			ServiceName: s["service_name"].(string),
			Datacenters: datacenters,
		})
	}

	for _, ni := range d.Get("node_identities").([]interface{}) {
		n := ni.(map[string]interface{})
		role.NodeIdentities = append(role.NodeIdentities, &consulapi.ACLNodeIdentity{
			NodeName:   n["node_name"].(string),
			Datacenter: n["datacenter"].(string),
		})
	}

	return role
}
