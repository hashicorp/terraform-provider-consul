// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
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
						fieldHeaderValue:              "header1",
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
				fieldHeaderValue:              "header1",
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
			name:      "error-missing-required",
			authField: fieldAuthLoginAWS,
			raw: map[string]interface{}{
				fieldAuthLoginAWS: []interface{}{
					map[string]interface{}{},
				},
			},
			expectParams: nil,
			wantErr:      true,
			expectErr: fmt.Errorf("required fields are unset: %v", []string{
				fieldAuthMethod,
			}),
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
			expectParams: map[string]interface{}{
				fieldAuthMethod:         "aws-auth",
				fieldAWSAccessKeyID:     "env-key-id",
				fieldAWSSecretAccessKey: "env-sa-key",
				fieldAWSRegion:          "us-west-2",
				// These should be empty strings from defaults
				fieldAWSSessionToken:          "",
				fieldAWSProfile:               "",
				fieldAWSSharedCredentialsFile: "",
				fieldAWSWebIdentityTokenFile:  "",
				fieldAWSRoleARN:               "",
				fieldAWSRoleSessionName:       "",
			},
			wantErr: false,
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

func TestAuthLoginAWS_generateAWSLoginData(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]interface{}
		wantErr bool
		wantLen int
	}{
		{
			name: "with-credentials",
			params: map[string]interface{}{
				fieldAWSAccessKeyID:     "key-id",
				fieldAWSSecretAccessKey: "sa-key",
				fieldAWSRegion:          "us-east-1",
			},
			wantErr: false,
			wantLen: 3,
		},
		{
			name: "with-all-fields",
			params: map[string]interface{}{
				fieldAWSAccessKeyID:     "key-id",
				fieldAWSSecretAccessKey: "sa-key",
				fieldAWSSessionToken:    "session-token",
				fieldAWSRegion:          "us-east-2",
				fieldHeaderValue:        "header1",
				fieldAWSRoleARN:         "arn:aws:iam::123456789012:role/test",
			},
			wantErr: false,
			wantLen: 6,
		},
		{
			name:    "no-credentials",
			params:  map[string]interface{}{},
			wantErr: true,
			wantLen: 0,
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
			got, err := l.generateAWSLoginData()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateAWSLoginData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("generateAWSLoginData() returned %d fields, want %d", len(got), tt.wantLen)
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
			want:           "test-secret-token-12345",
			expectReqCount: 1,
			wantErr:        false,
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
				handlerFunc: func(t *testLoginHandler, w http.ResponseWriter, req *http.Request) {
					// Should not be called
					t.Errorf("handler should not be called without auth method")
				},
			},
			want:           "",
			expectReqCount: 0,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAuthLogin(t, tt)
		})
	}
}

func TestAuthLoginAWS_Schema(t *testing.T) {
	schema := GetAWSLoginSchema(fieldAuthLoginAWS)

	if schema == nil {
		t.Fatal("GetAWSLoginSchema() returned nil")
	}

	if schema.Type != schema.TypeList {
		t.Errorf("expected TypeList, got %v", schema.Type)
	}

	if schema.MaxItems != 1 {
		t.Errorf("expected MaxItems=1, got %d", schema.MaxItems)
	}

	if schema.Optional != true {
		t.Error("expected Optional=true")
	}

	resource := schema.Elem.(*schema.Resource)
	if resource == nil {
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
		fieldHeaderValue,
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

func TestAuthLoginAWS_ConflictsWith(t *testing.T) {
	schema := GetAWSLoginSchema(fieldAuthLoginAWS)

	if len(schema.ConflictsWith) == 0 {
		t.Error("expected ConflictsWith to be populated")
	}

	// The schema should conflict with other auth methods
	// Since we only have AWS registered in tests, there might be no conflicts yet
	// But the mechanism should be in place
	t.Logf("ConflictsWith: %v", schema.ConflictsWith)
}
