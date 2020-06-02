---
layout: "consul"
page_title: "Consul: consul_agent_token"
sidebar_current: "docs-consul-resource-agent-token"
description: |-
  Provides access to Consul Agent Tokens. This can be used to set the [acl_tokens](https://www.consul.io/docs/agent/options#acl_tokens).
---

# consul_agent_service

Provides access to Consul Agent Tokens. This can be used to set the [acl_tokens](https://www.consul.io/docs/agent/options#acl_tokens).

!> To use this resource, Terraform must have access to the Consul HTTP API.
!> Read and delete are not implemented. Changes outside of terraform are therefore not detected.

## Example Usage

Creates a token and sets it to `node-1`:

```hcl
resource "consul_acl_token" "node_1" {
  description = "node-1 agent token"
  policies    = ...
}

resource "consul_agent_token" "node_1" {
  address     = "node-1:8500"
  type        = "agent"
  accessor_id = consul_acl_token.node_1.accessor_id
}

```

Creates individual agent tokens for all nodes and sets them:

```hcl
provider "consul" {
  address = "localhost:8500"
}

data "consul_nodes" "all" {}

output "test" {
  value = data.consul_nodes.all
}

resource "consul_acl_policy" "agent_policies" {
  count = length(data.consul_nodes.all.nodes)

  name  = "node-policy-${data.consul_nodes.all.nodes[count.index].name}"
  rules = <<-RULE
    node "${data.consul_nodes.all.nodes[count.index].name}" {
      policy = "write"
    }
  RULE
}

resource "consul_acl_token" "agent_tokens" {
  count = length(data.consul_nodes.all.nodes)

  description = "node token for ${data.consul_nodes.all.nodes[count.index].name}"
  policies    = [consul_acl_policy.agent_policies[count.index].name]
}

resource "consul_agent_token" "agent_tokens" {
  count = length(data.consul_nodes.all.nodes)

  address     = "${data.consul_nodes.all.nodes[count.index].address}:8500"
  type        = "agent"
  accessor_id = consul_acl_token.agent_tokens[count.index].accessor_id
}

```

## Argument Reference

The following arguments are supported:

* `address` - (Required) The address of the service.

* `type` - (Required) The type of the token. (default|agent|master|replication)

* `accessor_id` - (Optional) The token accessor ID.

## Attributes Reference

The following attributes are exported:

* `id` - An ID generated based on the `address` and `type`

* `address` - The address of the agent.

* `type` - The type of the token.

* `accessor_id` - The token accessor ID.
