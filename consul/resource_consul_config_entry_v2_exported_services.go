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

func resourceConsulV2ExportedServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulV2ExportedServicesCreate,
		Update: resourceConsulV2ExportedServicesUpdate,
		Read:   resourceConsulV2ExportedServicesRead,
		Delete: resourceConsulV2ExportedServicesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				Required:    true,
				Description: "The partition the config entry is associated with.",
				ForceNew:    true,
			},

			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The namespace the config entry is associated with.",
				ForceNew:    true,
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

func resourceConsulV2ExportedServicesCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceConsulV2ExportedServicesUpdate(d, meta)
}

func resourceConsulV2ExportedServicesUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	name := d.Get("name").(string)
	kind := d.Get("kind").(string)
	gvk := &GVK{
		Group:   "multicluster",
		Version: "v2",
		Kind:    kind,
	}
	var consumers []map[string]any
	for _, p := range d.Get("peer_consumers").([]interface{}) {
		consumers = append(consumers, map[string]any{"peer": p})
	}
	for _, ap := range d.Get("partition_consumers").([]interface{}) {
		consumers = append(consumers, map[string]any{"partition": ap})
	}
	for _, sg := range d.Get("sameness_group_consumers").([]interface{}) {
		consumers = append(consumers, map[string]any{"sameness_group": sg})
	}
	data := map[string]any{"consumers": consumers}
	services := d.Get("services").([]interface{})
	if len(services) > 0 {
		data["services"] = services
	}
	wReq := &V2WriteRequest{
		Metadata: nil,
		Data:     data,
		Owner:    nil,
	}
	resp, _, err := v2MulticlusterApply(client, gvk, name, wOpts, wReq)
	if err != nil || resp == nil {
		return fmt.Errorf("failed to write exported services config '%s': %v", name, err)
	}
	d.SetId(resp.ID.Type.Kind + resp.ID.Tenancy.Partition + resp.ID.Tenancy.Namespace + resp.ID.Name)
	sw := newStateWriter(d)
	sw.set("name", resp.ID.Name)
	sw.set("kind", resp.ID.Type.Kind)
	sw.set("partition", resp.ID.Tenancy.Partition)
	sw.set("namespace", resp.ID.Tenancy.Namespace)
	return resourceConsulV2ExportedServicesRead(d, meta)
}

func resourceConsulV2ExportedServicesRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Get("name").(string)
	kind := d.Get("kind").(string)
	gvk := &GVK{
		Group:   pbmulticluster.GroupName,
		Version: pbmulticluster.Version,
		Kind:    kind,
	}
	resp, err := v2MulticlusterRead(client, gvk, name, qOpts)
	if err != nil || resp == nil {
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
	sw := newStateWriter(d)
	sw.set("services", data.Services)
	sw.set("name", id.Name)
	sw.set("kind", id.Type.Kind)
	sw.set("partition", id.Tenancy.Partition)
	sw.set("namespace", id.Tenancy.Namespace)
	sw.set("partition_consumers", partitions)
	sw.set("peer_consumers", peers)
	sw.set("sameness_group_consumers", samenessgroups)
	return sw.error()
}

func resourceConsulV2ExportedServicesDelete(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	kind := d.Get("kind").(string)
	gvk := &GVK{
		Group:   pbmulticluster.GroupName,
		Version: pbmulticluster.Version,
		Kind:    kind,
	}
	name := d.Get("name").(string)
	return v2MulticlusterDelete(client, gvk, name, qOpts)
}
