package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	catalogNodesElem       = "nodes"
	catalogNodesDatacenter = "datacenter"
	catalogNodesQueryOpts  = "query_options"

	catalogNodesNodeID              = "id"
	catalogNodesNodeAddress         = "address"
	catalogNodesNodeMeta            = "meta"
	catalogNodesNodeName            = "name"
	catalogNodesNodeTaggedAddresses = "tagged_addresses"

	catalogNodesNodeIDs   = "node_ids"
	catalogNodesNodeNames = "node_names"

	catalogNodesAPITaggedLAN    = "lan"
	catalogNodesAPITaggedWAN    = "wan"
	catalogNodesSchemaTaggedLAN = "lan"
	catalogNodesSchemaTaggedWAN = "wan"
)

func dataSourceConsulNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulNodesRead,
		Schema: map[string]*schema.Schema{
			// Filters
			catalogNodesQueryOpts: schemaQueryOpts,

			// Out parameters
			catalogNodesDatacenter: {
				Computed: true,
				Type:     schema.TypeString,
			},
			catalogNodesNodeIDs: {
				Computed: true,
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			catalogNodesNodeNames: {
				Computed: true,
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			catalogNodesElem: {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						catalogNodesNodeID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						catalogNodesNodeName: {
							Type:     schema.TypeString,
							Computed: true,
						},
						catalogNodesNodeAddress: {
							Type:     schema.TypeString,
							Computed: true,
						},
						catalogNodesNodeMeta: {
							Type:     schema.TypeMap,
							Computed: true,
						},
						catalogNodesNodeTaggedAddresses: {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									catalogNodesSchemaTaggedLAN: {
										Type:     schema.TypeString,
										Computed: true,
									},
									catalogNodesSchemaTaggedWAN: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceConsulNodesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*consulapi.Client)

	// Parse out data source filters to populate Consul's query options
	queryOpts, err := getQueryOpts(d, client)
	if err != nil {
		return errwrap.Wrapf("unable to get query options for fetching catalog nodes: {{err}}", err)
	}

	nodes, meta, err := client.Catalog().Nodes(queryOpts)
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

		m[catalogNodesNodeAddress] = node.Address
		m[catalogNodesNodeID] = id
		m[catalogNodesNodeName] = node.Node
		m[catalogNodesNodeMeta] = node.Meta
		m[catalogNodesNodeTaggedAddresses] = node.TaggedAddresses

		l = append(l, m)
	}

	const idKeyFmt = "catalog-nodes-%s"
	d.SetId(fmt.Sprintf(idKeyFmt, queryOpts.Datacenter))

	d.Set(catalogNodesDatacenter, queryOpts.Datacenter)
	if err := d.Set(catalogNodesNodeIDs, nodeIDs); err != nil {
		return errwrap.Wrapf("Unable to store node IDs: {{err}}", err)
	}

	if err := d.Set(catalogNodesNodeNames, nodeNames); err != nil {
		return errwrap.Wrapf("Unable to store node names: {{err}}", err)
	}

	if err := d.Set(catalogNodesElem, l); err != nil {
		return errwrap.Wrapf("Unable to store nodes: {{err}}", err)
	}

	return nil
}
