// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

const keybasePrefix = "keybase:"

// keybaseLookupURL is a variable so tests can point it at a fake server.
var keybaseLookupURL = "https://keybase.io/_/api/1.0/user/lookup.json"

// encryptValueWithPGPKey encrypts value with pgpKey (a base64-encoded public key or "keybase:<username>")
// and returns the base64-encoded encrypted message.
func encryptValueWithPGPKey(pgpKey, value, description string) (string, error) {
	encryptionKey := pgpKey
	if strings.HasPrefix(pgpKey, keybasePrefix) {
		key, err := fetchKeybasePublicKey(strings.TrimPrefix(pgpKey, keybasePrefix))
		if err != nil {
			return "", fmt.Errorf("Error retrieving Public Key for %s: %w", pgpKey, err)
		}
		encryptionKey = key
	}

	encrypted, err := encryptWithPGPKey(encryptionKey, []byte(value))
	if err != nil {
		return "", fmt.Errorf("Error encrypting %s: %w", description, err)
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// encryptWithPGPKey encrypts value for the base64-encoded binary public key and returns a binary OpenPGP message.
func encryptWithPGPKey(base64Key string, value []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("error decoding given PGP key: %w", err)
	}
	entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewReader(data)))
	if err != nil {
		return nil, fmt.Errorf("error parsing given PGP key: %w", err)
	}

	var ciphertext bytes.Buffer
	// A nil config keeps the SEIPDv1 (non-AEAD), uncompressed output older gpg versions can decrypt.
	plaintext, err := openpgp.Encrypt(&ciphertext, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error setting up encryption for PGP message: %w", err)
	}
	if _, err := plaintext.Write(value); err != nil {
		return nil, fmt.Errorf("error encrypting PGP message: %w", err)
	}
	if err := plaintext.Close(); err != nil {
		return nil, fmt.Errorf("error encrypting PGP message: %w", err)
	}
	return ciphertext.Bytes(), nil
}

// fetchKeybasePublicKey returns the base64-encoded binary primary public key of a keybase user.
func fetchKeybasePublicKey(username string) (string, error) {
	query := url.Values{"usernames": {username}, "fields": {"public_keys"}}
	resp, err := cleanhttp.DefaultClient().Get(keybaseLookupURL + "?" + query.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out struct {
		Status struct {
			Name string
		}
		Them []struct {
			PublicKeys struct {
				Primary struct {
					Bundle string
				}
			} `json:"public_keys"`
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.Status.Name != "OK" {
		return "", fmt.Errorf("got non-OK response: %q", out.Status.Name)
	}
	if len(out.Them) == 0 || out.Them[0].PublicKeys.Primary.Bundle == "" {
		return "", fmt.Errorf("unable to fetch keys for user(s) %q from keybase", username)
	}

	entities, err := openpgp.ReadArmoredKeyRing(strings.NewReader(out.Them[0].PublicKeys.Primary.Bundle))
	if err != nil {
		return "", err
	}
	if len(entities) != 1 || entities[0] == nil {
		return "", fmt.Errorf("primary key could not be parsed for user %q", username)
	}

	var serialized bytes.Buffer
	if err := entities[0].Serialize(&serialized); err != nil {
		return "", fmt.Errorf("error serializing entity for user %q: %w", username, err)
	}
	return base64.StdEncoding.EncodeToString(serialized.Bytes()), nil
}
