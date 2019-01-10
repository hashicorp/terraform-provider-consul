---
layout: "consul"
page_title: "Consul: consul_acl_token"
sidebar_current: "docs-consul-resource-acl-token"
description: |-
  Allows Terraform to create an ACL token
---

# consul_acl_token

The `consul_acl_token` resource writes an ACL token into Consul.

## Example Usage

```hcl
resource "consul_acl_token" "test" {
  description = "my test token"
  policies = ["my_policy"]
  local = true
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the token.

* `policies` - (Optional) The list of policies attached to the token.

* `local` - (Optional) The flag to set the token local to the current datacenter.

## Attributes Reference

The following attributes are exported:

* `id` - The token accessor ID.

* `token` - The token secret ID
