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

### Basic usage

```hcl
resource "consul_acl_policy" "agent" {
  name  = "agent"
  rules = <<-RULE
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

### Set explictly the `accessor_id`

```hcl
resource "random_uuid" "test" { }

resource "consul_acl_token" "test_predefined_id" {
  accessor_id = ${random_uuid.test_uuid.result}
  description = "my test uuid token"
  policies = ["${consul_acl_policy.agent.name}"]
  local = true
}
```

## Argument Reference

The following arguments are supported:

* `accessor_id` - (Optional) The uuid of the token. If omitted, Consul will
  generate a random uuid.
* `description` - (Optional) The description of the token.
* `policies` - (Optional) The list of policies attached to the token.
* `local` - (Optional) The flag to set the token local to the current datacenter.
* `namespace` - (Optional, Enterprise Only) The namespace to create the token within.

## Attributes Reference

The following attributes are exported:

* `id` - The token accessor ID.
* `accessor_id` - The token accessor ID.
* `description` - The description of the token.
* `policies` - The list of policies attached to the token.
* `local` - The flag to set the token local to the current datacenter.


## Import

`consul_acl_token` can be imported. This is especially useful to manage the
anonymous and the master token with Terraform:

```
$ terraform import consul_acl_token.anonymous 00000000-0000-0000-0000-000000000002
$ terraform import consul_acl_token.master-token 624d94ca-bc5c-f960-4e83-0a609cf588be
```
