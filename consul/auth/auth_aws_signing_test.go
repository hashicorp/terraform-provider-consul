// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-provider-consul/consul/auth/utils"
)

const testCallerIdentityResponse = `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
    <Arn>%s</Arn>
    <UserId>AIDATESTUSERID</UserId>
    <Account>123456789012</Account>
  </GetCallerIdentityResult>
  <ResponseMetadata>
    <RequestId>00000000-0000-0000-0000-000000000000</RequestId>
  </ResponseMetadata>
</GetCallerIdentityResponse>`

// newFakeSTSServer answers GetCallerIdentity with the given ARN and counts requests.
func newFakeSTSServer(t *testing.T, arn string, status int) (*httptest.Server, *int32) {
	t.Helper()
	var count int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&count, 1)
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "Action=GetCallerIdentity") {
			t.Errorf("unexpected STS request body: %s", body)
		}
		if !strings.Contains(r.Header.Get("Authorization"), "Credential=AKIDTEST/") {
			t.Errorf("STS request not signed with test credentials: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(status)
		if status == http.StatusOK {
			fmt.Fprintf(w, testCallerIdentityResponse, arn)
			return
		}
		fmt.Fprint(w, `<ErrorResponse><Error><Type>Sender</Type><Code>InvalidClientTokenId</Code><Message>invalid</Message></Error></ErrorResponse>`)
	}))
	t.Cleanup(srv.Close)
	return srv, &count
}

func decodeLoginField(t *testing.T, data map[string]interface{}, key string) string {
	t.Helper()
	v, ok := data[key].(string)
	if !ok {
		t.Fatalf("login data field %q missing or not a string: %#v", key, data[key])
	}
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		t.Fatalf("login data field %q is not base64: %v", key, err)
	}
	return string(b)
}

// TestGenerateConsulAWSLoginData verifies the login data the provider sends to Consul for AWS IAM auth.
//
// Login flow:
//
//	provider ──(1 real call)──► AWS STS GetCallerIdentity   → learns IAM user/role name
//	provider builds + signs STS GetCallerIdentity request   (not sent)
//	provider builds + signs IAM GetUser/GetRole request     (not sent)
//	provider ──login data──► Consul ──replays signed requests──► AWS
//
// The one real STS call is redirected to a local fake server (via stsEndpoint) that returns a test ARN,
// so the test runs offline. It checks that the login data contains a correctly signed STS request
// (method, URL, body, SigV4 Authorization, session token, server ID header included in the signature)
// and the signed IAM GetUser/GetRole entity request headers Consul uses for binding rules.
func TestGenerateConsulAWSLoginData(t *testing.T) {
	creds := aws.Credentials{
		AccessKeyID:     "AKIDTEST",
		SecretAccessKey: "test-secret-key",
		SessionToken:    "test-session-token",
	}

	tests := []struct {
		name         string
		arn          string
		wantAction   string
		wantEntityKV string
	}{
		{
			name:         "iam-user",
			arn:          "arn:aws:iam::123456789012:user/test-user",
			wantAction:   "Action=GetUser",
			wantEntityKV: "UserName=test-user",
		},
		{
			name:         "assumed-role",
			arn:          "arn:aws:sts::123456789012:assumed-role/test-role/test-session",
			wantAction:   "Action=GetRole",
			wantEntityKV: "RoleName=test-role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, count := newFakeSTSServer(t, tt.arn, http.StatusOK)

			data, err := generateConsulAWSLoginData(creds, "us-east-1", srv.URL, "consul.example.com")
			if err != nil {
				t.Fatalf("generateConsulAWSLoginData() error: %v", err)
			}
			// Only the GetCallerIdentity lookup is sent; the login request itself is just signed.
			if got := atomic.LoadInt32(count); got != 1 {
				t.Errorf("expected 1 request to STS, got %d", got)
			}

			if got := data["iam_http_request_method"]; got != http.MethodPost {
				t.Errorf("iam_http_request_method = %v, want POST", got)
			}

			reqURL := decodeLoginField(t, data, "iam_request_url")
			if !strings.HasPrefix(reqURL, srv.URL) {
				t.Errorf("iam_request_url = %q, want prefix %q", reqURL, srv.URL)
			}

			body, err := url.ParseQuery(decodeLoginField(t, data, "iam_request_body"))
			if err != nil {
				t.Fatalf("iam_request_body is not form encoded: %v", err)
			}
			if body.Get("Action") != "GetCallerIdentity" || body.Get("Version") != "2011-06-15" {
				t.Errorf("unexpected iam_request_body: %v", body)
			}

			var raw map[string][]string
			if err := json.Unmarshal([]byte(decodeLoginField(t, data, "iam_request_headers")), &raw); err != nil {
				t.Fatalf("iam_request_headers is not JSON: %v", err)
			}
			headers := http.Header{}
			for k, v := range raw {
				for _, s := range v {
					headers.Add(k, s)
				}
			}

			auth := headers.Get("Authorization")
			if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256 ") ||
				!strings.Contains(auth, "Credential=AKIDTEST/") ||
				!strings.Contains(auth, "/us-east-1/sts/aws4_request") {
				t.Errorf("unexpected Authorization header: %q", auth)
			}
			if !strings.Contains(strings.ToLower(auth), strings.ToLower(utils.IAMServerIDHeaderName)) {
				t.Errorf("server ID header is not part of the signature: %q", auth)
			}
			if got := headers.Get("X-Amz-Security-Token"); got != creds.SessionToken {
				t.Errorf("X-Amz-Security-Token = %q, want %q", got, creds.SessionToken)
			}
			if got := headers.Get(utils.IAMServerIDHeaderName); got != "consul.example.com" {
				t.Errorf("%s = %q, want %q", utils.IAMServerIDHeaderName, got, "consul.example.com")
			}

			if got := headers.Get(utils.GetEntityMethodHeaderName); got != http.MethodPost {
				t.Errorf("%s = %q, want POST", utils.GetEntityMethodHeaderName, got)
			}
			if got := headers.Get(utils.GetEntityURLHeaderName); !strings.Contains(got, "iam.amazonaws.com") {
				// TestGenerateConsulAWSLoginData_STSError verifies that a rejected GetCallerIdentity lookup
				// (e.g. invalid credentials) is returned as an error instead of producing login data.
				t.Errorf("%s = %q, want IAM endpoint", utils.GetEntityURLHeaderName, got)
			}
			entityBody := headers.Get(utils.GetEntityBodyHeaderName)
			if !strings.Contains(entityBody, tt.wantAction) || !strings.Contains(entityBody, tt.wantEntityKV) {
				t.Errorf("%s = %q, want %q and %q", utils.GetEntityBodyHeaderName, entityBody, tt.wantAction, tt.wantEntityKV)
			}
			if got := headers.Get(utils.GetEntityHeadersHeaderName); !strings.Contains(got, "AWS4-HMAC-SHA256") {
				t.Errorf("%s does not contain a signed Authorization header: %q", utils.GetEntityHeadersHeaderName, got)
			}
		})
	}
}

func TestGenerateConsulAWSLoginData_STSError(t *testing.T) {
	srv, _ := newFakeSTSServer(t, "", http.StatusForbidden)
	creds := aws.Credentials{AccessKeyID: "AKIDTEST", SecretAccessKey: "test-secret-key"}

	if _, err := generateConsulAWSLoginData(creds, "us-east-1", srv.URL, ""); err == nil {
		t.Fatal("expected an error when STS rejects GetCallerIdentity")
	}
}
