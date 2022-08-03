package consul

import (
	"context"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulPartition() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulNamespaceCreate,
		Read:   resourceConsulNamespaceRead,
		Update: resourceConsulNamespaceUpdate,
		Delete: resourceConsulNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceConsulPartitionCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	partition := getPartitionFromResourceData(d)
	partition, _, err := client.Partitions().Create(context.Background(), partition, wOpts)
	if err != nil {
		return fmt.Errorf("failed to create namespace: %v", err)
	}
	d.SetId(partition.Name)
	return resourceConsulPartitionRead(d, meta)
}

func resourceConsulPartitionRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	name := d.Id()

	partition, _, err := client.Partitions().Read(context.Background(), name, qOpts)
	if partition == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to read partition '%s': %v", name, err)
	}

	sw := newStateWriter(d)
	sw.set("name", partition.Name)
	sw.set("description", partition.Description)

	return sw.error()
}

func resourceConsulPartitionUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	partition := getPartitionFromResourceData(d)
	partition, _, err := client.Partitions().Update(context.Background(), partition, wOpts)
	if err != nil {
		return fmt.Errorf("failed to update partition '%s': %v", partition.Name, err)
	}

	return resourceConsulPartitionRead(d, meta)
}

func resourceConsulPartitionDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)

	_, err := client.Partitions().Delete(context.Background(), d.Id(), wOpts)
	if err != nil {
		return fmt.Errorf("failed to delete namespace '%s': %v", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getPartitionFromResourceData(d *schema.ResourceData) *consulapi.Partition {
	return &consulapi.Partition{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}
