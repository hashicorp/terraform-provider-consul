---
layout: "consul"
page_title: "Consul: consul_license"
sidebar_current: "docs-consul-license"
description: |-
  Manage the Consul Enterprise license.
---

# consul_license

~> **NOTE:** This feature requires Consul Enterprise.

The `consul_license` resource provides datacenter-level management of
the Consul Enterprise license. If ACLs are enabled then a token with operator
privileges may be required in order to use this command.

## Example Usage

```hcl
resource "consul_license" "license" {
  license = file("license.hclic")
}
```

## Argument Reference

The following arguments are supported:

* `datacenter` - (Optional) The datacenter to use. This overrides the
  agent's default datacenter and the datacenter in the provider setup.
* `license` - (Required) The Consul license to use.

## Attributes Reference

The following attributes are exported:

* `valid` - Whether the license is valid.
* `license_id` - The ID of the license used.
* `customer_id` - The ID of the customer the license is attached to.
* `installation_id` - The ID of the current installation.
* `issue_time` - The date the license was issued.
* `start_time` - The start time of the license.
* `expiration_time` - The expiration time of the license.
* `product` - The product for which the license is valid.
* `flags` - The metadata attached to the license.
* `features` - The features for which the license is valid.
* `warnings` - A list of warning messages regarding the license validity.
