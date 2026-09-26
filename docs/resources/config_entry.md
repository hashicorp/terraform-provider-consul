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

~> **NOTE:** Because the schema in a `consul_config_entry` resource can vary
widely between the various configuration entry kinds, it is necessary to explicitly
define every attribute to avoid Terraform reporting a diff on the resource.

## Example Usage

```hcl
resource "consul_config_entry" "proxy_defaults" {
  kind = "proxy-defaults"
  # Note that only "global" is currently supported for proxy-defaults and that
  # Consul will override this attribute if you set it to anything else.
  name = "global"

  config_json = jsonencode({
    Config = {
      local_connect_timeout_ms = 1000
      handshake_timeout_ms     = 10000
    }
  })
}

resource "consul_config_entry" "web" {
  name = "web"
  kind = "service-defaults"

  config_json = jsonencode({
    Protocol    = "http"
  })
}

resource "consul_config_entry" "admin" {
  name = "admin"
  kind = "service-defaults"

  config_json = jsonencode({
    Protocol    = "http"
  })
}

resource "consul_config_entry" "service_resolver" {
  kind = "service-resolver"
  name = consul_config_entry.web.name

  config_json = jsonencode({
    DefaultSubset = "v1"

    Subsets = {
      "v1" = {
        Filter = "Service.Meta.version == v1"
      }
      "v2" = {
        Filter = "Service.Meta.version == v2"
      }
    }
  })
}

resource "consul_config_entry" "service_splitter" {
  kind = "service-splitter"
  name = consul_config_entry.service_resolver.name

  config_json = jsonencode({
    Splits = [
      {
        Weight        = 90
        ServiceSubset = "v1"
      },
      {
        Weight        = 10
        ServiceSubset = "v2"
      },
    ]
  })
}

resource "consul_config_entry" "service_router" {
  kind = "service-router"
  name = "web"

  config_json = jsonencode({
    Routes = [
      {
        Match = {
          HTTP = {
            PathPrefix = "/admin"
          }
        }

        Destination = {
          Service = "admin"
        }
      },
      # NOTE: a default catch-all will send unmatched traffic to "web"
    ]
  })
}

resource "consul_config_entry" "ingress_gateway" {
  name = "us-east-ingress"
  kind = "ingress-gateway"

  config_json = jsonencode({
    TLS = {
      Enabled = true
    }
    Listeners = [{
      Port     = 8000
      Protocol = "http"
      Services = [{ Name  = "*" }]
    }]
  })
}

resource "consul_config_entry" "terminating_gateway" {
  name = "us-west-gateway"
  kind = "terminating-gateway"

  config_json = jsonencode({
    Services = [{ Name = "billing" }]
  })
}
```

### `service-intentions` config entry

```hcl
resource "consul_config_entry" "service_intentions" {
  name = "api-service"
  kind = "service-intentions"

  config_json = jsonencode({
    Sources = [
      {
        Action     = "allow"
        Name       = "frontend-webapp"
        Precedence = 9
        Type       = "consul"
      },
      {
        Action     = "allow"
        Name       = "nightly-cronjob"
        Precedence = 9
        Type       = "consul"
      }
    ]
  })
}
```

```hcl
resource "consul_config_entry" "sd" {
  name = "fort-knox"
  kind = "service-defaults"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "jwt_provider" {
	name = "test-provider"
	kind = "jwt-provider"

	config_json = jsonencode({
		Issuer = "test-issuer"
		JSONWebKeySet = {
			Remote = {
				URI = "https://127.0.0.1:9091"
				FetchAsynchronously = true
			}
		}
		Forwarding = {
			HeaderName = "test-token"
		}
	})
}

resource "consul_config_entry" "service_intentions" {
  name = consul_config_entry.sd.name
  kind = "service-intentions"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "contractor-webapp"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              Methods   = ["GET", "HEAD"]
              PathExact = "/healtz"
            }
            JWT = {
              Providers = [
                {
                  Name = consul_config_entry.jwt_provider.name
                }
              ]
            }
          }
        ]
        Precedence = 9
        Type       = "consul"
      },
      {
        Name        = "admin-dashboard-webapp",
        Permissions = [
          {
            Action = "deny",
            HTTP = {
              PathPrefix= "/debugz"
            }
          },
          {
            Action= "allow"
            HTTP = {
              PathPrefix= "/"
            }
          }
        ],
        Precedence = 9
        Type       = "consul"
      }
    ]
  })
}
```

### `exported-services` config entry

```hcl
resource "consul_config_entry" "exported_services" {
	name = "test"
	kind = "exported-services"

	config_json = jsonencode({
		Services = [{
			Name = "test"
			Namespace = "default"
			Consumers = [{
				Partition = "default"
			}]
		}]
	})
}
```

### `mesh` config entry

```hcl
resource "consul_config_entry" "mesh" {
	#name      = "mesh"  ### The "name" argument field doesn't align with the consul configuration entry object which may cause throw error when customize to any other name instead of "mesh". Ref. https://developer.hashicorp.com/consul/docs/connect/config-entries/mesh#available-fields
	kind      = "mesh"
  partition = "default"

	config_json = jsonencode({
		TransparentProxy = {
			MeshDestinationsOnly = true
		}
	})
}
```

### `jwt-provider` config entry

```hcl
resource "consul_config_entry" "jwt_provider" {
  name = "provider-name"
  kind = "jwt-provider"

  config_json = jsonencode({
    Issuer = "https://your.issuer.com"
    JSONWebKeySet = {
      Remote = {
        URI = "https://your-remote.jwks.com"
        FetchAsynchronously = true
        CacheDuration = "10s"
      }
    }
    Forwarding = {
      HeaderName = "test-token"
    }
  })
}
```


## Argument Reference

The following arguments are supported:

* `kind` - (Required) The kind of configuration entry to register.

* `name` - (Required) The name of the configuration entry being registered.

* `partition` - (Optional, Enterprise Only) The partition the config entry is associated with.

* `namespace` - (Optional, Enterprise Only) The namespace to create the config entry within.

* `config_json` - (Optional) An arbitrary map of configuration values.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the configuration entry.

* `kind` - The kind of the configuration entry.

* `name` - The name of the configuration entry.

* `partition` - The partition the config entry is associated with.

* `namespace` - The namespace to create the config entry within.

* `config_json` - A map of configuration values.


## Import

`consul_config_entry` can be imported using the syntax `<kind>/<name>` if the
config entry is in the default partition and default namespace, or
`<partition>/<namespace>/<kind>/<name>` for config entries in a non-default
partition or namespace:

```
$ terraform import consul_config_entry.service_splitter 816a195f-6cb1-2e8d-92af-3011ae706318
```
