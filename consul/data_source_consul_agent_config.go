package consul

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceConsulAgentConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsulAgentConfigRead,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Description: "The datacenter the agent is running in",
				Computed:    true,
			},

			"node_id": {
				Type:        schema.TypeString,
				Description: "The ID of the node the agent is running on",
				Computed:    true,
			},

			"node_name": {
				Type:        schema.TypeString,
				Description: "The name of the node the agent is running on",
				Computed:    true,
			},

			"server": {
				Type:        schema.TypeBool,
				Description: "If the agent is a server or not",
				Computed:    true,
			},

			"revision": {
				Type:        schema.TypeString,
				Description: "The VCS revision of the build of Consul that is running",
				Computed:    true,
			},

			"version": {
				Type:        schema.TypeString,
				Description: "The version of the build of Consul that is running",
				Computed:    true,
			},
		},
	}
}

func dataSourceConsulAgentConfigRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).Client()
	if err != nil {
		return err
	}
	agentSelf, err := client.Agent().Self()
	if err != nil {
		return err
	}

	config, ok := agentSelf["Config"]
	if !ok {
		return fmt.Errorf("Config key not present on agent self endpoint")
	}

	// We use the ID of the node as the datasource ID, as the datasource
	// queries config from the agent running on that registered node, so it
	// is the best we can do to get a consistent identifier
	d.SetId(fmt.Sprintf("agent-%s", config["NodeId"]))

	d.Set("datacenter", config["Datacenter"])
	d.Set("node_id", config["NodeID"])
	d.Set("node_name", config["NodeName"])
	d.Set("server", config["Server"])
	d.Set("revision", config["Revision"])
	d.Set("version", config["Version"])

	return nil
}
