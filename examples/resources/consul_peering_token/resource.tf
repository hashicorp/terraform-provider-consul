# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "consul_peering_token" "token" {
  peer_name = "eu-cluster"
}
