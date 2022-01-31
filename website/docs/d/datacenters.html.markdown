---
layout: "consul"
page_title: "Consul: consul_datacenters"
sidebar_current: "docs-consul-data-source-datacenters"
description: |-
  Provides the list of all known datacenters.
---

# consul_datacenters

The `consul_datacenters` data source returns the list of all knwown Consul
datacenters.

## Example Usage

```hcl
data "consul_datacenters" "all" {}

# Register a prepared query in each of the datacenters
resource "consul_prepared_query" "myapp-query" {
  for_each = toset(data.consul_datacenters.all.datacenters)

  name         = "myquery"
  datacenter   = each.key
  only_passing = true
  near         = "_agent"

  service = "myapp"
  tags    = ["active", "!standby"]

  failover {
    nearest_n   = 3
    datacenters = ["us-west1", "us-east-2", "asia-east1"]
  }

  dns {
    ttl = "30s"
  }
}
```

## Argument Reference

This data source has no arguments.

## Attributes Reference

The following attributes are exported:

* `datacenters` - The list of datacenters known. The datacenters will be sorted
  in ascending order based on the estimated median round trip time from the server
  to the servers in that datacenter.
