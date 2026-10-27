// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestSetDeprecatedConfig covers configurations the deprecated 'config' map of strings cannot hold. TF_ACC is set
// because SDK v2 panics on invalid d.Set calls only then, which is how acceptance tests would hit it.
func TestSetDeprecatedConfig(t *testing.T) {
	t.Setenv("TF_ACC", "1")

	schemas := map[string]map[string]*schema.Schema{
		"consul_acl_auth_method":             resourceConsulACLAuthMethod().Schema,
		"consul_acl_auth_method data source": dataSourceConsulACLAuthMethod().Schema,
		"consul_certificate_authority":       resourceConsulCertificateAuthority().Schema,
	}
	configs := map[string]map[string]interface{}{
		"bool":   {"VerboseOIDCLogging": true},
		"number": {"PrivateKeyBits": float64(2048)},
		"nested": {"ClaimMappings": map[string]interface{}{"first_name": "first_name"}},
		"list":   {"BoundAudiences": []interface{}{"a", "b"}},
	}

	for schemaName, s := range schemas {
		for configName, config := range configs {
			t.Run(schemaName+"/"+configName, func(t *testing.T) {
				d := schema.TestResourceDataRaw(t, s, map[string]interface{}{})
				if err := d.Set("config_json", `{"set":"by the caller"}`); err != nil {
					t.Fatal(err)
				}
				sw := newStateWriter(d)
				if err := sw.setDeprecatedConfig(config); err != nil {
					t.Fatalf("with config_json set: unexpected error: %v", err)
				}
				if err := sw.error(); err != nil {
					t.Fatalf("with config_json set: unexpected state error: %v", err)
				}
				if got := d.Get("config").(map[string]interface{}); len(got) != 0 {
					t.Errorf("config = %v, want an empty placeholder", got)
				}

				d = schema.TestResourceDataRaw(t, s, map[string]interface{}{})
				err := newStateWriter(d).setDeprecatedConfig(config)
				if err == nil || !strings.Contains(err.Error(), "use 'config_json' instead") {
					t.Fatalf("without config_json: error = %v, want a hint to use config_json", err)
				}
			})
		}

		t.Run(schemaName+"/strings", func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, map[string]interface{}{})
			if err := newStateWriter(d).setDeprecatedConfig(map[string]interface{}{"Host": "https://example.com", "Empty": nil}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := d.Get("config").(map[string]interface{}); got["Host"] != "https://example.com" {
				t.Errorf("config = %v, want Host to be set", got)
			}
		})
	}
}
