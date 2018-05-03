---
layout: "consul"
page_title: "Provider: Consul - Upgrading"
sidebar_current: "docs-consul-upgrading"
description: |-
  Detailed guidelines for upgrading between versions of the Consul Terraform Provider.
---

# Upgrading the Consul Terraform Provider

This page includes details on our compatibility promise and guidelines to
follow when upgrading between versions of the provider. Whenever possible,
we recommend verifying upgrades in isolated test environments.

## Upgrading to 1.1.0

There were several major deprecation notices introduced in 1.1.0.

### Removal of consul_agent_self

The `consul_agent_self` data source will be removed in the next major version
of the provider. As a result, we recommend moving to the new [`consul_agent_config`](/docs/providers/consul/d/agent_config.html) provider.

In the case of information being retrieved from the internal data structures utilized
by the previous resource, Consul still provides this data via API but [promised no
compatibility across versions](https://www.consul.io/docs/upgrade-specific.html#config-section-of-agent-self-endpoint-has-changed),
therefore it is being removed from this provider.