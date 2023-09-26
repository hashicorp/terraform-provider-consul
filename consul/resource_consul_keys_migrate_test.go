// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestConsulKeysMigrateState(t *testing.T) {
	startTestServer(t)

	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		Expected     map[string]string
		Meta         interface{}
	}{
		"v0.6.9 and earlier, with old values hash function": {
			StateVersion: 0,
			Attributes: map[string]string{
				"key.#":             "2",
				"key.12345.name":    "hello",
				"key.12345.path":    "foo/bar",
				"key.12345.value":   "world",
				"key.12345.default": "",
				"key.12345.delete":  "false",
				"key.6789.name":     "temp",
				"key.6789.path":     "foo/foo",
				"key.6789.value":    "",
				"key.6789.default":  "",
				"key.6789.delete":   "true",
			},
			Expected: map[string]string{
				"key.#":                  "2",
				"key.1529757638.default": "",
				"key.1529757638.delete":  "true",
				"key.3609964659.name":    "hello",
				"key.1529757638.path":    "foo/foo",
				"key.1529757638.value":   "",
				"key.1529757638.flags":   "0",
				"key.1529757638.cas":     "0",
				"key.3609964659.path":    "foo/bar",
				"key.3609964659.default": "",
				"key.3609964659.delete":  "false",
				"key.1529757638.name":    "temp",
				"key.3609964659.value":   "world",
				"key.3609964659.cas":     "0",
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "consul",
			Attributes: tc.Attributes,
		}
		is, err := resourceConsulKeys().MigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.Expected {
			if is.Attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, v, k, is.Attributes[k], is.Attributes)
			}
		}
	}
}

func TestConsulKeysMigrateState_empty(t *testing.T) {
	var is *terraform.InstanceState
	var meta interface{}

	// should handle nil
	is, err := resourceConsulKeys().MigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
	if is != nil {
		t.Fatalf("expected nil instancestate, got: %#v", is)
	}

	// should handle non-nil but empty
	is = &terraform.InstanceState{}

	if _, err = resourceConsulKeys().MigrateState(0, is, meta); err != nil {
		t.Fatalf("err: %#v", err)
	}
}
