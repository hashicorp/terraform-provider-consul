package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	iamauth "github.com/hashicorp/consul-awsauth"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-provider-consul/consul/auth/utils"
)

// resolveAWSCredentials resolves AWS credentials using aws-sdk-go-v2 credential chain
// It handles: static credentials, environment variables, profiles, web identity tokens, instance metadata, and assume role
func resolveAWSCredentials(ctx context.Context, params map[string]interface{}) (aws.Credentials, error) {
	// Check for static credentials first
	accessKey := utils.GetStringParam(params, utils.FieldAWSAccessKeyID)
	secretKey := utils.GetStringParam(params, utils.FieldAWSSecretAccessKey)
	sessionToken := utils.GetStringParam(params, utils.FieldAWSSessionToken)

	// If both access key and secret key are provided, use static credentials
	if utils.IsNonEmpty(accessKey) && utils.IsNonEmpty(secretKey) {
		return aws.Credentials{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
			SessionToken:    sessionToken,
		}, nil
	}

	// Get region from params, environment variables, or use default
	region := utils.GetStringParam(params, utils.FieldAWSRegion)
	if !utils.IsNonEmpty(region) {
		region = utils.GetStringEnv(utils.EnvVarAWSRegion, utils.EnvVarAWSDefaultRegion)
	}
	if !utils.IsNonEmpty(region) {
		region = utils.DefaultAWSRegion
	}

	// Build configuration using aws-sdk-go-v2
	cfgOpts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	// Handle profile if specified
	profile := utils.GetStringParam(params, utils.FieldAWSProfile)
	if utils.IsNonEmpty(profile) {
		cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
	}

	// Handle shared credentials file if specified
	credFile := utils.GetStringParam(params, utils.FieldAWSSharedCredentialsFile)
	if utils.IsNonEmpty(credFile) {
		cfgOpts = append(cfgOpts, config.WithSharedCredentialsFiles([]string{credFile}))
	}

	// Get role ARN for STS assume role
	roleARN := utils.GetStringParam(params, utils.FieldAWSRoleARN)

	// Load base configuration
	cfg, err := config.LoadDefaultConfig(ctx, cfgOpts...)
	if err != nil {
		return aws.Credentials{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Web identity token provider is automatically handled by config.LoadDefaultConfig
	// using AWS_WEB_IDENTITY_TOKEN_FILE and AWS_ROLE_ARN environment variables
	// No additional handling needed - fall through to default chain

	// Handle STS assume role if specified
	if utils.IsNonEmpty(roleARN) {
		stsClient := sts.NewFromConfig(cfg)
		assumeRoleProvider := stscreds.NewAssumeRoleProvider(stsClient, roleARN)

		creds, err := assumeRoleProvider.Retrieve(ctx)
		if err == nil {
			return creds, nil
		}
		// Fall through to default credential chain if assume role fails
	}

	// Use the default credential chain from the config
	// This will try: env vars, profile, instance metadata, ECS container credentials, etc.
	// The config.LoadDefaultConfig handles EC2 instance metadata provider automatically
	return cfg.Credentials.Retrieve(ctx)
}

// generateConsulAWSLoginData creates Consul-specific AWS login data
// Following the Consul CLI implementation exactly
func generateConsulAWSLoginData(creds aws.Credentials, region string, serverIDHeaderValue string) (map[string]interface{}, error) {
	// Create LoginInput for Consul AWS authentication
	// The library's LoginInput.Creds field accepts interface{} which can be
	// either SDK v1 *credentials.Credentials or SDK v2 aws.Credentials
	// We pass v2 credentials directly - the library will access the fields through interface{}

	loginInput := &iamauth.LoginInput{
		// Pass aws.Credentials struct (not pointer) as interface{}
		// The library can access AccessKeyID, SecretAccessKey, SessionToken fields through interface{}
		Creds:                  utils.ToV1Credentials(creds),
		IncludeIAMEntity:       utils.DefaultIncludeIAMEntity, // Required for binding rules with entity metadata
		STSRegion:              region,
		Logger:                 hclog.NewNullLogger(),
		ServerIDHeaderValue:    serverIDHeaderValue,
		ServerIDHeaderName:     utils.IAMServerIDHeaderName,      // Header name: "X-Consul-IAM-ServerID"
		GetEntityMethodHeader:  utils.GetEntityMethodHeaderName,  // Header name: "X-Consul-IAM-GetEntity-Method"
		GetEntityURLHeader:     utils.GetEntityURLHeaderName,     // Header name: "X-Consul-IAM-GetEntity-URL"
		GetEntityHeadersHeader: utils.GetEntityHeadersHeaderName, // Header name: "X-Consul-IAM-GetEntity-Headers"
		GetEntityBodyHeader:    utils.GetEntityBodyHeaderName,    // Header name: "X-Consul-IAM-GetEntity-Body"
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
