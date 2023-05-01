// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

		Description: `Allows Terraform to manage a Consul prepared query.

Managing prepared queries is done using Consul's REST API. This resource is useful to provide a consistent and declarative way of managing prepared queries in your Consul cluster using Terraform.`,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The datacenter to use. This overrides the agent's default datacenter and the datacenter in the provider setup.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the prepared query. Used to identify the prepared query during requests. Can be specified as an empty string to configure the query as a catch-all.",
			},

			"session": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Consul session to tie this query's lifetime to.  This is an advanced parameter that should not be used without a complete understanding of Consul sessions and the implications of their use (it is recommended to leave this blank in nearly all cases).  If this parameter is omitted the query will not expire.",
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Deprecated:  tokenDeprecationMessage,
				Description: "The ACL token to use when saving the prepared query. This overrides the token that the agent provides by default.",
			},

			"stored_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ACL token to store with the prepared query. This token will be used by default whenever the query is executed.",
			},

			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the service to query",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of required and/or disallowed tags.  If a tag is in this list it must be present.  If the tag is preceded with a "!" then it is disallowed.`,
			},

			"near": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allows specifying the name of a node to sort results near using Consul's distance sorting and network coordinates. The magic `_agent` value can be used to always sort nearest the node servicing the request.",
			},

			"only_passing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When `true`, the prepared query will only return nodes with passing health checks in the result.",
			},

			"connect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When `true` the prepared query will return connect proxy services for a queried service.  Conditions such as `tags` in the prepared query will be matched against the proxy service. Defaults to false.",
			},

			"ignore_check_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies a list of check IDs that should be ignored when filtering unhealthy instances. This is mostly useful in an emergency or as a temporary measure when a health check is found to be unreliable. Being able to ignore it in centrally-defined queries can be simpler than de-registering the check as an interim solution until the check can be fixed.",
			},

			"node_meta": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies a list of user-defined key/value pairs that will be used for filtering the query results to nodes with the given metadata values present.",
			},

			"service_meta": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies a list of user-defined key/value pairs that will be used for filtering the query results to services with the given metadata values present.",
			},

			"failover": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Options for controlling behavior when no healthy nodes are available in the local DC.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nearest_n": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Return results from this many datacenters, sorted in ascending order of estimated RTT.",
						},
						"datacenters": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Remote datacenters to return results from.",
						},
						"targets": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a sequential list of remote datacenters and cluster peers to failover to if there are no healthy service instances in the local datacenter. This option cannot be used with `nearest_n` or `datacenters`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"peer": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies a cluster peer to use for failover.",
									},
									"datacenter": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies a WAN federated datacenter to forward the query to.",
									},
								},
							},
						},
					},
				},
			},

			"dns": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Settings for controlling the DNS response details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ttl": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The TTL to send when returning DNS results.",
						},
					},
				},
			},

			"template": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Query templating options. This is used to make a single prepared query respond to many different requests",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of template matching to perform. Currently only `name_prefix_match` is supported.",
						},
						"regexp": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The regular expression to match with. When using `name_prefix_match`, this regex is applied against the query name.",
						},
					},
				},
			},
		},
	}
}

func resourceConsulPreparedQueryCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	pq := preparedQueryDefinitionFromResourceData(d)

	id, _, err := client.PreparedQuery().Create(pq, wOpts)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConsulPreparedQueryRead(d, meta)
}

func resourceConsulPreparedQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	pq := preparedQueryDefinitionFromResourceData(d)

	if _, err := client.PreparedQuery().Update(pq, wOpts); err != nil {
		return err
	}

	return resourceConsulPreparedQueryRead(d, meta)
}

func resourceConsulPreparedQueryRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)

	queries, _, err := client.PreparedQuery().Get(d.Id(), qOpts)
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

	sw := newStateWriter(d)
	sw.set("name", pq.Name)
	sw.set("session", pq.Session)
	sw.set("stored_token", pq.Token)
	sw.set("service", pq.Service.Service)
	sw.set("near", pq.Service.Near)
	sw.set("only_passing", pq.Service.OnlyPassing)
	sw.set("connect", pq.Service.Connect)
	sw.set("tags", pq.Service.Tags)
	sw.set("ignore_check_ids", pq.Service.IgnoreCheckIDs)
	sw.set("node_meta", pq.Service.NodeMeta)
	sw.set("service_meta", pq.Service.ServiceMeta)

	// Since failover and dns are implemented with an optionnal list instead of a
	// sub-resource, writing those attributes to the state is more involved that
	// it needs to.

	failover := make([]map[string]interface{}, 0)

	// First we must find whether the user wrote a failover block
	userWroteFailover := len(d.Get("failover").([]interface{})) != 0

	// We must write a failover block if the user wrote one or if one of the values
	// differ from the defaults
	if userWroteFailover || pq.Service.Failover.NearestN > 0 || len(pq.Service.Failover.Datacenters) > 0 || len(pq.Service.Failover.Targets) > 0 {
		targets := []interface{}{}
		for _, target := range pq.Service.Failover.Targets {
			targets = append(targets, map[string]interface{}{
				"peer":       target.Peer,
				"datacenter": target.Datacenter,
			})
		}
		failover = append(failover, map[string]interface{}{
			"nearest_n":   pq.Service.Failover.NearestN,
			"datacenters": pq.Service.Failover.Datacenters,
			"targets":     targets,
		})
	}

	// We can finally set the failover attribute
	sw.set("failover", failover)

	dns := make([]map[string]interface{}, 0)

	userWroteDNS := len(d.Get("dns").([]interface{})) != 0

	if userWroteDNS || pq.DNS.TTL != "" {
		dns = append(dns, map[string]interface{}{
			"ttl": pq.DNS.TTL,
		})
	}
	sw.set("dns", dns)

	template := make([]map[string]interface{}, 0)

	userWroteTemplate := len(d.Get("template").([]interface{})) != 0

	if userWroteTemplate || pq.Template.Type != "" {
		template = append(template, map[string]interface{}{
			"type":   pq.Template.Type,
			"regexp": pq.Template.Regexp,
		})
	}
	sw.set("template", template)

	return sw.error()
}

func resourceConsulPreparedQueryDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	if _, err := client.PreparedQuery().Delete(d.Id(), wOpts); err != nil {
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

	pq.Service.NodeMeta = make(map[string]string)
	for k, v := range d.Get("node_meta").(map[string]interface{}) {
		pq.Service.NodeMeta[k] = v.(string)
	}

	pq.Service.ServiceMeta = make(map[string]string)
	for k, v := range d.Get("service_meta").(map[string]interface{}) {
		pq.Service.ServiceMeta[k] = v.(string)
	}

	ignoreCheckIDs := d.Get("ignore_check_ids").([]interface{})
	pq.Service.IgnoreCheckIDs = make([]string, len(ignoreCheckIDs))
	for i, id := range ignoreCheckIDs {
		pq.Service.IgnoreCheckIDs[i] = id.(string)
	}

	if _, ok := d.GetOk("failover.0"); ok {
		failover := consulapi.QueryFailoverOptions{
			NearestN: d.Get("failover.0.nearest_n").(int),
		}

		dcs := d.Get("failover.0.datacenters").([]interface{})
		failover.Datacenters = make([]string, len(dcs))
		for i, v := range dcs {
			failover.Datacenters[i] = v.(string)
		}

		targets := d.Get("failover.0.targets").([]interface{})
		failover.Targets = make([]consulapi.QueryFailoverTarget, len(targets))
		for i, v := range targets {
			target := v.(map[string]interface{})
			failover.Targets[i] = consulapi.QueryFailoverTarget{
				Peer:       target["peer"].(string),
				Datacenter: target["datacenter"].(string),
			}
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
