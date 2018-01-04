## 1.1.0 (Unreleased)

NEW FEATURES:

* New data source: `consul_key_prefix` corresponds to the existing resource of the same name, allowing config to access a set of keys with a common prefix as a Terraform map for more convenient access [GH-34]

BUG FIXES:

* `consul_service` resource type now correctly detects when a service has been deleted outside of Terraform, flagging it for re-creation rather than returning an error [GH-33]
* `consul_catalog_service` data source now accepts the `tag` and `datacenter` arguments, as was described in documentation [GH-32]

## 1.0.0 (September 26, 2017)

BUG FIXES:

* d/consul_agent_self: The `enable_ui` config setting was always set to false, regardless of the actual agent configuration ([#16](https://github.com/terraform-providers/terraform-provider-consul/issues/16))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
