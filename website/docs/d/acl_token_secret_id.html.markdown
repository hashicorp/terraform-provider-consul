---
layout: "consul"
page_title: "Consul: consul_acl_token_secret_id"
sidebar_current: "docs-consul-data-source-acl-token-secret-id"
description: |-
  Provides the ACL Token secret ID.
---

# consul_acl_token_secret_id

The `consul_acl_token_secret` data source returns the secret ID associated to
the accessor ID. This can be useful to make systems that cannot use an auth
method to interface with Consul.

If you want to get other attributes of the Consul ACL token, please use the
`consul_acl_token` data source.

## Example Usage

```hcl
resource "consul_acl_policy" "test" {
	name = "test"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_token" "test" {
	description = "test"
	policies = ["${consul_acl_policy.test.name}"]
	local = true
}

data "consul_acl_token_secret_id" "read" {
    accessor_id = "${consul_acl_token.test.id}"
}

output "consul_acl_token_secret_id" {
  value = "${data.consul_acl_token.read.secret_id}"
}
```


## Argument Reference

The following arguments are supported:

* `accessor_id` - (Required) The accessor ID of the ACL token.

## Attributes Reference

The following attributes are exported:

* `secret_id` - The secret ID of the ACL token.
