---
layout: "consul"
page_title: "Consul: consul_acl_token"
sidebar_current: "docs-consul-data-source-acl-token"
description: |-
  Provides the ACL Token information.
---

# consul_acl_token

The `consul_acl_token` data source returns the information related to the
`consul_acl_token` resource with the exception of its secret ID.

If you want to get the secret ID associated with a token, use the
`consul_acl_token_secret_id` data source.

## Example Usage

```hcl
data "consul_acl_token" "test" {
  accessor_id = "..."
}

output "consul_acl_policies" {
  value = "${data.consul_acl_token.test.policies}"
}
```


## Argument Reference

The following arguments are supported:

* `accessor_id` - (Required) The accessor ID of the ACL token.

## Attributes Reference

The following attributes are exported:

* `description` - The description of the ACL token.
* `policies` - The policies associated with the ACL token.
* `local` - Whether the ACL token is local to the datacenter it was created within.
