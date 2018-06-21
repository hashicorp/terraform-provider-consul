---
layout: "consul"
page_title: "Consul: consul_service"
sidebar_current: "docs-consul-resource-service"
description: |-
  A high-level resource for creating a Service in Consul in the Consul catalog.
---

# consul_service

A high-level resource for creating a Service in Consul in the Consul catalog. This
is appropriate for registering [external services](https://www.consul.io/docs/guides/external.html) and
can be used to create services addressable by Consul that cannot be registered
with a [local agent](https://www.consul.io/docs/agent/basics.html).

If the Consul agent is running on the node where this service is registered, it is
not recommended to use this resource.

## Example Usage

Creating a new node with the service:

```hcl
resource "consul_service" "google" {
  name    = "google"
  node    = "${consul_node.compute.name}"
  port    = 80
  tags    = ["tag0"]
}

resource "consul_node" "compute" {
  name    = "compute-google"
  address = "www.google.com"
}
```

Utilizing an existing known node:

```hcl
resource "consul_service" "google" {
  name    = "google"
  node    = "google"
  port    = 443
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the service.

* `node` - (Required, string) The name of the node the to register the service on.

* `address` - (Optional, string) The address of the service. Defaults to the
  address of the node.

* `service_id` (Optional, string) - If the service ID is not provided, it will be defaulted to the value
of the `name` attribute.

* `port` - (Optional, int) The port of the service.

* `tags` - (Optional, set of strings) A list of values that are opaque to Consul,
  but can be used to distinguish between services or nodes.

* `datacenter` - (Optional) The datacenter to use. This overrides the datacenter in the
provider setup and the agent's default datacenter.

## Attributes Reference

The following attributes are exported:

* `service_id` - The ID of the service.
* `address` - The address of the service.
* `node` - The node the service is registered on.
* `name` - The name of the service.
* `port` - The port of the service.
* `tags` - The tags of the service.
