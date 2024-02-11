// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/consul/api"
	pbmulticluster "github.com/hashicorp/consul/proto-public/pbmulticluster/v2"
	"github.com/hashicorp/consul/proto-public/pbresource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/protobuf/encoding/protojson"

	multicluster "github.com/hashicorp/terraform-provider-consul/consul/tools/openapi"
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
	client, _, _ := getMulticlusterV2Client(d, meta)
	name := d.Get("name").(string)
	kind := d.Get("kind").(string)
	partition := d.Get("partition").(string)
	namespace := d.Get("namespace").(string)
	var consumers []multicluster.HashicorpConsulMulticlusterV2ExportedServicesConsumer
	peerConsumers := d.Get("peer_consumers").([]interface{})
	for _, p := range peerConsumers {
		peerString := p.(string)
		consumers = append(consumers, multicluster.HashicorpConsulMulticlusterV2ExportedServicesConsumer{
			Peer: &peerString,
		})
	}
	partitionConsumers := d.Get("partition_consumers").([]interface{})
	for _, ap := range partitionConsumers {
		partitionString := ap.(string)
		consumers = append(consumers, multicluster.HashicorpConsulMulticlusterV2ExportedServicesConsumer{
			Peer: &partitionString,
		})
	}
	samenessConsumers := d.Get("sameness_group_consumers").([]interface{})
	for _, sg := range samenessConsumers {
		sgString := sg.(string)
		consumers = append(consumers, multicluster.HashicorpConsulMulticlusterV2ExportedServicesConsumer{
			SamenessGroup: &sgString,
		})
	}
	services := d.Get("services").([]interface{})
	var servicesData []string
	for _, s := range services {
		servicesData = append(servicesData, s.(string))
	}
	resp, err := doWriteForKind(client, name, kind, namespace, partition, servicesData, consumers)
	if err != nil || resp == nil {
		return fmt.Errorf("failed to write exported services config '%s': %v", name, err)
	}

	// Probably should parse the response body to get this instead of just relying on OK response
	d.SetId(kind + partition + namespace + name)
	sw := newStateWriter(d)
	sw.set("name", name)
	sw.set("kind", kind)
	sw.set("partition", partition)
	sw.set("namespace", namespace)
	return resourceConsulV2ExportedServicesRead(d, meta)
}

func doWriteForKind(client *multicluster.Client, name string, kind string, namespace string, partition string, services []string, consumers []multicluster.HashicorpConsulMulticlusterV2ExportedServicesConsumer) (*http.Response, error) {
	group := "multicluster"
	gv := "v2"
	gvk := &multicluster.HashicorpConsulResourceType{
		Group:        &group,
		GroupVersion: &gv,
		Kind:         &kind,
	}
	id := &multicluster.HashicorpConsulResourceID{
		Name: &name,
		Tenancy: &multicluster.HashicorpConsulResourceTenancy{
			Namespace: &namespace,
			Partition: &partition,
		},
		Type: gvk,
	}

	var resp *http.Response
	var err error
	switch kind {
	case "ExportedServices":
		wParams := &multicluster.WriteExportedServicesParams{
			Peer:      nil,
			Namespace: &namespace,
			Ns:        nil, // why is this a thing?
			Partition: &partition,
		}
		body := multicluster.WriteExportedServicesJSONRequestBody{
			Data: &multicluster.HashicorpConsulMulticlusterV2ExportedServices{
				Consumers: &consumers,
				Services:  &services,
			},
			Id: id,
		}
		resp, err = client.WriteExportedServices(context.Background(), name, wParams, body, nil)
	case "NamespaceExportedServices":
		wParams := &multicluster.WriteNamespaceExportedServicesParams{
			Peer:      nil,
			Namespace: &namespace,
			Ns:        nil, // why is this a thing?
			Partition: &partition,
		}
		body := multicluster.WriteNamespaceExportedServicesJSONRequestBody{
			Data: &multicluster.HashicorpConsulMulticlusterV2NamespaceExportedServices{
				Consumers: &consumers,
			},
			Id: id,
		}
		resp, err = client.WriteNamespaceExportedServices(context.Background(), name, wParams, body, nil)
	case "PartitionExportedServices":
		wParams := &multicluster.WritePartitionExportedServicesParams{
			Peer:      nil,
			Partition: &partition,
		}
		body := multicluster.WritePartitionExportedServicesJSONRequestBody{
			Data: &multicluster.HashicorpConsulMulticlusterV2PartitionExportedServices{
				Consumers: &consumers,
			},
			Id: id,
		}
		resp, err = client.WritePartitionExportedServices(context.Background(), name, wParams, body)
	}
	return resp, err
}

func resourceConsulV2ExportedServicesRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Get("name").(string)
	kind := d.Get("kind").(string)
	gvk := &api.GVK{
		Group:   "multicluster",
		Version: "v2",
		Kind:    kind,
	}
	resp, err := client.Resource().Read(gvk, name, qOpts)
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
		return fmt.Errorf("Failed to unmarshal to proto message: %v", err)
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
	gvk := &api.GVK{
		Group:   "multicluster",
		Version: "v2",
		Kind:    "ExportedServices",
	}
	name := d.Get("name").(string)
	return client.Resource().Delete(gvk, name, qOpts)
}
