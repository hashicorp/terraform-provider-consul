//go:build integration
// +build integration

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"testing"
)

// TestAuthLoginAWS_generateAWSLoginData_Integration requires real AWS credentials
// Run with: go test -tags=integration ./consul/auth -v -run TestAuthLoginAWS_generateAWSLoginData_Integration
func TestAuthLoginAWS_generateAWSLoginData_Integration(t *testing.T) {
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
				fieldAWSAccessKeyID:      "key-id",
				fieldAWSSecretAccessKey:  "sa-key",
				fieldAWSSessionToken:     "session-token",
				fieldAWSRegion:           "us-east-2",
				fieldServerIDHeaderValue: "header1",
				fieldAWSRoleARN:          "arn:aws:iam::123456789012:role/test",
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
