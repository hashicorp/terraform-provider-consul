// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"bytes"
	"encoding/base64"
	"io"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/encryption"
)

// TestEncryptValueWithPGPKey_SameProfileAsSDKv1 compares live against terraform-plugin-sdk v1's helper/encryption.
// It is removed together with SDK v1; TestEncryptValueWithPGPKey_MatchesSDKv1Profile keeps the recorded result.
func TestEncryptValueWithPGPKey_SameProfileAsSDKv1(t *testing.T) {
	const secret = "6b4f9a2c-1d3e-4f5a-8b7c-9d0e1f2a3b4c"

	for _, k := range pgpTestKeys {
		t.Run(k.name, func(t *testing.T) {
			entity, pgpKey := newTestPGPKey(t, k.cfg)

			_, oldEncrypted, err := encryption.EncryptValue(pgpKey, secret, "ACL Token secret ID")
			if err != nil {
				t.Fatalf("SDK v1 EncryptValue() error: %v", err)
			}
			oldRaw, err := base64.StdEncoding.DecodeString(oldEncrypted)
			if err != nil {
				t.Fatal(err)
			}
			newRaw := encryptForTest(t, pgpKey)

			oldProfile := describePGPMessage(t, oldRaw, entity)
			newProfile := describePGPMessage(t, newRaw, entity)
			t.Logf("SDK v1 profile: %+v", oldProfile)
			if newProfile != oldProfile {
				t.Errorf("message profile differs from SDK v1:\n new: %+v\n old: %+v", newProfile, oldProfile)
			}

			for name, raw := range map[string][]byte{"SDK v1": oldRaw, "new": newRaw} {
				md, err := openpgp.ReadMessage(bytes.NewReader(raw), openpgp.EntityList{entity}, nil, nil)
				if err != nil {
					t.Fatalf("%s: decrypting: %v", name, err)
				}
				got, err := io.ReadAll(md.UnverifiedBody)
				if err != nil || string(got) != secret {
					t.Fatalf("%s: decrypted %q (err %v), want %q", name, got, err, secret)
				}
			}
		})
	}
}
