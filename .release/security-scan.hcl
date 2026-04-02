# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Reference: https://github.com/hashicorp/security-scanner/blob/main/CONFIG.md#binary (private repository)

binary {
  secrets {
    all = true
  }
  go_modules   = true
  osv          = true
  oss_index    = false
  nvd          = false

  # Triage items that are _safe_ to ignore here. Note that this list should be
	# periodically cleaned up to remove items that are no longer found by the scanner.

  triage {
		suppress {
			vulnerabilities = [
                "GO-2022-0635"
			]
		}
	}
}
