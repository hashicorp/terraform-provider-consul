---
layout: "consul"
page_title: "Consul: consul_acl_token"
sidebar_current: "docs-consul-data-source-acl-token"
description: |-
  Provides information about a Consul ACL Token.
---

# consul_acl_token

The `consul_acl_token` data source returns the information related to the
`consul_acl_token` resource with the exception of its secret ID.

If you want to get the secret ID associated with a token, use the
[`consul_acl_token_secret_id` data source](/docs/providers/consul/d/acl_token_secret_id.html).

## Example Usage

```hcl
data "consul_acl_token" "test" {
  accessor_id = "00000000-0000-0000-0000-000000000002"
}

output "consul_acl_policies" {
  value = "${data.consul_acl_token.test.policies}"
}
```


## Argument Reference

The following arguments are supported:

* `accessor_id` - (Required) The accessor ID of the ACL token.
* `namespace` - (Optional, Enterprise Only) The namespace to lookup the ACL token.

## Attributes Reference

The following attributes are exported:

* `description` - The description of the ACL token.
* `policies` - A list of policies associated with the ACL token. Each entry has an `id` and a `name` attribute.
* `roles` - The list of roles attached to the token.
* `service_identities` - The list of service identities attached to the token. Each entry has a `service_name` and a `datacenters` attribute.
* `node_identities` - The list of node identities attached to the token. Each entry has a `node_name` and a `datacenter` attributes.
* `local` - Whether the ACL token is local to the datacenter it was created within.
* `expiration_time` - If set this represents the point after which a token should be considered revoked and is eligible for destruction.
