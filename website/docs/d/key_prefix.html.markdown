---
layout: "consul"
page_title: "Consul: consul_key_prefix"
sidebar_current: "docs-consul-data-source-key-prefix"
description: |-
  Reads values from a "namespace" of Consul keys that share a
  common name prefix.
---

# consul_key_prefix

Allows Terraform to read values from a "namespace" of Consul keys that
share a common name prefix.

## Example Usage

```hcl
data "consul_key_prefix" "app" {
  datacenter = "nyc1"
  token      = "abcd"

  # Prefix to add to prepend to all of the subkey names below.
  path_prefix = "myapp/config/"

  # Read the ami subkey
  subkey {
    name    = "ami"
    path    = "app/launch_ami"
    default = "ami-1234"
  }
}

# Start our instance with the dynamic ami value
resource "aws_instance" "app" {
  ami = "${data.consul_key_prefix.app.var.ami}"

  # ...
}
```

```hcl
data "consul_key_prefix" "web" {
  datacenter = "nyc1"
  token      = "efgh"

  # Prefix to add to prepend to all of the subkey names below.
  path_prefix = "myapp/config/"
}

# Start our instance with the dynamic ami value
resource "aws_instance" "web" {
  ami = "${data.consul_key_prefix.web["app/launch_ami"]}"

  # ...
}
```

## Argument Reference

The following arguments are supported:

* `datacenter` - (Optional) The datacenter to use. This overrides the
  datacenter in the provider setup and the agent's default datacenter.

* `token` - (Optional) The ACL token to use. This overrides the
  token that the agent provides by default.

* `path_prefix` - (Required) Specifies the common prefix shared by all keys
  that will be read by this data source instance. In most cases, this will
  end with a slash to read a "folder" of subkeys.

* `subkey` - (Optional) Specifies a subkey in Consul to be read. Supported
  values documented below. Multiple blocks supported.

The `subkey` block supports the following:

* `name` - (Required) This is the name of the key. This value of the
  key is exposed as `var.<name>`. This is not the path of the subkey
  in Consul.

* `path` - (Required) This is the subkey path in Consul (which will be appended
  to the given `path_prefix`) to construct the full key that will be used
  to read the value.

* `default` - (Optional) This is the default value to set for `var.<name>`
  if the key does not exist in Consul. Defaults to an empty string.

## Attributes Reference

The following attributes are exported:

* `datacenter` - The datacenter the keys are being read from.
* `path_prefix` - the common prefix shared by all keys being read.
* `var.<name>` - For each name given, the corresponding attribute
  has the value of the key.
* `subkeys` - A map of the subkeys and values is set if no `subkey`
  block is provided.
