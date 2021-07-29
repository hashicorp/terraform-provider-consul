---
layout: "consul"
page_title: "Consul: consul_namespace"
sidebar_current: "docs-consul-namespace"
description: |-
  Manage a Consul namespace.
---

# consul_namespace

~> **NOTE:** This feature requires Consul Enterprise.

The `consul_namespace` resource provides isolated [Consul Enterprise Namespaces](https://www.consul.io/docs/enterprise/namespaces/index.html).

## Example Usage

```hcl
resource "consul_namespace" "production" {
  name        = "production"
  description = "Production namespace"

  meta = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The namespace name.
* `description` - (Optional) Free form namespace description.
* `policy_defaults` - (Optional) The list of default policies that should be
  applied to all tokens created in this namespace.
* `role_defaults` - (Optional) The list of default roles that should be applied
  to all tokens created in this namespace.
* `meta` - (Optional) Specifies arbitrary KV metadata to associate with the
  namespace.

## Attributes Reference

The following attributes are exported:


* `name` - The namespace name.
* `description` - The namespace description.
* `policy_defaults` - The list of default policies that will be
  applied to all tokens created in this namespace.
* `role_defaults` - The list of default roles that will be applied
  to all tokens created in this namespace.
* `meta` - Arbitrary KV metadata associated with the namespace.

## Import

`consul_namespace` can be imported. This is useful to manage attributes of the
default namespace that is created automatically:

```
$ terraform import consul_namespace.default default
```
