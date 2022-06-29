---
layout: "consul"
page_title: "Upgrading"
sidebar_current: "docs-consul-upgrading"
description: |-
  Detailed guidelines for upgrading between versions of the Consul Terraform Provider.
---

# Upgrading the Consul Terraform Provider

This page includes details on our compatibility promise and guidelines to
follow when upgrading between versions of the provider. Whenever possible,
we recommend verifying upgrades in isolated test environments.

## Upgrading to 2.11.0

The `consul_config_entry` resource now supports the `service-intentions` kind to
define configuration entries. The `consul_intention` resource used to manage
Consul intentions until Consul 1.9 have now been deprecated. It is recommended
to migrate `consul_intention` to `consul_config_entry`:

```terraform
# This was correct in version 2.10.0
resource "consul_intention" "database" {
  source_name      = "api"
  destination_name = "db"
  action           = "allow"
}

# This is now the correct configuration starting version 2.11.0
resource "consul_config_entry" "database" {
  name = "db"
  kind = "service-intentions"

  config_json = jsonencode({
    Sources = [{
      Action     = "allow"
      Name       = "api"
      Precedence = 9
      Type       = "consul"
    }]
  })
}
```

## Upgrading to 2.4.0

### Changes to consul_service

The `external` attribute introduced in 2.3.0 has been deprecated and no longer
updates the associated Consul Node meta information. The same functionality
can be achieved by setting the meta attribute on the `consul_node` resource:

``` terraform
# (Deprecated) External service configuration in 2.3.0
resource "consul_node" "compute" {
  name    = "compute-example"
  address = "www.hashicorptest.com"
}

resource "consul_service" "example1" {
  name = "example"
  node = "${consul_node.compute.name}"
  port = 80

  external = true
}

# Updated configuration as of 2.4.0
resource "consul_node" "compute" {
  name    = "compute-example"
  address = "www.hashicorptest.com"

  meta = {
    "external-node" = "true"
    "external-probe" = "true"
  }
}

resource "consul_service" "example1" {
  name = "example"
  node = "${consul_node.compute.name}"
  port = 80
}
```

## Upgrading to 2.0.0

There were several major deprecation notices introduced in 2.0.0. This
reviews the details of each and provides migration instructions to the
appropriate resources.

### Deprecation of consul_agent_self

The `consul_agent_self` data source will be removed in the next major version
of the provider. As a result, we recommend moving to the new [`consul_agent_config`](/docs/providers/consul/d/agent_config.html) data source.

The `consul_agent_config` resource returns far less attributes, so as a result
it may not provide all necessary functionality. Consul does still provide this data via API but [promises no
compatibility across versions](https://www.consul.io/docs/upgrade-specific.html#config-section-of-agent-self-endpoint-has-changed),
therefore it is being removed from this provider.

### Deprecation of consul_agent_service

The `consul_agent_service` resource will be removed in the next major version
of the provider. As a result, we recommend moving to the [`consul_service`](/docs/providers/consul/r/service.html) resource.

This resource has been updated to use the correct catalog APIs in place
of service registration APIs. The `consul_agent_service` resource previously also
used the service registration API designed for registration against an agent
running on a local node. Because Terraform is intended to be run externally to
the cluster, and for other internal reasons, this API was the incorrect one to use.

View migration instructions [here][migrate_service].

### Migrating to consul_service or consul_node resources

Migration to the `consul_service` resources are possible in two ways. Both require
the configuration to be modified.

**From `consul_agent_service` to `consul_service`:**

1. Rename `consul_agent_service` resources to `consul_service` in the Terraform configuration files.
1. Add the `node` attribute where the service is currently registered, retrievable
by [querying the catalog](https://www.consul.io/api/catalog.html#list-nodes-for-service) or using the UI. This new attribute is required.
1. For a small number of resources, the first class [`state rm`](https://www.terraform.io/docs/commands/state/rm.html) and [`import`](https://www.terraform.io/docs/import/usage.html) commands can
be used to first remove the old resource from the state, and then import it under the new resource
name.
1. For a large number of resources, edit the state file directly to rename every resource at
the same time (replace all instances of `consul_agent_service` with `consul_service`). This
requires understanding the consequences and guidelines for [editing state files](https://www.terraform.io/docs/backends/state.html#manual-state-pull-push),
so please read those.

After following these steps,  `terraform plan` should show no changes.

**From `consul_catalog_entry` to `consul_service` or `consul_node`:**

1. Copy the attributes from the `service {}` or `node {}` blocks into
new `consul_service` or `consul_node` resources in the Terraform
configuration files.
1. For a small number of resources, the first class [`state rm`](https://www.terraform.io/docs/commands/state/rm.html) and [`import`](https://www.terraform.io/docs/import/usage.html) commands can
be used to first remove the old resource from the state, and then import it under the new resource
name. The `node` attribute will need to be added to services.

### Modifications to consul_service

The `consul_service` resource has been modified to use catalog APIs in place
of service registration APIs for creating services in the Consul catalog. This should
be a functionally compatible change, and create or read services as prior. It now replaces `consul_catalog_entry` (the `service {}` block) and `consul_agent_service`. The `node` attribute is now required and the node must exist.

### Deprecation of consul_catalog_entry

The `consul_catalog_entry` resource will be removed in the next major version
of the provider. As a result, we recommend moving to the [`consul_service`](/docs/providers/consul/r/service.html) or [`consul_node`](/docs/providers/consul/r/node.html) resources.

These resources have been updated (or created) to use the correct catalog APIs as with `consul_catalog_entry`, but provide a first-class resource name.

View migration instructions [here][migrate_service].

### Renaming of Catalog Data Sources

`consul_catalog_nodes`, `consul_catalog_services`, and `consul_catalog_service` have been renamed to
`consul_nodes`, `consul_services`, and `consul_service` respectively. The prior naming will
continue to work, but in the long term it may be deprecated and removed. This is to present
a more consistent and intuitive naming convention for the resources.

[migrate_service]: /docs/providers/consul/guides/upgrading.html#migrating-to-consul_service-or-consul_node-resources "Migrate to consul_service or consul_node"
