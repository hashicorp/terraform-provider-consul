// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type samenessGroup struct{}

func (s *samenessGroup) GetKind() string {
	return consulapi.SamenessGroup
}

func (s *samenessGroup) GetDescription() string {
	return "The `consul_config_entry_sameness_group` resource configures a [sameness group](https://developer.hashicorp.com/consul/docs/connect/config-entries/sameness-group). Sameness groups associate services with identical names across partitions and cluster peers."
}

func (s *samenessGroup) GetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Specifies a name for the configuration entry.",
			Required:    true,
			ForceNew:    true,
		},
		"partition": {
			Type:        schema.TypeString,
			Description: "Specifies the local admin partition that the sameness group applies to.",
			Optional:    true,
			ForceNew:    true,
		},
		"default_for_failover": {
			Type:        schema.TypeBool,
			Description: "Determines whether the sameness group should be used to establish connections to services with the same name during failover scenarios. When this field is set to `true`, DNS queries and upstream requests automatically failover to services in the sameness group according to the order of the members in the `members` list.\n\nWhen this field is set to `false`, you can still use a sameness group for `failover` by configuring the failover block of a service resolver configuration entry.",
			Optional:    true,
		},
		"include_local": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"members": {
			Type:        schema.TypeList,
			Description: "Specifies the partitions and cluster peers that are members of the sameness group from the perspective of the local partition.\n\nThe local partition should be the first member listed. The order of the members determines their precedence during failover scenarios. If a member is listed but Consul cannot connect to it, failover proceeds with the next healthy member in the list.",
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"partition": {
						Type:        schema.TypeString,
						Description: "Specifies a partition in the local datacenter that is a member of the sameness group. When the value of this field is set to `*`, all local partitions become members of the sameness group.",
						Optional:    true,
					},
					"peer": {
						Type:        schema.TypeString,
						Description: "Specifies the name of a cluster peer that is a member of the sameness group.\n\nCluster peering connections must be established before adding a peer to the list of members.",
						Optional:    true,
					},
				},
			},
		},
		"meta": {
			Type:        schema.TypeMap,
			Description: "Specifies key-value pairs to add to the KV store.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (s *samenessGroup) Decode(d *schema.ResourceData) (consulapi.ConfigEntry, error) {
	configEntry := &consulapi.SamenessGroupConfigEntry{
		Kind:               consulapi.SamenessGroup,
		Name:               d.Get("name").(string),
		Partition:          d.Get("partition").(string),
		DefaultForFailover: d.Get("default_for_failover").(bool),
		IncludeLocal:       d.Get("include_local").(bool),
		Meta:               map[string]string{},
	}

	for k, v := range d.Get("meta").(map[string]interface{}) {
		configEntry.Meta[k] = v.(string)
	}

	for _, raw := range d.Get("members").([]interface{}) {
		m := raw.(map[string]interface{})
		configEntry.Members = append(configEntry.Members, consulapi.SamenessGroupMember{
			Partition: m["partition"].(string),
			Peer:      m["peer"].(string),
		})
	}

	return configEntry, nil
}

func (s *samenessGroup) Write(ce consulapi.ConfigEntry, d *schema.ResourceData, sw *stateWriter) error {
	sp, ok := ce.(*consulapi.SamenessGroupConfigEntry)
	if !ok {
		return fmt.Errorf("expected '%s' but got '%s'", consulapi.ServiceSplitter, ce.GetKind())
	}

	sw.set("name", sp.Name)
	sw.set("partition", sp.Partition)
	sw.set("default_for_failover", sp.DefaultForFailover)
	sw.set("include_local", sp.IncludeLocal)

	meta := map[string]interface{}{}
	for k, v := range sp.Meta {
		meta[k] = v
	}
	sw.set("meta", meta)

	members := make([]interface{}, 0)
	for _, m := range sp.Members {
		member := map[string]interface{}{
			"peer":      m.Peer,
			"partition": m.Partition,
		}
		members = append(members, member)
	}
	sw.set("members", members)

	return sw.error()
}
