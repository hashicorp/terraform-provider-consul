---
layout: "consul"
page_title: "Consul: consul_autopilot_health"
sidebar_current: "docs-consul-data-source-autopilot-health"
description: |-
  Provides health information of the autopilot.
---

# consul_autopilot_health

The `consul_autopilot_health` data source returns
[autopilot health information](https://www.consul.io/api/operator/autopilot.html#read-health)
about the current Consul cluster.

## Example Usage

```hcl
data "consul_autopilot_health" "read" {}

output "health" {
  value = "${data.consul_autopilot_health.read.healthy}"
}
```

## Argument Reference

The following arguments are supported:

* `datacenter` - (Optional) The datacenter to use. This overrides the agent's
  default datacenter and the datacenter in the provider setup.

## Attributes Reference

The following attributes are exported:

* `healthy` - Whether all the servers in the cluster are currently healthy
* `failure_tolerance` - The number of redundant healthy servers that could fail
without causing an outage
* `servers` - A list of server health information. See below for details on the
available information.

### Server health information
* `id` - The Raft ID of the server
* `name` - The node name of the server
* `address` - The address of the server
* `serf_status` - The status of the SerfHealth check of the server
* `version` - The Consul version of the server
* `leader` - Whether the server is currently leader
* `last_contact` - The time elapsed since the server's last contact with
the leader
* `last_term` - The server's last known Raft leader term
* `last_index` - The index of the server's last committed Raft log entry
* `healthy` - Whether the server is healthy according to the current Autopilot
configuration
* `voter` - Whether the server is a voting member of the Raft cluster
* `stable_since` - The time this server has been in its current ``Healthy``
state
