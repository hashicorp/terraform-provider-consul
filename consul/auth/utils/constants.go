// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

// AWS Authentication Configuration (aws-auth group)

// Environment Variables used for AWS authentication
const (
	// EnvVarAWSAccessKeyID is the environment variable for AWS access key ID
	EnvVarAWSAccessKeyID = "AWS_ACCESS_KEY_ID"
	// EnvVarAWSSecretAccessKey is the environment variable for AWS secret access key
	EnvVarAWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	// EnvVarAWSSessionToken is the environment variable for AWS session token
	EnvVarAWSSessionToken = "AWS_SESSION_TOKEN"
	// EnvVarAWSProfile is the environment variable for AWS profile
	EnvVarAWSProfile = "AWS_PROFILE"
	// EnvVarAWSSharedCredentialsFile is the environment variable for AWS shared credentials file
	EnvVarAWSSharedCredentialsFile = "AWS_SHARED_CREDENTIALS_FILE"
	// EnvVarAWSWebIdentityTokenFile is the environment variable for AWS web identity token file
	EnvVarAWSWebIdentityTokenFile = "AWS_WEB_IDENTITY_TOKEN_FILE"
	// EnvVarAWSRoleARN is the environment variable for AWS role ARN
	EnvVarAWSRoleARN = "AWS_ROLE_ARN"
	// EnvVarAWSRoleSessionName is the environment variable for AWS role session name
	EnvVarAWSRoleSessionName = "AWS_ROLE_SESSION_NAME"
	// EnvVarAWSRegion is the environment variable for AWS region
	EnvVarAWSRegion = "AWS_REGION"
	// EnvVarAWSDefaultRegion is the environment variable for AWS default region
	EnvVarAWSDefaultRegion = "AWS_DEFAULT_REGION"
	// EnvVarBearerToken is the environment variable for bearer token (used across all auth methods)
	EnvVarBearerToken = "CONSUL_LOGIN_BEARER_TOKEN"
)

// Schema Field Names for AWS authentication configuration
const (
	// FieldAuthLoginAWS is the field name for AWS authentication login
	FieldAuthLoginAWS = "auth_login_aws"
	// FieldAuthMethod is the field name for authentication method
	FieldAuthMethod = "auth_method"
	// FieldAWSAccessKeyID is the field name for AWS access key ID
	FieldAWSAccessKeyID = "aws_access_key_id"
	// FieldAWSSecretAccessKey is the field name for AWS secret access key
	FieldAWSSecretAccessKey = "aws_secret_access_key"
	// FieldAWSSessionToken is the field name for AWS session token
	FieldAWSSessionToken = "aws_session_token"
	// FieldAWSProfile is the field name for AWS profile
	FieldAWSProfile = "aws_profile"
	// FieldAWSSharedCredentialsFile is the field name for AWS shared credentials file
	FieldAWSSharedCredentialsFile = "aws_shared_credentials_file"
	// FieldAWSWebIdentityTokenFile is the field name for AWS web identity token file
	FieldAWSWebIdentityTokenFile = "aws_web_identity_token_file"
	// FieldAWSRoleARN is the field name for AWS role ARN
	FieldAWSRoleARN = "aws_role_arn"
	// FieldAWSRoleSessionName is the field name for AWS role session name
	FieldAWSRoleSessionName = "aws_role_session_name"
	// FieldAWSRegion is the field name for AWS region
	FieldAWSRegion = "aws_region"
	// FieldAWSSTSEndpoint is the field name for AWS STS endpoint
	FieldAWSSTSEndpoint = "aws_sts_endpoint"
	// FieldAWSIAMEndpoint is the field name for AWS IAM endpoint
	FieldAWSIAMEndpoint = "aws_iam_endpoint"
	// FieldServerIDHeaderValue is the field name for server ID header value
	FieldServerIDHeaderValue = "server_id_header_value"
	// FieldMeta is the field name for metadata
	FieldMeta = "meta"
	// FieldBearerToken is the field name for bearer token
	FieldBearerToken = "bearer_token"
)

// Default AWS Region used when no region is specified
const DefaultAWSRegion = "us-east-1"

// AWS STS Service Configuration
const (
	// STSGetCallerIdentityAction is the STS action for getting caller identity
	STSGetCallerIdentityAction = "sts:GetCallerIdentity"
	// IAMGetUserAction is the IAM action for getting user information
	IAMGetUserAction = "iam:GetUser"
)

// AWS IAM Authentication Header Names - must match Consul's expectations
// These headers are used in the signed AWS request for Consul authentication
const (
	// IAMServerIDHeaderName is the header name for Consul's Server ID verification
	IAMServerIDHeaderName = "X-Consul-IAM-ServerID"
	// GetEntityMethodHeaderName is the header name for the HTTP method in IAM entity resolution
	GetEntityMethodHeaderName = "X-Consul-IAM-GetEntity-Method"
	// GetEntityURLHeaderName is the header name for the URL in IAM entity resolution
	GetEntityURLHeaderName = "X-Consul-IAM-GetEntity-URL"
	// GetEntityHeadersHeaderName is the header name for request headers in IAM entity resolution
	GetEntityHeadersHeaderName = "X-Consul-IAM-GetEntity-Headers"
	// GetEntityBodyHeaderName is the header name for request body in IAM entity resolution
	GetEntityBodyHeaderName = "X-Consul-IAM-GetEntity-Body"
)

// AWS HTTP Methods used in signing
const (
	// HTTPMethodGET is the GET HTTP method
	HTTPMethodGET = "GET"
	// HTTPMethodPOST is the POST HTTP method
	HTTPMethodPOST = "POST"
)

// AWS SigV4 Signing Configuration
const (
	// AWSSignatureVersion is the AWS signature version used
	AWSSignatureVersion = "AWS4-HMAC-SHA256"
	// AWSServiceName is the AWS service name for SigV4 signing context
	AWSServiceName = "sts"
	// AWSRequestType is the type of AWS request
	AWSRequestType = "aws4_request"
)

// Context and Metadata Configuration
const (
	// DefaultIncludeIAMEntity determines if IAM entity metadata should be included by default
	// Set to true for binding rules that require entity information
	DefaultIncludeIAMEntity = true
)

// Consul Authentication Configuration
const (
	// ConsulAuthMethod is the authentication method name in Consul
	ConsulAuthMethod = "aws-iam"
	// ConsulAuthPath is the default authentication path
	ConsulAuthPath = "auth/aws-iam/login"
)

// HTTP Configuration
const (
	// HTTPContentTypeJSON is the JSON content type
	HTTPContentTypeJSON = "application/json"
	// HTTPContentTypeForm is the form content type
	HTTPContentTypeForm = "application/x-www-form-urlencoded"
)

// Error Messages and Descriptions
const (
	// ErrorCredentialsRequired is returned when credentials are required but not provided
	ErrorCredentialsRequired = "credentials are required"
	// ErrorLoginInputNil is returned when login input is nil
	ErrorLoginInputNil = "domain login input cannot be nil"
	// ErrorFailedToResolveCredentials is returned when credential resolution fails
	ErrorFailedToResolveCredentials = "failed to resolve AWS credentials"
	// ErrorFailedToGenerateLoginData is returned when login data generation fails
	ErrorFailedToGenerateLoginData = "failed to generate AWS login data"
	// ErrorFailedToMarshalLoginData is returned when marshaling login data fails
	ErrorFailedToMarshalLoginData = "failed to marshal login data"
	// ErrorFailedToUnmarshalLoginData is returned when unmarshaling login data fails
	ErrorFailedToUnmarshalLoginData = "failed to unmarshal login data"
)
