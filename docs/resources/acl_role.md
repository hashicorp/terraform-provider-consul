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

resource "consul_acl_role" "read" {
	name = "foo"
	description = "bar"

	policies = [
		"${consul_acl_policy.read-policy.id}"
	]

	service_identities {
		service_name = "foo"
	}

	templated_policies {
		template_name = "builtin/dns"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ACL role.
* `description` - (Optional) A free form human readable description of the role.
* `policies` - (Optional) The list of policies that should be applied to the role.
* `service_identities` - (Optional) The list of service identities that should be applied to the role.
* `node_identities` - (Optional) The list of node identities that should be applied to the role.
* `templated_policies` - (Optional) The list of templated policies that should be applied to the token.
* `namespace` - (Optional, Enterprise Only) The namespace to create the role within.
* `partition` - (Optional, Enterprise Only) The partition the ACL role is associated with.

The `service_identities` block supports:

* `service_name` - (Required) The name of the service.
* `datacenters` - (Optional) The datacenters the effective policy is valid within.

The `node_identities` block supports:

* `node_name` - (Required) The name of the node.
* `datacenter` - (Required) The datacenter of the node.

The `templated_policies` block supports the following arguments:

* `template_name` - (Optional) The name of the templated policy.
* `template_variables` - (Optional) The list of the templated policy variables.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role.
* `name` - The name of the ACL role.
* `description` - A free form human readable description of the role.
* `policies` - The list of policies that should be applied to the role.
* `service_identities` - The list of service identities that should be applied to the role.
* `node_identities` - The list of node identities that should be applied to the role.
* `templated_policies` - The list of templated policies that should be applied to the token.
* `namespace` - The namespace to create the role within.


## Import

`consul_acl_role` can be imported:

```
$ terraform import consul_acl_role.read 816a195f-6cb1-2e8d-92af-3011ae706318
```
