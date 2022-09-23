---
layout: "consul"
page_title: "Consul: consul_network_area_members"
sidebar_current: "docs-consul-data-source-network-area-members"
description: |-
  Provides the list of Consul servers present in a specific Network Area.
---

# consul_network_area_members

~> **NOTE:** This feature requires [Consul Enterprise](https://www.consul.io/docs/enterprise/index.html).

The `consul_network_area_members` data source provides a list of the Consul
servers present in a specific network area.

## Example Usage

```hcl
resource "consul_network_area" "dc2" {
	peer_datacenter = "dc2"
	retry_join      = ["1.2.3.4"]
	use_tls         = true
}

data "consul_network_area_members" "dc2" {
	uuid = consul_network_area.dc2.id
}

output "members" {
  value = data.consul_network_area_members.dc2.members
}
```

## Argument Reference

The following arguments are supported:

* `datacenter` - (Optional) The datacenter to use. This overrides the
  agent's default datacenter and the datacenter in the provider setup.
* `token` - (Optional) The ACL token to use. This overrides the
  token that the agent provides by default.
* `uuid` - (Required) The UUID of the area to list.

## Attributes Reference

The following attributes are exported:

* `datacenter` - The datacenter used to query the Network Area.
* `uuid` - The UUID of the Network Area being queried.
* `members` - The list of Consul servers in this network area
  * `id` - The node ID of the server.
  * `name` - The node name of the server, with its datacenter appended.
  * `address` - The IP address of the server.
  * `port` - The server RPC port the node.
  * `datacenter` - The node's Consul datacenter.
  * `role` - Role is always `"server"` since only Consul servers can participate
    in network areas.
  * `build` - The Consul version running on the node.
  * `protocol` - The protocol version being spoken by the node.
  * `status` - The current health status of the node, as determined by the
    network area distributed failure detector. This will be `"alive"`, `"leaving"`,
    or `"failed"`. A `"failed"` status means that other servers are not able to
    probe this server over its server RPC interface.
  * `rtt` - An estimated network round trip time from the server answering the
    query to the given server, in nanoseconds. This is computed using network
    coordinates.
