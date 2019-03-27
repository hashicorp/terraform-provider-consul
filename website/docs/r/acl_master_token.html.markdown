---
layout: "consul"
page_title: "Consul: consul_acl_master_token"
sidebar_current: "docs-consul-resource-acl-master-token"
description: |-
  Allows Terraform to create an ACL master token
---

# consul_acl_master_token

The `consul_acl_master_token` resource writes an ACL master token into Consul and
can be used to bootstrap ACLs support in a Consul cluster.

The Consul cluster should not have been bootstrapped to allow its creation.

!> **WARNING**: When using this resource, the Consul master token will be written
   as plaintext in the terraform state. You must make sure it is safely stored and
  protected against external access.

## Example Usage

```hcl
resource "consul_acl_master_token" "master" {}
```

## Attributes Reference

The following attributes are exported:

* `id` - The token accessor ID.
* `description` - The description of the token.
* `policies` - The list of policies attached to the token.
* `local` - The flag to set the token local to the current datacenter.
* `token` - The token secret ID.
