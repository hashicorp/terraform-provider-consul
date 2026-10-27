// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// baselineRelease is the released provider version users upgrade from. It is 2.23.0, the last release built on
// terraform-plugin-sdk v1, for the SDK v2 migration.
// To test a later upgrade, set it (or TF_ACC_UPGRADE_BASELINE for a single run) to the release users will
// upgrade from (e.g. the previous release); exclusions that mention the baseline release may then no longer apply.
const baselineRelease = "2.23.0"

// upgradeExclusion is why a config cannot be upgrade-tested; such cases are reported as skipped with the reason.
type upgradeExclusion string

const (
	included          upgradeExclusion = ""
	excludeApplyFails upgradeExclusion = "apply is expected to fail, so there is no state to upgrade"
	// TODO: include once the jwt-provider ClockSkewSeconds drift is fixed and the baseline release contains the fix.
	excludeDriftsOnBaseline upgradeExclusion = "the baseline release plans a ClockSkewSeconds change on its jwt-provider config entry right after apply, so there is no stable state to upgrade from"
	// TODO: include once the baseline release supports passive_health_check thresholds (added after 2.23.0).
	excludeNewerThanBaseline upgradeExclusion = "uses passive_health_check attributes the baseline release does not have"
	// TODO: include once both provider names (consul, consulremote) can be pinned to the baseline release in step 1.
	excludeSecondProviderName upgradeExclusion = "its second provider name (consulremote) maps to a different provider address than the baseline release"
	// TODO: include once the fixture creates the partitions and namespaces it uses.
	excludeMissingEnterpriseScopes upgradeExclusion = "uses partitions or namespaces that neither the config nor the test server creates"
)

