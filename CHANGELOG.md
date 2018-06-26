## 2.1.0 (Unreleased)

NEW FEATURES:

* New resource: `consul_intention`. This provides management of intentions in Consul Connect, used for service authorization.[GH-53]

## 2.0.0 (June 22, 2018)

NOTES:

* The `consul_catalog_entry` resource has been deprecated and will be removed in a future release. Please use the [upgrade guide in the documentation](https://www.terraform.io/docs/providers/consul/upgrading.html#upgrading-to-2-0-0) to migrate to `consul_service` or `consul_node` as appropriate. ([#49](https://github.com/terraform-providers/terraform-provider-consul/issues/49))
* The `consul_agent_service` resource has been deprecated and will be removed in a future release. Please use the [upgrade guide in the documentation](https://www.terraform.io/docs/providers/consul/upgrading.html#upgrading-to-2-0-0) to migrate to `consul_service`. ([#49](https://github.com/terraform-providers/terraform-provider-consul/issues/49))
* The `consul_agent_self` datasource has been deprecated and will be removed in a future release. Please use the [upgrade guide in the documentation](https://www.terraform.io/docs/providers/consul/upgrading.html#upgrading-to-2-0-0) to migrate to `consul_agent_config` if applicable. ([#49](https://github.com/terraform-providers/terraform-provider-consul/issues/49))

IMPROVEMENTS:

* The `consul_service` resource has been modified to use the Consul catalog APIs. The `node` attribute is now required, and nodes that do not exist will not be created automatically. Please see the upgrade guide in the documentation for more detail. ([#49](https://github.com/terraform-providers/terraform-provider-consul/issues/49))
* `consul_catalog_*` data sources have been renamed to remove catalog, for clarity. Both will work going forward, with the catalog version potentially being deprecated on a future date. ([#49](https://github.com/terraform-providers/terraform-provider-consul/issues/49))
* The provider now uses the post-1.0 version of the Consul API. ([#49](https://github.com/terraform-providers/terraform-provider-consul/issues/49))

## 1.1.0 (June 15, 2018)

NOTES:

* This will be the last release prior to significant changes coming to the provider that will deprecate
multiple resources. These changes are primarily to simplify overlap of resource functionality, ensure we are using the correct APIs/design provided by Consul for something like Terraform, and remove resources that can no longer be supported by the current version of the Consul API (1.0+). Read more [here](https://github.com/terraform-providers/terraform-provider-consul/issues/46).

IMPROVEMENTS:

* The provider now allows you to skip HTTPS certificate verification by supplying the `insecure_https` option. ([#31](https://github.com/terraform-providers/terraform-provider-consul/issues/31))

NEW FEATURES:

* New data source: `consul_agent_config`. This new datasource provides information similar to `consul_agent_self`,
but is designed to only expose configuration that Consul will not change without versioning upstream. ([#42](https://github.com/terraform-providers/terraform-provider-consul/issues/42))
* New data source: `consul_key_prefix` corresponds to the existing resource of the same name, allowing config to access a set of keys with a common prefix as a Terraform map for more convenient access ([#34](https://github.com/terraform-providers/terraform-provider-consul/issues/34))

BUG FIXES:

* `consul_catalog` resource now correctly re-creates resources deleted out-of-band. ([#30](https://github.com/terraform-providers/terraform-provider-consul/issues/30))
* `consul_service` resource type now correctly detects when a service has been deleted outside of Terraform, flagging it for re-creation rather than returning an error ([#33](https://github.com/terraform-providers/terraform-provider-consul/issues/33))
* `consul_catalog_service` data source now accepts the `tag` and `datacenter` arguments, as was described in documentation ([#32](https://github.com/terraform-providers/terraform-provider-consul/issues/32))

## 1.0.0 (September 26, 2017)

BUG FIXES:

* d/consul_agent_self: The `enable_ui` config setting was always set to false, regardless of the actual agent configuration ([#16](https://github.com/terraform-providers/terraform-provider-consul/issues/16))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
