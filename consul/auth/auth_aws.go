// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"context"
	"encoding/json"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-consul/consul/auth/utils"
)

func init() {
	field := utils.FieldAuthLoginAWS
	if err := globalAuthLoginRegistry.Register(field,
		func(r *schema.ResourceData) (AuthLogin, error) {
			a := &AuthLoginAWS{}
			return a.Init(r, field)
		}, GetAWSLoginSchema); err != nil {
		panic(err)
	}
}

// GetAWSLoginSchema for the AWS authentication engine.
func GetAWSLoginSchema(authField string) *schema.Schema {
	return getLoginSchema(
		authField,
		"Login to Consul using the AWS IAM auth method",
		GetAWSLoginSchemaResource,
	)
}

// GetAWSLoginSchemaResource for the AWS authentication engine.
func GetAWSLoginSchemaResource(authField string) *schema.Resource {
	return mustAddLoginSchema(&schema.Resource{
		Schema: map[string]*schema.Schema{
			utils.FieldAuthMethod: {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the Consul auth method to use for login.`,
			},
			utils.FieldBearerToken: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Pre-computed bearer token to use for login. If not provided, the provider will generate one using AWS credentials.`,
			},
			// static credential fields
			utils.FieldAWSAccessKeyID: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The AWS access key ID.`,
			},
			utils.FieldAWSSecretAccessKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The AWS secret access key.`,
			},
			utils.FieldAWSSessionToken: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The AWS session token.`,
			},
			utils.FieldAWSProfile: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the AWS profile.`,
			},
			utils.FieldAWSSharedCredentialsFile: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Path to the AWS shared credentials file.`,
			},
			utils.FieldAWSWebIdentityTokenFile: {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Path to the file containing an OAuth 2.0 access token or OpenID ` +
					`Connect ID token.`,
			},
			// STS assume role fields
			utils.FieldAWSRoleARN: {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The ARN of the AWS Role to assume. ` +
					`Used during STS AssumeRole`,
			},
			utils.FieldAWSRoleSessionName: {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the name to attach to the AWS role session. ` +
					`Used during STS AssumeRole`,
			},
			utils.FieldAWSRegion: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The AWS region.`,
			},
			utils.FieldAWSSTSEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The STS endpoint URL.`,
			},
			utils.FieldAWSIAMEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The IAM endpoint URL.`,
			},
			utils.FieldServerIDHeaderValue: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Consul Server ID header value to include in the STS signing request. This must match the ServerIDHeaderValue configured in the Consul auth method.`,
			},
			utils.FieldMeta: {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: `Specifies arbitrary KV metadata linked to the token. Can be useful to track origins.`,
			},
		},
	}, authField)
}

var _ AuthLogin = (*AuthLoginAWS)(nil)

// AuthLoginAWS for handling the Consul AWS IAM authentication method.
// Requires configuration provided by SchemaLoginAWS.
type AuthLoginAWS struct {
	AuthLoginCommon
}

func (l *AuthLoginAWS) Init(d *schema.ResourceData, authField string) (AuthLogin, error) {
	defaults := l.getDefaults()
	if err := l.AuthLoginCommon.Init(d, authField,
		func(data *schema.ResourceData, params map[string]interface{}) error {
			return l.setDefaultFields(d, defaults, params)
		},
		func(data *schema.ResourceData, params map[string]interface{}) error {
			return l.checkRequiredFields(d, params, utils.FieldAuthMethod)
		},
	); err != nil {
		return nil, err
	}

	return l, nil
}

// AuthMethodName returns the Consul auth method name.
func (l *AuthLoginAWS) AuthMethodName() string {
	if v, ok := l.params[utils.FieldAuthMethod].(string); ok {
		return v
	}
	return ""
}

