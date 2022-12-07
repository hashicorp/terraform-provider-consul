---
layout: "consul"
page_title: "Consul: consul_network_area"
sidebar_current: "docs-consul-resource-network-area"
description: |-
  Manage Network Areas.
---

# consul_network_area

~> **NOTE:** This feature requires [Consul Enterprise](https://www.consul.io/docs/enterprise/index.html).

The `consul_network_area` resource manages a relationship between servers in two
different Consul datacenters.

Unlike Consul's WAN feature, network areas use just the server RPC port for
communication, and relationships can be made between independent pairs of
datacenters, so not all servers need to be fully connected. This allows for
complex topologies among Consul datacenters like hub/spoke and more general trees.

## Example Usage

```hcl
resource "consul_network_area" "dc2" {
	peer_datacenter = "dc2"
	retry_join      = ["1.2.3.4"]
	use_tls         = true
}
```

## Argument Reference

The following arguments are supported:

* `datacenter` - (Optional) The datacenter to use. This overrides the
  agent's default datacenter and the datacenter in the provider setup.
* `token` - (Optional) The ACL token to use. This overrides the
  token that the agent provides by default.
* `peer_datacenter` - (Required) The name of the Consul datacenter that will be
  joined to form the area.
* `retry_join` - (Optional) Specifies a list of Consul servers to attempt to
  join. Servers can be given as `IP`, `IP:port`, `hostname`, or `hostname:port`.
* `use_tls` - (Optional) Specifies whether gossip over this area should be
  encrypted with TLS if possible. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `datacenter` - The datacenter being queried.
* `peer_datacenter` - The name of the Consul datacenter joined to form the area.
* `retry_join` - The list of Consul servers Consul attempts to join.
* `use_tls` - Whether the gossip over this area should be encrypted with TLS.
