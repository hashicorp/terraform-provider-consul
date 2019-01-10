---
layout: "consul"
page_title: "Consul: consul_acl_master_token"
sidebar_current: "docs-consul-resource-acl-master-token"
description: |-
  Allows Terraform to create an ACL master token
---

# consul_acl_token

The `consul_acl_master_token` resource writes an ACL master token into Consul.

The Consul cluster should not have been bootstrapped to allow its creation.

## Example Usage

```hcl
resource "consul_acl_master_token" "master" {
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
