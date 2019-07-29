---
layout: "consul"
page_title: "Consul: consul_acl_role"
sidebar_current: "docs-consul-resource-acl-role"
description: |-
  Allows Terraform to create an ACL role
---

# consul_acl_role

Starting with Consul 1.5.0, the consul_acl_role can be used to managed Consul ACL roles.


## Example Usage

```hcl
resource "consul_acl_policy" "read-policy" {
	name = "read-policy"
	rules = "node \"\" { policy = \"read\" }"
	datacenters = [ "dc1" ]
}

resource "consul_acl_role" "test" {
	name = "foo"
	description = "bar"

	policies = [
		"${consul_acl_policy.read-policy.id}"
	]

	service_identities {
		service_name = "foo"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ACL role.
* `description` - (Optional) A free form human readable description of the role.
* `policies` - (Optional) The list of policies that should be applied to the role.
* `service_identities` - (Optional) The list of service identities that should
be applied to the role.

The `service_identities` supports:

* `service_name` - (Required) The name of the service.
* `datacenters` - (Optional) The datacenters the effective policy is valid within.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role.
* `name` - The name of the ACL role.
* `description` - A free form human readable description of the role.
* `policies` - The list of policies that should be applied to the role.
* `service_identities` - The list of service identities that should
be applied to the role.
