---
layout: "consul"
page_title: "Provider: Consul"
sidebar_current: "docs-consul-index"
description: |-
  Consul is a tool for service discovery, configuration and orchestration. The Consul provider exposes resources used to interact with a Consul cluster. Configuration of the provider is optional, as it provides defaults for all arguments.
---

# Consul Provider

[Consul](https://www.consul.io) is a service networking platform which provides
service discovery, service mesh, and application configuration capabilities.
The Consul provider exposes resources used to interact with a
Consul cluster. Configuration of the provider is optional, as it provides
reasonable defaults for all arguments.

Use the navigation to the left to read about the available resources.

~> **NOTE:** The Consul provider should not be confused with the [Consul remote
state backend][consul-remote-state-backend], which is one of many backends that
can be used to store Terraform state. The Consul provider is instead used to
manage resources within Consul itself, such as adding external services or
working with the key/value store.

[consul-remote-state-backend]: /docs/backends/types/consul.html

## Example Usage

```hcl
# Configure the Consul provider
provider "consul" {
  address    = "demo.consul.io:80"
  datacenter = "nyc1"
}

# Access a key in Consul
data "consul_keys" "app" {
  key {
    name    = "ami"
    path    = "service/app/launch_ami"
    default = "ami-1234"
  }
}

# Use our variable from Consul
resource "aws_instance" "app" {
  ami = data.consul_keys.app.var.ami
}
```

## Compatibility

The Consul Terraform provider uses features of the latest version of Consul.
Some resources may not be supported by older versions of Consul.

The known compatibility between this provider and Consul is:

| Terraform provider version | Consul version |
| -------------------------- | -------------- |
| 2.15.0                     | >= 1.11.0      |
| 2.14.0                     | >= 1.10.0      |
| 2.13.0                     | >= 1.10.0      |

## Argument Reference

The following arguments are supported:

* `address` - (Optional) The HTTP(S) API address of the agent to use. Defaults to "127.0.0.1:8500".
* `scheme` - (Optional) The URL scheme of the agent to use ("http" or "https"). Defaults to "http".
* `http_auth` - (Optional) HTTP Basic Authentication credentials to be used when communicating with Consul, in the format of either `user` or `user:pass`. This may also be specified using the `CONSUL_HTTP_AUTH` environment variable.
* `datacenter` - (Optional) The datacenter to use. Defaults to that of the agent.
* `token` - (Optional) The ACL token to use by default when making requests to the agent. Can also be specified with `CONSUL_HTTP_TOKEN` or `CONSUL_TOKEN` as an environment variable.
* `ca_file` - (Optional) A path to a PEM-encoded certificate authority used to verify the remote agent's certificate.
* `ca_pem` - (Optional) PEM-encoded certificate authority used to verify the remote agent's certificate.
* `cert_file` - (Optional) A path to a PEM-encoded certificate provided to the remote agent; requires use of `key_file` or `key_pem`.
* `cert_pem` - (Optional) PEM-encoded certificate provided to the remote agent; requires use of `key_file` or `key_pem`.
* `key_file`- (Optional) A path to a PEM-encoded private key, required if `cert_file` or `cert_pem` is specified.
* `key_pem`- (Optional) PEM-encoded private key, required if `cert_file` or `cert_pem` is specified.
* `ca_path` - (Optional) A path to a directory of PEM-encoded certificate authority files to use to check the authenticity of client and server connections. Can also be specified with the `CONSUL_CAPATH` environment variable.
* `insecure_https`- (Optional) Boolean value to disable SSL certificate verification; setting this value to true is not recommended for production use. Only use this with scheme set to "https".
* `header` - (Optional) A configuration block, described below, that provides additional headers to be sent along with all requests to the Consul server. This block can be specified multiple times.

The `header` configuration block accepts the following arguments:

 * `name` - (Required) The name of the header.

 * `value` - (Required) The value of the header.

## Environment Variables

All environment variables listed in the [Consul environment variables](https://www.consul.io/docs/commands/index.html#environment-variables)
documentation are supported by the Terraform provider.
