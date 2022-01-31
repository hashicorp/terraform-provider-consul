---
layout: "consul"
page_title: "Consul: consul_datacenters"
sidebar_current: "docs-consul-data-source-datacenters"
description: |-
  Provides a list of datacenters in a given Consul cluster
---

# consul_datacenters

The `consul_datacenters` data source returns a list of Consul datacenters that are
available within the Consul cluster. 

## Example Usage

```hcl
data "consul_datacenters" "read-datacenters" {}

# Loop over all datacenters to create an entry in each
resource "example_resource" "app" {
  for_each = toset(data.consul_datacenters.read-datacenters.datacenters)

  # ...
}
```

## Argument Reference

This data source does not support any arguments

## Attributes Reference

The following attributes are exported:

* `datacenters` - A list of datacenters within the cluster
