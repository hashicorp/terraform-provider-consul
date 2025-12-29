// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// expectedRegisteredAuthLogin value should be modified when adding
// registering/de-registering AuthLogin resources.
const expectedRegisteredAuthLogin = 1

type authLoginTest struct {
	name               string
	authLogin          AuthLogin
	handler            *testLoginHandler
	want               string // expected token
	expectReqCount     int
	skipCheckReqParams bool
	expectReqParams    []map[string]interface{}
	wantErr            bool
	expectErr          error
	skipFunc           func(t *testing.T)
	preLoginFunc       func(t *testing.T)
	token              string
}

type authLoginInitTest struct {
	name         string
	authField    string
	raw          map[string]interface{}
	wantErr      bool
	envVars      map[string]string
	expectParams map[string]interface{}
	expectErr    error
}

type testLoginHandler struct {
	requestCount  int
	params        []map[string]interface{}
	excludeParams []string
	handlerFunc   func(t *testLoginHandler, w http.ResponseWriter, req *http.Request)
}

func (t *testLoginHandler) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		t.requestCount++

		switch req.Method {
		case http.MethodPut, http.MethodPost:
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		b, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var params map[string]interface{}
		if len(b) > 0 {
			if err := json.Unmarshal(b, &params); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		for _, p := range t.excludeParams {
			delete(params, p)
		}

		if len(params) > 0 {
			t.params = append(t.params, params)
		}

		t.handlerFunc(t, w, req)
	}
}

func testAuthLogin(t *testing.T, tt authLoginTest) {
	t.Helper()

	if tt.skipFunc != nil {
		tt.skipFunc(t)
	}

	if tt.preLoginFunc != nil {
		tt.preLoginFunc(t)
	}

	config := &consulapi.Config{
		Address: testHTTPServer(t, tt.handler.handler()),
	}

	c, err := consulapi.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}

	// Note: Consul API doesn't have SetToken() method on client
	// Tokens are set via the config or ACL().Login()

	got, err := tt.authLogin.Login(c)
	if (err != nil) != tt.wantErr {
		t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
		return
	}

	if err != nil && tt.expectErr != nil {
		if !reflect.DeepEqual(tt.expectErr, err) {
			t.Errorf("Login() expected error %#v, actual %#v", tt.expectErr, err)
		}
	}

	if tt.expectReqCount != tt.handler.requestCount {
		t.Errorf("Login() expected %d requests, actual %d", tt.expectReqCount, tt.handler.requestCount)
	}

	if !tt.skipCheckReqParams && !reflect.DeepEqual(tt.expectReqParams, tt.handler.params) {
		t.Errorf("Login() request params do not match expected %#v, actual %#v", tt.expectReqParams,
			tt.handler.params)
	}

	if got != tt.want {
		t.Errorf("Login() got = %#v, want %#v", got, tt.want)
	}
}

// testHTTPServer creates a test HTTP server for testing
func testHTTPServer(t *testing.T, handler http.HandlerFunc) string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	server := &http.Server{
		Handler: handler,
	}

	go func() {
		if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	t.Cleanup(func() {
		_ = server.Close()
		_ = ln.Close()
	})

	return ln.Addr().String()
}

// TestMustAddAuthLoginSchema_registered is only meant to validate that all
// expected AuthLogin(s) are registered. The expected count of all registered
// entries should be modified when registering/de-registering AuthLogin
// resources.
func TestMustAddAuthLoginSchema_registered(t *testing.T) {
	tests := []struct {
		name string
		s    map[string]*schema.Schema
	}{
		{
			name: "checkRegistered",
			s:    make(map[string]*schema.Schema),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MustAddAuthLoginSchema(tt.s)
			actual := len(tt.s)
			if expectedRegisteredAuthLogin != actual {
				t.Errorf("expected %d schema entries, actual %d", expectedRegisteredAuthLogin, actual)
			}
		})
	}
}

func TestGetAuthLogin_registered(t *testing.T) {
	registeredAuthLogins := globalAuthLoginRegistry.Values()
	actualRegistered := len(registeredAuthLogins)
	if expectedRegisteredAuthLogin != actualRegistered {
		t.Fatalf("expected %d registered AuthLogin, actual %d", expectedRegisteredAuthLogin, actualRegistered)
	}

	s := map[string]*schema.Schema{}
	MustAddAuthLoginSchema(s)

	for _, entry := range registeredAuthLogins {
		field := entry.Field()
		t.Run(field, func(t *testing.T) {
			raw := map[string]interface{}{
				field: []interface{}{
					map[string]interface{}{},
				},
			}

			r := &schema.ResourceData{}
			if err := r.Set(field, raw[field]); err != nil {
				t.Fatal(err)
			}

			_, err := GetAuthLogin(r)
			if err == nil {
				t.Errorf("GetAuthLogin() expected error for incomplete config")
			}
		})
	}
}

func assertAuthLoginInit(t *testing.T, tt authLoginInitTest, s map[string]*schema.Schema, authLogin AuthLogin) {
	t.Helper()

	for k, v := range tt.envVars {
		if err := os.Setenv(k, v); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = os.Unsetenv(k)
		})
	}

	r := schema.TestResourceDataRaw(t, s, tt.raw)

	got, err := authLogin.Init(r, tt.authField)
	if (err != nil) != tt.wantErr {
		t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
		return
	}

	if err != nil {
		if tt.expectErr != nil {
			if err.Error() != tt.expectErr.Error() {
				t.Errorf("Init() expected error %#v, actual %#v", tt.expectErr, err)
			}
		}
		return
	}

	if !reflect.DeepEqual(tt.expectParams, got.Params()) {
		t.Errorf("Init() params do not match expected %#v, actual %#v", tt.expectParams, got.Params())
	}
}
