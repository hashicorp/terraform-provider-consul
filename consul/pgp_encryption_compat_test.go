// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"bytes"
	"encoding/base64"
	"io"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

func encryptForTest(t *testing.T, pgpKey string) []byte {
	t.Helper()
	encrypted, err := encryptValueWithPGPKey(pgpKey, "6b4f9a2c-1d3e-4f5a-8b7c-9d0e1f2a3b4c", "ACL Token secret ID")
	if err != nil {
		t.Fatalf("encryptValueWithPGPKey() error: %v", err)
	}
	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		t.Fatalf("encrypted value is not base64: %v", err)
	}
	return raw
}

// pgpMessageProfile is everything about an encrypted message that affects whether a recipient can decrypt it.
type pgpMessageProfile struct {
	SessionKeyVersion int
	PubKeyAlgo        packet.PublicKeyAlgorithm
	Cipher            packet.CipherFunction
	SEIPDVersion      int
	MDC               bool
	Compressed        bool
	LiteralBinary     bool
	LiteralFileName   string
}

// pgpTestKeys are the key types checked against the SDK v1 behavior.
var pgpTestKeys = []struct {
	name string
	cfg  *packet.Config
}{
	{"rsa-2048", &packet.Config{Algorithm: packet.PubKeyAlgoRSA, RSABits: 2048}},
	{"rsa-1024", &packet.Config{Algorithm: packet.PubKeyAlgoRSA, RSABits: 1024}},
	{"ed25519-cv25519", &packet.Config{Algorithm: packet.PubKeyAlgoEdDSA, Curve: packet.Curve25519}},
	// gpg keys usually list AES-256 first; SDK v1 still picks AES-128 because it checks its own candidates in order.
	{"rsa-2048-prefers-aes256", &packet.Config{Algorithm: packet.PubKeyAlgoRSA, RSABits: 2048, DefaultCipher: packet.CipherAES256}},
	// gpg keys also advertise compression; SDK v1 never compresses.
	{"rsa-2048-prefers-zlib", &packet.Config{Algorithm: packet.PubKeyAlgoRSA, RSABits: 2048, DefaultCompressionAlgo: packet.CompressionZLIB}},
	// Modern (RFC 9580) keys advertise AEAD; SDK v1 never uses it, and older gpg versions cannot decrypt it.
	{"ed25519-cv25519-supports-aead", &packet.Config{Algorithm: packet.PubKeyAlgoEdDSA, Curve: packet.Curve25519, AEADConfig: &packet.AEADConfig{}}},
}

// describePGPMessage decrypts raw with entity's private keys and returns its profile.
func describePGPMessage(t *testing.T, raw []byte, entity *openpgp.Entity) pgpMessageProfile {
	t.Helper()
	var prof pgpMessageProfile
	r := packet.NewReader(bytes.NewReader(raw))

	p, err := r.Next()
	if err != nil {
		t.Fatalf("reading session key packet: %v", err)
	}
	ek, ok := p.(*packet.EncryptedKey)
	if !ok {
		t.Fatalf("first packet is %T, want *packet.EncryptedKey", p)
	}
	privateKeys := []*packet.PrivateKey{entity.PrivateKey}
	for _, s := range entity.Subkeys {
		privateKeys = append(privateKeys, s.PrivateKey)
	}
	for _, k := range privateKeys {
		if k != nil && k.KeyId == ek.KeyId {
			if err := ek.Decrypt(k, nil); err != nil {
				t.Fatalf("decrypting session key: %v", err)
			}
		}
	}
	if ek.Key == nil {
		t.Fatalf("message is not encrypted to any key of the test entity")
	}
	prof.SessionKeyVersion, prof.PubKeyAlgo, prof.Cipher = ek.Version, ek.Algo, ek.CipherFunc

	p, err = r.Next()
	if err != nil {
		t.Fatalf("reading encrypted data packet: %v", err)
	}
	se, ok := p.(*packet.SymmetricallyEncrypted)
	if !ok {
		t.Fatalf("second packet is %T, want *packet.SymmetricallyEncrypted", p)
	}
	prof.SEIPDVersion, prof.MDC = se.Version, se.IntegrityProtected

	rc, err := se.Decrypt(ek.CipherFunc, ek.Key)
	if err != nil {
		t.Fatalf("decrypting data packet: %v", err)
	}
	inner := packet.NewReader(rc)
	for {
		p, err := inner.Next()
		if err != nil {
			break
		}
		switch v := p.(type) {
		case *packet.Compressed:
			prof.Compressed = true
			inner.Push(v.Body)
		case *packet.LiteralData:
			prof.LiteralBinary, prof.LiteralFileName = v.IsBinary, v.FileName
			_, _ = io.ReadAll(v.Body)
		}
	}
	return prof
}

// TestEncryptValueWithPGPKey_MatchesSDKv1Profile pins the message format produced by terraform-plugin-sdk v1
// (helper/encryption, keybase/go-crypto). The expected values were recorded from SDK v1 with
// TestEncryptValueWithPGPKey_SameProfileAsSDKv1, so this keeps guarding compatibility after SDK v1 is removed.
func TestEncryptValueWithPGPKey_MatchesSDKv1Profile(t *testing.T) {
	sdkv1 := func(algo packet.PublicKeyAlgorithm, cipher packet.CipherFunction) pgpMessageProfile {
		return pgpMessageProfile{SessionKeyVersion: 3, PubKeyAlgo: algo, Cipher: cipher, SEIPDVersion: 1, MDC: true}
	}
	want := map[string]pgpMessageProfile{
		"rsa-2048":                      sdkv1(packet.PubKeyAlgoRSA, packet.CipherAES128),
		"rsa-1024":                      sdkv1(packet.PubKeyAlgoRSA, packet.CipherAES128),
		"ed25519-cv25519":               sdkv1(packet.PubKeyAlgoECDH, packet.CipherAES128),
		"rsa-2048-prefers-aes256":       sdkv1(packet.PubKeyAlgoRSA, packet.CipherAES128),
		"rsa-2048-prefers-zlib":         sdkv1(packet.PubKeyAlgoRSA, packet.CipherAES128),
		"ed25519-cv25519-supports-aead": sdkv1(packet.PubKeyAlgoECDH, packet.CipherAES128),
	}

	for _, k := range pgpTestKeys {
		t.Run(k.name, func(t *testing.T) {
			entity, pgpKey := newTestPGPKey(t, k.cfg)
			raw := encryptForTest(t, pgpKey)
			if got := describePGPMessage(t, raw, entity); got != want[k.name] {
				t.Errorf("message profile = %+v, want (SDK v1) %+v", got, want[k.name])
			}
		})
	}
}
