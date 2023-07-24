## 2.18.0 (July 24, 2023)

NEW FEATURES

* The provider can now use [Terraform Cloud Workload Identity](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials/workload-identity-tokens) to connect to Consul using the `use_terraform_cloud_workload_identity` argument in the provider configuration [[#342](https://github.com/hashicorp/terraform-provider-consul/issues/342)].
* The `consul_config_entry` resource now supports the `jwt-provider` configuration entries [[#349](https://github.com/hashicorp/terraform-provider-consul/issues/349)].
* The `consul_prepared_query` resource can now specify `targets` in their configuration [[#340](https://github.com/hashicorp/terraform-provider-consul/issues/340)].
* The `consul_prepared_query` resource can now specify `remove_empty_tags` in their configuration [[#312](https://github.com/hashicorp/terraform-provider-consul/issues/312)].
* The `consul_certificate_authority` now supports the `config_json` argument to set complex configuration [[#341](https://github.com/hashicorp/terraform-provider-consul/issues/341)].

CHANGES:

* The `name` argument of the `consul_prepared_query` resource is now optional [[#312](https://github.com/hashicorp/terraform-provider-consul/issues/312)].
* The `consul_acl_role` now return an explicit error message if a policy is specified using its name instead of its ID [[#345](https://github.com/hashicorp/terraform-provider-consul/issues/345)].

## 2.17.0 (December 13, 2022)

CHANGES:

* The `token` argument in the resources and datasources has been deprecated. The `token` argument in the provider configuration should be used instead ([[#332](https://github.com/hashicorp/terraform-provider-consul/issues/332)] and [[#298](https://github.com/hashicorp/terraform-provider-consul/issues/298)]).

NEW FEATURES:

* The `consul_config_entry` datasource can now be used to read the configuration of a given config entry ([[#318](https://github.com/hashicorp/terraform-provider-consul/issues/318)]).

IMPROVEMENTS:

* The `consul_node` resource can now be imported ([[#323](https://github.com/hashicorp/terraform-provider-consul/issues/323)]).
* The `consul_config_entry` resource can now be imported ([Gh-319] and [[#316](https://github.com/hashicorp/terraform-provider-consul/issues/316)]).
* The `consul_peering` resource has been updated to support the changes in Consul 1.14 ([[#328](https://github.com/hashicorp/terraform-provider-consul/issues/328)]).
* The `consul_peering`, `consul_peerings` datasources have been updated to support the changes in Consul 1.14 ([[#328](https://github.com/hashicorp/terraform-provider-consul/issues/328)]).
* The `consul_config_entry` resource now support the new peer parameters introduced in Consul 1.14 ([[#328](https://github.com/hashicorp/terraform-provider-consul/issues/328)]).

## 2.16.0 (September 24, 2022)

NEW FEATURES:

* The `consul_peering` and `consul_peering_token` resources can now be used to manage Consul Cluster Peering configuration ([#309](https://github.com/hashicorp/terraform-provider-consul/pull/309)).
* The `consul_peering` and `consul_peerings` datasources can now be used to introspect Consul Cluster Peering configuration ([#309](https://github.com/hashicorp/terraform-provider-consul/pull/309)).

IMPROVEMENTS:

* The `consul_acl_token_secret_id` datasource now supports reading tokens in a Consul Admin Partition ([#315](https://github.com/hashicorp/terraform-provider-consul/pull/315)).

## 2.15.1 (April 11, 2022)

BUG FIXES:

* The support of Admin Partition has been fixed for `consul_config_entry`: a new `partition` argument is now present and should be used instead of setting `Partition` in `config_json`.

## 2.15.0 (March 21, 2022)

CHANGES:

* The `consul_license` resource is now deprecated and will be removed in a future version of the provider ([#292](https://github.com/hashicorp/terraform-provider-consul/issues/292)).

NEW FEATURES:

* The `consul_admin_partition` resource can now be used to manage Consul Admin Partitions ([#292](https://github.com/hashicorp/terraform-provider-consul/issues/292)).
* The `consul_datacenters` datasource can now be used to get the list of known datacenters ([#290](https://github.com/hashicorp/terraform-provider-consul/issues/290) and [#293](https://github.com/hashicorp/terraform-provider-consul/issues/293)).
* The `consul_config_entry` now supports the `mesh` and `exported-services` kinds ([#292](https://github.com/hashicorp/terraform-provider-consul/issues/292)).

IMPROVEMENTS:

* The `consul_acl_auth_method`, `consul_acl_binding_rule`, `consul_acl_policy`, `consul_acl_role`, `consul_acl_token`, `consul_key_prefix`, `consul_keys`, `consul_namespace` and `consul_node` can now manage resources in Admin Partitions ([#292](https://github.com/hashicorp/terraform-provider-consul/issues/292)).

* The `consul_acl_auth_method`, `consul_acl_policy`, `consul_acl_role`, `consul_acl_token`, `consul_key_prefix` and `consul_keys` datasources can now look for a resource in a specific Admin Partition ([#292](https://github.com/hashicorp/terraform-provider-consul/issues/292)).

## 2.14.0 (October 01, 2021)

IMPROVEMENTS:

* The `consul_acl_role` resource now has a `node_identities` argument ([#284](https://github.com/hashicorp/terraform-provider-consul/issues/284) and [#287](https://github.com/hashicorp/terraform-provider-consul/issues/287)).
* The `consul_acl_role` datasource now has a `node_identities` attribute ([#284](https://github.com/hashicorp/terraform-provider-consul/issues/284) and [#287](https://github.com/hashicorp/terraform-provider-consul/issues/287)).
* The `consul_acl_token` resource now supports `service_identities`, `node_identities`, and `expiration_time` arguments ([#284](https://github.com/hashicorp/terraform-provider-consul/issues/284) and [#287](https://github.com/hashicorp/terraform-provider-consul/issues/287)).
* The `consul_acl_token` datasource now supports `roles`, `service_identities`, `node_identities`, and `expiration_time` attributes ([#284](https://github.com/hashicorp/terraform-provider-consul/issues/284) and [#287](https://github.com/hashicorp/terraform-provider-consul/issues/287)).

## 2.13.0 (August 19, 2021)

NEW FEATURES:

* The `consul_namespace_policy_attachment` and the `consul_namespace_role_attachment` can now be used to attach a default policy or role to an already existing Consul namespace ([#247](https://github.com/hashicorp/terraform-provider-consul/issues/247) and [#267](https://github.com/hashicorp/terraform-provider-consul/issues/267)).
* Additional headers can be set to be sent with each requests to the Consul server using the `header` block in the provider configuration ([#245](https://github.com/hashicorp/terraform-provider-consul/issues/245)).

IMPROVEMENTS:

* The `consul_namespace` resource can now be imported ([#247](https://github.com/hashicorp/terraform-provider-consul/issues/247) and [#263](https://github.com/hashicorp/terraform-provider-consul/issues/263)).
* The `tags` attribute in the `consul_services` datasource that was previously documented but missing is now present ([#274](https://github.com/hashicorp/terraform-provider-consul/issues/274)).

BUG FIXES:

* The `consul_service` now properly detect drift in the health-check configuration ([#237](https://github.com/hashicorp/terraform-provider-consul/issues/237) and [#262](https://github.com/hashicorp/terraform-provider-consul/issues/262)).
* Detecting incorrect configuration in `consul_acl_auth_method` is now deferred until read to unsure the configuration has been interpolated ([#260](https://github.com/hashicorp/terraform-provider-consul/issues/260) and [#261](https://github.com/hashicorp/terraform-provider-consul/issues/261)).


## 2.12.0 (May 12, 2021)

NEW FEATURES:

* The darwin/arm64 platform is now supported ([#253](https://github.com/hashicorp/terraform-provider-consul/issues/253) and [#254](https://github.com/hashicorp/terraform-provider-consul/issues/254)).
* The `consul_acl_token_role_attachment` can be used to attach a role to a token created outside the Terraform configuration ([#252](https://github.com/hashicorp/terraform-provider-consul/issues/252)).

IMPROVEMENTS:

* The `consul_acl_token_secret_id` datasource can now look for a token in a namespace ([#242](https://github.com/hashicorp/terraform-provider-consul/issues/242)).
* The `Content-Type` header is now present in all `PUT` and `POST` HTTP requests sent by the provider ([#255](https://github.com/hashicorp/terraform-provider-consul/issues/255)).
* The `consul_config_entry` can now asssociate a config entry with a namespace ([#246](https://github.com/hashicorp/terraform-provider-consul/issues/246) and [#256](https://github.com/hashicorp/terraform-provider-consul/issues/256)).
* The `consul_key_prefix` can now be used to manage keys at the root of the key-value store ([#258](https://github.com/hashicorp/terraform-provider-consul/issues/258)).

BUG FIXES:

* The `consul_acl_auth_method` now correctly detects changes in the configuration ([#240](https://github.com/hashicorp/terraform-provider-consul/issues/240) and [#244](https://github.com/hashicorp/terraform-provider-consul/issues/244)).
* All resources and datasources now properly inherit the `datacenter`, `token` and `namespace` configuration from the provider when they are not set in the resource ([#8](https://github.com/hashicorp/terraform-provider-consul/issues/8) and [#259](https://github.com/hashicorp/terraform-provider-consul/issues/259))

## 2.11.0 (January 14, 2021)

NEW FEATURES:

* The `consul_config_entry` can now be used to manage Consul service intentions ([[#232](https://github.com/hashicorp/terraform-provider-consul/issues/232)] and [[#235](https://github.com/hashicorp/terraform-provider-consul/issues/235)]).

## 2.10.1 (October 22, 2020)

BUG FIXES:

* The same API client is now reused across all operations ([[#233](https://github.com/hashicorp/terraform-provider-consul/issues/233)]).

## 2.10.0 (September 18, 2020)

NEW FEATURES:

* The TLS configuration of the provider can now be given directly as strings instead of using files ([[#220](https://github.com/hashicorp/terraform-provider-consul/issues/220)] and [[#5](https://github.com/hashicorp/terraform-provider-consul/issues/5)]).
* The `consul_intention` resource now has a `datacenter` argument ([[#213](https://github.com/hashicorp/terraform-provider-consul/issues/213)] and [[#214](https://github.com/hashicorp/terraform-provider-consul/issues/214)]).
* The `consul_intention` resource can now be imported ([[#222](https://github.com/hashicorp/terraform-provider-consul/issues/222)] and [[#225](https://github.com/hashicorp/terraform-provider-consul/issues/225)]).

BUG FIXES:

* The `consul_acl_binding_rule` now delegates the validation of `bind_type` to Consul and supports the `node` bind type ([[#217](https://github.com/hashicorp/terraform-provider-consul/issues/217)] and [[#218](https://github.com/hashicorp/terraform-provider-consul/issues/218)]).
* The `CONSUL_HTTP_SSL` environment variable can now be used to force the use of SSL like it does for the Consul CLI ([[#215](https://github.com/hashicorp/terraform-provider-consul/issues/215)] and [[#219](https://github.com/hashicorp/terraform-provider-consul/issues/219)]).
* The `flags` attribute has been removed from the `consul_license` resource to make it work with Consul 1.8 ([[#223](https://github.com/hashicorp/terraform-provider-consul/issues/223)] and [[#227](https://github.com/hashicorp/terraform-provider-consul/issues/227)]).

## 2.9.0 (July 23, 2020)

NEW FEATURES:

* The new `consul_certificate_authority` can be used to manage the Consul Connect Certificate Authority ([[#205](https://github.com/hashicorp/terraform-provider-consul/issues/205)]).

* `consul_acl_auth_method` now supports the `display_name`, `max_token_ttl`, `token_locality` and `namespace_rule` attributes ([[#204](https://github.com/hashicorp/terraform-provider-consul/issues/204)]).

* The `consul_service` and `consul_service_health` data sources now support filter expressions ([[#203](https://github.com/hashicorp/terraform-provider-consul/issues/203)]).

* The `consul_config_entry` resource now support Ingress and Terminating Gateways ([[#199](https://github.com/hashicorp/terraform-provider-consul/issues/199)] and [[#202](https://github.com/hashicorp/terraform-provider-consul/issues/202)]).

* The `consul_service` resource now has a `enable_tag_override` attribute ([[#201](https://github.com/hashicorp/terraform-provider-consul/issues/201)]).

* The `consul_acl_auth_method` resource now has a `config_json` attribute to use an arbitrary complex configuration. The `config` attribute has been deprecated and will be removed in a future release ([[#208](https://github.com/hashicorp/terraform-provider-consul/issues/208)] and [[#209](https://github.com/hashicorp/terraform-provider-consul/issues/209)]).

* The `consul_acl_auth_method` data source now has a `config_json` attribute. The `config` attribute has been deprecated as it will be blank when the configuration is too complex and it will be removed in a future release ([[#208](https://github.com/hashicorp/terraform-provider-consul/issues/208)] and [[#209](https://github.com/hashicorp/terraform-provider-consul/issues/209)]).

## 2.8.0 (May 21, 2020)

NEW FEATURES:

* The `ignore_check_ids`, `node_meta` and `service_meta` have been added to the `consul_prepared_query` resource ([[#192](https://github.com/hashicorp/terraform-provider-consul/issues/192)] and [[#193](https://github.com/hashicorp/terraform-provider-consul/issues/193)]).

BUG FIXES:

* The `subkey` attribute of the `consul_key_prefix` resource now detects external changes ([[#189](https://github.com/hashicorp/terraform-provider-consul/issues/189)]).

## 2.7.0 (March 26, 2020)

NEW FEATURES:

* The `consul_acl_role` can now be imported ([[#182](https://github.com/hashicorp/terraform-provider-consul/issues/182)]).
* Roles can be attached to `consul_acl_token` using the new `roles` attribute ([[#178](https://github.com/hashicorp/terraform-provider-consul/issues/178)] and [[#180](https://github.com/hashicorp/terraform-provider-consul/issues/180)]]).
* The `consul_namespace` resource can now be used to manage namespaces in a Consul Enterprise cluster ([[#183](https://github.com/hashicorp/terraform-provider-consul/issues/183)]).
* The `consul_acl_auth_method`, `consul_acl_binding_rule`, `consul_acl_policy`, `consul_acl_role`, `consul_acl_token`, `consul_intention`, `consul_key_prefix`, `consul_keys` and `consul_service` can now be associated to a namespace in a Consul Enterprise cluster ([[#183](https://github.com/hashicorp/terraform-provider-consul/issues/183)]).
* The `consul_acl_auth_method`, `consul_acl_policy`, `consul_acl_role`, `consul_acl_token`, `consul_key_prefix`, `consul_keys`, `consul_nodes`, `consul_service`, `consul_services` datasources can now be used to query a specific namespace in a Consul Enterprise cluster ([[#183](https://github.com/hashicorp/terraform-provider-consul/issues/183)]).
* The `consul_license` resource can now be used to manage automatically the license of a Consul Enterprise cluster ([[#172](https://github.com/hashicorp/terraform-provider-consul/issues/172)] and [[#173](https://github.com/hashicorp/terraform-provider-consul/issues/173)]).
* The `consul_network_area` and `consul_network_area_members` can now be used to manage the network areas of a Consul Enterprise cluster ([[#175](https://github.com/hashicorp/terraform-provider-consul/issues/175)]).
* The `consul_network_segments` can now be used to manage the network segments of a Consul Enterprise cluster ([[#175](https://github.com/hashicorp/terraform-provider-consul/issues/175)]).

BUG FIXES:

* Importing `consul_key_prefix` no longer delete and replace all keys ([[#169](https://github.com/hashicorp/terraform-provider-consul/issues/169)] and [[#171](https://github.com/hashicorp/terraform-provider-consul/issues/171)]).

## 2.6.1 (December 07, 2019)

IMPROVEMENTS:

* The `consul_keys` diffs are now easier to read.

BUG FIXES:

* The `CONSUL_CLIENT_CERT`, `CONSUL_CLIENT_KEY`, `CONSUL_CACERT` and `CONSUL_CAPATH` are now used to set the TLS configuration of the provider.
* The `consul_keys` can now create keys with an empty string as value.

## 2.6.0 (October 25, 2019)

NEW FEATURES:

* The `consul_acl_role`, `consul_acl_auth_method` and `consul_acl_binding_rule` can now be used to manage the Consul ACL system ([[#128](https://github.com/hashicorp/terraform-provider-consul/issues/128)] and [[#123](https://github.com/hashicorp/terraform-provider-consul/issues/123)]).

* The `consul_acl_token_policy_attachment` can be used to attach a policy to a token created outside the Terraform configuration ([[#130](https://github.com/hashicorp/terraform-provider-consul/issues/130)] and [[#125](https://github.com/hashicorp/terraform-provider-consul/issues/125)]).

* The `consul_config_entry` can now be used to manage Consul Configuration Entries ([[#127](https://github.com/hashicorp/terraform-provider-consul/issues/127)]).

* The `consul_acl_auth_method`, `consul_acl_policy`, `consul_acl_role` datasources can now be used to retrieve information about the Consul ACL objects ([[#153](https://github.com/hashicorp/terraform-provider-consul/issues/153)]).

* The `consul_acl_token` can now be used to read public token information ([[#137](https://github.com/hashicorp/terraform-provider-consul/issues/137)] and [[#126](https://github.com/hashicorp/terraform-provider-consul/issues/126)]).

* The `consul_acl_token_secret_id` can now be used to read a token secret ID ([[#137](https://github.com/hashicorp/terraform-provider-consul/issues/137)] and [[#126](https://github.com/hashicorp/terraform-provider-consul/issues/126)]).

IMPROVEMENTS:

* The `consul_service` resource can now set the service metadata ([[#122](https://github.com/hashicorp/terraform-provider-consul/issues/122)]).

* The `consul_service` datasource now returns the service metadata ([[#148](https://github.com/hashicorp/terraform-provider-consul/issues/148)] and [[#132](https://github.com/hashicorp/terraform-provider-consul/issues/132)]).

BUG FIXES:

* The `consul_prepared_query` now handles default values correctly for the `failover`, `dns` and `template` blocks ([[#119](https://github.com/hashicorp/terraform-provider-consul/issues/119)] and [[#121](https://github.com/hashicorp/terraform-provider-consul/issues/121)])

* The `consul_service` resource correctly associates a service instance with its health-checks ([[#147](https://github.com/hashicorp/terraform-provider-consul/issues/147)] and [[#146](https://github.com/hashicorp/terraform-provider-consul/issues/146)]).

* Fix the `check_id` and `status` attribute of health-checks in the `consul_service` resources that would always mark the plan as dirty ([[#142](https://github.com/hashicorp/terraform-provider-consul/issues/142)]).

## 2.5.0 (June 03, 2019)

NEW FEATURES:

* The Consul Terraform provider is now compatible with Terraform 0.12 ([[#118](https://github.com/hashicorp/terraform-provider-consul/issues/118)] and [[#88](https://github.com/hashicorp/terraform-provider-consul/issues/88)]).


## 2.4.0 (May 29, 2019)

NEW FEATURES:

* New resource: the `consul_service_health` can now be used to fetch healthy instances of a service ([[#87](https://github.com/hashicorp/terraform-provider-consul/issues/87)] and [[#89](https://github.com/hashicorp/terraform-provider-consul/issues/89)])

IMPROVEMENTS:

*  The `consul_prepared_query` resource now supports Consul Connect ([[#107](https://github.com/hashicorp/terraform-provider-consul/issues/107)])
*  The `consul_acl_token` and `consul_acl_policy` resources are now importable ([[#103](https://github.com/hashicorp/terraform-provider-consul/issues/103)])

BUG FIXES:

* Tokens attribute nested in a resource attribute are now marked as sensitive so they won't appear in the ouput and the logs ([[#106](https://github.com/hashicorp/terraform-provider-consul/issues/106)])
* The default datacenter is left empty if it cannot be read from the agent and is not set in the provider configuration ([[#99](https://github.com/hashicorp/terraform-provider-consul/issues/99)], [[#97](https://github.com/hashicorp/terraform-provider-consul/issues/97)] and [[#105](https://github.com/hashicorp/terraform-provider-consul/issues/105)])
* The attributes `failover`, `dns` and `template` of the `consul_prepared_query` resource are now set correctly ([[#109](https://github.com/hashicorp/terraform-provider-consul/issues/109)] and [[#108](https://github.com/hashicorp/terraform-provider-consul/issues/108)])
* The `consul_acl_token` resource can now be updated and does not crashes Terraform anymore ([[#102](https://github.com/hashicorp/terraform-provider-consul/issues/102)])
* The `consul_node` resource now detect external changes made to its `address` and `meta` attributes ([[#104](https://github.com/hashicorp/terraform-provider-consul/issues/104)])
* The `external` attribute of the `consul_service` resource has been deprecated ([[#104](https://github.com/hashicorp/terraform-provider-consul/issues/104)])
* The `local` attribute is now correctly marked as requiring the creation of a new ACL token in the `consul_acl_token` resource ([[#117](https://github.com/hashicorp/terraform-provider-consul/issues/117)])

## 2.3.0 (April 09, 2019)

NEW FEATURES:

* New resources: `consul_acl_policy` and `consul_acl_token` can now be used to manage Consul ACLs with Terraform. ([[!60](https://github.com/hashicorp/terraform-provider-consul/pull/60)])
* New resource: the `consul_autopilot_config` resource can now be used to manage the [Consul Autopilot](https://learn.hashicorp.com/consul/day-2-operations/advanced-operations/autopilot) configuration ([[!86](https://github.com/hashicorp/terraform-provider-consul/pull/86)]).
* New datasource: The `consul_autopilot_health` datasource returns the [autopilot health information](https://www.consul.io/api/operator/autopilot.html#read-health) of the Consul cluster ([[!84](https://github.com/hashicorp/terraform-provider-consul/pull/84)])

IMPROVEMENTS:

* `consul_service` can now manage health-checks associated with the service. ([[!64](https://github.com/hashicorp/terraform-provider-consul/pull/64)] and [[#54](https://github.com/hashicorp/terraform-provider-consul/issues/54)])
* The `ca_path` attribute of the provider configuration can now be used to indicate a directory containing certificate files. ([[!80](https://github.com/hashicorp/terraform-provider-consul/pull/80)] and [[!79](https://github.com/hashicorp/terraform-provider-consul/issues/79)])
* The `consul_prepared_query` resource can now be imported. ([[!94](https://github.com/hashicorp/terraform-provider-consul/pull/94)])
* The `consul_key_prefix` resource can now be imported. ([[!78](https://github.com/hashicorp/terraform-provider-consul/pull/78)] and [[#77](https://github.com/hashicorp/terraform-provider-consul/issues/77)])
* `consul_keys` and `consul_key_prefix` can now manage flags associated with each key. ([[!71](https://github.com/hashicorp/terraform-provider-consul/pull/71)] and [[#59](https://github.com/hashicorp/terraform-provider-consul/issues/59)])

BUG FIXES:

* `consul_intention`, `consul_node` and `consul_service` now correctly re-creates
resources deleted out-of-band ([#81](https://github.com/hashicorp/terraform-provider-consul/issues/81) and [!69](https://github.com/hashicorp/terraform-provider-consul/pull/69)).
* Consul tokens no longer appear in the logs and the standard output. ([[!73](https://github.com/hashicorp/terraform-provider-consul/pull/73)] and [[#50](https://github.com/hashicorp/terraform-provider-consul/issues/50)])

## 2.2.0 (October 03, 2018)

IMPROVEMENTS:

* The `consul_node` resource now supports setting node metadata via the `meta` attribute. ([#65](https://github.com/hashicorp/terraform-provider-consul/issues/65))


## 2.1.0 (June 26, 2018)

NEW FEATURES:

* New resource: `consul_intention`. This provides management of intentions in Consul Connect, used for service authorization.  ([#53](https://github.com/hashicorp/terraform-provider-consul/issues/53))

## 2.0.0 (June 22, 2018)

NOTES:

* The `consul_catalog_entry` resource has been deprecated and will be removed in a future release. Please use the [upgrade guide in the documentation](https://www.terraform.io/docs/providers/consul/upgrading.html#upgrading-to-2-0-0) to migrate to `consul_service` or `consul_node` as appropriate. ([#49](https://github.com/hashicorp/terraform-provider-consul/issues/49))
* The `consul_agent_service` resource has been deprecated and will be removed in a future release. Please use the [upgrade guide in the documentation](https://www.terraform.io/docs/providers/consul/upgrading.html#upgrading-to-2-0-0) to migrate to `consul_service`. ([#49](https://github.com/hashicorp/terraform-provider-consul/issues/49))
* The `consul_agent_self` datasource has been deprecated and will be removed in a future release. Please use the [upgrade guide in the documentation](https://www.terraform.io/docs/providers/consul/upgrading.html#upgrading-to-2-0-0) to migrate to `consul_agent_config` if applicable. ([#49](https://github.com/hashicorp/terraform-provider-consul/issues/49))

IMPROVEMENTS:

* The `consul_service` resource has been modified to use the Consul catalog APIs. The `node` attribute is now required, and nodes that do not exist will not be created automatically. Please see the upgrade guide in the documentation for more detail. ([#49](https://github.com/hashicorp/terraform-provider-consul/issues/49))
* `consul_catalog_*` data sources have been renamed to remove catalog, for clarity. Both will work going forward, with the catalog version potentially being deprecated on a future date. ([#49](https://github.com/hashicorp/terraform-provider-consul/issues/49))
* The provider now uses the post-1.0 version of the Consul API. ([#49](https://github.com/hashicorp/terraform-provider-consul/issues/49))

## 1.1.0 (June 15, 2018)

NOTES:

* This will be the last release prior to significant changes coming to the provider that will deprecate
multiple resources. These changes are primarily to simplify overlap of resource functionality, ensure we are using the correct APIs/design provided by Consul for something like Terraform, and remove resources that can no longer be supported by the current version of the Consul API (1.0+). Read more [here](https://github.com/hashicorp/terraform-provider-consul/issues/46).

IMPROVEMENTS:

* The provider now allows you to skip HTTPS certificate verification by supplying the `insecure_https` option. ([#31](https://github.com/hashicorp/terraform-provider-consul/issues/31))

NEW FEATURES:

* New data source: `consul_agent_config`. This new datasource provides information similar to `consul_agent_self`,
but is designed to only expose configuration that Consul will not change without versioning upstream. ([#42](https://github.com/hashicorp/terraform-provider-consul/issues/42))
* New data source: `consul_key_prefix` corresponds to the existing resource of the same name, allowing config to access a set of keys with a common prefix as a Terraform map for more convenient access ([#34](https://github.com/hashicorp/terraform-provider-consul/issues/34))

BUG FIXES:

* `consul_catalog` resource now correctly re-creates resources deleted out-of-band. ([#30](https://github.com/hashicorp/terraform-provider-consul/issues/30))
* `consul_service` resource type now correctly detects when a service has been deleted outside of Terraform, flagging it for re-creation rather than returning an error ([#33](https://github.com/hashicorp/terraform-provider-consul/issues/33))
* `consul_catalog_service` data source now accepts the `tag` and `datacenter` arguments, as was described in documentation ([#32](https://github.com/hashicorp/terraform-provider-consul/issues/32))

## 1.0.0 (September 26, 2017)

BUG FIXES:

* d/consul_agent_self: The `enable_ui` config setting was always set to false, regardless of the actual agent configuration ([#16](https://github.com/hashicorp/terraform-provider-consul/issues/16))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
