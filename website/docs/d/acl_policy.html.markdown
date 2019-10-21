---
layout: "consul"
page_title: "Consul: consul_acl_policy"
sidebar_current: "docs-consul-data-source-acl-policy"
description: |-
  Provides information about a Consul ACL Poliy.
---

# consul_acl_policy

The `consul_acl_policy` data source returns the information related to a
[Consul ACL Policy](https://www.consul.io/docs/acl/acl-system.html#acl-policies).


## Example Usage

```hcl
data "consul_acl_policy" "agent" {
  name  = "agent"
}

output "consul_acl_policy" {
  value = data.consul_acl_policy.agent.rules
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ACL Policy.

## Attributes Reference

The following attributes are exported:

* `description` - The description of the ACL Policy.
* `rules` - The rules associated with the ACL Policy.
* `datacenters` - The datacenters associated with the ACL Policy.
