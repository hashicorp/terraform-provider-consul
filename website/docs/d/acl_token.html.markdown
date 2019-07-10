---
layout: "consul"
page_title: "Consul: consul_acl_token"
sidebar_current: "docs-consul-data-source-acl-token"
description: |-
  Provides the ACL Token information, in particular the secret_id.
---

# consul_acl_token

The `consul_acl_token` data source returns the information related to the `consul_acl_token` resource and in particular the related `secret_id` attribute.

## Example Usage

```hcl
data "consul_acl_token" "test" {}

output "consul_acl_secret" {
  value = "${data.consul_acl_token.test.secret_id}"
}
```


## Argument Reference

The following arguments are supported:

* `accessor_id` - (required) The accessor_id of the generated token.

## Attributes Reference

The following attributes are exported:

* `secret_id` - The secret_id of the acl token
* `description` - The decsription of the acl token
* `policies` - The policies of the acl token
* `local` - the local Boolean attribute of the acl token