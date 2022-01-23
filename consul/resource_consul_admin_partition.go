package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulAdminPartition() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulAdminPartitionCreate,
		Read:   resourceConsulAdminPartitionRead,
		Update: resourceConsulAdminPartitionUpdate,
		Delete: resourceConsulAdminPartitionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The partition name. This must be a valid DNS hostname label.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Free form partition description.",
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceConsulAdminPartitionCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	partitions := client.Partitions()
	name := d.Get("name").(string)

	partition := &api.Partition{
		Name:        name,
		Description: d.Get("description").(string),
	}

	_, _, err := partitions.Create(context.TODO(), partition, wOpts)
	if err != nil {
		return fmt.Errorf("failed to create partition %q: %w", name, err)
	}

	d.SetId(name)

	return resourceConsulAdminPartitionRead(d, meta)
}

func resourceConsulAdminPartitionRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	partitions := client.Partitions()
	name := d.Id()

	partition, _, err := partitions.Read(context.TODO(), name, qOpts)
	if err != nil {
		return fmt.Errorf("failed to read partition %q: %w", name, err)
	}
	if partition == nil {
		// The partition has been removed
		d.SetId("")
		return nil
	}

	sw := newStateWriter(d)
	sw.set("name", partition.Name)
	sw.set("description", partition.Description)

	return sw.error()
}

func resourceConsulAdminPartitionUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	partitions := client.Partitions()
	name := d.Get("name").(string)

	partition := &api.Partition{
		Name:        name,
		Description: d.Get("description").(string),
	}

	_, _, err := partitions.Update(context.TODO(), partition, wOpts)
	if err != nil {
		return fmt.Errorf("failed to update partition %q: %w", name, err)
	}

	return resourceConsulAdminPartitionRead(d, meta)
}

func resourceConsulAdminPartitionDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	partitions := client.Partitions()
	name := d.Get("name").(string)

	if _, err := partitions.Delete(context.TODO(), name, wOpts); err != nil {
		return fmt.Errorf("failed to delete partition %q: %w", name, err)
	}

	return nil
}
