---
layout: "consul"
page_title: "Consul: consul_acl_token_policy_attachment"
sidebar_current: "docs-consul-resource-acl-tp-attachment"
description: |-
  Allows Terraform to create a link between an ACL token and a policy
---

# consul_acl_token_policy_attachment

The `consul_acl_token_attachment` resource links a Consul Token and an ACL
policy. The link is implemented through an update to the Consul ACL token.

~> **NOTE:** This resource is only useful to attach policies to an ACL token
that has been created outside the current Terraform configuration, like the
anonymous or the master token. If the token you need to attach a policy to has
been created in the current Terraform configuration and will only be used in it,
you should use the `policies` attribute of [`consul_acl_token`](/docs/providers/consul/r/acl_token.html).

## Example Usage

### Attach a policy to the anonymous token

```hcl
resource "consul_acl_policy" "agent" {
  name  = "agent"
  rules = <<-RULE
    node_prefix "" {
      policy = "read"
    }
    RULE
}

resource "consul_acl_token_policy_attachment" "attachment" {
    token_id = "00000000-0000-0000-0000-000000000002"
    policy   = "${consul_acl_policy.agent.name}"
}
```

### Attach a policy to a token created in another Terraform configuration

#### In `first_configuration/main.tf`

```hcl
resource "consul_acl_token" "test" {
  accessor_id = "9b20de68-3ea2-4b70-b4f1-506afad062a4"
  description = "my test token"
  local = true

  lifecycle {
    ignore_changes = ["policies"]
  }
}
```

#### In `second_configuration/main.tf`

```hcl
resource "consul_acl_policy" "agent" {
  name  = "agent"
  rules = <<-RULE
    node_prefix "" {
      policy = "read"
    }
    RULE
}

resource "consul_acl_token_policy_attachment" "attachment" {
    token_id = "9b20de68-3ea2-4b70-b4f1-506afad062a4"
    policy   = "${consul_acl_policy.agent.name}"
}
```
**NOTE**: consul_acl_token would attempt to enforce an empty set of policies,
because its policies attribute is empty. For this reason it is necessary to add
the lifecycle clause to prevent Terraform from attempting to empty the set of
policies associated to the token.

## Argument Reference

The following arguments are supported:

* `token_id` - (Required) The id of the token.
* `policy` - (Required) The name of the policy attached to the token.

## Attributes Reference

The following attributes are exported:

* `id` - The attachment ID.
* `token_id` - The id of the token.
* `policy` - The name of the policy attached to the token.


## Import

`consul_acl_token_policy_attachment` can be imported. This is especially useful to manage the
policies attached to the anonymous and the master tokens with Terraform:

```
$ terraform import consul_acl_token_policy_attachment.anonymous 00000000-0000-0000-0000-000000000002:policy_name
$ terraform import consul_acl_token_policy_attachment.master-token 624d94ca-bc5c-f960-4e83-0a609cf588be:policy_name
```
