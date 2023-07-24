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
* `display_name` - (Optional) An optional name to use instead of the name
  attribute when displaying information about this auth method.
* `max_token_ttl` - (Optional) The maximum life of any token created by this
  auth method.
* `token_locality` - (Optional) The kind of token that this auth method
  produces. This can be either 'local' or 'global'.
* `description` - (Optional) A free form human readable description of the auth method.
* `config_json` - (Required) The raw configuration for this ACL auth method.
* `config` - (Optional) The raw configuration for this ACL auth method. This
  attribute is deprecated and will be removed in a future version. `config_json`
  should be used instead.
* `namespace` - (Optional, Enterprise Only) The namespace in which to create the auth method.
* `partition` - (Optional, Enterprise Only) The partition the ACL auth method is associated with.
* `namespace_rule` - (Optional, Enterprise Only) A set of rules that control
  which namespace tokens created via this auth method will be created within.

Each `namespace_rule` can have the following attributes:
* `selector` - (Optional) Specifies the expression used to match this namespace
  rule against valid identities returned from an auth method validation.
  Defaults to `""`.
* `bind_namespace` - (Required) If the namespace rule's `selector` matches then
  this is used to control the namespace where the token is created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the the auth method.
* `name` - The name of the ACL auth method.
* `type` - The type of the ACL auth method.
* `display_name` - An optional name to use instead of the name attribute when
  displaying information about this auth method.
* `max_token_ttl` - The maximum life of any token created by this auth method.
* `token_locality` - The kind of token that this auth method produces. This can
  be either 'local' or 'global'.
* `description` - A free form human readable description of the auth method.
* `config_json` - The raw configuration for this ACL auth method.
* `config` - The raw configuration for this ACL auth method. This attribute is
  deprecated and will be removed in a future version. If the configuration is
  too complex to be represented as a map of strings it will be blank.
  `config_json` should be used instead.
* `namespace` - (Enterprise Only) The namespace in which to create the auth method.
* `namespace_rule` - (Enterprise Only) A set of rules that control which
  namespace tokens created via this auth method will be created within.
