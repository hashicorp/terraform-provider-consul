---
layout: "consul"
page_title: "Consul: consul_acl_agent_token"
sidebar_current: "docs-consul-resource-acl-agent-token"
description: |-
  Allows Terraform to create an ACL agent token
---

# consul_acl_token

The `consul_acl_token` resource writes an ACL token into Consul.

## Example Usage

```hcl
resource "consul_acl_policy" "agent" {
  name = "agent"
  rules = "node_prefix \"\" { policy = \"write\" } service_prefix \"\" { policy = \"read\" }"
  datacenters = [ "dc1" ]
}

resource "consul_acl_agent_token" "agent" {
  description = "test"
  policies = ["${consul_acl_policy.agent.name}"]
  local = true
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
