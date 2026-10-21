// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// newTestPGPKey returns a new PGP entity and its base64-encoded public key, the format users pass as pgp_key.
func newTestPGPKey(t *testing.T, cfg *packet.Config) (*openpgp.Entity, string) {
	t.Helper()
	entity, err := openpgp.NewEntity("test", "", "test@example.com", cfg)
	if err != nil {
		t.Fatalf("generating key: %v", err)
	}
	var pub bytes.Buffer
	if err := entity.Serialize(&pub); err != nil {
		t.Fatalf("serializing public key: %v", err)
	}
	return entity, base64.StdEncoding.EncodeToString(pub.Bytes())
}

func TestEncryptValueWithPGPKey(t *testing.T) {
	const secret = "6b4f9a2c-1d3e-4f5a-8b7c-9d0e1f2a3b4c"

	tests := []struct {
		name string
		cfg  *packet.Config
	}{
		{"rsa-2048", &packet.Config{Algorithm: packet.PubKeyAlgoRSA, RSABits: 2048}},
		{"rsa-4096", &packet.Config{Algorithm: packet.PubKeyAlgoRSA, RSABits: 4096}},
		{"ed25519-cv25519", &packet.Config{Algorithm: packet.PubKeyAlgoEdDSA, Curve: packet.Curve25519}},
		// Legacy key size; documents whether a library change starts rejecting existing user keys.
		{"rsa-1024", &packet.Config{Algorithm: packet.PubKeyAlgoRSA, RSABits: 1024}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, pgpKey := newTestPGPKey(t, tt.cfg)

			encrypted, err := encryptValueWithPGPKey(pgpKey, secret, "ACL Token secret ID")
			if err != nil {
				t.Fatalf("encryptValueWithPGPKey() error: %v", err)
			}

			raw, err := base64.StdEncoding.DecodeString(encrypted)
			if err != nil {
				t.Fatalf("encrypted value is not base64: %v", err)
			}
			// Users decode with `base64 -d | gpg -d`, so the payload must be a binary (not armored) OpenPGP message.
			if strings.HasPrefix(string(raw), "-----BEGIN") || len(raw) == 0 || raw[0]&0x80 == 0 {
				t.Fatalf("encrypted value is not a binary OpenPGP message: %q", raw[:min(len(raw), 16)])
			}

			md, err := openpgp.ReadMessage(bytes.NewReader(raw), openpgp.EntityList{entity}, nil, nil)
			if err != nil {
				t.Fatalf("decrypting: %v", err)
			}
			got, err := io.ReadAll(md.UnverifiedBody)
			if err != nil {
				t.Fatalf("reading decrypted body: %v", err)
			}
			if string(got) != secret {
				t.Fatalf("decrypted value = %q, want %q", got, secret)
			}
			assertGPGCompatibleMessage(t, raw)
		})
	}
}

// assertGPGCompatibleMessage checks the message uses the classic format every gpg version decrypts:
// a public-key encrypted session key followed by SEIPDv1 (MDC, no AEAD).
func assertGPGCompatibleMessage(t *testing.T, raw []byte) {
	t.Helper()
	r := packet.NewReader(bytes.NewReader(raw))

	p, err := r.Next()
	if err != nil {
		t.Fatalf("reading first packet: %v", err)
	}
	ek, ok := p.(*packet.EncryptedKey)
	if !ok {
		t.Fatalf("first packet is %T, want *packet.EncryptedKey", p)
	}
	if ek.Version != 3 {
		t.Errorf("encrypted session key packet version = %d, want 3", ek.Version)
	}

	p, err = r.Next()
	if err != nil {
		t.Fatalf("reading second packet: %v", err)
	}
	se, ok := p.(*packet.SymmetricallyEncrypted)
	if !ok {
		t.Fatalf("second packet is %T, want *packet.SymmetricallyEncrypted", p)
	}
	if se.Version != 1 || !se.IntegrityProtected {
		t.Errorf("encrypted data packet: version=%d integrityProtected=%v, want SEIPDv1 with MDC", se.Version, se.IntegrityProtected)
	}
}

