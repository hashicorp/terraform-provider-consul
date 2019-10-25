---
layout: "consul"
page_title: "Consul: consul_acl_role"
sidebar_current: "docs-consul-data-source-acl-role"
description: |-
  Provides information about a Consul ACL Role.
---

# consul_acl_role

The `consul_acl_role` data source returns the information related to a
[Consul ACL Role](https://www.consul.io/api/acl/roles.html).

## Example Usage

```hcl
data "consul_acl_role" "test" {
  name = "example-role"
}

output "consul_acl_role" {
  value = data.consul_acl_role.test.id
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ACL Role.

## Attributes Reference

The following attributes are exported:

* `description` - The description of the ACL Role.
* `policies` - The list of policies associated with the ACL Role. Each entry has
  an `id` and a `name` attribute.
* `service_identities` - The list of service identities associated with the ACL
  Role. Each entry has a `service_name` attribute and a list of `datacenters`.
