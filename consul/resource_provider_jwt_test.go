package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestGetOptions_TokenFallback(t *testing.T) {
	// mock ResourceData with no token override
	d := schema.TestResourceDataRaw(t, make(map[string]*schema.Schema), make(map[string]interface{}))

	// config with a global provider-level JWT token
	config := &Config{
		Token:      "jprovider-level-jwt-token",
		Datacenter: "dc1",
	}

	qOpts, wOpts := getOptions(d, config)

	if qOpts.Token != "jprovider-level-jwt-token" {
		t.Errorf("Expected QueryOptions.Token to be 'jprovider-level-jwt-token', got '%s'", qOpts.Token)
	}

	if wOpts.Token != "jprovider-level-jwt-token" {
		t.Errorf("Expected WriteOptions.Token to be 'jprovider-level-jwt-token', got '%s'", wOpts.Token)
	}
}

func TestGetOptions_ResourceTokenOverride(t *testing.T) {
	// mock ResourceData with a resource-level token override
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"token": {Type: schema.TypeString, Optional: true},
	}, map[string]any{
		"token": "resource-level-token",
	})

	config := &Config{
		Token:      "provider-level-jwt-token",
		Datacenter: "dc1",
	}

	qOpts, wOpts := getOptions(d, config)

	// Test that the resource-level token takes precedence
	if qOpts.Token != "resource-level-token" {
		t.Errorf("Expected QueryOptions.Token to be 'resource-level-token', got '%s'", qOpts.Token)
	}

	if wOpts.Token != "resource-level-token" {
		t.Errorf("Expected WriteOptions.Token to be 'resource-level-token', got '%s'", wOpts.Token)
	}
}
