---
layout: "consul"
page_title: "Consul: consul_acl_token_role_attachment"
sidebar_current: "docs-consul-resource-acl-tr-attachment"
description: |-
  Allows Terraform to create a link between an ACL token and a role
---

# consul_acl_token_role_attachment

The `consul_acl_token_role_attachment` resource links a Consul Token and an ACL
role. The link is implemented through an update to the Consul ACL token.

~> **NOTE:** This resource is only useful to attach roles to an ACL token
that has been created outside the current Terraform configuration, like the
anonymous or the master token. If the token you need to attach a policy to has
been created in the current Terraform configuration and will only be used in it,
you should use the `roles` attribute of [`consul_acl_token`](/docs/providers/consul/r/acl_token.html).

## Example Usage

### Attach a role to the anonymous token

```hcl
resource "consul_acl_role" "role" {
  name = "foo"
  description = "Foo"

  service_identities {
    service_name = "foo"
  }
}

resource "consul_acl_token_role_attachment" "attachment" {
  token_id = "00000000-0000-0000-0000-000000000002"
  role_id  = consul_acl_role.role.id
}
```

### Attach a policy to a token created in another Terraform configuration

#### In `first_configuration/main.tf`

```hcl
resource "consul_acl_token" "test" {
  accessor_id = "5914ee49-eb8d-4837-9767-9299ec155000"
  description = "my test token"
  local = true

  lifecycle {
    ignore_changes = ["roles"]
  }
}
```

#### In `second_configuration/main.tf`

```hcl
resource "consul_acl_role" "role" {
  name = "foo"
  description = "Foo"

  service_identities {
    service_name = "foo"
  }
}

resource "consul_acl_token_role_attachment" "attachment" {
  token_id = "00000000-0000-0000-0000-000000000002"
  role_id  = consul_acl_role.role.id
}
```
**NOTE**: `consul_acl_token` would attempt to enforce an empty set of roles,
because its `roles` attribute is empty. For this reason it is necessary to add
the lifecycle clause to prevent Terraform from attempting to clear the set of
roles associated to the token.

## Argument Reference

The following arguments are supported:

* `token_id` - (Required) The id of the token.
* `role_id` - (Required) The id of the role to attach to the token.

## Attributes Reference

The following attributes are exported:

* `id` - The attachment ID.
* `token_id` - The id of the token.
* `role_id` - The id of the role attached to the token.


## Import

`consul_acl_token_role_attachment` can be imported. This is especially useful to manage the
policies attached to the anonymous and the master tokens with Terraform:

```
$ terraform import consul_acl_token_role_attachment.anonymous token_id:role_id
```
