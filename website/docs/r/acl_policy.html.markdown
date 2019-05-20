---
layout: "consul"
page_title: "Consul: consul_acl_policy"
sidebar_current: "docs-consul-resource-acl-policy"
description: |-
  Allows Terraform to create an ACL policy
---

# consul_acl_policy

The `consul_acl_policy` resource writes an ACL policy into Consul.

## Example Usage

```hcl
resource "consul_acl_policy" "test" {
  name        = "my_policy"
  datacenters = ["dc1"]
  rules       = <<RULE
    node_prefix "" {
      policy = "read"
    }
    RULE
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy.
* `description` - (Optional) The description of the policy.
* `rules` - (Required) The rules of the policy.
* `datacenters` - (Optional) The datacenters of the policy.
* `request_options` - (Optional) See below.

The `request_options` block supports the following:

* `datacenter` - (Optional) Specify the Consul Datacenter to use when performing the
  request.  This defaults to the datacenter local to the `consul` provider configuration
  but may be overwritten to query a remote datacenter if necessary.

* `token` - (Optional) Specify the Consul ACL token to use when performing the
  request.  This defaults to the same API token configured by the `consul`
  provider but may be overriden if necessary.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the policy.
* `name` - The name of the policy.
* `description` - The description of the policy.
* `rules` - The rules of the policy.
* `datacenters` - The datacenters of the policy.

## Import

`consul_acl_policy` can be imported:

```
$ terraform import consul_acl_policy.my-policy 1c90ef03-a6dd-6a8c-ac49-042ad3752896
```
