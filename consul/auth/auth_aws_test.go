// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestAuthLoginAWS_Init(t *testing.T) {
	tests := []authLoginInitTest{
		{
			name:      "basic",
			authField: fieldAuthLoginAWS,
			raw: map[string]interface{}{
				fieldAuthLoginAWS: []interface{}{
					map[string]interface{}{
						"namespace":                   "ns1",
						"partition":                   "part1",
						fieldAuthMethod:               "aws-auth",
						fieldAWSAccessKeyID:           "key-id",
						fieldAWSSecretAccessKey:       "sa-key",
						fieldAWSSessionToken:          "session-token",
						fieldAWSIAMEndpoint:           "iam.us-east-2.amazonaws.com",
						fieldAWSSTSEndpoint:           "sts.us-east-2.amazonaws.com",
						fieldAWSRegion:                "us-east-2",
						fieldAWSSharedCredentialsFile: "credentials",
						fieldAWSProfile:               "profile1",
						fieldAWSRoleARN:               "role-arn",
						fieldAWSRoleSessionName:       "session1",
						fieldAWSWebIdentityTokenFile:  "web-token",
						fieldServerIDHeaderValue:      "header1",
					},
				},
			},
			expectParams: map[string]interface{}{
				"namespace":                   "ns1",
				"partition":                   "part1",
				fieldAuthMethod:               "aws-auth",
				fieldAWSAccessKeyID:           "key-id",
				fieldAWSSecretAccessKey:       "sa-key",
				fieldAWSSessionToken:          "session-token",
				fieldAWSIAMEndpoint:           "iam.us-east-2.amazonaws.com",
				fieldAWSSTSEndpoint:           "sts.us-east-2.amazonaws.com",
				fieldAWSRegion:                "us-east-2",
				fieldAWSSharedCredentialsFile: "credentials",
				fieldAWSProfile:               "profile1",
				fieldAWSRoleARN:               "role-arn",
				fieldAWSRoleSessionName:       "session1",
				fieldAWSWebIdentityTokenFile:  "web-token",
				fieldServerIDHeaderValue:      "header1",
				fieldBearerToken:              "",
				fieldMeta:                     map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:         "error-missing-resource",
			authField:    fieldAuthLoginAWS,
			expectParams: nil,
			wantErr:      true,
			expectErr:    fmt.Errorf("resource data missing field %q", fieldAuthLoginAWS),
		},
		{
			name:      "with-env-vars",
			authField: fieldAuthLoginAWS,
			raw: map[string]interface{}{
				fieldAuthLoginAWS: []interface{}{
					map[string]interface{}{
						fieldAuthMethod: "aws-auth",
					},
				},
			},
			envVars: map[string]string{
				envVarAWSAccessKeyID:     "env-key-id",
				envVarAWSSecretAccessKey: "env-sa-key",
				envVarAWSRegion:          "us-west-2",
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
				fieldAuthMethod: "aws-iam-auth",
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
					authField:   fieldAuthLoginAWS,
					initialized: true,
					params: map[string]interface{}{
						fieldAuthMethod:  "aws-auth",
						fieldBearerToken: "test-bearer-token",
					},
				},
			},
			handler: &testLoginHandler{
				handlerFunc: func(t *testLoginHandler, w http.ResponseWriter, req *http.Request) {
					// Mock successful response
					response := map[string]interface{}{
						"SecretID": "test-secret-token-12345",
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
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
					authField:   fieldAuthLoginAWS,
					initialized: true,
					params: map[string]interface{}{
						fieldBearerToken: "test-bearer-token",
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
	s := GetAWSLoginSchema(fieldAuthLoginAWS)

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

	requiredFields := []string{fieldAuthMethod}
	for _, field := range requiredFields {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("missing required field: %s", field)
		} else if !resource.Schema[field].Required {
			t.Errorf("field %s should be required", field)
		}
	}

	optionalFields := []string{
		fieldAWSAccessKeyID,
		fieldAWSSecretAccessKey,
		fieldAWSSessionToken,
		fieldAWSProfile,
		fieldAWSSharedCredentialsFile,
		fieldAWSWebIdentityTokenFile,
		fieldAWSRoleARN,
		fieldAWSRoleSessionName,
		fieldAWSRegion,
		fieldAWSSTSEndpoint,
		fieldAWSIAMEndpoint,
		fieldServerIDHeaderValue,
		fieldBearerToken,
		fieldMeta,
		"namespace",
		"partition",
	}
	for _, field := range optionalFields {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("missing optional field: %s", field)
		}
	}
}
