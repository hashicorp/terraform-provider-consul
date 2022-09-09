---
layout: "consul"
page_title: "Consul: consul_node"
sidebar_current: "docs-consul-resource-node"
description: |-
  Provides access to Node data in Consul. This can be used to define a node.
---

# consul_node

Provides access to Node data in Consul. This can be used to define a
node. Currently, defining health checks is not supported.

## Example Usage

```hcl
resource "consul_node" "foobar" {
  address = "192.168.10.10"
  name    = "foobar"
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Required) The address of the node being added to, or referenced in the catalog.
* `name` - (Required) The name of the node being added to, or referenced in the catalog.
* `datacenter` - (Optional) The datacenter to use. This overrides the agent's default datacenter and the datacenter in the provider setup.
* `meta` - (Optional, map) Key/value pairs that are associated with the node.
* `partition` - (Optional, Enterprise Only) The partition the node is associated with.

## Attributes Reference

The following attributes are exported:

* `address` - The address of the node.
* `name` - The name of the node.
* `meta` - (Optional, map) Key/value pairs that are associated with the node.

## Import

`consul_node` can be imported:

```
$ terraform import consul_node.example $name
```
