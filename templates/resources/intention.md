---
layout: "consul"
page_title: "Consul: consul_intention"
sidebar_current: "docs-consul-resource-intention"
description: |-
    A resource that can create intentions for Consul Connect.
---

# consul_intention

[Intentions](https://www.consul.io/docs/connect/intentions.html) are used to define
rules for which services may connect to one another when using [Consul Connect](https://www.consul.io/docs/connect/index.html).

~> **NOTE:** This resource is appropriate for managing legacy intentions in
Consul version 1.8 and earlier. As of Consul 1.9, intentions should be managed
using the [`service-intentions`](https://www.consul.io/docs/connect/intentions)
configuration entry. It is recommended to migrate from the `consul_intention`
resource to `consul_config_entry` when running Consul 1.9 and later.

It is appropriate to either reference existing services, or specify non-existent services
that will be created in the future when creating intentions. This resource can be used
in conjunction with the `consul_service` datasource when referencing services
registered on nodes that have a running Consul agent.

## Example Usage

Create a simplest intention with static service names:

```hcl
resource "consul_intention" "database" {
  source_name      = "api"
  destination_name = "db"
  action           = "allow"
}
```

Referencing a known service via a datasource:

```hcl
resource "consul_intention" "database" {
  source_name      = "api"
  destination_name = "${consul_service.pg.name}"
  action           = "allow"
}

data "consul_service" "pg" {
  name = "postgresql"
}
```

## Argument Reference

The following arguments are supported:

* `source_name` - (Required, string) The name of the source service for the intention. This
service does not have to exist.

* `source_namespace` - (Optional, Enterprise Only) The source namespace of the
  intention.

* `destination_name` - (Required, string) The name of the destination service for the intention. This
service does not have to exist.

* `destination_namespace` - (Optional, Enterprise Only) The destination
  namespace of the intention.

* `action` - (Required, string) The intention action. Must be one of `allow` or `deny`.

* `meta` - (Optional, map) Key/value pairs that are opaque to Consul and are associated
with the intention.

* `description` - (Optional, string) Optional description that can be used by Consul
tooling, but is not used internally.

* `datacenter` - (Optional) The datacenter to use. This overrides the
  agent's default datacenter and the datacenter in the provider setup.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the intention.
* `source_name` - The source for the intention.
* `source_namespace` - The source namespace of the intention.
* `destination_name` - The destination for the intention.
* `destination_namespace` - The destination namespace of the intention.
* `action` - The intention action.
* `description` - A description of the intention.
* `meta` - Key/value pairs associated with the intention.
* `datacenter` - The datacenter in which the intention is created.

## Import

`consul_intention` can be imported:

```
$ terraform import consul_intention.database 657a57d6-0d56-57e2-31cb-e9f1ed3c18dd
```
