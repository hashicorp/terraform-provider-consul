---
layout: "consul"
page_title: "Consul: consul_agent_config"
sidebar_current: "docs-consul-data-source-agent-config"
description: |-
  Provides the configuration information of the local Consul agent.
---

# consul_agent_config

-> **Note:** The `consul_agent_config` resource differs from [`consul_agent_self`](/docs/providers/consul/d/agent_self.html),
providing less information but utilizing stable APIs. `consul_agent_self` will be
deprecated in a future release.

The `consul_agent_config` data source returns
[configuration data](https://www.consul.io/api/agent.html#read-configuration)
from the agent specified in the `provider`.

## Example Usage

```hcl
data "consul_agent_config" "remote_agent" {}

output "info" {
  consul_version = "${data.consul_agent_config.version}"
}
```

## Attributes Reference

The following attributes are exported:

* `datacenter` - The datacenter the agent is running in
* `node_id` - The ID of the node the agent is running on
* `node_name` - The name of the node the agent is running on
* `server` - Boolean if the agent is a server or not
* `revision` - The first 9 characters of the VCS revision of the build of Consul that is running
* `version` - The version of the build of Consul that is running
