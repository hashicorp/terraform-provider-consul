---
layout: "consul"
page_title: "Consul: consul_certificate_authority"
sidebar_current: "docs-consul-resource-certificate-authority"
description: |-
    A resource that manage the Consul Connect Certificate Authority
---

# certificate_authority

The `consul_certificate_authority` resource can be used to manage the configuration of
the Certificate Authority used by [Consul Connect](https://www.consul.io/docs/connect/ca).

## Example Usage

Use the built-in CA with specific TTL:

```hcl
resource "consul_certificate_authority" "connect" {
  connect_provider = "consul"

  config = {
    LeafCertTTL         = "24h"
    RotationPeriod      = "2160h"
    IntermediateCertTTL = "8760h"
  }
}
```

Use Vault to manage and sign certificates:

```hcl
resource "consul_certificate_authority" "connect" {
  connect_provider = "vault"

  config = {
    address = "http://localhost:8200"
    token = "..."
    root_pki_path = "connect-root"
    intermediate_pki_path = "connect-intermediate"
  }
}
```

Use the [AWS Certificate Manager Private Certificate Authority](https://aws.amazon.com/certificate-manager/private-certificate-authority/):

```hcl
resource "consul_certificate_authority" "connect" {
  connect_provider = "aws-pca"

  config = {
    existing_arn = "arn:aws:acm-pca:region:account:certificate-authority/12345678-1234-1234-123456789012"
  }
}
```

## Argument Reference

The following arguments are supported:

* `connect_provider` - (Required, string) Specifies the CA provider type to use.
* `config` - (Required, map) The raw configuration to use for the chosen provider.


## Attributes Reference

The following attributes are exported:

* `connect_provider` - Specifies the CA provider type to use.
* `config` - The raw configuration to use for the chosen provider.

## Import

`certificate_authority` can be imported:

```
$ terraform import certificate_authority.connect connect-ca
```
