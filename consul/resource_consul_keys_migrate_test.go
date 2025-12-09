// Copyright IBM Corp. 2014, 2025
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
				"key.3630941097.default": "",
				"key.3630941097.delete":  "true",
				"key.3630941097.name":    "temp",
				"key.3630941097.path":    "foo/foo",
				"key.3630941097.value":   "",
				"key.3975462262.path":    "foo/bar",
				"key.3975462262.default": "",
				"key.3975462262.delete":  "false",
				"key.3975462262.name":    "hello",
				"key.3975462262.value":   "world",
				"key.3975462262.flags":   "0",
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
