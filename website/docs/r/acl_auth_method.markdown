---
layout: "consul"
page_title: "Consul: consul_acl_auth_method"
sidebar_current: "docs-consul-resource-acl-auth-method"
description: |-
  Allows Terraform to create an ACL auth method
---

# consul_acl_auth_method

Starting with Consul 1.5.0, the consul_acl_auth_method resource can be used to
managed [Consul ACL auth methods](https://www.consul.io/docs/acl/auth-methods).


## Example Usage

Define a `kubernetes` auth method:
```hcl
resource "consul_acl_auth_method" "minikube" {
  name        = "minikube"
  type        = "kubernetes"
  description = "dev minikube cluster"

  config_json = jsonencode({
    Host              = "https://192.0.2.42:8443"
    CACert            = "-----BEGIN CERTIFICATE-----\n...-----END CERTIFICATE-----\n"
    ServiceAccountJWT = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9..."
  })
}
```

Define a `jwt` auth method:
```hcl
resource "consul_acl_auth_method" "minikube" {
  name        = "auth_method"
  type        = "jwt"

  config_json = jsonencode({
    JWKSURL          = "https://example.com/identity/oidc/.well-known/keys"
    JWTSupportedAlgs = "RS256"
    BoundIssuer      = "https://example.com"
    ClaimMappings    = {
      subject = "subject"
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ACL auth method.
* `type` - (Required) The type of the ACL auth method.
* `description` - (Optional) A free form human readable description of the auth method.
* `config_json` - (Required) The raw configuration for this ACL auth method.
* `config` - (Optional) The raw configuration for this ACL auth method. This
  attribute is deprecated and will be removed in a future version. `config_json`
  should be used instead.
* `namespace` - (Optional, Enterprise Only) The namespace to create the policy within.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the the auth method.
* `name` - The name of the ACL auth method.
* `type` - The type of the ACL auth method.
* `description` - A free form human readable description of the auth method.
* `config_json` - The raw configuration for this ACL auth method.
* `config` - The raw configuration for this ACL auth method. This attribute is
  deprecated and will be removed in a future version. If the configuration is
  too complex to be represented as a map of strings it will be blank.
  `config_json` should be used instead.
