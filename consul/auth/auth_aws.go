// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"encoding/json"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	envVarAWSAccessKeyID           = "AWS_ACCESS_KEY_ID"
	envVarAWSSecretAccessKey       = "AWS_SECRET_ACCESS_KEY"
	envVarAWSSessionToken          = "AWS_SESSION_TOKEN"
	envVarAWSProfile               = "AWS_PROFILE"
	envVarAWSSharedCredentialsFile = "AWS_SHARED_CREDENTIALS_FILE"
	envVarAWSWebIdentityTokenFile  = "AWS_WEB_IDENTITY_TOKEN_FILE"
	envVarAWSRoleARN               = "AWS_ROLE_ARN"
	envVarAWSRoleSessionName       = "AWS_ROLE_SESSION_NAME"
	envVarAWSRegion                = "AWS_REGION"
	envVarAWSDefaultRegion         = "AWS_DEFAULT_REGION"

	fieldAuthLoginAWS             = "auth_login_aws"
	fieldAuthMethod               = "auth_method"
	fieldAWSAccessKeyID           = "aws_access_key_id"
	fieldAWSSecretAccessKey       = "aws_secret_access_key"
	fieldAWSSessionToken          = "aws_session_token"
	fieldAWSProfile               = "aws_profile"
	fieldAWSSharedCredentialsFile = "aws_shared_credentials_file"
	fieldAWSWebIdentityTokenFile  = "aws_web_identity_token_file"
	fieldAWSRoleARN               = "aws_role_arn"
	fieldAWSRoleSessionName       = "aws_role_session_name"
	fieldAWSRegion                = "aws_region"
	fieldAWSSTSEndpoint           = "aws_sts_endpoint"
	fieldAWSIAMEndpoint           = "aws_iam_endpoint"
	fieldServerIDHeaderValue      = "server_id_header_value"
	fieldMeta                     = "meta"
	fieldBearerToken              = "bearer_token"
)

func init() {
	field := fieldAuthLoginAWS
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
			fieldAuthMethod: {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the Consul auth method to use for login.`,
			},
			fieldBearerToken: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Pre-computed bearer token to use for login. If not provided, the provider will generate one using AWS credentials.`,
			},
			// static credential fields
			fieldAWSAccessKeyID: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The AWS access key ID.`,
			},
			fieldAWSSecretAccessKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The AWS secret access key.`,
			},
			fieldAWSSessionToken: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The AWS session token.`,
			},
			fieldAWSProfile: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the AWS profile.`,
			},
			fieldAWSSharedCredentialsFile: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Path to the AWS shared credentials file.`,
			},
			fieldAWSWebIdentityTokenFile: {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Path to the file containing an OAuth 2.0 access token or OpenID ` +
					`Connect ID token.`,
			},
			// STS assume role fields
			fieldAWSRoleARN: {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The ARN of the AWS Role to assume. ` +
					`Used during STS AssumeRole`,
			},
			fieldAWSRoleSessionName: {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the name to attach to the AWS role session. ` +
					`Used during STS AssumeRole`,
			},
			fieldAWSRegion: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The AWS region.`,
			},
			fieldAWSSTSEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The STS endpoint URL.`,
			},
			fieldAWSIAMEndpoint: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The IAM endpoint URL.`,
			},
			fieldServerIDHeaderValue: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Consul Server ID header value to include in the STS signing request. This must match the ServerIDHeaderValue configured in the Consul auth method.`,
			},
			fieldMeta: {
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
			return l.checkRequiredFields(d, params, fieldAuthMethod)
		},
	); err != nil {
		return nil, err
	}

	return l, nil
}

