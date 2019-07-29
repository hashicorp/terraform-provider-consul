---
layout: "consul"
page_title: "Consul: consul_acl_auth_method"
sidebar_current: "docs-consul-resource-acl-auth-method"
description: |-
  Allows Terraform to create an ACL auth method
---

# consul_acl_auth_method

Starting with Consul 1.5.0, the consul_acl_auth_method resource can be used to
managed Consul ACL auth methods.


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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ACL auth method.
* `type` - (Required) The type of the ACL auth method.
* `description` - (Optional) A free form human readable description of the auth method.
* `config` - (Required) The raw configuration for this ACL auth method.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the the auth method.
* `name` - The name of the ACL auth method.
* `type` - The type of the ACL auth method.
* `description` - A free form human readable description of the auth method.
* `config` - The raw configuration for this ACL auth method.
