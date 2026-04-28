// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"os"
	"strings"
)

// GetStringParam safely extracts and trims a string parameter from params map
// Returns empty string if key doesn't exist, value is nil, or trimmed value is empty
func GetStringParam(params map[string]interface{}, key string) string {
	if params == nil {
		return ""
	}

	val, exists := params[key]
	if !exists {
		return ""
	}

	// Handle nil value
	if val == nil {
		return ""
	}

	// Type assert to string
	strVal, ok := val.(string)
	if !ok {
		return ""
	}

	// Trim whitespace and return
	return strings.TrimSpace(strVal)
}

// GetStringEnv safely retrieves and trims an environment variable
// Falls back to fallback values if primary is empty
func GetStringEnv(primary, fallback string) string {
	if val := strings.TrimSpace(os.Getenv(primary)); val != "" {
		return val
	}
	if fallback != "" {
		return strings.TrimSpace(os.Getenv(fallback))
	}
	return ""
}

// IsNonEmpty checks if a string is not empty after trimming
func IsNonEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}
