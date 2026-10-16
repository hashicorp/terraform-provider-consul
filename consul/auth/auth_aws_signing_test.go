// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
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
	srv := httptest.NewServer(fakeSTSHandler(t, arn, status, &count))
	t.Cleanup(srv.Close)
	return srv, &count
}

func fakeSTSHandler(t *testing.T, arn string, status int, count *int32) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(count, 1)
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
	})
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

			data, err := generateConsulAWSLoginData(creds, nil, "us-east-1", srv.URL, "consul.example.com")
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

	if _, err := generateConsulAWSLoginData(creds, nil, "us-east-1", srv.URL, ""); err == nil {
		t.Fatal("expected an error when STS rejects GetCallerIdentity")
	}
}

// isolateAWSEnv clears AWS settings so the developer's environment and ~/.aws do not affect the test.
func isolateAWSEnv(t *testing.T) {
	t.Helper()
	missing := filepath.Join(t.TempDir(), "missing")
	t.Setenv("AWS_CONFIG_FILE", missing)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", missing)
	for _, v := range []string{
		"AWS_PROFILE", "AWS_DEFAULT_PROFILE", "AWS_REGION", "AWS_DEFAULT_REGION",
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN",
		"AWS_CA_BUNDLE", "AWS_USE_FIPS_ENDPOINT", "AWS_USE_DUALSTACK_ENDPOINT", "AWS_ENDPOINT_URL", "AWS_ENDPOINT_URL_STS", "AWS_ENDPOINT_URL_IAM",
	} {
		t.Setenv(v, "")
	}
}

func staticCredsParams(region string) map[string]interface{} {
	return map[string]interface{}{
		utils.FieldAWSAccessKeyID:     "AKIDTEST",
		utils.FieldAWSSecretAccessKey: "test-secret-key",
		utils.FieldAWSRegion:          region,
	}
}

// TestGenerateConsulAWSLoginData_CABundle verifies AWS_CA_BUNDLE is honored for the STS GetCallerIdentity lookup,
// e.g. when STS is reached through a TLS-intercepting proxy with a private CA.
func TestGenerateConsulAWSLoginData_CABundle(t *testing.T) {
	isolateAWSEnv(t)
	var count int32
	srv := httptest.NewTLSServer(fakeSTSHandler(t, "arn:aws:iam::123456789012:user/test-user", http.StatusOK, &count))
	t.Cleanup(srv.Close)

	caFile := filepath.Join(t.TempDir(), "ca.pem")
	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw})
	if err := os.WriteFile(caFile, caPEM, 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("without AWS_CA_BUNDLE the private CA is not trusted", func(t *testing.T) {
		creds, awsCfg, err := resolveAWSCredentials(context.Background(), staticCredsParams("us-east-1"))
		if err != nil {
			t.Fatalf("resolveAWSCredentials() error: %v", err)
		}
		_, err = generateConsulAWSLoginData(creds, awsCfg, "us-east-1", srv.URL, "")
		if err == nil || !strings.Contains(err.Error(), "certificate") {
			t.Fatalf("expected a certificate error, got: %v", err)
		}
	})

	t.Run("with AWS_CA_BUNDLE login data is generated", func(t *testing.T) {
		t.Setenv("AWS_CA_BUNDLE", caFile)
		creds, awsCfg, err := resolveAWSCredentials(context.Background(), staticCredsParams("us-east-1"))
		if err != nil {
			t.Fatalf("resolveAWSCredentials() error: %v", err)
		}
		before := atomic.LoadInt32(&count)
		if _, err := generateConsulAWSLoginData(creds, awsCfg, "us-east-1", srv.URL, ""); err != nil {
			t.Fatalf("generateConsulAWSLoginData() error: %v", err)
		}
		if got := atomic.LoadInt32(&count) - before; got != 1 {
			t.Errorf("expected 1 request to STS, got %d", got)
		}
	})
}

