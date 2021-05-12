package consul

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceConsulNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulNodesRead,
		Schema: map[string]*schema.Schema{
			// Filters
			"query_options": schemaQueryOpts(),

			// Out parameters
			"datacenter": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"node_ids": {
				Computed: true,
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"node_names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nodes": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"meta": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"tagged_addresses": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceConsulNodesRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	// Parse out data source filters to populate Consul's query options
	getQueryOpts(qOpts, d, meta)

	nodes, meta, err := client.Catalog().Nodes(qOpts)
	if err != nil {
		return err
	}

	l := make([]interface{}, 0, len(nodes))

	nodeNames := make([]interface{}, 0, len(nodes))
	nodeIDs := make([]interface{}, 0, len(nodes))

	for _, node := range nodes {
		const defaultNodeAttrs = 4
		m := make(map[string]interface{}, defaultNodeAttrs)
		id := node.ID
		if id == "" {
			id = node.Node
		}

		nodeIDs = append(nodeIDs, id)
		nodeNames = append(nodeNames, node.Node)

		m["address"] = node.Address
		m["id"] = id
		m["name"] = node.Node
		m["meta"] = node.Meta
		m["tagged_addresses"] = node.TaggedAddresses

		l = append(l, m)
	}

	const idKeyFmt = "catalog-nodes-%s"
	d.SetId(fmt.Sprintf(idKeyFmt, qOpts.Datacenter))

	d.Set("datacenter", qOpts.Datacenter)
	if err := d.Set("node_ids", nodeIDs); err != nil {
		return errwrap.Wrapf("Unable to store node IDs: {{err}}", err)
	}

	if err := d.Set("node_names", nodeNames); err != nil {
		return errwrap.Wrapf("Unable to store node names: {{err}}", err)
	}

	if err := d.Set("nodes", l); err != nil {
		return errwrap.Wrapf("Unable to store nodes: {{err}}", err)
	}

	return nil
}
