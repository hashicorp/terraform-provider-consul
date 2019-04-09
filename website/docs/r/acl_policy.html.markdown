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

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the policy.
* `name` - The name of the policy.
* `description` - The description of the policy.
* `rules` - The rules of the policy.
* `datacenters` - The datacenters of the policy.