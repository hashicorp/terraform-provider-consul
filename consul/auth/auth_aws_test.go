// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-consul/consul/auth/utils"
)

func TestAuthLoginAWS_Init(t *testing.T) {
	tests := []authLoginInitTest{
		{
			name:      "basic",
			authField: utils.FieldAuthLoginAWS,
			raw: map[string]interface{}{
				utils.FieldAuthLoginAWS: []interface{}{
					map[string]interface{}{
						"namespace":                         "ns1",
						"partition":                         "part1",
						utils.FieldAuthMethod:               "aws-auth",
						utils.FieldAWSAccessKeyID:           "key-id",
						utils.FieldAWSSecretAccessKey:       "sa-key",
						utils.FieldAWSSessionToken:          "session-token",
						utils.FieldAWSIAMEndpoint:           "iam.us-east-2.amazonaws.com",
						utils.FieldAWSSTSEndpoint:           "sts.us-east-2.amazonaws.com",
						utils.FieldAWSRegion:                "us-east-2",
						utils.FieldAWSSharedCredentialsFile: "credentials",
						utils.FieldAWSProfile:               "profile1",
						utils.FieldAWSRoleARN:               "role-arn",
						utils.FieldAWSRoleSessionName:       "session1",
						utils.FieldAWSWebIdentityTokenFile:  "web-token",
						utils.FieldServerIDHeaderValue:      "header1",
					},
				},
			},
			expectParams: map[string]interface{}{
				"namespace":                         "ns1",
				"partition":                         "part1",
				utils.FieldAuthMethod:               "aws-auth",
				utils.FieldAWSAccessKeyID:           "key-id",
				utils.FieldAWSSecretAccessKey:       "sa-key",
				utils.FieldAWSSessionToken:          "session-token",
				utils.FieldAWSIAMEndpoint:           "iam.us-east-2.amazonaws.com",
				utils.FieldAWSSTSEndpoint:           "sts.us-east-2.amazonaws.com",
				utils.FieldAWSRegion:                "us-east-2",
				utils.FieldAWSSharedCredentialsFile: "credentials",
				utils.FieldAWSProfile:               "profile1",
				utils.FieldAWSRoleARN:               "role-arn",
				utils.FieldAWSRoleSessionName:       "session1",
				utils.FieldAWSWebIdentityTokenFile:  "web-token",
				utils.FieldServerIDHeaderValue:      "header1",
				utils.FieldBearerToken:              "",
				utils.FieldMeta:                     map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:         "error-missing-resource",
			authField:    utils.FieldAuthLoginAWS,
			expectParams: nil,
			wantErr:      true,
			expectErr:    fmt.Errorf("resource data missing field %q", utils.FieldAuthLoginAWS),
		},
		{
			name:      "with-env-vars",
			authField: utils.FieldAuthLoginAWS,
			raw: map[string]interface{}{
				utils.FieldAuthLoginAWS: []interface{}{
					map[string]interface{}{
						utils.FieldAuthMethod: "aws-auth",
					},
				},
			},
			envVars: map[string]string{
				utils.EnvVarAWSAccessKeyID:     "env-key-id",
				utils.EnvVarAWSSecretAccessKey: "env-sa-key",
				utils.EnvVarAWSRegion:          "us-west-2",
			},
			expectParams: nil, // Don't check params - AWS env vars from host will be included
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := map[string]*schema.Schema{
				tt.authField: GetAWSLoginSchema(tt.authField),
			}
			assertAuthLoginInit(t, tt, s, &AuthLoginAWS{})
		})
	}
}

func TestAuthLoginAWS_AuthMethodName(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		want   string
	}{
		{
			name: "with-auth-method",
			params: map[string]interface{}{
				utils.FieldAuthMethod: "aws-iam-auth",
			},
			want: "aws-iam-auth",
		},
		{
			name:   "without-auth-method",
			params: map[string]interface{}{},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &AuthLoginAWS{
				AuthLoginCommon: AuthLoginCommon{
					params:      tt.params,
					initialized: true,
				},
			}
			if got := l.AuthMethodName(); got != tt.want {
				t.Errorf("AuthMethodName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthLoginAWS_Login(t *testing.T) {
	tests := []authLoginTest{
		{
			name: "successful-login-with-bearer-token",
			authLogin: &AuthLoginAWS{
				AuthLoginCommon: AuthLoginCommon{
					authField:   utils.FieldAuthLoginAWS,
					initialized: true,
					params: map[string]interface{}{
						utils.FieldAuthMethod:  "aws-auth",
						utils.FieldBearerToken: "test-bearer-token",
					},
				},
			},
			handler: &testLoginHandler{
				handlerFunc: func(t *testLoginHandler, w http.ResponseWriter, req *http.Request) {
					// Mock successful response
					response := map[string]interface{}{
						"SecretID": "test-secret-token-12345",
					}
					w.Header().Set("Content-Type", utils.HTTPContentTypeJSON)
					json.NewEncoder(w).Encode(response)
					_ = json.NewEncoder(w).Encode(response)
				},
			},
			want:               "test-secret-token-12345",
			expectReqCount:     1,
			skipCheckReqParams: true, // Skip param check since AWS auth auto-generates token
			wantErr:            false,
		},
		{
			name: "error-no-auth-method",
			authLogin: &AuthLoginAWS{
				AuthLoginCommon: AuthLoginCommon{
					authField:   utils.FieldAuthLoginAWS,
					initialized: true,
					params: map[string]interface{}{
						utils.FieldBearerToken: "test-bearer-token",
					},
				},
			},
			handler: &testLoginHandler{
				handlerFunc: func(h *testLoginHandler, w http.ResponseWriter, req *http.Request) {
					// Should not be called - return error response
					w.WriteHeader(http.StatusBadRequest)
				},
			},
			want:               "",
			expectReqCount:     0,
			skipCheckReqParams: true,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAuthLogin(t, tt)
		})
	}
}

func TestAuthLoginAWS_Schema(t *testing.T) {
	s := GetAWSLoginSchema(utils.FieldAuthLoginAWS)

	if s == nil {
		t.Fatal("GetAWSLoginSchema() returned nil")
	}

	if s.Type != schema.TypeList {
		t.Errorf("expected TypeList, got %v", s.Type)
	}

	if s.MaxItems != 1 {
		t.Errorf("expected MaxItems=1, got %d", s.MaxItems)
	}

	if s.Optional != true {
		t.Error("expected Optional=true")
	}

	resource, ok := s.Elem.(*schema.Resource)
	if !ok {
		t.Fatal("schema.Elem is not a Resource")
	}

	requiredFields := []string{utils.FieldAuthMethod}
	for _, field := range requiredFields {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("missing required field: %s", field)
		} else if !resource.Schema[field].Required {
			t.Errorf("field %s should be required", field)
		}
	}

	optionalFields := []string{
		utils.FieldAWSAccessKeyID,
		utils.FieldAWSSecretAccessKey,
		utils.FieldAWSSessionToken,
		utils.FieldAWSProfile,
		utils.FieldAWSSharedCredentialsFile,
		utils.FieldAWSWebIdentityTokenFile,
		utils.FieldAWSRoleARN,
		utils.FieldAWSRoleSessionName,
		utils.FieldAWSRegion,
		utils.FieldAWSSTSEndpoint,
		utils.FieldAWSIAMEndpoint,
		utils.FieldServerIDHeaderValue,
		utils.FieldMeta,
		"namespace",
		"partition",
	}
	for _, field := range optionalFields {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("missing optional field: %s", field)
		}
	}
}