// Login using the AWS IAM authentication method.
func (l *AuthLoginAWS) Login(client *consulapi.Client) (string, error) {
	if err := l.validate(); err != nil {
		return "", err
	}

	authMethod := l.AuthMethodName()
	if authMethod == "" {
		return "", fmt.Errorf("auth_method is required")
	}

	// Check if bearer token is directly provided
	var bearerToken string
	if v, ok := l.params[utils.FieldBearerToken].(string); ok && v != "" {
		bearerToken = v
	} else {
		// Generate bearer token from AWS credentials
		loginData, err := l.generateAWSLoginData()
		if err != nil {
			return "", fmt.Errorf("failed to generate AWS login data: %w", err)
		}

		// Encode the login data as the bearer token
		// Consul expects the JSON string directly (not base64-encoded)
		jsonData, err := json.Marshal(loginData)
		if err != nil {
			return "", fmt.Errorf("failed to marshal login data: %w", err)
		}
		bearerToken = string(jsonData)
	}

	// Extract metadata if provided
	meta := make(map[string]string)
	if metaRaw, ok := l.params[utils.FieldMeta].(map[string]interface{}); ok {
		for k, v := range metaRaw {
			if strVal, ok := v.(string); ok {
				meta[k] = strVal
			}
		}
	}

	// Use standard login - the X-Consul-IAM-GetEntity-Method header is already
	// included in the signed AWS request (in the bearer token)
	return l.login(client, authMethod, bearerToken, meta)
}

func (l *AuthLoginAWS) getDefaults() authDefaults {
	defaults := authDefaults{
		{
			field:      utils.FieldAWSAccessKeyID,
			envVars:    []string{utils.EnvVarAWSAccessKeyID},
			defaultVal: "",
		},
		{
			field:      utils.FieldAWSSecretAccessKey,
			envVars:    []string{utils.EnvVarAWSSecretAccessKey},
			defaultVal: "",
		},
		{
			field:      utils.FieldAWSSessionToken,
			envVars:    []string{utils.EnvVarAWSSessionToken},
			defaultVal: "",
		},
		{
			field:      utils.FieldAWSProfile,
			envVars:    []string{utils.EnvVarAWSProfile},
			defaultVal: "",
		},
		{
			field:      utils.FieldAWSSharedCredentialsFile,
			envVars:    []string{utils.EnvVarAWSSharedCredentialsFile},
			defaultVal: "",
		},
		{
			field:      utils.FieldAWSWebIdentityTokenFile,
			envVars:    []string{utils.EnvVarAWSWebIdentityTokenFile},
			defaultVal: "",
		},
		{
			field:      utils.FieldAWSRoleARN,
			envVars:    []string{utils.EnvVarAWSRoleARN},
			defaultVal: "",
		},
		{
			field:      utils.FieldAWSRoleSessionName,
			envVars:    []string{utils.EnvVarAWSRoleSessionName},
			defaultVal: "",
		},
		{
			field:      utils.FieldAWSRegion,
			envVars:    []string{utils.EnvVarAWSRegion, utils.EnvVarAWSDefaultRegion},
			defaultVal: "",
		},
	}

	return defaults
}

// generateAWSLoginData creates login data from configured AWS parameters.
// This generates a properly signed STS GetCallerIdentity request for Consul.
func (l *AuthLoginAWS) generateAWSLoginData() (map[string]interface{}, error) {
	ctx := context.Background()

	// Resolve AWS credentials using aws-sdk-go-v2 credential chain
	creds, err := resolveAWSCredentials(ctx, l.params)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve AWS credentials: %w", err)
	}

	// Get region
	region, _ := l.params[utils.FieldAWSRegion].(string)
	if region == "" {
		region = utils.DefaultAWSRegion
	}

	// Get server ID header value (optional)
	serverIDHeaderValue, _ := l.params[utils.FieldServerIDHeaderValue].(string)

	// Generate Consul-specific AWS login data
	loginData, err := generateConsulAWSLoginData(creds, region, serverIDHeaderValue)
	if err != nil {
		return nil, fmt.Errorf("failed to generate login data: %w", err)
	}

	return loginData, nil
}

// getCredentialsConfig builds an awsutil.CredentialsConfig from the configured parameters.
// This handles all AWS credential sources: static credentials, profiles, instance metadata,
// web identity tokens, assume role, etc.
