---
layout: "consul"
page_title: "Consul: consul_namespace_policy_attachment"
sidebar_current: "docs-consul-resource-namespace-policy-attachment"
description: |-
  Allows Terraform to add a policy as a default for a namespace
---

# consul_namespace_policy_attachment

~> **NOTE:** This feature requires Consul Enterprise.

The `consul_namespace_policy_attachment` resource links a Consul Namespace and an ACL
policy. The link is implemented through an update to the Consul Namespace.

~> **NOTE:** This resource is only useful to attach policies to a namespace
that has been created outside the current Terraform configuration, like the
`default` namespace. If the namespace you need to attach a policy to has
been created in the current Terraform configuration and will only be used in it,
you should use the `policy_defaults` attribute of [`consul_namespace`](/docs/providers/consul/r/namespace.html).

## Example Usage

### Attach a policy to the default namespace

```hcl
resource "consul_acl_policy" "agent" {
  name  = "agent"
  rules = <<-RULE
    node_prefix "" {
      policy = "read"
    }
    RULE
}

resource "consul_namespace_policy_attachment" "attachment" {
    namespace = "default"
    policy    = consul_acl_policy.agent.name
}
```

### Attach a policy to a namespace created in another Terraform configuration

#### In `first_configuration/main.tf`

```hcl
resource "consul_namespace" "qa" {
  name = "qa"

  lifecycle {
    ignore_changes = [policy_defaults]
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

resource "consul_namespace_policy_attachment" "attachment" {
    namespace = "qa"
    policy    = consul_acl_policy.agent.name
}
```
**NOTE**: consul_acl_namespace would attempt to enforce an empty set of default
policies, because its `policy_defaults` attribute is empty. For this reason it
is necessary to add the lifecycle clause to prevent Terraform from attempting to
empty the set of policies associated to the namespace.

## Argument Reference

The following arguments are supported:

* `namespace` - (Required) The namespace to attach the policy to.
* `policy` - (Required) The name of the policy attached to the namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The attachment ID.
* `namespace` - The name of the namespace.
* `policy` - The name of the policy attached to the namespace.


## Import

`consul_namespace_policy_attachment` can be imported. This is especially useful
to manage the policies attached to the `default` namespace:

```
$ terraform import consul_namespace_policy_attachment.default default:policy_name
```