// newFakeKeybase serves the keybase user lookup API with the given response and records the request query.
func newFakeKeybase(t *testing.T, response interface{}) *http.Request {
	t.Helper()
	var got http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = *r
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	t.Cleanup(srv.Close)

	orig := keybaseLookupURL
	keybaseLookupURL = srv.URL + "/_/api/1.0/user/lookup.json"
	t.Cleanup(func() { keybaseLookupURL = orig })
	return &got
}

func keybaseResponse(status, bundle string) map[string]interface{} {
	return map[string]interface{}{
		"status": map[string]interface{}{"name": status},
		"them": []interface{}{
			map[string]interface{}{"public_keys": map[string]interface{}{"primary": map[string]interface{}{"bundle": bundle}}},
		},
	}
}

func armoredPublicKey(t *testing.T, entity *openpgp.Entity) string {
	t.Helper()
	var buf bytes.Buffer
	w, err := armor.Encode(&buf, openpgp.PublicKeyType, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := entity.Serialize(w); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestEncryptValueWithPGPKey_Keybase(t *testing.T) {
	const secret = "6b4f9a2c-1d3e-4f5a-8b7c-9d0e1f2a3b4c"
	entity, _ := newTestPGPKey(t, &packet.Config{Algorithm: packet.PubKeyAlgoRSA, RSABits: 2048})
	req := newFakeKeybase(t, keybaseResponse("OK", armoredPublicKey(t, entity)))

	encrypted, err := encryptValueWithPGPKey("keybase:test_user", secret, "ACL Token secret ID")
	if err != nil {
		t.Fatalf("encryptValueWithPGPKey() error: %v", err)
	}

	if got := req.URL.Query().Get("usernames"); got != "test_user" {
		t.Errorf("usernames = %q, want test_user", got)
	}
	if got := req.URL.Query().Get("fields"); got != "public_keys" {
		t.Errorf("fields = %q, want public_keys", got)
	}

	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		t.Fatalf("encrypted value is not base64: %v", err)
	}
	md, err := openpgp.ReadMessage(bytes.NewReader(raw), openpgp.EntityList{entity}, nil, nil)
	if err != nil {
		t.Fatalf("decrypting: %v", err)
	}
	got, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		t.Fatalf("reading decrypted body: %v", err)
	}
	if string(got) != secret {
		t.Fatalf("decrypted value = %q, want %q", got, secret)
	}
}

func TestEncryptValueWithPGPKey_KeybaseErrors(t *testing.T) {
	tests := []struct {
		name     string
		response interface{}
		wantErr  string
	}{
		{"non-OK status", keybaseResponse("NOT_FOUND", ""), `got non-OK response: "NOT_FOUND"`},
		{"user without key", keybaseResponse("OK", ""), "unable to fetch keys"},
		{"no user", map[string]interface{}{"status": map[string]interface{}{"name": "OK"}, "them": []interface{}{}}, "unable to fetch keys"},
		{"invalid bundle", keybaseResponse("OK", "not a key"), "Error retrieving Public Key for keybase:test_user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newFakeKeybase(t, tt.response)
			_, err := encryptValueWithPGPKey("keybase:test_user", "secret", "test")
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

// TestEncryptValueWithPGPKey_KeybaseUsernameEscaped verifies the username cannot inject extra query parameters.
func TestEncryptValueWithPGPKey_KeybaseUsernameEscaped(t *testing.T) {
	req := newFakeKeybase(t, keybaseResponse("NOT_FOUND", ""))
	_, _ = encryptValueWithPGPKey("keybase:a&fields=x", "secret", "test")

	if got := req.URL.Query().Get("usernames"); got != "a&fields=x" {
		t.Errorf("usernames = %q, want %q", got, "a&fields=x")
	}
	if got := req.URL.Query()["fields"]; len(got) != 1 || got[0] != "public_keys" {
		t.Errorf("fields = %q, want [public_keys]", got)
	}
}

func TestEncryptValueWithPGPKey_InvalidKey(t *testing.T) {
	tests := []struct {
		name   string
		pgpKey string
	}{
		{"not-base64", "not a valid key!"},
		{"base64-non-key", base64.StdEncoding.EncodeToString([]byte("hello world"))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := encryptValueWithPGPKey(tt.pgpKey, "secret", "test"); err == nil {
				t.Fatal("expected an error for an invalid PGP key")
			}
		})
	}
}
