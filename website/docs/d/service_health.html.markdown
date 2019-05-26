---
layout: "consul"
page_title: "Consul: consul_service_health"
sidebar_current: "docs-consul-data-source-health"
description: |-
  Filter service instances based on health status
---

# consul_service_health

`consul_service_health` can be used to get the list of the instances that
are currently healthy, according to their associated  health-checks.
The result includes the list of service instances, the node associated to each
instance and its health-checks.

## Example Usage

```hcl
provider "consul" {}

data "consul_service_health" "vault" {
  service = "vault"
  passing = true
}

provider "vault" {
  address = "https://${data.consul_service_health.vault.results.0.service.0.address}:${data.consul_service_health.vault.results.0.service.0.port}"
}
```

## Argument Reference

The following arguments are supported:

* `datacenter` - (Optional) The Consul datacenter to query.

* `name` - (Required) The service name to select.

* `near` - (Optional) Specifies a node name to sort the node list in ascending order
  based on the estimated round trip time from that node.

* `tag` - (Optional) A single tag that can be used to filter the list to return
   based on a single matching tag.

* `node_meta` - (Optional) Filter the results to nodes with the specified key/value
  pairs.

* `passing` - (Optional) Whether to return only nodes with all checks in the
  passing state. Defaults to `true`.

## Attributes Reference

The following attributes are exported:

* `datacenter` - The datacenter the keys are being read from to.
* `name` - The name of the service.
* `near` - The node to which the result must be sorted to.
* `tag` - The name of the tag used to filter the list.
* `node_meta` - The list of metadata to filter the nodes.
* `passing` - Whether to return only nodes with all checks in the
  passing state.
* `results` - A list of entries and details about each endpoint advertising a
  service.  Each element in the list has three attributes: `node`, `service` and
  `checks`.  The list of the attributes of each one is detailed below.



The following is a list of the per-entry `node` attributes:

* `id` - The Node ID of the Consul node advertising the service.
* `name` - The name of the node.
* `address` - The address of the node.
* `datacenter` - The datacenter in which the node is running.
* [`tagged_addresses`](https://www.consul.io/docs/agent/http/catalog.html#TaggedAddresses) -
  List of explicit LAN and WAN IP addresses for the agent.
* `meta` - Node meta data tag information, if any.


The following is a list of the per-entry `service` attributes:

* `id` - The ID of the service.
* `name` - The name of the service.
* `tags` - The list of tags associated with this instance.
* `address` - The address of this instance.
* `port` - The port of this instance.
* `meta` - Service metadata tag information, if any.


`checks` is a list of the health-checks associated to the entry with the
following attributes:

* `id` - The ID of this health-check.
* `node` - The name of the node associated with this health-check.
* `name` - The name of this health-check.
* `status` - The status of this health-check.
* `notes` - A human readable description of the current state of the health-check.
* `output` - The output of the health-check.
* `service_id` - The ID of the service associated to this health-check.
* `service_name` - The name of the service associated with this health-check.
* `service_tags` - The list of tags associated with this health-check.