// TestGenerateConsulAWSLoginData_FIPS verifies AWS_USE_FIPS_ENDPOINT produces FIPS STS/IAM requests,
// with the IAM request still signed for us-east-1 (IAM is global).
func TestGenerateConsulAWSLoginData_FIPS(t *testing.T) {
	isolateAWSEnv(t)
	t.Setenv("AWS_USE_FIPS_ENDPOINT", "true")

	creds, awsCfg, err := resolveAWSCredentials(context.Background(), staticCredsParams("us-west-2"))
	if err != nil {
		t.Fatalf("resolveAWSCredentials() error: %v", err)
	}
	if awsCfg == nil {
		t.Fatal("expected an AWS config to be loaded")
	}

	// Answer the GetCallerIdentity lookup locally instead of calling AWS.
	var requestedHosts []string
	awsCfg.HTTPClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		requestedHosts = append(requestedHosts, r.URL.Host)
		body := fmt.Sprintf(testCallerIdentityResponse, "arn:aws:sts::123456789012:assumed-role/test-role/test-session")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/xml"}},
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Request:    r,
		}, nil
	})}

	data, err := generateConsulAWSLoginData(creds, awsCfg, "us-west-2", "", "")
	if err != nil {
		t.Fatalf("generateConsulAWSLoginData() error: %v", err)
	}

	if len(requestedHosts) != 1 || requestedHosts[0] != "sts-fips.us-west-2.amazonaws.com" {
		t.Errorf("GetCallerIdentity lookup sent to %v, want [sts-fips.us-west-2.amazonaws.com]", requestedHosts)
	}
	if reqURL := decodeLoginField(t, data, "iam_request_url"); !strings.HasPrefix(reqURL, "https://sts-fips.us-west-2.amazonaws.com") {
		t.Errorf("iam_request_url = %q, want sts-fips.us-west-2.amazonaws.com", reqURL)
	}

	var raw map[string][]string
	if err := json.Unmarshal([]byte(decodeLoginField(t, data, "iam_request_headers")), &raw); err != nil {
		t.Fatalf("iam_request_headers is not JSON: %v", err)
	}
	headers := http.Header(raw)
	if auth := headers.Get("Authorization"); !strings.Contains(auth, "/us-west-2/sts/aws4_request") {
		t.Errorf("STS request not signed for us-west-2: %q", auth)
	}
	if got := headers.Get(utils.GetEntityURLHeaderName); !strings.HasPrefix(got, "https://iam-fips.amazonaws.com") {
		t.Errorf("%s = %q, want iam-fips.amazonaws.com", utils.GetEntityURLHeaderName, got)
	}
	var entityHeaders http.Header
	if err := json.Unmarshal([]byte(headers.Get(utils.GetEntityHeadersHeaderName)), &entityHeaders); err != nil {
		t.Fatalf("%s is not JSON: %v", utils.GetEntityHeadersHeaderName, err)
	}
	if auth := entityHeaders.Get("Authorization"); !strings.Contains(auth, "/us-east-1/iam/aws4_request") {
		t.Errorf("IAM request not signed for us-east-1: %q", auth)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// TestResolveAWSCredentials_StaticCredsConfigLoadFailure verifies static credentials keep working
// (as before) when only the shared-config profile is missing, and that a warning is logged.
func TestResolveAWSCredentials_StaticCredsConfigLoadFailure(t *testing.T) {
	isolateAWSEnv(t)
	t.Setenv("AWS_PROFILE", "does-not-exist")

	var logs bytes.Buffer
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	creds, awsCfg, err := resolveAWSCredentials(context.Background(), staticCredsParams("us-east-1"))
	if err != nil {
		t.Fatalf("resolveAWSCredentials() error: %v", err)
	}
	if creds.AccessKeyID != "AKIDTEST" {
		t.Errorf("AccessKeyID = %q, want AKIDTEST", creds.AccessKeyID)
	}
	if awsCfg != nil {
		t.Errorf("expected no AWS config when loading fails, got %+v", awsCfg)
	}
	if !strings.Contains(logs.String(), "[WARN] Failed to load AWS config") {
		t.Errorf("expected a [WARN] log, got: %q", logs.String())
	}
}

// TestResolveAWSCredentials_StaticCredsSecuritySettingsFailClosed verifies that, with static credentials,
// explicitly configured CA bundle / FIPS settings are never silently dropped when the AWS config fails to load.
func TestResolveAWSCredentials_StaticCredsSecuritySettingsFailClosed(t *testing.T) {
	missingFile := filepath.Join(t.TempDir(), "missing-ca.pem")

	tests := []struct {
		name string
		env  map[string]string
	}{
		{"unreadable CA bundle and FIPS", map[string]string{"AWS_CA_BUNDLE": missingFile, "AWS_USE_FIPS_ENDPOINT": "true"}},
		{"unreadable CA bundle", map[string]string{"AWS_CA_BUNDLE": missingFile}},
		{"invalid FIPS value", map[string]string{"AWS_USE_FIPS_ENDPOINT": "not-a-bool"}},
		{"missing profile with FIPS", map[string]string{"AWS_PROFILE": "does-not-exist", "AWS_USE_FIPS_ENDPOINT": "true"}},
		{"missing profile with CA bundle", map[string]string{"AWS_PROFILE": "does-not-exist", "AWS_CA_BUNDLE": missingFile}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateAWSEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			_, awsCfg, err := resolveAWSCredentials(context.Background(), staticCredsParams("us-east-1"))
			if err == nil {
				t.Fatalf("expected an error, got config %+v", awsCfg)
			}
			if !strings.Contains(err.Error(), "failed to load AWS config") {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}

	// A missing source_profile surfaces as a wrapped "profile not exist" error, but the selected
	// profile exists and its use_fips_endpoint must not be silently dropped.
	t.Run("profile FIPS with missing source profile", func(t *testing.T) {
		isolateAWSEnv(t)
		configFile := filepath.Join(t.TempDir(), "config")
		profile := "[profile fips]\nuse_fips_endpoint = true\nrole_arn = arn:aws:iam::123456789012:role/test-role\nsource_profile = does-not-exist\n"
		if err := os.WriteFile(configFile, []byte(profile), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("AWS_CONFIG_FILE", configFile)
		t.Setenv("AWS_PROFILE", "fips")

		_, awsCfg, err := resolveAWSCredentials(context.Background(), staticCredsParams("us-east-1"))
		if err == nil {
			t.Fatalf("expected an error, got config %+v", awsCfg)
		}
		if !strings.Contains(err.Error(), "failed to load AWS config") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
