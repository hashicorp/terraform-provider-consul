---
layout: "consul"
page_title: "Consul: consul_admin_partition"
sidebar_current: "docs-consul-admin-partition"
description: |-
  Manage a Consul Admin Partition.
---

# consul_admin_partition

~> **NOTE:** This feature requires Consul Enterprise.

The `consul_admin_partition` resource manages [Consul Enterprise Admin Partitions](https://www.consul.io/docs/enterprise/admin-partitions).

## Example Usage

```hcl
resource "consul_admin_partition" "na_west" {
  name        = "na-west"
  description = "Partition for North America West"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The partition name. This must be a valid DNS hostname label.
* `description` - (Optional) Free form partition description.
*  `disable_gossip` - (Optional). Set to `true` to disable the gossip pool for the partition. Defaults to`false`."

## Attributes Reference

The following attributes are exported:

* `name` - The partition name.
* `description` - The partition description.
* `disable_gossip` - If `true`, the gossip pool is disabled for the partition. Defaults to`false`."

## Import

`consul_admin_partition` can be imported:

```
$ terraform import consul_admin_partition.na_west na-west
```
