// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"bytes"
	"fmt"
	"hash/crc32"
)

// hashcodeString is terraform-plugin-sdk v1's helper/hashcode.String, kept so set hashes don't change.
func hashcodeString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// hashcodeStrings is terraform-plugin-sdk v1's helper/hashcode.Strings.
func hashcodeStrings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", hashcodeString(buf.String()))
}
