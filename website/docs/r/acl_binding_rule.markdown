---
layout: "consul"
page_title: "Consul: consul_acl_binding_rule"
sidebar_current: "docs-consul-resource-acl-binding-rule"
description: |-
  Allows Terraform to create an ACL binding rule
---

# consul_acl_binding_rule

Starting with Consul 1.5.0, the consul_acl_binding_rule resource can be used to
managed Consul ACL binding rules.


## Example Usage

```hcl
resource "consul_acl_auth_method" "minikube" {
	name        = "minikube"
    type        = "kubernetes"
    description = "dev minikube cluster"

	config = {
        Host = "https://192.0.2.42:8443"
		CACert = "-----BEGIN CERTIFICATE-----\n...-----END CERTIFICATE-----\n"
        ServiceAccountJWT = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9..."
    }
}

resource "consul_acl_binding_rule" "test" {
	auth_method = "${consul_acl_auth_method.minikube.name}"
	description = "foobar"
	selector    = "serviceaccount.namespace==default"
	bind_type   = "service"
	bind_name   = "minikube"
}
```

## Argument Reference

The following arguments are supported:

* `auth_method` - (Required) The name of the ACL auth method this rule apply.
* `description` - (Optional) A free form human readable description of the
binding rule.
* `selector` - (Optional) The expression used to math this rule against valid
identities returned from an auth method validation.
* `bind_type` - (Required) Specifies the way the binding rule affects a token
created at login.
* `bind_name` - (Required) The name to bind to a token at login-time.
* `namespace` - (Optional, Enterprise Only) The namespace to create the binding
rule within.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the the binding rule.
* `auth_method` - The name of the ACL auth method this rule apply.
* `description` - A free form human readable description of the
binding rule.
* `selector` - The expression used to math this rule against valid
identities returned from an auth method validation.
* `bind_type` - Specifies the way the binding rule affects a token
created at login.
* `bind_name` - The name to bind to a token at login-time.