// TestAccConsulProviderUpgrade tests upgrading the provider release, as users do: state created by the baseline
// release must be read by this build without planning any change. Added for the terraform-plugin-sdk v1 to v2
// migration, which changes how existing state is read.
// It downloads the baseline release and takes several minutes, so it only runs when TF_ACC_UPGRADE is set.
func TestAccConsulProviderUpgrade(t *testing.T) {
	if os.Getenv("TF_ACC_UPGRADE") == "" {
		t.Skip("set TF_ACC_UPGRADE=1 to run the upgrade test from the baseline release")
	}
	baseline := baselineRelease
	if v := os.Getenv("TF_ACC_UPGRADE_BASELINE"); v != "" {
		baseline = v
	}
	t.Logf("upgrading from hashicorp/consul %s", baseline)
	tests := []struct {
		name       string
		config     string
		edition    string // "ce" or "ent" when the config only applies to one Consul edition
		minVersion string
		remoteDC   bool // needs the second datacenter (dc2)
		exclusion  upgradeExclusion
	}{
		{"ACLAuthMethodConfigBasic_NoConfig", testResourceACLAuthMethodConfigBasic_NoConfig, "", "", false, excludeApplyFails},
		{"ACLAuthMethodConfigBasic", testResourceACLAuthMethodConfigBasic, "", "", false, included},
		{"ACLAuthMethodConfigBasic_Update", testResourceACLAuthMethodConfigBasic_Update, "", "", false, included},
		{"ACLAuthMethodConfigBasic_ConfigJSON", testResourceACLAuthMethodConfigBasic_ConfigJSON, "", "", false, included},
		{"ACLAuthMethodConfigBasic_ConfigJSONUpdate", testResourceACLAuthMethodConfigBasic_ConfigJSONUpdate, "", "", false, included},
		{"ACLAuthMethodConfigBasicConfigJSON", testResourceACLAuthMethodConfigBasicConfigJSON, "", "", false, included},
		{"ACLAuthMethod_interpolation", testResourceACLAuthMethod_interpolation, "", "", false, included},
		{"ACLAuthMethodNamespaceCE", testResourceACLAuthMethodNamespaceCE, "", "", false, excludeApplyFails},
		{"ACLAuthMethodNamespaceEE", testResourceACLAuthMethodNamespaceEE, "ent", "", false, included},
		{"ACLAuthMethodNamespaceEE_namespaceRule", testResourceACLAuthMethodNamespaceEE_namespaceRule, "ent", "", false, included},
		{"ACLBindingRuleConfigBasic", testResourceACLBindingRuleConfigBasic, "", "", false, included},
		{"ACLBindingRuleConfigBasic_Update", testResourceACLBindingRuleConfigBasic_Update, "", "", false, included},
		{"ACLBindingRuleConfig_node", testResourceACLBindingRuleConfig_node, "", "", false, included},
		{"ACLBindingRuleConfig_templatedPolicyWithNoVariables", testResourceACLBindingRuleConfig_templatedPolicyWithNoVariables, "", "1.17.0", false, included},
		{"ACLBindingRuleConfig_templatedPolicyWithVariables", testResourceACLBindingRuleConfig_templatedPolicyWithVariables, "", "1.17.0", false, included},
		{"ACLBindingRuleConfig_wrongType", testResourceACLBindingRuleConfig_wrongType, "", "", false, excludeApplyFails},
		{"ACLBindingRuleConfig_namespaceCE", testResourceACLBindingRuleConfig_namespaceCE, "", "", false, excludeApplyFails},
		{"ACLBindingRuleConfig_namespaceEE", testResourceACLBindingRuleConfig_namespaceEE, "ent", "", false, included},
		{"ACLPolicyConfigBasic", testResourceACLPolicyConfigBasic, "", "", false, included},
		{"ACLPolicyConfigBasicUpdate", testResourceACLPolicyConfigBasicUpdate, "", "", false, included},
		{"ACLPolicyNamespaceCE", testResourceACLPolicyNamespaceCE, "", "", false, excludeApplyFails},
		{"ACLPolicyNamespaceEE", testResourceACLPolicyNamespaceEE, "ent", "", false, included},
		{"ACLRolePolicyAttachmentConfigBasic", testResourceACLRolePolicyAttachmentConfigBasic, "", "", false, included},
		{"ACLRolePolicyAttachmentConfigUpdate", testResourceACLRolePolicyAttachmentConfigUpdate, "", "", false, included},
		{"ACLRolePolicyAttachmentConfigMissingRole", testResourceACLRolePolicyAttachmentConfigMissingRole, "", "", false, excludeApplyFails},
		{"ACLRoleConfigBasic", testResourceACLRoleConfigBasic, "", "", false, included},
		{"ACLRoleConfigUpdate", testResourceACLRoleConfigUpdate, "", "", false, included},
		{"ACLRoleConfigUpdateTemplatedPolicies", testResourceACLRoleConfigUpdateTemplatedPolicies, "", "1.17.0", false, included},
		{"ACLRoleConfigPolicyName", testResourceACLRoleConfigPolicyName, "", "", false, included},
		{"ACLRoleNamespaceCE", testResourceACLRoleNamespaceCE, "", "", false, excludeApplyFails},
		{"ACLRoleNamespaceEE", testResourceACLRoleNamespaceEE, "ent", "", false, included},
		{"ACLTokenPolicyAttachmentConfigBasic", testResourceACLTokenPolicyAttachmentConfigBasic, "", "", false, included},
		{"ACLTokenPolicyAttachmentConfigUpdate", testResourceACLTokenPolicyAttachmentConfigUpdate, "", "", false, included},
		{"ACLTokenRoleAttachmentConfigBasic", testResourceACLTokenRoleAttachmentConfigBasic, "", "", false, included},
		{"ACLTokenRoleAttachmentConfigUpdate", testResourceACLTokenRoleAttachmentConfigUpdate, "", "", false, included},
		{"ACLTokenConfigBasic", testResourceACLTokenConfigBasic, "", "", false, included},
		{"ACLTokenConfigUpdate", testResourceACLTokenConfigUpdate, "", "", false, included},
		{"ACLTokenConfigUpdateTemplatedPolicies", testResourceACLTokenConfigUpdateTemplatedPolicies, "", "1.17.0", false, included},
		{"ACLTokenConfigRole", testResourceACLTokenConfigRole, "", "", false, included},
		{"ACLTokenConfigNamespaceCE", testResourceACLTokenConfigNamespaceCE, "", "", false, excludeApplyFails},
		{"ACLTokenConfigNamespaceEE", testResourceACLTokenConfigNamespaceEE, "ent", "", false, included},
		{"ConsulAdminPartitionBasic", testAccConsulAdminPartitionBasic, "ent", "", false, included},
		{"ConsulAgentServiceConfig", testAccConsulAgentServiceConfig, "", "", false, included},
		{"ConsulAutopilotConfigBasic", testAccConsulAutopilotConfigBasic, "", "", false, included},
		{"ConsulAutopilotConfig", testAccConsulAutopilotConfig, "", "", false, included},
		{"ConsulAutopilotConfigParseDuration", testAccConsulAutopilotConfigParseDuration, "", "", false, excludeApplyFails},
		{"ConsulAutopilotDatacenter", testAccConsulAutopilotDatacenter, "", "", true, included},
		{"ConsulCatalogEntryConfig", testAccConsulCatalogEntryConfig, "", "", false, included},
		{"ConsulCertificateAuthorityConfig", testAccConsulCertificateAuthorityConfig, "", "", false, included},
		{"ConsulCertificateAuthorityConfigBoth", testAccConsulCertificateAuthorityConfigBoth, "", "", false, excludeApplyFails},
		{"ConsulCertificateAuthorityConfigJSON", testAccConsulCertificateAuthorityConfigJSON, "", "", false, included},
		{"ConsulConfigEntryCE_ServiceDefaults", testAccConsulConfigEntryCE_ServiceDefaults, "ce", "", false, included},
		{"ConsulConfigEntryCE_ServiceDefaultsOptionalField", testAccConsulConfigEntryCE_ServiceDefaultsOptionalField, "ce", "", false, included},
		{"ConsulConfigEntryCE_ProxyDefaults", testAccConsulConfigEntryCE_ProxyDefaults, "ce", "", false, included},
		{"ConsulConfigEntryCE_ServiceRouter", testAccConsulConfigEntryCE_ServiceRouter, "ce", "", false, included},
		{"ConsulConfigEntryCE_ServiceSplitter", testAccConsulConfigEntryCE_ServiceSplitter, "ce", "", false, included},
		{"ConsulConfigEntryCE_ServiceResolver", testAccConsulConfigEntryCE_ServiceResolver, "ce", "", false, included},
		{"ConsulConfigEntryCE_IngressGateway", testAccConsulConfigEntryCE_IngressGateway, "ce", "", false, included},
		{"ConsulConfigEntryCE_TerminatingGateway", testAccConsulConfigEntryCE_TerminatingGateway, "ce", "", false, included},
		{"ConsulConfigEntryCE_ServiceConfigL4", testAccConsulConfigEntryCE_ServiceConfigL4, "ce", "", false, excludeDriftsOnBaseline},
		{"ConsulConfigEntryCE_ServiceConfigL7", testAccConsulConfigEntryCE_ServiceConfigL7, "ce", "", false, excludeDriftsOnBaseline},
		{"ConsulConfigEntryCE_ServiceConfigL7b", testAccConsulConfigEntryCE_ServiceConfigL7b, "ce", "", false, excludeDriftsOnBaseline},
		{"ConsulConfigEntryCE_ServiceConfigL7gRPC", testAccConsulConfigEntryCE_ServiceConfigL7gRPC, "ce", "", false, included},
		{"ConsulConfigEntryCE_ServiceConfigL7Mixed", testAccConsulConfigEntryCE_ServiceConfigL7Mixed, "ce", "", false, included},
		{"ConsulConfigEntryCE_HTTPRoute", testAccConsulConfigEntryCE_HTTPRoute, "ce", "", false, included},
		{"ConsulConfigEntryCE_exportedServicesCE", TestAccConsulConfigEntryCE_exportedServicesCE, "ce", "", false, included},
		{"ConsulConfigEntryCE_meshCE", TestAccConsulConfigEntryCE_meshCE, "ce", "", false, included},
		{"ConsulConfigEntryCE_jwtRemote", TestAccConsulConfigEntryCE_jwtRemote, "ce", "", false, excludeDriftsOnBaseline},
		{"ConsulConfigEntryCE_jwtLocal", TestAccConsulConfigEntryCE_jwtLocal, "ce", "", false, excludeDriftsOnBaseline},
		{"ConsulConfigEntryEE_ServiceDefaults", testAccConsulConfigEntryEE_ServiceDefaults, "ent", "", false, included},
		{"ConsulConfigEntryEE_ServiceDefaultsPartition", testAccConsulConfigEntryEE_ServiceDefaultsPartition, "ent", "", false, included},
		{"ConsulConfigEntryEE_ProxyDefaults", testAccConsulConfigEntryEE_ProxyDefaults, "ent", "", false, included},
		{"ConsulConfigEntryEE_ServiceRouter", testAccConsulConfigEntryEE_ServiceRouter, "ent", "", false, included},
		{"ConsulConfigEntryEE_ServiceSplitter", testAccConsulConfigEntryEE_ServiceSplitter, "ent", "", false, included},
		{"ConsulConfigEntryEE_ServiceResolver", testAccConsulConfigEntryEE_ServiceResolver, "ent", "", false, included},
		{"ConsulConfigEntryEE_IngressGateway", testAccConsulConfigEntryEE_IngressGateway, "ent", "", false, included},
		{"ConsulConfigEntryEE_TerminatingGateway", testAccConsulConfigEntryEE_TerminatingGateway, "ent", "", false, included},
		{"ConsulConfigEntryEE_ServiceConfigAdminPartition", testAccConsulConfigEntryEE_ServiceConfigAdminPartition, "ent", "", false, included},
		{"ConsulConfigEntryEE_ServiceConfigL4", testAccConsulConfigEntryEE_ServiceConfigL4, "ent", "", false, excludeDriftsOnBaseline},
		{"ConsulConfigEntryEE_ServiceConfigL7", testAccConsulConfigEntryEE_ServiceConfigL7, "ent", "", false, excludeDriftsOnBaseline},
		{"ConsulConfigEntryEE_ServiceConfigL7b", testAccConsulConfigEntryEE_ServiceConfigL7b, "ent", "", false, excludeDriftsOnBaseline},
		{"ConsulConfigEntryEE_ServiceConfigL7gRPC", testAccConsulConfigEntryEE_ServiceConfigL7gRPC, "ent", "", false, included},
		{"ConsulConfigEntryEE_ServiceConfigL7Mixed", testAccConsulConfigEntryEE_ServiceConfigL7Mixed, "ent", "", false, included},
		{"ConsulConfigEntryEE_DefaultNamespace", testAccConsulConfigEntryEE_DefaultNamespace, "ent", "", false, included},
		{"ConsulConfigEntryEE_Namespace", testAccConsulConfigEntryEE_Namespace, "ent", "", false, included},
		{"ConsulConfigEntryEE_exportedServicesEE", TestAccConsulConfigEntryEE_exportedServicesEE, "ent", "", false, included},
		{"ConsulConfigEntryEE_meshEE", TestAccConsulConfigEntryEE_meshEE, "ent", "", false, included},
		{"ConsulConfigEntryEE_jwtRemote", TestAccConsulConfigEntryEE_jwtRemote, "ent", "", false, excludeDriftsOnBaseline},
		{"ConsulConfigEntryEE_jwtLocal", TestAccConsulConfigEntryEE_jwtLocal, "ent", "", false, excludeDriftsOnBaseline},
		{"ConsulServiceDefaultsConfigEntryWithDestination", testConsulServiceDefaultsConfigEntryWithDestination, "ce", "", false, included},
		{"ConsulServiceDefaultsConfigEntryWithUpstreamConfig", testConsulServiceDefaultsConfigEntryWithUpstreamConfig, "ce", "", false, included},
		{"ConsulServiceDefaultsPassiveHealthCheckConfigEntryWithDestination", testConsulServiceDefaultsPassiveHealthCheckConfigEntryWithDestination, "ce", "", false, included},
		{"ConsulServiceDefaultsPassiveHealthCheckConfigEntryWithUpstreamConfig", testConsulServiceDefaultsPassiveHealthCheckConfigEntryWithUpstreamConfig, "ce", "", false, excludeNewerThanBaseline},
		{"ConsulServiceDefaultsConfigEntryWithDestinationEE", testConsulServiceDefaultsConfigEntryWithDestinationEE, "ent", "", false, excludeMissingEnterpriseScopes},
		{"ConsulServiceDefaultsConfigEntryWithUpstreamConfigEE", testConsulServiceDefaultsConfigEntryWithUpstreamConfigEE, "ent", "", false, excludeMissingEnterpriseScopes},
		{"ConsulServiceDefaultsPassiveHealthCheckConfigEntryWithDestinationEE", testConsulServiceDefaultsPassiveHealthCheckConfigEntryWithDestinationEE, "ent", "", false, excludeMissingEnterpriseScopes},
		{"ConsulServiceDefaultsPassiveHealthCheckConfigEntryWithUpstreamConfigEE", testConsulServiceDefaultsPassiveHealthCheckConfigEntryWithUpstreamConfigEE, "ent", "", false, excludeNewerThanBaseline},
		{"ConsulConfigEntryServiceIntentionsCE", testConsulConfigEntryServiceIntentionsCE, "ce", "", false, included},
		{"ConsulConfigEntryServiceIntentionsEE", testConsulConfigEntryServiceIntentionsEE, "ent", "", false, included},
		{"ConsulConfigEntryServiceResolverCEWithRedirect", testConsulConfigEntryServiceResolverCEWithRedirect, "ce", "", false, included},
		{"ConsulConfigEntryServiceResolverCEWithFailover", testConsulConfigEntryServiceResolverCEWithFailover, "ce", "", false, included},
		{"ConsulConfigEntryServiceResolverEEWithRedirect", testConsulConfigEntryServiceResolverEEWithRedirect, "ent", "", false, included},
		{"ConsulConfigEntryServiceResolverEEWithFailover", testConsulConfigEntryServiceResolverEEWithFailover, "ent", "", false, included},
		{"ConsulConfigEntryServiceRouterCE_Empty", testConsulConfigEntryServiceRouterCE_Empty, "ce", "", false, included},
		{"ConsulConfigEntryServiceRouterCE", testConsulConfigEntryServiceRouterCE, "ce", "", false, included},
		{"ConsulConfigEntryServiceRouterCE_noMatch", testConsulConfigEntryServiceRouterCE_noMatch, "ce", "", false, included},
		{"ConsulConfigEntryServiceRouterCE_noDestination", testConsulConfigEntryServiceRouterCE_noDestination, "ce", "", false, included},
		{"ConsulConfigEntryServiceRouterEE_Empty", testConsulConfigEntryServiceRouterEE_Empty, "ent", "", false, included},
		{"ConsulConfigEntryServiceRouterEE", testConsulConfigEntryServiceRouterEE, "ent", "", false, excludeMissingEnterpriseScopes},
		{"ConsulConfigEntryServiceRouterEE_noMatch", testConsulConfigEntryServiceRouterEE_noMatch, "ent", "", false, excludeMissingEnterpriseScopes},
		{"ConsulConfigEntryServiceRouterEE_noDestination", testConsulConfigEntryServiceRouterEE_noDestination, "ent", "", false, excludeMissingEnterpriseScopes},
		{"ConsulConfigEntryServiceSplitterCE", testConsulConfigEntryServiceSplitter("", ""), "ce", "", false, included},
		{"ConsulConfigEntryServiceSplitterEE", testConsulConfigEntryServiceSplitter("default", "default"), "ent", "", false, included},
		{"ConsulExportedServicesV2Basic", testAccConsulExportedServicesV2Basic, "ent", "1.18.0", false, included},
		{"ConsulNamespaceExportedServicesV2Basic", testAccConsulNamespaceExportedServicesV2Basic, "ent", "1.18.0", false, included},
		{"ConsulPartitionExportedServicesV2Basic", testAccConsulPartitionExportedServicesV2Basic, "ent", "1.18.0", false, included},
		{"ConsulIntentionConfigBasic", testAccConsulIntentionConfigBasic, "", "", false, included},
		{"ConsulIntentionConfigBadAction", testAccConsulIntentionConfigBadAction, "", "", false, excludeApplyFails},
		{"ConsulIntentionConfigNamespaceCE", testAccConsulIntentionConfigNamespaceCE, "", "", false, excludeApplyFails},
		{"ConsulIntentionConfigNamespaceEE", testAccConsulIntentionConfigNamespaceEE, "ent", "", false, included},
		{"ConsulKeyPrefixNoKeys", testAccConsulKeyPrefixNoKeys, "", "", false, included},
		{"ConsulKeyPrefixConfig", testAccConsulKeyPrefixConfig, "", "", false, included},
		{"ConsulKeyPrefixConfig_Update", testAccConsulKeyPrefixConfig_Update, "", "", false, included},
		{"ConsulKeyPrefixConfig_root", testAccConsulKeyPrefixConfig_root, "", "", false, included},
		{"ConsulKeyPrefixConfig_Import", testAccConsulKeyPrefixConfig_Import, "", "", false, included},
		{"ConsulKeyPrefixConfig_namespaceCE", testAccConsulKeyPrefixConfig_namespaceCE, "", "", false, excludeApplyFails},
		{"ConsulKeyPrefixConfig_namespaceEE", testAccConsulKeyPrefixConfig_namespaceEE, "ent", "", false, included},
		{"ConsulKeyPrefixConfig_deleted", testAccConsulKeyPrefixConfig_deleted, "", "", false, included},
		{"ConsulKeyPrefixConfig_datacenter", testAccConsulKeyPrefixConfig_datacenter, "", "", true, included},
		{"ConsulKeysConfig", testAccConsulKeysConfig, "", "", false, included},
		{"ConsulKeysConfig_Update", testAccConsulKeysConfig_Update, "", "", false, included},
		{"ConsulKeysEmptyValue", testAccConsulKeysEmptyValue, "", "", false, included},
		{"ConsulKeysNamespaceCE", testAccConsulKeysNamespaceCE, "", "", false, excludeApplyFails},
		{"ConsulKeysNamespaceEE", testAccConsulKeysNamespaceEE, "ent", "", false, included},
		{"ConsulKeysDatacenter", testAccConsulKeysDatacenter, "", "", true, included},
		{"ConsulLicense", testAccConsulLicense, "", "", false, excludeApplyFails},
		{"NamespacePolicyConfig", testResourceNamespacePolicyConfig, "ent", "", false, included},
		{"NamespacePolicyConfigUpdate", testResourceNamespacePolicyConfigUpdate, "ent", "", false, included},
		{"NamespaceRoleConfig", testResourceNamespaceRoleConfig, "ent", "", false, included},
		{"NamespaceRoleConfigUpdate", testResourceNamespaceRoleConfigUpdate, "ent", "", false, included},
		{"ConsulNamespace", testAccConsulNamespace, "ent", "", false, included},
		{"ConsulNamespace_Update", testAccConsulNamespace_Update, "ent", "", false, included},
		{"ConsulNetworkAreaBasic", testAccConsulNetworkAreaBasic, "ent", "", false, included},
		{"ConsulNetworkAreaBasic_update1", testAccConsulNetworkAreaBasic_update1, "ent", "", false, included},
		{"ConsulNetworkAreaBasic_update2", testAccConsulNetworkAreaBasic_update2, "ent", "", false, included},
		{"ConsulNetworkAreaDatacenter", testAccConsulNetworkAreaDatacenter, "ent", "", true, included},
		{"ConsulNodeConfigBasic", testAccConsulNodeConfigBasic, "", "", false, included},
		{"ConsulNodeConfigNodeMeta", testAccConsulNodeConfigNodeMeta, "", "", false, included},
		{"ConsulNodeConfigNodeMeta_Update", testAccConsulNodeConfigNodeMeta_Update, "", "", false, included},
		{"ConsulNodeConfigDatacenter", testAccConsulNodeConfigDatacenter, "", "", true, included},
		{"ConsulPeeringBasic", testAccConsulPeeringBasic, "", "", false, excludeSecondProviderName},
		{"ConsulPeeringTokenBasic", testAccConsulPeeringTokenBasic, "", "", false, included},
		{"AcConsulPeeringTokenServerExternalAddresses", testAcConsulPeeringTokenServerExternalAddresses, "", "", false, included},
		{"ConsulPreparedQueryConfig", testAccConsulPreparedQueryConfig, "", "", false, included},
		{"ConsulPreparedQueryConfigUpdate1", testAccConsulPreparedQueryConfigUpdate1, "", "", false, included},
		{"ConsulPreparedQueryConfigUpdate2", testAccConsulPreparedQueryConfigUpdate2, "", "", false, included},
		{"ConsulPreparedQueryBlocks", testAccConsulPreparedQueryBlocks, "", "", false, included},
		{"ConsulPreparedQueryBlocks2", testAccConsulPreparedQueryBlocks2, "", "", false, included},
		{"ConsulPreparedQueryBlocks3", testAccConsulPreparedQueryBlocks3, "", "", false, included},
		{"ConsulPreparedQueryBlocks4", testAccConsulPreparedQueryBlocks4, "", "", false, included},
		{"ConsulPreparedQueryBlocks5", testAccConsulPreparedQueryBlocks5, "", "", false, included},
		{"ConsulPreparedQueryDatacenter", testAccConsulPreparedQueryDatacenter, "", "", true, included},
		{"ConsulServiceConfigBasic", testAccConsulServiceConfigBasic, "", "", false, included},
		{"ConsulServiceConfigBasicMeta", testAccConsulServiceConfigBasicMeta, "", "", false, included},
		{"ConsulServiceConfigBasicNewTags", testAccConsulServiceConfigBasicNewTags, "", "", false, included},
		{"ConsulServiceConfigBasicAddress", testAccConsulServiceConfigBasicAddress, "", "", false, included},
		{"ConsulServiceConfigServiceID", testAccConsulServiceConfigServiceID, "", "", false, included},
		{"ConsulServiceCheck", testAccConsulServiceCheck, "", "", false, included},
		{"ConsulServiceCheckID", testAccConsulServiceCheckID, "", "", false, included},
		{"ConsulServiceCheckOrder", testAccConsulServiceCheckOrder, "", "", false, included},
		{"DataConsulServiceSameServiceMultipleNodes", testAccDataConsulServiceSameServiceMultipleNodes, "", "", false, included},
		{"ConsulServiceConfigNoNode", testAccConsulServiceConfigNoNode, "", "", false, excludeApplyFails},
		{"ConsulServiceDontOverrideNodeMeta", testAccConsulServiceDontOverrideNodeMeta, "", "", false, included},
		{"ConsulServiceMultipleInstances", testAccConsulServiceMultipleInstances, "", "", false, included},
		{"ConsulServiceNamespaceCE", testAccConsulServiceNamespaceCE, "", "", false, excludeApplyFails},
		{"ConsulServiceNamespaceEE", testAccConsulServiceNamespaceEE, "ent", "", false, included},
		{"ConsulServiceDatacenter", testAccConsulServiceDatacenter, "", "", true, included},
		{"ConsulServiceConfigBasicWeights", testAccConsulServiceConfigBasicWeights, "", "", false, included},
		{"ConsulServiceConfigBasicWeightsOptional", testAccConsulServiceConfigBasicWeightsOptional, "", "", false, included},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.exclusion != included {
				t.Skip(string(tt.exclusion))
			}

			start := startTestServer
			if tt.remoteDC {
				start = startRemoteDatacenterTestServer
			}
			providers, client := start(t)

			switch tt.edition {
			case "ce":
				skipTestOnConsulEnterpriseEdition(t)
			case "ent":
				skipTestOnConsulCommunityEdition(t)
			}
			if tt.minVersion != "" {
				if skip, err := skipIfConsulVersionLT(client, tt.minVersion)(); err != nil {
					t.Fatal(err)
				} else if skip {
					t.Skipf("requires Consul %s or later", tt.minVersion)
				}
			}

			resource.Test(t, resource.TestCase{
				Steps: []resource.TestStep{
					{
						ExternalProviders: map[string]resource.ExternalProvider{
							"consul": {Source: "hashicorp/consul", VersionConstraint: baseline},
						},
						Config: tt.config,
					},
					{
						ProviderFactories: providers,
						Config:            tt.config,
						PlanOnly:          true,
					},
				},
			})
		})
	}
}
