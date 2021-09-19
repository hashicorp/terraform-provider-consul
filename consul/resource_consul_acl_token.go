package consul

import (
	"fmt"
	"log"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceConsulACLToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulACLTokenCreate,
		Read:   resourceConsulACLTokenRead,
		Update: resourceConsulACLTokenUpdate,
		Delete: resourceConsulACLTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"accessor_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Optional:    true,
				Description: "The token id.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The token description.",
			},
			"policies": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of policies.",
			},
			"roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of roles",
			},
			"service_identities": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of service identities that should be applied to the token.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the service.",
						},
						"datacenters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the datacenters the effective policy is valid within.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"node_identities": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of node identities that should be applied to the token.",
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
			"local": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Description: "Flag to set the token local to the current datacenter.",
			},
			"expiration_time": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.ValidateRFC3339TimeString,
				Description:  "If set this represents the point after which a token should be considered revoked and is eligible for destruction.",
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceConsulACLTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	log.Printf("[DEBUG] Creating ACL token")

	aclToken := getToken(d)

	token, _, err := client.ACL().TokenCreate(aclToken, wOpts)
	if err != nil {
		return fmt.Errorf("error creating ACL token: %s", err)
	}

	log.Printf("[DEBUG] Created ACL token %q", token.AccessorID)

	d.SetId(token.AccessorID)

	return resourceConsulACLTokenRead(d, meta)
}

func resourceConsulACLTokenRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	id := d.Id()
	log.Printf("[DEBUG] Reading ACL token %q", id)

	aclToken, _, err := client.ACL().TokenRead(id, qOpts)
	if err != nil {
		if strings.Contains(err.Error(), "ACL not found") {
			log.Printf("[WARN] ACL token not found, removing from state")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Failed to read token '%s': %v", id, err)
	}

	log.Printf("[DEBUG] Read ACL token %q", id)

	roles := make([]string, 0, len(aclToken.Roles))
	for _, roleLink := range aclToken.Roles {
		roles = append(roles, roleLink.Name)
	}

	policies := make([]string, 0, len(aclToken.Policies))
	for _, policyLink := range aclToken.Policies {
		policies = append(policies, policyLink.Name)
	}

	var expirationTime string
	if aclToken.ExpirationTime != nil {
		expirationTime = aclToken.ExpirationTime.Format(time.RFC3339)
	}

	serviceIdentities := make([]interface{}, len(aclToken.ServiceIdentities))
	for i, si := range aclToken.ServiceIdentities {
		serviceIdentities[i] = map[string]interface{}{
			"service_name": si.ServiceName,
			"datacenters":  si.Datacenters,
		}
	}
	nodeIdentities := make([]interface{}, len(aclToken.NodeIdentities))
	for i, ni := range aclToken.NodeIdentities {
		nodeIdentities[i] = map[string]interface{}{
			"node_name":  ni.NodeName,
			"datacenter": ni.Datacenter,
		}
	}

	sw := newStateWriter(d)
	sw.set("accessor_id", aclToken.AccessorID)
	sw.set("description", aclToken.Description)
	sw.set("policies", policies)
	sw.set("roles", roles)
	sw.set("service_identities", serviceIdentities)
	sw.set("node_identities", nodeIdentities)
	sw.set("local", aclToken.Local)
	sw.set("expiration_time", expirationTime)

	return sw.error()
}

func resourceConsulACLTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	id := d.Id()
	log.Printf("[DEBUG] Updating ACL token %q", id)

	aclToken := getToken(d)
	aclToken.AccessorID = id

	_, _, err := client.ACL().TokenUpdate(aclToken, wOpts)
	if err != nil {
		return fmt.Errorf("error updating ACL token %q: %s", id, err)
	}
	log.Printf("[DEBUG] Updated ACL token %q", id)

	return resourceConsulACLTokenRead(d, meta)
}

func resourceConsulACLTokenDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	id := d.Id()

	log.Printf("[DEBUG] Deleting ACL token %q", id)
	_, err := client.ACL().TokenDelete(id, wOpts)
	if err != nil {
		return fmt.Errorf("error deleting ACL token %q: %s", id, err)
	}
	log.Printf("[DEBUG] Deleted ACL token %q", id)

	return nil
}

func getToken(d *schema.ResourceData) *consulapi.ACLToken {
	aclToken := &consulapi.ACLToken{
		AccessorID:  d.Get("accessor_id").(string),
		Description: d.Get("description").(string),
		Local:       d.Get("local").(bool),
	}

	iPolicies := d.Get("policies").(*schema.Set).List()
	policyLinks := make([]*consulapi.ACLTokenPolicyLink, 0, len(iPolicies))
	for _, iPolicy := range iPolicies {
		policyLinks = append(policyLinks, &consulapi.ACLTokenPolicyLink{
			Name: iPolicy.(string),
		})
	}
	aclToken.Policies = policyLinks

	iRoles := d.Get("roles").(*schema.Set).List()
	roleLinks := make([]*consulapi.ACLTokenRoleLink, 0, len(iRoles))
	for _, iRole := range iRoles {
		roleLinks = append(roleLinks, &consulapi.ACLTokenRoleLink{
			Name: iRole.(string),
		})
	}
	aclToken.Roles = roleLinks

	serviceIdentities := []*consulapi.ACLServiceIdentity{}
	for _, si := range d.Get("service_identities").([]interface{}) {
		s := si.(map[string]interface{})

		datacenters := []string{}
		for _, d := range s["datacenters"].([]interface{}) {
			datacenters = append(datacenters, d.(string))
		}

		serviceIdentities = append(serviceIdentities, &consulapi.ACLServiceIdentity{
			ServiceName: s["service_name"].(string),
			Datacenters: datacenters,
		})
	}
	aclToken.ServiceIdentities = serviceIdentities

	nodeIdentities := []*consulapi.ACLNodeIdentity{}
	for _, ni := range d.Get("node_identities").([]interface{}) {
		n := ni.(map[string]interface{})

		nodeIdentities = append(nodeIdentities, &consulapi.ACLNodeIdentity{
			NodeName:   n["node_name"].(string),
			Datacenter: n["datacenter"].(string),
		})
	}
	aclToken.NodeIdentities = nodeIdentities

	expirationTime := d.Get("expiration_time").(string)
	if expirationTime != "" {
		// the string has already been validated so there is no need to check
		// the error here
		t, _ := time.Parse(time.RFC3339, expirationTime)
		aclToken.ExpirationTime = &t
	}

	return aclToken
}
