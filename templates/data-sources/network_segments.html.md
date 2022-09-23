---
layout: "consul"
page_title: "Consul: consul_network_segments"
sidebar_current: "docs-consul-data-source-network-segments"
description: |-
  Provides the list of Network Segments.
---

# consul_network_segments

~> **NOTE:** This feature requires [Consul Enterprise](https://www.consul.io/docs/enterprise/index.html).

The `consul_network_segment` data source can be used to retrieve the network
segments defined in the configuration.

## Example Usage

```hcl
data "consul_network_segments" "segments" {}

output "segments" {
  value = data.consul_network_segments.segments.segments
}
```

## Argument Reference

The following arguments are supported:

* `datacenter` - (Optional) The datacenter to use. This overrides the
  agent's default datacenter and the datacenter in the provider setup.
* `token` - (Optional) The ACL token to use. This overrides the
  token that the agent provides by default.

## Attributes Reference

The following attributes are exported:

* `datacenter` - The datacenter the segments are being read from.
* `segments` - The list of network segments.
