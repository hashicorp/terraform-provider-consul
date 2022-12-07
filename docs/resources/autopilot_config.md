---
layout: "consul"
page_title: "Consul: consul_autopilot_config"
sidebar_current: "docs-consul-resource-autopilot-config"
description: |-
  Provides access to the Autopilot Configuration of Consul.
---

# consul_autopilot_config

Provides access to the [Autopilot Configuration](https://www.consul.io/docs/guides/autopilot.html)
of Consul to automatically manage Consul servers.

It includes to automatically cleanup dead servers, monitor the status of the Raft
cluster and stable server introduction.

## Example Usage

```hcl
resource "consul_autopilot_config" "config" {
	cleanup_dead_servers      =  false
	last_contact_threshold    =  "1s"
	max_trailing_logs         =  500
}
```

## Argument Reference

The following arguments are supported:

* `datacenter` - (Optional) The datacenter to use. This overrides the agent's
  default datacenter and the datacenter in the provider setup.

* `cleanup_dead_servers` - (Optional) Whether to remove failing servers when a
replacement comes online. Defaults to true.

* `last_contact_threshold` - (Optional) The time after which a server is
considered as unhealthy and will be removed. Defaults to `"200ms"`.

* `max_trailing_logs` - (Optional) The maximum number of Raft log entries a
server can trail the leader. Defaults to 250.

* `server_stabilization_time` - (Optional) The period to wait for a server to be
healthy and stable before being promoted to a full, voting member. Defaults to
`"10s"`.

* `redundancy_zone_tag` - (Optional) The [redundancy zone](https://www.consul.io/docs/guides/autopilot.html#redundancy-zones)
tag to use. Consul will try to keep one voting server by zone to take advantage
of isolated failure domains. Defaults to an empty string.

* `disable_upgrade_migration` - (Optional) Whether to disable [upgrade migrations](https://www.consul.io/docs/guides/autopilot.html#redundancy-zones).
Defaults to false.

* `upgrade_version_tag` - (Optional) The tag to override the version information
used during a migration. Defaults to an empty string.


## Attributes Reference

The following attributes are exported:

* `datacenter` - The datacenter used.

* `cleanup_dead_servers` - Whether to remove failing servers.

* `last_contact_threshold` - The time after which a server is considered as
unhealthy and will be removed.

* `max_trailing_logs` - The maximum number of Raft log entries a server can trail
the leader.

* `server_stabilization_time` - The period to wait for a server to be healthy and
stable before being promoted to a full, voting member.

* `redundancy_zone_tag` - The [redundancy zone](https://www.consul.io/docs/guides/autopilot.html#redundancy-zones)
tag used. Consul will try to keep one voting server by zone to take advantage of
isolated failure domains.

* `disable_upgrade_migration` - Whether to disable [upgrade migrations](https://www.consul.io/docs/guides/autopilot.html#redundancy-zones).

* `upgrade_version_tag` - The tag to override the version information used during
a migration.
