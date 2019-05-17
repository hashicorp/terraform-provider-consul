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
resource "consul_acl_policy" "agent" {
  name  = "agent"
  rules = <<RULE
    node_prefix "" {
      policy = "read"
    }
    RULE
}

resource "consul_acl_token" "test" {
  description = "my test token"
  policies = ["${consul_acl_policy.agent.name}"]
  local = true
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the token.
* `policies` - (Optional) The list of policies attached to the token.
* `local` - (Optional) The flag to set the token local to the current datacenter.
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

* `id` - The token accessor ID.
* `secret` - The token secret ID.
* `description` - The description of the token.
* `policies` - The list of policies attached to the token.
* `local` - The flag to set the token local to the current datacenter.
