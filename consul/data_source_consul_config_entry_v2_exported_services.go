// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"encoding/json"
	"fmt"

	pbmulticluster "github.com/hashicorp/consul/proto-public/pbmulticluster/v2"
	"github.com/hashicorp/consul/proto-public/pbresource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/protobuf/encoding/protojson"
)

func dataSourceConsulConfigEntryV2ExportedServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulV2ExportedServicesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the config entry to read.",
			},

			"kind": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The kind of exported services config (ExportedServices, NamespaceExportedServices, PartitionExportedServices).",
			},

			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The partition the config entry is associated with.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The namespace the config entry is associated with.",
			},

			"services": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The exported services.",
			},

			"partition_consumers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The exported service partition consumers.",
			},
			"peer_consumers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The exported service peer consumers.",
			},
			"sameness_group_consumers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The exported service sameness group consumers.",
			},
		},
	}
}

func dataSourceConsulV2ExportedServicesRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Get("name").(string)
	kind := d.Get("kind").(string)
	gvk := &GVK{
		Group:   pbmulticluster.GroupName,
		Version: pbmulticluster.Version,
		Kind:    kind,
	}
	resp, err := v2MulticlusterRead(client, gvk, name, qOpts)
	if err != nil || resp == nil || resp["id"] == nil || resp["data"] == nil {
		return fmt.Errorf("exported services config not found: %s", name)
	}
	respData, err := json.Marshal(resp["data"])
	if err != nil {
		return fmt.Errorf("failed to unmarshal response data: %v", err)
	}
	data := &pbmulticluster.ExportedServices{}
	if err = protojson.Unmarshal(respData, data); err != nil {
		return fmt.Errorf("failed to unmarshal to proto message: %v", err)
	}
	respID, err := json.Marshal(resp["id"])
	if err != nil {
		return fmt.Errorf("failed to unmarshal response id: %v", err)
	}
	id := &pbresource.ID{}
	if err = protojson.Unmarshal(respID, id); err != nil {
		return fmt.Errorf("failed to unmarshal to proto message: %v", err)
	}
	var partitions []string
	var peers []string
	var samenessgroups []string
	for _, e := range data.Consumers {
		switch v := e.ConsumerTenancy.(type) {
		case *pbmulticluster.ExportedServicesConsumer_Peer:
			peers = append(peers, v.Peer)
		case *pbmulticluster.ExportedServicesConsumer_Partition:
			partitions = append(partitions, v.Partition)
		case *pbmulticluster.ExportedServicesConsumer_SamenessGroup:
			samenessgroups = append(samenessgroups, v.SamenessGroup)
		default:
			return fmt.Errorf("unknown exported service consumer type: %T", v)
		}
	}
	d.SetId(id.Uid)
	sw := newStateWriter(d)
	sw.set("services", data.Services)
	sw.set("partition_consumers", partitions)
	sw.set("peer_consumers", peers)
	sw.set("sameness_group_consumers", samenessgroups)
	return sw.error()
}
