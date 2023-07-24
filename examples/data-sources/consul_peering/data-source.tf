# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "consul_peering" "basic" {
  peer_name = "peered-cluster"
}
