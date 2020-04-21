---
layout: "consul"
page_title: "Consul: consul_acl_token_secret_id"
sidebar_current: "docs-consul-data-source-acl-tok-secret-id"
description: |-
  Provides the ACL Token secret ID.
---

# consul_acl_token_secret_id

~> **Warning:** When using this is resource, the ACL Token secret ID will be
written to the Terraform state. It is strongly recommended to use the `pgp_key`
attribute and to make sure the remote state has strong access controls before
using this resource.

The `consul_acl_token_secret` data source returns the secret ID associated to
the accessor ID. This can be useful to make systems that cannot use an auth
method to interface with Consul.

If you want to get other attributes of the Consul ACL token, please use the
[`consul_acl_token` data source](/docs/providers/consul/d/acl_token.html).

## Example Usage

```hcl
resource "consul_acl_policy" "test" {
	name        = "test"
	rules       = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies    = [consul_acl_policy.test.name]
	local       = true
}

data "consul_acl_token_secret_id" "read" {
    accessor_id = consul_acl_token.test.id
	pgp_key     = "keybase:my_username"
}

output "consul_acl_token_secret_id" {
  value = data.consul_acl_token.read.encrypted_secret_id
}
```


## Argument Reference

The following arguments are supported:

* `accessor_id` - (Required) The accessor ID of the ACL token.
* `pgp_key` - (Optional) Either a base-64 encoded PGP public key, or a keybase
  username in the form `keybase:some_person_that_exists`. **If you do not set this
  argument, the token secret ID will be written as plain text in the Terraform
  state.**

## Attributes Reference

The following attributes are exported:

* `secret_id` - The secret ID of the ACL token if `pgp_key` has not been set.
* `encrypted_secret_id` - The encrypted secret ID of the ACL token if `pgp_key`
  has been set. You can decrypt the secret by using the command line, for example
  with: `terraform output encrypted_secret | base64 --decode | keybase pgp decrypt`.
