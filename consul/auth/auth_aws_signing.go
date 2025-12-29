package auth

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	iamauth "github.com/hashicorp/consul-awsauth"
	"github.com/hashicorp/go-hclog"
)

// AWS IAM auth method header names - must match Consul's expectations
const (
	IAMServerIDHeaderName  string = "X-Consul-IAM-ServerID"
	GetEntityMethodHeader  string = "X-Consul-IAM-GetEntity-Method"
	GetEntityURLHeader     string = "X-Consul-IAM-GetEntity-URL"
	GetEntityHeadersHeader string = "X-Consul-IAM-GetEntity-Headers"
	GetEntityBodyHeader    string = "X-Consul-IAM-GetEntity-Body"
)

// generateConsulAWSLoginData creates Consul-specific AWS login data
// Following the Consul CLI implementation exactly
func generateConsulAWSLoginData(creds *credentials.Credentials, region string, serverIDHeaderValue string) (map[string]interface{}, error) {
	// Create LoginInput exactly as Consul CLI does
	// The GetEntityXXXHeader parameters are HEADER NAMES, not values
	// The library will use these names when creating the signed request
	loginInput := &iamauth.LoginInput{
		Creds:                  creds,
		IncludeIAMEntity:       true, // Required for binding rules with entity metadata
		STSRegion:              region,
		STSEndpoint:            "", // Use default
		Logger:                 hclog.NewNullLogger(),
		ServerIDHeaderValue:    serverIDHeaderValue,
		ServerIDHeaderName:     IAMServerIDHeaderName,  // Header name: "X-Consul-IAM-ServerID"
		GetEntityMethodHeader:  GetEntityMethodHeader,  // Header name: "X-Consul-IAM-GetEntity-Method"
		GetEntityURLHeader:     GetEntityURLHeader,     // Header name: "X-Consul-IAM-GetEntity-URL"
		GetEntityHeadersHeader: GetEntityHeadersHeader, // Header name: "X-Consul-IAM-GetEntity-Headers"
		GetEntityBodyHeader:    GetEntityBodyHeader,    // Header name: "X-Consul-IAM-GetEntity-Body"
	}

	// Generate login data using the library's function
	loginData, err := iamauth.GenerateLoginData(loginInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AWS login data: %w", err)
	}

	// Convert to map[string]interface{} for our interface
	var result map[string]interface{}
	jsonBytes, err := json.Marshal(loginData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login data: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal login data: %w", err)
	}

	return result, nil
}
