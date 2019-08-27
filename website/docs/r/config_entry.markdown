---
layout: "consul"
page_title: "Consul: config_entry"
sidebar_current: "docs-consul-resource-config-entry"
description: |-
  Registers a configuration entry in Consul.
---

# consul_config_entry

The [Configuration Entry](https://www.consul.io/docs/agent/config_entries.html)
resource can be used to provide cluster-wide defaults for various aspects of
Consul.

## Example Usage

```hcl
resource "consul_config_entry" "service-defaults" {
	name = "foo"
	kind = "service-defaults"

	protocol = "https"
}
```

```hcl
resource "consul_config_entry" "proxy-defaults" {
	name = "global"
	kind = "proxy-defaults"

	config = {
        local_connect_timeout_ms = 1000
        handshake_timeout_ms = 1000
	}
}
```

## Argument Reference

The following arguments are supported:

* `kind` - (Required) The kind of configuration entry to register. Can be
`proxy-defaults` or `service-defaults`.

* `name` - (Required) The name of the configuration entry being registred. If
`kind` is `proxy-defaults`, `name` must be `global`.

* `config` - (Optional) An arbitrary map of configuration values used by Connect
proxies. Can only be set when `kind` is `proxy-defaults`.

* `protocol` - (Optional) The protocol of the service. Can only be set when
`kind` is `service-defaults`.

* `token` - (Optional) ACL token.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the configuration entry.

* `kind` - The kind of the configuration entry, `proxy-defaults` or
`service-defaults`.

* `name` - The name of the configuration entry.

* `config` - A map of configuration values.

* `protocol` - The protocol of the service.

* `token` - ACL token.
