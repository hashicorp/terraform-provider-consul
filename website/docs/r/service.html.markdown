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

```hcl
resource "consul_service" "google" {
  address = "www.google.com"
  name    = "google"
  port    = 80
  tags    = ["tag0", "tag1"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the service.

* `address` - (Optional, string) The address of the service. Defaults to the
  address of the agent.

* `service_id` (Optional, string) - If the service ID is not provided, it will be defaulted to the value
of the `name` attribute.

* `port` - (Optional, int) The port of the service.

* `tags` - (Optional, set of strings) A list of values that are opaque to Consul,
  but can be used to distinguish between services or nodes.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service.
* `address` - The address of the service.
* `name` - The name of the service.
* `port` - The port of the service.
* `tags` - The tags of the service.
