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

Register a health-check:

```hcl
resource "consul_service" "redis" {
  name = "redis"
  node = "redis"
  port = 6379

  check {
    check_id                          = "service:redis1"
    name                              = "Redis health check"
    status                            = "passing"
    http                              = "https://www.hashicorptest.com"
    tls_skip_verify                   = false
    method                            = "PUT"
    interval                          = "5s"
    timeout                           = "1s"
    deregister_critical_service_after = "30s"

    header {
      name  = "foo"
      value = ["test"]
    }

    header {
      name  = "bar"
      value = ["test"]
    }
  }
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

* `checks` - (Optional, list of checks) Health-checks to register to monitor the
  service. The list of attributes for each health-check is detailled below.

* `tags` - (Optional, set of strings) A list of values that are opaque to Consul,
  but can be used to distinguish between services or nodes.

* `datacenter` - (Optional) The datacenter to use. This overrides the
  agent's default datacenter and the datacenter in the provider setup.

* `meta` - (Optional) A map of arbitrary KV metadata linked to the service
  instance.

* `namespace` - (Optional, Enterprise Only) The namespace to create the service within.

The following attributes are available for each health-check:

* `check_id` - (Optional, string) An ID, *unique per agent*. Will default to *name*
  if not set.
* `name` - (Required) The name of the health-check.
* `notes` - (Optional, string) An opaque field meant to hold human readable text.
* `status` - (Optional, string) The initial health-check status.
* `tcp` - (Optional, string) The TCP address and port to connect to for a TCP check.
* `http` - (Optional, string) The HTTP endpoint to call for an HTTP check.
* `header` - (Optional, set of headers) The headers to send for an HTTP check.
  The attributes of each header is given below.
* `tls_skip_verify` - (Optional, boolean) Whether to deactivate certificate
  verification for HTTP health-checks. Defaults to `false`.
* `method` - (Optional, string) The method to use for HTTP health-checks. Defaults
  to `GET`.
* `interval` - (Required, string) The interval to wait between each health-check
  invocation.
* `timeout` - (Required, string) The timeout value for HTTP checks.
* `deregister_critical_service_after` - (Optional, string) The time after which
  the service is automatically deregistered when in the `critical` state.
  Defaults to `30s`.

Each `header` must have the following attributes:
* `name` - (Required, string) The name of the header.
* `value` - (Required, list of strings) The header's list of values.

## Attributes Reference

The following attributes are exported:

* `service_id` - The ID of the service.
* `address` - The address of the service.
* `node` - The node the service is registered on.
* `name` - The name of the service.
* `port` - The port of the service.
* `tags` - The tags of the service.
* `checks` - The list of health-checks associated with the service.
* `datacenter` - The datacenter of the service.
* `meta` - A map of arbitrary KV metadata linked to the service instance.