// AuthMethodName returns the Consul auth method name.
func (l *AuthLoginAWS) AuthMethodName() string {
	if v, ok := l.params[fieldAuthMethod].(string); ok {
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
	if v, ok := l.params[fieldBearerToken].(string); ok && v != "" {
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
	if metaRaw, ok := l.params[fieldMeta].(map[string]interface{}); ok {
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
			field:      fieldAWSAccessKeyID,
			envVars:    []string{envVarAWSAccessKeyID},
			defaultVal: "",
		},
		{
			field:      fieldAWSSecretAccessKey,
			envVars:    []string{envVarAWSSecretAccessKey},
			defaultVal: "",
		},
		{
			field:      fieldAWSSessionToken,
			envVars:    []string{envVarAWSSessionToken},
			defaultVal: "",
		},
		{
			field:      fieldAWSProfile,
			envVars:    []string{envVarAWSProfile},
			defaultVal: "",
		},
		{
			field:      fieldAWSSharedCredentialsFile,
			envVars:    []string{envVarAWSSharedCredentialsFile},
			defaultVal: "",
		},
		{
			field:      fieldAWSWebIdentityTokenFile,
			envVars:    []string{envVarAWSWebIdentityTokenFile},
			defaultVal: "",
		},
		{
			field:      fieldAWSRoleARN,
			envVars:    []string{envVarAWSRoleARN},
			defaultVal: "",
		},
		{
			field:      fieldAWSRoleSessionName,
			envVars:    []string{envVarAWSRoleSessionName},
			defaultVal: "",
		},
		{
			field:      fieldAWSRegion,
			envVars:    []string{envVarAWSRegion, envVarAWSDefaultRegion},
			defaultVal: "",
		},
	}

	return defaults
}

// generateAWSLoginData creates login data from configured AWS parameters.
// This generates a properly signed STS GetCallerIdentity request for Consul.
func (l *AuthLoginAWS) generateAWSLoginData() (map[string]interface{}, error) {
	// Build credentials config from parameters
	config, err := l.getCredentialsConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials config: %w", err)
	}

	// Generate credential chain (handles profiles, instance metadata, env vars, etc.)
	creds, err := config.GenerateCredentialChain()
	if err != nil {
		return nil, fmt.Errorf("failed to generate credential chain: %w", err)
	}

	// Get header value for X-Consul-IAM-ServerID (optional)
	// If provided, this must match the ServerIDHeaderValue configured in the Consul auth method
	// If not provided, no ServerID validation will be performed by Consul
	var serverIDHeaderValue string
	if v, ok := l.params[fieldServerIDHeaderValue].(string); ok {
		serverIDHeaderValue = v
	}

	// Use our custom Consul-specific signing function instead of awsutil.GenerateLoginData
	// This ensures the X-Consul-IAM-GetEntity-Method header is included BEFORE signing
	region := config.Region
	if region == "" {
		region = "us-east-1"
	}

	loginData, err := generateConsulAWSLoginData(creds, region, serverIDHeaderValue)
	if err != nil {
		return nil, fmt.Errorf("failed to generate login data: %w", err)
	}

	return loginData, nil
}

// getCredentialsConfig builds an awsutil.CredentialsConfig from the configured parameters.
// This handles all AWS credential sources: static credentials, profiles, instance metadata,
// web identity tokens, assume role, etc.
func (l *AuthLoginAWS) getCredentialsConfig() (*awsutil.CredentialsConfig, error) {
	config, err := awsutil.NewCredentialsConfig()
	if err != nil {
		return nil, err
	}

	// Map all our parameters to awsutil config
	if v, ok := l.params[fieldAWSAccessKeyID].(string); ok && v != "" {
		config.AccessKey = v
	}
	if v, ok := l.params[fieldAWSSecretAccessKey].(string); ok && v != "" {
		config.SecretKey = v
	}
	if v, ok := l.params[fieldAWSSessionToken].(string); ok && v != "" {
		config.SessionToken = v
	}
	if v, ok := l.params[fieldAWSProfile].(string); ok && v != "" {
		config.Profile = v
	}
	if v, ok := l.params[fieldAWSSharedCredentialsFile].(string); ok && v != "" {
		config.Filename = v
	}
	if v, ok := l.params[fieldAWSWebIdentityTokenFile].(string); ok && v != "" {
		config.WebIdentityTokenFile = v
	}
	if v, ok := l.params[fieldAWSRoleARN].(string); ok && v != "" {
		config.RoleARN = v
	}
	if v, ok := l.params[fieldAWSRoleSessionName].(string); ok && v != "" {
		config.RoleSessionName = v
	}
	if v, ok := l.params[fieldAWSRegion].(string); ok && v != "" {
		config.Region = v
	}
	if v, ok := l.params[fieldAWSSTSEndpoint].(string); ok && v != "" {
		config.STSEndpoint = v
	}
	if v, ok := l.params[fieldAWSIAMEndpoint].(string); ok && v != "" {
		config.IAMEndpoint = v
	}

	config.Logger = hclog.NewNullLogger()

	return config, nil
}
