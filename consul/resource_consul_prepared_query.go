package consul

import (
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConsulPreparedQuery() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulPreparedQueryCreate,
		Update: resourceConsulPreparedQueryUpdate,
		Read:   resourceConsulPreparedQueryRead,
		Delete: resourceConsulPreparedQueryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"session": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"stored_token": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"service": {
				Type:     schema.TypeString,
				Required: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"near": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"only_passing": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"connect": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"failover": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nearest_n": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"datacenters": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"dns": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ttl": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"template": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"regexp": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceConsulPreparedQueryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	wo := &consulapi.WriteOptions{
		Datacenter: d.Get("datacenter").(string),
		Token:      d.Get("token").(string),
	}

	pq := preparedQueryDefinitionFromResourceData(d)

	id, _, err := client.PreparedQuery().Create(pq, wo)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConsulPreparedQueryRead(d, meta)
}

func resourceConsulPreparedQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	wo := &consulapi.WriteOptions{
		Datacenter: d.Get("datacenter").(string),
		Token:      d.Get("token").(string),
	}

	pq := preparedQueryDefinitionFromResourceData(d)

	if _, err := client.PreparedQuery().Update(pq, wo); err != nil {
		return err
	}

	return resourceConsulPreparedQueryRead(d, meta)
}

func resourceConsulPreparedQueryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	qo := &consulapi.QueryOptions{
		Datacenter: d.Get("datacenter").(string),
		Token:      d.Get("token").(string),
	}

	queries, _, err := client.PreparedQuery().Get(d.Id(), qo)
	if err != nil {
		// Check for a 404/not found, these are returned as errors.
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	if len(queries) != 1 {
		d.SetId("")
		return nil
	}
	pq := queries[0]

	if err = d.Set("name", pq.Name); err != nil {
		return fmt.Errorf("Failed to set 'name': %v", err)
	}
	if err = d.Set("session", pq.Session); err != nil {
		return fmt.Errorf("Failed to set 'session': %v", err)
	}
	if err = d.Set("stored_token", pq.Token); err != nil {
		return fmt.Errorf("Failed to set 'stored_token': %v", err)
	}
	if err = d.Set("service", pq.Service.Service); err != nil {
		return fmt.Errorf("Failed to set 'service': %v", err)
	}
	if err = d.Set("near", pq.Service.Near); err != nil {
		return fmt.Errorf("Failed to set 'near': %v", err)
	}
	if err = d.Set("only_passing", pq.Service.OnlyPassing); err != nil {
		return fmt.Errorf("Failed to set 'only_passing': %v", err)
	}
	if err = d.Set("connect", pq.Service.Connect); err != nil {
		return fmt.Errorf("Failed to set 'connect': %v", err)
	}
	if err = d.Set("tags", pq.Service.Tags); err != nil {
		return fmt.Errorf("Failed to set 'tags': %v", err)
	}

	failover := make([]map[string]interface{}, 0)
	if pq.Service.Failover.NearestN > 0 {
		failover = append(failover, map[string]interface{}{
			"nearest_n":   pq.Service.Failover.NearestN,
			"datacenters": pq.Service.Failover.Datacenters,
		})
	}
	if err = d.Set("failover", failover); err != nil {
		return fmt.Errorf("Failed to set 'failover': %v", err)
	}

	dns := make([]map[string]interface{}, 0)
	if pq.DNS.TTL != "" {
		dns = append(dns, map[string]interface{}{
			"ttl": pq.DNS.TTL,
		})
	}
	if err = d.Set("dns", dns); err != nil {
		return fmt.Errorf("Failed to set 'dns': %v", err)
	}

	template := make([]map[string]interface{}, 0)
	if pq.Template.Type != "" {
		template = append(template, map[string]interface{}{
			"type":   pq.Template.Type,
			"regexp": pq.Template.Regexp,
		})
	}
	if err = d.Set("template", template); err != nil {
		return fmt.Errorf("Failed to set 'template': %v", err)
	}

	return nil
}

func resourceConsulPreparedQueryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)
	writeOpts := &consulapi.WriteOptions{
		Datacenter: d.Get("datacenter").(string),
		Token:      d.Get("token").(string),
	}

	if _, err := client.PreparedQuery().Delete(d.Id(), writeOpts); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func preparedQueryDefinitionFromResourceData(d *schema.ResourceData) *consulapi.PreparedQueryDefinition {
	pq := &consulapi.PreparedQueryDefinition{
		ID:      d.Id(),
		Name:    d.Get("name").(string),
		Session: d.Get("session").(string),
		Token:   d.Get("stored_token").(string),
		Service: consulapi.ServiceQuery{
			Service:     d.Get("service").(string),
			Near:        d.Get("near").(string),
			OnlyPassing: d.Get("only_passing").(bool),
			Connect:     d.Get("connect").(bool),
		},
	}

	tags := d.Get("tags").(*schema.Set).List()
	pq.Service.Tags = make([]string, len(tags))
	for i, v := range tags {
		pq.Service.Tags[i] = v.(string)
	}

	if _, ok := d.GetOk("failover.0"); ok {
		failover := consulapi.QueryDatacenterOptions{
			NearestN: d.Get("failover.0.nearest_n").(int),
		}

		dcs := d.Get("failover.0.datacenters").([]interface{})
		failover.Datacenters = make([]string, len(dcs))
		for i, v := range dcs {
			failover.Datacenters[i] = v.(string)
		}

		pq.Service.Failover = failover
	}

	if _, ok := d.GetOk("template.0"); ok {
		pq.Template = consulapi.QueryTemplate{
			Type:   d.Get("template.0.type").(string),
			Regexp: d.Get("template.0.regexp").(string),
		}
	}

	if _, ok := d.GetOk("dns.0"); ok {
		pq.DNS = consulapi.QueryDNSOptions{
			TTL: d.Get("dns.0.ttl").(string),
		}
	}

	return pq
}
